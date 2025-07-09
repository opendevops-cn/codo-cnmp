package biz

import (
	"codo-cnmp/common/idip"
	"codo-cnmp/common/utils"
	"codo-cnmp/internal/conf"
	"context"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"time"
)

const (
	IdiPDuration      = 10 * time.Second
	IdipLBCommand     = 7005
	IdipEntityCommand = 7004
)

type GameServer struct {
	ID                uint32 // 主键ID
	ServerName        string // 进程名称
	ServerType        string // 进程类型
	ServerTypeDesc    string // 进程类型描述
	Pod               string // 命名空间
	WorkloadType      string // 工作负载类型
	Workload          string // 工作负载名称
	ServerVersion     string // server version
	EntityNum         uint32 // entity num
	OnlineNum         uint32 // online num
	EntityLockStatus  uint32 // entity_lock_status
	LbLockStatus      uint32 // lb_lock_status
	CodeVersionGame   string // code_version_game
	CodeVersionConfig string // code_version_config
	CodeVersionScript string // code_version_script
	ClusterName       string // 集群名称
	Namespace         string // 命名空间
	BigArea           string // 大区名称
	GameAppId         string // 游戏应用ID
}

type ListGameServerRequest struct {
	ClusterName      string // 集群名称
	Namespace        string // 命名空间
	Keyword          string // 查询关键字
	Page             uint32 // 页码
	PageSize         uint32 // 页大小
	ListAll          bool   // 查询全部
	ServerType       string // 进程类型
	EntityLockStatus uint32 // entity锁定状态
	LbLockStatus     uint32 // lb 锁定状态
}

type GameServerType struct {
	Name string // 进程类型名称
}

type ListGameServerTypeRequest struct {
	ClusterName string // 集群名称
	Namespace   string // 命名空间
	Keyword     string // 查询关键字
	Page        uint32 // 页码
	PageSize    uint32 // 页大小
	ListAll     bool   // 查询全部
}

type DeleteGameServerRequest struct {
	ClusterName string
	Namespace   string
	Pod         string
}

type ManageEntityRequest struct {
	ClusterName string
	Namespace   string
	ServerName  string // 进程名称
	Lock        bool   // 锁定/解锁
}

type ManageLBRequest struct {
	ClusterName string
	Namespace   string
	ServerName  string
	Lock        bool // 锁定/解锁
}

type BatchManageEntityRequest struct {
	ClusterName string
	Namespace   string
	ServerNames []string
	Lock        bool // 锁定/解锁
}

type BatchManageLBRequest struct {
	ClusterName string
	Namespace   string
	ServerNames []string
	Lock        bool
}

type IGameServerUseCase interface {
	// ListGameServer 进程列表
	ListGameServer(ctx context.Context, req *ListGameServerRequest) ([]*GameServer, uint32, error)
	// ListGameServeType 进程类型列表
	ListGameServeType(ctx context.Context, req *ListGameServerTypeRequest) ([]*GameServerType, uint32, error)
	// GetGameServerByServerName 获取进程详情
	GetGameServerByServerName(ctx context.Context, clusterName, namespace, serverName string) (*GameServer, error)
	// UpdateGameServer 更新进程信息
	UpdateGameServer(ctx context.Context, req *GameServer) error
	// BatchUpdateGameServer 批量更新进程信息
	BatchUpdateGameServer(ctx context.Context, req []*GameServer) error
	// DeleteGameServerByPod 删除进程
	DeleteGameServerByPod(ctx context.Context, req *DeleteGameServerRequest) error
	// CreateGameServer 创建进程
	CreateGameServer(ctx context.Context, req *GameServer) error
	// ExistsGameServer 进程是否存在
	ExistsGameServer(ctx context.Context, req *GameServer) (bool, error)
	// ManageEntity 管理entity
	ManageEntity(ctx context.Context, req *ManageEntityRequest) (*idip.ResponseBody, error)
	// ManageLB 管理lb
	ManageLB(ctx context.Context, req *ManageLBRequest) (*idip.ResponseBody, error)
	// BatchManageEntity 批量管理entity
	BatchManageEntity(ctx context.Context, req *BatchManageEntityRequest) error
	// BatchManageLB 批量管理lb
	BatchManageLB(ctx context.Context, req *BatchManageLBRequest) error
}

