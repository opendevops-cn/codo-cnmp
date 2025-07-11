// Code generated by protoc-gen-go-http. DO NOT EDIT.
// versions:
// - protoc-gen-go-http v2.7.3
// - protoc             v5.27.2
// source: pb/sidecarset.v1.proto

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

const OperationSidecarSetListSidecarSet = "/sidecarset.v1.SidecarSet/ListSidecarSet"
const OperationSidecarSetGetSidecarSet = "/sidecarset.v1.SidecarSet/GetSidecarSet"
const OperationSidecarSetUpdateSideCarSetStrategy = "/sidecarset.v1.SidecarSet/UpdateSideCarSetStrategy"
const OperationSidecarSetDeleteSidecarSet = "/sidecarset.v1.SidecarSet/DeleteSidecarSet"

type SidecarSetHTTPServer interface {
	// ListSidecarSet查看-云原生管理-SideCarSet-列表
	ListSidecarSet(context.Context, *ListSidecarSetRequest) (*ListSidecarSetResponse, error)
	// GetSidecarSet查看-云原生管理-SideCarSet-详情
	GetSidecarSet(context.Context, *GetSidecarSetRequest) (*GetSidecarSetResponse, error)
	// UpdateSideCarSetStrategy管理-云原生管理-SideCarSet-更新策略
	UpdateSideCarSetStrategy(context.Context, *UpdateSideCarSetStrategyRequest) (*UpdateSideCarSetStrategyResponse, error)
	// DeleteSidecarSet管理-云原生管理-SideCarSet-删除
	DeleteSidecarSet(context.Context, *DeleteSidecarSetRequest) (*DeleteSidecarSetResponse, error)
}

func NewSidecarSetHTTPServerMiddleware() middleware.Middleware {
	return selector.Server(
		selector.Server().Path(OperationSidecarSetListSidecarSet).Build(),
		selector.Server().Path(OperationSidecarSetGetSidecarSet).Build(),
		selector.Server().Path(OperationSidecarSetUpdateSideCarSetStrategy).Build(),
		selector.Server().Path(OperationSidecarSetDeleteSidecarSet).Build(),
	).Path(
		OperationSidecarSetListSidecarSet,
		OperationSidecarSetGetSidecarSet,
		OperationSidecarSetUpdateSideCarSetStrategy,
		OperationSidecarSetDeleteSidecarSet,
	).Build()
}

func RegisterSidecarSetHTTPServer(s *http.Server, srv SidecarSetHTTPServer) {
	r := s.Route("/")
	r.GET("/api/v1/sidecarset/list", _SidecarSet_ListSidecarSet0_HTTP_Handler(srv))
	r.GET("/api/v1/sidecarset/detail", _SidecarSet_GetSidecarSet0_HTTP_Handler(srv))
	r.POST("/api/v1/sidecarset/upgrade_strategy/update", _SidecarSet_UpdateSideCarSetStrategy0_HTTP_Handler(srv))
	r.POST("/api/v1/sidecarset/delete", _SidecarSet_DeleteSidecarSet0_HTTP_Handler(srv))
}

func GenerateSidecarSetHTTPServerRouteInfo() []route.Route {
	routes := make([]route.Route, 0, 4)
	routes = append(routes, route.Route{
		Method:  "GET",
		Path:    "/api/v1/sidecarset/list",
		Comment: "查看-云原生管理-SideCarSet-列表",
	})
	routes = append(routes, route.Route{
		Method:  "GET",
		Path:    "/api/v1/sidecarset/detail",
		Comment: "查看-云原生管理-SideCarSet-详情",
	})
	routes = append(routes, route.Route{
		Method:  "POST",
		Path:    "/api/v1/sidecarset/upgrade_strategy/update",
		Comment: "管理-云原生管理-SideCarSet-更新策略",
	})
	routes = append(routes, route.Route{
		Method:  "POST",
		Path:    "/api/v1/sidecarset/delete",
		Comment: "管理-云原生管理-SideCarSet-删除",
	})
	return routes
}

func _SidecarSet_ListSidecarSet0_HTTP_Handler(srv SidecarSetHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		stdCtx := kcontext.SetKHTTPContextWithContext(ctx, ctx)
		var in ListSidecarSetRequest
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationSidecarSetListSidecarSet)

		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.ListSidecarSet(ctx, req.(*ListSidecarSetRequest))
		})
		out, err := h(stdCtx, &in)
		if err != nil {
			return err
		}
		reply := out.(*ListSidecarSetResponse)
		return ctx.Result(200, reply)
	}
}

func _SidecarSet_GetSidecarSet0_HTTP_Handler(srv SidecarSetHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		stdCtx := kcontext.SetKHTTPContextWithContext(ctx, ctx)
		var in GetSidecarSetRequest
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationSidecarSetGetSidecarSet)

		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.GetSidecarSet(ctx, req.(*GetSidecarSetRequest))
		})
		out, err := h(stdCtx, &in)
		if err != nil {
			return err
		}
		reply := out.(*GetSidecarSetResponse)
		return ctx.Result(200, reply)
	}
}

