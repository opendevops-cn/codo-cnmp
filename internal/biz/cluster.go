package biz

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"

	"codo-cnmp/internal/conf"
	"github.com/google/uuid"
	"k8s.io/client-go/tools/clientcmd"

	"codo-cnmp/common/consts"
	"codo-cnmp/common/k8s"
	"codo-cnmp/common/utils"
	"codo-cnmp/pb"
	"github.com/ccheers/xpkg/sync/errgroup"
	"github.com/ccheers/xpkg/xmsgbus"
	"github.com/go-kratos/kratos/v2/log"
	kruiseclientset "github.com/openkruise/kruise-api/client/clientset/versioned"
	kruisegameclientset "github.com/openkruise/kruise-game/pkg/client/clientset/versioned"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/version"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/metrics/pkg/client/clientset/versioned"
)

var ErrClusterNameExist = fmt.Errorf("集群名称已存在")

type ClusterRepo interface {
	// ListClusters 查看集群列表
	ListClusters(ctx context.Context, query *QueryClusterReq) (ClusterItems, error)

	// CreateCluster 创建集群
	CreateCluster(ctx context.Context, data *ClusterItem) (uint32, error)

	// DeleteCluster 删除集群
	DeleteCluster(ctx context.Context, id uint32) error

	// ExistCluster 检查集群名称是否重复
	ExistCluster(ctx context.Context, name string) (bool, error)

	// CountCluster 统计集群数量
	CountCluster(ctx context.Context, query *QueryClusterReq) (uint32, error)

	// FetchAllClusters 获取所有集群
	FetchAllClusters(ctx context.Context) (ClusterItems, error)

	// UpdateCluster 更新集群
	UpdateCluster(ctx context.Context, data *ClusterItem) error

	// UpdateClusterV2 更新集群
	UpdateClusterV2(ctx context.Context, data *ClusterItem) error

	// GetClusterByID 获取集群
	GetClusterByID(ctx context.Context, id uint32) (*ClusterItem, error)

	// GetClusterByName 获取集群
	GetClusterByName(ctx context.Context, name string) (*ClusterItem, error)
	UpdateClusterState(ctx context.Context, req *UpdateClusterStateRequest) error
	UpdateClusterBasicInfo(ctx context.Context, req *UpdateClusterBasicRequest) error
}

type UpdateClusterStateRequest struct {
	Id           uint32
	ClusterState ClusterState
	HealthState  []HealthState
}

type UpdateClusterBasicRequest struct {
	Id            uint32
	ServerVersion string
	Platform      string
	BuildDate     string
	CpuTotal      float32
	MemoryTotal   float32
	CpuUsage      float32
	MemoryUsage   float32
	NodeCount     int
	ClusterState  ClusterState
	HealthState   []HealthState
}

type DownLoadKubeConfigRequest struct {
	Id uint32
}

type IClusterUseCase interface {
	GetClientConfigByClusterID(ctx context.Context, clusterID uint32) (*k8s.ClientConfig, error)
	GetClientConfigByClusterName(ctx context.Context, clusterName string) (*k8s.ClientConfig, error)
	ListClusters(ctx context.Context, query *QueryClusterReq) (ClusterItems, uint32, error)
	CreateCluster(ctx context.Context, data *ClusterItem) (uint32, error)
	DeleteCluster(ctx context.Context, id uint32) error
	FetchAllClusters(ctx context.Context) (ClusterItems, error)
	UpdateCluster(ctx context.Context, data *ClusterItem) error
	GetClusterByID(ctx context.Context, id uint32) (*ClusterItem, error)
	GetClusterByName(ctx context.Context, name string) (*ClusterItem, error)
	GetDynamicClientByClusterName(ctx context.Context, clusterName string) (*dynamic.DynamicClient, error)
	GetClientSetByClusterName(ctx context.Context, clusterName string) (*kubernetes.Clientset, error)
	GetKruiseClientSetByClusterName(ctx context.Context, clusterName string) (*kruiseclientset.Clientset, error)
	GetKruiseGameClientSetByClusterName(ctx context.Context, clusterName string) (*kruisegameclientset.Clientset, error)
	GetMetricsClientSetByClusterName(ctx context.Context, clusterName string) (*versioned.Clientset, error)
	PingIdip(ctx context.Context, idip string) (bool, error)
	UpdateClusterState(ctx context.Context, req *UpdateClusterStateRequest) error
	UpdateClusterBasicInfo(ctx context.Context, req *UpdateClusterBasicRequest) error
	DownLoadKubeConfig(ctx context.Context, req *DownLoadKubeConfigRequest) (string, error)
}

type ClusterCreativeEvent struct {
	ClusterID uint32
}

func (x *ClusterCreativeEvent) Topic() string {
	return "cluster.creative"
}

type ClusterDeletionEvent struct {
	ClusterID uint32
}

func (x *ClusterDeletionEvent) Topic() string {
	return "cluster.deletion"
}

type ClusterUseCase struct {
	clusterRepo ClusterRepo
	nodeRepo    NodeRepo
	log         *log.Helper

	msgbus      xmsgbus.IMsgBus
	tm          xmsgbus.ITopicManager
	otelOptions *xmsgbus.OTELOptions

	agent       IAgentUseCase
	mesh        IMeshUseCase
	conf        *conf.Bootstrap
	roleBinding IRoleBindingUseCase

	pubClusterCreative xmsgbus.IPublisher[*ClusterCreativeEvent]
	pubClusterDeletion xmsgbus.IPublisher[*ClusterDeletionEvent]

	// 缓存各种clientSet
	clientSetCache        map[string]*kubernetes.Clientset          // key: clusterName
	dynamicClientCache    map[string]*dynamic.DynamicClient         // key: clusterName
	metricsClientCache    map[string]*versioned.Clientset           // key: clusterName
	kruiseClientCache     map[string]*kruiseclientset.Clientset     // key: clusterName
	kruiseGameClientCache map[string]*kruisegameclientset.Clientset // key: clusterName
	cacheMutex            sync.RWMutex                              // 保护缓存的读写锁
}

