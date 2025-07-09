package middleware

import (
	"codo-cnmp/common/consts"
	"codo-cnmp/common/utils"
	"codo-cnmp/common/xerr"
	"codo-cnmp/internal/biz"
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	kmiddleware "github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"
	khttp "github.com/go-kratos/kratos/v2/transport/http"
	"github.com/go-redis/redis/v8"
	"github.com/golang-jwt/jwt/v5"
	"github.com/opendevops-cn/codo-golang-sdk/adapter/kratos/transport/websocket"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
	rbacv1 "k8s.io/api/rbac/v1"
	"net"
	"net/http"
	"strings"
)

var (
	ErrKratosTransportNotFound         = fmt.Errorf("kratos: transport not found in ctx")
	ErrKratosHTTPContextNotFound       = fmt.Errorf("kratos: http.context not found in ctx")
	ErrKratosTransportNotHTTPTransport = fmt.Errorf("kratos: transport not a http.Transporter")
	AuthFailed                         = xerr.NewErrCodeMsg(xerr.ErrUnAuthorization, "auth failed")
)

func extraHTTPRequestFromKratosContext(ctx context.Context) (*http.Request, error) {
	httpTr, err := extraKratosHTTPTransport(ctx)
	if err != nil {
		return nil, err
	}
	return httpTr.Request(), nil
}
func extraKratosHTTPTransport(ctx context.Context) (khttp.Transporter, error) {
	tr, ok := transport.FromServerContext(ctx)
	if !ok {
		return nil, ErrKratosTransportNotFound
	}
	httpTr, ok := tr.(khttp.Transporter)
	if !ok {
		return nil, ErrKratosTransportNotHTTPTransport
	}
	return httpTr, nil
}
func extraKratosHTTPContext(ctx context.Context) (khttp.Context, error) {
	httpCtx, ok := ctx.(khttp.Context)
	if !ok {
		return nil, ErrKratosHTTPContextNotFound
	}
	return httpCtx, nil
}

type UserClaimsData struct {
	UserID      string `json:"user_id"`
	Username    string `json:"username"`
	Nickname    string `json:"nickname"`
	Email       string `json:"email"`
	IsSuperuser bool   `json:"is_superuser"`
}

type UserClaims struct {
	Data UserClaimsData `json:"data"`
	jwt.RegisteredClaims
}

// DecodeAuthToken 解析和验证 JWT Token
func (x *CasbinCheckMiddleware) DecodeAuthToken(authKey string) (*UserClaims, error) {
	// 解析token
	token, err := jwt.ParseWithClaims(authKey, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		// 验证签名算法
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("加密算法错误: %v", token.Header["alg"])
		}
		// 返回签名密钥
		return []byte(x.JwtSecret), nil
	})
	if err != nil {
		return nil, err
	}

	// 验证token是否有效
	if claims, ok := token.Claims.(*UserClaims); ok && token.Valid {
		// 如果token有效且解析无误，返回解析到的自定义数据
		return claims, nil
	}
	return nil, jwt.ErrSignatureInvalid
}

func (x *CasbinCheckMiddleware) DecodeAuthTokenWithoutSignature(authKey string) (*UserClaims, error) {
	// 解析token并且跳过签名验证
	token, _, err := jwt.NewParser(jwt.WithoutClaimsValidation()).ParseUnverified(authKey, &UserClaims{})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*UserClaims); ok {
		return claims, nil
	}
	return nil, err
}

type CasbinCheckMiddleware struct {
	TokenVerify bool
	JwtSecret   string
	redis       *redis.Client
	userGroup   *biz.UserGroupUseCase
	role        *biz.RoleUseCase
	cluster     *biz.ClusterUseCase
	contextKeys struct {
		isSuperuser         string
		userID              string
		userName            string
		aclList             string
		clusterRoles        string
		clusterRoleBindings string
		clientIP            string
		traceID             string
	}
	logger *log.Helper
}

