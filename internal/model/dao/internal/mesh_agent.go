// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// MeshAgentDao is the data access object for table mesh_agent.
type MeshAgentDao struct {
	table   string           // table is the underlying table name of the DAO.
	group   string           // group is the database configuration group name of current DAO.
	columns MeshAgentColumns // columns contains all the column names of Table for convenient usage.
}

// MeshAgentColumns defines and stores column names for table mesh_agent.
type MeshAgentColumns struct {
	Id        string //
	CreatedAt string //
	UpdatedAt string //
	DeletedAt string //
	Name      string //
	Ip        string //
	UpdateBy  string //
	AgentId   string //
}

// meshAgentColumns holds the columns for table mesh_agent.
var meshAgentColumns = MeshAgentColumns{
	Id:        "id",
	CreatedAt: "created_at",
	UpdatedAt: "updated_at",
	DeletedAt: "deleted_at",
	Name:      "name",
	Ip:        "ip",
	UpdateBy:  "update_by",
	AgentId:   "agent_id",
}

// NewMeshAgentDao creates and returns a new DAO object for table data access.
func NewMeshAgentDao() *MeshAgentDao {
	return &MeshAgentDao{
		group:   "default",
		table:   "mesh_agent",
		columns: meshAgentColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *MeshAgentDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *MeshAgentDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *MeshAgentDao) Columns() MeshAgentColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *MeshAgentDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *MeshAgentDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *MeshAgentDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