// 清除特定集群的所有clientSet缓存
func (x *ClusterUseCase) clearClientSetCache(clusterName string) {
	x.cacheMutex.Lock()
	defer x.cacheMutex.Unlock()
	delete(x.clientSetCache, clusterName)
	delete(x.dynamicClientCache, clusterName)
	delete(x.metricsClientCache, clusterName)
	delete(x.kruiseClientCache, clusterName)
	delete(x.kruiseGameClientCache, clusterName)
}

// 获取或创建kubernetes.Clientset
func (x *ClusterUseCase) getOrCreateClientSet(ctx context.Context, clusterName string) (*kubernetes.Clientset, error) {
	// 先尝试从缓存获取
	x.cacheMutex.RLock()
	if clientSet, exists := x.clientSetCache[clusterName]; exists {
		x.cacheMutex.RUnlock()
		return clientSet, nil
	}
	x.cacheMutex.RUnlock()

	// 缓存中没有，创建新的
	cfg, err := x.getClusterConfig(ctx, clusterName)
	if err != nil {
		return nil, err
	}

	clientSet, err := cfg.CreateClientSet()
	if err != nil {
		return nil, err
	}

	// 缓存结果
	x.cacheMutex.Lock()
	x.clientSetCache[clusterName] = clientSet
	x.cacheMutex.Unlock()

	return clientSet, nil
}

// 获取或创建dynamic.DynamicClient
func (x *ClusterUseCase) getOrCreateDynamicClient(ctx context.Context, clusterName string) (*dynamic.DynamicClient, error) {
	x.cacheMutex.RLock()
	if dynamicClient, exists := x.dynamicClientCache[clusterName]; exists {
		x.cacheMutex.RUnlock()
		return dynamicClient, nil
	}
	x.cacheMutex.RUnlock()

	cfg, err := x.getClusterConfig(ctx, clusterName)
	if err != nil {
		return nil, err
	}

	dynamicClient, err := cfg.CreateDynamicClient()
	if err != nil {
		return nil, err
	}

	x.cacheMutex.Lock()
	x.dynamicClientCache[clusterName] = dynamicClient
	x.cacheMutex.Unlock()

	return dynamicClient, nil
}

// 获取或创建metrics clientSet
func (x *ClusterUseCase) getOrCreateMetricsClient(ctx context.Context, clusterName string) (*versioned.Clientset, error) {
	x.cacheMutex.RLock()
	if metricsClient, exists := x.metricsClientCache[clusterName]; exists {
		x.cacheMutex.RUnlock()
		return metricsClient, nil
	}
	x.cacheMutex.RUnlock()

	cfg, err := x.getClusterConfig(ctx, clusterName)
	if err != nil {
		return nil, err
	}

	metricsClient, err := cfg.CreateMetricsClientSet()
	if err != nil {
		return nil, err
	}

	x.cacheMutex.Lock()
	x.metricsClientCache[clusterName] = metricsClient
	x.cacheMutex.Unlock()

	return metricsClient, nil
}

// 获取或创建kruise clientSet
func (x *ClusterUseCase) getOrCreateKruiseClient(ctx context.Context, clusterName string) (*kruiseclientset.Clientset, error) {
	x.cacheMutex.RLock()
	if kruiseClient, exists := x.kruiseClientCache[clusterName]; exists {
		x.cacheMutex.RUnlock()
		return kruiseClient, nil
	}
	x.cacheMutex.RUnlock()

	cfg, err := x.getClusterConfig(ctx, clusterName)
	if err != nil {
		return nil, err
	}

	kruiseClient, err := cfg.CreateKruiseClientSet()
	if err != nil {
		return nil, err
	}

	x.cacheMutex.Lock()
	x.kruiseClientCache[clusterName] = kruiseClient
	x.cacheMutex.Unlock()

	return kruiseClient, nil
}

// 获取或创建kruise game clientSet
func (x *ClusterUseCase) getOrCreateKruiseGameClient(ctx context.Context, clusterName string) (*kruisegameclientset.Clientset, error) {
	x.cacheMutex.RLock()
	if kruiseGameClient, exists := x.kruiseGameClientCache[clusterName]; exists {
		x.cacheMutex.RUnlock()
		return kruiseGameClient, nil
	}
	x.cacheMutex.RUnlock()

	cfg, err := x.getClusterConfig(ctx, clusterName)
	if err != nil {
		return nil, err
	}

	kruiseGameClient, err := cfg.CreateKruiseGameClientSet()
	if err != nil {
		return nil, err
	}

	x.cacheMutex.Lock()
	x.kruiseGameClientCache[clusterName] = kruiseGameClient
	x.cacheMutex.Unlock()

	return kruiseGameClient, nil
}

func (x *ClusterUseCase) DownLoadKubeConfig(ctx context.Context, req *DownLoadKubeConfigRequest) (string, error) {
	clusterItem, err := x.clusterRepo.GetClusterByID(ctx, req.Id)
	if err != nil {
		return "", fmt.Errorf("获取集群信息失败: %w", err)
	}
	// 获取当前用户Id
	userId, err := utils.GetUserIDFromCtx(ctx)
	if err != nil {
		return "", fmt.Errorf("获取用户Id失败: %w", err)
	}
	// 生成token
	signKey := x.conf.APP.GetSECRET()
	if signKey == "" {
		return "", fmt.Errorf("获取signKey失败")
	}
	token, err := utils.GenerateToken(userId, req.Id, signKey)
	if err != nil {
		return "", fmt.Errorf("生成token失败: %w", err)
	}

	// 替换apiServer地址
	ApiServer := x.conf.PROXY.GetAPI_SERVER()
	if ApiServer == "" {
		return "", fmt.Errorf("服务地址未配置")
	}

	newApiServer := fmt.Sprintf("%s/%s", ApiServer, token)

	kubeConfig := clusterItem.ImportDetail.KubeConfig
	config, err := clientcmd.Load([]byte(kubeConfig))
	if err != nil {
		return "", fmt.Errorf("解析kubeconfig失败: %w", err)
	}
	// 获取当前上下文
	Context := config.Contexts[config.CurrentContext]
	if Context == nil {
		return "", fmt.Errorf("kubeConfig上下文缺失")
	}
	for _, cluster := range config.Clusters {
		cluster.Server = newApiServer
		cluster.InsecureSkipTLSVerify = true
		cluster.CertificateAuthorityData = nil
	}

	// 获取认证信息
	authInfo := config.AuthInfos[Context.AuthInfo]
	if authInfo == nil {
		return "", fmt.Errorf("kubeConfig认证信息缺失")
	}
	// 生成自定义证书并替换kubeConfig
	certData, keyData, err := utils.GenerateCert()
	if err != nil {
		return "", fmt.Errorf("生成证书失败: %w", err)
	}
	authInfo.ClientCertificateData = certData
	authInfo.ClientKeyData = keyData

	// 生成kubeConfig
	newKubeConfig, err := clientcmd.Write(*config)
	if err != nil {
		return "", fmt.Errorf("生成kubeConfig失败: %w", err)
	}
	return string(newKubeConfig), nil
}

