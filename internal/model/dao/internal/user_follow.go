// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// UserFollowDao is the data access object for table user_follow.
type UserFollowDao struct {
	table   string            // table is the underlying table name of the DAO.
	group   string            // group is the database configuration group name of current DAO.
	columns UserFollowColumns // columns contains all the column names of Table for convenient usage.
}

// UserFollowColumns defines and stores column names for table user_follow.
type UserFollowColumns struct {
	Id          string // 主键ID
	UserId      string // 用户ID
	FollowType  string // 关注类型
	FollowValue string // 关注对象
	CreatedAt   string // 创建时间
	UpdatedAt   string // 更新时间
	DeletedAt   string // 删除时间
	ClusterName string //
}

// userFollowColumns holds the columns for table user_follow.
var userFollowColumns = UserFollowColumns{
	Id:          "id",
	UserId:      "user_id",
	FollowType:  "follow_type",
	FollowValue: "follow_value",
	CreatedAt:   "created_at",
	UpdatedAt:   "updated_at",
	DeletedAt:   "deleted_at",
	ClusterName: "cluster_name",
}

// NewUserFollowDao creates and returns a new DAO object for table data access.
func NewUserFollowDao() *UserFollowDao {
	return &UserFollowDao{
		group:   "default",
		table:   "user_follow",
		columns: userFollowColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *UserFollowDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *UserFollowDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *UserFollowDao) Columns() UserFollowColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *UserFollowDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *UserFollowDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *UserFollowDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
