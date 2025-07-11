// Code generated by protoc-gen-go-http. DO NOT EDIT.
// versions:
// - protoc-gen-go-http v2.7.3
// - protoc             v5.27.2
// source: pb/svc.v1.proto

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

const OperationSVCListSvc = "/svc.v1.SVC/ListSvc"
const OperationSVCCreateSvc = "/svc.v1.SVC/CreateSvc"
const OperationSVCUpdateSvc = "/svc.v1.SVC/UpdateSvc"
const OperationSVCDeleteSvc = "/svc.v1.SVC/DeleteSvc"
const OperationSVCGetSvcDetail = "/svc.v1.SVC/GetSvcDetail"

type SVCHTTPServer interface {
	// ListSvc查看-云原生管理-Service-列表
	ListSvc(context.Context, *ListSvcRequest) (*ListSvcResponse, error)
	// CreateSvc管理-云原生管理-Service-创建
	CreateSvc(context.Context, *CreateSvcRequest) (*CreateSvcResponse, error)
	// UpdateSvc管理-云原生管理-Service-编辑
	UpdateSvc(context.Context, *CreateSvcRequest) (*CreateSvcResponse, error)
	// DeleteSvc管理-云原生管理-Service-删除
	DeleteSvc(context.Context, *DeleteSvcRequest) (*DeleteSvcResponse, error)
	// GetSvcDetail查看-云原生管理-Service-详情
	GetSvcDetail(context.Context, *SvcDetailRequest) (*SvcDetailResponse, error)
}

func NewSVCHTTPServerMiddleware() middleware.Middleware {
	return selector.Server(
		selector.Server().Path(OperationSVCListSvc).Build(),
		selector.Server().Path(OperationSVCCreateSvc).Build(),
		selector.Server().Path(OperationSVCUpdateSvc).Build(),
		selector.Server().Path(OperationSVCDeleteSvc).Build(),
		selector.Server().Path(OperationSVCGetSvcDetail).Build(),
	).Path(
		OperationSVCListSvc,
		OperationSVCCreateSvc,
		OperationSVCUpdateSvc,
		OperationSVCDeleteSvc,
		OperationSVCGetSvcDetail,
	).Build()
}

func RegisterSVCHTTPServer(s *http.Server, srv SVCHTTPServer) {
	r := s.Route("/")
	r.GET("/api/v1/svc/list", _SVC_ListSvc0_HTTP_Handler(srv))
	r.POST("/api/v1/svc/create", _SVC_CreateSvc0_HTTP_Handler(srv))
	r.POST("/api/v1/svc/update", _SVC_UpdateSvc0_HTTP_Handler(srv))
	r.POST("/api/v1/svc/delete", _SVC_DeleteSvc0_HTTP_Handler(srv))
	r.GET("/api/v1/svc/detail", _SVC_GetSvcDetail0_HTTP_Handler(srv))
}

func GenerateSVCHTTPServerRouteInfo() []route.Route {
	routes := make([]route.Route, 0, 5)
	routes = append(routes, route.Route{
		Method:  "GET",
		Path:    "/api/v1/svc/list",
		Comment: "查看-云原生管理-Service-列表",
	})
	routes = append(routes, route.Route{
		Method:  "POST",
		Path:    "/api/v1/svc/create",
		Comment: "管理-云原生管理-Service-创建",
	})
	routes = append(routes, route.Route{
		Method:  "POST",
		Path:    "/api/v1/svc/update",
		Comment: "管理-云原生管理-Service-编辑",
	})
	routes = append(routes, route.Route{
		Method:  "POST",
		Path:    "/api/v1/svc/delete",
		Comment: "管理-云原生管理-Service-删除",
	})
	routes = append(routes, route.Route{
		Method:  "GET",
		Path:    "/api/v1/svc/detail",
		Comment: "查看-云原生管理-Service-详情",
	})
	return routes
}

func _SVC_ListSvc0_HTTP_Handler(srv SVCHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		stdCtx := kcontext.SetKHTTPContextWithContext(ctx, ctx)
		var in ListSvcRequest
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationSVCListSvc)

		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.ListSvc(ctx, req.(*ListSvcRequest))
		})
		out, err := h(stdCtx, &in)
		if err != nil {
			return err
		}
		reply := out.(*ListSvcResponse)
		return ctx.Result(200, reply)
	}
}

func _SVC_CreateSvc0_HTTP_Handler(srv SVCHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		stdCtx := kcontext.SetKHTTPContextWithContext(ctx, ctx)
		var in CreateSvcRequest
		if err := ctx.Bind(&in); err != nil {
			return err
		}
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationSVCCreateSvc)
		auditRule := audit.NewAudit(
			"Service",
			"创建svc",
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
						Const: "service",
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
			return srv.CreateSvc(ctx, req.(*CreateSvcRequest))
		})
		out, err := h(stdCtx, &in)
		if err != nil {
			return err
		}
		reply := out.(*CreateSvcResponse)
		return ctx.Result(200, reply)
	}
}