func (x *ClusterUseCase) UpdateClusterBasicInfo(ctx context.Context, req *UpdateClusterBasicRequest) error {
	return x.clusterRepo.UpdateClusterBasicInfo(ctx, req)
}

func (x *ClusterUseCase) UpdateClusterState(ctx context.Context, req *UpdateClusterStateRequest) error {
	return x.clusterRepo.UpdateClusterState(ctx, req)
}

func (x *ClusterUseCase) getClusterConfig(ctx context.Context, clusterName string) (*k8s.ClientConfig, error) {
	clusterItem, err := x.clusterRepo.GetClusterByName(ctx, clusterName)
	if err != nil {
		return nil, fmt.Errorf("获取集群信息失败: %w", err)
	}
	return &k8s.ClientConfig{
		ImportType:  pb.ImportType(clusterItem.ImportType),
		KubeConfig:  clusterItem.ImportDetail.KubeConfig,
		Token:       clusterItem.ImportDetail.Token,
		APIServer:   clusterItem.ImportDetail.ApiServer,
		CaData:      clusterItem.ImportDetail.CaData,
		Context:     clusterItem.ImportDetail.Context,
		Agent:       clusterItem.ImportDetail.Agent,
		AgentImage:  clusterItem.ImportDetail.AgentImage,
		AgentProxy:  clusterItem.ImportDetail.AgentProxy,
		ClusterName: clusterName,
		ConnectType: clusterItem.ConnectType,
		MeshAddr:    clusterItem.MeshAddr,
	}, nil
}

func (x *ClusterUseCase) GetDynamicClientByClusterName(ctx context.Context, clusterName string) (*dynamic.DynamicClient, error) {
	//cfg, err := x.getClusterConfig(ctx, clusterName)
	//if err != nil {
	//	return nil, err
	//}
	//clientSet, err := cfg.CreateDynamicClient()
	//if err != nil {
	//	return nil, err
	//}
	//return clientSet, nil
	return x.getOrCreateDynamicClient(ctx, clusterName)
}

func (x *ClusterUseCase) PingIdip(ctx context.Context, idip string) (bool, error) {
	return utils.CheckConnection(idip)
}

func (x *ClusterUseCase) GetMetricsClientSetByClusterName(ctx context.Context, clusterName string) (*versioned.Clientset, error) {
	//cfg, err := x.getClusterConfig(ctx, clusterName)
	//if err != nil {
	//	return nil, err
	//}
	//metricsClient, err := cfg.CreateMetricsClientSet()
	//if err != nil {
	//	return nil, err
	//}
	//return metricsClient, nil
	return x.getOrCreateMetricsClient(ctx, clusterName)
}

func (x *ClusterUseCase) GetClientSetByClusterName(ctx context.Context, clusterName string) (*kubernetes.Clientset, error) {
	//cfg, err := x.getClusterConfig(ctx, clusterName)
	//if err != nil {
	//	return nil, err
	//}
	//clientSet, err := cfg.CreateClientSet()
	//if err != nil {
	//	return nil, err
	//}
	//return clientSet, nil
	return x.getOrCreateClientSet(ctx, clusterName)
}

func (x *ClusterUseCase) GetKruiseClientSetByClusterName(ctx context.Context, clusterName string) (*kruiseclientset.Clientset, error) {
	//cfg, err := x.getClusterConfig(ctx, clusterName)
	//if err != nil {
	//	return nil, err
	//}
	//kruiseClientSet, err := cfg.CreateKruiseClientSet()
	//if err != nil {
	//	return nil, err
	//}
	//return kruiseClientSet, nil
	return x.getOrCreateKruiseClient(ctx, clusterName)
}

func (x *ClusterUseCase) GetKruiseGameClientSetByClusterName(ctx context.Context, clusterName string) (*kruisegameclientset.Clientset, error) {
	//cfg, err := x.getClusterConfig(ctx, clusterName)
	//if err != nil {
	//	return nil, err
	//}
	//kruiseGameClientSet, err := cfg.CreateKruiseGameClientSet()
	//if err != nil {
	//	return nil, err
	//}
	//return kruiseGameClientSet, nil
	return x.getOrCreateKruiseGameClient(ctx, clusterName)
}

func NewIClusterUseCase(x *ClusterUseCase) IClusterUseCase {
	return x
}

// NewClusterUseCase 创建 ClusterUseCase
func NewClusterUseCase(ctx context.Context, bus xmsgbus.IMsgBus, tm xmsgbus.ITopicManager, repo ClusterRepo, nodeRepo NodeRepo, logger log.Logger,
	agent IAgentUseCase, mesh IMeshUseCase, conf *conf.Bootstrap, roleBinding *RoleBindingUseCase,
) (*ClusterUseCase, func()) {
	x := &ClusterUseCase{
		clusterRepo:        repo,
		nodeRepo:           nodeRepo,
		log:                log.NewHelper(log.With(logger, "module", "biz/cluster")),
		otelOptions:        xmsgbus.NewOTELOptions(),
		msgbus:             bus,
		tm:                 tm,
		agent:              agent,
		mesh:               mesh,
		pubClusterCreative: xmsgbus.NewPublisher[*ClusterCreativeEvent](bus, tm, xmsgbus.NewOTELOptions()),
		pubClusterDeletion: xmsgbus.NewPublisher[*ClusterDeletionEvent](bus, tm, xmsgbus.NewOTELOptions()),
		conf:               conf,
		roleBinding:        roleBinding,

		// 初始化clientSet缓存
		clientSetCache:        make(map[string]*kubernetes.Clientset),
		dynamicClientCache:    make(map[string]*dynamic.DynamicClient),
		metricsClientCache:    make(map[string]*versioned.Clientset),
		kruiseClientCache:     make(map[string]*kruiseclientset.Clientset),
		kruiseGameClientCache: make(map[string]*kruisegameclientset.Clientset),
	}
	return x, x.Init(ctx)
}

