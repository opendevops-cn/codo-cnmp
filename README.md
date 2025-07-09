# 容器管理平台

# 目录结构
```text
.
├── Dockerfile
├── LICENSE
├── Makefile
├── README.md
├── etc # 配置文件
│   └── config.yaml
├── go.mod
├── go.sum
├── internal
│   ├── biz # 核心业务层
│   │   ├── README.md
│   │   ├── biz.go
│   │   └── greeter.go
│   ├── conf # 配置定义层
│   │   ├── conf.pb.go
│   │   └── conf.proto
│   ├── data # 数据存储层
│   │   ├── README.md
│   │   ├── data.go
│   │   └── greeter.go
│   ├── dep # 第三方依赖
│   │   ├── data.go
│   │   ├── otel.go
│   │   └── provider.go
│   ├── server # 运输服务层
│   │   ├── grpc.go
│   │   ├── http.go
│   │   └── server.go
│   └── service # 协议转换层
│       ├── README.md
│       ├── greeter.go
│       └── service.go
├── main.go # 函数入口
├── openapi.yaml # openaiv3 接口文档
├── pb # protobuf 定义
│   ├── greeter.pb.go
│   ├── greeter.proto
│   ├── greeter_grpc.pb.go
│   └── greeter_http.pb.go
├── third_party # proto 第三方依赖
│   ├── README.md
│   ├── google
│   │   ├── api
│   │   └── protobuf
│   ├── openapi
│   │   └── v3
│   └── validate
│       ├── README.md
│       └── validate.proto
├── wire.go # wire 文件
└── wire_gen.go
```

# 快速部署

## 修改配置
```shell
## 方法1：直接修改配置文件
cp etc/config.yaml.example etc/config.yaml # 复制配置文件
vi etc/config.yaml # 修改配置

## 方法2：使用环境变量
cp .env.example .env # 复制环境变量文件
vi .env # 修改环境变量
```

## 构建镜像
```shell
# 设置私仓token
vi git_key

# 构建镜像
docker-compose -f docker-compose.yaml build
```

## 启动
```shell
docker-compose -f docker-compose.yaml up -d
```

## License

Everything is [GPL v3.0](https://www.gnu.org/licenses/gpl-3.0.html).
