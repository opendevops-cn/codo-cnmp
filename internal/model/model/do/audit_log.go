// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// AuditLog is the golang structure of table audit_log for DAO operations like Where/Data.
type AuditLog struct {
	g.Meta        `orm:"table:audit_log, do:true"`
	Id            interface{} // 主键ID
	Username      interface{} // 操作用户名
	ClientIp      interface{} // 客户端IP
	Module        interface{} // 模块
	Cluster       interface{} // 集群名称
	Namespace     interface{} // 命名空间
	ResourceType  interface{} // 对象类型
	ResourceName  interface{} // 对象名称
	Action        interface{} // 操作类型
	RequestPath   interface{} // 请求路径
	RequestBody   interface{} // 请求内容
	ResponseBody  interface{} // 响应内容
	Status        interface{} // 操作状态 0:成功 1:失败
	Duration      interface{} // 操作耗时(ms)
	OperationTime *gtime.Time // 操作时间
	CreatedAt     *gtime.Time // 创建时间
	UpdatedAt     *gtime.Time // 更新时间
	DeletedAt     *gtime.Time // 删除时间
	TraceId       interface{} // traceID
}
