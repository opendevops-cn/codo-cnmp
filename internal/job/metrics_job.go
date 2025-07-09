package job

import (
	"codo-cnmp/internal/biz"
	"codo-cnmp/internal/conf"
	"codo-cnmp/pb"
	"context"
	"fmt"
	"github.com/ccheers/xpkg/sync/errgroup"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-redis/redis/v8"
	"github.com/opendevops-cn/codo-golang-sdk/client/xvm"
	"sync"
	"time"
)

type MetricType string

const (
	MetricServerName       MetricType = "server_name"
	MetricEntityCount      MetricType = "entity_count"
	MetricLockEntityStatus MetricType = "lock_entity_status"
	MetricLockLbStatus     MetricType = "lock_lb_status"
	MetricOnlineNumber     MetricType = "online_number"
	ServerNameLabel                   = "server_name"
	ClusterLabel                      = "_k8s_"
	NamespaceLabel                    = "_namespace_"
	PodLabel                          = "_pod_name_"
	EntityCountLabel                  = "entity_count"
	OnlineNumberLabel                 = "online_number"
	LockEntityStatusLabel             = "lock_entity_status"
	LockLbStatusLabel                 = "lock_lb_status"
)

type MetricsJob struct {
	log    *log.Helper
	gs     *biz.GameServerUseCase
	client xvm.IMetricsClient
	redis  *redis.Client
}

type MetricInfo struct {
	ServerName       string
	EntityCount      uint32
	OnlineNumber     uint32
	LockEntityStatus uint32
	LockLbStatus     uint32
	ClusterName      string
	Namespace        string
	PodName          string
}

func NewMetricsJob(gs *biz.GameServerUseCase, logger log.Logger, bc *conf.Bootstrap, redis *redis.Client) *MetricsJob {
	// 创建 metrics client
	grafanaConfig := bc.GRAFANA
	client, err := xvm.NewMetricsClient(
		grafanaConfig.GetADDR(),
		xvm.WithClientOptionBasicAuth(grafanaConfig.GetUSER(), grafanaConfig.GetPASSWORD()))
	if err != nil {
		fmt.Println("创建 metrics client 失败: ", err)
		return nil
	}

	return &MetricsJob{
		log:    log.NewHelper(log.With(logger, "module", "job")),
		gs:     gs,
		client: client,
		redis:  redis,
	}
}

const metricsJobLockKey = "metrics_job_lock"
const metricsJobLockExpiration = 10 * time.Second

// acquireLock 加锁
func (x *MetricsJob) acquireLock(ctx context.Context) (bool, error) {
	ok, err := x.redis.SetNX(ctx, metricsJobLockKey, "locked", metricsJobLockExpiration).Result()
	if err != nil {
		x.log.WithContext(ctx).Error(fmt.Sprintf("设置任务锁失败: %v", err))
		return false, err
	}
	return ok, nil
}

// releaseLock 释放锁
func (x *MetricsJob) releaseLock(ctx context.Context) error {
	// Release the lock by deleting the key
	_, err := x.redis.Del(ctx, metricsJobLockKey).Result()
	if err != nil {
		x.log.WithContext(ctx).Error(fmt.Sprintf("释放任务锁失败: %v", err))
		return err
	}
	return nil
}

func (x *MetricsJob) Run(ctx context.Context) error {
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
	metrics := []MetricType{
		MetricServerName,
		MetricEntityCount,
		MetricLockEntityStatus,
		MetricLockLbStatus,
		MetricOnlineNumber,
	}

	// 使用 errgroup 并发查询指标
	eg := errgroup.WithContext(ctx)
	var mu sync.RWMutex
	metricMap := make(map[string]*MetricInfo)

	// 并发查询每个指标
	for _, metric := range metrics {
		metric := metric // 避免闭包问题
		eg.Go(func(ctx context.Context) error {
			return x.queryAndProcessMetric(ctx, metric, &mu, metricMap)
		})
	}

	// 等待所有查询完成
	if err := eg.Wait(); err != nil {
		return fmt.Errorf("查询指标等待线程结束失败: %w", err)
	}

	// 转换为列表并记录日志
	return x.processResults(ctx, metricMap)
}

func (x *MetricsJob) queryAndProcessMetric(ctx context.Context, metricType MetricType, mu *sync.RWMutex, metricMap map[string]*MetricInfo) error {
	// 设置查询超时
	queryCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	result, err := x.client.Query(queryCtx, string(metricType), xvm.WithQueryOptionLimit(1))
	if err != nil {
		return fmt.Errorf("查询指标 %s 失败: %w", metricType, err)
	}
	// 处理不同类型的结果
	switch {
	case result.Data.IsVector():
		// 处理向量类型
		for _, metric := range result.Data.VectorValue() {
			if err := x.processVectorMetric(metricType, metric, mu, metricMap); err != nil {
				x.log.WithContext(ctx).Errorf("处理向量指标失败: %v", err)
			}
		}

	case result.Data.IsMatrix():
		// 处理矩阵类型
		for _, metric := range result.Data.MatrixValue() {
			if err := x.processMatrixMetric(metricType, metric, mu, metricMap); err != nil {
				x.log.WithContext(ctx).Errorf("处理矩阵指标失败: %v", err)
			}
		}

	case result.Data.IsScalar():
		// 处理标量类型
		for _, value := range result.Data.ScalarValue() {
			x.log.WithContext(ctx).Infof("指标 %s 的标量值: %v", metricType, value.Value)
		}

	case result.Data.IsString():
		// 处理字符串类型
		for _, value := range result.Data.StringValue() {
			x.log.WithContext(ctx).Infof("指标 %s 的字符串值: %v", metricType, value.Value)
		}

	default:
		x.log.WithContext(ctx).Warnf("未知的指标类型: %s", result.Data.ResultType)
	}
	return nil
}

