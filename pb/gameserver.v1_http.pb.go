// Code generated by protoc-gen-go-http. DO NOT EDIT.
// versions:
// - protoc-gen-go-http v2.7.3
// - protoc             v5.27.2
// source: pb/gameserver.v1.proto

package pb

import (
	context "context"
	audit "github.com/Ccheers/protoc-gen-go-kratos-http/audit"
	kcontext "github.com/Ccheers/protoc-gen-go-kratos-http/kcontext"
	route "github.com/Ccheers/protoc-gen-go-kratos-http/route"
	middleware "github.com/go-kratos/kratos/v2/middleware"
	selector "github.com/go-kratos/kratos/v2/middleware/selector"
	http "github.com/go-kratos/kratos/v2/transport/http"
	binding "github.com/go-kratos/kratos/v2/transport/http/binding"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the kratos package it is being compiled against.
var _ = new(context.Context)
var _ = binding.EncodeURL

const _ = http.SupportPackageIsVersion1

type _ = middleware.Middleware
type _ = selector.Builder
type _ = route.Route
type _ = audit.Audit

var _ = kcontext.SetKHTTPContextWithContext

const OperationGameServerListGameServer = "/gameserver.v1.GameServer/ListGameServer"
const OperationGameServerListGameServerType = "/gameserver.v1.GameServer/ListGameServerType"
const OperationGameServerManageGameServerEntity = "/gameserver.v1.GameServer/ManageGameServerEntity"
const OperationGameServerManageGameServerLB = "/gameserver.v1.GameServer/ManageGameServerLB"
const OperationGameServerBatchManageGameServerEntity = "/gameserver.v1.GameServer/BatchManageGameServerEntity"
const OperationGameServerBatchManageGameServerLB = "/gameserver.v1.GameServer/BatchManageGameServerLB"

type GameServerHTTPServer interface {
	// ListGameServer查看-云原生管理-游戏进程-列表
	ListGameServer(context.Context, *ListGameServerRequest) (*ListGameServerResponse, error)
	// ListGameServerType查看-云原生管理-游戏进程-进程类型
	ListGameServerType(context.Context, *ListGameServerTypeRequest) (*ListGameServerTypeResponse, error)
	// ManageGameServerEntity管理-云原生管理-游戏进程-Entity
	ManageGameServerEntity(context.Context, *ManageGameServerEntityRequest) (*ManageGameServerEntityResponse, error)
	// ManageGameServerLB管理-云原生管理-游戏进程-LB
	ManageGameServerLB(context.Context, *ManageGameServerEntityRequest) (*ManageGameServerEntityResponse, error)
	// BatchManageGameServerEntity管理-云原生管理-游戏进程-批量管理
	BatchManageGameServerEntity(context.Context, *BatchManageGameServerEntityRequest) (*BatchManageGameServerEntityResponse, error)
	// BatchManageGameServerLB管理-云原生管理-LB-批量管理
	BatchManageGameServerLB(context.Context, *BatchManageGameServerEntityRequest) (*BatchManageGameServerEntityResponse, error)
}

func NewGameServerHTTPServerMiddleware() middleware.Middleware {
	return selector.Server(
		selector.Server().Path(OperationGameServerListGameServer).Build(),
		selector.Server().Path(OperationGameServerListGameServerType).Build(),
		selector.Server().Path(OperationGameServerManageGameServerEntity).Build(),
		selector.Server().Path(OperationGameServerManageGameServerLB).Build(),
		selector.Server().Path(OperationGameServerBatchManageGameServerEntity).Build(),
		selector.Server().Path(OperationGameServerBatchManageGameServerLB).Build(),
	).Path(
		OperationGameServerListGameServer,
		OperationGameServerListGameServerType,
		OperationGameServerManageGameServerEntity,
		OperationGameServerManageGameServerLB,
		OperationGameServerBatchManageGameServerEntity,
		OperationGameServerBatchManageGameServerLB,
	).Build()
}

func RegisterGameServerHTTPServer(s *http.Server, srv GameServerHTTPServer) {
	r := s.Route("/")
	r.GET("/api/v1/gameserver/list", _GameServer_ListGameServer0_HTTP_Handler(srv))
	r.GET("/api/v1/gameserver/type/list", _GameServer_ListGameServerType0_HTTP_Handler(srv))
	r.POST("/api/v1/gameserver/entity/manage", _GameServer_ManageGameServerEntity0_HTTP_Handler(srv))
	r.POST("/api/v1/gameserver/lb/manage", _GameServer_ManageGameServerLB0_HTTP_Handler(srv))
	r.POST("/api/v1/gameserver/entity/batch/manage", _GameServer_BatchManageGameServerEntity0_HTTP_Handler(srv))
	r.POST("/api/v1/gameserver/lb/batch/manage", _GameServer_BatchManageGameServerLB0_HTTP_Handler(srv))
}

func GenerateGameServerHTTPServerRouteInfo() []route.Route {
	routes := make([]route.Route, 0, 6)
	routes = append(routes, route.Route{
		Method:  "GET",
		Path:    "/api/v1/gameserver/list",
		Comment: "查看-云原生管理-游戏进程-列表",
	})
	routes = append(routes, route.Route{
		Method:  "GET",
		Path:    "/api/v1/gameserver/type/list",
		Comment: "查看-云原生管理-游戏进程-进程类型",
	})
	routes = append(routes, route.Route{
		Method:  "POST",
		Path:    "/api/v1/gameserver/entity/manage",
		Comment: "管理-云原生管理-游戏进程-Entity",
	})
	routes = append(routes, route.Route{
		Method:  "POST",
		Path:    "/api/v1/gameserver/lb/manage",
		Comment: "管理-云原生管理-游戏进程-LB",
	})
	routes = append(routes, route.Route{
		Method:  "POST",
		Path:    "/api/v1/gameserver/entity/batch/manage",
		Comment: "管理-云原生管理-游戏进程-批量管理",
	})
	routes = append(routes, route.Route{
		Method:  "POST",
		Path:    "/api/v1/gameserver/lb/batch/manage",
		Comment: "管理-云原生管理-LB-批量管理",
	})
	return routes
}

func _GameServer_ListGameServer0_HTTP_Handler(srv GameServerHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		stdCtx := kcontext.SetKHTTPContextWithContext(ctx, ctx)
		var in ListGameServerRequest
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationGameServerListGameServer)

		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.ListGameServer(ctx, req.(*ListGameServerRequest))
		})
		out, err := h(stdCtx, &in)
		if err != nil {
			return err
		}
		reply := out.(*ListGameServerResponse)
		return ctx.Result(200, reply)
	}
}

