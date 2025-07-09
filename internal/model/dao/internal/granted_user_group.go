// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// GrantedUserGroupDao is the data access object for table granted_user_group.
type GrantedUserGroupDao struct {
	table   string                  // table is the underlying table name of the DAO.
	group   string                  // group is the database configuration group name of current DAO.
	columns GrantedUserGroupColumns // columns contains all the column names of Table for convenient usage.
}

// GrantedUserGroupColumns defines and stores column names for table granted_user_group.
type GrantedUserGroupColumns struct {
	Id          string //
	CreatedAt   string //
	UpdatedAt   string //
	DeletedAt   string //
	Name        string //
	UserGroupId string //
	RoleDetail  string //
}

// grantedUserGroupColumns holds the columns for table granted_user_group.
var grantedUserGroupColumns = GrantedUserGroupColumns{
	Id:          "id",
	CreatedAt:   "created_at",
	UpdatedAt:   "updated_at",
	DeletedAt:   "deleted_at",
	Name:        "name",
	UserGroupId: "user_group_id",
	RoleDetail:  "role_detail",
}

// NewGrantedUserGroupDao creates and returns a new DAO object for table data access.
func NewGrantedUserGroupDao() *GrantedUserGroupDao {
	return &GrantedUserGroupDao{
		group:   "default",
		table:   "granted_user_group",
		columns: grantedUserGroupColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *GrantedUserGroupDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *GrantedUserGroupDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *GrantedUserGroupDao) Columns() GrantedUserGroupColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *GrantedUserGroupDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *GrantedUserGroupDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *GrantedUserGroupDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
