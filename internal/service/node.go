package service

import (
	"context"
	"fmt"
	"sort"

	"codo-cnmp/common/utils"
	"codo-cnmp/internal/biz"
	pb "codo-cnmp/pb"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
)

type NodeService struct {
	pb.UnimplementedNodeServer
	uc  *biz.NodeUseCase
	pod *biz.PodUseCase
	uf  *biz.UserFollowUseCase
}

func NewNodeService(uc *biz.NodeUseCase, pod *biz.PodUseCase, uf *biz.UserFollowUseCase) *NodeService {
	return &NodeService{
		uc:  uc,
		uf:  uf,
		pod: pod,
	}
}

// getUserFollowMap 获取用户关注的CloneSet列表
func (x *NodeService) getUserFollowMap(ctx context.Context, userID uint32) (map[string]bool, error) {
	userFollows, _, err := x.uf.ListUserFollow(ctx, &biz.ListUserFollowRequest{
		UserFollowCommonParams: biz.UserFollowCommonParams{
			UserID:     userID,
			FollowType: pb.FollowType_Node,
		},
		ListAll: true,
	})
	if err != nil {
		return nil, err
	}

	followMap := make(map[string]bool)
	for _, follow := range userFollows {
		followKey := fmt.Sprintf("%s.%s", follow.ClusterName, follow.FollowValue)
		followMap[followKey] = true
	}
	return followMap, nil
}

// setFollowedStatus 设置关注状态
func (x *NodeService) setFollowedStatus(ctx context.Context, clusterName string, items []*pb.NodeItem) error {
	userID, err := utils.GetUserIDFromCtx(ctx)
	if err != nil {
		return err
	}
	followMap, err := x.getUserFollowMap(ctx, userID)
	if err != nil {
		return err
	}
	for _, item := range items {
		followKey := fmt.Sprintf("%s.%s", clusterName, item.Name)
		item.IsFollowed = followMap[followKey]
	}
	return nil
}

func (x *NodeService) convertDO2DTO(node *biz.NodeItem) *pb.NodeItem {
	healthStateList := make([]pb.NodeHealthState, 0, len(node.HealthState))
	for _, item := range node.HealthState {
		healthStateList = append(healthStateList, pb.NodeHealthState(item))
	}
	addresses := node.Addresses
	var internalIp string
	for _, address := range addresses {
		if address.Type == corev1.NodeInternalIP {
			internalIp = address.Address
		}
	}
	if node.Roles == nil {
		node.Roles = make([]string, 0)
	}

	dto := &pb.NodeItem{
		Name:                    node.Name,
		State:                   pb.NodeStatus(node.Status),
		HealthState:             healthStateList,
		Roles:                   node.Roles,
		InternalIp:              internalIp,
		KubeletVersion:          node.NodeInfo.KernelVersion,
		OsImage:                 node.NodeInfo.OSImage,
		ContainerRuntimeVersion: node.NodeInfo.ContainerRuntimeVersion,
		CreateTime:              setTime(node.CreationTimestamp),
		Uid:                     node.UID,
		Unschedulable:           node.Spec.Unschedulable,
		PodTotal:                0, // 废弃, 使用 ResourceUsage 替代
		PodUsage:                0, // 废弃, 使用 ResourceUsage 替代
		IsFollowed:              false,
		OperatingSystem:         node.NodeInfo.OperatingSystem,
		ResourceUsage:           x.convertCapacity2ResourceUsage(node),
	}
	return dto
}

func (x *NodeService) convertNodeConditions(node *biz.NodeItem) []*pb.NodeCondition {
	nodeConditions := make([]*pb.NodeCondition, 0, len(node.Conditions))
	for _, condition := range node.Conditions {
		nodeConditions = append(nodeConditions, &pb.NodeCondition{
			Type:               string(condition.Type),
			Status:             string(condition.Status),
			Reason:             condition.Reason,
			Message:            condition.Message,
			LastHeartbeatTime:  uint64(condition.LastHeartbeatTime.UnixNano() / 1e6),
			LastTransitionTime: uint64(condition.LastTransitionTime.UnixNano() / 1e6),
		})
	}
	return nodeConditions

}

func (x *NodeService) convertCapacity2ResourceUsage(node *biz.NodeItem) *pb.ResourceUsage {
	formattedCpuCapacity := utils.ConvertCPUToCores(*resource.NewQuantity(node.Capacity.Cpu().MilliValue(), resource.DecimalSI))
	formattedMemoryCapacity := utils.ConvertMemoryToGiB(*node.Capacity.Memory())

	formattedAllocatedCpu := utils.ConvertCPUToCores(*resource.NewQuantity(node.Allocatable.Cpu().MilliValue(), resource.DecimalSI))
	formattedAllocatedMemory := utils.ConvertMemoryToGiB(*node.Allocatable.Memory())

	return &pb.ResourceUsage{
		NodeCpuUsage:          node.CpuUsage,
		NodeMemoryUsage:       node.MemoryUsage,
		NodeCpuTotal:          float32(formattedCpuCapacity),
		NodeMemoryTotal:       float32(formattedMemoryCapacity),
		PodTotal:              uint32(node.Allocatable.Pods().Value()),
		PodUsage:              node.NonTerminatedPods,
		NodeAllocatableCpu:    float32(formattedAllocatedCpu),
		NodeAllocatableMemory: float32(formattedAllocatedMemory),
	}

}

