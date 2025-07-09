// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// UserFollow is the golang structure for table user_follow.
type UserFollow struct {
	Id          uint64      `json:"id"           orm:"id"           ` // 主键ID
	UserId      uint64      `json:"user_id"      orm:"user_id"      ` // 用户ID
	FollowType  int         `json:"follow_type"  orm:"follow_type"  ` // 关注类型
	FollowValue string      `json:"follow_value" orm:"follow_value" ` // 关注对象
	CreatedAt   *gtime.Time `json:"created_at"   orm:"created_at"   ` // 创建时间
	UpdatedAt   *gtime.Time `json:"updated_at"   orm:"updated_at"   ` // 更新时间
	DeletedAt   *gtime.Time `json:"deleted_at"   orm:"deleted_at"   ` // 删除时间
	ClusterName string      `json:"cluster_name" orm:"cluster_name" ` //
}
