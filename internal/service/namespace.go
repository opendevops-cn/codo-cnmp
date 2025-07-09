package service

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"strconv"

	"codo-cnmp/common/consts"
	"codo-cnmp/common/utils"
	"codo-cnmp/internal/biz"
	pb "codo-cnmp/pb"
	"github.com/go-redis/redis/v8"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
)

type NameSpaceService struct {
	pb.UnimplementedNameSpaceServer
	uc      *biz.NameSpaceUseCase
	uf      *biz.UserFollowUseCase
	redis   *redis.Client
	rb      *biz.RoleBindingUseCase
	cluster *biz.ClusterUseCase
	role    *biz.RoleUseCase
}

func NewNameSpaceService(uc *biz.NameSpaceUseCase, uf *biz.UserFollowUseCase, rb *biz.RoleBindingUseCase, redis *redis.Client, cluster *biz.ClusterUseCase, role *biz.RoleUseCase) *NameSpaceService {
	return &NameSpaceService{
		uc:      uc,
		uf:      uf,
		redis:   redis,
		rb:      rb,
		cluster: cluster,
		role:    role,
	}
}

// getUserFollowMap 获取用户关注的CloneSet列表
func (x *NameSpaceService) getUserFollowMap(ctx context.Context, userID uint32) (map[string]bool, error) {
	userFollows, _, err := x.uf.ListUserFollow(ctx, &biz.ListUserFollowRequest{
		UserFollowCommonParams: biz.UserFollowCommonParams{
			UserID:     userID,
			FollowType: pb.FollowType_Namespace,
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
func (x *NameSpaceService) setFollowedStatus(ctx context.Context, clusterName string, items []*pb.NameSpaceItem) error {
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

// isDefaultNameSpace 检查命名空间是否为默认命名空间
func (x *NameSpaceService) isDefaultNameSpace(name string) bool {
	return biz.IsDefaultNameSpaces(name)
}

func (x *NameSpaceService) convertDO2DTO(ns *corev1.Namespace) *pb.NameSpaceItem {
	if ns.APIVersion == "" {
		ns.APIVersion = "v1"
	}
	if ns.Kind == "" {
		ns.Kind = "Namespace"
	}
	yamlStr, _ := utils.ConvertResourceTOYaml(ns)
	dto := &pb.NameSpaceItem{
		Name:        ns.Name,
		Description: ns.GetAnnotations()["description"],
		CreateTime:  uint64(ns.CreationTimestamp.Time.UnixNano() / 1e6),
		State:       string(ns.Status.Phase),
		Uid:         string(ns.UID),
		IsDefault:   x.isDefaultNameSpace(ns.Name),
		Labels:      ns.Labels,
		Annotations: ns.Annotations,
		Yaml:        yamlStr,
	}
	return dto
}

func (x *NameSpaceService) CreateNameSpace(ctx context.Context, req *pb.CreateNameSpaceRequest) (*pb.CreateNameSpaceResponse, error) {
	err := x.uc.CreateNameSpace(ctx, &biz.CreateNameSpaceRequest{
		ClusterName: req.ClusterName,
		Name:        req.Name,
		Description: req.Description,
		Labels:      req.Labels,
		Annotations: req.Annotations,
	})
	if err != nil {
		return nil, err
	}

	return &pb.CreateNameSpaceResponse{}, nil
}

// GetUserGroupIds 获取用户组ID
func (x *NameSpaceService) GetUserGroupIds(ctx context.Context) ([]uint32, error) {
	userID, err := utils.GetUserIDFromCtx(ctx)
	if err != nil {
		return nil, err
	}
	// 将用户ID转换为字符串
	userId := strconv.FormatUint(uint64(userID), 10)
	var userGroupIds []uint32
	key := fmt.Sprintf(consts.UserGroupIDCacheKey, userId)
	result, err := x.redis.Get(ctx, key).Result()
	if err != nil {
		return nil, fmt.Errorf("获取用户组ID失败: %w", err)
	}
	err = json.Unmarshal([]byte(result), &userGroupIds)
	if err != nil {
		return nil, fmt.Errorf("解析用户组ID失败: %w", err)
	}
	return userGroupIds, nil
}

// GetRoleBindings 获取角色绑定
func (x *NameSpaceService) GetRoleBindings(ctx context.Context, userGroupIds []uint32) ([]*biz.RoleBindingItem, error) {
	var items []*biz.RoleBindingItem
	for _, userGroupId := range userGroupIds {
		roleBindingItems, _, err := x.rb.List(ctx, &biz.ListRoleBindingRequest{
			RoleBindingCommonParams: biz.RoleBindingCommonParams{
				UserGroupID: userGroupId,
			},
		})
		if err != nil {
			return nil, fmt.Errorf("获取角色失败: %w", err)
		}
		items = append(items, roleBindingItems...)

	}
	return items, nil
}

// SetRoles 设置角色
func (x *NameSpaceService) SetRoles(ctx context.Context, clusterName string, userGroupIds []uint32, list []*pb.NameSpaceItem) ([]*pb.NameSpaceItem, error) {
	roleBindings, err := x.GetRoleBindings(ctx, userGroupIds)
	if err != nil {
		return nil, err
	}

	// 预先获取所有需要的 cluster 信息
	clusterCache := make(map[uint32]*biz.ClusterItem)
	for _, rb := range roleBindings {
		if _, ok := clusterCache[rb.ClusterID]; !ok {
			clusterItem, err := x.cluster.GetClusterByID(ctx, rb.ClusterID)
			if err != nil {
				continue
			}
			clusterCache[rb.ClusterID] = clusterItem
		}
	}

	// 预先获取所有需要的角色信息
	roleCache := make(map[uint32]*pb.RoleItem)
	roleIDSet := make(map[uint32]struct{})
	for _, rb := range roleBindings {
		roleIDSet[rb.RoleID] = struct{}{}
	}
	for roleID := range roleIDSet {
		roleItem, err := x.role.GetRole(ctx, roleID)
		if err != nil {
			continue
		}
		roleCache[roleID] = &pb.RoleItem{
			Id:          roleItem.ID,
			Name:        roleItem.Name,
			Description: roleItem.Description,
			CreateTime:  setTime(roleItem.CreateTime),
			UpdateTime:  setTime(roleItem.UpdateTime),
			IsDefault:   roleItem.ISDefault,
			Yaml:        roleItem.YamlStr,
			UpdateBy:    roleItem.UpdateBy,
			RoleType:    pb.RoleType(roleItem.RoleType),
		}
	}

	// 处理每个 namespace 的角色
	for _, item := range list {
		roleMap := make(map[uint32]struct{})
		for _, rb := range roleBindings {
			cluster, ok := clusterCache[rb.ClusterID]
			if !ok {
				continue
			}
			if cluster.Name == clusterName && (rb.Namespace == item.Name || rb.Namespace == "all" || rb.Namespace == "*") {
				roleMap[rb.RoleID] = struct{}{}
			}
		}

		roleItems := make([]*pb.RoleItem, 0, len(roleMap))
		for roleID := range roleMap {
			if role, ok := roleCache[roleID]; ok {
				roleItems = append(roleItems, role)
			}
		}
		item.Roles = roleItems
	}
	return list, nil
}

// filterUserPermissionNameSpace 过滤用户无权限的命名空间
func (x *NameSpaceService) filterUserPermissionNameSpace(ctx context.Context, clusterName string, list []*pb.NameSpaceItem) ([]*pb.NameSpaceItem, error) {
	// 如果是超级用户，则不进行过滤
	if IsSuperUser(ctx) {
		return list, nil
	}
	clusterRoleBindingsValue := ctx.Value(consts.ContextClusterRoleBindingsKey)
	if clusterRoleBindingsValue == nil {
		return nil, nil
	}
	clusterBindings, ok := clusterRoleBindingsValue.([]map[string][]string)
	if !ok || len(clusterBindings) == 0 {
		return nil, nil
	}

	// 先创建一个现有namespace的映射，避免后续重复遍历list
	existingNamespaces := make(map[string]*pb.NameSpaceItem, len(list))
	for _, item := range list {
		existingNamespaces[item.Name] = item
	}

	// 创建一个map来记录有权限的namespace
	authorizedNamespaces := make(map[string]struct{})

	for _, binding := range clusterBindings {
		namespaces, exists := binding[clusterName]
		if !exists {
			continue
		}

		// 检查是否有通配符权限或特定命名空间权限
		if utils.Contains(namespaces, "*") {
			return list, nil
		}

		// 记录所有有权限的namespace
		for _, ns := range namespaces {
			if _, exists := existingNamespaces[ns]; exists {
				authorizedNamespaces[ns] = struct{}{}
			}
		}
	}

	// 过滤结果
	result := make([]*pb.NameSpaceItem, 0, len(authorizedNamespaces))
	for _, item := range list {
		if _, ok := authorizedNamespaces[item.Name]; ok {
			result = append(result, item)
		}
	}

	return result, nil
}

// ListNameSpace 获取命名空间列表
func (x *NameSpaceService) ListNameSpace(ctx context.Context, req *pb.ListNameSpaceRequest) (*pb.ListNameSpaceResponse, error) {
	namespaces, total, err := x.uc.ListNameSpace(ctx, &biz.ListNameSpaceRequest{
		ClusterName: req.ClusterName,
		Keyword:     req.Keyword,
		Page:        req.Page,
		PageSize:    req.PageSize,
		ListAll:     utils.IntToBool(int(req.ListAll)),
	})
	if errors.IsNotFound(err) {
		return &pb.ListNameSpaceResponse{
			List:  []*pb.NameSpaceItem{},
			Total: 0,
		}, nil
	}
	if err != nil {
		return nil, err
	}
	list := make([]*pb.NameSpaceItem, 0)
	for _, ns := range namespaces {
		list = append(list, x.convertDO2DTO(ns))
	}
	if err := x.setFollowedStatus(ctx, req.ClusterName, list); err != nil {
		return nil, err
	}
	sort.SliceStable(list, func(i, j int) bool {
		return list[i].IsFollowed && !list[j].IsFollowed
	})
	if req.AuthFilter != nil {
		authFilter := utils.IntToBool(int(*req.AuthFilter))
		if authFilter {
			list, err = x.filterUserPermissionNameSpace(ctx, req.ClusterName, list)
			if err != nil {
				return nil, err
			}
			total = uint32(len(list))
		}
	}
	//
	userGroupIds, err := x.GetUserGroupIds(ctx)
	if err == nil {
		list, err = x.SetRoles(ctx, req.ClusterName, userGroupIds, list)
		if err != nil {
			return nil, err
		}
	}

	return &pb.ListNameSpaceResponse{
		List:  list,
		Total: total,
	}, nil
}

func (x *NameSpaceService) DeleteNameSpace(ctx context.Context, req *pb.DeleteNameSpaceRequest) (*pb.DeleteNameSpaceResponse, error) {
	err := x.uc.DeleteNameSpace(ctx, &biz.DeleteNameSpaceRequest{
		ClusterName: req.ClusterName,
		Name:        req.Name,
	})
	if err != nil {
		return nil, err
	}
	return &pb.DeleteNameSpaceResponse{}, nil
}

func (x *NameSpaceService) UpdateNameSpace(ctx context.Context, req *pb.UpdateNameSpaceRequest) (*pb.DeleteNameSpaceResponse, error) {
	err := x.uc.UpdateNameSpace(ctx, &biz.UpdateNameSpaceRequest{
		ClusterName: req.ClusterName,
		Name:        req.Name,
		Description: req.Description,
		Labels:      req.Labels,
		Annotations: req.Annotations,
	})
	if err != nil {
		return nil, err
	}
	return &pb.DeleteNameSpaceResponse{}, nil
}

func (x *NameSpaceService) CreateNameSpaceByYaml(ctx context.Context, req *pb.CreateNameSpaceByYamlRequest) (*pb.CreateNameSpaceResponse, error) {
	err := x.uc.CreateNameSpaceByYaml(ctx, &biz.CreateNameSpaceByYamlRequest{
		ClusterName: req.ClusterName,
		Yaml:        req.Yaml,
	})
	if err != nil {
		return nil, err
	}
	return &pb.CreateNameSpaceResponse{}, nil
}

func (x *NameSpaceService) UpdateNameSpaceByYaml(ctx context.Context, req *pb.CreateNameSpaceByYamlRequest) (*pb.CreateNameSpaceResponse, error) {
	err := x.uc.UpdateNameSpaceByYaml(ctx, &biz.CreateNameSpaceByYamlRequest{
		ClusterName: req.ClusterName,
		Yaml:        req.Yaml,
	})
	if err != nil {
		return nil, err
	}
	return &pb.CreateNameSpaceResponse{}, nil
}

func (x *NameSpaceService) RefreshNameSpaceYaml(ctx context.Context, req *pb.GetNameSpaceYamlRequest) (*pb.GetNameSpaceYamlResponse, error) {
	namespace, err := x.uc.GetNameSpace(ctx, &biz.GetNamespaceRequest{
		NamespaceCommonParams: biz.NamespaceCommonParams{
			Name:        req.Name,
			ClusterName: req.ClusterName,
		},
	})
	if err != nil {
		return nil, err
	}
	if namespace.APIVersion == "" {
		namespace.APIVersion = "v1"
	}
	if namespace.Kind == "" {
		namespace.Kind = "Namespace"
	}
	yamlStr, err := utils.ConvertResourceTOYaml(namespace)
	if err != nil {
		yamlStr = ""
	}
	return &pb.GetNameSpaceYamlResponse{
		Yaml: yamlStr,
	}, nil
}

func (x *NameSpaceService) GetNameSpaceDetail(ctx context.Context, req *pb.GetNameSpaceDetailRequest) (*pb.GetNameSpaceDetailResponse, error) {
	namespace, err := x.uc.GetNameSpace(ctx, &biz.GetNamespaceRequest{
		NamespaceCommonParams: biz.NamespaceCommonParams{
			Name:        req.Name,
			ClusterName: req.ClusterName,
		},
	})
	if err != nil {
		return nil, err
	}
	dto := x.convertDO2DTO(namespace)
	if namespace.APIVersion == "" {
		namespace.APIVersion = "v1"
	}
	if namespace.Kind == "" {
		namespace.Kind = "Namespace"
	}
	yamlStr, err := utils.ConvertResourceTOYaml(namespace)
	if err != nil {
		yamlStr = ""
	}
	return &pb.GetNameSpaceDetailResponse{
		Name:        namespace.Name,
		Description: dto.Description,
		CreateTime:  dto.CreateTime,
		State:       dto.State,
		Uid:         dto.Uid,
		IsDefault:   dto.IsDefault,
		Yaml:        yamlStr,
		Labels:      dto.Labels,
		Annotations: dto.Annotations,
	}, nil
}
