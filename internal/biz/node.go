package biz

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"time"

	"codo-cnmp/common/utils"
	pb "codo-cnmp/pb"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-redis/redis/v8"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/yaml"
)

type NodeHealthState pb.NodeHealthState
type NodeState pb.NodeState
type Operation = pb.NodeOperation

type ListNodeRequest struct {
	ID          uint32 `json:"id"`
	Page        int32  `json:"page"`
	PageSize    int32  `json:"page_size"`
	Keyword     string `json:"keyword"`
	ClusterID   uint32 `json:"cluster_id"`
	ClusterName string `json:"cluster_name"`
	Name        string `json:"name"`
	ListAll     bool   `json:"list_all"`
}

type NodeItem struct {
	ID                uint32                 `json:"id"`
	ClusterID         uint32                 `json:"cluster_id"`
	Name              string                 `json:"name"`
	UID               string                 `json:"uid"`
	ResourceVersion   string                 `json:"resource_version"`
	Labels            map[string]string      `json:"labels"`
	Annotations       map[string]string      `json:"annotations"`
	Conditions        []corev1.NodeCondition `json:"conditions"`
	Capacity          corev1.ResourceList    `json:"capacity"`
	Allocatable       corev1.ResourceList    `json:"allocatable"`
	Addresses         []corev1.NodeAddress   `json:"addresses"`
	CreationTimestamp string                 `json:"creation_timestamp"`
	CpuUsage          float32                `json:"cpu_usage"`
	MemoryUsage       float32                `json:"memory_usage"`
	Status            NodeState              `json:"status"`
	Roles             []string               `json:"roles"`
	NodeInfo          corev1.NodeSystemInfo  `json:"system_info"`
	HealthState       []NodeHealthState      `json:"health_state"`
	Spec              corev1.NodeSpec        `json:"spec"`
	NonTerminatedPods uint32                 `json:"non_terminated_pods"`
}

type NodeDetail struct {
	NodeItem NodeItem      `json:"node_item"`
	PodList  []*corev1.Pod `json:"pod_list"`
	Yaml     string        `json:"yaml"`
}

type CreateNodeRequest struct {
	ClusterID uint32   `json:"cluster_id"`
	NodeItem  NodeItem `json:"node_item"`
}

type UpdateNodeRequest struct {
	ClusterName string   `json:"cluster_name"`
	NodeItem    NodeItem `json:"node_item"`
}

type CreateOrUpdateNodeByYamlRequest struct {
	ClusterName string `json:"cluster_name"`
	Yaml        string `json:"yaml"`
}

type HandleNodeRequest struct {
	ClusterName string `json:"cluster_name"`
	Name        string `json:"name"`
	Operation   Operation
}

type CheckEvictionRequest struct {
	ClusterName string `json:"cluster_name"`
	Name        string `json:"name"`
}

type CheckEvictionResponse struct {
	// 等待驱逐pod数量
	ReadyToEvictCount uint32 `json:"ready_to_evict_count"`
	// 忽略驱逐pod数量
	IgnoreEvictCount uint32 `json:"ignore_evict_count"`
}

type NodeRepo interface {
	// ListNodes 查看节点列表
	ListNodes(ctx context.Context, req *ListNodeRequest) ([]*NodeItem, error)

	// CreateNode 创建节点
	CreateNode(ctx context.Context, req *CreateNodeRequest) (uint32, error)

	// DeleteNode 删除节点
	DeleteNode(ctx context.Context, id uint32) error

	// UpdateNode 更新节点
	UpdateNode(ctx context.Context, req *NodeItem) error

	// ExistNode 判断节点是否存在
	ExistNode(ctx context.Context, name string, clusterID uint32) (bool, error)

	// CountNode 统计节点数量
	CountNode(ctx context.Context, query *ListNodeRequest) (uint32, error)

	// CreateOrUpdateNode 创建或者更新节点
	CreateOrUpdateNode(ctx context.Context, req *NodeItem) (uint32, error)

	// GetNode 获取节点
	GetNode(ctx context.Context, req *ListNodeRequest) (*NodeItem, error)
}