func (x *CasbinCheckMiddleware) GetUserGroupIDS(ctx context.Context, userID string) ([]int, error) {
	// 从redis中获取mg用户组ID
	var roleIDs []int
	key := fmt.Sprintf(consts.UserGroupIDCacheKey, userID)
	result, err := x.redis.Get(ctx, key).Result()
	if err != nil {
		return roleIDs, xerr.NewErrCodeMsg(xerr.ErrNotAllowed, "用户非超管且不属于任何用户组")
	}
	err = json.Unmarshal([]byte(result), &roleIDs)
	if err != nil {
		return nil, err
	}
	return roleIDs, nil
}

func (x *CasbinCheckMiddleware) CodoLogin(r *http.Request) (*UserClaims, error) {
	// 1. 从cookie中获取auth_key
	authKeyCookie, err := r.Cookie(consts.CoDoCookieAuthKey)
	if err != nil || authKeyCookie == nil {
		return nil, AuthFailed
	}
	// 2 .解析auth_key
	var userClaims *UserClaims
	if x.TokenVerify {
		userClaims, err = x.DecodeAuthToken(authKeyCookie.Value)
	} else {
		// 否则，跳过签名验证
		userClaims, err = x.DecodeAuthTokenWithoutSignature(authKeyCookie.Value)
	}
	if err != nil || userClaims.Data.UserID == "" {
		return nil, AuthFailed
	}
	return userClaims, nil

}

func NewCasbinCheckMiddleware(redis *redis.Client, userGroup *biz.UserGroupUseCase, role *biz.RoleUseCase, cluster *biz.ClusterUseCase, logger log.Logger) *CasbinCheckMiddleware {
	return &CasbinCheckMiddleware{
		TokenVerify: false,
		redis:       redis,
		userGroup:   userGroup,
		role:        role,
		cluster:     cluster,
		contextKeys: struct {
			isSuperuser         string
			userID              string
			userName            string
			aclList             string
			clusterRoles        string
			clusterRoleBindings string
			clientIP            string
			traceID             string
		}{
			isSuperuser:         consts.ContextSuperUserKey,
			userID:              consts.ContextUserIDKey,
			userName:            consts.ContextUserNameKey,
			aclList:             consts.ContextAclListKey,
			clusterRoles:        consts.ContextClusterRolesKey,
			clusterRoleBindings: consts.ContextClusterRoleBindingsKey,
			clientIP:            consts.ContextClientIPKey,
			traceID:             consts.ContextTraceIDKey,
		},
		logger: log.NewHelper(log.With(logger, "module", "middleware/casbin")),
	}
}

// getRemoteIP 获取客户端IP
func getClientIPFromRequest(r *http.Request) string {
	// 1. 从 X-Forwarded-For 获取
	if forwardedFor := r.Header.Get("X-Forwarded-For"); forwardedFor != "" {
		// X-Forwarded-For 可能包含多个 IP，第一个是客户端真实 IP
		ips := strings.Split(forwardedFor, ",")
		if len(ips) > 0 {
			return strings.TrimSpace(ips[0])
		}
	}
	// 2. 从 X-Real-IP 获取
	if ip := r.Header.Get("X-Real-IP"); ip != "" {
		return ip
	}

	// 3. 从 RemoteAddr 获取
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		// 如果分割失败，说明可能没有端口号，直接返回整个地址
		return r.RemoteAddr
	}
	return ip
}

type TextMapCarrierWrapper struct {
	Header http.Header
}

func (x *TextMapCarrierWrapper) Get(key string) string {
	return x.Header.Get(key)
}

func (x *TextMapCarrierWrapper) Set(key, value string) {
	x.Header.Set(key, value)
}

func (x *TextMapCarrierWrapper) Keys() []string {
	keys := make([]string, 0, len(x.Header))
	for key := range x.Header {
		keys = append(keys, key)
	}
	return keys
}

