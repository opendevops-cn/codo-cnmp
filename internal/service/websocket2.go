package service

import (
	"context"
	"fmt"
	"io"
	"sync"

	"codo-cnmp/internal/biz"
	"codo-cnmp/internal/middleware"
	"codo-cnmp/pb"
	"github.com/go-kratos/kratos/v2/log"
	websocket2 "github.com/gorilla/websocket"
	"github.com/opendevops-cn/codo-golang-sdk/adapter/kratos/transport/websocket"
	"golang.org/x/sync/errgroup"
	"k8s.io/client-go/tools/remotecommand"
)

type PodTerminalWebsocketService struct {
	pod        *biz.PodUseCase
	log        *log.Helper
	middleware *middleware.CasbinCheckMiddleware
}

func NewPodTerminalWebsocketService(logger log.Logger, pod *biz.PodUseCase, middleware *middleware.CasbinCheckMiddleware) *PodTerminalWebsocketService {
	l := log.NewHelper(log.With(logger, "module", "service/websocket"))
	return &PodTerminalWebsocketService{
		log:        l,
		pod:        pod,
		middleware: middleware,
	}
}

type TerminalSession struct {
	sizeChan chan remotecommand.TerminalSize
	msgChan  chan *pb.PodCommandResponse
	cmdChan  chan string
	ctx      context.Context
	cancel   context.CancelFunc
	log      *log.Helper
}

func NewTerminalSession(ctx context.Context, msgChan chan *pb.PodCommandResponse, logger *log.Helper) *TerminalSession {
	ctx, cancel := context.WithCancel(ctx)
	return &TerminalSession{
		sizeChan: make(chan remotecommand.TerminalSize, 1),
		msgChan:  msgChan,
		cmdChan:  make(chan string, 100),
		ctx:      ctx,
		cancel:   cancel,
		log:      logger,
	}
}

func (t *TerminalSession) Read(p []byte) (n int, err error) {
	select {
	case cmd := <-t.cmdChan:
		return copy(p, cmd), nil
	case <-t.ctx.Done():
		// session的context取消时，发送退出命令
		t.log.Info("会话终端开始退出，发送退出命令")
		return copy(p, "\u0004"), io.EOF
	}
}

func (t *TerminalSession) Write(p []byte) (n int, err error) {
	select {
	case t.msgChan <- &pb.PodCommandResponse{Output: string(p)}:
		return len(p), nil
	case <-t.ctx.Done():
		return 0, io.EOF
	}
}

func (t *TerminalSession) Next() *remotecommand.TerminalSize {
	select {
	case size := <-t.sizeChan:
		return &size
	case <-t.ctx.Done():
		return nil
	}
}

func (t *TerminalSession) Close() error {
	t.log.Info("开始关闭终端会话")
	t.cancel() // 取消session的context，触发退出命令发送
	close(t.cmdChan)
	close(t.sizeChan)
	return nil
}

type TerminalReplier struct {
	session  *TerminalSession
	pod      *biz.PodUseCase
	log      *log.Helper
	executor remotecommand.Executor
	msgChan  chan *pb.PodCommandResponse

	closeOnce sync.Once
	ctx       context.Context
	cancel    func()
}

func NewTerminalReplier(ctx context.Context, pod *biz.PodUseCase, logger *log.Helper) *TerminalReplier {
	ctx, cancel := context.WithCancel(ctx)
	return &TerminalReplier{
		pod:     pod,
		log:     logger,
		msgChan: make(chan *pb.PodCommandResponse, 100),
		ctx:     ctx,
		cancel:  cancel,
	}
}