// INodeUseCase 节点业务逻辑接口
type INodeUseCase interface {
	// ListNodes 查看节点列表
	ListNodes(ctx context.Context, req *ListNodeRequest) ([]*NodeItem, uint32, error)
	// CreateNode 创建节点
	CreateNode(ctx context.Context, data *CreateNodeRequest) (uint32, error)
	// ExistNode 节点是否存在
	ExistNode(ctx context.Context, name string, clusterID uint32) (bool, error)
	// UpdateNode 更新节点
	UpdateNode(ctx context.Context, req *NodeItem) error
	// CreateOrUpdateNode 创建或者更新节点
	CreateOrUpdateNode(ctx context.Context, req *NodeItem) (uint32, error)
	// GetNodeDetail 获取节点详情
	GetNodeDetail(ctx context.Context, req *ListNodeRequest) (*NodeDetail, error)
	// UpdateNodeForAPI 更新节点
	UpdateNodeForAPI(ctx context.Context, req *UpdateNodeRequest) error
	// CreateOrUpdateNodeByYaml 使用Yaml创建或者更新节点
	CreateOrUpdateNodeByYaml(ctx context.Context, req *CreateOrUpdateNodeByYamlRequest) error
	// CordonNode 封锁节点
	CordonNode(ctx context.Context, req *HandleNodeRequest) (bool, error)
	// DrainNode 驱逐节点
	DrainNode(ctx context.Context, req *HandleNodeRequest) (bool, error)
	// HandleNode 处理节点
	HandleNode(ctx context.Context, req *HandleNodeRequest) (bool, error)
	// ListNodesForAPI 查看节点列表
	ListNodesForAPI(ctx context.Context, req *ListNodeRequest) ([]*NodeItem, uint32, error)
	// GetNodeForAPI 获取节点
	GetNodeForAPI(ctx context.Context, req *ListNodeRequest) (*NodeItem, error)

	// SyncNodePods 同步节点上的Pods信息
	SyncNodePods(ctx context.Context) error
	// GetNodePods 获取集群节点上的Pods信息
	GetNodePods(ctx context.Context, clusterID uint32) (map[string][]corev1.Pod, error)
}

func NewINodeUseCase(x *NodeUseCase) INodeUseCase {
	return x
}

type NodeUseCase struct {
	cluster IClusterUseCase
	pod     IPodUseCase
	repo    NodeRepo
	log     *log.Helper
	redis   *redis.Client
}

// NewNodeUseCase 创建 NodeUseCase
func NewNodeUseCase(repo NodeRepo, cluster IClusterUseCase, pod IPodUseCase, logger log.Logger, redis *redis.Client) *NodeUseCase {
	x := &NodeUseCase{
		repo:    repo,
		cluster: cluster,
		pod:     pod,
		log:     log.NewHelper(log.With(logger, "module", "biz/node")),
		redis:   redis,
	}
	return x
}

// ListNodes 查看节点列表
func (x *NodeUseCase) ListNodes(ctx context.Context, req *ListNodeRequest) ([]*NodeItem, uint32, error) {
	cluster, err := x.cluster.GetClusterByName(ctx, req.ClusterName)
	if err != nil {
		return nil, 0, err
	}
	req.ClusterID = cluster.ID
	nodes, err := x.repo.ListNodes(ctx, req)
	if err != nil {
		return nil, 0, err
	}
	req.ListAll = true
	total, err := x.repo.CountNode(ctx, req)
	if err != nil {
		return nil, 0, err
	}
	return nodes, total, nil
}

// CreateNode 创建节点
func (x *NodeUseCase) CreateNode(ctx context.Context, data *CreateNodeRequest) (uint32, error) {
	// 创建节点
	return x.repo.CreateNode(ctx, data)
}

// ExistNode 判断节点是否存在
func (x *NodeUseCase) ExistNode(ctx context.Context, name string, clusterID uint32) (bool, error) {
	return x.repo.ExistNode(ctx, name, clusterID)
}

// UpdateNode 更新节点
func (x *NodeUseCase) UpdateNode(ctx context.Context, req *NodeItem) error {
	return x.repo.UpdateNode(ctx, req)
}

// CreateOrUpdateNode 创建或者更新节点
func (x *NodeUseCase) CreateOrUpdateNode(ctx context.Context, req *NodeItem) (uint32, error) {
	return x.repo.CreateOrUpdateNode(ctx, req)
}

