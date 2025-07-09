// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// GameServer is the golang structure of table game_server for DAO operations like Where/Data.
type GameServer struct {
	g.Meta            `orm:"table:game_server, do:true"`
	Id                interface{} // 自增主键
	EntityNum         interface{} // 实例对象数
	ServerName        interface{} // 进程ID
	OnlineNum         interface{} // 在线用户数
	LockEntityStatus  interface{} // entity锁定状态
	LockLbStatus      interface{} // lb锁定状态
	ServerType        interface{} // 进程类型
	Server            interface{} // 进程名称
	ClusterName       interface{} // 集群名
	Namespace         interface{} // 命名空间
	Pod               interface{} // pod名称
	ServerVersion     interface{} // 服务版本
	CodeVersionGame   interface{} // 游戏代码版本
	CodeVersionConfig interface{} // 配置代码版本
	CreatedAt         *gtime.Time // 创建时间
	UpdatedAt         *gtime.Time // 更新时间
	DeletedAt         *gtime.Time // 删除时间
	CodeVersionScript interface{} //
	Workload          interface{} // 工作负载
	WorkloadType      interface{} // 工作负载类型
	ServerTypeDesc    interface{} // 进程类型中文名
	BigArea           interface{} // 大区编号
	GameAppId         interface{} // 游戏应用编号
}
