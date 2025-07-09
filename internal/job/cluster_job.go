package job

import (
	"codo-cnmp/internal/biz"
	"context"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-redis/redis/v8"
	"time"
)

type ClusterJob struct {
	cluster *biz.ClusterUseCase
	node    *biz.NodeUseCase
	log     *log.Helper
	redis   *redis.Client
}

func (x *ClusterJob) CronSpec() string {
	return "@every 3m"
}

const clusterJobLockKey = "cluster_job_lock"
const clusterJobLockExpiration = 1 * time.Minute

// acquireLock 加锁
func (x *ClusterJob) acquireLock(ctx context.Context) (bool, error) {
	ok, err := x.redis.SetNX(ctx, clusterJobLockKey, "locked", clusterJobLockExpiration).Result()
	if err != nil {
		return false, fmt.Errorf("failed to acquire lock: %v", err)
	}
	return ok, nil
}

// releaseLock 释放锁
func (x *ClusterJob) releaseLock(ctx context.Context) error {
	_, err := x.redis.Del(ctx, clusterJobLockKey).Result()
	if err != nil {
		return fmt.Errorf("failed to release lock: %v", err)
	}
	return nil
}

func (x *ClusterJob) Run(ctx context.Context) error {
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
		if err := x.releaseLock(ctx); err != nil {
			x.log.WithContext(ctx).Error(fmt.Sprintf("释放任务锁失败: %v", err))
		}
	}()
	clusters, err := x.cluster.FetchAllClusters(ctx)
	if err != nil {
		x.log.WithContext(ctx).Error(fmt.Sprintf("获取集群列表失败: %v", err))
		return err
	}

	for _, cluster := range clusters {
		x.log.WithContext(ctx).Infof("开始同步集群 %s", cluster.Name)
		err := x.cluster.HandleClusterCreative(ctx, cluster.ID)
		if err != nil {
			x.log.WithContext(ctx).Error(fmt.Sprintf("同步集群 %s 失败: %v", cluster.Name, err))
		}
		x.log.WithContext(ctx).Info(fmt.Sprintf("同步集群 %s 完成", cluster.Name))
	}
	x.log.WithContext(ctx).Info("同步集群完成")

	return nil
}

func NewClusterJob(cluster *biz.ClusterUseCase, node *biz.NodeUseCase, logger log.Logger, redis *redis.Client) *ClusterJob {
	return &ClusterJob{
		cluster: cluster,
		node:    node,
		log:     log.NewHelper(log.With(logger, "module", "job/cluster")),
		redis:   redis,
	}
}
