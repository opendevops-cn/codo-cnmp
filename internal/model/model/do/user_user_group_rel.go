// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// UserUserGroupRel is the golang structure of table user_user_group_rel for DAO operations like Where/Data.
type UserUserGroupRel struct {
	g.Meta      `orm:"table:user_user_group_rel, do:true"`
	Id          interface{} //
	CreatedAt   *gtime.Time //
	UpdatedAt   *gtime.Time //
	DeletedAt   *gtime.Time //
	UserId      interface{} //
	UserGroupId interface{} //
}
