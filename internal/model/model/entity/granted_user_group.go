// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// GrantedUserGroup is the golang structure for table granted_user_group.
type GrantedUserGroup struct {
	Id          uint64      `json:"id"            orm:"id"            ` //
	CreatedAt   *gtime.Time `json:"created_at"    orm:"created_at"    ` //
	UpdatedAt   *gtime.Time `json:"updated_at"    orm:"updated_at"    ` //
	DeletedAt   *gtime.Time `json:"deleted_at"    orm:"deleted_at"    ` //
	Name        string      `json:"name"          orm:"name"          ` //
	UserGroupId string      `json:"user_group_id" orm:"user_group_id" ` //
	RoleDetail  string      `json:"role_detail"   orm:"role_detail"   ` //
}
