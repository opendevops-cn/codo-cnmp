package data

import (
	"codo-cnmp/internal/biz"
	"codo-cnmp/internal/model/dao"
	"codo-cnmp/internal/model/model/entity"
	"context"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
	"time"
)

type GameServerRepo struct {
	data *Data
	log  *log.Helper
}

// BatchUpdateGameServer 批量更新游戏进程
func (x *GameServerRepo) BatchUpdateGameServer(ctx context.Context, servers []*biz.GameServer) error {
	if len(servers) == 0 {
		return nil
	}

	// 更新耗时
	start := time.Now()
	defer func() {
		x.log.Infof("批量更新游戏服务器耗时: %s", time.Since(start))
	}()

	// 使用事务处理
	return g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		for _, server := range servers {
			// 转换为数据库模型
			vo := x.convertDO2VO(server)

			// 构建更新字段
			updates := g.Map{}
			if vo.Workload != "" {
				updates["workload"] = vo.Workload
			}
			if vo.WorkloadType != "" {
				updates["workload_type"] = vo.WorkloadType
			}
			if vo.EntityNum >= 0 {
				updates["entity_num"] = vo.EntityNum
			}
			if vo.OnlineNum >= 0 {
				updates["online_num"] = vo.OnlineNum
			}
			if vo.ServerVersion != "" {
				updates["server_version"] = vo.ServerVersion
			}
			if vo.CodeVersionGame != "" {
				updates["code_version_game"] = vo.CodeVersionGame
			}
			if vo.CodeVersionConfig != "" {
				updates["code_version_config"] = vo.CodeVersionConfig
			}
			if vo.CodeVersionScript != "" {
				updates["code_version_script"] = vo.CodeVersionScript
			}
			if vo.LockEntityStatus >= 0 {
				updates["lock_entity_status"] = vo.LockEntityStatus
			}
			if vo.LockLbStatus >= 0 {
				updates["lock_lb_status"] = vo.LockLbStatus
			}
			if vo.ServerName != "" {
				updates["server_name"] = vo.ServerName
			}

			// 如果没有需要更新的字段，跳过
			if len(updates) == 0 {
				continue
			}

			// 执行更新
			result, err := dao.GameServer.Ctx(ctx).TX(tx).
				Where(dao.GameServer.Columns().ClusterName, server.ClusterName).
				Where(dao.GameServer.Columns().Namespace, server.Namespace).
				Where(dao.GameServer.Columns().Pod, server.Pod).
				Update(updates)

			if err != nil {
				x.log.WithContext(ctx).Errorf("更新游戏进程失败, 集群: [%s], 命名空间: [%s], pod: [%s], 原因： %v", server.ClusterName, server.Namespace, server.Pod, err)
				return fmt.Errorf("更新游戏进程失败, 集群: [%s], 命名空间: [%s], pod: [%s], 原因： %v", server.ClusterName, server.Namespace, server.Pod, err)
			}

			_, err = result.RowsAffected()
			if err != nil {
				return fmt.Errorf("更新游戏进程失败: %w", err)
			}

		}

		return nil
	})
}