func _SVC_UpdateSvc0_HTTP_Handler(srv SVCHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		stdCtx := kcontext.SetKHTTPContextWithContext(ctx, ctx)
		var in CreateSvcRequest
		if err := ctx.Bind(&in); err != nil {
			return err
		}
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationSVCUpdateSvc)
		auditRule := audit.NewAudit(
			"Service",
			"编辑svc",
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
						Const: "service",
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
			return srv.UpdateSvc(ctx, req.(*CreateSvcRequest))
		})
		out, err := h(stdCtx, &in)
		if err != nil {
			return err
		}
		reply := out.(*CreateSvcResponse)
		return ctx.Result(200, reply)
	}
}

func _SVC_DeleteSvc0_HTTP_Handler(srv SVCHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		stdCtx := kcontext.SetKHTTPContextWithContext(ctx, ctx)
		var in DeleteSvcRequest
		if err := ctx.Bind(&in); err != nil {
			return err
		}
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationSVCDeleteSvc)
		auditRule := audit.NewAudit(
			"Service",
			"删除svc",
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
						Const: "service",
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
			return srv.DeleteSvc(ctx, req.(*DeleteSvcRequest))
		})
		out, err := h(stdCtx, &in)
		if err != nil {
			return err
		}
		reply := out.(*DeleteSvcResponse)
		return ctx.Result(200, reply)
	}
}

func _SVC_GetSvcDetail0_HTTP_Handler(srv SVCHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		stdCtx := kcontext.SetKHTTPContextWithContext(ctx, ctx)
		var in SvcDetailRequest
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationSVCGetSvcDetail)

		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.GetSvcDetail(ctx, req.(*SvcDetailRequest))
		})
		out, err := h(stdCtx, &in)
		if err != nil {
			return err
		}
		reply := out.(*SvcDetailResponse)
		return ctx.Result(200, reply)
	}
}

type SVCHTTPClient interface {
	ListSvc(ctx context.Context, req *ListSvcRequest, opts ...http.CallOption) (rsp *ListSvcResponse, err error)
	CreateSvc(ctx context.Context, req *CreateSvcRequest, opts ...http.CallOption) (rsp *CreateSvcResponse, err error)
	UpdateSvc(ctx context.Context, req *CreateSvcRequest, opts ...http.CallOption) (rsp *CreateSvcResponse, err error)
	DeleteSvc(ctx context.Context, req *DeleteSvcRequest, opts ...http.CallOption) (rsp *DeleteSvcResponse, err error)
	GetSvcDetail(ctx context.Context, req *SvcDetailRequest, opts ...http.CallOption) (rsp *SvcDetailResponse, err error)
}

type SVCHTTPClientImpl struct {
	cc *http.Client
}

func NewSVCHTTPClient(client *http.Client) SVCHTTPClient {
	return &SVCHTTPClientImpl{client}
}

func (c *SVCHTTPClientImpl) ListSvc(ctx context.Context, in *ListSvcRequest, opts ...http.CallOption) (*ListSvcResponse, error) {
	var out ListSvcResponse
	pattern := "/api/v1/svc/list"
	path := binding.EncodeURL(pattern, in, true)
	opts = append(opts, http.Operation(OperationSVCListSvc))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "GET", path, nil, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *SVCHTTPClientImpl) CreateSvc(ctx context.Context, in *CreateSvcRequest, opts ...http.CallOption) (*CreateSvcResponse, error) {
	var out CreateSvcResponse
	pattern := "/api/v1/svc/create"
	path := binding.EncodeURL(pattern, in, false)
	opts = append(opts, http.Operation(OperationSVCCreateSvc))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "POST", path, in, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *SVCHTTPClientImpl) UpdateSvc(ctx context.Context, in *CreateSvcRequest, opts ...http.CallOption) (*CreateSvcResponse, error) {
	var out CreateSvcResponse
	pattern := "/api/v1/svc/update"
	path := binding.EncodeURL(pattern, in, false)
	opts = append(opts, http.Operation(OperationSVCUpdateSvc))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "POST", path, in, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *SVCHTTPClientImpl) DeleteSvc(ctx context.Context, in *DeleteSvcRequest, opts ...http.CallOption) (*DeleteSvcResponse, error) {
	var out DeleteSvcResponse
	pattern := "/api/v1/svc/delete"
	path := binding.EncodeURL(pattern, in, false)
	opts = append(opts, http.Operation(OperationSVCDeleteSvc))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "POST", path, in, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *SVCHTTPClientImpl) GetSvcDetail(ctx context.Context, in *SvcDetailRequest, opts ...http.CallOption) (*SvcDetailResponse, error) {
	var out SvcDetailResponse
	pattern := "/api/v1/svc/detail"
	path := binding.EncodeURL(pattern, in, true)
	opts = append(opts, http.Operation(OperationSVCGetSvcDetail))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "GET", path, nil, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, nil
}