func _GameServer_ListGameServerType0_HTTP_Handler(srv GameServerHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		stdCtx := kcontext.SetKHTTPContextWithContext(ctx, ctx)
		var in ListGameServerTypeRequest
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationGameServerListGameServerType)

		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.ListGameServerType(ctx, req.(*ListGameServerTypeRequest))
		})
		out, err := h(stdCtx, &in)
		if err != nil {
			return err
		}
		reply := out.(*ListGameServerTypeResponse)
		return ctx.Result(200, reply)
	}
}

func _GameServer_ManageGameServerEntity0_HTTP_Handler(srv GameServerHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		stdCtx := kcontext.SetKHTTPContextWithContext(ctx, ctx)
		var in ManageGameServerEntityRequest
		if err := ctx.Bind(&in); err != nil {
			return err
		}
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationGameServerManageGameServerEntity)
		auditRule := audit.NewAudit(
			"进程列表",
			"管理Entity",
			[]audit.Meta{
				{
					Key: "cluster",
					Value: audit.MetaValue{
						Extract: "cluster_name",
					},
				},
				{
					Key: "namespace",
					Value: audit.MetaValue{
						Extract: "namespace",
					},
				},
				{
					Key: "kind",
					Value: audit.MetaValue{
						Const: "进程",
					},
				},
				{
					Key: "name",
					Value: audit.MetaValue{
						Extract: "server_name",
					},
				},
			},
		)
		auditInfo, err := audit.ExtractFromRequest(ctx.Request(), auditRule)
		if err != nil {
			return err
		}
		stdCtx = kcontext.SetKHTTPAuditContextWithContext(stdCtx, auditInfo)

		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.ManageGameServerEntity(ctx, req.(*ManageGameServerEntityRequest))
		})
		out, err := h(stdCtx, &in)
		if err != nil {
			return err
		}
		reply := out.(*ManageGameServerEntityResponse)
		return ctx.Result(200, reply)
	}
}

