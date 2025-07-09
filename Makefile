# 下面是一些自定义的 Makefile
#include ../../../../Makefile
GOHOSTOS:=$(shell go env GOHOSTOS)
GOPATH:=$(shell go env GOPATH)
VERSION=$(shell git describe --tags --always)
GOVERSION := $(shell go env GOVERSION)
GITCOMMIT := $(shell git rev-parse --short HEAD)
BUILT := $(shell date -u '+%Y-%m-%dT%H:%M:%SZ')
OSARCH := `go env GOOS`/`go env GOARCH`
ENV:=pre
KPIDAY:=30
GO_BIN_OUTPUT:=./bin/codo-cnmp
GO_CODE_DIR:=.
GOLDFLAGS = -X 'codo-cnmp/meta.Version=${BUIlDTAG}' \
	-X 'codo-cnmp/meta.GoVersion=${GOVERSION}' \
	-X 'codo-cnmp/meta.GitCommit=${GITCOMMIT}' \
	-X 'codo-cnmp/meta.Built=${BUILT}' \
	-X 'codo-cnmp/meta.ENV=${ENV}' \
	-X 'codo-cnmp/meta.OsArch=${OSARCH}'
ifeq ($(GOHOSTOS), windows)
	#the `find.exe` is different from `find` in bash/shell.
	#to see https://docs.microsoft.com/en-us/windows-server/administration/windows-commands/find.
	#changed to use git-bash.exe to run find cli or other cli friendly, caused of every developer has a Git.
	#Git_Bash= $(subst cmd\,bin\bash.exe,$(dir $(shell where git)))
#	Git_Bash=$(subst \,/,$(subst mingw64\bin\,bin\bash.exe,$(dir $(shell where git |head -n 1))))
	Git_Bash=/usr/bin/bash
	INTERNAL_PROTO_FILES=$(shell $(Git_Bash) -c "find ./internal -name *.proto")
	API_PROTO_FILES=$(shell $(Git_Bash) -c "find ./pb -name *.proto")
else
	INTERNAL_PROTO_FILES=$(shell find ./internal -name *.proto)
	API_PROTO_FILES=$(shell find ./pb -name *.proto)
endif
.PHONY: init
# init env
init:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.31.0
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.3.0
	go install github.com/google/wire/cmd/wire@latest
	go install github.com/zeromicro/go-zero/tools/goctl@v1.5.2
	go install github.com/zeromicro/goctl-swagger@latest
	go install github.com/git-chglog/git-chglog/cmd/git-chglog@latest
	go install mvdan.cc/gofumpt@latest
	go install helm.sh/helm/v3/cmd/helm@latest
	go install github.com/Ccheers/protoc-gen-zeroapi@v0.0.10
	go install github.com/Ccheers/protoc-gen-zerorpc@v0.0.5
	go install github.com/envoyproxy/protoc-gen-validate@v1.0.2
	go install github.com/google/gnostic/cmd/protoc-gen-openapi@latest
	go install github.com/Ccheers/grpc-gateway/protoc-gen-openapiv2@v0.0.2-unofficial
	go install github.com/favadi/protoc-go-inject-tag@latest
	go install go.uber.org/nilaway/cmd/nilaway@latest
	go install github.com/go-kratos/kratos/cmd/kratos/v2@latest
	go install github.com/Ccheers/protoc-gen-go-kratos-http@v0.0.11
	go install golang.org/x/tools/cmd/goimports@latest
	go clone
.PHONY: config
# generate internal proto
config:
	protoc --proto_path=./internal \
	       --proto_path=./third_party \
 	       --go_out=paths=source_relative:./internal \
	       ./internal/conf/conf.proto
	protoc-go-inject-tag -input="./internal/conf/conf.pb.go"
.PHONY: lint
# 格式化代码
lint:
	nilaway ./...
.PHONY: fmt
# 格式化代码
fmt:
	gofumpt -l -w .
	goimports -l -w .