func (x *ClusterUseCase) Init(ctx context.Context) context.CancelFunc {
	ctx, cancel := context.WithCancel(ctx)
	eg := errgroup.WithCancel(ctx)

	eg.Go(func(ctx context.Context) error {
		sub := xmsgbus.NewSubscriber((*ClusterCreativeEvent)(nil).Topic(), "ClusterUseCase",
			x.msgbus, x.otelOptions, x.tm,
			xmsgbus.WithHandleFunc(func(ctx context.Context, msg *ClusterCreativeEvent) error {
				// handle event
				return x.HandleClusterCreative(ctx, msg.ClusterID)
			}),
		)
		defer sub.Close(ctx)
		for {
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
				if err := sub.Handle(ctx); err != nil {
					if errors.Is(err, context.DeadlineExceeded) {
						continue
					}
					x.log.Errorf("subcribe error: %v", err)
				}
			}
		}
	})
	eg.Go(func(ctx context.Context) error {
		sub := xmsgbus.NewSubscriber((*ClusterDeletionEvent)(nil).Topic(), "ClusterUseCase",
			x.msgbus, x.otelOptions, x.tm,
			xmsgbus.WithHandleFunc(func(ctx context.Context, msg *ClusterCreativeEvent) error {
				// handle event
				return x.HandleClusterDeletion(ctx, msg.ClusterID)
			}),
		)
		defer sub.Close(ctx)
		for {
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
				if err := sub.Handle(ctx); err != nil {
					if errors.Is(err, context.DeadlineExceeded) {
						continue
					}
					x.log.Errorf("subcribe error: %v", err)
				}
			}
		}
	})

	return func() {
		cancel()
		_ = eg.Wait()
	}
}

type (
	HealthState  pb.HealthState
	ClusterState pb.ClusterState
)

type Link struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

// ClusterItem 集群信息
type ClusterItem struct {
	ID            uint32
	Name          string
	Description   string
	ClusterState  ClusterState
	HealthState   []HealthState
	CreateTime    string
	BuildDate     string
	ServerVersion string
	CpuUsage      float32
	CpuTotal      float32
	MemoryTotal   float32
	Platform      string
	MemoryUsage   float32
	ImportType    ImportType
	ImportDetail  ImportDetail
	NodeCount     int
	UID           string
	IsFollowed    bool
	IDIP          string
	AppId         string
	AppSecret     string
	Ops           []string
	ConnectType   pb.ConnectType
	DstAgentId    uint32
	MeshAddr      string
	Links         []*pb.Link
}

type ClusterItems []*ClusterItem

// QueryClusterReq 查询集群请求
type QueryClusterReq struct {
	ID       uint32 `json:"id"`
	Page     uint32 `json:"page"`
	PageSize uint32 `json:"page_size"`
	Keyword  string `json:"keyword"`
	ListAll  bool   `json:"list_all"`
}

// ImportType 导入类型
type ImportType pb.ImportType

type ImportDetail struct {
	// kubeConfig详情
	KubeConfig string `json:"kube_config"`
	// token详情
	Token string `json:"token"`
	// apiServer地址
	ApiServer string `json:"api_server"`
	// ca证书数据
	CaData string `json:"ca_data"`
	// context名称
	Context string `json:"context"`
	// agent名称
	Agent string `json:"agent"`
	// agent镜像
	AgentImage string `json:"agent_image"`
	// agent proxy地址
	AgentProxy string `json:"agent_proxy"`
}

// CreateClusterReq 创建集群请求
type CreateClusterReq struct {
	ID           uint32
	Name         string
	Description  string
	ImportType   ImportType
	ImportDetail ImportDetail
}

// UpdateClusterReq 更新集群请求
type UpdateClusterReq struct {
	Id            uint32
	Name          string
	Description   string
	ServerVersion string
	Platform      string
	BuildDate     string
	CpuUsage      string
	MemoryUsage   string
	HealthState   uint8
	NodeState     uint8
}

// DeleteClusterReq 删除集群请求
type DeleteClusterReq struct {
	Id uint32
}

// GetClientConfigByClusterID 获取集群配置
func (x *ClusterUseCase) GetClientConfigByClusterID(ctx context.Context, clusterID uint32) (*k8s.ClientConfig, error) {
	clusterItem, err := x.GetClusterByID(ctx, clusterID)
	if err != nil {
		return nil, fmt.Errorf("获取集群信息失败: %v", err)
	}
	cfg := k8s.ClientConfig{
		ImportType:  pb.ImportType(clusterItem.ImportType),
		KubeConfig:  clusterItem.ImportDetail.KubeConfig,
		Token:       clusterItem.ImportDetail.Token,
		APIServer:   clusterItem.ImportDetail.ApiServer,
		CaData:      clusterItem.ImportDetail.CaData,
		Context:     clusterItem.ImportDetail.Context,
		Agent:       clusterItem.ImportDetail.Agent,
		AgentImage:  clusterItem.ImportDetail.AgentImage,
		AgentProxy:  clusterItem.ImportDetail.AgentProxy,
		ClusterName: clusterItem.Name,
		ConnectType: clusterItem.ConnectType,
		MeshAddr:    clusterItem.MeshAddr,
	}
	return &cfg, nil
}

func (x *ClusterUseCase) GetClientConfigByClusterName(ctx context.Context, clusterName string) (*k8s.ClientConfig, error) {
	clusterItem, err := x.GetClusterByName(ctx, clusterName)
	if err != nil {
		return nil, fmt.Errorf("获取集群信息失败: %w", err)
	}
	cfg := k8s.ClientConfig{
		ImportType:  pb.ImportType(clusterItem.ImportType),
		KubeConfig:  clusterItem.ImportDetail.KubeConfig,
		Token:       clusterItem.ImportDetail.Token,
		APIServer:   clusterItem.ImportDetail.ApiServer,
		CaData:      clusterItem.ImportDetail.CaData,
		Context:     clusterItem.ImportDetail.Context,
		Agent:       clusterItem.ImportDetail.Agent,
		AgentImage:  clusterItem.ImportDetail.AgentImage,
		AgentProxy:  clusterItem.ImportDetail.AgentProxy,
		ClusterName: clusterName,
		ConnectType: clusterItem.ConnectType,
		MeshAddr:    clusterItem.MeshAddr,
	}
	return &cfg, nil
}

