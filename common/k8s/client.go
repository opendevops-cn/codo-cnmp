package k8s

import (
	"codo-cnmp/common/consts"
	"codo-cnmp/common/utils"
	"codo-cnmp/common/xerr"
	"codo-cnmp/pb"
	"fmt"
	kruiseclientset "github.com/openkruise/kruise-api/client/clientset/versioned"
	kruisegameclientset "github.com/openkruise/kruise-game/pkg/client/clientset/versioned"
	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/apiserver/pkg/authorization/authorizer"
	apirequest "k8s.io/apiserver/pkg/endpoints/request"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/metrics/pkg/client/clientset/versioned"
	"net/http"
	"strings"
)

type ClientConfig struct {
	ImportType pb.ImportType
	KubeConfig string `json:"kube_config"`
	// token详情
	Token string `json:"token"`
	// apiServer地址
	APIServer string `json:"api_server"`
	// 是否使用自签名证书
	Insecure bool `json:"insecure"`
	// ca证书数据
	CaData string `json:"ca_data"`
	// context名称
	Context string `json:"context"`
	// agent名称
	Agent string `json:"agent"`
	// agent镜像
	AgentImage string `json:"agent_image"`
	// agent proxy地址
	AgentProxy string `json:"agent_proxy"`
	// ClusterName 集群名称
	ClusterName string `json:"cluster_name"`
	// ConnectType 连接类型
	ConnectType pb.ConnectType `json:"connect_type"`
	// MeshAddr 网格地址
	MeshAddr string `json:"mesh_addr"`
}

// isSuperUser checks if the request context contains superUser privileges.
func isSuperUser(req *http.Request) bool {
	superUserValueFlag := req.Context().Value(consts.ContextSuperUserKey)
	if superUserValueFlag != nil {
		if superUser, ok := superUserValueFlag.(bool); ok && superUser {
			return true
		}
	}
	return false
}

// RBACAuthorizer RBAC授权器
type RBACAuthorizer struct {
	roles []rbacv1.ClusterRole
}

// NewRBACAuthorizer 创建新的RBAC授权器
func NewRBACAuthorizer(roles []rbacv1.ClusterRole) *RBACAuthorizer {
	return &RBACAuthorizer{roles: roles}
}

// Authorize 检查是否允许访问
func (r *RBACAuthorizer) Authorize(attrs authorizer.Attributes) (bool, string, error) {
	for _, role := range r.roles {
		for _, rule := range role.Rules {
			if RulesAllow(attrs, rule) {
				return true, "", nil
			}
		}
	}
	return false, "没有操作权限", nil
}

// NewRequestInfoResolver creates a new RequestInfoFactory.
func NewRequestInfoResolver() *apirequest.RequestInfoFactory {
	return &apirequest.RequestInfoFactory{
		APIPrefixes:          sets.NewString("api", "apis"),
		GrouplessAPIPrefixes: sets.NewString("api"),
	}
}

// RulesAllow checks if the given rules allow the given request attributes.
// 参考: k8s.io/kubernetes/plugin/pkg/auth/authorizer/rbac
func RulesAllow(requestAttributes authorizer.Attributes, rules ...rbacv1.PolicyRule) bool {
	for i := range rules {
		if RuleAllows(requestAttributes, &rules[i]) {
			return true
		}
	}

	return false
}

// RuleAllows checks if the given rule allows the given request attributes.
// 参考: k8s.io/kubernetes/plugin/pkg/auth/authorizer/rbac
func RuleAllows(requestAttributes authorizer.Attributes, rule *rbacv1.PolicyRule) bool {
	if requestAttributes.IsResourceRequest() {
		combinedResource := requestAttributes.GetResource()
		if len(requestAttributes.GetSubresource()) > 0 {
			combinedResource = requestAttributes.GetResource() + "/" + requestAttributes.GetSubresource()
		}

		return utils.VerbMatches(rule, requestAttributes.GetVerb()) &&
			utils.APIGroupMatches(rule, requestAttributes.GetAPIGroup()) &&
			utils.ResourceMatches(rule, combinedResource, requestAttributes.GetSubresource()) &&
			utils.ResourceNameMatches(rule, requestAttributes.GetName())
	}

	return utils.VerbMatches(rule, requestAttributes.GetVerb()) &&
		utils.NonResourceURLMatches(rule, requestAttributes.GetPath())
}