func (x *CasbinCheckMiddleware) WSServer() websocket.WSMiddlewareFunc {
	pg := propagation.TraceContext{}
	return func(handleFunc websocket.WSPreHandleFunc) websocket.WSPreHandleFunc {
		return func(ctx context.Context, r *http.Request) error {
			wrapper := &TextMapCarrierWrapper{Header: r.Header}
			traceCtx := pg.Extract(ctx, wrapper)
			// 获取 Trace ID
			spanContext := trace.SpanContextFromContext(traceCtx)
			traceID := spanContext.TraceID().String()

			// 1. 验证用户登录
			userClaims, err := x.CodoLogin(r)
			if err != nil {
				return AuthFailed
			}

			// 2. 构建上下文值
			ctx = x.buildContext(ctx, userClaims)

			// 3. 客户端IP
			clientIP := getClientIPFromRequest(r)
			ctx = x.buildClientIPContext(ctx, clientIP)

			// 4. Trace ID
			ctx = x.buildTraceIDContext(ctx, traceID)

			// 5. 处理超级管理员
			if userClaims.Data.IsSuperuser {
				return handleFunc(ctx, r)
			}

			// 6. 获取用户acl列表
			//ctx = x.buildContextWithAclList(ctx, userClaims.Data.UserID)

			// 7. 获取用户集群角色
			ctx = x.BuildContextWithClusterRoles(ctx, userClaims.Data.UserID)

			// 8. 继续执行后续的中间件和服务
			return handleFunc(ctx, r)
		}

	}
}

func (x *CasbinCheckMiddleware) buildContext(ctx context.Context, claims *UserClaims) context.Context {
	ctx = context.WithValue(ctx, x.contextKeys.isSuperuser, claims.Data.IsSuperuser)
	ctx = context.WithValue(ctx, x.contextKeys.userID, claims.Data.UserID)
	ctx = context.WithValue(ctx, x.contextKeys.userName, claims.Data.Username)
	return ctx
}

func (x *CasbinCheckMiddleware) buildClientIPContext(ctx context.Context, clientIP string) context.Context {
	ctx = context.WithValue(ctx, x.contextKeys.clientIP, clientIP)
	return ctx
}

func (x *CasbinCheckMiddleware) buildTraceIDContext(ctx context.Context, traceID string) context.Context {
	ctx = context.WithValue(ctx, x.contextKeys.traceID, traceID)
	return ctx
}

func (x *CasbinCheckMiddleware) buildContextWithAclList(ctx context.Context, userID string) context.Context {
	userGroupIDs, err := x.GetUserGroupIDS(ctx, userID)
	if err != nil {
		x.logger.WithContext(ctx).Errorf("查询用户组失败: %v", err)
		return ctx
	}

	allAccessList, err := x.GetAclList(ctx, userGroupIDs)
	if err != nil {
		x.logger.WithContext(ctx).Errorf("获取用户权限失败: %v", err)
		return ctx
	}

	return context.WithValue(ctx, x.contextKeys.aclList, allAccessList)
}

func (x *CasbinCheckMiddleware) BuildContextWithClusterRoles(ctx context.Context, userID string) context.Context {
	userGroupIDs, err := x.GetUserGroupIDS(ctx, userID)
	if err != nil {
		x.logger.WithContext(ctx).Errorf("查询用户所属用户组失败: %v", err)
		return ctx
	}

	clusterRoleBindings, clusterRoles, err := x.GetClusterRoleAndBindings(ctx, userGroupIDs)
	if err != nil {
		x.logger.WithContext(ctx).Errorf("获取用户集群角色失败: %v", err)
		return ctx
	}
	ctx = context.WithValue(ctx, x.contextKeys.clusterRoleBindings, clusterRoleBindings)
	ctx = context.WithValue(ctx, x.contextKeys.clusterRoles, clusterRoles)
	return ctx
}