// ListClusters 查看集群列表
func (x *ClusterUseCase) ListClusters(ctx context.Context, query *QueryClusterReq) (ClusterItems, uint32, error) {
	clusters, err := x.clusterRepo.ListClusters(ctx, query)
	if err != nil {
		return nil, 0, err
	}
	total, err := x.clusterRepo.CountCluster(ctx, query)
	if err != nil {
		return nil, 0, err
	}
	return clusters, total, nil
}

// BuildMeshItem 构建MeshItem
func (x *ClusterUseCase) BuildMeshItem(ctx context.Context, data *ClusterItem) (*MeshItem, error) {
	// 查询目标Agent
	dstAgent, err := x.agent.GetAgent(ctx, &GetAgentRequest{Id: int(data.DstAgentId)})
	if err != nil {
		x.log.WithContext(ctx).Errorf("获取目标Agent失败: %v", err)
		return nil, fmt.Errorf("获取目标Agent失败: %v", err)
	}
	// 查询源Agent
	srcAgentId, err := x.agent.GetSrcAgentID(ctx)
	if err != nil {
		x.log.WithContext(ctx).Errorf("获取源Agent失败: %v", err)
		return nil, fmt.Errorf("获取源Agent失败: %v", err)
	}
	srcAgentPort, err := x.agent.GetSrcAgentPort(ctx, &GetAgentPortRequest{AgentId: srcAgentId})
	if err != nil {
		x.log.WithContext(ctx).Errorf("获取源Agent端口失败: %v", err)
		return nil, fmt.Errorf("获取源Agent端口失败: %v", err)
	}

	srcWhiteList, err := x.agent.GetSrcAgentWhiteList(ctx)
	if err != nil {
		x.log.WithContext(ctx).Errorf("获取源Agent白名单失败: %v", err)
		return nil, fmt.Errorf("获取源Agent白名单失败: %v", err)
	}
	// 获取目标Agent的apiServer地址, 只取IP+端口， 忽略协议

	dstServiceAddr := strings.Split(data.ImportDetail.ApiServer, "//")[1]
	return &MeshItem{
		ServiceName:    fmt.Sprintf("腾讯云test-%s", dstAgent.Name),
		WhiteIpList:    srcWhiteList,
		SrcAgentId:     srcAgentId,
		SrcAgentPort:   srcAgentPort,
		DstAgentId:     dstAgent.AgentId,
		DstServiceAddr: dstServiceAddr,
	}, nil
}

// CreateMeshAddr 创建Mesh
func (x *ClusterUseCase) CreateMeshAddr(ctx context.Context, meshItem *MeshItem) (string, error) {
	return x.mesh.CreateMesh(ctx, meshItem)
}

// CreateCluster 创建集群
func (x *ClusterUseCase) CreateCluster(ctx context.Context, data *ClusterItem) (uint32, error) {
	// 创建k8s client
	var meshAddr string
	switch data.ConnectType {
	case pb.ConnectType_Direct:
		meshAddr = ""
	case pb.ConnectType_Mesh:
		// 构建MeshItem
		meshItem, err := x.BuildMeshItem(ctx, data)
		if err != nil {
			return 0, fmt.Errorf("构建MeshItem失败: %v", err)
		}
		// 创建Mesh
		meshAddr, err = x.mesh.CreateMesh(ctx, meshItem)
		if err != nil {
			return 0, fmt.Errorf("创建Mesh失败: %v", err)
		}
		if meshAddr == "" {
			return 0, fmt.Errorf("创建Mesh失败: apiServer为空")
		}
		data.MeshAddr = meshAddr
	default:
		return 0, fmt.Errorf("未知的连接类型: %v", data.ConnectType)
	}

	cfg := k8s.ClientConfig{
		ImportType:  pb.ImportType(data.ImportType),
		KubeConfig:  data.ImportDetail.KubeConfig,
		Token:       data.ImportDetail.Token,
		APIServer:   data.ImportDetail.ApiServer,
		CaData:      data.ImportDetail.CaData,
		Context:     data.ImportDetail.Context,
		Agent:       data.ImportDetail.Agent,
		AgentImage:  data.ImportDetail.AgentImage,
		AgentProxy:  data.ImportDetail.AgentProxy,
		ClusterName: data.Name,
		ConnectType: data.ConnectType,
		MeshAddr:    meshAddr,
	}
	// 检查集群名称是否重复
	exist, err := x.clusterRepo.ExistCluster(ctx, data.Name)
	if err != nil {
		return 0, err
	}
	if exist {
		return 0, ErrClusterNameExist
	}

	data.UID = uuid.New().String()
	if data.ConnectType == pb.ConnectType_Direct {
		// 连接方式为直连时，请求k8s集群
		clientSet, err := cfg.CreateClientSet()
		if err != nil {
			return 0, err
		}

		// 获取k8s版本
		versionInfo, err := clientSet.Discovery().ServerVersion()
		if err != nil {
			return 0, fmt.Errorf("获取k8s版本失败: %v", err)
		}
		data.ServerVersion = versionInfo.String()
		data.Platform = versionInfo.Platform
		BuildDate := versionInfo.BuildDate
		buildDate, _ := utils.UTCToLocal(BuildDate)
		data.BuildDate = buildDate
	}

	// 创建集群
	id, err := x.clusterRepo.CreateCluster(ctx, data)
	if err != nil {
		return 0, err
	}
	// 发送创建集群事件
	_ = x.pubClusterCreative.Publish(ctx, &ClusterCreativeEvent{
		ClusterID: id,
	})
	return id, nil
}

// DeleteCluster 删除集群
func (x *ClusterUseCase) DeleteCluster(ctx context.Context, id uint32) error {
	// 获取集群信息用于清理缓存
	clusterItem, err := x.clusterRepo.GetClusterByID(ctx, id)
	if err == nil {
		// 清理缓存
		x.clearClientSetCache(clusterItem.Name)
	}
	err = x.clusterRepo.DeleteCluster(ctx, id)
	if err != nil {
		return fmt.Errorf("删除集群失败：%v", err)
	}
	err = x.pubClusterDeletion.Publish(ctx, &ClusterDeletionEvent{
		ClusterID: id,
	})
	return err
}

