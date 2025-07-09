package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"time"
)

type ListMeshRequest struct {
	Page     uint32
	PageSize uint32
	ListAll  bool
}

type MeshItem struct {
	Id             int       `json:"id"`
	ServiceName    string    `json:"service_name"`
	WhiteIpList    []string  `json:"white_ip_list"`
	SrcAgentId     string    `json:"src_agent_id"`
	SrcAgentPort   int       `json:"src_agent_port"`
	DstAgentId     string    `json:"dst_agent_id"`
	DstServiceAddr string    `json:"dst_service_addr"`
	HeartbeatAt    time.Time `json:"heartbeat_at"`
	CreatedAt      time.Time `json:"created_at"`
}

type IMeshUseCase interface {
	// ListMesh 获取Mesh列表
	ListMesh(ctx context.Context, req *ListMeshRequest) ([]*MeshItem, uint32, error)
	// CreateMesh 注册Mesh
	CreateMesh(ctx context.Context, req *MeshItem) (string, error)
	// DeleteMesh 删除Mesh
	DeleteMesh(ctx context.Context, id string) error
}

type IMeshRepository interface {
	// ListMesh 获取Mesh列表
	ListMesh(ctx context.Context, req *ListMeshRequest) ([]*MeshItem, uint32, error)
	// CreateMesh 注册Mesh
	CreateMesh(ctx context.Context, req *MeshItem) (string, error)
	// DeleteMesh 删除Mesh
	DeleteMesh(ctx context.Context, id string) error
}

type MeshUseCase struct {
	repo IMeshRepository
	log  *log.Helper
}

func (x *MeshUseCase) DeleteMesh(ctx context.Context, id string) error {
	return x.repo.DeleteMesh(ctx, id)
}

func (x *MeshUseCase) CreateMesh(ctx context.Context, req *MeshItem) (string, error) {
	return x.repo.CreateMesh(ctx, req)
}

func (x *MeshUseCase) ListMesh(ctx context.Context, req *ListMeshRequest) ([]*MeshItem, uint32, error) {
	return x.repo.ListMesh(ctx, req)
}

func NewMeshUseCase(x IMeshRepository, logger log.Logger) *MeshUseCase {
	return &MeshUseCase{repo: x, log: log.NewHelper(log.With(logger, "module", "biz/mesh"))}
}

func NewIMeshUseCase(x *MeshUseCase) IMeshUseCase {
	return x
}
