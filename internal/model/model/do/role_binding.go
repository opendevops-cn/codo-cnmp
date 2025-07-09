// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// RoleBinding is the golang structure of table role_binding for DAO operations like Where/Data.
type RoleBinding struct {
	g.Meta      `orm:"table:role_binding, do:true"`
	Id          interface{} //
	CreatedAt   *gtime.Time //
	UpdatedAt   *gtime.Time //
	DeletedAt   *gtime.Time //
	UserGroupId interface{} //
	ClusterId   interface{} //
	RoleId      interface{} //
	Namespace   interface{} //
}