// ListNode 获取节点列表
func (x *NodeService) ListNode(ctx context.Context, req *pb.ListNodeRequest) (*pb.ListNodeResponse, error) {
	nodes, total, err := x.uc.ListNodesForAPI(ctx, &biz.ListNodeRequest{
		ClusterName: req.ClusterName,
		Page:        int32(req.Page),
		PageSize:    int32(req.PageSize),
		Keyword:     req.Keyword,
		ListAll:     utils.IntToBool(int(req.ListAll)),
	})
	if err != nil {
		return nil, err
	}
	list := make([]*pb.NodeItem, 0, len(nodes))
	for _, node := range nodes {
		nodeItem := x.convertDO2DTO(node)
		baseResourceUsage := nodeItem.ResourceUsage
		// 获取节点上的pod列表
		pods, _, err := x.pod.ListPods(ctx, &biz.ListPodRequest{
			ClusterName: req.ClusterName,
			NodeName:    node.Name,
			ListAll:     true,
		})
		if err != nil {
			return nil, err
		}
		detailedResourceUsage, err := x.calculateResourceUsage(pods, baseResourceUsage)
		if err != nil {
			return nil, err
		}
		nodeItem.ResourceUsage = detailedResourceUsage
		list = append(list, nodeItem)
	}
	if err := x.setFollowedStatus(ctx, req.ClusterName, list); err != nil {
		return nil, err
	}
	sort.SliceStable(list, func(i, j int) bool {
		return list[i].IsFollowed && !list[j].IsFollowed
	})
	return &pb.ListNodeResponse{
		List:  list,
		Total: total,
	}, nil
}

// addResourceList 累加资源列表
func addResourceList(total corev1.ResourceList, add corev1.ResourceList) {
	if total == nil {
		total = corev1.ResourceList{}
	}
	for name, quantity := range add {
		if _, exists := total[name]; !exists {
			total[name] = resource.Quantity{}
		}
		current := total[name]
		current.Add(quantity)
		total[name] = current
	}
}

// NodeResources 节点资源使用情况
type NodeResources struct {
	RequestCPU    *resource.Quantity
	RequestMemory *resource.Quantity
	LimitCPU      *resource.Quantity
	LimitMemory   *resource.Quantity
}

// calculateResourceUsage 计算节点资源使用情况
func (x *NodeService) calculateResourceUsage(pods []*corev1.Pod, baseResourceUsage *pb.ResourceUsage) (*pb.ResourceUsage, error) {

	var podUsage uint32
	// 初始化资源总量
	resources := &NodeResources{
		RequestCPU:    resource.NewQuantity(0, resource.DecimalSI),
		RequestMemory: resource.NewQuantity(0, resource.BinarySI),
		LimitCPU:      resource.NewQuantity(0, resource.DecimalSI),
		LimitMemory:   resource.NewQuantity(0, resource.BinarySI),
	}

	for _, pod := range pods {
		if pod.DeletionTimestamp == nil {
			podUsage++
		}
		for _, container := range pod.Spec.Containers {
			// 累加请求资源
			if cpu := container.Resources.Requests.Cpu(); cpu != nil {
				resources.RequestCPU.Add(*cpu)
			}
			if memory := container.Resources.Requests.Memory(); memory != nil {
				resources.RequestMemory.Add(*memory)
			}

			// 累加限制资源
			if cpu := container.Resources.Limits.Cpu(); cpu != nil {
				resources.LimitCPU.Add(*cpu)
			}
			if memory := container.Resources.Limits.Memory(); memory != nil {
				resources.LimitMemory.Add(*memory)
			}
		}
	}

	resourceUsage := baseResourceUsage
	resourceUsage.PodUsage = podUsage
	resourceUsage.RequestsCpuUsage = float32(utils.ConvertCPUToCores(*resource.NewQuantity(resources.RequestCPU.MilliValue(), resource.DecimalSI)))
	resourceUsage.RequestsCpuTotal = resourceUsage.NodeCpuTotal

	resourceUsage.RequestsMemoryUsage = float32(utils.ConvertMemoryToGiB(*resources.RequestMemory))
	resourceUsage.RequestsMemoryTotal = resourceUsage.NodeMemoryTotal

	resourceUsage.LimitsCpuUsage = float32(utils.ConvertCPUToCores(*resource.NewQuantity(resources.LimitCPU.MilliValue(), resource.DecimalSI)))
	resourceUsage.LimitsCpuTotal = resourceUsage.NodeCpuTotal

	resourceUsage.LimitsMemoryUsage = float32(utils.ConvertMemoryToGiB(*resources.LimitMemory))
	resourceUsage.LimitsMemoryTotal = resourceUsage.NodeMemoryTotal

	return resourceUsage, nil
}

