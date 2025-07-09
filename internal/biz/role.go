package biz

import (
	"codo-cnmp/common/utils"
	"codo-cnmp/pb"
	"context"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-redis/redis/v8"
)

// RoleItem 角色列表项
type RoleItem struct {
	ID          uint32      `json:"id"`            // ID
	Name        string      `json:"name"`          // 角色名称
	ISDefault   bool        `json:"is_default"`    // 是否默认角色
	RoleType    pb.RoleType `json:"role_type"`     // 角色类型
	Description string      `json:"user_group_id"` // 描述
	CreateTime  string      `json:"create_time"`   // 创建时间
	UpdateTime  string      `json:"update_time"`   // 更新时间
	YamlStr     string      `json:"yaml_string"`   // yaml字符串
	UpdateBy    string      `json:"update_by"`     // 更新人
}

type RoleCommonParams struct {
}

type ListRoleRequest struct {
	// 模糊搜索关键字
	Keyword string `json:"keyword"`
	// 页码
	Page uint32 `json:"page"`
	// 每页数量
	PageSize uint32 `json:"page_size"`
	// 是否查询所有
	ListAll bool `json:"list_all"`
}

// ListRoleResponseItem 角色列表项
type ListRoleResponseItem struct {
	ID                  uint32 `json:"id"`                    // ID
	Name                string `json:"name"`                  // 用户组名称
	GrantedClusterCount uint32 `json:"granted_cluster_count"` // 已授权集群数量
	GrantedRoleCount    uint32 `json:"granted_role_count"`    // 已授权角色数量
	UserGroupID         string `json:"user_group_id"`         // 用户组ID
	UpdateTime          string `json:"update_time"`           // 更新时间
}

// CreateRoleRequest 创建角色请求
type CreateRoleRequest struct {
	RoleItem
}

type UpdateRoleRequest struct {
	RoleItem
}

type RoleBindingRequest struct {
	// 页码
	Page uint32 `json:"page"`
	// 每页数量
	PageSize uint32 `json:"page_size"`
	// 是否查询所有
	ListAll bool   `json:"list_all"`
	RoleID  uint32 `json:"role_id"`
}

type IRoleRepo interface {
	// List 获取角色列表
	List(ctx context.Context, req *ListRoleRequest) ([]*RoleItem, error)
	// Count 获取用户组列表数量
	Count(ctx context.Context, req *ListRoleRequest) (uint32, error)
	// Create 创建角色
	Create(ctx context.Context, req *RoleItem) error
	// Update 更新角色
	Update(ctx context.Context, req *RoleItem) error
	// Delete 删除角色
	Delete(ctx context.Context, userGroupID uint32) error
	// GetRoleByID  获取角色
	GetRoleByID(ctx context.Context, roleID uint32) (*RoleItem, error)
	// GetRoleByName 获取角色
	GetRoleByName(ctx context.Context, roleName string) (*RoleItem, error)
	// ExistRoleByName  检查角色是否存在
	ExistRoleByName(ctx context.Context, roleName string) (bool, error)
}

type IRoleUseCase interface {
	// ListRoles  获取角色列表
	ListRoles(ctx context.Context, req *ListRoleRequest) ([]*RoleItem, uint32, error)
	// CreateRole 创建角色
	CreateRole(ctx context.Context, req *RoleItem) error
	// DeleteRole 删除角色
	DeleteRole(ctx context.Context, roleID uint32) error
	// UpdateRole 更新角色
	UpdateRole(ctx context.Context, req *RoleItem) error
	// GetRole 获取角色
	GetRole(ctx context.Context, roleID uint32) (*RoleItem, error)
	// ListRoleBindings 获取角色绑定列表
	ListRoleBindings(ctx context.Context, req *ListRoleBindingRequest) ([]*RoleBindingItem, uint32, error)
	// UpdateRoleBinding 更新角色绑定
	UpdateRoleBinding(ctx context.Context, data []*RoleBindingItem) error
	// ExistByRoleName  检查角色是否存在
	ExistByRoleName(ctx context.Context, roleName string) (bool, error)
}

type RoleUseCase struct {
	Repo         IRoleRepo
	RoleBindRepo IRoleBindingRepo
	log          *log.Helper
	redis        *redis.Client
}

func (x *RoleUseCase) ExistByRoleName(ctx context.Context, roleName string) (bool, error) {
	return x.Repo.ExistRoleByName(ctx, roleName)
}

func (x *RoleUseCase) ListRoleBindings(ctx context.Context, req *ListRoleBindingRequest) ([]*RoleBindingItem, uint32, error) {
	results, count, err := x.RoleBindRepo.ListByRoleID(ctx, &ListRoleBindingRequest{
		Page:     req.Page,
		PageSize: req.PageSize,
		ListAll:  req.ListAll,
		RoleBindingCommonParams: RoleBindingCommonParams{
			RoleID: req.RoleID,
		},
	})
	if err != nil {
		return nil, 0, err
	}
	return results, count, nil
}

func (x *RoleUseCase) UpdateRoleBinding(ctx context.Context, data []*RoleBindingItem) error {
	return x.RoleBindRepo.ManageRoleBindingByRoleID(ctx, data[0].RoleID, data)
}

func NewRoleUseCase(repo IRoleRepo, roleBindRepo IRoleBindingRepo, logger log.Logger, redis *redis.Client) *RoleUseCase {
	return &RoleUseCase{Repo: repo, log: log.NewHelper(logger), redis: redis, RoleBindRepo: roleBindRepo}
}

func NewIRoleUseCase(x *RoleUseCase) IRoleUseCase {
	return x
}

func (x *RoleUseCase) ListRoles(ctx context.Context, req *ListRoleRequest) ([]*RoleItem, uint32, error) {
	roles, err := x.Repo.List(ctx, req)
	if err != nil {
		return nil, 0, err
	}
	count, err := x.Repo.Count(ctx, req)
	if err != nil {
		return nil, 0, err
	}
	return roles, count, nil
}

func (x *RoleUseCase) CreateRole(ctx context.Context, req *RoleItem) error {
	_, err := utils.ParseK8sClusterRoleYAML(req.YamlStr)
	if err != nil {
		x.log.Errorf("序列化对象为ClusterRole YAML 失败: %v", err)
		return fmt.Errorf("yaml格式错误")
	}
	if req.UpdateBy == "" {
		userName, err := utils.GetUserNameFromCtx(ctx)
		if err != nil {
			return fmt.Errorf("创建角色失败，未获取到操作人: %w", err)
		}
		req.UpdateBy = userName
		return x.Repo.Create(ctx, req)

	}
	return x.Repo.Create(ctx, req)

}

func (x *RoleUseCase) DeleteRole(ctx context.Context, roleID uint32) error {
	return x.Repo.Delete(ctx, roleID)
}

func (x *RoleUseCase) UpdateRole(ctx context.Context, req *RoleItem) error {
	_, err := utils.ParseK8sClusterRoleYAML(req.YamlStr)
	if err != nil {
		x.log.Errorf("序列化对象为ClusterRole YAML 失败: %v", err)
		return fmt.Errorf("yaml格式错误: %v", err)
	}
	userName, err := utils.GetUserNameFromCtx(ctx)
	if err != nil {
		return fmt.Errorf("更新角色失败，未获取到操作人")
	}
	req.UpdateBy = userName
	return x.Repo.Update(ctx, req)
}

func (x *RoleUseCase) GetRole(ctx context.Context, RoleID uint32) (*RoleItem, error) {
	return x.Repo.GetRoleByID(ctx, RoleID)
}
