package data

import (
	"codo-cnmp/common/consts"
	"codo-cnmp/internal/biz"
	"codo-cnmp/internal/conf"
	"codo-cnmp/internal/model/dao"
	"codo-cnmp/internal/model/model/entity"
	"context"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/gogf/gf/v2/database/gdb"
)

type AgentRepo struct {
	data *Data
	log  *log.Helper
	bc   *conf.Bootstrap
}

func (x *AgentRepo) GetSrcAgentWhiteList(ctx context.Context) ([]string, error) {
	return x.bc.GetMESH().GetWHITE_IP_LIST(), nil
}

func (x *AgentRepo) GetSrcAgentID(ctx context.Context) (string, error) {
	return x.bc.GetMESH().GetSRC_AGENT_ID(), nil
}

// GetSrcAgentPort 获取代理端口. 代理存储在redis，用于代理服务的监听
func (x *AgentRepo) GetSrcAgentPort(ctx context.Context, req *biz.GetAgentPortRequest) (int, error) {
	var port int
	redisClient := x.data.redis
	key := fmt.Sprintf(consts.AgentPortCacheKey, req.AgentId)
	port, err := redisClient.Get(ctx, key).Int()
	if err != nil {
		port = int(x.bc.MESH.GetSRC_AGENT_PORT())
	}
	// 递增端口
	port += 1

	// 设置新端口
	_, err = redisClient.Set(ctx, key, port, 0).Result()
	if err != nil {
		x.log.WithContext(ctx).Errorf("设置代理端口失败: %v", err)
		return 0, err
	}
	return port, nil
}

func (x *AgentRepo) convertQuery(db *gdb.Model, query *biz.ListAgentRequest) *gdb.Model {
	if query.Keyword != "" {
		db = db.WhereLike(dao.Role.Columns().Name, "%"+query.Keyword+"%")
	}
	if query.ListAll {
		return db
	}
	if query.Page <= 0 {
		query.Page = 1
	}
	if query.PageSize <= 0 {
		query.PageSize = 10
	}
	return db.Page(int(query.Page), int(query.PageSize))

}

func (x *AgentRepo) convertVO2DO(t *entity.ProxyAgent) *biz.AgentItem {
	return &biz.AgentItem{
		Id:      int(t.Id),
		Name:    t.Name,
		AgentId: t.AgentId,
	}

}

func (x *AgentRepo) ListAgent(ctx context.Context, req *biz.ListAgentRequest) ([]*biz.AgentItem, uint32, error) {
	model := dao.ProxyAgent.Ctx(ctx)
	db := x.convertQuery(model, req)
	total, err := db.Count()
	if err != nil {
		return nil, 0, err
	}
	var agents []*biz.AgentItem
	err = db.Scan(&agents)
	if err != nil {
		return nil, 0, err
	}
	return agents, uint32(total), nil

}

func (x *AgentRepo) GetAgent(ctx context.Context, req *biz.GetAgentRequest) (*biz.AgentItem, error) {
	model := dao.ProxyAgent.Ctx(ctx)
	var agent entity.ProxyAgent
	err := model.Where(dao.ProxyAgent.Columns().Id, req.Id).Scan(&agent)
	if err != nil {
		return nil, err
	}
	return x.convertVO2DO(&agent), nil
}

func (x *AgentRepo) CreateAgent(ctx context.Context, req *biz.CreateAgentRequest) error {
	model := dao.ProxyAgent.Ctx(ctx)
	item := &biz.AgentItem{
		Name:    req.Name,
		AgentId: req.AgentId,
	}
	// 查询是否存在
	agent, err := model.Where(dao.ProxyAgent.Columns().AgentId, req.AgentId).One()
	if err != nil {
		return err
	}
	if agent != nil {
		return fmt.Errorf("agentId %s 已存在", req.AgentId)
	}

	vo := x.convertDO2VO(item)
	_, err = model.Insert(dao.ProxyAgent.Table, vo)
	return err
}

func (x *AgentRepo) convertDO2VO(data *biz.AgentItem) *entity.ProxyAgent {
	return &entity.ProxyAgent{
		Name:    data.Name,
		AgentId: data.AgentId,
	}

}

func (x *AgentRepo) UpdateAgent(ctx context.Context, req *biz.UpdateAgentRequest) error {
	model := dao.ProxyAgent.Ctx(ctx)
	item := &biz.AgentItem{
		Id:      req.Id,
		Name:    req.Name,
		AgentId: req.AgentId,
	}
	vo := x.convertDO2VO(item)
	_, err := model.Where(dao.ProxyAgent.Columns().Id, req.Id).Update(vo)
	return err
}

func (x *AgentRepo) DeleteAgent(ctx context.Context, req *biz.DeleteAgentRequest) error {
	model := dao.ProxyAgent.Ctx(ctx)
	_, err := model.Where(dao.ProxyAgent.Columns().Id, req.Id).Delete()
	return err
}

func NewAgentRepo(data *Data, logger log.Logger, bc *conf.Bootstrap) *AgentRepo {
	return &AgentRepo{
		data: data,
		bc:   bc,
		log:  log.NewHelper(log.With(logger, "module", "data/agent")),
	}
}

func NewIAgentRepo(x *AgentRepo) biz.IAgentRepository {
	return x
}