.PHONY: gen_dao
# 生成 dao 层，依赖于 common/hack/config.yaml.bak
gen_dao:
	cd common/hack && gf gen dao
.PHONY: changelog
# 生成 changelog
changelog:
	git-chglog -o ./CHANGELOG.md
# 查看最近写了多少代码
.PHONY: build
build:
	go build -ldflags "-s -w $(GOLDFLAGS)" -o $(GO_BIN_OUTPUT) $(GO_CODE_DIR)

.PHONY: kpi
kpi:
ifeq ($(shell uname), Darwin)
	git log --since="`date -v -$(KPIDAY)d +"%Y-%m-%d"`" --before="`date +"%Y-%m-%d"`" --author="`git config --get user.name`" --pretty=tformat: --numstat -- . ":(exclude)*.json" ":(exclude)*.yaml" ":(exclude)*_route.go" ":(exclude)*_handler.go" ":(exclude)app/*/cmd/rpc/client/*" ":(exclude)app/*/cmd/*/pb/*" ":(exclude)app/*/cmd/api/internal/handler/*" ":(exclude)app/*/cmd/rpc/internal/server/*" ":(exclude)app/*/model/*" | awk '{ add += $$1; subs += $$2; loc += $$1 - $$2 } END { printf "added lines: %s removed lines : %s total lines: %s\n",add,subs,loc }'
else
	git log --since="`date -d "$(KPIDAY) day ago" +"%Y-%m-%d"`" --before="`date +"%Y-%m-%d"`" --author="`git config --get user.name`" --pretty=tformat: --numstat -- . ":(exclude)*.json" ":(exclude)*.yaml" ":(exclude)*_route.go" ":(exclude)*_handler.go" ":(exclude)app/*/cmd/rpc/client/*" ":(exclude)app/*/cmd/*/pb/*" ":(exclude)app/*/cmd/api/internal/handler/*" ":(exclude)app/*/cmd/rpc/internal/server/*" ":(exclude)app/*/model/*" | awk '{ add += $$1; subs += $$2; loc += $$1 - $$2 } END { printf "added lines: %s removed lines : %s total lines: %s\n",add,subs,loc }'
endif
.PHONY: api
# 生成 api 层，需要在 app/{service}/cmd/rpc/pb 目录执行
api:
	protoc --proto_path=. \
	--proto_path=./third_party \
	--go-kratos-http_out=paths=source_relative:. \
	--go_out=paths=source_relative:. \
	--go-grpc_out=paths=source_relative:. \
    --openapiv2_out=enums_as_ints=false,json_names_for_fields=false:. \
    --openapi_out=naming=proto:. \
    --validate_out=lang=go,paths=source_relative:. \
	$(API_PROTO_FILES)
ifeq ($(GOHOSTOS), windows)
	find ./pb -name '*.pb.go' -exec  sed -i -e "s/,omitempty/,optional/g" {} \;
	# find ./pb -name '*.swagger.json' -exec  sed -i -e "s/\"title\": \"/\"description\": \"/g" {} \;
else
	find ./pb -name '*.pb.go' -exec  sed -i "" -e "s/,omitempty/,optional/g" {} \;
	# find ./pb -name '*.swagger.json' -exec  sed -i -e "s/\"title\": \"/\"description\": \"/g" {} \;
endif
	protoc-go-inject-tag -input="./pb/*.pb.go"
# show help
help:
	@echo ''
	@echo 'Usage:'
	@echo ' make [target]'
	@echo ''
	@echo 'Targets:'
	@awk '/^[a-zA-Z\-_0-9]+:/ { \
	helpMessage = match(lastLine, /^# (.*)/); \
		if (helpMessage) { \
			helpCommand = substr($$1, 0, index($$1, ":")-1); \
			helpMessage = substr(lastLine, RSTART + 2, RLENGTH); \
			printf "\033[36m%-22s\033[0m %s\n", helpCommand,helpMessage; \
		} \
	} \
	{ lastLine = $$0 }' $(MAKEFILE_LIST)
.DEFAULT_GOAL := help