// GetNodeDetail 获取节点详情
func (x *NodeUseCase) GetNodeDetail(ctx context.Context, req *ListNodeRequest) (*NodeDetail, error) {
	node, err := x.GetNodeForAPI(ctx, req)
	if err != nil {
		return nil, err
	}
	yamlStr, err := x.GetNodeYaml(ctx, req)
	if err != nil {
		return nil, err
	}
	nodeDetail := &NodeDetail{
		NodeItem: *node,
		Yaml:     yamlStr,
		PodList:  []*corev1.Pod{},
	}
	return nodeDetail, nil
}

func (x *NodeUseCase) UpdateNodeForAPI(ctx context.Context, req *UpdateNodeRequest) error {
	clientSet, err := x.cluster.GetClientSetByClusterName(ctx, req.ClusterName)
	if err != nil {
		return fmt.Errorf("创建k8s client失败: %w", err)
	}
	node, err := clientSet.CoreV1().Nodes().Get(ctx, req.NodeItem.Name, metav1.GetOptions{})
	if err != nil {
		x.log.WithContext(ctx).Errorf("查询节点失败: %v", err)
		return fmt.Errorf("查询节点失败: %w", err)
	}
	// 更新标签
	if req.NodeItem.Labels != nil {
		node.Labels = req.NodeItem.Labels
	}
	// 更新注解
	if req.NodeItem.Annotations != nil {
		node.Annotations = req.NodeItem.Annotations
	}
	// 更新污点
	if req.NodeItem.Spec.Taints != nil {
		node.Spec.Taints = req.NodeItem.Spec.Taints
	}
	_, err = clientSet.CoreV1().Nodes().Update(ctx, node, metav1.UpdateOptions{})
	if err != nil {
		x.log.WithContext(ctx).Errorf("更新节点失败: %v", err)
		return fmt.Errorf("更新节点失败: %w", err)
	}
	return nil
}

// CreateOrUpdateNodeByYaml 创建或者更新节点
func (x *NodeUseCase) CreateOrUpdateNodeByYaml(ctx context.Context, req *CreateOrUpdateNodeByYamlRequest) error {
	clientSet, err := x.cluster.GetClientSetByClusterName(ctx, req.ClusterName)
	if err != nil {
		return fmt.Errorf("创建k8s client失败: %w", err)
	}
	var node corev1.Node
	decoder := yaml.NewYAMLOrJSONDecoder(strings.NewReader(req.Yaml), 4096)
	if err := decoder.Decode(&node); err != nil {
		return fmt.Errorf("解析yaml失败: %v", err)
	}
	// 获取Node
	_, err = clientSet.CoreV1().Nodes().Get(ctx, node.Name, metav1.GetOptions{})
	if errors.IsNotFound(err) {
		// 创建Node
		_, err := clientSet.CoreV1().Nodes().Create(ctx, &node, metav1.CreateOptions{})
		if err != nil {
			x.log.WithContext(ctx).Errorf("创建节点失败: %v", err)
			return fmt.Errorf("创建节点失败: %w", err)
		}
	}
	if err != nil {
		x.log.WithContext(ctx).Errorf("查询节点失败: %v", err)
		return fmt.Errorf("查询节点失败: %w", err)
	}
	// 更新Node
	_, err = clientSet.CoreV1().Nodes().Update(ctx, &node, metav1.UpdateOptions{})
	if err != nil {
		x.log.WithContext(ctx).Errorf("更新节点失败: %v", err)
		return fmt.Errorf("更新节点失败: %w", err)
	}
	return nil

}

// GetNodeYaml 调用K8s API 获取节点yaml
func (x *NodeUseCase) GetNodeYaml(ctx context.Context, req *ListNodeRequest) (string, error) {
	var yamlStr string
	clientSet, err := x.cluster.GetClientSetByClusterName(ctx, req.ClusterName)
	if err != nil {
		return yamlStr, fmt.Errorf("创建k8s client失败: %w", err)
	}
	node, err := clientSet.CoreV1().Nodes().Get(ctx, req.Name, metav1.GetOptions{})
	if err != nil {
		return yamlStr, fmt.Errorf("查询节点失败: %v", err)
	}
	if node.Kind == "" {
		node.Kind = "Node"
	}
	if node.APIVersion == "" {
		node.APIVersion = "v1"
	}
	yamlStr, err = utils.ConvertResourceTOYaml(node)
	return yamlStr, err
}

