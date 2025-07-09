package job

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"codo-cnmp/internal/biz"
	"codo-cnmp/internal/dep"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-redis/redis/v8"
)

type UserJob struct {
	apiGW *dep.CODOAPIGateway
	redis *redis.Client
	log   *log.Helper
	user  *biz.UserUseCase
}

func (x *UserJob) CronSpec() string {
	return "@every 3m"
}

const userJobLockKey = "user_job_lock"
const userJobLockExpiration = 1 * time.Minute

// acquireLock 加锁
func (x *UserJob) acquireLock(ctx context.Context) (bool, error) {
	ok, err := x.redis.SetNX(ctx, userJobLockKey, "locked", userJobLockExpiration).Result()
	if err != nil {
		x.log.WithContext(ctx).Error(fmt.Sprintf("设置任务锁失败: %v", err))
		return false, err
	}
	return ok, nil
}

// releaseLock 释放锁
func (x *UserJob) releaseLock(ctx context.Context) error {
	// Release the lock by deleting the key
	_, err := x.redis.Del(ctx, userJobLockKey).Result()
	if err != nil {
		x.log.WithContext(ctx).Error(fmt.Sprintf("释放任务锁失败: %v", err))
		return err
	}
	return nil
}

func (x *UserJob) Run(ctx context.Context) error {
	locked, err := x.acquireLock(ctx)
	if err != nil {
		x.log.WithContext(ctx).Error(fmt.Sprintf("获取任务锁失败: %v", err))
		return err
	}

	if !locked {
		x.log.WithContext(ctx).Info("任务正在执行中，跳过本次执行")
		return nil
	}

	defer func() {
		// Always release the lock after the job finishes
		if err := x.releaseLock(ctx); err != nil {
			x.log.WithContext(ctx).Error(fmt.Sprintf("释放任务锁失败: %v", err))
		}
	}()

	response, err := x.apiGW.SendRequest(ctx, "GET", "/api/p/v4/user/?page_number=1&page_size=300&order_by=id&order=ascend", nil, nil)
	if err != nil {
		x.log.WithContext(ctx).Error(fmt.Sprintf("获取用户列表失败: %v", err))
		return err
	}
	if response.StatusCode != http.StatusOK {
		x.log.WithContext(ctx).Error(fmt.Sprintf("获取用户列表失败: %v", response.Status))
		return err
	}
	defer response.Body.Close()
	data, err := io.ReadAll(response.Body)
	if err != nil {
		x.log.WithContext(ctx).Error(fmt.Sprintf("读取响应体失败: %v", err))
		return err
	}
	// 解析响应体并存入 Redis
	type Response struct {
		Code  int        `json:"code"`
		Count int        `json:"count"`
		Data  []biz.User `json:"data"`
		Msg   string     `json:"msg"`
	}
	var resp Response
	err = json.Unmarshal(data, &resp)
	if err != nil {
		x.log.WithContext(ctx).Error(fmt.Sprintf("解析响应体失败: %v", err))
		return err
	}
	if resp.Code != 0 {
		x.log.WithContext(ctx).Error(fmt.Sprintf("获取用户列表失败: %v", resp.Msg))
		return fmt.Errorf(resp.Msg)
	}
	userList := resp.Data
	//users := make([]*biz.RoleUser, 0)
	//for _, user := range userList {
	//	users = append(users, &user)
	//}
	// 写入db
	//res, err := x.user.BulkInsertOrUpdateUser(ctx, users)
	//if err != nil {
	//	x.log.WithContext(ctx).Error(fmt.Sprintf("写入用户列表失败: %v", err))
	//}
	//if !res {
	//	x.log.WithContext(ctx).Error(fmt.Sprintf("写入用户列表失败: %v", err))
	//}
	err = x.user.StoreUserSnapshot(ctx, userList)
	if err != nil {
		x.log.WithContext(ctx).Error(fmt.Sprintf("角色列表保存到 Redis 失败: %v", err))
		return err
	}
	x.log.WithContext(ctx).Info("角色列表更新成功")
	return err
}

func NewUserJob(apiGW *dep.CODOAPIGateway, logger log.Logger, redis *redis.Client, user *biz.UserUseCase) *UserJob {
	return &UserJob{
		apiGW: apiGW,
		log:   log.NewHelper(log.With(logger, "module", "job/user")),
		redis: redis,
		user:  user,
	}
}