// FetchAllClusters 获取所有集群
func (x *ClusterUseCase) FetchAllClusters(ctx context.Context) (ClusterItems, error) {
	return x.clusterRepo.FetchAllClusters(ctx)
}

// UpdateCluster 更新集群
func (x *ClusterUseCase) UpdateCluster(ctx context.Context, data *ClusterItem) error {
	// 清理旧缓存，因为配置可能已经改变
	x.clearClientSetCache(data.Name)
	var meshAddr string
	// 不复用Mesh，重新创建Mesh
	switch data.ConnectType {
	case pb.ConnectType_Direct:
		meshAddr = ""
	case pb.ConnectType_Mesh:
		// 构建MeshItem
		meshItem, err := x.BuildMeshItem(ctx, data)
		if err != nil {
			return fmt.Errorf("构建MeshItem失败: %v", err)
		}
		// 创建Mesh
		meshAddr, err = x.mesh.CreateMesh(ctx, meshItem)
		if err != nil {
			return fmt.Errorf("创建Mesh失败: %v", err)
		}
		if meshAddr == "" {
			return fmt.Errorf("创建Mesh失败: apiServer为空")
		}
		data.MeshAddr = meshAddr
	default:

	}
	return x.clusterRepo.UpdateClusterV2(ctx, data)
}

// GetClusterByID 获取集群
func (x *ClusterUseCase) GetClusterByID(ctx context.Context, id uint32) (*ClusterItem, error) {
	return x.clusterRepo.GetClusterByID(ctx, id)
}

// GetClusterByName 获取集群
func (x *ClusterUseCase) GetClusterByName(ctx context.Context, name string) (*ClusterItem, error) {
	return x.clusterRepo.GetClusterByName(ctx, name)
}

func (x *ClusterUseCase) initializeClients(ctx context.Context, clusterID uint32) (*k8s.ClientConfig, *kubernetes.Clientset, *versioned.Clientset, error) {
	cfg, err := x.GetClientConfigByClusterID(ctx, clusterID)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("获取集群配置失败: %v", err)
	}

	// clientSet, err := cfg.CreateClientSet()
	clientSet, err := x.getOrCreateClientSet(ctx, cfg.ClusterName)
	if err != nil {
		return nil, nil, nil, err
	}

	// metricsClient, err := cfg.CreateMetricsClientSet()
	metricsClient, err := x.getOrCreateMetricsClient(ctx, cfg.ClusterName)
	if err != nil {
		return nil, nil, nil, err
	}

	return cfg, clientSet, metricsClient, nil
}

func (x *ClusterUseCase) processNodes(ctx context.Context, clusterID uint32, nodes *corev1.NodeList, metricsClient *versioned.Clientset) *ClusterMetrics {
	var wg sync.WaitGroup
	metrics := &ClusterMetrics{
		NodeCount: len(nodes.Items),
	}

	for i := range nodes.Items {
		wg.Add(1)
		go func(node *corev1.Node) {
			defer wg.Done()
			x.processNode(ctx, clusterID, node, metricsClient, metrics)
		}(&nodes.Items[i])
	}
	wg.Wait()
	return metrics
}

// HandleClusterDeletion 处理集群删除事件
func (x *ClusterUseCase) HandleClusterDeletion(ctx context.Context, clusterID uint32) error {
	return x.roleBinding.DeleteByClusterId(ctx, clusterID)
}

// HandleClusterCreative 处理集群创建事件
func (x *ClusterUseCase) HandleClusterCreative(ctx context.Context, clusterID uint32) error {
	cfg, clientSet, metricsClient, err := x.initializeClients(ctx, clusterID)
	if err != nil {
		return err
	}
	var healthStateList []HealthState
	isAPIServerHealthy, err := x.isAPIServerHealthy(ctx, cfg)
	if err != nil || !isAPIServerHealthy {
		healthStateList = append(healthStateList, HealthState(pb.HealthState_APIServerUnHealthy))
		return x.UpdateClusterState(ctx, &UpdateClusterStateRequest{
			Id:           clusterID,
			ClusterState: x.determineClusterState(healthStateList),
			HealthState:  healthStateList,
		})
	}

	versionInfo, err := x.getServerVersion(ctx, clientSet)
	if err != nil {
		return err
	}

	nodes, err := x.getNodes(ctx, clientSet)
	if err != nil {
		return err
	}

	clusterMetrics := x.processNodes(ctx, clusterID, nodes, metricsClient)

	healthStateList, err = x.checkClusterHealth(ctx, clientSet, cfg, clusterMetrics)
	if err != nil {
		return err
	}

	return x.updateCluster(ctx, clusterID, versionInfo, clusterMetrics, healthStateList)
}

// IsClusterHealthy 检查集群健康状态
func (x *ClusterUseCase) isClusterHealthy(ctx context.Context, cfg *k8s.ClientConfig, clientSet *kubernetes.Clientset) ([]HealthState, error) {
	healthStateList := make([]HealthState, 0)
	components, err := clientSet.CoreV1().ComponentStatuses().List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("获取组件状态失败: %v", err)
	}
	for _, component := range components.Items {
		if len(component.Conditions) == 0 {
			x.log.WithContext(ctx).Warnf("组件 %s 没有状态信息", component.Name)
			continue
		}

		if component.Conditions[0].Type == corev1.ComponentHealthy && component.Conditions[0].Status != corev1.ConditionTrue {
			switch component.Name {
			case "scheduler":
				healthStateList = append(healthStateList, HealthState(pb.HealthState_SchedulerUnhealthy))
				x.log.WithContext(ctx).Warnf("Scheduler 不健康: %s", component.Conditions[0].Message)
			case "controller-manager":
				healthStateList = append(healthStateList, HealthState(pb.HealthState_ControllerManagerUnhealthy))
				x.log.WithContext(ctx).Warnf("Controller Manager 不健康: %s", component.Conditions[0].Message)
			case "etcd-0":
				healthStateList = append(healthStateList, HealthState(pb.HealthState_EtcdUnhealthy))
				x.log.WithContext(ctx).Warnf("Etcd 不健康: %s", component.Conditions[0].Message)
			default:
				x.log.WithContext(ctx).Warnf("未知组件 %s 不健康: %s", component.Name, component.Conditions[0].Message)
			}
		} else {
			x.log.WithContext(ctx).Infof("组件 %s 健康", component.Name)
		}
	}
	// 检查APIServer是否健康
	apiServerHealthy, err := x.isAPIServerHealthy(ctx, cfg)
	if err != nil {
		x.log.WithContext(ctx).Warnf("检查 API Server 健康状态失败: %v", err)
		healthStateList = append(healthStateList, HealthState(pb.HealthState_APIServerUnHealthy))
	} else if !apiServerHealthy {
		healthStateList = append(healthStateList, HealthState(pb.HealthState_APIServerUnHealthy))
		x.log.WithContext(ctx).Warn("API Server 不健康")
	}

	// 检查节点是否ready
	nodeReady, err := x.checkNodeReady(ctx, clientSet)
	if err != nil {
		return nil, err
	}
	if !nodeReady {
		healthStateList = append(healthStateList, HealthState(pb.HealthState_NodeNotReady))
	}
	return healthStateList, nil
}

