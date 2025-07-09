package biz

import (
	"codo-cnmp/common/utils"
	"context"
	"math"
	"strconv"

	"codo-cnmp/pb"
	"github.com/go-kratos/kratos/v2/log"
)

// ClusterOverview 集群概览统计信息
type ClusterOverview struct {
	Total        uint32 `json:"total"`        // 集群总数
	RunningTotal uint32 `json:"runningTotal"` // 运行中的集群数量
	ErrorTotal   uint32 `json:"errorTotal"`   // 异常状态的集群数量
	OfflineTotal uint32 `json:"offlineTotal"` // 离线状态的集群数量
}

// NodeOverview 节点概览统计信息
type NodeOverview struct {
	Total        uint32 `json:"total"`        // 节点总数
	RunningTotal uint32 `json:"runningTotal"` // 运行中的节点数量
	ErrorTotal   uint32 `json:"errorTotal"`   // 异常状态的节点数量
}

// ResourceOverview 资源概览基础结构
type ResourceOverview struct {
	Total            float32 `json:"total"`            // 资源总量
	UnallocatedTotal float32 `json:"unallocatedTotal"` // 未分配资源量
	AllocatedTotal   float32 `json:"allocatedTotal"`   // 已分配资源量
}

// CpuOverview CPU资源概览
type CpuOverview ResourceOverview

// MemoryOverview 内存资源概览
type MemoryOverview ResourceOverview

// Overview 系统整体概览信息
type Overview struct {
	Cluster ClusterOverview `json:"cluster"` // 集群统计信息
	Node    NodeOverview    `json:"node"`    // 节点统计信息
	CPU     CpuOverview     `json:"cpu"`     // CPU资源统计信息
	Memory  MemoryOverview  `json:"memory"`  // 内存资源统计信息
}

// OverviewResponse 概览接口响应结构
type OverviewResponse struct {
	ClusterTotal         Overview      `json:"clusterTotal"`         // 集群整体统计信息
	FollowedClusterItems []ClusterItem `json:"myFollowClusterItems"` // 用户关注的集群列表
	ErrorClusterItems    []ClusterItem `json:"errorClusterItems"`    // 异常状态的集群列表
}

type IOverViewUseCase interface {
	// OverView 集群概览
	OverView(ctx context.Context) (OverviewResponse, error)
}

type OverViewUseCase struct {
	cluster *ClusterUseCase
	node    *NodeUseCase
	uf      *UserFollowUseCase
	log     *log.Helper
}

func NewOverViewUseCase(cluster *ClusterUseCase, node *NodeUseCase, uf *UserFollowUseCase, logger log.Logger) *OverViewUseCase {
	return &OverViewUseCase{
		cluster: cluster,
		node:    node,
		uf:      uf,
		log:     log.NewHelper(log.With(logger, "module", "biz/overview")),
	}
}

func NewIOverViewUseCase(x *OverViewUseCase) IOverViewUseCase {
	return x
}

// getUserFollowMap 获取用户关注的集群映射
func (x *OverViewUseCase) getUserFollowMap(ctx context.Context, userID uint32) (map[string]bool, error) {
	userFollows, _, err := x.uf.ListUserFollow(ctx, &ListUserFollowRequest{
		UserFollowCommonParams: UserFollowCommonParams{
			UserID:     userID,
			FollowType: pb.FollowType_Cluster,
		},
		ListAll: true,
	})
	if err != nil {
		return nil, err
	}

	followMap := make(map[string]bool)
	for _, follow := range userFollows {
		followMap[follow.FollowValue] = true
	}
	return followMap, nil
}

// enrichClusterWithFollowStatus 为集群列表添加关注状态
func (x *OverViewUseCase) enrichClusterWithFollowStatus(items ClusterItems, followMap map[string]bool) {
	for _, item := range items {
		item.IsFollowed = followMap[strconv.Itoa(int(item.ID))]
	}
}

