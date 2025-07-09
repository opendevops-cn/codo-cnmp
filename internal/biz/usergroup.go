package biz

import (
	"context"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-redis/redis/v8"
	"strconv"
)

// GrantedUserGroupItem 用户组授权列表项
type GrantedUserGroupItem struct {
	ID          uint32       `json:"id"`            // ID
	Name        string       `json:"name"`          // 用户组名称
	UserGroupID uint32       `json:"user_group_id"` // 用户组ID
	RoleDetail  []RoleDetail `json:"role_detail"`   // 权限详情
	UpdateTime  string       `json:"update_time"`   // 更新时间
}

// RoleDetail 角色详情
type RoleDetail struct {
	ClusterID uint32 `json:"cluster_id"` // 集群ID
	Namespace string `json:"namespace"`  // 命名空间
	RoleID    uint32 `json:"role_id"`    // 角色ID
}

type ListGrantedUserGroupRequest struct {
	// 模糊搜索关键字
	Keyword string `json:"keyword"`
	// 页码
	Page uint32 `json:"page"`
	// 每页数量
	PageSize uint32 `json:"page_size"`
	// 是否查询所有
	ListAll bool `json:"list_all"`
}

// ListGrantedUserGroupResponseItem 用户组授权列表项
type ListGrantedUserGroupResponseItem struct {
	ID                  uint32 `json:"id"`            // ID
	Name                string `json:"name"`          // 用户组名称
	GrantedClusterCount uint32 `json:"cluster_count"` // 已授权集群数量
	GrantedRoleCount    uint32 `json:"role_count"`    // 已授权角色数量
	UserGroupID         uint32 `json:"user_group_id"` // 用户组ID
	UpdateTime          string `json:"update_time"`   // 更新时间
}

type CreateGrantedUserGroupRequest struct {
	// 用户组名称
	Name string `json:"name"`
	// 用户组ID列表
	UserGroupIDS []uint32 `json:"user_group_ids"`
	// 权限详情
	RoleDetail []RoleDetail `json:"role_detail"`
}

type UserGroupItem struct {
	ID uint32 `json:"id"`
	// 用户组名称
	Name string `json:"name"`
	// 用户组ID列表
	UserGroupID uint32 `json:"user_group_id"`
	// 权限详情
	RoleDetail []RoleDetail `json:"role_detail"`
}

type ListUserGroupRequest struct {
	// 模糊搜索关键字
	Keyword string `json:"keyword"`
	// 页码
	Page uint32 `json:"page"`
	// 每页数量
	PageSize uint32 `json:"page_size"`
	// 是否查询所有
	ListAll bool `json:"list_all"`
}

type GetGrantedUserGroupRequest struct {
	Page        uint32 `json:"page"`
	PageSize    uint32 `json:"page_size"`
	ListAll     bool   `json:"list_all"`
	UserGroupId uint32 `json:"user_group_ids"`
}

type ListUserRequest struct {
	Keyword  string `json:"keyword"`
	Page     uint32 `json:"page"`
	PageSize uint32 `json:"page_size"`
	ListAll  bool   `json:"list_all"`
}

type IGrantedUserGroupRepo interface {
	// List 获取用户组授权列表
	List(ctx context.Context, req *ListGrantedUserGroupRequest) ([]*GrantedUserGroupItem, error)
	// Count 获取用户组授权列表数量
	Count(ctx context.Context, req *ListGrantedUserGroupRequest) (uint32, error)
	// BulkCreate 创建用户组授权
	BulkCreate(ctx context.Context, req *CreateGrantedUserGroupRequest) error
	// Update 更新用户组授权
	Update(ctx context.Context, req *CreateGrantedUserGroupRequest) error
	// Delete 删除用户组授权
	Delete(ctx context.Context, userGroupID uint32) error
	// Get 获取用户组授权详情
	Get(ctx context.Context, userGroupID uint32) (*GrantedUserGroupItem, error)
}

