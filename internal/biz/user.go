package biz

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"codo-cnmp/common/consts"
	"github.com/ccheers/xpkg/lru"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-redis/redis/v8"
)

type GroupUser struct {
	UserID   uint32 `json:"user_id"`
	Username string `json:"username"`
	Nickname string `json:"nickname"`
	Email    string `json:"email"`
	Source   string `json:"source"`
}

type User struct {
	UserID    uint32 `json:"id"`
	Username  string `json:"username"`
	Nickname  string `json:"nickname"`
	SuperUser string `json:"superuser"`
}

type RoleUser struct {
	UserID   uint32 `json:"id"`
	Username string `json:"username"`
	Nickname string `json:"nickname"`
	Email    string `json:"email"`
	Source   string `json:"source"`
}

type IUserRepo interface {
	CreateUser(ctx context.Context, data *RoleUser) (bool, error)
	BulkInsertOrUpdateUser(ctx context.Context, data []*RoleUser) (bool, error)
}

type IUserUseCase interface {
	CreateOrUpdateUser(ctx context.Context, data *RoleUser) (bool, error)
	BulkInsertOrUpdateUser(ctx context.Context, data []*RoleUser) (bool, error)
}

type UserUseCase struct {
	repo  IUserRepo
	log   *log.Helper
	redis *redis.Client

	cache lru.ILRUCache
}

func NewUserUseCase(repo IUserRepo, logger log.Logger, redis *redis.Client) *UserUseCase {
	return &UserUseCase{
		repo:  repo,
		log:   log.NewHelper(logger),
		redis: redis,
		cache: lru.NewLRUCache(1024),
	}
}

func (x *UserUseCase) BulkInsertOrUpdateUser(ctx context.Context, data []*RoleUser) (bool, error) {
	return x.repo.BulkInsertOrUpdateUser(ctx, data)
}

func (x *UserUseCase) CreateUser(ctx context.Context, data *RoleUser) (bool, error) {
	return x.repo.CreateUser(ctx, data)
}

func (x *UserUseCase) CreateOrUpdateUser(ctx context.Context, data *RoleUser) (bool, error) {
	return x.repo.CreateUser(ctx, data)
}

func (x *UserUseCase) StoreUserSnapshot(ctx context.Context, userList []User) error {
	// 写入redis
	jsonData, err := json.Marshal(userList)
	if err != nil {
		x.log.WithContext(ctx).Error(fmt.Sprintf("序列化角色列表失败: %v", err))
		return err
	}
	err = x.redis.Set(ctx, consts.UserCacheKey, jsonData, 0).Err()
	if err != nil {
		x.log.WithContext(ctx).Error(fmt.Sprintf("角色列表保存到 Redis 失败: %v", err))
		return err
	}
	return nil
}

func (x *UserUseCase) GetUserByUsername(ctx context.Context, username string) (*User, error) {
	return lru.FuncCacheCall(ctx, x.cache, fmt.Sprintf("get_user_by_username_%s", username), func(ctx context.Context) (*User, error) {
		// 从 Redis 中获取用户信息
		data, err := x.redis.Get(ctx, consts.UserCacheKey).Result()
		if err != nil {
			x.log.WithContext(ctx).Error(fmt.Sprintf("获取用户列表失败: %v", err))
			return nil, err
		}
		var userList []User
		err = json.Unmarshal([]byte(data), &userList)
		if err != nil {
			x.log.WithContext(ctx).Error(fmt.Sprintf("解析用户列表失败: %v", err))
			return nil, err
		}
		for _, user := range userList {
			if user.Username == username {
				return &user, nil
			}
		}
		return nil, fmt.Errorf("用户不存在")
	}, time.Minute)
}

func (x *UserUseCase) GetUserByNickname(ctx context.Context, username string) (*User, error) {
	return lru.FuncCacheCall(ctx, x.cache, fmt.Sprintf("get_user_by_nickname_%s", username), func(ctx context.Context) (*User, error) {
		// 从 Redis 中获取用户信息
		data, err := x.redis.Get(ctx, consts.UserCacheKey).Result()
		if err != nil {
			x.log.WithContext(ctx).Error(fmt.Sprintf("获取用户列表失败: %v", err))
			return nil, err
		}
		var userList []User
		err = json.Unmarshal([]byte(data), &userList)
		if err != nil {
			x.log.WithContext(ctx).Error(fmt.Sprintf("解析用户列表失败: %v", err))
			return nil, err
		}
		for _, user := range userList {
			if user.Nickname == username {
				return &user, nil
			}
		}
		return nil, fmt.Errorf("用户不存在")
	}, time.Minute)
}

func NewIUserUseCase(x *UserUseCase) IUserUseCase {
	return x
}