// AuthorizeRequest 统一的鉴权处理函数
func AuthorizeRequest(req *http.Request, clusterName string) error {
	// 首先检查是否为超级管理员
	if isSuperUser(req) {
		return nil
	}

	// 对于 GET 请求直接放行
	if req.Method == http.MethodGet {
		return nil
	}

	// 只对写操作进行权限校验
	if req.Method != http.MethodPost &&
		req.Method != http.MethodPut &&
		req.Method != http.MethodPatch &&
		req.Method != http.MethodDelete {
		return nil
	}

	requestInfoResolver := NewRequestInfoResolver()
	reqCopy := req.Clone(req.Context())
	requestInfo, err := requestInfoResolver.NewRequestInfo(reqCopy)
	if err != nil {
		return xerr.NewErrCodeMsg(xerr.ServerCommonError, "系统内部错误")
	}

	return authorizeRequestInfo(req, clusterName, requestInfo)
}

func authorizeRequestInfo(req *http.Request, clusterName string, requestInfo *apirequest.RequestInfo) error {
	// 从上下文获取角色配置
	clusterRoleBindingsValue := req.Context().Value(consts.ContextClusterRoleBindingsKey)
	if clusterRoleBindingsValue == nil {
		return xerr.NewForbiddenErrMsg("未找到集群角色绑定关系")
	}
	clusterBindings, ok := clusterRoleBindingsValue.([]map[string][]string)
	if !ok {
		return xerr.NewForbiddenErrMsg("集群角色绑定关系格式错误")
	}
	if len(clusterBindings) == 0 {
		return xerr.NewForbiddenErrMsg("没有操作权限")
	}
	hasClusterPermission := false
	for _, binding := range clusterBindings {
		namespaces, exists := binding[clusterName]
		if !exists || len(namespaces) == 0 {
			continue
		}
		// 如果 namespaces包含 "*" 或者 requestInfo.Namespace 在 namespaces 中 或者 requestInfo.Namespace 为空，则有权限
		if requestInfo.Namespace == "" || utils.Contains(namespaces, "*") || utils.Contains(namespaces, requestInfo.Namespace) {
			hasClusterPermission = true
			break
		}
	}
	if !hasClusterPermission {
		return xerr.NewForbiddenErrMsg("没有操作权限")
	}
	clusterRolesValue := req.Context().Value(consts.ContextClusterRolesKey)
	if clusterRolesValue == nil {
		return xerr.NewForbiddenErrMsg("未找到角色配置")
	}
	roles, ok := clusterRolesValue.([]rbacv1.ClusterRole)
	if !ok {
		return xerr.NewForbiddenErrMsg("角色配置格式错误")
	}
	if len(roles) == 0 {
		return xerr.NewForbiddenErrMsg("没有操作权限")
	}
	attrs := &authorizer.AttributesRecord{
		User:            nil,
		Verb:            requestInfo.Verb,
		Namespace:       requestInfo.Namespace,
		APIGroup:        requestInfo.APIGroup,
		APIVersion:      requestInfo.APIVersion,
		Resource:        requestInfo.Resource,
		Subresource:     requestInfo.Subresource,
		Name:            requestInfo.Name,
		ResourceRequest: requestInfo.IsResourceRequest,
		Path:            requestInfo.Path,
	}
	newAuthorizer := NewRBACAuthorizer(roles)
	// 进行授权检查
	allowed, _, err := newAuthorizer.Authorize(attrs)
	if err != nil {
		return xerr.NewErrCodeMsg(xerr.ServerCommonError, fmt.Sprintf("授权检查失败: %v", err))
	}

	if !allowed {
		return xerr.NewForbiddenErrMsg("没有操作权限")
	}
	return nil
}

