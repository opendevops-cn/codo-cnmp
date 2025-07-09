// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// UserUserGroupRel is the golang structure for table user_user_group_rel.
type UserUserGroupRel struct {
	Id          uint64      `json:"id"            orm:"id"            ` //
	CreatedAt   *gtime.Time `json:"created_at"    orm:"created_at"    ` //
	UpdatedAt   *gtime.Time `json:"updated_at"    orm:"updated_at"    ` //
	DeletedAt   *gtime.Time `json:"deleted_at"    orm:"deleted_at"    ` //
	UserId      uint64      `json:"user_id"       orm:"user_id"       ` //
	UserGroupId uint64      `json:"user_group_id" orm:"user_group_id" ` //
}