func (x *NodeUseCase) CordonNode(ctx context.Context, req *HandleNodeRequest) (bool, error) {
	clientSet, err := x.cluster.GetClientSetByClusterName(ctx, req.ClusterName)
	if err != nil {
		return false, fmt.Errorf("创建k8s client失败: %w", err)
	}
	node, err := clientSet.CoreV1().Nodes().Get(ctx, req.Name, metav1.GetOptions{})
	if err != nil {
		x.log.WithContext(ctx).Errorf("查询节点失败: %v", err)
		return false, fmt.Errorf("查询节点失败: %w", err)
	}
	switch req.Operation {
	case pb.NodeOperation_NodeCordon:
		// 封锁
		node.Spec.Unschedulable = true
	case pb.NodeOperation_NodeUncordon:
		// 解封
		node.Spec.Unschedulable = false
	default:
		return false, fmt.Errorf("不支持的操作: %v", req.Operation)
	}
	_, err = clientSet.CoreV1().Nodes().Update(ctx, node, metav1.UpdateOptions{})
	if err != nil {
		x.log.WithContext(ctx).Errorf("更新节点失败, 操作: %v, 错误: %v", req.Operation, err)
		return false, fmt.Errorf("更新节点失败: %w", err)
	}
	return true, nil
}

// DrainNode 驱逐节点
func (x *NodeUseCase) DrainNode(ctx context.Context, req *HandleNodeRequest) (bool, error) {
	clientSet, err := x.cluster.GetClientSetByClusterName(ctx, req.ClusterName)
	if err != nil {
		return false, fmt.Errorf("创建k8s client失败: %w", err)
	}
	node, err := clientSet.CoreV1().Nodes().Get(ctx, req.Name, metav1.GetOptions{})
	if err != nil {
		x.log.WithContext(ctx).Errorf("查询节点失败: %v", err)
		return false, fmt.Errorf("查询节点失败: %w", err)
	}
	// 查询当前节点所有pods
	pods, _, err := x.pod.ListPods(ctx, &ListPodRequest{
		ClusterName: req.ClusterName,
		NodeName:    node.Name,
		ListAll:     true,
	})
	if err != nil {
		return false, err
	}
	for _, pod := range pods {
		isDaemonSet := false
		for _, owner := range pod.OwnerReferences {
			if owner.Kind == "DaemonSet" {
				isDaemonSet = true
				break
			}
		}
		if isDaemonSet {
			continue
		}
		_, err = x.pod.EvictPod(ctx, &EvictPodRequest{
			PodName:     pod.Name,
			Namespace:   pod.Namespace,
			ClusterName: req.ClusterName,
		})
		if err != nil {
			return false, err
		}
	}
	return true, nil
}

// HandleNode 处理节点
func (x *NodeUseCase) HandleNode(ctx context.Context, req *HandleNodeRequest) (bool, error) {
	switch req.Operation {
	case pb.NodeOperation_NodeCordon:
		// 封锁节点
		return x.CordonNode(ctx, req)
	case pb.NodeOperation_NodeUncordon:
		// 解封节点
		return x.CordonNode(ctx, req)
	case pb.NodeOperation_NodeDrain:
		// 驱逐节点
		// 先封锁
		req.Operation = pb.NodeOperation_NodeCordon
		if _, err := x.CordonNode(ctx, req); err != nil {
			return false, err
		}
		// 再驱逐
		req.Operation = pb.NodeOperation_NodeDrain
		return x.DrainNode(ctx, req)
	case pb.NodeOperation_NodeUndrain:
		return false, fmt.Errorf("不支持的操作: %v", req.Operation)
	default:
		return false, fmt.Errorf("不支持的操作: %v", req.Operation)
	}
}