func _SidecarSet_UpdateSideCarSetStrategy0_HTTP_Handler(srv SidecarSetHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		stdCtx := kcontext.SetKHTTPContextWithContext(ctx, ctx)
		var in UpdateSideCarSetStrategyRequest
		if err := ctx.Bind(&in); err != nil {
			return err
		}
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationSidecarSetUpdateSideCarSetStrategy)
		auditRule := audit.NewAudit(
			"sidecarSet",
			"编辑升级策略",
			[]audit.Meta{
				{
					Key: "cluster",
					Value: audit.MetaValue{
						Extract: "cluster_name",
					},
				},
				{
					Key:   "namespace",
					Value: audit.MetaValue{},
				},
				{
					Key: "kind",
					Value: audit.MetaValue{
						Const: "sidecarSet",
					},
				},
				{
					Key: "name",
					Value: audit.MetaValue{
						Extract: "name",
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
			return srv.UpdateSideCarSetStrategy(ctx, req.(*UpdateSideCarSetStrategyRequest))
		})
		out, err := h(stdCtx, &in)
		if err != nil {
			return err
		}
		reply := out.(*UpdateSideCarSetStrategyResponse)
		return ctx.Result(200, reply)
	}
}

func _SidecarSet_DeleteSidecarSet0_HTTP_Handler(srv SidecarSetHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		stdCtx := kcontext.SetKHTTPContextWithContext(ctx, ctx)
		var in DeleteSidecarSetRequest
		if err := ctx.Bind(&in); err != nil {
			return err
		}
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationSidecarSetDeleteSidecarSet)
		auditRule := audit.NewAudit(
			"sidecarSet",
			"删除",
			[]audit.Meta{
				{
					Key: "cluster",
					Value: audit.MetaValue{
						Extract: "cluster_name",
					},
				},
				{
					Key:   "namespace",
					Value: audit.MetaValue{},
				},
				{
					Key: "kind",
					Value: audit.MetaValue{
						Const: "sidecarSet",
					},
				},
				{
					Key: "name",
					Value: audit.MetaValue{
						Extract: "name",
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
			return srv.DeleteSidecarSet(ctx, req.(*DeleteSidecarSetRequest))
		})
		out, err := h(stdCtx, &in)
		if err != nil {
			return err
		}
		reply := out.(*DeleteSidecarSetResponse)
		return ctx.Result(200, reply)
	}
}

type SidecarSetHTTPClient interface {
	ListSidecarSet(ctx context.Context, req *ListSidecarSetRequest, opts ...http.CallOption) (rsp *ListSidecarSetResponse, err error)
	GetSidecarSet(ctx context.Context, req *GetSidecarSetRequest, opts ...http.CallOption) (rsp *GetSidecarSetResponse, err error)
	UpdateSideCarSetStrategy(ctx context.Context, req *UpdateSideCarSetStrategyRequest, opts ...http.CallOption) (rsp *UpdateSideCarSetStrategyResponse, err error)
	DeleteSidecarSet(ctx context.Context, req *DeleteSidecarSetRequest, opts ...http.CallOption) (rsp *DeleteSidecarSetResponse, err error)
}

type SidecarSetHTTPClientImpl struct {
	cc *http.Client
}

func NewSidecarSetHTTPClient(client *http.Client) SidecarSetHTTPClient {
	return &SidecarSetHTTPClientImpl{client}
}

func (c *SidecarSetHTTPClientImpl) ListSidecarSet(ctx context.Context, in *ListSidecarSetRequest, opts ...http.CallOption) (*ListSidecarSetResponse, error) {
	var out ListSidecarSetResponse
	pattern := "/api/v1/sidecarset/list"
	path := binding.EncodeURL(pattern, in, true)
	opts = append(opts, http.Operation(OperationSidecarSetListSidecarSet))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "GET", path, nil, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *SidecarSetHTTPClientImpl) GetSidecarSet(ctx context.Context, in *GetSidecarSetRequest, opts ...http.CallOption) (*GetSidecarSetResponse, error) {
	var out GetSidecarSetResponse
	pattern := "/api/v1/sidecarset/detail"
	path := binding.EncodeURL(pattern, in, true)
	opts = append(opts, http.Operation(OperationSidecarSetGetSidecarSet))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "GET", path, nil, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *SidecarSetHTTPClientImpl) UpdateSideCarSetStrategy(ctx context.Context, in *UpdateSideCarSetStrategyRequest, opts ...http.CallOption) (*UpdateSideCarSetStrategyResponse, error) {
	var out UpdateSideCarSetStrategyResponse
	pattern := "/api/v1/sidecarset/upgrade_strategy/update"
	path := binding.EncodeURL(pattern, in, false)
	opts = append(opts, http.Operation(OperationSidecarSetUpdateSideCarSetStrategy))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "POST", path, in, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *SidecarSetHTTPClientImpl) DeleteSidecarSet(ctx context.Context, in *DeleteSidecarSetRequest, opts ...http.CallOption) (*DeleteSidecarSetResponse, error) {
	var out DeleteSidecarSetResponse
	pattern := "/api/v1/sidecarset/delete"
	path := binding.EncodeURL(pattern, in, false)
	opts = append(opts, http.Operation(OperationSidecarSetDeleteSidecarSet))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "POST", path, in, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, nil
}
