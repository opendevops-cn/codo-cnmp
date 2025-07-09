// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// User is the golang structure for table user.
type User struct {
	Id        uint64      `json:"id"         orm:"id"         ` //
	CreatedAt *gtime.Time `json:"created_at" orm:"created_at" ` //
	UpdatedAt *gtime.Time `json:"updated_at" orm:"updated_at" ` //
	DeletedAt *gtime.Time `json:"deleted_at" orm:"deleted_at" ` //
	Username  string      `json:"username"   orm:"username"   ` //
	Nickname  string      `json:"nickname"   orm:"nickname"   ` //
	UserId    uint64      `json:"user_id"    orm:"user_id"    ` //
	Email     string      `json:"email"      orm:"email"      ` //
	Source    string      `json:"source"     orm:"source"     ` //
}