// CheckNodeReady 检查节点是否ready
func (x *ClusterUseCase) checkNodeReady(ctx context.Context, clientSet *kubernetes.Clientset) (bool, error) {
	// TODO: 优化节点列表获取, 目前只获取前100个节点
	nodes, err := clientSet.CoreV1().Nodes().List(ctx, metav1.ListOptions{Limit: 100})
	if err != nil {
		return false, fmt.Errorf("获取节点列表失败: %v", err)
	}
	for _, node := range nodes.Items {
		for _, condition := range node.Status.Conditions {
			if condition.Type == corev1.NodeReady && condition.Status != corev1.ConditionTrue {
				return false, nil
			}
		}
	}
	return true, nil
}

// IsAPIServerHealthy 检查APIServer健康状态
func (x *ClusterUseCase) isAPIServerHealthy(ctx context.Context, cfg *k8s.ClientConfig) (bool, error) {
	restConfig, err := cfg.BuildRestConfig()
	if err != nil {
		return false, fmt.Errorf("构建RestConfig失败: %v", err)
	}
	apiServerUrl := fmt.Sprintf("%s/healthz", restConfig.Host)
	// 创建一个自定义的http.Client，支持HTTPS
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true, // 跳过服务器证书验证
			},
			DisableKeepAlives: true, // 禁用长连接
		},
		Timeout: time.Second * 3, // 设置超时时间
	}
	resp, err := client.Get(apiServerUrl)
	if err != nil {
		return false, fmt.Errorf("APIServer健康状态检查失败: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return false, fmt.Errorf("APIServer健康状态检查失败: %s", resp.Status)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, fmt.Errorf("APIServer健康状态检查失败: %v", err)
	}
	if strings.ToLower(string(body)) != "ok" {
		return false, fmt.Errorf("APIServer健康状态检查失败: %s", string(body))
	}
	return true, nil
}

// getServerVersion 获取k8s版本
func (x *ClusterUseCase) getServerVersion(ctx context.Context, clientSet *kubernetes.Clientset) (*version.Info, error) {
	versionInfo, err := clientSet.Discovery().ServerVersion()
	if err != nil {
		return nil, fmt.Errorf("获取k8s版本失败: %v", err)
	}
	return versionInfo, nil
}

// getNodes 获取节点列表
func (x *ClusterUseCase) getNodes(ctx context.Context, clientSet *kubernetes.Clientset) (*corev1.NodeList, error) {
	nodes, err := clientSet.CoreV1().Nodes().List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("获取节点列表失败: %v", err)
	}
	return nodes, nil
}

type ClusterMetrics struct {
	CpuTotal    float32
	MemoryTotal float32
	CpuUsage    float32
	MemoryUsage float32
	NodeCount   int
}

func (m *ClusterMetrics) updateWithNodeMetrics(nodeMetrics NodeMetrics) {
	m.CpuTotal += nodeMetrics.CpuTotal
	m.MemoryTotal += nodeMetrics.MemoryTotal
	m.CpuUsage += nodeMetrics.CpuUsage
	m.MemoryUsage += nodeMetrics.MemoryUsage
}

type NodeMetrics struct {
	CpuTotal    float32
	MemoryTotal float32
	CpuUsage    float32
	MemoryUsage float32
}

// getMetrics 获取节点资源使用情况
func (x *ClusterUseCase) processNode(ctx context.Context, clusterID uint32, node *corev1.Node, metricsClient *versioned.Clientset, metrics *ClusterMetrics) {
	nodeMetrics := x.getNodeMetrics(ctx, node, metricsClient)
	nodeState, healthState := x.calculateNodeState(node, nodeMetrics)
	roles := x.getNodeRoles(node)

	x.updateNodeInRepo(ctx, clusterID, node, nodeMetrics, nodeState, healthState, roles)

	metrics.updateWithNodeMetrics(nodeMetrics)
}

func (x *ClusterUseCase) getNodeMetrics(ctx context.Context, node *corev1.Node, metricsClient *versioned.Clientset) NodeMetrics {
	metrics, err := metricsClient.MetricsV1beta1().NodeMetricses().Get(ctx, node.Name, metav1.GetOptions{})
	if err != nil {
		x.log.WithContext(ctx).Error(fmt.Sprintf("获取节点 %s 资源使用情况失败: %v", node.Name, err))
		return NodeMetrics{}
	}

	cpuUsage := metrics.Usage.Cpu().MilliValue()
	memoryUsage := utils.ConvertMemoryToGiB(*metrics.Usage.Memory())

	return NodeMetrics{
		CpuTotal:    float32(node.Status.Allocatable.Cpu().MilliValue()),
		MemoryTotal: float32(utils.ConvertMemoryToGiB(*node.Status.Allocatable.Memory())),
		CpuUsage:    float32(utils.ConvertCPUToCores(*resource.NewQuantity(cpuUsage, resource.DecimalSI))),
		MemoryUsage: float32(memoryUsage),
	}
}

func (x *ClusterUseCase) isResourceUsageHigh(usage, total, threshold int) bool {
	if total == 0 {
		return false
	}
	return (usage * 100 / total) > threshold
}

