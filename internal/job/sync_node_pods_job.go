package job

import (
	"context"

	"codo-cnmp/internal/biz"
	"github.com/go-kratos/kratos/v2/log"
)

type SyncNodePodsJob struct {
	log    *log.Helper
	nodeUC biz.INodeUseCase
}

func NewSyncNodePodsJob(logger log.Logger, nodeUC biz.INodeUseCase) *SyncNodePodsJob {
	return &SyncNodePodsJob{
		log:    log.NewHelper(log.With(logger, "module", "job/node")),
		nodeUC: nodeUC,
	}
}

func (x *SyncNodePodsJob) CronSpec() string {
	return "@every 1m"
}

func (x *SyncNodePodsJob) Run(ctx context.Context) error {
	return x.nodeUC.SyncNodePods(ctx)
}
