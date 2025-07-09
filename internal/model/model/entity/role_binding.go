// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// RoleBinding is the golang structure for table role_binding.
type RoleBinding struct {
	Id          uint64      `json:"id"            orm:"id"            ` //
	CreatedAt   *gtime.Time `json:"created_at"    orm:"created_at"    ` //
	UpdatedAt   *gtime.Time `json:"updated_at"    orm:"updated_at"    ` //
	DeletedAt   *gtime.Time `json:"deleted_at"    orm:"deleted_at"    ` //
	UserGroupId uint64      `json:"user_group_id" orm:"user_group_id" ` //
	ClusterId   uint64      `json:"cluster_id"    orm:"cluster_id"    ` //
	RoleId      uint64      `json:"role_id"       orm:"role_id"       ` //
	Namespace   string      `json:"namespace"     orm:"namespace"     ` //
}
