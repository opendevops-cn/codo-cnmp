package job

import (
	"context"

	"github.com/google/wire"
)

type ICronJobList []ICronJob

type ICronJob interface {
	CronSpec() string
	Run(context.Context) error
}

var ProviderSet = wire.NewSet(NewJobs, NewClusterJob, NewUserGroupJob, NewMetricsJob, NewUserJob, NewSyncNodePodsJob)

func NewJobs(clusterJob *ClusterJob, userGroupJob *UserGroupJob, userJob *UserJob, syncNodePodsJob *SyncNodePodsJob) ICronJobList {
	return ICronJobList{clusterJob, userGroupJob, userJob, syncNodePodsJob}
}
