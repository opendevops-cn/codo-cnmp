package dep

import (
	"context"
	"github.com/ccheers/xpkg/sync/try_lock"
	"github.com/ccheers/xpkg/xmsgbus"
	r "github.com/ccheers/xpkg/xmsgbus/impl/redis"
	"github.com/go-redis/redis/v8"
	"github.com/opendevops-cn/codo-golang-sdk/tools/cascmd"
)

func NewMsgBus(client *redis.Client) xmsgbus.IMsgBus {
	return r.NewMsgBus(client)
}

func NewTopicManager(ctx context.Context, bus xmsgbus.IMsgBus, cas try_lock.CASCommand, storage xmsgbus.ISharedStorage) xmsgbus.ITopicManager {
	return xmsgbus.NewTopicManager(ctx, bus, cas, storage)
}

func NewSharedStorage(client *redis.Client) xmsgbus.ISharedStorage {
	return r.NewSharedStorage(client)
}

func NewCAS(client *redis.Client) try_lock.CASCommand {
	return cascmd.NewCasCmd(client)
}