func (x *CasbinCheckMiddleware) Server() kmiddleware.Middleware {
	return func(handler kmiddleware.Handler) kmiddleware.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			// 1. 提取Http请求
			httpTr, err := extraKratosHTTPTransport(ctx)
			if err != nil {
				return nil, err
			}

			// 2. 验证用户登录
			userClaims, err := x.CodoLogin(httpTr.Request())
			if err != nil {
				return nil, AuthFailed
			}

			// 3. 构建上下文值
			ctx = x.buildContext(ctx, userClaims)

			// 4. 处理超级管理员
			if userClaims.Data.IsSuperuser {
				return handler(ctx, req)
			}

			// 5. 获取用户acl列表
			//ctx = x.buildContextWithAclList(ctx, userClaims.Data.UserID)

			// 6. 获取用户ClusterRoles
			ctx = x.BuildContextWithClusterRoles(ctx, userClaims.Data.UserID)

			// 7. 继续执行后续的中间件和服务
			return handler(ctx, req)
		}
	}
}

func (x *CasbinCheckMiddleware) GetAclList(ctx context.Context, userGroupIDS []int) ([]string, error) {
	aclList := make([]string, 0)
	for _, userGroupID := range userGroupIDS {
		bindings, err := x.userGroup.RoleBindingRepo.List(ctx,
			&biz.ListRoleBindingRequest{
				RoleBindingCommonParams: biz.RoleBindingCommonParams{
					UserGroupID: uint32(userGroupID),
				},
				ListAll: true,
			})
		if err != nil {
			continue
		}
		if len(bindings) == 0 {
			continue
		}
		for _, binding := range bindings {
			roleItem, err := x.role.GetRole(ctx, binding.RoleID)
			if err != nil {
				continue
			}
			yamlStr := roleItem.YamlStr
			clusterID := binding.ClusterID
			clusterRole, err := utils.ParseK8sClusterRoleYAML(yamlStr)
			if err != nil {
				continue
			}
			clusterObj, err := x.cluster.GetClusterByID(ctx, clusterID)
			if err != nil {
				continue
			}
			accessList, err := utils.AggregateToACL(clusterRole, clusterObj.Name, binding.Namespace)
			if err != nil {
				continue
			}
			aclList = append(aclList, accessList...)
		}
	}
	return aclList, nil
}

func (x *CasbinCheckMiddleware) GetClusterRoleAndBindings(ctx context.Context, userGroupIDS []int) ([]map[string][]string, []rbacv1.ClusterRole, error) {
	clusterRoles := make([]rbacv1.ClusterRole, 0)
	clusterBindings := make([]map[string][]string, 0)
	for _, userGroupID := range userGroupIDS {
		bindings, err := x.userGroup.RoleBindingRepo.List(ctx,
			&biz.ListRoleBindingRequest{
				RoleBindingCommonParams: biz.RoleBindingCommonParams{
					UserGroupID: uint32(userGroupID),
				},
				ListAll: true,
			})
		if err != nil {
			continue
		}
		if len(bindings) == 0 {
			continue
		}
		// 为每个用户组创建一个集群绑定映射
		clusterBinding := make(map[string][]string)
		for _, binding := range bindings {
			if binding.Namespace == "all" || binding.Namespace == "" {
				// 如果没有指定命名空间，则默认绑定所有命名空间
				binding.Namespace = "*"
			}
			// 获取角色信息
			roleItem, err := x.role.GetRole(ctx, binding.RoleID)
			if err != nil {
				continue
			}

			// 解析集群角色
			clusterRole, err := utils.ParseK8sClusterRoleYAML(roleItem.YamlStr)
			if err != nil {
				continue
			}
			clusterRoles = append(clusterRoles, *clusterRole)

			// 获取集群信息
			clusterObj, err := x.cluster.GetClusterByID(ctx, binding.ClusterID)
			if err != nil {
				continue
			}
			// 保存集群和命名空间的绑定关系
			if namespaces, ok := clusterBinding[clusterObj.Name]; !ok {
				clusterBinding[clusterObj.Name] = []string{binding.Namespace}
			} else if binding.Namespace != "" && !utils.Contains(namespaces, binding.Namespace) {
				clusterBinding[clusterObj.Name] = append(namespaces, binding.Namespace)
			}
		}
		// 只有当有绑定关系时才添加到结果中
		if len(clusterBinding) > 0 {
			clusterBindings = append(clusterBindings, clusterBinding)
		}
	}
	return clusterBindings, clusterRoles, nil
}
