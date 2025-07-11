// Code generated by protoc-gen-go-http. DO NOT EDIT.
// versions:
// - protoc-gen-go-http v2.7.3
// - protoc             v5.27.2
// source: pb/limitrange.v1.proto

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

const OperationLimitRangeListLimitRange = "/limitrange.v1.LimitRange/ListLimitRange"
const OperationLimitRangeCreateLimitRange = "/limitrange.v1.LimitRange/CreateLimitRange"
const OperationLimitRangeUpdateLimitRange = "/limitrange.v1.LimitRange/UpdateLimitRange"
const OperationLimitRangeCreateOrUpdateLimitRange = "/limitrange.v1.LimitRange/CreateOrUpdateLimitRange"
const OperationLimitRangeGetLimitRangeDetail = "/limitrange.v1.LimitRange/GetLimitRangeDetail"
const OperationLimitRangeDeleteLimitRange = "/limitrange.v1.LimitRange/DeleteLimitRange"

type LimitRangeHTTPServer interface {
	// ListLimitRange查看-云原生管理-LimitRange-列表
	ListLimitRange(context.Context, *ListLimitRangeRequest) (*ListLimitRangeResponse, error)
	// CreateLimitRange管理-云原生管理-LimitRange-创建
	CreateLimitRange(context.Context, *CreateLimitRangeRequest) (*CreateLimitRangeResponse, error)
	// UpdateLimitRange管理-云原生管理-LimitRange-编辑
	UpdateLimitRange(context.Context, *CreateLimitRangeRequest) (*CreateLimitRangeResponse, error)
	// CreateOrUpdateLimitRange管理-云原生管理-LimitRange-创建或编辑
	CreateOrUpdateLimitRange(context.Context, *CreateLimitRangeRequest) (*CreateLimitRangeResponse, error)
	// GetLimitRangeDetail查看-云原生管理-LimitRange-详情
	GetLimitRangeDetail(context.Context, *LimitRangeDetailRequest) (*LimitRangeDetailResponse, error)
	// DeleteLimitRange管理-云原生管理-LimitRange-删除
	DeleteLimitRange(context.Context, *DeleteLimitRangeRequest) (*DeleteLimitRangeResponse, error)
}

func NewLimitRangeHTTPServerMiddleware() middleware.Middleware {
	return selector.Server(
		selector.Server().Path(OperationLimitRangeListLimitRange).Build(),
		selector.Server().Path(OperationLimitRangeCreateLimitRange).Build(),
		selector.Server().Path(OperationLimitRangeUpdateLimitRange).Build(),
		selector.Server().Path(OperationLimitRangeCreateOrUpdateLimitRange).Build(),
		selector.Server().Path(OperationLimitRangeGetLimitRangeDetail).Build(),
		selector.Server().Path(OperationLimitRangeDeleteLimitRange).Build(),
	).Path(
		OperationLimitRangeListLimitRange,
		OperationLimitRangeCreateLimitRange,
		OperationLimitRangeUpdateLimitRange,
		OperationLimitRangeCreateOrUpdateLimitRange,
		OperationLimitRangeGetLimitRangeDetail,
		OperationLimitRangeDeleteLimitRange,
	).Build()
}

func RegisterLimitRangeHTTPServer(s *http.Server, srv LimitRangeHTTPServer) {
	r := s.Route("/")
	r.GET("/api/v1/limitrange/list", _LimitRange_ListLimitRange0_HTTP_Handler(srv))
	r.POST("/api/v1/limitrange/create", _LimitRange_CreateLimitRange0_HTTP_Handler(srv))
	r.POST("/api/v1/limitrange/update", _LimitRange_UpdateLimitRange0_HTTP_Handler(srv))
	r.POST("/api/v1/limitrange/create_or_update", _LimitRange_CreateOrUpdateLimitRange0_HTTP_Handler(srv))
	r.GET("/api/v1/limitrange/detail", _LimitRange_GetLimitRangeDetail0_HTTP_Handler(srv))
	r.POST("/api/v1/limitrange/delete", _LimitRange_DeleteLimitRange0_HTTP_Handler(srv))
}

