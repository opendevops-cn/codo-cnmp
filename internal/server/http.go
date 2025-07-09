package server

import (
	"time"

	"codo-cnmp/common/result"
	"codo-cnmp/internal/conf"
	middleware2 "codo-cnmp/internal/middleware"
	"codo-cnmp/internal/service"
	pb "codo-cnmp/pb"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/logging"
	"github.com/go-kratos/kratos/v2/middleware/metrics"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/opendevops-cn/codo-golang-sdk/adapter/kratos/middleware/ktracing"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
)

// NewHTTPServer new an HTTP server.
func NewHTTPServer(bc *conf.Bootstrap, greeter *service.GreeterService, cluster *service.ClusterService,
	namespace *service.NameSpaceService, node *service.NodeService, pod *service.PodService,
	deployment *service.DeploymentService, event *service.EventService, cloneSet *service.CloneSetService,
	gameServerSet *service.GameServerSetService, usergroup *service.UserGroupService, role *service.RoleService,
	statefulSet *service.StatefulSetService, hpa *service.HpaService, configMap *service.ConfigMapService,
	secret *service.SecretService, userFlow *service.UserFollowService, auditLog *service.AuditLogService,
	daemonSet *service.DaemonSetService, gameServer *service.GameServerService, ezRollout *service.EzRolloutService,
	resource *service.ResourceService, resourceQuota *service.ResourceQuotaService, limitRange *service.LimitRangeService,
	logger log.Logger, mp metric.MeterProvider, tp trace.TracerProvider, svc *service.SvcService, ingress *service.IngressService,
	crd *service.CRDService, apiGroup *service.ApiGroupService, sidecarSet *service.SideCarSetService, sc *service.ScService,
	pv *service.PvService, pvc *service.PvcService, agent *service.AgentService, crr *service.CRRService,
	ingressClass *service.IngressClassService,
	// middleware
	casbinCheckMiddleware *middleware2.CasbinCheckMiddleware, auditMiddleware *middleware2.AuditMiddleware,
) (*http.Server, error) {
	c := bc.APP
	meter := mp.Meter("APP.ADDR")

	counter, err := metrics.DefaultRequestsCounter(meter, metrics.DefaultServerRequestsCounterName)
	if err != nil {
		return nil, err
	}
	seconds, err := metrics.DefaultSecondsHistogram(meter, metrics.DefaultServerSecondsHistogramName)
	if err != nil {
		return nil, err
	}

	opts := []http.ServerOption{
		http.Middleware(
			recovery.Recovery(),
			ktracing.Server(
				ktracing.WithTracerProvider(tp),
			),
			logging.Server(logger),
			metrics.Server(
				metrics.WithRequests(counter),
				metrics.WithSeconds(seconds),
			),
			// 注册校验中间件
			middleware2.NewValidateMiddleware().Server(),
			// 注册casbin中间件
			casbinCheckMiddleware.Server(),
			// 注册审计中间件
			auditMiddleware.Server(),
		),
		http.ResponseEncoder(result.KHttpResult),
		http.ErrorEncoder(result.KHttpError),
	}
	network := c.NETWORK
	if network == "" {
		opts = append(opts, http.Network("tcp"))
	} else {
		opts = append(opts, http.Network(network))
	}
	if c.ADDR != "" {
		opts = append(opts, http.Address(c.ADDR))
	} else {
		opts = append(opts, http.Address("0.0.0.0:8000"))
	}
	timeout := c.TIMEOUT
	if timeout == 0 {
		opts = append(opts, http.Timeout(30*time.Second))
	} else {
		opts = append(opts, http.Timeout(time.Duration(timeout)*time.Second))
	}
	srv := http.NewServer(opts...)
	pb.RegisterGreeterHTTPServer(srv, greeter)
	pb.RegisterClusterHTTPServer(srv, cluster)
	pb.RegisterNameSpaceHTTPServer(srv, namespace)
	pb.RegisterNodeHTTPServer(srv, node)
	pb.RegisterPodHTTPServer(srv, pod)
	pb.RegisterDeploymentHTTPServer(srv, deployment)
	pb.RegisterEventHTTPServer(srv, event)
	pb.RegisterCloneSetHTTPServer(srv, cloneSet)
	pb.RegisterGameServerSetHTTPServer(srv, gameServerSet)
	pb.RegisterUserGroupServiceHTTPServer(srv, usergroup)
	pb.RegisterRoleServiceHTTPServer(srv, role)
	pb.RegisterStatefulSetHTTPServer(srv, statefulSet)
	pb.RegisterHPAHTTPServer(srv, hpa)
	pb.RegisterConfigMapHTTPServer(srv, configMap)
	pb.RegisterSecretHTTPServer(srv, secret)
	pb.RegisterUserFollowHTTPServer(srv, userFlow)
	pb.RegisterAuditLogHTTPServer(srv, auditLog)
	pb.RegisterDaemonSetHTTPServer(srv, daemonSet)
	pb.RegisterGameServerHTTPServer(srv, gameServer)
	pb.RegisterEzRolloutHTTPServer(srv, ezRollout)
	pb.RegisterResourceHTTPServer(srv, resource)
	pb.RegisterResourceQuotaHTTPServer(srv, resourceQuota)
	pb.RegisterLimitRangeHTTPServer(srv, limitRange)
	pb.RegisterSVCHTTPServer(srv, svc)
	pb.RegisterIngressHTTPServer(srv, ingress)
	pb.RegisterCRDHTTPServer(srv, crd)
	pb.RegisterAPIGroupHTTPServer(srv, apiGroup)
	pb.RegisterSidecarSetHTTPServer(srv, sidecarSet)
	pb.RegisterStorageClassHTTPServer(srv, sc)
	pb.RegisterPersistentVolumeHTTPServer(srv, pv)
	pb.RegisterPersistentVolumeClaimHTTPServer(srv, pvc)
	pb.RegisterAgentHTTPServer(srv, agent)
	pb.RegisterCRRHTTPServer(srv, crr)
	pb.RegisterIngressClassHTTPServer(srv, ingressClass)
	return srv, nil
}
