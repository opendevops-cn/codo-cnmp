package server

import (
	"codo-cnmp/internal/job"
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/robfig/cron/v3"
)

type CronServer struct {
	cron *cron.Cron
}

func NewCronServer(list job.ICronJobList, logger log.Logger) (*CronServer, error) {
	l := log.NewHelper(log.With(logger, "module", "server/cron"))
	_cron := cron.New()
	for _, j := range list {
		_, err := _cron.AddFunc(j.CronSpec(), func() {
			err := j.Run(context.TODO())
			if err != nil {
				l.WithContext(context.Background()).Errorf("任务执行失败: %v", err)
			} else {
				l.WithContext(context.Background()).Infof("任务执行成功")
			}
		})
		if err != nil {
			return nil, err
		}
	}
	return &CronServer{cron: _cron}, nil
}

func (x *CronServer) Start(ctx context.Context) error {
	x.cron.Start()
	return nil
}

func (x *CronServer) Stop(ctx context.Context) error {
	select {
	case <-x.cron.Stop().Done():
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
	return nil
}
