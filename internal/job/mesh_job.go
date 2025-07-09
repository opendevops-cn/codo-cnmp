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
	"time"
)

type MeshJob struct {
	apiGW *dep.CODOAPIGateway
	redis *redis.Client
	log   *log.Helper
}

func (x *MeshJob) CronSpec() string {
	return "@every 1m"
}

const meshJobLockKey = "mesh_job_lock"
const meshJobLockExpiration = 1 * time.Minute

// acquireLock 加锁
func (x *MeshJob) acquireLock(ctx context.Context) (bool, error) {
	ok, err := x.redis.SetNX(ctx, meshJobLockKey, "locked", meshJobLockExpiration).Result()
	if err != nil {
		x.log.WithContext(ctx).Error(fmt.Sprintf("设置任务锁失败: %v", err))
		return false, err
	}
	return ok, nil
}

type MeshResponse struct {
	Status int `json:"status"`
	Code   int `json:"code"`
	Data   struct {
		Total int        `json:"total"`
		List  []MeshItem `json:"list"`
	} `json:"data"`
}

// releaseLock 释放锁
func (x *MeshJob) releaseLock(ctx context.Context) error {
	// Release the lock by deleting the key
	_, err := x.redis.Del(ctx, meshJobLockKey).Result()
	if err != nil {
		x.log.WithContext(ctx).Error(fmt.Sprintf("释放任务锁失败: %v", err))
		return err
	}
	return nil
}

type MeshItem struct {
	Id             int       `json:"id"`
	ServiceName    string    `json:"service_name"`
	WhiteIpList    []string  `json:"white_ip_list"`
	SrcAgentId     string    `json:"src_agent_id"`
	SrcAgentPort   int       `json:"src_agent_port"`
	DstAgentId     string    `json:"dst_agent_id"`
	DstServiceAddr string    `json:"dst_service_addr"`
	HeartbeatAt    time.Time `json:"heartbeat_at"`
	CreatedAt      time.Time `json:"created_at"`
}

func (x *MeshJob) Run(ctx context.Context) error {
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

	response, err := x.apiGW.SendRequest(ctx, "GET", "/api/agent/v1/manager/agent/mesh/list?page=1&page_size=200", nil, nil)
	if err != nil {
		x.log.WithContext(ctx).Error(fmt.Sprintf("获取组网列表失败: %v", err))
		return err
	}
	if response.StatusCode != http.StatusOK {
		x.log.WithContext(ctx).Error(fmt.Sprintf("获取组网列表失败: %v", response.Status))
		return err
	}
	defer response.Body.Close()
	data, err := io.ReadAll(response.Body)
	if err != nil {
		x.log.WithContext(ctx).Error(fmt.Sprintf("读取响应体失败: %v", err))
		return err
	}
	var resp MeshResponse
	err = json.Unmarshal(data, &resp)
	if err != nil {
		x.log.WithContext(ctx).Error(fmt.Sprintf("解析响应体失败: %v", err))
		return err
	}
	if resp.Code != 200 {
		x.log.WithContext(ctx).Error(fmt.Sprintf("获取组网列表失败: %v", resp.Status))
		return fmt.Errorf("获取组网列表失败: %v", resp.Status)
	}
	list := resp.Data.List
	// 转为json存入 Redis
	res, err := json.Marshal(list)
	if err != nil {
		x.log.WithContext(ctx).Error(fmt.Sprintf("序列化组网列表失败: %v", err))
		return err
	}
	x.redis.Set(ctx, consts.MeshCacheKey, res, 0)
	x.log.WithContext(ctx).Info("组网列表已存入 Redis")
	return err
}

func NewMeshJob(apiGW *dep.CODOAPIGateway, logger log.Logger, redis *redis.Client, user *biz.UserUseCase) *MeshJob {
	return &MeshJob{
		apiGW: apiGW,
		log:   log.NewHelper(log.With(logger, "module", "job/user")),
		redis: redis,
	}
}