// handleRequest 该方法主要用于检查请求是否有权限，并根据请求方法调用相应的处理函数。
func handleRequest(req *http.Request, clusterName string, requestInfo *apirequest.RequestInfo, rt http.RoundTripper) (*http.Response, error) {
	// 从上下文获取角色配置
	clusterRoleBindingsValue := req.Context().Value(consts.ContextClusterRoleBindingsKey)
	if clusterRoleBindingsValue == nil {
		return nil, xerr.NewForbiddenErrMsg("未找到集群角色绑定关系")
	}
	clusterBindings, ok := clusterRoleBindingsValue.([]map[string][]string)
	if !ok {
		return nil, xerr.NewForbiddenErrMsg("集群角色绑定关系格式错误")
	}
	if len(clusterBindings) == 0 {
		return nil, xerr.NewForbiddenErrMsg("没有操作权限")
	}
	hasClusterPermission := false
	for _, binding := range clusterBindings {
		namespaces, exists := binding[clusterName]
		if !exists || len(namespaces) == 0 {
			continue
		}
		// 如果 namespaces包含 "*" 或者 requestInfo.Namespace 在 namespaces 中 或者 requestInfo.Namespace 为空，则有权限
		if requestInfo.Namespace == "" || utils.Contains(namespaces, "*") || utils.Contains(namespaces, requestInfo.Namespace) {
			hasClusterPermission = true
			break
		}
	}
	if !hasClusterPermission {
		return nil, xerr.NewForbiddenErrMsg("没有操作权限")
	}
	clusterRolesValue := req.Context().Value(consts.ContextClusterRolesKey)
	if clusterRolesValue == nil {
		return nil, xerr.NewForbiddenErrMsg("未找到角色配置")
	}
	roles, ok := clusterRolesValue.([]rbacv1.ClusterRole)
	if !ok {
		return nil, xerr.NewForbiddenErrMsg("角色配置格式错误")
	}
	if len(roles) == 0 {
		return nil, xerr.NewForbiddenErrMsg("没有操作权限")
	}
	attrs := &authorizer.AttributesRecord{
		User:            nil,
		Verb:            requestInfo.Verb,
		Namespace:       requestInfo.Namespace,
		APIGroup:        requestInfo.APIGroup,
		APIVersion:      requestInfo.APIVersion,
		Resource:        requestInfo.Resource,
		Subresource:     requestInfo.Subresource,
		Name:            requestInfo.Name,
		ResourceRequest: requestInfo.IsResourceRequest,
		Path:            requestInfo.Path,
	}
	newAuthorizer := NewRBACAuthorizer(roles)
	// 进行授权检查
	allowed, _, err := newAuthorizer.Authorize(attrs)
	if err != nil {
		return nil, xerr.NewErrCodeMsg(xerr.ServerCommonError, fmt.Sprintf("授权检查失败: %v", err))
	}

	if !allowed {
		return nil, xerr.NewForbiddenErrMsg("没有操作权限")
	}
	return rt.RoundTrip(req)
}

// UpdateAPIServer 更新 kubeConfig 中的 APIServer 地址
func (cfg *ClientConfig) UpdateAPIServer(newAPIServer string) (string, error) {
	// 解析 kubeConfig 字符串
	config, err := clientcmd.Load([]byte(cfg.KubeConfig))
	if err != nil {
		return "", fmt.Errorf("failed to load kubeconfig: %v", err)
	}

	// 遍历 clusters 并替换 server 地址
	for _, cluster := range config.Clusters {
		cluster.Server = newAPIServer
		cluster.InsecureSkipTLSVerify = true
		cluster.CertificateAuthorityData = nil
	}

	// 转换回字符串格式
	modifiedConfig, err := clientcmd.Write(*config)
	if err != nil {
		return "", fmt.Errorf("failed to serialize kubeconfig: %v", err)
	}

	return string(modifiedConfig), nil
}

// BuildRestConfig 根据不同的配置类型创建 rest.Config
func (cfg *ClientConfig) BuildRestConfig() (*rest.Config, error) {
	// 创建 rest.Config
	var (
		c   *rest.Config
		err error
	)

	if cfg.ConnectType == pb.ConnectType_Mesh {
		// 如果是组网连接，则使用mesh地址，保持原有协议
		var newHost string
		parts := strings.SplitN(cfg.APIServer, "://", 2)
		if len(parts) == 2 {
			newHost = parts[0] + "://" + cfg.MeshAddr
		} else {
			newHost = cfg.MeshAddr
		}
		newKubeConfig, err := cfg.UpdateAPIServer(newHost)
		if err != nil {
			return nil, fmt.Errorf("替换APIServer 失败: %v", err)
		}
		cfg.KubeConfig = newKubeConfig
	}

	switch cfg.ImportType {
	case pb.ImportType_KubeConfig:
		c, err = cfg.BuildRestConfigFromKubeConfigString()
	case pb.ImportType_Token:
		c, err = cfg.BuildRestConfigFromToken()
	default:
		return nil, fmt.Errorf("不支持导入类型: %d", cfg.ImportType)
	}
	if err != nil {
		return nil, fmt.Errorf("创建RestConfig失败: %w", err)
	}
	// 自定义Transport, 用于权限校验
	c.WrapTransport = func(rt http.RoundTripper) http.RoundTripper {
		return httpRoundTripperWrapper(func(req *http.Request) (*http.Response, error) {
			// 首先检查是否为超级管理员，如果是则直接放行
			if isSuperUser(req) {
				return rt.RoundTrip(req)
			}
			// 非超级管理员的请求处理

			switch req.Method {
			case http.MethodGet:
				// get请求直接放行
				return rt.RoundTrip(req)
			case http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete:
				requestInfoResolver := NewRequestInfoResolver()
				reqCopy := req.Clone(req.Context())
				requestInfo, err := requestInfoResolver.NewRequestInfo(reqCopy)
				if err != nil {
					return nil, xerr.NewErrCodeMsg(xerr.ServerCommonError, "系统内部错误")
				}
				return handleRequest(req, cfg.ClusterName, requestInfo, rt)
			default:
				return rt.RoundTrip(req)
			}
		})
	}
	return c, nil
}

