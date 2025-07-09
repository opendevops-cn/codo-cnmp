package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"codo-cnmp/initialization"
	"codo-cnmp/internal/dep"
	"codo-cnmp/migrate"
	"github.com/opendevops-cn/codo-golang-sdk/adapter/kratos/transport/websocket"
	configsdk "github.com/opendevops-cn/codo-golang-sdk/config"

	"codo-cnmp/internal/server"
	"github.com/go-kratos/kratos/v2"

	"codo-cnmp/internal/conf"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport/http"

	_ "go.uber.org/automaxprocs"
)

// go build -ldflags "-X main.Version=x.y.z"
var (
	// Version is the version of the compiled software.
	Version string
	// flagconf is the config flag.
	flagconf string

	id, _ = os.Hostname()
)

func init() {
	flag.StringVar(&flagconf, "f", "etc/config.yaml", "config path, eg: -conf config.yaml")
}

func newApp(bc *conf.Bootstrap, logger log.Logger, hs *http.Server,
	pprof *server.PprofServer, prom *server.PrometheusServer,
	cs *server.CronServer, informer *server.InformerServerWrapper, websocket *websocket.Server,
	metricsInformer *server.MetricsInformerServerWrapper, migrate *migrate.Migration,
	registration *initialization.RegistrationMeta, proxy *server.APIServerProxy,
) (*kratos.App, error) {
	// 只注册 grpc 的 endpoint
	endpoint, err := hs.Endpoint()
	if err != nil {
		return nil, err
	}
	app := bc.APP
	err = migrate.Run(context.Background())
	if err != nil {
		fmt.Println(err)
	}
	err = registration.Run(context.Background())
	if err != nil {
		fmt.Println(err)
	}
	crd := dep.NewCRD()
	crd.Register()

	return kratos.New(
		kratos.ID(id),
		kratos.Name(app.NAME),
		kratos.Version(app.VERSION),
		kratos.Metadata(map[string]string{}),
		kratos.Logger(logger),
		kratos.Server(
			// gs,
			hs,              // http 服务
			pprof,           // 性能
			prom,            // 指标
			cs,              // 定时任务
			informer,        // informer
			metricsInformer, // metrics informer
			// agent,           // agent
			websocket, // websocket
			proxy,     // proxy
		),
		kratos.Endpoint(endpoint),
	), nil
}

func IsValid(bc *conf.Bootstrap) bool {
	return bc.APP != nil && bc.APP.ADDR != "" && // 确保 APP 不为 nil 且包含有效地址
		bc.DB != nil && bc.DB.DB_HOST != "" && // 检查数据库配置
		bc.REDIS != nil && bc.REDIS.R_HOST != "" // 检查 Redis 配置
}

func main() {
	flag.Parse()
	var bc conf.Bootstrap
	// 优先加载环墶变量
	err := configsdk.LoadConfig(&bc,
		configsdk.WithEnv(""),
	)
	if err != nil || !IsValid(&bc) {
		fmt.Println("加载环境变量失败，尝试加载配置文件")
		err = configsdk.LoadConfig(&bc,
			configsdk.WithYaml(flagconf),
		)
	}
	if err != nil || !IsValid(&bc) {
		panic("加载配置失败: 环境变量和配置文件均无效")
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	app, cleanup, err := wireApp(ctx, &bc)
	if err != nil {
		panic(err)
	}
	defer cleanup()

	// start and wait for stop signal
	if err := app.Run(); err != nil {
		panic(err)
	}
}
