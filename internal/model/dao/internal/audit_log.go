// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// AuditLogDao is the data access object for table audit_log.
type AuditLogDao struct {
	table   string          // table is the underlying table name of the DAO.
	group   string          // group is the database configuration group name of current DAO.
	columns AuditLogColumns // columns contains all the column names of Table for convenient usage.
}

// AuditLogColumns defines and stores column names for table audit_log.
type AuditLogColumns struct {
	Id            string // 主键ID
	Username      string // 操作用户名
	ClientIp      string // 客户端IP
	Module        string // 模块
	Cluster       string // 集群名称
	Namespace     string // 命名空间
	ResourceType  string // 对象类型
	ResourceName  string // 对象名称
	Action        string // 操作类型
	RequestPath   string // 请求路径
	RequestBody   string // 请求内容
	ResponseBody  string // 响应内容
	Status        string // 操作状态 0:成功 1:失败
	Duration      string // 操作耗时(ms)
	OperationTime string // 操作时间
	CreatedAt     string // 创建时间
	UpdatedAt     string // 更新时间
	DeletedAt     string // 删除时间
	TraceId       string // traceID
}

// auditLogColumns holds the columns for table audit_log.
var auditLogColumns = AuditLogColumns{
	Id:            "id",
	Username:      "username",
	ClientIp:      "client_ip",
	Module:        "module",
	Cluster:       "cluster",
	Namespace:     "namespace",
	ResourceType:  "resource_type",
	ResourceName:  "resource_name",
	Action:        "action",
	RequestPath:   "request_path",
	RequestBody:   "request_body",
	ResponseBody:  "response_body",
	Status:        "status",
	Duration:      "duration",
	OperationTime: "operation_time",
	CreatedAt:     "created_at",
	UpdatedAt:     "updated_at",
	DeletedAt:     "deleted_at",
	TraceId:       "trace_id",
}

// NewAuditLogDao creates and returns a new DAO object for table data access.
func NewAuditLogDao() *AuditLogDao {
	return &AuditLogDao{
		group:   "default",
		table:   "audit_log",
		columns: auditLogColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *AuditLogDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *AuditLogDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *AuditLogDao) Columns() AuditLogColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *AuditLogDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *AuditLogDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *AuditLogDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
