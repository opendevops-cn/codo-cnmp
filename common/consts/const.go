package consts

const (
	CPUUsedThreshold              = 80
	MemoryUsedThreshold           = 80
	UserCacheKey                  = "codo:cnmp:user"
	UserGroupCacheKey             = "codo:cnmp:usergroup"
	ContextSuperUserKey           = "isSuperUser"
	ContextUserIDKey              = "userID"
	ContextUserNameKey            = "userName"
	ContextAclListKey             = "aclList"
	ContextClusterRoleBindingsKey = "clusterRoleBindings"
	ContextClientIPKey            = "clientIP"
	ContextTraceIDKey             = "traceID"
	ContextClusterRolesKey        = "clusterRoles"
	CoDoCookieAuthKey             = "auth_key"
	UserGroupIDCacheKey           = "user:%s:mg_user_group_ids"
	MeshCacheKey                  = "codo:cnmp:mesh"
	AgentPortCacheKey             = "codo:cnmp:agent:%s"
)