func _GameServer_ManageGameServerLB0_HTTP_Handler(srv GameServerHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		stdCtx := kcontext.SetKHTTPContextWithContext(ctx, ctx)
		var in ManageGameServerEntityRequest
		if err := ctx.Bind(&in); err != nil {
			return err
		}
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationGameServerManageGameServerLB)
		auditRule := audit.NewAudit(
			"进程列表",
			"管理LB",
			[]audit.Meta{
				{
					Key: "cluster",
					Value: audit.MetaValue{
						Extract: "cluster_name",
					},
				},
				{
					Key: "namespace",
					Value: audit.MetaValue{
						Extract: "namespace",
					},
				},
				{
					Key: "kind",
					Value: audit.MetaValue{
						Const: "进程",
					},
				},
				{
					Key: "name",
					Value: audit.MetaValue{
						Extract: "server_name",
					},
				},
			},
		)
		auditInfo, err := audit.ExtractFromRequest(ctx.Request(), auditRule)
		if err != nil {
			return err
		}
		stdCtx = kcontext.SetKHTTPAuditContextWithContext(stdCtx, auditInfo)

		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.ManageGameServerLB(ctx, req.(*ManageGameServerEntityRequest))
		})
		out, err := h(stdCtx, &in)
		if err != nil {
			return err
		}
		reply := out.(*ManageGameServerEntityResponse)
		return ctx.Result(200, reply)
	}
}

func _GameServer_BatchManageGameServerEntity0_HTTP_Handler(srv GameServerHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		stdCtx := kcontext.SetKHTTPContextWithContext(ctx, ctx)
		var in BatchManageGameServerEntityRequest
		if err := ctx.Bind(&in); err != nil {
			return err
		}
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationGameServerBatchManageGameServerEntity)
		auditRule := audit.NewAudit(
			"进程列表",
			"批量管理Entity",
			[]audit.Meta{
				{
					Key: "cluster",
					Value: audit.MetaValue{
						Extract: "cluster_name",
					},
				},
				{
					Key: "namespace",
					Value: audit.MetaValue{
						Extract: "namespace",
					},
				},
				{
					Key: "kind",
					Value: audit.MetaValue{
						Const: "进程",
					},
				},
				{
					Key: "name",
					Value: audit.MetaValue{
						Extract: "server_name",
					},
				},
			},
		)
		auditInfo, err := audit.ExtractFromRequest(ctx.Request(), auditRule)
		if err != nil {
			return err
		}
		stdCtx = kcontext.SetKHTTPAuditContextWithContext(stdCtx, auditInfo)

		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.BatchManageGameServerEntity(ctx, req.(*BatchManageGameServerEntityRequest))
		})
		out, err := h(stdCtx, &in)
		if err != nil {
			return err
		}
		reply := out.(*BatchManageGameServerEntityResponse)
		return ctx.Result(200, reply)
	}
}

func _GameServer_BatchManageGameServerLB0_HTTP_Handler(srv GameServerHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		stdCtx := kcontext.SetKHTTPContextWithContext(ctx, ctx)
		var in BatchManageGameServerEntityRequest
		if err := ctx.Bind(&in); err != nil {
			return err
		}
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationGameServerBatchManageGameServerLB)
		auditRule := audit.NewAudit(
			"进程列表",
			"批量管理LB",
			[]audit.Meta{
				{
					Key: "cluster",
					Value: audit.MetaValue{
						Extract: "cluster_name",
					},
				},
				{
					Key: "namespace",
					Value: audit.MetaValue{
						Extract: "namespace",
					},
				},
				{
					Key: "kind",
					Value: audit.MetaValue{
						Const: "进程",
					},
				},
				{
					Key: "name",
					Value: audit.MetaValue{
						Extract: "server_name",
					},
				},
			},
		)
		auditInfo, err := audit.ExtractFromRequest(ctx.Request(), auditRule)
		if err != nil {
			return err
		}
		stdCtx = kcontext.SetKHTTPAuditContextWithContext(stdCtx, auditInfo)

		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.BatchManageGameServerLB(ctx, req.(*BatchManageGameServerEntityRequest))
		})
		out, err := h(stdCtx, &in)
		if err != nil {
			return err
		}
		reply := out.(*BatchManageGameServerEntityResponse)
		return ctx.Result(200, reply)
	}
}

