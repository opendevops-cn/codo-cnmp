package data

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/gogf/gf/v2/database/gdb"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewData, NewGreeterRepo, NewClusterRepo, NewIClusterRepo, NewNodeRepo, NewINodeRepo,
	NewRoleRepo, NewIGrantedRoleRepo, NewRoleBindingRepo, NewIRoleBindingRepoRepo, NewUserRepo, NewIUserRepo,
	NewUserGroupRepoRepoV2, NewIUserGroupRepoV2, NewUserFollowRepo, NewIUserFollowRepo, NewAuditLogRepo,
	NewIAuditLogRepo, NewGameServerRepo, NewIGameServerRepo, NewIMeshRepo, NewMeshRepo, NewAgentRepo, NewIAgentRepo)

// Data .
type Data struct {
	db    gdb.DB
	redis *redis.Client
}

// NewData .
func NewData(logger log.Logger, db gdb.DB, redisClient *redis.Client) (*Data, func(), error) {
	cleanup := func() {
		helper := log.NewHelper(logger)
		helper.Info("closing the data resources")

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()
		err := db.Close(ctx)
		if err != nil {
			helper.Errorf("failed to close the database: %v", err)
		}

		err = redisClient.Close()
		if err != nil {
			helper.Errorf("failed to close the redis: %v", err)
		}
	}

	// 初始化数据库连接
	return &Data{
		db:    db,
		redis: redisClient,
	}, cleanup, nil
}
