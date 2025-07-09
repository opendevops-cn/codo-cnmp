package dep

import (
	"context"
	"fmt"
	"time"

	"codo-cnmp/internal/conf"

	"github.com/go-redis/redis/extra/redisotel/v8"
	"github.com/go-redis/redis/v8"
	"github.com/gogf/gf/v2/database/gdb"
	"go.opentelemetry.io/otel/trace"

	_ "github.com/gogf/gf/contrib/drivers/mysql/v2"
	_ "github.com/gogf/gf/contrib/nosql/redis/v2"
)

func NewRedis(ctx context.Context, bc *conf.Bootstrap, provider trace.TracerProvider) *redis.Client {
	redisConf := bc.REDIS
	addr := fmt.Sprintf("%s:%d", redisConf.R_HOST, redisConf.R_PORT)
	dialTimeout := redisConf.DIAL_TIMEOUT
	if redisConf.DIAL_TIMEOUT == 0 {
		dialTimeout = 30
	}
	readTimeout := redisConf.READ_TIMEOUT
	if readTimeout == 0 {
		readTimeout = 30
	}
	writeTimeout := redisConf.WRITE_TIMEOUT
	if writeTimeout == 0 {
		writeTimeout = 30
	}
	network := redisConf.NETWORK
	if network == "" {
		network = "tcp"
	}
	client := redis.NewClient(&redis.Options{
		Network:      network,
		Addr:         addr,
		Password:     redisConf.R_PASSWORD,
		DB:           int(redisConf.R_DB),
		MaxRetries:   5,
		DialTimeout:  time.Duration(dialTimeout) * time.Second,
		ReadTimeout:  time.Duration(readTimeout) * time.Second,
		WriteTimeout: time.Duration(writeTimeout) * time.Second,
		PoolSize:     128,
	})

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-time.After(3 * time.Second):
			}
			pingCtx, cancel := context.WithTimeout(ctx, time.Duration(dialTimeout)*time.Second)
			client.Ping(pingCtx)
			cancel()
		}
	}()

	client.AddHook(redisotel.NewTracingHook(redisotel.WithTracerProvider(provider)))
	fmt.Println("redis ping >>>>>>>>>>>>>>>>", client.Ping(ctx))
	return client
}

func NewMysql(bc *conf.Bootstrap, _ trace.TracerProvider) (gdb.DB, error) {
	// database/gdb/gdb_core_underlying.go:176 直接拿的全局 trace provider ， 需要保证全局 trace provider 已经初始化, 所以引入
	DBConf := bc.DB
	dbType := DBConf.DB_TYPE
	if dbType == "" {
		dbType = "mysql"
	}
	dBMaxIdleConns := DBConf.DB_MaxIdleConns
	if dBMaxIdleConns == 0 {
		dBMaxIdleConns = 10
	}
	dBMaxOpenConns := DBConf.DB_MaxOpenConns
	if dBMaxOpenConns == 0 {
		dBMaxOpenConns = 100
	}
	dBConnMaxLifetime := DBConf.DB_ConnMaxLifetime
	if dBConnMaxLifetime == 0 {
		dBConnMaxLifetime = 300
	}
	configNode := gdb.ConfigNode{
		//Link:             mysqlConf.Link,
		Type:             dbType,
		Host:             DBConf.DB_HOST,
		Port:             DBConf.DB_PORT,
		Name:             DBConf.DB_NAME,
		User:             DBConf.DB_USER,
		Pass:             DBConf.DB_PASSWORD,
		Prefix:           DBConf.DB_TABLE_PREFIX,
		Debug:            DBConf.DEBUG,
		Timezone:         "Asia/Shanghai",
		MaxIdleConnCount: int(dBMaxIdleConns),
		MaxOpenConnCount: int(dBMaxOpenConns),
		MaxConnLifeTime:  time.Duration(dBMaxOpenConns) * time.Second,
	}
	gdb.SetConfigGroup(gdb.DefaultGroupName, gdb.ConfigGroup{configNode})
	db, err := gdb.NewByGroup(gdb.DefaultGroupName)
	if err != nil {
		return nil, err
	}
	err = db.PingMaster()
	if err != nil {
		return nil, err
	}

	return db, nil
}
