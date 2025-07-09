package server

import (
	"bufio"
	"bytes"
	"context"
	"crypto/md5"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"path/filepath"
	"regexp"
	"runtime/debug"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"codo-cnmp/common/utils"
	"github.com/google/uuid"
	"github.com/moby/spdystream/spdy"
	"k8s.io/apimachinery/pkg/util/httpstream"
	k8sproxy "k8s.io/apimachinery/pkg/util/proxy"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd/api"

	"codo-cnmp/common/consts"
	"codo-cnmp/common/k8s"
	"codo-cnmp/internal/biz"
	"codo-cnmp/internal/conf"
	"codo-cnmp/internal/dep"
	"codo-cnmp/internal/middleware"
	"codo-cnmp/pb"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-redis/redis/v8"
	apirequest "k8s.io/apiserver/pkg/endpoints/request"
	"k8s.io/client-go/tools/clientcmd"
)

var ansiEscape = regexp.MustCompile(`\x1b\[[0-9;]*[a-zA-Z]`)

type APIServerProxy struct {
	conf         *conf.Bootstrap
	started      uint32
	listener     net.Listener
	cluster      *biz.ClusterUseCase // 集群业务逻辑
	casbin       *middleware.CasbinCheckMiddleware
	redis        *redis.Client
	audit        *biz.AuditLogUseCase
	kafka        dep.IKafka
	log          *log.Helper
	proxyHandler *APIServerProxyHandler
}

type APIServerProxyHandler struct {
	// 集群业务逻辑
	cluster      *biz.ClusterUseCase
	requestInfo  *apirequest.RequestInfoFactory
	casbin       *middleware.CasbinCheckMiddleware
	redis        *redis.Client
	audit        *biz.AuditLogUseCase
	kafka        dep.IKafka
	log          *log.Helper
	bc           *conf.Bootstrap
	proxyManager *K8sProxyManager
}

func NewAPIServerHandler(cluster *biz.ClusterUseCase, casbin *middleware.CasbinCheckMiddleware, redis *redis.Client,
	audit *biz.AuditLogUseCase, kafka dep.IKafka, logger *log.Helper, conf *conf.Bootstrap,
) *APIServerProxyHandler {
	return &APIServerProxyHandler{
		requestInfo:  k8s.NewRequestInfoResolver(),
		cluster:      cluster,
		casbin:       casbin,
		redis:        redis,
		audit:        audit,
		kafka:        kafka,
		log:          logger,
		bc:           conf,
		proxyManager: NewK8sProxyManager(logger), // 添加代理管理器
	}
}

// RequestInfo 结构体用于存储请求相关信息
type RequestInfo struct {
	UserID       uint32
	UserName     string
	ClientIP     string
	Method       string
	Path         string
	ClusterName  string
	Namespace    string
	Resource     string
	ResourceName string
	StartTime    time.Time
	LogType      string
	ResponseBody string
}

func (x *APIServerProxyHandler) IsSuperUser(ctx context.Context, userId uint32) (bool, string, error) {
	superUserFlag := "0"
	data, err := x.redis.Get(ctx, consts.UserCacheKey).Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		return false, "", err
	}

	if data != "" {
		// 从缓存中获取数据
		var users []*biz.User
		err := json.Unmarshal([]byte(data), &users)
		if err != nil {
			return false, "", err
		}
		for _, user := range users {
			if user.UserID == userId {
				return user.SuperUser == superUserFlag, user.Username, nil
			}
		}
	}
	return false, "", nil
}

// authenticate 从 token 中解析用户信息
func (x *APIServerProxyHandler) authenticate(r *http.Request) (context.Context, uint32, uint32, string, error) {
	strs := strings.Split(r.URL.Path, "/")
	token := strs[3]
	r.URL.Path = "/" + strings.Join(strs[4:], "/")
	r.RequestURI = r.URL.RequestURI()

	secret := x.bc.APP.GetSECRET()
	if secret == "" {
		return nil, 0, 0, "", errors.New("系统未配置签名密钥")
	}

	bs, err := utils.DecodeTokenWithCheckSumEscape64(token, []byte(secret))
	if err != nil {
		return nil, 0, 0, "", fmt.Errorf("用户token无效: %v", err)
	}

	userID := utils.BytesToUInt32(bs[0:4])
	clusterID := utils.BytesToUInt32(bs[4:8])
	isSuperUser, userName, err := x.IsSuperUser(r.Context(), userID)
	if err != nil {
		return nil, 0, 0, "", fmt.Errorf("查询用户失败: %v", err)
	}

	ctx := x.casbin.BuildContextWithClusterRoles(r.Context(), strconv.Itoa(int(userID)))
	ctx = context.WithValue(ctx, consts.ContextSuperUserKey, isSuperUser)
	return ctx, userID, clusterID, userName, nil
}

