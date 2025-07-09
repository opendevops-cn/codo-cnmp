package service

import (
	"bytes"
	"codo-cnmp/common/consts"
	"codo-cnmp/common/utils"
	"codo-cnmp/common/xerr"
	"codo-cnmp/internal/dep"
	"codo-cnmp/internal/middleware"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	websocket2 "github.com/gorilla/websocket"
	e "github.com/pkg/errors"
	pkgerr "github.com/pkg/errors"
	"io"
	"net/http"
	"time"

	"codo-cnmp/internal/biz"
	"codo-cnmp/pb"
	"github.com/ccheers/xpkg/sync/errgroup"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/opendevops-cn/codo-golang-sdk/adapter/kratos/transport/websocket"
	"k8s.io/client-go/tools/remotecommand"
)

type PodLogWebsocketService struct {
	log *log.Helper
	pod *biz.PodUseCase

	middleware *middleware.CasbinCheckMiddleware
}

func (x *PodLogWebsocketService) Build(ctx context.Context) (websocket.Handler, error) {
	return websocket.NewWebSocket[*pb.TailPodLogsResponse, *pb.TailPodLogsRequest](
		NewPodLogReplier(x.pod, x.log),
		websocket.WithWSMiddlewareFunc(x.middleware.WSServer()),
	), nil
}

func (x *PodLogWebsocketService) Path() string {
	return "/api/v1/ws/pod/log"
}

func NewWebsocketService(logger log.Logger, pod *biz.PodUseCase, middleware *middleware.CasbinCheckMiddleware) *PodLogWebsocketService {
	l := log.NewHelper(log.With(logger, "module", "service/websocket"))
	return &PodLogWebsocketService{
		log:        l,
		pod:        pod,
		middleware: middleware,
	}
}

type PodLogReplier struct {
	ch  chan string
	pod *biz.PodUseCase
	log *log.Helper
}

func NewPodLogReplier(pod *biz.PodUseCase, logger *log.Helper) *PodLogReplier {
	return &PodLogReplier{pod: pod, ch: make(chan string, 32), log: logger}
}

func (x *PodLogReplier) Apply(ctx context.Context, req *pb.TailPodLogsRequest) error {
	x.log.WithContext(ctx).Infof("用户连接到 Pod 日志: %v", req)
	return x.pod.TailPodLogs(ctx, x.ch, &biz.TailPodLogRequest{
		PodCommonParams: biz.PodCommonParams{
			ClusterName: req.ClusterName,
			Namespace:   req.Namespace,
			PodName:     req.PodName,
		},
		ContainerName: req.ContainerName,
		TailLines:     req.TailLines,
	})
}

func (x *PodLogReplier) Reply(ctx context.Context) (*pb.TailPodLogsResponse, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case line := <-x.ch:
		return &pb.TailPodLogsResponse{Log: line}, nil
	}
}

func (x *PodLogReplier) Close(ctx context.Context) error {
	close(x.ch)
	return nil
}

type PodCommandWebsocketService struct {
	log        *log.Helper
	pod        *biz.PodUseCase
	auditLog   *biz.AuditLogUseCase
	middleware *middleware.CasbinCheckMiddleware
	kafka      dep.IKafka
}

func WebSocketErrorEncodeFunc(conn *websocket2.Conn, err error) {
	// 错误返回
	newErr := errors.Unwrap(err)
	if newErr != nil {
		err = newErr
	}
	errCode := xerr.ServerCommonError
	causeErr := pkgerr.Cause(err) // err类型
	var e *xerr.CodeError
	if errors.As(causeErr, &e) { // 自定义错误类型
		// 自定义CodeError
		errCode = e.GetErrCode()
	}
	if errCode == xerr.ErrNotAllowed {
		errCode = 403
	}

	_ = conn.WriteJSON(map[string]interface{}{
		"code": errCode,
		"msg":  err.Error(),
	})
}

func (x *PodCommandWebsocketService) Build(ctx context.Context) (websocket.Handler, error) {
	return websocket.NewWebSocket[*pb.PodCommandResponse, *pb.ExecPodCommandRequest](
		NewPodCommandReplier(ctx, x.pod, x.auditLog, x.log, x.kafka),
		websocket.WithWSMiddlewareFunc(x.middleware.WSServer(), func(handleFunc websocket.WSPreHandleFunc) websocket.WSPreHandleFunc {
			return func(ctx context.Context, request *http.Request) error {
				err := handleFunc(ctx, request)
				return err
			}
		}),
		websocket.WithReplyErrorEncodeFunc(func(conn *websocket2.Conn, err error) {
			WebSocketErrorEncodeFunc(conn, err)
		}),
	), nil
}

