package service

import (
	"context"
	"fmt"
	"sort"
	"strconv"

	"codo-cnmp/common/consts"
	"codo-cnmp/common/utils"
	"codo-cnmp/internal/biz"
	"codo-cnmp/internal/event"
	pb "codo-cnmp/pb"
	"github.com/ccheers/xpkg/sync/errgroup"
)

// ClusterService implements the Cluster interface.
type ClusterService struct {
	pb.UnimplementedClusterServer
	uc           *biz.ClusterUseCase
	oc           *biz.OverViewUseCase
	uf           *biz.UserFollowUseCase
	eventManager *event.ClusterEventManager
}

// NewClusterService creates a new ClusterService object.
func NewClusterService(uc *biz.ClusterUseCase, oc *biz.OverViewUseCase, uf *biz.UserFollowUseCase, eventManager *event.ClusterEventManager) *ClusterService {
	return &ClusterService{uc: uc, oc: oc, uf: uf, eventManager: eventManager}
}

// getUserFollowMap 获取用户关注的集群映射
func (x *ClusterService) getUserFollowMap(ctx context.Context, userID uint32) (map[string]bool, error) {
	userFollows, _, err := x.uf.ListUserFollow(ctx, &biz.ListUserFollowRequest{
		UserFollowCommonParams: biz.UserFollowCommonParams{
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

func IsSuperUser(ctx context.Context) bool {
	superUserValueFlag := ctx.Value(consts.ContextSuperUserKey)
	if superUserValueFlag != nil {
		if superUser, ok := superUserValueFlag.(bool); ok && superUser {
			return true
		}
	}
	return false
}

// 过滤用户有权限的集群
func (x *ClusterService) filterUserPermissionCluster(ctx context.Context, items []*pb.ClusterItem) ([]*pb.ClusterItem, error) {
	// 如果是超级用户，则不进行过滤
	if IsSuperUser(ctx) {
		return items, nil
	}
	clusterRoleBindingsValue := ctx.Value(consts.ContextClusterRoleBindingsKey)
	if clusterRoleBindingsValue == nil {
		return nil, nil
	}
	clusterBindings, ok := clusterRoleBindingsValue.([]map[string][]string)
	if !ok {
		return nil, nil
	}
	if len(clusterBindings) == 0 {
		return nil, nil
	}
	// 创建授权集群名称的快速查找map
	authorizedClusters := make(map[string]struct{})
	for _, binding := range clusterBindings {
		for clusterName := range binding {
			authorizedClusters[clusterName] = struct{}{}
		}
	}

	// 过滤有权限的集群
	result := make([]*pb.ClusterItem, 0, len(items))
	for _, item := range items {
		if _, hasPermission := authorizedClusters[item.Name]; hasPermission {
			result = append(result, item)
		}
	}

	return result, nil
}

// setFollowedStatus 设置关注状态
func (x *ClusterService) setFollowedStatus(ctx context.Context, items []*pb.ClusterItem) error {
	userID, err := utils.GetUserIDFromCtx(ctx)
	if err != nil {
		return err
	}
	followMap, err := x.getUserFollowMap(ctx, userID)
	if err != nil {
		return err
	}
	for _, item := range items {
		item.IsFollowed = followMap[strconv.Itoa(int(item.Id))]
	}
	return nil
}

func (x *ClusterService) convertDO2DTO(cluster *biz.ClusterItem) *pb.ClusterItem {
	healthStateList := make([]pb.HealthState, 0, len(cluster.HealthState))
	for _, item := range cluster.HealthState {
		healthStateList = append(healthStateList, pb.HealthState(item))
	}
	ops := make([]string, 0, len(cluster.Ops))
	for _, item := range cluster.Ops {
		ops = append(ops, item)
	}
	dto := &pb.ClusterItem{
		Id:            cluster.ID,
		ImportType:    pb.ImportType(cluster.ImportType),
		Name:          cluster.Name,
		Description:   cluster.Description,
		ClusterState:  pb.ClusterState(cluster.ClusterState),
		HealthState:   healthStateList,
		BuildDate:     setTime(cluster.BuildDate),
		ServerVersion: cluster.ServerVersion,
		CpuUsage:      cluster.CpuUsage,
		MemoryUsage:   cluster.MemoryUsage,
		CpuTotal:      cluster.CpuTotal,
		MemoryTotal:   cluster.MemoryTotal,
		NodeCount:     uint32(cluster.NodeCount),
		Uid:           cluster.UID,
		Idip:          cluster.IDIP,
		IsFollowed:    cluster.IsFollowed,
		AppId:         cluster.AppId,
		AppSecret:     cluster.AppSecret,
		Ops:           ops,
		DstAgentId:    cluster.DstAgentId,
		ConnectType:   cluster.ConnectType,
		Links:         cluster.Links,
	}
	return dto
}

// ListCluster 集群列表
func (x *ClusterService) ListCluster(ctx context.Context, request *pb.ListClusterRequest) (*pb.ListClusterResponse, error) {
	clusters, count, err := x.uc.ListClusters(ctx, &biz.QueryClusterReq{
		Page:     request.Page,
		PageSize: request.PageSize,
		Keyword:  request.Keyword,
		ListAll:  utils.IntToBool(int(request.ListAll)),
	})
	if err != nil {
		return nil, err
	}

	list := make([]*pb.ClusterItem, 0, len(clusters))
	for _, cluster := range clusters {
		list = append(list, x.convertDO2DTO(cluster))
	}

	if err := x.setFollowedStatus(ctx, list); err != nil {
		return nil, err
	}
	sort.SliceStable(list, func(i, j int) bool {
		return list[i].IsFollowed && !list[j].IsFollowed
	})

	// 过滤用户有权限的集群
	if request.AuthFilter != nil {
		authFilter := utils.IntToBool(int(*request.AuthFilter))
		if authFilter {
			filteredList, err := x.filterUserPermissionCluster(ctx, list)
			if err != nil {
				return nil, err
			}
			list = filteredList
			count = uint32(len(filteredList))
		}
	}

	return &pb.ListClusterResponse{
		List:  list,
		Total: count,
	}, nil
}

// CreateCluster 创建集群
func (x *ClusterService) CreateCluster(ctx context.Context, request *pb.ImportClusterRequest) (*pb.ImportClusterResponse, error) {
	id, err := x.uc.CreateCluster(ctx, &biz.ClusterItem{
		Name:        request.Name,
		Description: request.Description,
		ImportType:  biz.ImportType(request.ImportType),
		ImportDetail: biz.ImportDetail{
			KubeConfig: request.ImportDetail.KubeConfig,
			Token:      request.ImportDetail.Token,
			ApiServer:  request.ImportDetail.ApiServer,
			CaData:     request.ImportDetail.CaData,
			Context:    request.ImportDetail.Context,
			Agent:      request.ImportDetail.Agent,
			AgentImage: request.ImportDetail.AgentImage,
			AgentProxy: request.ImportDetail.AgentProxy,
		},
		IDIP:        request.Idip,
		AppId:       request.AppId,
		AppSecret:   request.AppSecret,
		Ops:         request.Ops,
		DstAgentId:  request.DstAgentId,
		ConnectType: request.ConnectType,
		Links:       request.GetLinks(),
	})
	if err != nil {
		return nil, err
	}
	var eg errgroup.Group
	// 创建一个新的 context，确保父 context 不会被取消
	taskCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	// 异步处理集群添加事件，不等待任务完成
	eg.Go(func(ctx context.Context) error {
		if err := x.eventManager.OnClusterAdd(taskCtx, request.Name); err != nil {
			// 错误处理可以记录日志，但不影响主流程
			return fmt.Errorf("处理集群添加事件失败: %v", err)
		}
		return nil
	})
	//if err := x.eventManager.OnClusterAdd(ctx, request.Name); err != nil {
	//	return nil, fmt.Errorf("处理集群添加事件失败: %v", err)
	//}
	return &pb.ImportClusterResponse{
		Id: id,
	}, nil
}

// DeleteCluster 删除集群
func (x *ClusterService) DeleteCluster(ctx context.Context, request *pb.DeleteClusterRequest) (*pb.DeleteClusterResponse, error) {
	clusterObj, err := x.uc.GetClusterByID(ctx, request.Id)
	if err != nil {
		return nil, err
	}
	if err := x.eventManager.OnClusterDelete(ctx, clusterObj.Name); err != nil {
		return nil, fmt.Errorf("处理集群删除事件失败: %v", err)
	}
	err = x.uc.DeleteCluster(ctx, request.Id)
	if err != nil {
		return nil, err
	}
	return &pb.DeleteClusterResponse{}, nil
}

// convertErrorCluster
func (x *ClusterService) convertErrorClusters(errorClusters []biz.ClusterItem) []*pb.ClusterItem {
	list := make([]*pb.ClusterItem, 0, len(errorClusters))
	for _, cluster := range errorClusters {
		list = append(list, x.convertDO2DTO(&cluster))
	}
	return list
}

// convertFollowedClusters
func (x *ClusterService) convertFollowedClusters(followedClusters []biz.ClusterItem) []*pb.ClusterItem {
	// 获取所有集群
	list := make([]*pb.ClusterItem, 0, len(followedClusters))
	for _, cluster := range followedClusters {
		list = append(list, x.convertDO2DTO(&cluster))
	}
	return list
}

func (x *ClusterService) convertToOverview(data *biz.OverviewResponse) *pb.OverView {
	return &pb.OverView{
		Cluster: &pb.ClusterOverview{
			Total:        data.ClusterTotal.Cluster.Total,
			RunningTotal: data.ClusterTotal.Cluster.RunningTotal,
			ErrorTotal:   data.ClusterTotal.Cluster.ErrorTotal,
			OfflineTotal: data.ClusterTotal.Cluster.OfflineTotal,
		},
		Node: &pb.NodeOverview{
			Total:        data.ClusterTotal.Node.Total,
			RunningTotal: data.ClusterTotal.Node.RunningTotal,
			ErrorTotal:   data.ClusterTotal.Node.ErrorTotal,
		},
		Cpu: &pb.CpuOverview{
			Total:            data.ClusterTotal.CPU.Total,
			UnallocatedTotal: data.ClusterTotal.CPU.UnallocatedTotal,
			AllocatedTotal:   data.ClusterTotal.CPU.AllocatedTotal,
		},
		Memory: &pb.MemoryOverview{
			Total:            data.ClusterTotal.Memory.Total,
			UnallocatedTotal: data.ClusterTotal.Memory.UnallocatedTotal,
			AllocatedTotal:   data.ClusterTotal.Memory.AllocatedTotal,
		},
	}
}

// OverviewCluster 集群概览
func (x *ClusterService) OverviewCluster(ctx context.Context, request *pb.ClusterOverviewRequest) (*pb.ClusterOverviewResponse, error) {
	// 获取概览数据
	overViewResp, err := x.oc.OverView(ctx)
	if err != nil {
		return nil, err
	}

	// 获取关注的集群
	followedClusters := x.convertFollowedClusters(overViewResp.FollowedClusterItems)

	// 获取异常集群
	errorClusters := x.convertErrorClusters(overViewResp.ErrorClusterItems)

	// 构建响应
	return &pb.ClusterOverviewResponse{
		Overview: x.convertToOverview(&overViewResp),
		Follow: &pb.FollowClusterItem{
			List: followedClusters,
		},
		Error: &pb.ErrorClusterItem{
			List: errorClusters,
		},
	}, nil

}

func (x *ClusterService) UpdateCluster(ctx context.Context, request *pb.UpdateClusterRequest) (*pb.UpdateClusterResponse, error) {
	links := request.GetLinks()
	if links == nil {
		links = make([]*pb.Link, 0)
	}
	ops := request.GetOps()
	if ops == nil {
		ops = make([]string, 0)
	}
	err := x.uc.UpdateCluster(ctx, &biz.ClusterItem{
		ID:          request.Id,
		Name:        request.Name,
		Description: request.Description,
		IDIP:        request.Idip,
		AppId:       request.AppId,
		AppSecret:   request.AppSecret,
		Ops:         ops,
		DstAgentId:  request.DstAgentId,
		ConnectType: request.ConnectType,
		Links:       links,
		ImportDetail: biz.ImportDetail{
			KubeConfig: request.ImportDetail.KubeConfig,
			Token:      request.ImportDetail.Token,
			ApiServer:  request.ImportDetail.ApiServer,
			CaData:     request.ImportDetail.CaData,
			Context:    request.ImportDetail.Context,
			Agent:      request.ImportDetail.Agent,
			AgentImage: request.ImportDetail.AgentImage,
			AgentProxy: request.ImportDetail.AgentProxy,
		},
	})
	if err != nil {
		return nil, err
	}
	if err := x.eventManager.OnClusterAdd(ctx, request.Name); err != nil {
		return nil, fmt.Errorf("处理集群添加事件失败: %v", err)
	}
	return &pb.UpdateClusterResponse{}, nil
}

func (x *ClusterService) PingIdip(ctx context.Context, request *pb.PingIdipRequest) (*pb.PingIdipResponse, error) {
	result, err := x.uc.PingIdip(ctx, request.Idip)
	if err != nil {
		return nil, err
	}
	return &pb.PingIdipResponse{
		Connected: result,
	}, nil
}

func (x *ClusterService) GetClusterDetail(ctx context.Context, request *pb.GetClusterDetailRequest) (*pb.GetClusterDetailResponse, error) {
	cluster, err := x.uc.GetClusterByID(ctx, request.Id)
	if err != nil {
		return nil, err
	}
	dto := x.convertDO2DTO(cluster)
	dto.ImportDetail = x.convertImportDetail(cluster.ImportDetail)
	return &pb.GetClusterDetailResponse{
		Detail: dto,
	}, nil
}

func (x *ClusterService) convertLinks(links []biz.Link) []*pb.Link {
	list := make([]*pb.Link, 0, len(links))
	for _, link := range links {
		list = append(list, &pb.Link{
			Name: link.Name,
			Url:  link.Url,
		})
	}
	return list
}

func (x *ClusterService) convertImportDetail(detail biz.ImportDetail) *pb.ImportDetail {
	return &pb.ImportDetail{
		KubeConfig: detail.KubeConfig,
		Token:      detail.Token,
		ApiServer:  detail.ApiServer,
		CaData:     detail.CaData,
		Context:    detail.Context,
		Agent:      detail.Agent,
		AgentImage: detail.AgentImage,
		AgentProxy: detail.AgentProxy,
	}
}

func (x *ClusterService) DownloadKubeConfig(ctx context.Context, request *pb.DownloadKubeConfigRequest) (*pb.DownloadKubeConfigResponse, error) {
	kubeConfig, err := x.uc.DownLoadKubeConfig(ctx, &biz.DownLoadKubeConfigRequest{
		Id: request.Id,
	})
	if err != nil {
		return nil, err
	}
	return &pb.DownloadKubeConfigResponse{
		KubeConfig: kubeConfig,
	}, nil
}