// authorize 校验请求权限
func (x *APIServerProxyHandler) authorize(r *http.Request, clusterName string) error {
	if err := k8s.AuthorizeRequest(r, clusterName); err != nil {
		return fmt.Errorf("authorization failed: %v", err)
	}
	return nil
}

type K8sProxyManager struct {
	transports sync.Map // clusterID -> http.RoundTripper
	configs    sync.Map // clusterID -> *rest.Config
	mu         sync.RWMutex
	log        *log.Helper
}

func NewK8sProxyManager(log *log.Helper) *K8sProxyManager {
	return &K8sProxyManager{
		log: log,
	}
}

func (m *K8sProxyManager) GetConfig(clusterID, kubeConfig string) (*rest.Config, error) {
	configHash := m.generateConfigHash(kubeConfig)
	cacheKey := fmt.Sprintf("%s_%s", clusterID, configHash)

	if config, ok := m.configs.Load(cacheKey); ok {
		return config.(*rest.Config), nil
	}

	// 创建新的config
	k8sConfig, err := clientcmd.Load([]byte(kubeConfig))
	if err != nil {
		return nil, err
	}

	restConfig, err := clientcmd.BuildConfigFromKubeconfigGetter("", func() (*api.Config,
		error,
	) {
		return k8sConfig, nil
	})
	if err != nil {
		return nil, err
	}

	// 设置Transport复用
	restConfig.Wrap(func(rt http.RoundTripper) http.RoundTripper {
		return newWrapRoundTripper(rt, restConfig, m.log)
	})

	// 缓存config
	m.configs.Store(cacheKey, restConfig)
	return restConfig, nil
}

func (m *K8sProxyManager) generateConfigHash(kubeConfig string) string {
	hash := md5.Sum([]byte(kubeConfig))
	return fmt.Sprintf("%x", hash)
}

func (m *K8sProxyManager) Shutdown() {
	m.configs.Range(func(key, value interface{}) bool {
		config := value.(*rest.Config)

		// 获取Transport并尝试关闭
		if tp, err := rest.TransportFor(config); err == nil {
			// 如果是http.Transport类型，关闭空闲连接
			if httpTransport, ok := tp.(*http.Transport); ok {
				httpTransport.CloseIdleConnections()
			}
			// 如果实现了Closer接口，调用Close
			if closer, ok := tp.(io.Closer); ok {
				closer.Close()
			}
		}

		m.configs.Delete(key)
		return true
	})

	// 清理transports缓存
	m.transports.Range(func(key, value interface{}) bool {
		if tp, ok := value.(http.RoundTripper); ok {
			if httpTransport, ok := tp.(*http.Transport); ok {
				httpTransport.CloseIdleConnections()
			}
			if closer, ok := tp.(io.Closer); ok {
				closer.Close()
			}
		}
		m.transports.Delete(key)
		return true
	})
}

func (x *APIServerProxyHandler) Shutdown(ctx context.Context) error {
	x.log.Info("开始清理APIServerProxyHandler缓存...")
	if x.proxyManager != nil {
		x.proxyManager.Shutdown()
	}
	x.log.Info("APIServerProxyHandler缓存清理完成")
	return nil
}

func (m *K8sProxyManager) GetTransportForProxy(clusterID, kubeConfig string) (http.RoundTripper, error) {
	config, err := m.GetConfig(clusterID, kubeConfig)
	if err != nil {
		return nil, err
	}

	return rest.TransportFor(config)
}