func (x *PodCommandWebsocketService) Path() string {
	return "/api/v1/ws/pod/command"
}

func NewPodCommandWebsocketService(ctx context.Context, logger log.Logger, pod *biz.PodUseCase, auditLog *biz.AuditLogUseCase, middleware *middleware.CasbinCheckMiddleware, kafka dep.IKafka) *PodCommandWebsocketService {
	l := log.NewHelper(log.With(logger, "module", "service/websocket"))
	return &PodCommandWebsocketService{
		log:        l,
		pod:        pod,
		auditLog:   auditLog,
		middleware: middleware,
		kafka:      kafka,
	}
}

type PodCommandReplier struct {
	executor remotecommand.Executor
	pod      *biz.PodUseCase
	auditLog *biz.AuditLogUseCase
	log      *log.Helper

	cancel func()
	eg     *errgroup.Group

	//stdin *os.File
	//stdout *os.File
	stdin             io.WriteCloser
	stdout            io.ReadCloser
	stderr            bytes.Buffer
	terminalSizeQueue remotecommand.TerminalSizeQueue
	resizeChan        chan remotecommand.TerminalSize
	pingChan          chan struct{}
	kafka             dep.IKafka

	packChan    chan *pb.PodCommandResponse
	packErrChan chan error

	pipeClose func()
}

func NewPodCommandReplier(ctx context.Context, pod *biz.PodUseCase, auditLog *biz.AuditLogUseCase, logger *log.Helper, kafka dep.IKafka) *PodCommandReplier {
	resizeChan := make(chan remotecommand.TerminalSize, 1)
	return &PodCommandReplier{
		executor:          nil,
		pod:               pod,
		auditLog:          auditLog,
		log:               logger,
		cancel:            nil,
		eg:                nil,
		stdin:             nil,
		stdout:            nil,
		stderr:            bytes.Buffer{},
		terminalSizeQueue: NewTerminalSizeQueue(ctx, resizeChan),
		resizeChan:        resizeChan,
		kafka:             kafka,
		packChan:          make(chan *pb.PodCommandResponse, 1024),
		packErrChan:       make(chan error, 128),
		pipeClose:         nil,
	}
}

type TerminalSizeQueue struct {
	ctx        context.Context
	resizeChan chan remotecommand.TerminalSize
}

func NewTerminalSizeQueue(ctx context.Context, resizeChan chan remotecommand.TerminalSize) *TerminalSizeQueue {
	return &TerminalSizeQueue{
		resizeChan: resizeChan,
		ctx:        ctx,
	}
}

func (t *TerminalSizeQueue) Next() *remotecommand.TerminalSize {
	select {
	case size := <-t.resizeChan:
		return &size
	case <-t.ctx.Done():
		return nil
	}
}

// PodExecInfo 包含需要审计的容器信息
type PodExecInfo struct {
	ClusterName   string
	Namespace     string
	PodName       string
	ContainerName string
	UserName      string
	Shell         string
	ClientIP      string
}

type auditBuffer struct {
	src       io.WriteCloser
	buffer    bytes.Buffer
	ctx       context.Context
	auditUc   *biz.AuditLogUseCase
	podInfo   *PodExecInfo
	startTime time.Time
	log       *log.Helper
	traceId   string
	kafka     dep.IKafka
}

func (x *auditBuffer) Close() error {
	return x.src.Close()
}

func (x *auditBuffer) Write(p []byte) (n int, err error) {
	// 检查是否需要记录审计日志（当遇到换行符时）
	go x.recordAudit(p)
	n, err = x.src.Write(p)
	// 同时写入到 buffer
	return n, err
}

