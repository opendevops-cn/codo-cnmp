// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// UserUserGroupRelDao is the data access object for table user_user_group_rel.
type UserUserGroupRelDao struct {
	table   string                  // table is the underlying table name of the DAO.
	group   string                  // group is the database configuration group name of current DAO.
	columns UserUserGroupRelColumns // columns contains all the column names of Table for convenient usage.
}

// UserUserGroupRelColumns defines and stores column names for table user_user_group_rel.
type UserUserGroupRelColumns struct {
	Id          string //
	CreatedAt   string //
	UpdatedAt   string //
	DeletedAt   string //
	UserId      string //
	UserGroupId string //
}

// userUserGroupRelColumns holds the columns for table user_user_group_rel.
var userUserGroupRelColumns = UserUserGroupRelColumns{
	Id:          "id",
	CreatedAt:   "created_at",
	UpdatedAt:   "updated_at",
	DeletedAt:   "deleted_at",
	UserId:      "user_id",
	UserGroupId: "user_group_id",
}

// NewUserUserGroupRelDao creates and returns a new DAO object for table data access.
func NewUserUserGroupRelDao() *UserUserGroupRelDao {
	return &UserUserGroupRelDao{
		group:   "default",
		table:   "user_user_group_rel",
		columns: userUserGroupRelColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *UserUserGroupRelDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *UserUserGroupRelDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *UserUserGroupRelDao) Columns() UserUserGroupRelColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *UserUserGroupRelDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *UserUserGroupRelDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *UserUserGroupRelDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