type IGrantedUserGroupUseCase interface {
	// ListGrantedUserGroups 获取用户组授权列表
	ListGrantedUserGroups(ctx context.Context, req *ListGrantedUserGroupRequest) ([]*ListGrantedUserGroupResponseItem, uint32, error)
	// CountGrantedUserGroups 获取用户组授权列表数量
	CountGrantedUserGroups(ctx context.Context, req *ListGrantedUserGroupRequest) (uint32, error)
	// CreateGrantedUserGroup 创建用户组授权
	CreateGrantedUserGroup(ctx context.Context, req *CreateGrantedUserGroupRequest) error
	// DeleteGrantedUserGroup 删除用户组授权
	DeleteGrantedUserGroup(ctx context.Context, userGroupID uint32) error
	// UpdateGrantedUserGroup 更新用户组授权
	UpdateGrantedUserGroup(ctx context.Context, req *CreateGrantedUserGroupRequest) error
	// GetGrantedUserGroup 获取用户组授权详情
	GetGrantedUserGroup(ctx context.Context, req *ListRoleBindingRequest) (*GrantedUserGroupItem, error)
	// ListUserGroups 获取用户组列表
	ListUserGroups(ctx context.Context, req *ListUserGroupRequest) ([]*UserGroup, uint32, error)
	// GetUsersByUserGroupID 获取用户组下的用户列表
	GetUsersByUserGroupID(ctx context.Context, userGroupID uint32) ([]*GroupUser, uint32, error)
	// GetUserGroup 查询用户组详情
	GetUserGroup(ctx context.Context, userGroupID uint32) (*UserGroup, error)
	// ListUsers 查询用户
	ListUsers(ctx context.Context, req *ListUserRequest) ([]*User, uint32, error)
}

type UserGroupUseCase struct {
	RoleBindingRepo IRoleBindingRepo
	UserGroupV2Repo IUserGroupV2Repo
	log             *log.Helper
	redis           *redis.Client
}

func (x *UserGroupUseCase) ListUsers(ctx context.Context, req *ListUserRequest) ([]*User, uint32, error) {
	return x.UserGroupV2Repo.ListUsers(ctx, req)
}

func (x *UserGroupUseCase) GetUserGroup(ctx context.Context, userGroupID uint32) (*UserGroup, error) {
	userGroups, _, err := x.UserGroupV2Repo.ListUserGroups(ctx, &ListUserGroupRequest{})
	if err != nil {
		return nil, err
	}
	for _, userGroup := range userGroups {
		if userGroup.ID == int(userGroupID) {
			return userGroup, nil
		}
	}
	return nil, fmt.Errorf("用户组不存在")
}

func (x *UserGroupUseCase) GetUsersByUserGroupID(ctx context.Context, userGroupID uint32) ([]*GroupUser, uint32, error) {
	return x.UserGroupV2Repo.GetUsersByUserGroupID(ctx, userGroupID)
}

// UserGroup 用户组
type UserGroup struct {
	ID         int    `json:"id"`
	RoleName   string `json:"role_name"`
	Details    string `json:"details"`
	Status     string `json:"status"`
	RoleType   string `json:"role_type"`
	RoleSubs   []int  `json:"role_subs"`
	CreateTime string `json:"create_time"`
	UpdateTime string `json:"update_time"`
}

// ListUserGroups 获取用户组列表
func (x *UserGroupUseCase) ListUserGroups(ctx context.Context, req *ListUserGroupRequest) ([]*UserGroup, uint32, error) {
	return x.UserGroupV2Repo.ListUserGroups(ctx, req)
}

func (x *UserGroupUseCase) GetGrantedUserGroup(ctx context.Context, req *ListRoleBindingRequest) (*GrantedUserGroupItem, error) {
	//return x.Repo.Get(ctx, userGroupID)
	bindings, err := x.RoleBindingRepo.List(ctx, req)
	if err != nil {
		return nil, err
	}
	roleDetail := make([]RoleDetail, 0)
	for _, binding := range bindings {
		roleDetail = append(roleDetail, RoleDetail{
			ClusterID: binding.ClusterID,
			Namespace: binding.Namespace,
			RoleID:    binding.RoleID,
		})
	}
	return &GrantedUserGroupItem{
		UserGroupID: req.UserGroupID,
		RoleDetail:  roleDetail,
	}, nil
}

