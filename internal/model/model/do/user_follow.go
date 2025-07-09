// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// UserFollow is the golang structure of table user_follow for DAO operations like Where/Data.
type UserFollow struct {
	g.Meta      `orm:"table:user_follow, do:true"`
	Id          interface{} // 主键ID
	UserId      interface{} // 用户ID
	FollowType  interface{} // 关注类型
	FollowValue interface{} // 关注对象
	CreatedAt   *gtime.Time // 创建时间
	UpdatedAt   *gtime.Time // 更新时间
	DeletedAt   *gtime.Time // 删除时间
	ClusterName interface{} //
}