func (x *ClusterUseCase) determineClusterState(healthStateList []HealthState) ClusterState {
	for _, state := range healthStateList {
		if state != HealthState(pb.HealthState_Healthy) {
			return ClusterState(pb.ClusterState_ClusterError)
		}
	}
	return ClusterState(pb.ClusterState_ClusterReady)
}

func (x *ClusterUseCase) updateCluster(ctx context.Context, clusterID uint32, versionInfo *version.Info, metrics *ClusterMetrics, healthStateList []HealthState) error {
	buildDate, _ := utils.UTCToLocal(versionInfo.BuildDate)
	// 转为核保留两位小数
	cpuTotal := utils.ConvertCPUToCores(*resource.NewQuantity(int64(metrics.CpuTotal), resource.DecimalSI))
	return x.UpdateClusterBasicInfo(ctx, &UpdateClusterBasicRequest{
		Id:            clusterID,
		ServerVersion: versionInfo.String(),
		Platform:      versionInfo.Platform,
		BuildDate:     buildDate,
		CpuTotal:      float32(cpuTotal),
		MemoryTotal:   metrics.MemoryTotal,
		CpuUsage:      metrics.CpuUsage,
		MemoryUsage:   metrics.MemoryUsage,
		NodeCount:     metrics.NodeCount,
		ClusterState:  x.determineClusterState(healthStateList),
		HealthState:   healthStateList,
	})
}

func (x *ClusterUseCase) calculateNodeState(node *corev1.Node, metrics NodeMetrics) (NodeState, []NodeHealthState) {
	var nodeState NodeState
	healthStateList := make([]NodeHealthState, 0)

	for _, condition := range node.Status.Conditions {
		switch condition.Type {
		case corev1.NodeReady:
			if condition.Status == corev1.ConditionTrue {
				nodeState = NodeState(pb.NodeState_NodeReady)
			} else {
				nodeState = NodeState(pb.NodeState_NodeError)
			}
		case corev1.NodeMemoryPressure:
			if condition.Status != corev1.ConditionFalse {
				healthStateList = append(healthStateList, NodeHealthState(pb.NodeHealthState_MemoryPressure))
			}
		case corev1.NodeDiskPressure:
			if condition.Status != corev1.ConditionFalse {
				healthStateList = append(healthStateList, NodeHealthState(pb.NodeHealthState_DiskPressure))
			}
		case corev1.NodePIDPressure:
			if condition.Status != corev1.ConditionFalse {
				healthStateList = append(healthStateList, NodeHealthState(pb.NodeHealthState_PIDPressure))
			}
		case corev1.NodeNetworkUnavailable:
			if condition.Status != corev1.ConditionFalse {
				healthStateList = append(healthStateList, NodeHealthState(pb.NodeHealthState_NetworkUnavailable))
			}
		}
	}

	if x.isResourceUsageHigh(int(metrics.CpuUsage), int(metrics.CpuTotal), consts.CPUUsedThreshold) {
		healthStateList = append(healthStateList, NodeHealthState(pb.NodeHealthState_CpuUsageHigh))
		nodeState = NodeState(pb.NodeState_NodeError)
	}

	if x.isResourceUsageHigh(int(metrics.MemoryUsage), int(metrics.MemoryTotal), consts.MemoryUsedThreshold) {
		healthStateList = append(healthStateList, NodeHealthState(pb.NodeHealthState_MemoryUsageHigh))
		nodeState = NodeState(pb.NodeState_NodeError)
	}

	if len(healthStateList) == 0 {
		healthStateList = append(healthStateList, NodeHealthState(pb.NodeHealthState_Healthy))
	}

	return nodeState, healthStateList
}

func (x *ClusterUseCase) getNodeRoles(node *corev1.Node) []string {
	roles := make([]string, 0)
	nodeRolePrefix := "node-role.kubernetes.io/"
	for label := range node.Labels {
		if strings.HasPrefix(label, nodeRolePrefix) {
			role := strings.TrimPrefix(label, nodeRolePrefix)
			roles = append(roles, role)
		}
	}
	return roles
}

func (x *ClusterUseCase) updateNodeInRepo(ctx context.Context, clusterID uint32, node *corev1.Node, metrics NodeMetrics, nodeState NodeState, healthState []NodeHealthState, roles []string) {
	_, err := x.nodeRepo.CreateOrUpdateNode(ctx, &NodeItem{
		ClusterID:         clusterID,
		Name:              node.Name,
		Conditions:        node.Status.Conditions,
		Capacity:          node.Status.Capacity,
		Allocatable:       node.Status.Allocatable,
		Addresses:         node.Status.Addresses,
		CreationTimestamp: node.CreationTimestamp.Format(time.DateTime),
		CpuUsage:          metrics.CpuUsage,
		MemoryUsage:       metrics.MemoryUsage,
		Status:            nodeState,
		Labels:            node.Labels,
		NodeInfo:          node.Status.NodeInfo,
		UID:               string(node.GetUID()),
		ResourceVersion:   node.ResourceVersion,
		Spec:              node.Spec,
		Annotations:       node.Annotations,
		HealthState:       healthState,
		Roles:             roles,
	})
	if err != nil {
		x.log.WithContext(ctx).Debugf("保存节点数据失败, 节点名称: %s, 错误信息: %v", node.Name, err)
	}
}

func (x *ClusterUseCase) checkClusterHealth(ctx context.Context, clientSet *kubernetes.Clientset, cfg *k8s.ClientConfig, metrics *ClusterMetrics) ([]HealthState, error) {
	var healthStateList []HealthState

	if x.isResourceUsageHigh(int(metrics.CpuUsage), int(metrics.CpuTotal), consts.CPUUsedThreshold) {
		healthStateList = append(healthStateList, HealthState(pb.HealthState_CpuUsageHigh))
	}

	if x.isResourceUsageHigh(int(metrics.MemoryUsage), int(metrics.MemoryTotal), consts.MemoryUsedThreshold) {
		healthStateList = append(healthStateList, HealthState(pb.HealthState_MemoryUsageHigh))
	}

	if len(healthStateList) == 0 {
		healthStateList = []HealthState{HealthState(pb.HealthState_Healthy)}
	}

	return healthStateList, nil
}