// BatchUpdateGameServerWithChunks 分块批量更新游戏服务器
func (x *GameServerRepo) BatchUpdateGameServerWithChunks(ctx context.Context,
	servers []*biz.GameServer, chunkSize int) error {

	if len(servers) == 0 {
		return nil
	}

	if chunkSize <= 0 {
		chunkSize = 100 // 默认块大小
	}

	var errs []error
	for i := 0; i < len(servers); i += chunkSize {
		end := i + chunkSize
		if end > len(servers) {
			end = len(servers)
		}

		chunk := servers[i:end]
		if err := x.BatchUpdateGameServer(ctx, chunk); err != nil {
			errs = append(errs, fmt.Errorf("chunk %d-%d: %w", i, end, err))
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("batch update errors: %v", errs)
	}

	return nil
}

func (x *GameServerRepo) convertCountQuery(ctx context.Context, req *biz.ListGameServerRequest) *gdb.Model {
	model := dao.GameServer.Ctx(ctx)
	if req.ClusterName != "" {
		model = model.Where(dao.GameServer.Columns().ClusterName, req.ClusterName)
	}
	if req.Namespace != "" {
		model = model.Where(dao.GameServer.Columns().Namespace, req.Namespace)
	}
	if req.ServerType != "" {
		model = model.Where(dao.GameServer.Columns().ServerType, req.ServerType)
	}
	if req.EntityLockStatus > 0 {
		model = model.Where(dao.GameServer.Columns().LockEntityStatus, req.EntityLockStatus)
	}
	if req.LbLockStatus > 0 {
		model = model.Where(dao.GameServer.Columns().LockLbStatus, req.LbLockStatus)
	}
	if req.Keyword != "" {
		model = model.WhereOrLike(dao.GameServer.Columns().ServerName, "%"+req.Keyword+"%").
			WhereOrLike(dao.GameServer.Columns().Pod, "%"+req.Keyword+"%").
			WhereOrLike(dao.GameServer.Columns().ServerVersion, "%"+req.Keyword+"%").
			WhereOrLike(dao.GameServer.Columns().CodeVersionGame, "%"+req.Keyword+"%").
			WhereOrLike(dao.GameServer.Columns().CodeVersionConfig, "%"+req.Keyword+"%").
			WhereOrLike(dao.GameServer.Columns().CodeVersionScript, "%"+req.Keyword+"%")
	}
	if req.ListAll {
		return model
	}
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}
	return model.Page(int(req.Page), int(req.PageSize))
}

func (x *GameServerRepo) convertQuery(ctx context.Context, req *biz.ListGameServerRequest) *gdb.Model {
	model := dao.GameServer.Ctx(ctx)
	if req.ClusterName != "" {
		model = model.Where(dao.GameServer.Columns().ClusterName, req.ClusterName)
	}
	if req.Namespace != "" {
		model = model.Where(dao.GameServer.Columns().Namespace, req.Namespace)
	}
	if req.ServerType != "" {
		model = model.Where(dao.GameServer.Columns().ServerTypeDesc, req.ServerType)
	}
	// 处理 EntityLockStatus 和 LbLockStatus 为布尔类型
	if req.EntityLockStatus > 0 {
		model = model.Where(dao.GameServer.Columns().LockEntityStatus, req.EntityLockStatus)
	}
	if req.LbLockStatus > 0 {
		model = model.Where(dao.GameServer.Columns().LockLbStatus, req.LbLockStatus)
	}

	if req.Keyword != "" {
		model = model.Where("("+
			dao.GameServer.Columns().ServerName+" LIKE ? OR "+
			dao.GameServer.Columns().Pod+" LIKE ? OR "+
			dao.GameServer.Columns().ServerVersion+" LIKE ? OR "+
			dao.GameServer.Columns().CodeVersionGame+" LIKE ? OR "+
			dao.GameServer.Columns().CodeVersionConfig+" LIKE ? OR "+
			dao.GameServer.Columns().CodeVersionScript+" LIKE ?"+
			")",
			"%"+req.Keyword+"%",
			"%"+req.Keyword+"%",
			"%"+req.Keyword+"%",
			"%"+req.Keyword+"%",
			"%"+req.Keyword+"%",
			"%"+req.Keyword+"%",
		)
	}
	if req.ListAll {
		return model
	}
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}
	return model.Page(int(req.Page), int(req.PageSize))
}

