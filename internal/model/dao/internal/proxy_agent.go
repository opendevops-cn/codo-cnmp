// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// ProxyAgentDao is the data access object for table proxy_agent.
type ProxyAgentDao struct {
	table   string            // table is the underlying table name of the DAO.
	group   string            // group is the database configuration group name of current DAO.
	columns ProxyAgentColumns // columns contains all the column names of Table for convenient usage.
}

// ProxyAgentColumns defines and stores column names for table proxy_agent.
type ProxyAgentColumns struct {
	Id        string //
	CreatedAt string //
	UpdatedAt string //
	DeletedAt string //
	Name      string //
	AgentId   string //
}

// proxyAgentColumns holds the columns for table proxy_agent.
var proxyAgentColumns = ProxyAgentColumns{
	Id:        "id",
	CreatedAt: "created_at",
	UpdatedAt: "updated_at",
	DeletedAt: "deleted_at",
	Name:      "name",
	AgentId:   "agent_id",
}

// NewProxyAgentDao creates and returns a new DAO object for table data access.
func NewProxyAgentDao() *ProxyAgentDao {
	return &ProxyAgentDao{
		group:   "default",
		table:   "proxy_agent",
		columns: proxyAgentColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *ProxyAgentDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *ProxyAgentDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *ProxyAgentDao) Columns() ProxyAgentColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *ProxyAgentDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *ProxyAgentDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *ProxyAgentDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
