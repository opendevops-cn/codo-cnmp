package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-redis/redis/v8"
)

type UserGroupV2 struct {
	UserGroupID uint64 `json:"id"`
	Name        string `json:"name"`
}

type IUserGroupV2Repo interface {
	CreateUserGroup(ctx context.Context, data *UserGroupV2) (bool, error)
	BulkInsertOrUpdateUserGroup(ctx context.Context, data []*UserGroupV2) (bool, error)
	ListUserGroups(ctx context.Context, req *ListUserGroupRequest) ([]*UserGroup, uint32, error)
	GetUsersByUserGroupID(ctx context.Context, groupID uint32) ([]*GroupUser, uint32, error)
	ListUsers(ctx context.Context, req *ListUserRequest) ([]*User, uint32, error)
}

type IUserGroupV2UseCase interface {
	CreateOrUpdateUserGroup(ctx context.Context, data *UserGroupV2) (bool, error)
	BulkInsertOrUpdateUserGroup(ctx context.Context, data []*UserGroupV2) (bool, error)
}

type UserGroupV2UseCase struct {
	repo  IUserGroupV2Repo
	log   *log.Helper
	redis *redis.Client
}

func (x *UserGroupV2UseCase) BulkInsertOrUpdateUserGroup(ctx context.Context, data []*UserGroupV2) (bool, error) {
	return x.repo.BulkInsertOrUpdateUserGroup(ctx, data)
}

func (x *UserGroupV2UseCase) CreateOrUpdateUserGroup(ctx context.Context, data *UserGroupV2) (bool, error) {
	//TODO implement me
	panic("implement me")
}

func (x *UserGroupV2UseCase) CreateUserGroup(ctx context.Context, data *UserGroupV2) (bool, error) {
	return x.repo.CreateUserGroup(ctx, data)
}

func NewUserGroupV2UseCase(repo IUserGroupV2Repo, logger log.Logger, redis *redis.Client) *UserGroupV2UseCase {
	return &UserGroupV2UseCase{repo: repo, log: log.NewHelper(logger), redis: redis}
}

func NewIUserGroupV2UseCase(x *UserGroupV2UseCase) IUserGroupV2UseCase {
	return x
}