// ListNodesForAPI 查看节点列表 调用K8s API 获取节点列表
func (x *NodeUseCase) ListNodesForAPI(ctx context.Context, req *ListNodeRequest) ([]*NodeItem, uint32, error) {
	clientSet, err := x.cluster.GetClientSetByClusterName(ctx, req.ClusterName)
	if err != nil {
		x.log.WithContext(ctx).Errorf("创建k8s client失败: %v", err)
		return nil, 0, fmt.Errorf("创建k8s client失败: %w", err)
	}
	if req.Page == 0 {
		req.Page = 1
	}
	if req.PageSize == 0 {
		req.PageSize = 10
	}

	var (
		allFilteredNodes = make([]*NodeItem, 0)
		limit            = int64(req.PageSize)
		continueToken    = ""
	)

	cluster, err := x.cluster.GetClusterByName(ctx, req.ClusterName)
	if err != nil {
		return nil, 0, fmt.Errorf("获取集群失败: %w", err)
	}

	clusterNodePods, err := x.GetNodePods(ctx, cluster.ID)
	if err != nil {
		x.log.WithContext(ctx).Errorf("获取节点Pods失败: %v, custerID=%d", err, cluster.ID)
	}

	for {
		nodeListOptions := metav1.ListOptions{
			Limit:    limit,
			Continue: continueToken,
		}
		nodes, err := clientSet.CoreV1().Nodes().List(ctx, nodeListOptions)
		if err != nil {
			x.log.WithContext(ctx).Errorf("查询nods失败: %v", err)
			return nil, 0, fmt.Errorf("查询nods失败: %w", err)
		}
		filteredNodes := filterNodes(nodes, req.Keyword)
		for _, node := range filteredNodes.Items {
			nodeItem := &NodeItem{
				ID:                0,
				ClusterID:         0,
				Name:              node.Name,
				UID:               string(node.UID),
				ResourceVersion:   node.ResourceVersion,
				Labels:            node.Labels,
				Annotations:       node.Annotations,
				Conditions:        node.Status.Conditions,
				Capacity:          node.Status.Capacity,
				Allocatable:       node.Status.Allocatable,
				Addresses:         node.Status.Addresses,
				CreationTimestamp: node.CreationTimestamp.Format(time.DateTime),
				CpuUsage:          0,
				MemoryUsage:       0,
				Status:            0,
				Roles:             nil,
				NodeInfo:          node.Status.NodeInfo,
				HealthState:       nil,
				Spec:              node.Spec,
				NonTerminatedPods: 0,
			}
			// 补充自定义的状态和健康状态
			nodeObj, err := x.repo.GetNode(ctx, &ListNodeRequest{
				Name:      node.Name,
				ClusterID: req.ClusterID,
			})
			if err != nil {
				x.log.WithContext(ctx).Error("查询节点失败", err)
				continue
			} else {
				nodeItem.ID = nodeObj.ID
				nodeItem.Status = nodeObj.Status
				nodeItem.HealthState = nodeObj.HealthState
				nodeItem.CpuUsage = nodeObj.CpuUsage
				nodeItem.MemoryUsage = nodeObj.MemoryUsage
			}
			nodePods, ok := clusterNodePods[nodeItem.Name]
			if ok {
				nodeItem.NonTerminatedPods = uint32(len(nodePods))
			}
			allFilteredNodes = append(allFilteredNodes, nodeItem)
		}

		if nodes.Continue == "" {
			break
		}
		// 更新 continueToken，继续获取下一页
		continueToken = nodes.Continue
	}

	// 如果 ListAll 为 true，直接返回所有结果
	if req.ListAll {
		return allFilteredNodes, uint32(len(allFilteredNodes)), nil
	}
	if len(allFilteredNodes) == 0 {
		return allFilteredNodes, 0, nil
	}
	paginatedNodes, total := utils.K8sPaginate(allFilteredNodes, uint32(req.Page), uint32(req.PageSize))
	return paginatedNodes, total, nil
}

// filterNodes 过滤节点
func filterNodes(nodes *corev1.NodeList, keyword string) *corev1.NodeList {
	result := &corev1.NodeList{}
	// 编译正则表达式，(?i)表示忽略大小写
	pattern := "(?i).*" + regexp.QuoteMeta(keyword) + ".*"
	for _, node := range nodes.Items {
		if utils.MatchString(pattern, node.Name) {
			result.Items = append(result.Items, node)
		}
	}
	return result
}

