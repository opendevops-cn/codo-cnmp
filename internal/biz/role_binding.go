package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
)

type RoleBindingCommonParams struct {
	// 用户组ID
	UserGroupID uint32 `json:"user_group_id"`
	// 集群ID
	ClusterID uint32 `json:"cluster_id"`
	// 角色ID
	RoleID uint32 `json:"role_id"`
}

type RoleBindingItem struct {
	// ID
	ID uint32 `json:"id"`
	// 更新时间
	UpdateTime string `json:"update_time"`
	// namespace
	Namespace string `json:"namespace"`
	RoleBindingCommonParams
}

type ListRoleBindingRequest struct {
	// 页码
	Page uint32 `json:"page"`
	// 每页数量
	PageSize uint32 `json:"page_size"`
	// 是否查询所有
	ListAll bool `json:"list_all"`
	RoleBindingCommonParams
}

type CreateRoleBindingRequest struct {
	RoleBindingCommonParams
}

type UpdateRoleBindingRequest struct {
	RoleBindingCommonParams
}

type BulkCreateRoleBindingRequest struct {
	Items []RoleBindingItem
}

type BulkUpdateRoleBindingRequest struct {
	Items []RoleBindingItem
}

type IRoleBindingRepo interface {
	// List 获取用户组集群角色关联列表
	List(ctx context.Context, query *ListRoleBindingRequest) ([]*RoleBindingItem, error)
	// ListByUserGroupID 根据用户组ID查询角色绑定关系 聚合查询
	ListByUserGroupID(ctx context.Context, query *ListRoleBindingRequest) ([]*ListGrantedUserGroupResponseItem, uint32, error)
	// ListByRoleID 根据角色ID查询角色绑定关系 聚合查询
	ListByRoleID(ctx context.Context, query *ListRoleBindingRequest) ([]*RoleBindingItem, uint32, error)
	// Create 创建用户组集群角色关联
	Create(ctx context.Context, item *RoleBindingItem) error
	// Update 更新用户组集群角色关联
	Update(ctx context.Context, item *RoleBindingItem) error
	// DeleteByUserGroupID 根据用户组ID删除用户组集群角色关联
	DeleteByUserGroupID(ctx context.Context, id uint32) error
	// DeleteByClusterId 删除集群权限
	DeleteByClusterId(ctx context.Context, clusterID uint32) error
	// DeleteByNamespace 根据namespace删除集群权限
	DeleteByNamespace(ctx context.Context, clusterID uint32, namespace string) error
	// Count  获取用户组集群角色关联数量
	Count(ctx context.Context, query *ListRoleBindingRequest) (uint32, error)
	// BulkCreate 批量创建用户组集群角色关联
	BulkCreate(ctx context.Context, data []*RoleBindingItem) error
	// BulkUpdate 批量更新用户组集群角色关联
	BulkUpdate(ctx context.Context, data []*RoleBindingItem) error
	// ManageRoleBindingByUserGroupID 根据用户组ID管理用户组集群角色关联
	ManageRoleBindingByUserGroupID(ctx context.Context, userGroupID uint32, items []*RoleBindingItem) error
	// ManageRoleBindingByRoleID 根据角色ID管理用户组集群角色关联
	ManageRoleBindingByRoleID(ctx context.Context, roleID uint32, items []*RoleBindingItem) error
}

type IRoleBindingUseCase interface {
	// List 获取用户组集群角色关联列表
	List(ctx context.Context, req *ListRoleBindingRequest) ([]*RoleBindingItem, uint32, error)
	// Create 创建用户组集群角色关联
	Create(ctx context.Context, req *CreateRoleBindingRequest) error
	// Update 更新用户组集群角色关联
	Update(ctx context.Context, req *UpdateRoleBindingRequest) error
	// Delete 删除用户组集群角色关联
	Delete(ctx context.Context, id uint32) error
	// DeleteByClusterId 删除集群权限
	DeleteByClusterId(ctx context.Context, clusterID uint32) error
	// DeleteByNamespace 根据namespace删除集群权限
	DeleteByNamespace(ctx context.Context, clusterID uint32, namespace string) error
	// BatchCreate 批量创建用户组集群角色关联
	BatchCreate(ctx context.Context, req []*RoleBindingItem) error
	// BatchUpdate 批量更新用户组集群角色关联
	BatchUpdate(ctx context.Context, req []*RoleBindingItem) error
	// ManageRoleBindingByUserGroupID 根据用户组ID管理用户组集群角色关联
	ManageRoleBindingByUserGroupID(ctx context.Context, userGroupID uint32, items []*RoleBindingItem) error
}

type RoleBindingUseCase struct {
	repo IRoleBindingRepo
	log  *log.Helper
}

func (x *RoleBindingUseCase) DeleteByNamespace(ctx context.Context, clusterID uint32, namespace string) error {
	return x.repo.DeleteByNamespace(ctx, clusterID, namespace)
}

func (x *RoleBindingUseCase) DeleteByClusterId(ctx context.Context, clusterID uint32) error {
	return x.repo.DeleteByClusterId(ctx, clusterID)
}

func (x *RoleBindingUseCase) ManageRoleBindingByUserGroupID(ctx context.Context, userGroupID uint32, items []*RoleBindingItem) error {
	return x.repo.ManageRoleBindingByUserGroupID(ctx, userGroupID, items)
}

func (x *RoleBindingUseCase) BatchCreate(ctx context.Context, req []*RoleBindingItem) error {
	return x.repo.BulkCreate(ctx, req)
}

func (x *RoleBindingUseCase) BatchUpdate(ctx context.Context, req []*RoleBindingItem) error {
	return x.repo.BulkUpdate(ctx, req)
}

func NewRoleBindingUseCase(repo IRoleBindingRepo, logger log.Logger) *RoleBindingUseCase {
	return &RoleBindingUseCase{
		repo: repo,
		log:  log.NewHelper(logger),
	}
}

func NewIRoleBindingUseCase(x *RoleBindingUseCase) IRoleBindingUseCase {
	return x
}

func (x *RoleBindingUseCase) List(ctx context.Context, query *ListRoleBindingRequest) ([]*RoleBindingItem, uint32, error) {
	items, err := x.repo.List(ctx, query)
	if err != nil {
		return nil, 0, err
	}
	count, err := x.repo.Count(ctx, query)
	if err != nil {
		return nil, 0, err
	}
	return items, count, nil
}

func (x *RoleBindingUseCase) Create(ctx context.Context, req *CreateRoleBindingRequest) error {
	r := &RoleBindingItem{
		RoleBindingCommonParams: req.RoleBindingCommonParams,
	}
	return x.repo.Create(ctx, r)
}

func (x *RoleBindingUseCase) Update(ctx context.Context, req *UpdateRoleBindingRequest) error {
	r := &RoleBindingItem{
		RoleBindingCommonParams: req.RoleBindingCommonParams,
	}
	return x.repo.Update(ctx, r)
}

func (x *RoleBindingUseCase) Delete(ctx context.Context, id uint32) error {
	return x.repo.DeleteByUserGroupID(ctx, id)
}