type IGameServerRepo interface {
	// ListGameServer 进程列表
	ListGameServer(ctx context.Context, req *ListGameServerRequest) ([]*GameServer, uint32, error)
	// ListGameServeType 进程类型列表
	ListGameServeType(ctx context.Context, req *ListGameServerTypeRequest) ([]*GameServerType, uint32, error)
	// GetGameServerByServerName 获取进程详情
	GetGameServerByServerName(ctx context.Context, clusterName, namespace, serverName string) (*GameServer, error)
	// UpdateGameServer 更新进程信息
	UpdateGameServer(ctx context.Context, req *GameServer) error
	// BatchUpdateGameServer 批量更新进程信息
	BatchUpdateGameServer(ctx context.Context, req []*GameServer) error
	// DeleteGameServerByPod 删除进程
	DeleteGameServerByPod(ctx context.Context, req *DeleteGameServerRequest) error
	// CreateGameServer 创建进程
	CreateGameServer(ctx context.Context, req *GameServer) error
	// ExistsGameServer 进程是否存在
	ExistsGameServer(ctx context.Context, req *GameServer) (bool, error)
}

type GameServerUseCase struct {
	cluster IClusterUseCase
	repo    IGameServerRepo
	log     *log.Helper
	bc      *conf.Bootstrap
}

func (x *GameServerUseCase) BatchUpdateGameServer(ctx context.Context, req []*GameServer) error {
	return x.repo.BatchUpdateGameServer(ctx, req)
}

func (x *GameServerUseCase) BatchManageEntity(ctx context.Context, req *BatchManageEntityRequest) error {
	for _, serverName := range req.ServerNames {
		responseBody, err := x.ManageEntity(ctx, &ManageEntityRequest{
			ClusterName: req.ClusterName,
			Namespace:   req.Namespace,
			ServerName:  serverName,
			Lock:        req.Lock,
		})
		if err != nil {
			x.log.WithContext(ctx).Errorf("管理entity失败: 集群：[%s], 命名空间：[%s], 进程名称：[%s], 原因：%v", req.ClusterName, req.Namespace, serverName, err)
		}
		x.log.WithContext(ctx).Infof("管理entity成功: 集群：[%s], 命名空间：[%s], 进程名称：[%s], Lock操作：[%v], 结果：%v", req.ClusterName, req.Namespace, serverName, req.Lock, responseBody)
	}
	return nil
}

func (x *GameServerUseCase) BatchManageLB(ctx context.Context, req *BatchManageLBRequest) error {
	for _, serverName := range req.ServerNames {
		responseBody, err := x.ManageLB(ctx, &ManageLBRequest{
			ClusterName: req.ClusterName,
			Namespace:   req.Namespace,
			ServerName:  serverName,
			Lock:        req.Lock,
		})
		if err != nil {
			x.log.WithContext(ctx).Errorf("管理lb失败: 集群：[%s], 命名空间：[%s], 进程名称：[%s], 原因：%v", req.ClusterName, req.Namespace, serverName, err)
		}
		x.log.WithContext(ctx).Infof("管理lb成功: 集群：[%s], 命名空间：[%s], 进程名称：[%s], Lock操作：[%v],"+
			" 结果：%v", req.ClusterName, req.Namespace, serverName, req.Lock, responseBody)
	}
	return nil
}