type GameServerHTTPClient interface {
	ListGameServer(ctx context.Context, req *ListGameServerRequest, opts ...http.CallOption) (rsp *ListGameServerResponse, err error)
	ListGameServerType(ctx context.Context, req *ListGameServerTypeRequest, opts ...http.CallOption) (rsp *ListGameServerTypeResponse, err error)
	ManageGameServerEntity(ctx context.Context, req *ManageGameServerEntityRequest, opts ...http.CallOption) (rsp *ManageGameServerEntityResponse, err error)
	ManageGameServerLB(ctx context.Context, req *ManageGameServerEntityRequest, opts ...http.CallOption) (rsp *ManageGameServerEntityResponse, err error)
	BatchManageGameServerEntity(ctx context.Context, req *BatchManageGameServerEntityRequest, opts ...http.CallOption) (rsp *BatchManageGameServerEntityResponse, err error)
	BatchManageGameServerLB(ctx context.Context, req *BatchManageGameServerEntityRequest, opts ...http.CallOption) (rsp *BatchManageGameServerEntityResponse, err error)
}

type GameServerHTTPClientImpl struct {
	cc *http.Client
}

func NewGameServerHTTPClient(client *http.Client) GameServerHTTPClient {
	return &GameServerHTTPClientImpl{client}
}

func (c *GameServerHTTPClientImpl) ListGameServer(ctx context.Context, in *ListGameServerRequest, opts ...http.CallOption) (*ListGameServerResponse, error) {
	var out ListGameServerResponse
	pattern := "/api/v1/gameserver/list"
	path := binding.EncodeURL(pattern, in, true)
	opts = append(opts, http.Operation(OperationGameServerListGameServer))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "GET", path, nil, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *GameServerHTTPClientImpl) ListGameServerType(ctx context.Context, in *ListGameServerTypeRequest, opts ...http.CallOption) (*ListGameServerTypeResponse, error) {
	var out ListGameServerTypeResponse
	pattern := "/api/v1/gameserver/type/list"
	path := binding.EncodeURL(pattern, in, true)
	opts = append(opts, http.Operation(OperationGameServerListGameServerType))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "GET", path, nil, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *GameServerHTTPClientImpl) ManageGameServerEntity(ctx context.Context, in *ManageGameServerEntityRequest, opts ...http.CallOption) (*ManageGameServerEntityResponse, error) {
	var out ManageGameServerEntityResponse
	pattern := "/api/v1/gameserver/entity/manage"
	path := binding.EncodeURL(pattern, in, false)
	opts = append(opts, http.Operation(OperationGameServerManageGameServerEntity))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "POST", path, in, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *GameServerHTTPClientImpl) ManageGameServerLB(ctx context.Context, in *ManageGameServerEntityRequest, opts ...http.CallOption) (*ManageGameServerEntityResponse, error) {
	var out ManageGameServerEntityResponse
	pattern := "/api/v1/gameserver/lb/manage"
	path := binding.EncodeURL(pattern, in, false)
	opts = append(opts, http.Operation(OperationGameServerManageGameServerLB))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "POST", path, in, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *GameServerHTTPClientImpl) BatchManageGameServerEntity(ctx context.Context, in *BatchManageGameServerEntityRequest, opts ...http.CallOption) (*BatchManageGameServerEntityResponse, error) {
	var out BatchManageGameServerEntityResponse
	pattern := "/api/v1/gameserver/entity/batch/manage"
	path := binding.EncodeURL(pattern, in, false)
	opts = append(opts, http.Operation(OperationGameServerBatchManageGameServerEntity))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "POST", path, in, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *GameServerHTTPClientImpl) BatchManageGameServerLB(ctx context.Context, in *BatchManageGameServerEntityRequest, opts ...http.CallOption) (*BatchManageGameServerEntityResponse, error) {
	var out BatchManageGameServerEntityResponse
	pattern := "/api/v1/gameserver/lb/batch/manage"
	path := binding.EncodeURL(pattern, in, false)
	opts = append(opts, http.Operation(OperationGameServerBatchManageGameServerLB))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "POST", path, in, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, nil
}
