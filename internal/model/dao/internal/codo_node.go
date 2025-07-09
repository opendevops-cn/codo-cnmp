// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// CodoNodeDao is the data access object for table codo_node.
type CodoNodeDao struct {
	table   string          // table is the underlying table name of the DAO.
	group   string          // group is the database configuration group name of current DAO.
	columns CodoNodeColumns // columns contains all the column names of Table for convenient usage.
}

// CodoNodeColumns defines and stores column names for table codo_node.
type CodoNodeColumns struct {
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
}

// codoNodeColumns holds the columns for table codo_node.
var codoNodeColumns = CodoNodeColumns{
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
}

// NewCodoNodeDao creates and returns a new DAO object for table data access.
func NewCodoNodeDao() *CodoNodeDao {
	return &CodoNodeDao{
		group:   "default",
		table:   "codo_node",
		columns: codoNodeColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *CodoNodeDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *CodoNodeDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *CodoNodeDao) Columns() CodoNodeColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *CodoNodeDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *CodoNodeDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *CodoNodeDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
