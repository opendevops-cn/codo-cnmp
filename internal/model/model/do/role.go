// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// Role is the golang structure of table role for DAO operations like Where/Data.
type Role struct {
	g.Meta      `orm:"table:role, do:true"`
	Id          interface{} //
	CreatedAt   *gtime.Time //
	UpdatedAt   *gtime.Time //
	DeletedAt   *gtime.Time //
	Name        interface{} //
	RoleType    interface{} //
	IsDefault   interface{} //
	Description interface{} //
	Yaml        interface{} //
	UpdateBy    interface{} //
}
