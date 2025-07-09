// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// RoleBindingDao is the data access object for table role_binding.
type RoleBindingDao struct {
	table   string             // table is the underlying table name of the DAO.
	group   string             // group is the database configuration group name of current DAO.
	columns RoleBindingColumns // columns contains all the column names of Table for convenient usage.
}

// RoleBindingColumns defines and stores column names for table role_binding.
type RoleBindingColumns struct {
	Id          string //
	CreatedAt   string //
	UpdatedAt   string //
	DeletedAt   string //
	UserGroupId string //
	ClusterId   string //
	RoleId      string //
	Namespace   string //
}

// roleBindingColumns holds the columns for table role_binding.
var roleBindingColumns = RoleBindingColumns{
	Id:          "id",
	CreatedAt:   "created_at",
	UpdatedAt:   "updated_at",
	DeletedAt:   "deleted_at",
	UserGroupId: "user_group_id",
	ClusterId:   "cluster_id",
	RoleId:      "role_id",
	Namespace:   "namespace",
}

// NewRoleBindingDao creates and returns a new DAO object for table data access.
func NewRoleBindingDao() *RoleBindingDao {
	return &RoleBindingDao{
		group:   "default",
		table:   "role_binding",
		columns: roleBindingColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *RoleBindingDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *RoleBindingDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *RoleBindingDao) Columns() RoleBindingColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *RoleBindingDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *RoleBindingDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *RoleBindingDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