func GenerateLimitRangeHTTPServerRouteInfo() []route.Route {
	routes := make([]route.Route, 0, 6)
	routes = append(routes, route.Route{
		Method:  "GET",
		Path:    "/api/v1/limitrange/list",
		Comment: "查看-云原生管理-LimitRange-列表",
	})
	routes = append(routes, route.Route{
		Method:  "POST",
		Path:    "/api/v1/limitrange/create",
		Comment: "管理-云原生管理-LimitRange-创建",
	})
	routes = append(routes, route.Route{
		Method:  "POST",
		Path:    "/api/v1/limitrange/update",
		Comment: "管理-云原生管理-LimitRange-编辑",
	})
	routes = append(routes, route.Route{
		Method:  "POST",
		Path:    "/api/v1/limitrange/create_or_update",
		Comment: "管理-云原生管理-LimitRange-创建或编辑",
	})
	routes = append(routes, route.Route{
		Method:  "GET",
		Path:    "/api/v1/limitrange/detail",
		Comment: "查看-云原生管理-LimitRange-详情",
	})
	routes = append(routes, route.Route{
		Method:  "POST",
		Path:    "/api/v1/limitrange/delete",
		Comment: "管理-云原生管理-LimitRange-删除",
	})
	return routes
}

func _LimitRange_ListLimitRange0_HTTP_Handler(srv LimitRangeHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		stdCtx := kcontext.SetKHTTPContextWithContext(ctx, ctx)
		var in ListLimitRangeRequest
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationLimitRangeListLimitRange)

		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.ListLimitRange(ctx, req.(*ListLimitRangeRequest))
		})
		out, err := h(stdCtx, &in)
		if err != nil {
			return err
		}
		reply := out.(*ListLimitRangeResponse)
		return ctx.Result(200, reply)
	}
}

func _LimitRange_CreateLimitRange0_HTTP_Handler(srv LimitRangeHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		stdCtx := kcontext.SetKHTTPContextWithContext(ctx, ctx)
		var in CreateLimitRangeRequest
		if err := ctx.Bind(&in); err != nil {
			return err
		}
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationLimitRangeCreateLimitRange)
		auditRule := audit.NewAudit(
			"LimitRange",
			"创建",
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
						Const: "LimitRange",
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
			return srv.CreateLimitRange(ctx, req.(*CreateLimitRangeRequest))
		})
		out, err := h(stdCtx, &in)
		if err != nil {
			return err
		}
		reply := out.(*CreateLimitRangeResponse)
		return ctx.Result(200, reply)
	}
}

func _LimitRange_UpdateLimitRange0_HTTP_Handler(srv LimitRangeHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		stdCtx := kcontext.SetKHTTPContextWithContext(ctx, ctx)
		var in CreateLimitRangeRequest
		if err := ctx.Bind(&in); err != nil {
			return err
		}
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationLimitRangeUpdateLimitRange)
		auditRule := audit.NewAudit(
			"LimitRange",
			"编辑",
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
						Const: "LimitRange",
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
			return srv.UpdateLimitRange(ctx, req.(*CreateLimitRangeRequest))
		})
		out, err := h(stdCtx, &in)
		if err != nil {
			return err
		}
		reply := out.(*CreateLimitRangeResponse)
		return ctx.Result(200, reply)
	}
}

func _LimitRange_CreateOrUpdateLimitRange0_HTTP_Handler(srv LimitRangeHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		stdCtx := kcontext.SetKHTTPContextWithContext(ctx, ctx)
		var in CreateLimitRangeRequest
		if err := ctx.Bind(&in); err != nil {
			return err
		}
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationLimitRangeCreateOrUpdateLimitRange)
		auditRule := audit.NewAudit(
			"LimitRange",
			"创建或编辑",
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
						Const: "LimitRange",
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
			return srv.CreateOrUpdateLimitRange(ctx, req.(*CreateLimitRangeRequest))
		})
		out, err := h(stdCtx, &in)
		if err != nil {
			return err
		}
		reply := out.(*CreateLimitRangeResponse)
		return ctx.Result(200, reply)
	}
}

func _LimitRange_GetLimitRangeDetail0_HTTP_Handler(srv LimitRangeHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		stdCtx := kcontext.SetKHTTPContextWithContext(ctx, ctx)
		var in LimitRangeDetailRequest
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationLimitRangeGetLimitRangeDetail)

		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.GetLimitRangeDetail(ctx, req.(*LimitRangeDetailRequest))
		})
		out, err := h(stdCtx, &in)
		if err != nil {
			return err
		}
		reply := out.(*LimitRangeDetailResponse)
		return ctx.Result(200, reply)
	}
}