// GetNodeForAPI 获取节点信息 调用K8s API 获取节点信息
func (x *NodeUseCase) GetNodeForAPI(ctx context.Context, req *ListNodeRequest) (*NodeItem, error) {
	// todo 分页
	clientSet, err := x.cluster.GetClientSetByClusterName(ctx, req.ClusterName)
	if err != nil {
		return nil, fmt.Errorf("创建k8s client失败: %w", err)
	}
	node, err := clientSet.CoreV1().Nodes().Get(ctx, req.Name, metav1.GetOptions{})
	if err != nil {
		x.log.WithContext(ctx).Errorf("查询节点失败: %v", err)
		return nil, fmt.Errorf("查询节点失败: %w", err)
	}
	nodeItem := &NodeItem{
		Name:              node.Name,
		UID:               string(node.UID),
		Spec:              node.Spec,
		Labels:            node.Labels,
		Annotations:       node.Annotations,
		Conditions:        node.Status.Conditions,
		Addresses:         node.Status.Addresses,
		Capacity:          node.Status.Capacity,
		Allocatable:       node.Status.Allocatable,
		NodeInfo:          node.Status.NodeInfo,
		ResourceVersion:   node.ResourceVersion,
		CreationTimestamp: node.CreationTimestamp.Format(time.DateTime),
	}
	// 补充自定义的状态和健康状态
	nodeObj, err := x.repo.GetNode(ctx, &ListNodeRequest{
		Name: node.Name,
	})
	if err != nil {
		x.log.WithContext(ctx).Errorf("查询节点失败: %v", err)
		return nil, fmt.Errorf("查询节点失败: %w", err)
	}
	nodeItem.ID = nodeObj.ID
	nodeItem.Status = nodeObj.Status
	nodeItem.HealthState = nodeObj.HealthState
	nodeItem.CpuUsage = nodeObj.CpuUsage
	nodeItem.MemoryUsage = nodeObj.MemoryUsage
	return nodeItem, nil
}

func (x *NodeUseCase) SyncNodePods(ctx context.Context) error {
	const (
		lockKey = "codo:cnmp:sync_node_pods:lock"
	)

	ok, err := x.redis.SetNX(ctx, lockKey, "1", time.Minute*5).Result()
	if !ok {
		x.log.Warnf("获取锁失败, 锁已被占用, 锁key: %s", lockKey)
		return nil
	}
	defer x.redis.Del(ctx, lockKey)
	clusters, err := x.cluster.FetchAllClusters(ctx)
	if err != nil {
		return fmt.Errorf("获取集群失败: %w", err)
	}

	for _, cluster := range clusters {
		clientSet, err := x.cluster.GetClientSetByClusterName(ctx, cluster.Name)
		if err != nil {
			return fmt.Errorf("创建k8s client失败: %w", err)
		}
		nodeList, err := clientSet.CoreV1().Nodes().List(ctx, metav1.ListOptions{})
		if err != nil {
			return fmt.Errorf("查询节点失败: %w", err)
		}

		mm := make(map[string][]corev1.Pod)
		for _, node := range nodeList.Items {
			pods, err := clientSet.CoreV1().Pods("").List(ctx, metav1.ListOptions{
				FieldSelector: fmt.Sprintf("spec.nodeName=%s", node.Name),
			})
			if err != nil {
				x.log.Errorf("查询节点下的Pods失败: %v", err)
				continue
			}
			mm[node.Name] = pods.Items
		}
		bs, _ := json.Marshal(mm)
		x.redis.Set(ctx, fmt.Sprintf("codo:cnmp:cluster_nodes:%d", cluster.ID), bs, time.Hour)
	}
	return nil
}

func (x *NodeUseCase) GetNodePods(ctx context.Context, clusterID uint32) (map[string][]corev1.Pod, error) {
	key := fmt.Sprintf("codo:cnmp:cluster_nodes:%d", clusterID)
	val, err := x.redis.Get(ctx, key).Bytes()
	if err != nil {
		return nil, fmt.Errorf("获取节点Pods失败: %w", err)
	}
	var mm map[string][]corev1.Pod
	if err := json.Unmarshal(val, &mm); err != nil {
		return nil, fmt.Errorf("解析节点Pods失败: %w", err)
	}
	return mm, nil
}