func (r *TerminalReplier) handleConnect(ctx context.Context, info *pb.ExecPodCommandRequest_ConnectInfo) error {
	executor, err := r.pod.ExecPod(ctx, &biz.ExecPodRequest{
		PodCommonParams: biz.PodCommonParams{
			ClusterName: info.ClusterName,
			Namespace:   info.Namespace,
			PodName:     info.PodName,
		},
		ContainerName: info.ContainerName,
		Shell:         info.Shell,
	})
	if err != nil {
		return fmt.Errorf("创建终端会话失败: %w", err)
	}

	r.executor = executor
	r.session = NewTerminalSession(ctx, r.msgChan, r.log)

	eg, ctx := errgroup.WithContext(r.session.ctx)

	eg.Go(func() error {
		err := executor.StreamWithContext(ctx, remotecommand.StreamOptions{
			Stdin:             r.session,
			Stdout:            r.session,
			Stderr:            r.session,
			Tty:               true,
			TerminalSizeQueue: r.session,
		})
		if err != nil {
			r.log.WithContext(ctx).Errorf("终端流异常: %v", err)
			return err
		}
		return nil
	})

	go func(ctx context.Context) {
		if err := eg.Wait(); err != nil {
			r.log.WithContext(ctx).Errorf("终端会话错误: %v", err)
		}
	}(ctx)

	return nil
}

func (r *TerminalReplier) handleCommand(command string) error {
	if r.session == nil {
		return fmt.Errorf("会话终端未建立")
	}
	select {
	case r.session.cmdChan <- command:
		return nil
	default:
		return fmt.Errorf("会话终端异常")
	}
}

func (r *TerminalReplier) handleResize(info *pb.ExecPodCommandRequest_ResizeWindowInfo) error {
	if r.session == nil {
		return fmt.Errorf("会话终端未建立")
	}
	select {
	case r.session.sizeChan <- remotecommand.TerminalSize{
		Width:  uint16(info.Cols),
		Height: uint16(info.Rows),
	}:
		return nil
	default:
		return fmt.Errorf("会话终端大小调整异常")
	}
}

func (r *TerminalReplier) Apply(ctx context.Context, req *pb.ExecPodCommandRequest) error {
	switch req.OperationType {
	case pb.ExecPodCommandRequest_Connect:
		return r.handleConnect(ctx, req.ConnectInfo)
	case pb.ExecPodCommandRequest_Command:
		return r.handleCommand(req.CommandInfo.Command)
	case pb.ExecPodCommandRequest_ResizeWindow:
		return r.handleResize(req.ResizeWindowInfo)
	case pb.ExecPodCommandRequest_Disconnect:
		return r.Close(ctx)
	default:
		return fmt.Errorf("未知操作类型")
	}
}

func (r *TerminalReplier) Reply(ctx context.Context) (*pb.PodCommandResponse, error) {
	select {
	case msg, ok := <-r.msgChan:
		if !ok {
			return nil, fmt.Errorf("消息通道已关闭")
		}
		return msg, nil
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}
func (r *TerminalReplier) Close(ctx context.Context) error {
	r.log.WithContext(ctx).Info("开始关闭终端会话replier")

	if r.session != nil {
		if err := r.session.Close(); err != nil {
			r.log.WithContext(ctx).Errorf("关闭终端出错: %v", err)
		}
	}

	if r.cancel != nil {
		r.cancel()
	}

	close(r.msgChan)
	r.session = nil
	r.executor = nil
	r.log.WithContext(ctx).Info("终端会话replier已关闭")
	return nil
}

func (x *PodTerminalWebsocketService) Build(ctx context.Context) (websocket.Handler, error) {
	return websocket.NewWebSocket[*pb.PodCommandResponse, *pb.ExecPodCommandRequest](
		NewTerminalReplier(ctx, x.pod, x.log),
		websocket.WithWSMiddlewareFunc(x.middleware.WSServer()),
		websocket.WithReplyErrorEncodeFunc(func(conn *websocket2.Conn, err error) {
			WebSocketErrorEncodeFunc(conn, err)
		}),
	), nil
}

func (x *PodTerminalWebsocketService) Path() string {
	return "/api/v1/ws/pod/command"
}
