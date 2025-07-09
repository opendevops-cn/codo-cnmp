// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// GameServer is the golang structure for table game_server.
type GameServer struct {
	Id                int         `json:"id"                  orm:"id"                  ` // 自增主键
	EntityNum         int         `json:"entity_num"          orm:"entity_num"          ` // 实例对象数
	ServerName        string      `json:"server_name"         orm:"server_name"         ` // 进程ID
	OnlineNum         int         `json:"online_num"          orm:"online_num"          ` // 在线用户数
	LockEntityStatus  int         `json:"lock_entity_status"  orm:"lock_entity_status"  ` // entity锁定状态
	LockLbStatus      int         `json:"lock_lb_status"      orm:"lock_lb_status"      ` // lb锁定状态
	ServerType        string      `json:"server_type"         orm:"server_type"         ` // 进程类型
	Server            string      `json:"server"              orm:"server"              ` // 进程名称
	ClusterName       string      `json:"cluster_name"        orm:"cluster_name"        ` // 集群名
	Namespace         string      `json:"namespace"           orm:"namespace"           ` // 命名空间
	Pod               string      `json:"pod"                 orm:"pod"                 ` // pod名称
	ServerVersion     string      `json:"server_version"      orm:"server_version"      ` // 服务版本
	CodeVersionGame   string      `json:"code_version_game"   orm:"code_version_game"   ` // 游戏代码版本
	CodeVersionConfig string      `json:"code_version_config" orm:"code_version_config" ` // 配置代码版本
	CreatedAt         *gtime.Time `json:"created_at"          orm:"created_at"          ` // 创建时间
	UpdatedAt         *gtime.Time `json:"updated_at"          orm:"updated_at"          ` // 更新时间
	DeletedAt         *gtime.Time `json:"deleted_at"          orm:"deleted_at"          ` // 删除时间
	CodeVersionScript string      `json:"code_version_script" orm:"code_version_script" ` //
	Workload          string      `json:"workload"            orm:"workload"            ` // 工作负载
	WorkloadType      string      `json:"workload_type"       orm:"workload_type"       ` // 工作负载类型
	ServerTypeDesc    string      `json:"server_type_desc"    orm:"server_type_desc"    ` // 进程类型中文名
	BigArea           string      `json:"big_area"            orm:"big_area"            ` // 大区编号
	GameAppId         string      `json:"game_app_id"         orm:"game_app_id"         ` // 游戏应用编号
}
