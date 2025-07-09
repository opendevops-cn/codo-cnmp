// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// AuditLog is the golang structure for table audit_log.
type AuditLog struct {
	Id            uint64      `json:"id"             orm:"id"             ` // 主键ID
	Username      string      `json:"username"       orm:"username"       ` // 操作用户名
	ClientIp      string      `json:"client_ip"      orm:"client_ip"      ` // 客户端IP
	Module        string      `json:"module"         orm:"module"         ` // 模块
	Cluster       string      `json:"cluster"        orm:"cluster"        ` // 集群名称
	Namespace     string      `json:"namespace"      orm:"namespace"      ` // 命名空间
	ResourceType  string      `json:"resource_type"  orm:"resource_type"  ` // 对象类型
	ResourceName  string      `json:"resource_name"  orm:"resource_name"  ` // 对象名称
	Action        string      `json:"action"         orm:"action"         ` // 操作类型
	RequestPath   string      `json:"request_path"   orm:"request_path"   ` // 请求路径
	RequestBody   string      `json:"request_body"   orm:"request_body"   ` // 请求内容
	ResponseBody  string      `json:"response_body"  orm:"response_body"  ` // 响应内容
	Status        int         `json:"status"         orm:"status"         ` // 操作状态 0:成功 1:失败
	Duration      string      `json:"duration"       orm:"duration"       ` // 操作耗时(ms)
	OperationTime *gtime.Time `json:"operation_time" orm:"operation_time" ` // 操作时间
	CreatedAt     *gtime.Time `json:"created_at"     orm:"created_at"     ` // 创建时间
	UpdatedAt     *gtime.Time `json:"updated_at"     orm:"updated_at"     ` // 更新时间
	DeletedAt     *gtime.Time `json:"deleted_at"     orm:"deleted_at"     ` // 删除时间
	TraceId       string      `json:"trace_id"       orm:"trace_id"       ` // traceID
}