// BuildRestConfigFromKubeConfigString 通过 kubeConfig 创建 clientConfig
func (cfg *ClientConfig) BuildRestConfigFromKubeConfigString() (*rest.Config, error) {
	overrideClientConfig, err := clientcmd.NewClientConfigFromBytes([]byte(cfg.KubeConfig))
	if err != nil {
		return nil, fmt.Errorf("解析kubeConfig失败: %w", err)
	}

	clientConfig, err := overrideClientConfig.ClientConfig()
	if err != nil {
		return nil, fmt.Errorf("创建restConfig失败: %w", err)
	}

	return clientConfig, nil
}

// BuildRestConfigFromToken 通过 token 字符串创建 clientConfig
func (cfg *ClientConfig) BuildRestConfigFromToken() (*rest.Config, error) {
	return &rest.Config{
		Host:        cfg.APIServer,
		BearerToken: cfg.Token,
		TLSClientConfig: rest.TLSClientConfig{
			Insecure: true, // 如果集群使用自签名证书或没有提供证书
		},
	}, nil
}

// CreateDynamicClient 创建 dynamic client
func (cfg *ClientConfig) CreateDynamicClient() (*dynamic.DynamicClient, error) {
	clientConfig, err := cfg.BuildRestConfig()
	if err != nil {
		return nil, err
	}

	dynamicClient, err := dynamic.NewForConfig(clientConfig)
	if err != nil {
		return nil, fmt.Errorf("创建dynamicClient失败: %w", err)
	}

	return dynamicClient, nil
}

// CreateClientSet 创建 clientSet
func (cfg *ClientConfig) CreateClientSet() (*kubernetes.Clientset, error) {
	clientConfig, err := cfg.BuildRestConfig()
	if err != nil {
		return nil, err
	}

	clientSet, err := kubernetes.NewForConfig(clientConfig)
	if err != nil {
		return nil, fmt.Errorf("创建clientSet失败: %w", err)
	}

	return clientSet, nil
}

// CreateMetricsClientSet 创建 metrics clientSet
func (cfg *ClientConfig) CreateMetricsClientSet() (*versioned.Clientset, error) {
	clientConfig, err := cfg.BuildRestConfig()
	if err != nil {
		return nil, err
	}

	clientSet, err := versioned.NewForConfig(clientConfig)
	if err != nil {
		return nil, fmt.Errorf("创建metrics clientSet失败: %w", err)
	}

	return clientSet, nil
}

// CreateKruiseClientSet 创建 kruise clientSet
func (cfg *ClientConfig) CreateKruiseClientSet() (*kruiseclientset.Clientset, error) {
	clientConfig, err := cfg.BuildRestConfig()
	if err != nil {
		return nil, err
	}
	clientSet := kruiseclientset.NewForConfigOrDie(clientConfig)
	return clientSet, nil
}

// CreateKruiseGameClientSet 创建 kruise game clientSet
func (cfg *ClientConfig) CreateKruiseGameClientSet() (*kruisegameclientset.Clientset, error) {
	clientConfig, err := cfg.BuildRestConfig()
	if err != nil {
		return nil, err
	}
	clientSet := kruisegameclientset.NewForConfigOrDie(clientConfig)
	return clientSet, nil
}

type httpRoundTripperWrapper func(req *http.Request) (*http.Response, error)

func (x httpRoundTripperWrapper) RoundTrip(request *http.Request) (*http.Response, error) {
	return x(request)
}