func (x *UserGroupUseCase) UpdateGrantedUserGroup(ctx context.Context, req *CreateGrantedUserGroupRequest) error {
	//return x.Repo.Update(ctx, req)
	items := make([]*RoleBindingItem, 0)
	for _, detail := range req.RoleDetail {
		items = append(items, &RoleBindingItem{
			Namespace: detail.Namespace,
			RoleBindingCommonParams: RoleBindingCommonParams{
				UserGroupID: req.UserGroupIDS[0],
				ClusterID:   detail.ClusterID,
				RoleID:      detail.RoleID,
			},
		})
	}
	return x.RoleBindingRepo.ManageRoleBindingByUserGroupID(ctx, req.UserGroupIDS[0], items)
}

// List2Map 获取用户组列表并转换为map
func (x *UserGroupUseCase) List2Map(ctx context.Context) (map[string]string, error) {
	res := make(map[string]string)
	userGroups, _, err := x.ListUserGroups(ctx, &ListUserGroupRequest{})
	if err != nil {
		return res, err
	}
	for _, userGroup := range userGroups {
		// int 转 string
		userGroupID := strconv.Itoa(userGroup.ID)
		res[userGroupID] = userGroup.RoleName
	}
	return res, nil
}

func NewIGrantedUserGroupRepo(x *UserGroupUseCase) IGrantedUserGroupUseCase {
	return x
}

// NewUserGroupUseCase 创建用户组用例
func NewUserGroupUseCase(RoleBindingRepo IRoleBindingRepo, userGroupV2Repo IUserGroupV2Repo, logger log.Logger, redis *redis.Client) *UserGroupUseCase {
	x := &UserGroupUseCase{
		RoleBindingRepo: RoleBindingRepo,
		log:             log.NewHelper(logger),
		redis:           redis,
		UserGroupV2Repo: userGroupV2Repo,
	}
	return x
}

func (x *UserGroupUseCase) ListGrantedUserGroups(ctx context.Context, req *ListGrantedUserGroupRequest) ([]*ListGrantedUserGroupResponseItem, uint32, error) {
	list, total, err := x.RoleBindingRepo.ListByUserGroupID(ctx, &ListRoleBindingRequest{ListAll: req.ListAll, Page: req.Page, PageSize: req.PageSize})
	if err != nil {
		return nil, 0, err
	}
	respList := make([]*ListGrantedUserGroupResponseItem, 0)
	userGroupMap, err := x.List2Map(ctx)
	if err != nil {
		return nil, 0, err
	}
	for _, item := range list {
		userGroupID := strconv.Itoa(int(item.UserGroupID))
		userGroupName := userGroupMap[userGroupID]
		respList = append(respList, &ListGrantedUserGroupResponseItem{
			ID:                  item.ID,
			Name:                userGroupName,
			GrantedClusterCount: item.GrantedClusterCount,
			GrantedRoleCount:    item.GrantedRoleCount,
			UserGroupID:         item.UserGroupID,
			UpdateTime:          item.UpdateTime,
		})
	}
	return respList, total, nil
}

func (x *UserGroupUseCase) CountGrantedUserGroups(ctx context.Context, req *ListGrantedUserGroupRequest) (uint32, error) {
	//TODO implement me
	panic("implement me")
}

func (x *UserGroupUseCase) CreateGrantedUserGroup(ctx context.Context, req *CreateGrantedUserGroupRequest) error {
	//return x.Repo.BulkCreate(ctx, req)
	userGroupIDs := req.UserGroupIDS
	roleDetail := req.RoleDetail
	items := make([]*RoleBindingItem, 0)
	for _, userGroupID := range userGroupIDs {
		for _, binding := range roleDetail {
			items = append(items, &RoleBindingItem{
				Namespace: binding.Namespace,
				RoleBindingCommonParams: RoleBindingCommonParams{
					UserGroupID: userGroupID,
					ClusterID:   binding.ClusterID,
					RoleID:      binding.RoleID,
				}})
		}
	}
	return x.RoleBindingRepo.BulkCreate(ctx, items)
}

func (x *UserGroupUseCase) DeleteGrantedUserGroup(ctx context.Context, userGroupID uint32) error {
	//return x.Repo.Delete(ctx, userGroupID)
	return x.RoleBindingRepo.DeleteByUserGroupID(ctx, userGroupID)
}