// recordAudit 记录审计日志
func (x *auditBuffer) recordAudit(p []byte) {
	x.buffer.Write(p)
	delemiters := []byte{'\n', '\r'}
	delimiter := byte('0')
	for _, d := range delemiters {
		if bytes.Contains(p, []byte{d}) {
			delimiter = d
			break
		}
	}
	if delimiter == '0' {
		return
	}
	command, err := x.buffer.ReadString(delimiter)
	if err != nil {
		return
	}
	// 跳过特殊字符（如 Ctrl+D）
	if command == "\u0004\n" {
		return
	}
	if command == "" {
		return
	}
	requestBody := struct {
		Container string `json:"container"`
		Shell     string `json:"shell"`
		Command   string `json:"command"`
	}{
		Container: x.podInfo.ContainerName,
		Shell:     x.podInfo.Shell,
		Command:   command,
	}

	// 转换为 JSON
	requestBodyJSON, err := json.Marshal(requestBody)
	if err != nil {
		x.log.Errorf("序列化请求体失败: %v", err)
		return
	}
	auditLog := &biz.CreateAuditLogRequest{
		AuditLogItem: biz.AuditLogItem{
			UserName:      x.podInfo.UserName,
			ClientIP:      x.podInfo.ClientIP,
			Module:        "pod",
			Action:        "控制台",
			Cluster:       x.podInfo.ClusterName,
			Namespace:     x.podInfo.Namespace,
			ResourceType:  "pod",
			ResourceName:  x.podInfo.PodName,
			RequestBody:   string(requestBodyJSON),
			RequestPath:   "/api/v1/ws/pod/command",
			Status:        int(pb.OperationStatus_Success),
			Duration:      "0",
			OperationTime: x.startTime.Format(time.DateTime),
			TraceID:       x.traceId,
		},
	}

	if err := x.auditUc.CreateAuditLog(x.ctx, auditLog); err != nil {
		x.log.Errorf("记录容器审计日志失败: %v", err)
	}
	// 发送审计日志到 kafka
	auditLogBytes, err := json.Marshal(auditLog)
	if err != nil {
		x.log.Errorf("序列化审计日志失败: %v", err)
	}
	// 发送消息到 Kafka
	if x.kafka != nil {
		go func(ctx context.Context, auditLogBytes []byte) {
			defer func() {
				if r := recover(); r != nil {
					x.log.Errorf("发送消息到 Kafka 时发生 panic: %v", r)
					x.kafka.Close(ctx)
				}
			}()
			err := x.kafka.SendMessage(ctx, auditLogBytes)
			if err != nil {
				x.log.Errorf("发送消息到 Kafka 失败: %v", err)
			} else {
				x.log.Info("发送消息到 Kafka 成功")
			}
		}(x.ctx, auditLogBytes)
	}
}

// 创建新的 auditBuffer
func newAuditBuffer(src io.WriteCloser, podInfo *PodExecInfo, auditUc *biz.AuditLogUseCase, traceId string, logger *log.Helper, kafka dep.IKafka) *auditBuffer {
	return &auditBuffer{
		src:       src,
		auditUc:   auditUc,
		podInfo:   podInfo,
		traceId:   traceId,
		startTime: time.Now(),
		log:       logger,
		kafka:     kafka,
	}
}