func (x *GameServerRepo) convertVO2DO(vo *entity.GameServer) *biz.GameServer {
	return &biz.GameServer{
		Workload:          vo.Workload,
		WorkloadType:      vo.WorkloadType,
		EntityNum:         uint32(vo.EntityNum),
		OnlineNum:         uint32(vo.OnlineNum),
		ServerName:        vo.ServerName,
		Pod:               vo.Pod,
		ServerType:        vo.ServerType,
		ServerVersion:     vo.ServerVersion,
		CodeVersionGame:   vo.CodeVersionGame,
		CodeVersionConfig: vo.CodeVersionConfig,
		CodeVersionScript: vo.CodeVersionScript,
		EntityLockStatus:  uint32(vo.LockEntityStatus),
		LbLockStatus:      uint32(vo.LockLbStatus),
		ClusterName:       vo.ClusterName,
		Namespace:         vo.Namespace,
		ID:                uint32(vo.Id),
		ServerTypeDesc:    vo.ServerTypeDesc,
		GameAppId:         vo.GameAppId,
		BigArea:           vo.BigArea,
	}
}

func (x *GameServerRepo) convertDO2VO(do *biz.GameServer) *entity.GameServer {
	return &entity.GameServer{
		Workload:          do.Workload,
		WorkloadType:      do.WorkloadType,
		EntityNum:         int(do.EntityNum),
		OnlineNum:         int(do.OnlineNum),
		ServerName:        do.ServerName,
		Pod:               do.Pod,
		ServerType:        do.ServerType,
		ServerVersion:     do.ServerVersion,
		CodeVersionGame:   do.CodeVersionGame,
		CodeVersionConfig: do.CodeVersionConfig,
		CodeVersionScript: do.CodeVersionScript,
		LockEntityStatus:  int(do.EntityLockStatus),
		LockLbStatus:      int(do.LbLockStatus),
		ClusterName:       do.ClusterName,
		Namespace:         do.Namespace,
		ServerTypeDesc:    do.ServerTypeDesc,
		BigArea:           do.BigArea,
		GameAppId:         do.GameAppId,
	}
}

