// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// GameServerDao is the data access object for table game_server.
type GameServerDao struct {
	table   string            // table is the underlying table name of the DAO.
	group   string            // group is the database configuration group name of current DAO.
	columns GameServerColumns // columns contains all the column names of Table for convenient usage.
}

// GameServerColumns defines and stores column names for table game_server.
type GameServerColumns struct {
	Id                string // 自增主键
	EntityNum         string // 实例对象数
	ServerName        string // 进程ID
	OnlineNum         string // 在线用户数
	LockEntityStatus  string // entity锁定状态
	LockLbStatus      string // lb锁定状态
	ServerType        string // 进程类型
	Server            string // 进程名称
	ClusterName       string // 集群名
	Namespace         string // 命名空间
	Pod               string // pod名称
	ServerVersion     string // 服务版本
	CodeVersionGame   string // 游戏代码版本
	CodeVersionConfig string // 配置代码版本
	CreatedAt         string // 创建时间
	UpdatedAt         string // 更新时间
	DeletedAt         string // 删除时间
	CodeVersionScript string //
	Workload          string // 工作负载
	WorkloadType      string // 工作负载类型
	ServerTypeDesc    string // 进程类型中文名
	BigArea           string // 大区编号
	GameAppId         string // 游戏应用编号
}

// gameServerColumns holds the columns for table game_server.
var gameServerColumns = GameServerColumns{
	Id:                "id",
	EntityNum:         "entity_num",
	ServerName:        "server_name",
	OnlineNum:         "online_num",
	LockEntityStatus:  "lock_entity_status",
	LockLbStatus:      "lock_lb_status",
	ServerType:        "server_type",
	Server:            "server",
	ClusterName:       "cluster_name",
	Namespace:         "namespace",
	Pod:               "pod",
	ServerVersion:     "server_version",
	CodeVersionGame:   "code_version_game",
	CodeVersionConfig: "code_version_config",
	CreatedAt:         "created_at",
	UpdatedAt:         "updated_at",
	DeletedAt:         "deleted_at",
	CodeVersionScript: "code_version_script",
	Workload:          "workload",
	WorkloadType:      "workload_type",
	ServerTypeDesc:    "server_type_desc",
	BigArea:           "big_area",
	GameAppId:         "game_app_id",
}

// NewGameServerDao creates and returns a new DAO object for table data access.
func NewGameServerDao() *GameServerDao {
	return &GameServerDao{
		group:   "default",
		table:   "game_server",
		columns: gameServerColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *GameServerDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *GameServerDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *GameServerDao) Columns() GameServerColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *GameServerDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *GameServerDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *GameServerDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
