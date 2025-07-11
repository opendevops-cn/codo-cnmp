// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// UserGroupDao is the data access object for table user_group.
type UserGroupDao struct {
	table   string           // table is the underlying table name of the DAO.
	group   string           // group is the database configuration group name of current DAO.
	columns UserGroupColumns // columns contains all the column names of Table for convenient usage.
}

// UserGroupColumns defines and stores column names for table user_group.
type UserGroupColumns struct {
	Id          string //
	CreatedAt   string //
	UpdatedAt   string //
	DeletedAt   string //
	Name        string //
	UserGroupId string //
}

// userGroupColumns holds the columns for table user_group.
var userGroupColumns = UserGroupColumns{
	Id:          "id",
	CreatedAt:   "created_at",
	UpdatedAt:   "updated_at",
	DeletedAt:   "deleted_at",
	Name:        "name",
	UserGroupId: "user_group_id",
}

// NewUserGroupDao creates and returns a new DAO object for table data access.
func NewUserGroupDao() *UserGroupDao {
	return &UserGroupDao{
		group:   "default",
		table:   "user_group",
		columns: userGroupColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *UserGroupDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *UserGroupDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *UserGroupDao) Columns() UserGroupColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *UserGroupDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *UserGroupDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *UserGroupDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
