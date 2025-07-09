// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// UserGroupClusterRoleRef is the golang structure of table user_group_cluster_role_ref for DAO operations like Where/Data.
type UserGroupClusterRoleRef struct {
	g.Meta      `orm:"table:user_group_cluster_role_ref, do:true"`
	Id          interface{} //
	CreatedAt   *gtime.Time //
	UpdatedAt   *gtime.Time //
	DeletedAt   *gtime.Time //
	UserGroupId interface{} //
	ClusterId   interface{} //
	RoleId      interface{} //
	Namespace   interface{} //
}
