// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// GrantedUserGroup is the golang structure of table granted_user_group for DAO operations like Where/Data.
type GrantedUserGroup struct {
	g.Meta      `orm:"table:granted_user_group, do:true"`
	Id          interface{} //
	CreatedAt   *gtime.Time //
	UpdatedAt   *gtime.Time //
	DeletedAt   *gtime.Time //
	Name        interface{} //
	UserGroupId interface{} //
	RoleDetail  interface{} //
}
