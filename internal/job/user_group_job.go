package job

import (
	"codo-cnmp/common/consts"
	"codo-cnmp/internal/biz"
	"codo-cnmp/internal/dep"
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-redis/redis/v8"
	"io"
	"net/http"
	"strconv"
	"time"
)

type UserGroupJob struct {
	apiGW *dep.CODOAPIGateway
	redis *redis.Client
	log   *log.Helper
}

type UserGroupUser struct {
	UserID   uint32 `json:"user_id"`
	Username string `json:"username"`
	Nickname string `json:"nickname"`
	Email    string `json:"email"`
	RoleID   uint32 `json:"role_id"`
}

func (x *UserGroupJob) CronSpec() string {
	return "@every 2m"
}

const userGroupJobLockKey = "user_group_job_lock"

const userGroupJobLockExpiration = 10 * time.Minute

// acquireLock 加锁
func (x *UserGroupJob) acquireLock(ctx context.Context) (bool, error) {
	ok, err := x.redis.SetNX(ctx, userGroupJobLockKey, "locked", userGroupJobLockExpiration).Result()
	if err != nil {
		return false, fmt.Errorf("failed to acquire lock: %v", err)
	}
	return ok, nil
}

// releaseLock 释放锁
func (x *UserGroupJob) releaseLock(ctx context.Context) error {
	// Release the lock by deleting the key
	_, err := x.redis.Del(ctx, userGroupJobLockKey).Result()
	if err != nil {
		return fmt.Errorf("failed to release lock: %v", err)
	}
	return nil
}

func (x *UserGroupJob) GetUsersByUserGroupID(ctx context.Context, userGroupID uint32) ([]*UserGroupUser, error) {
	users := make([]*UserGroupUser, 0)
	response, err := x.apiGW.SendRequest(ctx, "GET", fmt.Sprintf("/api/p/v4/role_user/?role_id=%d", userGroupID), nil, nil)
	if err != nil {
		x.log.WithContext(ctx).Error(fmt.Sprintf("获取用户列表失败: %v", err))
		return users, err
	}
	if response.StatusCode != http.StatusOK {
		x.log.WithContext(ctx).Error(fmt.Sprintf("获取用户列表失败, 状态码异常: %v", response.Status))
		return users, err
	}
	defer response.Body.Close()
	data, err := io.ReadAll(response.Body)
	if err != nil {
		x.log.WithContext(ctx).Error(fmt.Sprintf("读取响应体失败: %v", err))
		return users, err
	}
	// 解析响应体并存入 Redis
	type UserResponse struct {
		Code  int             `json:"code"`
		Count int             `json:"count"`
		Data  []UserGroupUser `json:"data"`
		Msg   string          `json:"msg"`
	}
	var userResp UserResponse
	err = json.Unmarshal(data, &userResp)
	if err != nil {
		x.log.WithContext(ctx).Error(fmt.Sprintf("解析响应体失败: %v", err))
		return users, err
	}
	if userResp.Code != 0 {
		x.log.WithContext(ctx).Error(fmt.Sprintf("获取用户列表失败: %v", userResp.Msg))
		return users, fmt.Errorf(userResp.Msg)
	}
	for _, user := range userResp.Data {
		users = append(users, &user)
	}
	return users, nil
}

