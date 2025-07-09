package biz

import (
	"codo-cnmp/pb"
	"context"
	"github.com/go-kratos/kratos/v2/log"
)

type UserFollowCommonParams struct {
	// 用户ID
	UserID uint32 `json:"user_id"`
	// 关注类型
	FollowType pb.FollowType
	// 关注对象
	FollowValue string
	// 集群名称
	ClusterName string
}

// UserFollowRequest 用户关注请求
type UserFollowRequest struct {
	UserFollowCommonParams
}

// DeleteUserFollowRequest 删除用户关注请求
type DeleteUserFollowRequest struct {
	UserFollowCommonParams
}

// ListUserFollowRequest 获取用户关注列表
type ListUserFollowRequest struct {
	UserFollowCommonParams
	// 页码
	Page uint32
	// 每页数量
	PageSize uint32
	// 是否查询所有
	ListAll bool
	// Keyword 关键字
	Keyword string
}

type UserFollowItem struct {
	// 主键
	ID uint32
	// 用户ID
	UserID uint32
	// 关注类型
	FollowType pb.FollowType
	// 关注对象
	FollowValue string
	// 创建时间
	CreatedTime string
	// 集群名称
	ClusterName string
}

type IUserFollowUseCase interface {
	// CreateUserFollow 创建用户关注
	CreateUserFollow(ctx context.Context, req *UserFollowRequest) error
	// DeleteUserFollow 删除用户关注
	DeleteUserFollow(ctx context.Context, req *DeleteUserFollowRequest) error
	// ListUserFollow 获取用户关注列表
	ListUserFollow(ctx context.Context, req *ListUserFollowRequest) ([]*UserFollowItem, uint32, error)
}

type IUserFollowRepo interface {
	Create(ctx context.Context, data *UserFollowRequest) error
	Delete(ctx context.Context, data *DeleteUserFollowRequest) error
	List(ctx context.Context, data *ListUserFollowRequest) ([]*UserFollowItem, uint32, error)
}

type UserFollowUseCase struct {
	repo IUserFollowRepo
	log  *log.Helper
}

func (x *UserFollowUseCase) ListUserFollow(ctx context.Context, req *ListUserFollowRequest) ([]*UserFollowItem, uint32, error) {
	return x.repo.List(ctx, req)
}

func (x *UserFollowUseCase) CreateUserFollow(ctx context.Context, data *UserFollowRequest) error {
	return x.repo.Create(ctx, data)
}

func (x *UserFollowUseCase) DeleteUserFollow(ctx context.Context, data *DeleteUserFollowRequest) error {
	return x.repo.Delete(ctx, data)
}

func NewUserFollowUseCase(repo IUserFollowRepo, logger log.Logger) *UserFollowUseCase {
	return &UserFollowUseCase{repo: repo, log: log.NewHelper(log.With(logger, "module", "biz/user_follow"))}
}

func NewIUserFollowUseCase(x *UserFollowUseCase) IUserFollowUseCase {
	return x
}