func (x *GameServerUseCase) ManageEntity(ctx context.Context, req *ManageEntityRequest) (*idip.ResponseBody, error) {
	gs, err := x.repo.GetGameServerByServerName(ctx, req.ClusterName, req.Namespace, req.ServerName)
	if err != nil {
		return nil, fmt.Errorf("查询游戏进程失败: %v", err)
	}
	cluster, err := x.cluster.GetClusterByName(ctx, req.ClusterName)
	if err != nil {
		return nil, fmt.Errorf("查询集群失败: %v", err)
	}
	requestBody := idip.NewRequestBody(
		idip.WithBigArea(gs.BigArea),
		idip.WithGameAppid(gs.GameAppId),
		idip.WithLock(req.Lock),
		idip.WithCmd(IdipEntityCommand),
		idip.WithServerName(req.ServerName))
	var decryptedSecret string
	if cluster.AppSecret != "" {
		decryptedSecret, err = utils.AESDecrypt(cluster.AppSecret, x.bc.APP.GetSECRET())
		if err != nil {
			return nil, fmt.Errorf("解密appSecret失败: %v", err)
		}
	}
	api, err := idip.NewIDIPBaseAPI(cluster.IDIP, gs.GameAppId, cluster.AppId, decryptedSecret, IdiPDuration, x.log)
	if err != nil {
		return nil, fmt.Errorf("创建idip api失败: %v", err)
	}
	responseBody, err := api.SendRequest(ctx, requestBody)
	if err != nil {
		return nil, err
	}
	return responseBody, nil
}

func (x *GameServerUseCase) ManageLB(ctx context.Context, req *ManageLBRequest) (*idip.ResponseBody, error) {
	gs, err := x.repo.GetGameServerByServerName(ctx, req.ClusterName, req.Namespace, req.ServerName)
	if err != nil {
		return nil, fmt.Errorf("查询游戏进程失败: %v", err)
	}
	cluster, err := x.cluster.GetClusterByName(ctx, req.ClusterName)
	if err != nil {
		return nil, fmt.Errorf("查询集群失败: %v", err)
	}
	requestBody := idip.NewRequestBody(idip.WithBigArea(gs.BigArea),
		idip.WithGameAppid(gs.GameAppId),
		idip.WithLock(req.Lock),
		idip.WithCmd(IdipLBCommand),
		idip.WithServerName(req.ServerName))
	var decryptedSecret string
	if cluster.AppSecret != "" {
		decryptedSecret, err = utils.AESDecrypt(cluster.AppSecret, x.bc.APP.GetSECRET())
		if err != nil {
			return nil, fmt.Errorf("解密appSecret失败: %v", err)
		}
	}
	api, err := idip.NewIDIPBaseAPI(cluster.IDIP, gs.GameAppId, cluster.AppId, decryptedSecret, IdiPDuration, x.log)
	if err != nil {
		return nil, fmt.Errorf("创建idip api失败: %v", err)
	}
	responseBody, err := api.SendRequest(ctx, requestBody)
	if err != nil {
		return nil, err
	}
	return responseBody, nil
}

func (x *GameServerUseCase) ExistsGameServer(ctx context.Context, req *GameServer) (bool, error) {
	return x.repo.ExistsGameServer(ctx, req)
}

func (x *GameServerUseCase) GetGameServerByServerName(ctx context.Context, clusterName, namespace, serverName string) (*GameServer, error) {
	return x.repo.GetGameServerByServerName(ctx, clusterName, namespace, serverName)
}

func (x *GameServerUseCase) UpdateGameServer(ctx context.Context, req *GameServer) error {
	return x.repo.UpdateGameServer(ctx, req)
}

func (x *GameServerUseCase) DeleteGameServerByPod(ctx context.Context, req *DeleteGameServerRequest) error {
	return x.repo.DeleteGameServerByPod(ctx, req)
}

func (x *GameServerUseCase) CreateGameServer(ctx context.Context, req *GameServer) error {
	return x.repo.CreateGameServer(ctx, req)
}

func (x *GameServerUseCase) ListGameServer(ctx context.Context, req *ListGameServerRequest) ([]*GameServer, uint32, error) {
	return x.repo.ListGameServer(ctx, req)
}

func (x *GameServerUseCase) ListGameServeType(ctx context.Context, req *ListGameServerTypeRequest) ([]*GameServerType, uint32, error) {
	return x.repo.ListGameServeType(ctx, req)
}

func NewIGameServerUseCase(x *GameServerUseCase) IGameServerUseCase {
	return x
}

func NewGameServerUseCase(repo IGameServerRepo, logger log.Logger, cluster IClusterUseCase, bc *conf.Bootstrap) *GameServerUseCase {
	return &GameServerUseCase{
		repo:    repo,
		cluster: cluster,
		log:     log.NewHelper(log.With(logger, "module", "biz/gameserver")),
		bc:      bc,
	}
}