func (x *UserGroupJob) Run(ctx context.Context) error {
	locked, err := x.acquireLock(ctx)
	if err != nil {
		x.log.WithContext(ctx).Error(fmt.Sprintf("获取任务锁失败: %v", err))
		return err
	}

	if !locked {
		x.log.WithContext(ctx).Infof("任务正在执行中，跳过本次执行: %v", time.Now().Format(time.DateTime))
		return nil
	}

	defer func() {
		if err := x.releaseLock(ctx); err != nil {
			x.log.WithContext(ctx).Error(fmt.Sprintf("释放任务锁失败: %v", err))
		}
	}()
	response, err := x.apiGW.SendRequest(ctx, "GET", "/api/p/v4/role/?page_number=1&page_size=300&role_type=normal&order_by=role_name&order=descend", nil, nil)
	if err != nil {
		x.log.WithContext(ctx).Error(fmt.Sprintf("获取角色列表失败: %v", err))
		return err
	}
	if response.StatusCode != http.StatusOK {
		x.log.WithContext(ctx).Error(fmt.Sprintf("获取角色列表失败: %v", response.Status))
		return err
	}
	defer response.Body.Close()
	data, err := io.ReadAll(response.Body)
	if err != nil {
		x.log.WithContext(ctx).Error(fmt.Sprintf("读取响应体失败: %v", err))
		return err
	}
	// 解析响应体并存入 Redis
	type UserGroupResponse struct {
		Code  int             `json:"code"`
		Count int             `json:"count"`
		Data  []biz.UserGroup `json:"data"`
		Msg   string          `json:"msg"`
	}
	var userGroupResp UserGroupResponse
	err = json.Unmarshal(data, &userGroupResp)
	if err != nil {
		x.log.WithContext(ctx).Error(fmt.Sprintf("解析响应体失败: %v", err))
		return err
	}
	if userGroupResp.Code != 0 {
		x.log.WithContext(ctx).Error(fmt.Sprintf("获取角色列表失败: %v", userGroupResp.Msg))
		return fmt.Errorf(userGroupResp.Msg)
	}
	userGroups := userGroupResp.Data

	// 创建一个 map 来存储用户 ID 到用户组ID 的映射
	userUserGroupMap := make(map[string][]int)

	// 遍历 roleList，构建用户 ID 到角色 ID 的映射
	for _, userGroup := range userGroups {
		users, err := x.GetUsersByUserGroupID(ctx, uint32(userGroup.ID))
		if err != nil {
			x.log.WithContext(ctx).Error(fmt.Sprintf("获取用户列表失败: %v", err))
			continue
		}
		for _, user := range users {
			uid := strconv.Itoa(int(user.UserID))
			if _, exists := userUserGroupMap[uid]; !exists {
				userUserGroupMap[uid] = []int{}
			}
			userUserGroupMap[uid] = append(userUserGroupMap[uid], userGroup.ID)
		}
	}

	// 将映射存入 Redis
	for userID, userGroupIDs := range userUserGroupMap {
		// 将角色 ID 列表转换为 JSON
		userGroupIDsJSON, err := json.Marshal(userGroupIDs)
		if err != nil {
			x.log.WithContext(ctx).Error(fmt.Sprintf("用户ID %s 的用户组ID 列表为 JSON 失败: %v", userID, err))
			continue
		}

		// 使用用户 ID 作为 key，角色 ID 列表作为 value 存入 Redis
		key := fmt.Sprintf(consts.UserGroupIDCacheKey, userID)
		err = x.redis.Set(ctx, key, userGroupIDsJSON, 0).Err()
		if err != nil {
			x.log.WithContext(ctx).Error(fmt.Sprintf("用户ID [%s] 的用户组ID列表到 Redis 失败: %v", userID, err))
			continue
		}
	}

	jsonData, err := json.Marshal(userGroups)
	if err != nil {
		x.log.WithContext(ctx).Error(fmt.Sprintf("序列化角色列表失败: %v", err))
		return err
	}
	err = x.redis.Set(ctx, consts.UserGroupCacheKey, jsonData, 0).Err()
	if err != nil {
		x.log.WithContext(ctx).Error(fmt.Sprintf("角色列表保存到 Redis 失败: %v", err))
		return err
	}
	x.log.WithContext(ctx).Info("角色列表更新成功")
	return nil
}

func NewUserGroupJob(apiGW *dep.CODOAPIGateway, logger log.Logger, redis *redis.Client) *UserGroupJob {
	return &UserGroupJob{
		apiGW: apiGW,
		log:   log.NewHelper(log.With(logger, "module", "job/user_group")),
		redis: redis,
	}
}