func (x *OverViewUseCase) getFollowedClusters(ctx context.Context) ([]ClusterItem, error) {
	// 获取所有集群
	clusters, _, err := x.cluster.ListClusters(ctx, &QueryClusterReq{
		ListAll: true,
	})
	if err != nil {
		return nil, err
	}

	// 获取用户ID
	userID, err := utils.GetUserIDFromCtx(ctx)
	if err != nil {
		return nil, err
	}

	// 获取用户关注映射
	followMap, err := x.getUserFollowMap(ctx, userID)
	if err != nil {
		return nil, err
	}

	// 筛选出已关注的集群
	followedClusters := make([]ClusterItem, 0)
	for _, cluster := range clusters {
		if followMap[strconv.Itoa(int(cluster.ID))] {
			cluster.IsFollowed = true
			followedClusters = append(followedClusters, *cluster)
		}
	}

	return followedClusters, nil
}

// OverView 集群概览
func (x *OverViewUseCase) OverView(ctx context.Context) (OverviewResponse, error) {
	var overview OverviewResponse
	// 获取集群列表
	clusterItems, err := x.cluster.FetchAllClusters(ctx)
	if err != nil {
		return overview, err
	}

	var (
		stats struct {
			cluster struct {
				total, running, error, offline int
			}
			node struct {
				total, running, error int
			}
			resources struct {
				cpuTotal, cpuAllocated float32
				memTotal, memAllocated float32
			}
		}
		errorClusterItems []ClusterItem
		followedClusters  []ClusterItem
	)
	for _, cluster := range clusterItems {
		stats.cluster.total++
		stats.node.total += cluster.NodeCount
		stats.resources.cpuTotal += cluster.CpuTotal
		stats.resources.cpuAllocated += cluster.CpuUsage
		stats.resources.memTotal += cluster.MemoryTotal
		stats.resources.memAllocated += cluster.MemoryUsage

		// 统计运行中、异常、离线
		for _, state := range cluster.HealthState {
			switch state {
			case HealthState(pb.HealthState_Healthy):
				stats.cluster.running++
			case HealthState(pb.HealthState_APIServerUnHealthy):
				stats.cluster.offline++
			default:
				stats.cluster.error++
				errorClusterItems = append(errorClusterItems, *cluster)
			}
			break
		}

		// 查询Nodes
		nodes, _, err := x.node.ListNodes(ctx, &ListNodeRequest{
			ClusterID:   cluster.ID,
			ListAll:     true,
			ClusterName: cluster.Name,
		})
		if err != nil {
			return overview, err
		}
		for _, node := range nodes {
			switch node.Status {
			case NodeState(pb.NodeState_NodeReady):
				stats.node.running++
			case NodeState(pb.NodeState_NodeError):
				stats.node.error++
			}
		}

	}

	// 辅助函数：四舍五入到小数点后2位
	roundToTwoDecimals := func(value float64) float32 {
		return float32(math.Round(value*100) / 100)
	}

	// 计算资源使用情况
	cpuUnallocated := roundToTwoDecimals(float64(stats.resources.cpuTotal - stats.resources.cpuAllocated))
	memTotal := roundToTwoDecimals(float64(stats.resources.memTotal))
	memUnallocated := roundToTwoDecimals(float64(stats.resources.memTotal - stats.resources.memAllocated))
	memAllocated := roundToTwoDecimals(float64(stats.resources.memAllocated))

	// 获取用户关注的集群列表
	followedClusters, err = x.getFollowedClusters(ctx)
	if err != nil {
		return overview, err
	}

	// 构建返回结果
	return OverviewResponse{
		ClusterTotal: Overview{
			Cluster: ClusterOverview{
				Total:        uint32(stats.cluster.total),
				RunningTotal: uint32(stats.cluster.running),
				ErrorTotal:   uint32(stats.cluster.error),
				OfflineTotal: uint32(stats.cluster.offline),
			},
			Node: NodeOverview{
				Total:        uint32(stats.node.total),
				RunningTotal: uint32(stats.node.running),
				ErrorTotal:   uint32(stats.node.error),
			},
			CPU: CpuOverview{
				Total:            stats.resources.cpuTotal,
				UnallocatedTotal: cpuUnallocated,
				AllocatedTotal:   stats.resources.cpuAllocated,
			},
			Memory: MemoryOverview{
				Total:            memTotal,
				UnallocatedTotal: memUnallocated,
				AllocatedTotal:   memAllocated,
			},
		},
		ErrorClusterItems:    errorClusterItems,
		FollowedClusterItems: followedClusters,
	}, nil
}
