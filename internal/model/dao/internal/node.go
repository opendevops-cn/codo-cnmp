// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// NodeDao is the data access object for table node.
type NodeDao struct {
	table   string      // table is the underlying table name of the DAO.
	group   string      // group is the database configuration group name of current DAO.
	columns NodeColumns // columns contains all the column names of Table for convenient usage.
}

// NodeColumns defines and stores column names for table node.
type NodeColumns struct {
	Id                string //
	Name              string //
	Conditions        string //
	Capacity          string //
	Allocatable       string //
	Addresses         string //
	CreationTimestamp string //
	CreatedAt         string //
	UpdatedAt         string //
	DeletedAt         string //
	ClusterId         string //
	CpuUsage          string //
	MemoryUsage       string //
	Status            string //
	Labels            string //
	Annotations       string //
	NodeInfo          string //
	Roles             string //
	Uid               string //
	ResourceVersion   string //
	HealthState       string //
	Spec              string //
}

// nodeColumns holds the columns for table node.
var nodeColumns = NodeColumns{
	Id:                "id",
	Name:              "name",
	Conditions:        "conditions",
	Capacity:          "capacity",
	Allocatable:       "allocatable",
	Addresses:         "addresses",
	CreationTimestamp: "creation_timestamp",
	CreatedAt:         "created_at",
	UpdatedAt:         "updated_at",
	DeletedAt:         "deleted_at",
	ClusterId:         "cluster_id",
	CpuUsage:          "cpu_usage",
	MemoryUsage:       "memory_usage",
	Status:            "status",
	Labels:            "labels",
	Annotations:       "annotations",
	NodeInfo:          "node_info",
	Roles:             "roles",
	Uid:               "uid",
	ResourceVersion:   "resource_version",
	HealthState:       "health_state",
	Spec:              "spec",
}

// NewNodeDao creates and returns a new DAO object for table data access.
func NewNodeDao() *NodeDao {
	return &NodeDao{
		group:   "default",
		table:   "node",
		columns: nodeColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *NodeDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *NodeDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *NodeDao) Columns() NodeColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *NodeDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *NodeDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *NodeDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