// GetNodeDetail 获取节点详情
func (x *NodeService) GetNodeDetail(ctx context.Context, req *pb.GetNodeDetailRequest) (*pb.GetNodeDetailResponse, error) {
	node, err := x.uc.GetNodeDetail(ctx, &biz.ListNodeRequest{
		ClusterName: req.ClusterName,
		Name:        req.Name,
	})
	if err != nil {
		return nil, err
	}
	taints := make([]*corev1.Taint, 0, len(node.NodeItem.Spec.Taints))
	for _, taint := range node.NodeItem.Spec.Taints {
		taints = append(taints, &taint)
	}
	nodeItem := x.convertDO2DTO(&node.NodeItem)
	baseResourceUsage := x.convertCapacity2ResourceUsage(&node.NodeItem)
	conditions := x.convertNodeConditions(&node.NodeItem)

	// 获取节点上的pod列表
	pods, _, err := x.pod.ListPods(ctx, &biz.ListPodRequest{
		ClusterName: req.ClusterName,
		NodeName:    req.Name,
		ListAll:     true,
	})
	if err != nil {
		return nil, err
	}

	// 使用抽离的公共方法计算完整的资源使用情况
	resourceUsage, err := x.calculateResourceUsage(pods, baseResourceUsage)
	if err != nil {
		return nil, err
	}

	return &pb.GetNodeDetailResponse{
		Labels:        node.NodeItem.Labels,
		Annotations:   node.NodeItem.Annotations,
		Taints:        taints,
		NodeItem:      nodeItem,
		Pods:          pods,
		Yaml:          node.Yaml,
		SystemInfo:    &node.NodeItem.NodeInfo,
		ResourceUsage: resourceUsage,
		NodeCondition: conditions,
		Unschedulable: node.NodeItem.Spec.Unschedulable,
	}, nil

}

// UpdateNode 更新节点
func (x *NodeService) UpdateNode(ctx context.Context, req *pb.UpdateNodeRequest) (*pb.UpdateNodeResponse, error) {
	var spec corev1.NodeSpec
	taints := make([]corev1.Taint, 0)
	item := &biz.UpdateNodeRequest{
		ClusterName: req.ClusterName,
		NodeItem: biz.NodeItem{
			Name:        req.Name,
			Labels:      req.Labels,
			Annotations: req.Annotations,
		},
	}

	if len(req.Taints) > 0 {
		for _, data := range req.Taints {
			taints = append(taints, *data)
		}
		spec.Taints = taints
		item.NodeItem.Spec = spec
	} else {
		spec.Taints = taints
		item.NodeItem.Spec = spec
	}

	err := x.uc.UpdateNodeForAPI(ctx, item)
	if err != nil {
		return nil, err
	}
	return &pb.UpdateNodeResponse{}, nil
}

// CreateOrUpdateNodeByYaml 创建或更新节点
func (x *NodeService) CreateOrUpdateNodeByYaml(ctx context.Context, req *pb.CreateOrUpdateNodeByYamlRequest) (*pb.CreateOrUpdateNodeByYamlResponse, error) {
	err := x.uc.CreateOrUpdateNodeByYaml(ctx, &biz.CreateOrUpdateNodeByYamlRequest{
		ClusterName: req.ClusterName,
		Yaml:        req.Yaml,
	})
	if err != nil {
		return nil, err
	}
	return &pb.CreateOrUpdateNodeByYamlResponse{}, nil
}

// HandleNode 节点操作
func (x *NodeService) HandleNode(ctx context.Context, req *pb.HandleNodeRequest) (*pb.HandleNodeResponse, error) {
	result, err := x.uc.HandleNode(ctx, &biz.HandleNodeRequest{
		ClusterName: req.ClusterName,
		Name:        req.Name,
		Operation:   req.Operation,
	})
	if err != nil {
		return nil, err
	}
	return &pb.HandleNodeResponse{
		Success: result,
	}, nil
}

// CheckEviction 节点驱逐检查
func (x *NodeService) CheckEviction(ctx context.Context, req *pb.CheckEvictionRequest) (*pb.CheckEvictionResponse, error) {
	pods, _, err := x.pod.ListPods(ctx, &biz.ListPodRequest{
		ClusterName: req.ClusterName,
		NodeName:    req.Name,
		ListAll:     true,
	})
	if err != nil {
		return nil, err
	}
	result := &pb.CheckEvictionItem{
		ReadyToEvictPodsCount: 0,
		IgnoreEvictPodsCount:  0,
	}
	for _, pod := range pods {
		isDaemonSetPod := false
		for _, ownerRef := range pod.OwnerReferences {
			if ownerRef.Kind == "DaemonSet" {
				isDaemonSetPod = true
				break
			}
		}
		if isDaemonSetPod {
			result.IgnoreEvictPodsCount++
		} else {
			result.ReadyToEvictPodsCount++
		}
	}
	return &pb.CheckEvictionResponse{Detail: result}, nil
}
