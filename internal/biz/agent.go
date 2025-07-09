package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
)

type ListAgentRequest struct {
	Page     uint32
	PageSize uint32
	ListAll  bool
	Keyword  string
}

type AgentCommonParams struct {
	AgentId string `json:"agent_id"`
	Name    string `json:"name"`
}

type AgentItem struct {
	Id      int    `json:"id"`
	AgentId string `json:"agent_id"`
	Name    string `json:"name"`
}

type CreateAgentRequest struct {
	AgentCommonParams
}

type UpdateAgentRequest struct {
	Id int
	AgentCommonParams
}

type DeleteAgentRequest struct {
	Id int `json:"id"`
}

type GetAgentRequest struct {
	Id int `json:"id"`
}

type GetAgentPortRequest struct {
	AgentId string `json:"agent_id"`
}

type IAgentUseCase interface {
	// ListAgent 获取Agent列表
	ListAgent(ctx context.Context, req *ListAgentRequest) ([]*AgentItem, uint32, error)
	// GetAgent 注册Agent
	GetAgent(ctx context.Context, req *GetAgentRequest) (*AgentItem, error)
	// CreateAgent 注册Agent
	CreateAgent(ctx context.Context, req *CreateAgentRequest) error
	// UpdateAgent 更新Agent
	UpdateAgent(ctx context.Context, req *UpdateAgentRequest) error
	// DeleteAgent 删除Agent
	DeleteAgent(ctx context.Context, req *DeleteAgentRequest) error
	// GetSrcAgentPort 获取源Agent端口
	GetSrcAgentPort(ctx context.Context, req *GetAgentPortRequest) (int, error)
	// GetSrcAgentID 获取源Agent端口
	GetSrcAgentID(ctx context.Context) (string, error)
	// GetSrcAgentWhiteList 获取源Agent白名单
	GetSrcAgentWhiteList(ctx context.Context) ([]string, error)
}

type IAgentRepository interface {
	// ListAgent 获取Agent列表
	ListAgent(ctx context.Context, req *ListAgentRequest) ([]*AgentItem, uint32, error)
	// GetAgent 获取Agent
	GetAgent(ctx context.Context, req *GetAgentRequest) (*AgentItem, error)
	// CreateAgent 注册Agent
	CreateAgent(ctx context.Context, req *CreateAgentRequest) error
	// UpdateAgent 更新Agent
	UpdateAgent(ctx context.Context, req *UpdateAgentRequest) error
	// DeleteAgent 删除Agent
	DeleteAgent(ctx context.Context, req *DeleteAgentRequest) error
	// GetSrcAgentPort 获取源Agent端口
	GetSrcAgentPort(ctx context.Context, req *GetAgentPortRequest) (int, error)
	// GetSrcAgentID 获取源Agent端口
	GetSrcAgentID(ctx context.Context) (string, error)
	// GetSrcAgentWhiteList 获取源Agent白名单
	GetSrcAgentWhiteList(ctx context.Context) ([]string, error)
}

type AgentUseCase struct {
	repo IAgentRepository
	log  *log.Helper
}

func (x *AgentUseCase) GetSrcAgentPort(ctx context.Context, req *GetAgentPortRequest) (int, error) {
	return x.repo.GetSrcAgentPort(ctx, req)
}

func (x *AgentUseCase) GetSrcAgentID(ctx context.Context) (string, error) {
	return x.repo.GetSrcAgentID(ctx)
}

func (x *AgentUseCase) GetSrcAgentWhiteList(ctx context.Context) ([]string, error) {
	return x.repo.GetSrcAgentWhiteList(ctx)
}

func (x *AgentUseCase) ListAgent(ctx context.Context, req *ListAgentRequest) ([]*AgentItem, uint32, error) {
	return x.repo.ListAgent(ctx, req)
}

func (x *AgentUseCase) GetAgent(ctx context.Context, req *GetAgentRequest) (*AgentItem, error) {
	return x.repo.GetAgent(ctx, req)
}

func (x *AgentUseCase) CreateAgent(ctx context.Context, req *CreateAgentRequest) error {
	return x.repo.CreateAgent(ctx, req)
}

func (x *AgentUseCase) UpdateAgent(ctx context.Context, req *UpdateAgentRequest) error {
	return x.repo.UpdateAgent(ctx, req)
}

func (x *AgentUseCase) DeleteAgent(ctx context.Context, req *DeleteAgentRequest) error {
	return x.repo.DeleteAgent(ctx, req)
}

func NewAgentUseCase(x IAgentRepository, logger log.Logger) *AgentUseCase {
	return &AgentUseCase{repo: x, log: log.NewHelper(log.With(logger, "module", "biz/agent"))}
}

func NewIAgentUseCase(x *AgentUseCase) IAgentUseCase {
	return x
}
