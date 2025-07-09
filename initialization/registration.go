package initialization

import (
	"codo-cnmp/internal/dep"
	pb "codo-cnmp/pb"
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"io"

	"github.com/Ccheers/protoc-gen-go-kratos-http/route"
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(NewRegistrationMeta)

type MenuItem struct {
	Name    string `json:"name"`
	Details string `json:"details"`
}

var menuList = []MenuItem{
	{
		Name:    "CNMP_AUDIT_LOG",
		Details: "操作审计",
	},
	{
		Name:    "CNMP_PLAY_SERVER",
		Details: "进程列表",
	},
	{
		Name:    "CNMP_HPA_EZ",
		Details: "版本伸缩",
	},
	{
		Name:    "CNMP_AUTH_ROLE",
		Details: "角色管理",
	},
	{
		Name:    "CNMP_AUTH_USER",
		Details: "用户授权",
	},
	{
		Name:    "CNMP_SECRET",
		Details: "Secret",
	},
	{
		Name:    "CNMP_CONFIG_MAP",
		Details: "ConfigMap",
	},
	{
		Name:    "CNMP_HPA",
		Details: "弹性伸缩",
	},
	{
		Name:    "CNMP_POD",
		Details: "Pod",
	},
	{
		Name:    "CNMP_DAEMON_SET",
		Details: "DaemonSet",
	},
	{
		Name:    "CNMP_STATEFUL_SET",
		Details: "StatefulSet",
	},
	{
		Name:    "CNMP_GAME_SERVER_SET",
		Details: "GameServerSet",
	},
	{
		Name:    "CNMP_CLONE_SET",
		Details: "CloneSet",
	},
	{
		Name:    "CNMP_DEPLOYMENT",
		Details: "Deployment",
	},
	{
		Name:    "CNMP_NODE",
		Details: "节点管理",
	},
	{
		Name:    "CNMP_NAMESPACE",
		Details: "命名空间",
	},
	{
		Name:    "CNMP_CLUSTER_LIST",
		Details: "集群列表",
	},
	{
		Name:    "CNMP_DASHBOARD",
		Details: "概览",
	},
}

type RegistrationMeta struct {
	apiGw *dep.CODOAPIGateway
	log   *log.Helper
}

type RegisterData struct {
	MethodType string `json:"method_type"`
	Path       string `json:"uri"`
	Name       string `json:"name"`
}

func NewRegistrationMeta(apiGw *dep.CODOAPIGateway, logger log.Logger) *RegistrationMeta {
	return &RegistrationMeta{
		apiGw: apiGw,
		log:   log.NewHelper(log.With(logger, "module", "initialization/registration")),
	}
}

func (r *RegistrationMeta) Run(ctx context.Context) error {
	routers := make([]route.Route, 0)
	// 集群相关路由
	clusterRoutes := pb.GenerateClusterHTTPServerRouteInfo()
	routers = append(routers, clusterRoutes...)

	// cloneSet相关路由
	cloneSetRoutes := pb.GenerateCloneSetHTTPServerRouteInfo()
	routers = append(routers, cloneSetRoutes...)

	// deployment相关路由
	deploymentRoutes := pb.GenerateDeploymentHTTPServerRouteInfo()
	routers = append(routers, deploymentRoutes...)

	// node相关路由
	nodeRoutes := pb.GenerateNodeHTTPServerRouteInfo()
	routers = append(routers, nodeRoutes...)

	// pod相关路由
	podRoutes := pb.GeneratePodHTTPServerRouteInfo()
	routers = append(routers, podRoutes...)

	// role相关路由
	roleRoutes := pb.GenerateRoleServiceHTTPServerRouteInfo()
	routers = append(routers, roleRoutes...)

	// daemonSet相关路由
	daemonSetRoutes := pb.GenerateDaemonSetHTTPServerRouteInfo()
	routers = append(routers, daemonSetRoutes...)

	// statefulSet相关路由
	statefulSetRoutes := pb.GenerateStatefulSetHTTPServerRouteInfo()
	routers = append(routers, statefulSetRoutes...)

	// event相关路由
	eventRoutes := pb.GenerateEventHTTPServerRouteInfo()
	routers = append(routers, eventRoutes...)

	// gameServerSet相关路由
	gameServerSetRoutes := pb.GenerateGameServerSetHTTPServerRouteInfo()
	routers = append(routers, gameServerSetRoutes...)

	// hpa相关路由
	hpaRoutes := pb.GenerateHPAHTTPServerRouteInfo()
	routers = append(routers, hpaRoutes...)

	// configMap相关路由
	configMapRoutes := pb.GenerateConfigMapHTTPServerRouteInfo()
	routers = append(routers, configMapRoutes...)

	// secret相关路由
	secretRoutes := pb.GenerateSecretHTTPServerRouteInfo()
	routers = append(routers, secretRoutes...)

	// 用户关注相关路由
	userFollowRoutes := pb.GenerateUserFollowHTTPServerRouteInfo()
	routers = append(routers, userFollowRoutes...)

	// 审计日志相关路由
	auditLogRoutes := pb.GenerateAuditLogHTTPServerRouteInfo()
	routers = append(routers, auditLogRoutes...)

	// 命名空间相关路由
	namespaceRoutes := pb.GenerateNameSpaceHTTPServerRouteInfo()
	routers = append(routers, namespaceRoutes...)

	// websocket相关路由
	websocketRoutes := pb.GenerateWebSocketHTTPServerRouteInfo()
	routers = append(routers, websocketRoutes...)

	// 授权管理相关路由
	userGroupRoutes := pb.GenerateUserGroupServiceHTTPServerRouteInfo()
	routers = append(routers, userGroupRoutes...)

	// 版本伸缩相关路由
	ezRolloutRoutes := pb.GenerateEzRolloutHTTPServerRouteInfo()
	routers = append(routers, ezRolloutRoutes...)

	// 进程管理相关路由
	gameServerRoutes := pb.GenerateGameServerHTTPServerRouteInfo()
	routers = append(routers, gameServerRoutes...)

	// LimitRange相关路由
	limitRangeRoutes := pb.GenerateLimitRangeHTTPServerRouteInfo()
	routers = append(routers, limitRangeRoutes...)

	// ResourceQuota相关路由
	resourceQuotaRoutes := pb.GenerateResourceQuotaHTTPServerRouteInfo()
	routers = append(routers, resourceQuotaRoutes...)

	// Resource相关路由
	resourceRoutes := pb.GenerateResourceHTTPServerRouteInfo()
	routers = append(routers, resourceRoutes...)

	// Svc相关路由
	svcRoutes := pb.GenerateSVCHTTPServerRouteInfo()
	routers = append(routers, svcRoutes...)

	// Ingress 相关路由
	ingressRoutes := pb.GenerateIngressHTTPServerRouteInfo()
	routers = append(routers, ingressRoutes...)

	// APiGroup 相关路由
	apiGroupRoutes := pb.GenerateAPIGroupHTTPServerRouteInfo()
	routers = append(routers, apiGroupRoutes...)

	// CRD 相关路由
	crdRoutes := pb.GenerateCRDHTTPServerRouteInfo()
	routers = append(routers, crdRoutes...)

	// Sidecar 相关路由
	sidecarRoutes := pb.GenerateSidecarSetHTTPServerRouteInfo()
	routers = append(routers, sidecarRoutes...)

	// StorageClass 相关路由
	scRoutes := pb.GenerateStorageClassHTTPServerRouteInfo()
	routers = append(routers, scRoutes...)

	// PersistentVolume 相关路由
	pvRoutes := pb.GeneratePersistentVolumeHTTPServerRouteInfo()
	routers = append(routers, pvRoutes...)

	// PersistentVolumeClaim 相关路由
	pvcRoutes := pb.GeneratePersistentVolumeClaimHTTPServerRouteInfo()
	routers = append(routers, pvcRoutes...)

	// agent相关路由
	agentRoutes := pb.GenerateAgentHTTPServerRouteInfo()
	routers = append(routers, agentRoutes...)

	// Crr相关路由
	crrRoutes := pb.GenerateCRRHTTPServerRouteInfo()
	routers = append(routers, crrRoutes...)

	// ingressClass相关路由
	ingressClassRoutes := pb.GenerateIngressClassHTTPServerRouteInfo()
	routers = append(routers, ingressClassRoutes...)

	funcList := make([]RegisterData, 0, len(routers))
	for _, v := range routers {
		funcList = append(funcList, RegisterData{
			MethodType: v.Method,
			Path:       fmt.Sprintf("/api/cnmp%s", v.Path),
			Name:       v.Comment,
		})

	}
	body := map[string]interface{}{
		"app_code":       "cnmp",
		"menu_list":      menuList,
		"component_list": []string{},
		"func_list":      funcList,
		"role_list":      []string{},
	}

	jsonStr, err := json.Marshal(body)
	if err != nil {
		r.log.WithContext(ctx).Errorf("序列化路由信息失败: %v", err)
		return err
	}

	response, err := r.apiGw.SendRequest(ctx, "POST", "/api/p/v4/authority/register/", jsonStr, nil)
	if err != nil {
		r.log.WithContext(ctx).Errorf("注册路由信息到API网关失败: %v", err)
		return err
	}
	defer response.Body.Close()
	data, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}
	r.log.WithContext(ctx).Infof("注册路由信息到API网关结果: %v", string(data))
	return nil
}