func (x *APIServerProxyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// 用户认证
	ctx, userID, clusterID, userName, err := x.authenticate(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	traceID := uuid.NewString()
	ctx = setTraceIDWithContext(ctx, traceID)
	r = r.WithContext(ctx)

	// 获取集群信息
	clusterItem, err := x.cluster.GetClusterByID(r.Context(), clusterID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// 权限校验
	if err := x.authorize(r, clusterItem.Name); err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	// 审计日志
	resolver := &apirequest.RequestInfoFactory{
		APIPrefixes:          sets.NewString("api", "apis"),
		GrouplessAPIPrefixes: sets.NewString("api"),
	}
	info, err := resolver.NewRequestInfo(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		host = ""
	}

	resourceType := info.Resource
	if info.Subresource != "" {
		resourceType = fmt.Sprintf("%s/%s", info.Resource, info.Subresource)
	}
	reqInfo := &RequestInfo{
		UserID:       userID,
		UserName:     userName,
		ClientIP:     host,
		Method:       r.Method,
		Path:         r.URL.Path,
		ClusterName:  clusterItem.Name,
		Namespace:    info.Namespace,
		Resource:     resourceType,
		ResourceName: info.Name,
		StartTime:    time.Now(),
	}

	// 证书认证
	kubeConfig := clusterItem.ImportDetail.KubeConfig
	config, err := clientcmd.Load([]byte(kubeConfig))
	if err != nil {
		http.Error(w, "kubeConfig无效", http.StatusBadRequest)
		return
	}

	APIServer := clusterItem.ImportDetail.ApiServer
	if clusterItem.ConnectType == pb.ConnectType_Mesh {
		if clusterItem.MeshAddr != "" {
			APIServer = fmt.Sprintf("https://%s", clusterItem.MeshAddr)
		}

		// 遍历 clusters 并替换 server 地址
		for _, cluster := range config.Clusters {
			cluster.Server = APIServer
			cluster.InsecureSkipTLSVerify = true
			cluster.CertificateAuthorityData = nil
		}

	}

	// 启动代理服务器
	tp, err := x.proxyManager.GetTransportForProxy(strconv.Itoa(int(clusterID)), kubeConfig)
	if err != nil {
		http.Error(w, "创建传输层失败: "+err.Error(), http.StatusInternalServerError)
		return
	}

	remoteURL, err := url.Parse(APIServer)
	if err != nil {
		panic(err)
	}
	remoteURL = remoteURL.JoinPath(r.URL.Path)
	remoteURL.RawQuery = r.URL.RawQuery

	// 记录日志
	x.recordAudit(ctx, reqInfo, "", time.Now().Format("2006-01-02 15:04:05"))

	if httpstream.IsUpgradeRequest(r) {
		// 流式请求, 记录流日志
		auditBuffer := NewAuditBuffer()
		w = NewResponseWriterWrapper(w, auditBuffer)
		x.handleAuditSPDY(ctx, reqInfo, auditBuffer)
	}

	httpProxy := k8sproxy.NewUpgradeAwareHandler(remoteURL, tp, true, false, newSimpleErrorResponder(x.log))
	httpProxy.ServeHTTP(w, r)
}

func (x *APIServerProxyHandler) handleAuditSPDY(ctx context.Context, reqInfo *RequestInfo, auditBuffer *AuditBuffer) {
	recoveryFn := func() {
		r := recover()
		if r != nil {
			x.log.Errorf("[handleAuditSPDY] panic, err=%v, stack=%s", r, string(debug.Stack()))
		}
	}
	go func() {
		defer recoveryFn()

		defer x.log.Infof("Audit buffer closed")
		framer, _ := spdy.NewFramer(io.Discard, auditBuffer)
		if framer == nil {
			return
		}
		var (
			mu  sync.Mutex
			buf bytes.Buffer
		)

		go func() {
			defer recoveryFn()

			var auditBuf bytes.Buffer
			for {
				select {
				case <-ctx.Done():
					return
				case <-time.After(time.Second):
				}
				mu.Lock()
				audit, err := buf.ReadString('\n')
				mu.Unlock()

				auditBuf.WriteString(audit)
				if errors.Is(err, io.EOF) {
					continue
				}
				if err != nil {
					continue
				}
				if auditBuf.String() != "" {
					// record audit info
					x.recordAudit(ctx, reqInfo, auditBuf.String(), time.Now().Format("2006-01-02 15:04:05"))
					auditBuf.Reset()
				}
			}
		}()

		for {
			select {
			case <-ctx.Done():
				return
			case <-time.After(time.Millisecond * 100):
				frame, _ := framer.ReadFrame()
				if frame == nil {
					continue
				}
				dataFrame, ok := frame.(*spdy.DataFrame)
				if !ok {
					continue
				}
				if dataFrame == nil {
					continue
				}
				mu.Lock()
				buf.Write(dataFrame.Data)
				mu.Unlock()
			}
		}
	}()
}

func NewAPIServerProxy(bc *conf.Bootstrap, cluster *biz.ClusterUseCase, casbin *middleware.CasbinCheckMiddleware,
	redis *redis.Client, audit *biz.AuditLogUseCase, kafka dep.IKafka, logger log.Logger,
) (*APIServerProxy, error) {
	svr := &APIServerProxy{
		conf:    bc,
		cluster: cluster,
		casbin:  casbin,
		redis:   redis,
		audit:   audit,
		kafka:   kafka,
		log:     log.NewHelper(log.With(logger, "module", "server/proxy")),
	}
	svcConfig := bc.PROXY
	addr := svcConfig.GetADDR()
	network := svcConfig.GetNETWORK()
	if network == "" {
		network = "tcp"
	}
	if svcConfig.GetENABLE() {
		// 创建 TLS 配置
		tlsConfig := &tls.Config{
			MinVersion: tls.VersionTLS12,
		}

		// 加载证书
		serverCrt := filepath.Join("cert", "server.crt")
		serverKey := filepath.Join("cert", "server.pem")
		cert, err := tls.LoadX509KeyPair(serverCrt, serverKey)
		if err != nil {
			return nil, fmt.Errorf("failed to load TLS cert: %v", err)
		}
		tlsConfig.Certificates = []tls.Certificate{cert}
		// 创建 TLS 监听器
		listener, err := tls.Listen(network, addr, tlsConfig)
		svr.listener = listener
	}
	return svr, nil
}

func (x *APIServerProxy) Start(ctx context.Context) error {
	if x.listener == nil {
		return nil
	}
	if !atomic.CompareAndSwapUint32(&x.started, 0, 1) {
		return nil
	}
	handler := http.NewServeMux()
	proxyHandler := NewAPIServerHandler(x.cluster, x.casbin, x.redis, x.audit, x.kafka, x.log, x.conf)
	handler.Handle("/", proxyHandler)
	return http.Serve(x.listener, handler)
}

func (x *APIServerProxy) Stop(ctx context.Context) error {
	if atomic.CompareAndSwapUint32(&x.started, 1, 0) {
		// 清理ProxyHandler缓存
		if x.proxyHandler != nil {
			if err := x.proxyHandler.Shutdown(ctx); err != nil {
				x.log.Errorf("清理ProxyHandler缓存失败: %v", err)
			}
		}
		return x.listener.Close()
	}
	x.log.WithContext(ctx).Info("k8s proxy stopped !!!")
	return nil
}

// recordAudit 记录审计日志
func (x *APIServerProxyHandler) recordAudit(ctx context.Context, reqInfo *RequestInfo, content, timestamp string) {
	traceID, _ := getTraceIDFromContext(ctx)
	requestBody := struct {
		Content string `json:"content,omitempty"`
		Method  string `json:"method"`
		Path    string `json:"path"`
	}{
		Content: content,
		Method:  reqInfo.Method,
		Path:    reqInfo.Path,
	}

	// 转换为 JSON
	requestBodyJSON, err := json.Marshal(requestBody)
	if err != nil {
		x.log.Errorf("序列化请求体失败: %v", err)
		return
	}

	// 创建审计日志
	auditLog := &biz.CreateAuditLogRequest{
		AuditLogItem: biz.AuditLogItem{
			UserName:      reqInfo.UserName,
			ClientIP:      reqInfo.ClientIP,
			Module:        "kubectl",
			Action:        reqInfo.Method,
			Cluster:       reqInfo.ClusterName,
			Namespace:     reqInfo.Namespace,
			ResourceType:  reqInfo.Resource,
			ResourceName:  reqInfo.ResourceName,
			RequestBody:   string(requestBodyJSON),
			RequestPath:   reqInfo.Path,
			ResponseBody:  reqInfo.ResponseBody,
			Status:        int(pb.OperationStatus_Success),
			Duration:      "0",
			OperationTime: timestamp, // 使用传入的时间戳
			TraceID:       traceID,
		},
	}

	// 记录到数据库
	if err := x.audit.CreateAuditLog(context.Background(), auditLog); err != nil {
		x.log.Errorf("记录审计日志失败: %v", err)
	}

	// 发送到 Kafka
	if x.kafka != nil {
		go func(auditLog *biz.CreateAuditLogRequest) {
			ctx := context.Background()
			defer func() {
				if r := recover(); r != nil {
					x.log.Errorf("发送消息到 Kafka 时发生 panic: %v", r)
				}
			}()

			auditLogBytes, err := json.Marshal(auditLog)
			if err != nil {
				x.log.Errorf("序列化审计日志失败: %v", err)
				return
			}

			if err := x.kafka.SendMessage(ctx, auditLogBytes); err != nil {
				x.log.Errorf("发送消息到 Kafka 失败: %v", err)
			} else {
				x.log.Debug("发送消息到 Kafka 成功")
			}
		}(auditLog)
	}
}

type AuditBuffer struct {
	mu  sync.Mutex
	buf bytes.Buffer
}

func NewAuditBuffer() *AuditBuffer {
	return &AuditBuffer{}
}

func (x *AuditBuffer) String() string {
	x.mu.Lock()
	defer x.mu.Unlock()
	str := x.buf.String()
	x.buf.Reset()
	return str
}

func (x *AuditBuffer) Read(p []byte) (n int, err error) {
	x.mu.Lock()
	defer x.mu.Unlock()
	return x.buf.Read(p)
}

func (x *AuditBuffer) Write(p []byte) (n int, err error) {
	x.mu.Lock()
	defer x.mu.Unlock()
	return x.buf.Write(p)
}

type ResponseWriterWrapper struct {
	http.ResponseWriter

	buf *AuditBuffer
}

func NewResponseWriterWrapper(responseWriter http.ResponseWriter, buf *AuditBuffer) *ResponseWriterWrapper {
	return &ResponseWriterWrapper{ResponseWriter: responseWriter, buf: buf}
}

func (x *ResponseWriterWrapper) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	rw := x.ResponseWriter
	for {
		switch t := rw.(type) {
		case http.Hijacker:
			conn, reader, err := t.Hijack()
			if err != nil {
				return nil, nil, err
			}
			wrapper := NewConnWrapper(conn, x.buf)
			return wrapper, reader, nil
		case interface {
			Unwrap() http.ResponseWriter
		}:
			rw = t.Unwrap()
		default:
			return nil, nil, fmt.Errorf("%w: 无法 Hijack", http.ErrNotSupported)
		}
	}
}

type ConnWrapper struct {
	conn net.Conn
	buf  *AuditBuffer
}

func NewConnWrapper(conn net.Conn, buf *AuditBuffer) *ConnWrapper {
	return &ConnWrapper{conn: conn, buf: buf}
}

func (x *ConnWrapper) Read(b []byte) (n int, err error) {
	n, err = x.conn.Read(b)
	return
}

func (x *ConnWrapper) Write(b []byte) (n int, err error) {
	n, err = x.conn.Write(b)
	x.buf.Write(b[:n])
	return
}

func (x *ConnWrapper) Close() error {
	return x.conn.Close()
}

func (x *ConnWrapper) LocalAddr() net.Addr {
	return x.conn.LocalAddr()
}

func (x *ConnWrapper) RemoteAddr() net.Addr {
	return x.conn.RemoteAddr()
}

func (x *ConnWrapper) SetDeadline(t time.Time) error {
	return x.conn.SetDeadline(t)
}

func (x *ConnWrapper) SetReadDeadline(t time.Time) error {
	return x.conn.SetReadDeadline(t)
}

func (x *ConnWrapper) SetWriteDeadline(t time.Time) error {
	return x.conn.SetWriteDeadline(t)
}

type wrapRoundTripper struct {
	log    *log.Helper
	rt     http.RoundTripper
	config *rest.Config
}

func newWrapRoundTripper(rt http.RoundTripper, config *rest.Config, log *log.Helper) *wrapRoundTripper {
	return &wrapRoundTripper{
		log:    log,
		rt:     rt,
		config: config,
	}
}

func (x *wrapRoundTripper) TLSClientConfig() *tls.Config {
	// 创建 HTTP 客户端
	tlsConfig := &tls.Config{
		InsecureSkipVerify: true, // 忽略 Kubernetes API 证书
		Certificates:       []tls.Certificate{},
	}

	certData := x.config.CertData
	keyData := x.config.KeyData
	// certData, keyData, _ := GenerateCert()
	// 创建 X509KeyPair
	cert, err3 := tls.X509KeyPair(certData, keyData)
	if err3 != nil {
		x.log.Errorf("failed to create X509KeyPair: %v", err3)
	} else {
		tlsConfig.Certificates = []tls.Certificate{cert}
	}
	return tlsConfig
}

func (x *wrapRoundTripper) RoundTrip(request *http.Request) (*http.Response, error) {
	// 首先检查是否为超级管理员，如果是则直接放行
	resp, err := x.rt.RoundTrip(request)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

type simpleErrorResponder struct {
	log *log.Helper
}

func newSimpleErrorResponder(log *log.Helper) *simpleErrorResponder {
	return &simpleErrorResponder{log: log}
}

func (x *simpleErrorResponder) Error(w http.ResponseWriter, req *http.Request, err error) {
	x.log.Errorf("[simpleErrorResponder] Error: %v", err)
	http.Error(w, err.Error(), http.StatusInternalServerError)
}

type traceIDKey struct{}

func setTraceIDWithContext(ctx context.Context, traceID string) context.Context {
	return context.WithValue(ctx, traceIDKey{}, traceID)
}

func getTraceIDFromContext(ctx context.Context) (string, bool) {
	traceID, ok := ctx.Value(traceIDKey{}).(string)
	return traceID, ok
}
