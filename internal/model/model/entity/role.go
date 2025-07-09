// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// Role is the golang structure for table role.
type Role struct {
	Id          uint64      `json:"id"          orm:"id"          ` //
	CreatedAt   *gtime.Time `json:"created_at"  orm:"created_at"  ` //
	UpdatedAt   *gtime.Time `json:"updated_at"  orm:"updated_at"  ` //
	DeletedAt   *gtime.Time `json:"deleted_at"  orm:"deleted_at"  ` //
	Name        string      `json:"name"        orm:"name"        ` //
	RoleType    int         `json:"role_type"   orm:"role_type"   ` //
	IsDefault   int         `json:"is_default"  orm:"is_default"  ` //
	Description string      `json:"description" orm:"description" ` //
	Yaml        string      `json:"yaml"        orm:"yaml"        ` //
	UpdateBy    string      `json:"update_by"   orm:"update_by"   ` //
}