func transferLBStatus(value float64) uint32 {
	switch value {
	case 0:
		return uint32(pb.LBLockStatus_LB_UNLOCK)
	case 1:
		return uint32(pb.LBLockStatus_LB_LOCKED)
	case 2:
		return uint32(pb.LBLockStatus_LB_HPA_LOCKED)
	default:
		return uint32(pb.LBLockStatus_UNKNOWN_ServerLockStatus)
	}
}

func transferEntityStatus(value float64) uint32 {
	switch value {
	case 0:
		return uint32(pb.EntityLockStatus_ENTITY_UNLOCK)
	case 1:
		return uint32(pb.EntityLockStatus_ENTITY_LOCKED)
	default:
		return uint32(pb.EntityLockStatus_UNKNOWN_EntityLockStatus)
	}
}

// 处理向量类型的指标
func (x *MetricsJob) processVectorMetric(metricType MetricType, metric xvm.VectorMetric, mu *sync.RWMutex, metricMap map[string]*MetricInfo) error {
	labels := metric.Labels
	serverName := labels[ServerNameLabel]
	if serverName == "" {
		return nil
	}
	mu.Lock()
	defer mu.Unlock()
	// 获取或创建 MetricInfo
	info, exists := metricMap[serverName]
	if !exists {
		info = &MetricInfo{
			ServerName:  serverName,
			ClusterName: labels[ClusterLabel],
			Namespace:   labels[NamespaceLabel],
			PodName:     labels[PodLabel],
		}
		metricMap[serverName] = info
	}

	// 获取指标值
	value := metric.Value.Value()

	// 根据指标类型更新值
	switch metricType {
	case MetricEntityCount:
		// sdk 返回的指标值可能为负数，这里将其置为 0
		if value < 0 {
			value = 0
		}
		info.EntityCount = uint32(value)
	case MetricOnlineNumber:
		if value < 0 {
			value = 0
		}
		info.OnlineNumber = uint32(value)
	case MetricLockEntityStatus:
		v := transferEntityStatus(value)
		info.LockEntityStatus = v
	case MetricLockLbStatus:
		v := transferLBStatus(value)
		info.LockLbStatus = v
	default:
		return fmt.Errorf("未知的指标类型: %s", metricType)
	}

	return nil
}

func (x *MetricsJob) processMatrixMetric(metricType MetricType, metric xvm.MatrixMetric, mu *sync.RWMutex, metricMap map[string]*MetricInfo) error {
	labels := metric.Labels
	serverName := labels[ServerNameLabel]
	if serverName == "" {
		return nil
	}
	mu.Lock()
	defer mu.Unlock()
	// 获取或创建 MetricInfo
	info, exists := metricMap[serverName]
	if !exists {
		info = &MetricInfo{
			ServerName:  serverName,
			ClusterName: labels[ClusterLabel],
			Namespace:   labels[NamespaceLabel],
			PodName:     labels[PodLabel],
		}
		metricMap[serverName] = info
	}

	// 使用最新的值
	if len(metric.Values) > 0 {
		latestValue := metric.Values[len(metric.Values)-1].Value()
		switch metricType {
		case MetricEntityCount:
			info.EntityCount = uint32(latestValue)
		case MetricOnlineNumber:
			info.OnlineNumber = uint32(latestValue)
		case MetricLockEntityStatus:
			v := transferEntityStatus(latestValue)
			info.LockEntityStatus = v
		case MetricLockLbStatus:
			v := transferLBStatus(latestValue)
			info.LockLbStatus = v
		default:
			return fmt.Errorf("未知的指标类型: %s", metricType)
		}
	}

	return nil
}

func (x *MetricsJob) processResults(ctx context.Context, metricMap map[string]*MetricInfo) error {
	gameServers := make([]*biz.GameServer, 0, len(metricMap))
	for _, info := range metricMap {
		// 保存到数据库
		gss := &biz.GameServer{
			ServerName:       info.ServerName,
			EntityNum:        info.EntityCount,
			OnlineNum:        info.OnlineNumber,
			EntityLockStatus: info.LockEntityStatus,
			LbLockStatus:     info.LockLbStatus,
			ClusterName:      info.ClusterName,
			Namespace:        info.Namespace,
			Pod:              info.PodName,
		}
		gameServers = append(gameServers, gss)
		//if err := x.gs.UpdateGameServer(ctx, data); err != nil {
		//	x.log.WithContext(ctx).Errorf("更新游戏进程指标失败: %v", err)
		//	continue
		//}
	}
	err := x.gs.BatchUpdateGameServer(ctx, gameServers)
	if err != nil {
		x.log.WithContext(ctx).Errorf("批量更新游戏进程指标失败: %v", err)
		return err
	}
	return nil
}

func (x *MetricsJob) CronSpec() string {
	return "@every 5s"
}