func (x *PodCommandReplier) Apply(ctx context.Context, req *pb.ExecPodCommandRequest) error {
	switch req.OperationType {
	case pb.ExecPodCommandRequest_Connect:
		if x.executor != nil {
			return fmt.Errorf("已打开终端")
		}
		ctx, x.cancel = context.WithCancel(ctx)
		x.eg = errgroup.WithCancel(ctx)
		connInfo := req.ConnectInfo
		executor, err := x.pod.ExecPod(ctx, &biz.ExecPodRequest{
			PodCommonParams: biz.PodCommonParams{
				ClusterName: connInfo.ClusterName,
				Namespace:   connInfo.Namespace,
				PodName:     connInfo.PodName,
			},
			ContainerName: connInfo.ContainerName,
			Shell:         connInfo.Shell,
		})
		if err != nil {
			x.log.WithContext(ctx).Errorf("进入pod终端失败: %v", err)
			defer x.Close(ctx)
			return e.Wrapf(xerr.NewErrCodeMsg(xerr.ErrWebsocketNotCONNECT, "打开连接失败"), "进入终端失败")
		}
		x.executor = executor
		stdinReader, stdinWriter := io.Pipe()
		stdoutReader, stdoutWriter := io.Pipe()
		userName, err := utils.GetUserNameFromCtx(ctx)
		if err != nil {
			userName = ""
		}
		// 从上下文获取客户端 IP
		clientIP := ""
		if ipValue, ok := ctx.Value(consts.ContextClientIPKey).(string); ok {
			clientIP = ipValue
		}
		podInfo := &PodExecInfo{
			ClusterName:   connInfo.ClusterName,
			Namespace:     connInfo.Namespace,
			PodName:       connInfo.PodName,
			ContainerName: connInfo.ContainerName,
			UserName:      userName,
			ClientIP:      clientIP,
		}
		traceID := ""
		if traceId, ok := ctx.Value(consts.ContextTraceIDKey).(string); ok {
			traceID = traceId
		}
		x.stdin = newAuditBuffer(stdinWriter, podInfo, x.auditLog, traceID, x.log, x.kafka)
		x.stdout = stdoutReader
		x.pipeClose = func() {
			if err := stdinReader.Close(); err != nil {
				fmt.Printf("stdinReader.Close() error: %v\n", err)
			}
			if err := stdinWriter.Close(); err != nil {
				fmt.Printf("stdinWriter.Close() error: %v\n", err)
			}
			if err := stdoutReader.Close(); err != nil {
				fmt.Printf("stdoutReader.Close() error: %v\n", err)
			}
			if err := stdoutWriter.Close(); err != nil {
				fmt.Printf("stdoutWriter.Close() error: %v\n", err)
			}
		}

		options := remotecommand.StreamOptions{
			Stdin:             stdinReader,
			Stdout:            stdoutWriter,
			TerminalSizeQueue: x.terminalSizeQueue,
			Tty:               true,
		}
		connectErr := make(chan error, 1)
		x.eg.Go(func(ctx context.Context) error {
			err := x.executor.StreamWithContext(ctx, options)
			if err != nil {
				connectErr <- err
			}
			return err
		})
		x.eg.Go(func(ctx context.Context) error {
			x.stdoutLoop(ctx)
			return nil
		})
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(time.Second * 1):
			return nil
		case err := <-connectErr:
			return err
		}
		// 监听 Ping 信号
		//x.pingChan = make(chan struct{}, 1)
		//x.eg.Go(func(ctx context.Context) error {
		//	return x.Ping(ctx)
		//})
	case pb.ExecPodCommandRequest_Command:
		if x.executor == nil {
			return fmt.Errorf("未建立连接")
		}
		_, err := x.stdin.Write([]byte(req.CommandInfo.Command))
		if err != nil {
			return err
		}
	case pb.ExecPodCommandRequest_ResizeWindow:
		if x.executor == nil {
			return fmt.Errorf("未建立连接")
		}
		size := remotecommand.TerminalSize{
			Width:  uint16(req.ResizeWindowInfo.Cols),
			Height: uint16(req.ResizeWindowInfo.Rows),
		}
		select {
		case x.resizeChan <- size:
		default:
			return fmt.Errorf("终端大小调整失败")
		}
		return nil
	case pb.ExecPodCommandRequest_Disconnect:
		if x.executor == nil {
			return fmt.Errorf("未建立连接")
		}
		_, err := x.stdin.Write([]byte("\u0004\n"))
		if err != nil {
			return err
		}
		x.executor = nil
	case pb.ExecPodCommandRequest_Ping:
		if x.executor == nil {
			return fmt.Errorf("未建立连接")
		}
		// 通过 Ping 来维持心跳
		select {
		case x.packChan <- &pb.PodCommandResponse{Output: ""}:
		default:
			//return x.Close(ctx)
		}
		return nil
	default:
		return fmt.Errorf("未知操作类型")
	}
	return nil
}

func (x *PodCommandReplier) stdoutLoop(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			x.stdoutRead(ctx)
		}
	}
}
func (x *PodCommandReplier) stdoutRead(ctx context.Context) {
	defer func() {
		recover()
	}()
	buf := make([]byte, 1024)
	n, err := x.stdout.Read(buf)
	if err != nil && !errors.Is(err, io.EOF) {
		x.packErrChan <- err
		return
	}
	if n > 0 {
		// 如果读取到数据，则返回数据并清空缓存
		data := string(buf[:n])
		x.packChan <- &pb.PodCommandResponse{Output: data}
	}
}

func (x *PodCommandReplier) Reply(ctx context.Context) (*pb.PodCommandResponse, error) {
	for {
		select {
		case <-ctx.Done():
			return &pb.PodCommandResponse{Output: ""}, nil
		case pack := <-x.packChan:
			return pack, nil
		case err := <-x.packErrChan:
			return nil, err
		}
	}
}

func (x *PodCommandReplier) Close(ctx context.Context) error {
	if x.stdin != nil {
		// 发送退出命令
		_, err := x.stdin.Write([]byte("\u0004\n"))
		if err != nil {
			return fmt.Errorf("发送退出命令时出错: %w", err)
		}
	}
	if x.pipeClose != nil {
		x.pipeClose() // 仅在 x.pipeClose 非空时调用
		x.pipeClose = nil
	}
	// 释放缓冲区
	close(x.packChan)
	close(x.packErrChan)

	if x.cancel != nil {
		x.cancel() // 仅在 x.cancel 非空时调用
		x.cancel = nil
	}
	if x.eg != nil {
		err := x.eg.Wait() // 仅在 x.eg 非空时调用
		x.eg = nil
		return err
	}
	return nil
}