// ExistsGameServer 判断游戏进程是否存在
// 查询条件：集群名、命名空间、pod
func (x *GameServerRepo) ExistsGameServer(ctx context.Context, req *biz.GameServer) (bool, error) {
	model := dao.GameServer.Ctx(ctx)
	model = model.Where(dao.GameServer.Columns().Pod, req.Pod).
		Where(dao.GameServer.Columns().Namespace, req.Namespace).
		Where(dao.GameServer.Columns().ClusterName, req.ClusterName)
	count, err := model.Count()
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (x *GameServerRepo) CreateGameServer(ctx context.Context, req *biz.GameServer) error {
	do := x.convertDO2VO(req)
	_, err := dao.GameServer.Ctx(ctx).Insert(do)
	if err != nil {
		return fmt.Errorf("创建游戏进程失败: %w", err)
	}
	return nil
}

func (x *GameServerRepo) UpdateGameServer(ctx context.Context, req *biz.GameServer) error {
	db := dao.GameServer.Ctx(ctx)
	vo := x.convertDO2VO(req)
	updates := make(g.Map)
	if vo.Workload != "" {
		updates["workload"] = vo.Workload
	}
	if vo.WorkloadType != "" {
		updates["workload_type"] = vo.WorkloadType
	}
	if vo.EntityNum >= 0 {
		updates["entity_num"] = vo.EntityNum
	}
	if vo.OnlineNum >= 0 {
		updates["online_num"] = vo.OnlineNum
	}
	if vo.ServerName != "" {
		updates["server_name"] = vo.ServerName
	}
	if vo.Pod != "" {
		updates["pod"] = vo.Pod
	}
	if vo.ServerType != "" {
		updates["server_type"] = vo.ServerType
	}
	if vo.ServerVersion != "" {
		updates["server_version"] = vo.ServerVersion
	}
	if vo.CodeVersionGame != "" {
		updates["code_version_game"] = vo.CodeVersionGame
	}
	if vo.CodeVersionConfig != "" {
		updates["code_version_config"] = vo.CodeVersionConfig
	}
	if vo.CodeVersionScript != "" {
		updates["code_version_script"] = vo.CodeVersionScript
	}
	if vo.LockEntityStatus >= 0 {
		updates["lock_entity_status"] = vo.LockEntityStatus
	}
	if vo.LockLbStatus >= 0 {
		updates["lock_lb_status"] = vo.LockLbStatus
	}
	if vo.ClusterName != "" {
		updates["cluster_name"] = vo.ClusterName
	}
	if vo.Namespace != "" {
		updates["namespace"] = vo.Namespace
	}
	_, err := db.Where(dao.GameServer.Columns().ClusterName, req.ClusterName).
		Where(dao.GameServer.Columns().Namespace, req.Namespace).
		Where(dao.GameServer.Columns().ServerName, req.ServerName).
		Where(dao.GameServer.Columns().Pod, req.Pod).Update(updates)
	return err

}

func (x *GameServerRepo) DeleteGameServerByPod(ctx context.Context, req *biz.DeleteGameServerRequest) error {
	model := dao.GameServer.Ctx(ctx).Where(dao.GameServer.Columns().Pod, req.Pod).
		Where(dao.GameServer.Columns().Namespace, req.Namespace).
		Where(dao.GameServer.Columns().ClusterName, req.ClusterName)
	_, err := model.Delete()
	return err
}

func (x *GameServerRepo) ListGameServer(ctx context.Context, req *biz.ListGameServerRequest) ([]*biz.GameServer, uint32, error) {
	model := x.convertQuery(ctx, req)
	var vos []*entity.GameServer
	count, err := model.Count()
	if err != nil {
		return nil, 0, err
	}
	err = model.Scan(&vos)
	if err != nil {
		return nil, 0, err
	}
	var dos []*biz.GameServer
	for _, vo := range vos {
		dos = append(dos, x.convertVO2DO(vo))
	}
	return dos, uint32(count), nil
}

func (x *GameServerRepo) ListGameServeType(ctx context.Context, req *biz.ListGameServerTypeRequest) ([]*biz.GameServerType, uint32, error) {
	model := dao.GameServer.Ctx(ctx).
		Fields(dao.GameServer.Columns().ServerTypeDesc).
		Distinct().
		Where(dao.GameServer.Columns().DeletedAt, nil)
	if req.ClusterName != "" {
		model = model.Where(dao.GameServer.Columns().ClusterName, req.ClusterName)
	}
	if req.Namespace != "" {
		model = model.Where(dao.GameServer.Columns().Namespace, req.Namespace)
	}
	var vos []*entity.GameServer
	count, err := model.Count()
	if err != nil {
		return nil, 0, err
	}
	err = model.Scan(&vos)
	if err != nil {
		return nil, 0, err
	}
	var dos []*biz.GameServerType
	for _, vo := range vos {
		dos = append(dos, &biz.GameServerType{
			Name: vo.ServerTypeDesc,
		})
	}
	return dos, uint32(count), nil
}

func (x *GameServerRepo) GetGameServerByServerName(ctx context.Context, clusterName, namespace, serverName string) (*biz.GameServer, error) {
	model := dao.GameServer.Ctx(ctx)
	var gameServer *entity.GameServer
	err := model.Where(dao.GameServer.Columns().ServerName, serverName).
		Where(dao.GameServer.Columns().Namespace, namespace).
		Where(dao.GameServer.Columns().ClusterName, clusterName).Scan(&gameServer)
	if err != nil {
		return nil, fmt.Errorf("查询游戏进程信息失败: %w", err)
	}
	if gameServer == nil {
		return nil, fmt.Errorf("游戏进程不存在")
	}
	return x.convertVO2DO(gameServer), nil
}

func NewGameServerRepo(data *Data, logger log.Logger) *GameServerRepo {
	return &GameServerRepo{
		data: data,
		log:  log.NewHelper(log.With(logger, "module", "data/gameserver")),
	}
}

func NewIGameServerRepo(x *GameServerRepo) biz.IGameServerRepo {
	return x
}