func _LimitRange_DeleteLimitRange0_HTTP_Handler(srv LimitRangeHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		stdCtx := kcontext.SetKHTTPContextWithContext(ctx, ctx)
		var in DeleteLimitRangeRequest
		if err := ctx.Bind(&in); err != nil {
			return err
		}
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationLimitRangeDeleteLimitRange)
		auditRule := audit.NewAudit(
			"LimitRange",
			"删除",
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
						Const: "LimitRange",
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
			return srv.DeleteLimitRange(ctx, req.(*DeleteLimitRangeRequest))
		})
		out, err := h(stdCtx, &in)
		if err != nil {
			return err
		}
		reply := out.(*DeleteLimitRangeResponse)
		return ctx.Result(200, reply)
	}
}

type LimitRangeHTTPClient interface {
	ListLimitRange(ctx context.Context, req *ListLimitRangeRequest, opts ...http.CallOption) (rsp *ListLimitRangeResponse, err error)
	CreateLimitRange(ctx context.Context, req *CreateLimitRangeRequest, opts ...http.CallOption) (rsp *CreateLimitRangeResponse, err error)
	UpdateLimitRange(ctx context.Context, req *CreateLimitRangeRequest, opts ...http.CallOption) (rsp *CreateLimitRangeResponse, err error)
	CreateOrUpdateLimitRange(ctx context.Context, req *CreateLimitRangeRequest, opts ...http.CallOption) (rsp *CreateLimitRangeResponse, err error)
	GetLimitRangeDetail(ctx context.Context, req *LimitRangeDetailRequest, opts ...http.CallOption) (rsp *LimitRangeDetailResponse, err error)
	DeleteLimitRange(ctx context.Context, req *DeleteLimitRangeRequest, opts ...http.CallOption) (rsp *DeleteLimitRangeResponse, err error)
}

type LimitRangeHTTPClientImpl struct {
	cc *http.Client
}

func NewLimitRangeHTTPClient(client *http.Client) LimitRangeHTTPClient {
	return &LimitRangeHTTPClientImpl{client}
}

func (c *LimitRangeHTTPClientImpl) ListLimitRange(ctx context.Context, in *ListLimitRangeRequest, opts ...http.CallOption) (*ListLimitRangeResponse, error) {
	var out ListLimitRangeResponse
	pattern := "/api/v1/limitrange/list"
	path := binding.EncodeURL(pattern, in, true)
	opts = append(opts, http.Operation(OperationLimitRangeListLimitRange))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "GET", path, nil, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *LimitRangeHTTPClientImpl) CreateLimitRange(ctx context.Context, in *CreateLimitRangeRequest, opts ...http.CallOption) (*CreateLimitRangeResponse, error) {
	var out CreateLimitRangeResponse
	pattern := "/api/v1/limitrange/create"
	path := binding.EncodeURL(pattern, in, false)
	opts = append(opts, http.Operation(OperationLimitRangeCreateLimitRange))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "POST", path, in, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *LimitRangeHTTPClientImpl) UpdateLimitRange(ctx context.Context, in *CreateLimitRangeRequest, opts ...http.CallOption) (*CreateLimitRangeResponse, error) {
	var out CreateLimitRangeResponse
	pattern := "/api/v1/limitrange/update"
	path := binding.EncodeURL(pattern, in, false)
	opts = append(opts, http.Operation(OperationLimitRangeUpdateLimitRange))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "POST", path, in, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *LimitRangeHTTPClientImpl) CreateOrUpdateLimitRange(ctx context.Context, in *CreateLimitRangeRequest, opts ...http.CallOption) (*CreateLimitRangeResponse, error) {
	var out CreateLimitRangeResponse
	pattern := "/api/v1/limitrange/create_or_update"
	path := binding.EncodeURL(pattern, in, false)
	opts = append(opts, http.Operation(OperationLimitRangeCreateOrUpdateLimitRange))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "POST", path, in, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *LimitRangeHTTPClientImpl) GetLimitRangeDetail(ctx context.Context, in *LimitRangeDetailRequest, opts ...http.CallOption) (*LimitRangeDetailResponse, error) {
	var out LimitRangeDetailResponse
	pattern := "/api/v1/limitrange/detail"
	path := binding.EncodeURL(pattern, in, true)
	opts = append(opts, http.Operation(OperationLimitRangeGetLimitRangeDetail))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "GET", path, nil, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *LimitRangeHTTPClientImpl) DeleteLimitRange(ctx context.Context, in *DeleteLimitRangeRequest, opts ...http.CallOption) (*DeleteLimitRangeResponse, error) {
	var out DeleteLimitRangeResponse
	pattern := "/api/v1/limitrange/delete"
	path := binding.EncodeURL(pattern, in, false)
	opts = append(opts, http.Operation(OperationLimitRangeDeleteLimitRange))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "POST", path, in, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, nil
}
