// Code generated by protoc-gen-go-http. DO NOT EDIT.
// versions:
// - protoc-gen-go-http v2.7.3
// - protoc             v5.27.2
// source: pb/ingress.v1.proto

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

const OperationIngressListIngress = "/ingress.v1.Ingress/ListIngress"
const OperationIngressListIngressHost = "/ingress.v1.Ingress/ListIngressHost"
const OperationIngressCreateIngress = "/ingress.v1.Ingress/CreateIngress"
const OperationIngressUpdateIngress = "/ingress.v1.Ingress/UpdateIngress"
const OperationIngressDeleteIngress = "/ingress.v1.Ingress/DeleteIngress"
const OperationIngressGetIngressDetail = "/ingress.v1.Ingress/GetIngressDetail"

type IngressHTTPServer interface {
	// ListIngress查看-云原生管理-Ingress-列表
	ListIngress(context.Context, *ListIngressRequest) (*ListIngressResponse, error)
	// ListIngressHost查看-云原生管理-Ingress域名-列表
	ListIngressHost(context.Context, *ListHostRequest) (*ListHostResponse, error)
	// CreateIngress管理-云原生管理-Ingress-创建
	CreateIngress(context.Context, *CreateIngressRequest) (*CreateIngressResponse, error)
	// UpdateIngress管理-云原生管理-Ingress-编辑
	UpdateIngress(context.Context, *CreateIngressRequest) (*CreateIngressResponse, error)
	// DeleteIngress管理-云原生管理-Ingress-删除
	DeleteIngress(context.Context, *DeleteIngressRequest) (*DeleteIngressResponse, error)
	// GetIngressDetail查看-云原生管理-Ingress-详情
	GetIngressDetail(context.Context, *IngressDetailRequest) (*IngressDetailResponse, error)
}

func NewIngressHTTPServerMiddleware() middleware.Middleware {
	return selector.Server(
		selector.Server().Path(OperationIngressListIngress).Build(),
		selector.Server().Path(OperationIngressListIngressHost).Build(),
		selector.Server().Path(OperationIngressCreateIngress).Build(),
		selector.Server().Path(OperationIngressUpdateIngress).Build(),
		selector.Server().Path(OperationIngressDeleteIngress).Build(),
		selector.Server().Path(OperationIngressGetIngressDetail).Build(),
	).Path(
		OperationIngressListIngress,
		OperationIngressListIngressHost,
		OperationIngressCreateIngress,
		OperationIngressUpdateIngress,
		OperationIngressDeleteIngress,
		OperationIngressGetIngressDetail,
	).Build()
}

func RegisterIngressHTTPServer(s *http.Server, srv IngressHTTPServer) {
	r := s.Route("/")
	r.GET("/api/v1/ingress/list", _Ingress_ListIngress0_HTTP_Handler(srv))
	r.GET("/api/v1/ingress/host/list", _Ingress_ListIngressHost0_HTTP_Handler(srv))
	r.POST("/api/v1/ingress/create", _Ingress_CreateIngress0_HTTP_Handler(srv))
	r.POST("/api/v1/ingress/update", _Ingress_UpdateIngress0_HTTP_Handler(srv))
	r.POST("/api/v1/ingress/delete", _Ingress_DeleteIngress0_HTTP_Handler(srv))
	r.GET("/api/v1/ingress/detail", _Ingress_GetIngressDetail0_HTTP_Handler(srv))
}

func GenerateIngressHTTPServerRouteInfo() []route.Route {
	routes := make([]route.Route, 0, 6)
	routes = append(routes, route.Route{
		Method:  "GET",
		Path:    "/api/v1/ingress/list",
		Comment: "查看-云原生管理-Ingress-列表",
	})
	routes = append(routes, route.Route{
		Method:  "GET",
		Path:    "/api/v1/ingress/host/list",
		Comment: "查看-云原生管理-Ingress域名-列表",
	})
	routes = append(routes, route.Route{
		Method:  "POST",
		Path:    "/api/v1/ingress/create",
		Comment: "管理-云原生管理-Ingress-创建",
	})
	routes = append(routes, route.Route{
		Method:  "POST",
		Path:    "/api/v1/ingress/update",
		Comment: "管理-云原生管理-Ingress-编辑",
	})
	routes = append(routes, route.Route{
		Method:  "POST",
		Path:    "/api/v1/ingress/delete",
		Comment: "管理-云原生管理-Ingress-删除",
	})
	routes = append(routes, route.Route{
		Method:  "GET",
		Path:    "/api/v1/ingress/detail",
		Comment: "查看-云原生管理-Ingress-详情",
	})
	return routes
}

func _Ingress_ListIngress0_HTTP_Handler(srv IngressHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		stdCtx := kcontext.SetKHTTPContextWithContext(ctx, ctx)
		var in ListIngressRequest
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationIngressListIngress)

		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.ListIngress(ctx, req.(*ListIngressRequest))
		})
		out, err := h(stdCtx, &in)
		if err != nil {
			return err
		}
		reply := out.(*ListIngressResponse)
		return ctx.Result(200, reply)
	}
}

func _Ingress_ListIngressHost0_HTTP_Handler(srv IngressHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		stdCtx := kcontext.SetKHTTPContextWithContext(ctx, ctx)
		var in ListHostRequest
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationIngressListIngressHost)

		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.ListIngressHost(ctx, req.(*ListHostRequest))
		})
		out, err := h(stdCtx, &in)
		if err != nil {
			return err
		}
		reply := out.(*ListHostResponse)
		return ctx.Result(200, reply)
	}
}

func _Ingress_CreateIngress0_HTTP_Handler(srv IngressHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		stdCtx := kcontext.SetKHTTPContextWithContext(ctx, ctx)
		var in CreateIngressRequest
		if err := ctx.Bind(&in); err != nil {
			return err
		}
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationIngressCreateIngress)
		auditRule := audit.NewAudit(
			"Service",
			"创建ingress",
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
						Const: "Service",
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
			return srv.CreateIngress(ctx, req.(*CreateIngressRequest))
		})
		out, err := h(stdCtx, &in)
		if err != nil {
			return err
		}
		reply := out.(*CreateIngressResponse)
		return ctx.Result(200, reply)
	}
}

func _Ingress_UpdateIngress0_HTTP_Handler(srv IngressHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		stdCtx := kcontext.SetKHTTPContextWithContext(ctx, ctx)
		var in CreateIngressRequest
		if err := ctx.Bind(&in); err != nil {
			return err
		}
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationIngressUpdateIngress)
		auditRule := audit.NewAudit(
			"Service",
			"编辑ingress",
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
						Const: "Service",
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
			return srv.UpdateIngress(ctx, req.(*CreateIngressRequest))
		})
		out, err := h(stdCtx, &in)
		if err != nil {
			return err
		}
		reply := out.(*CreateIngressResponse)
		return ctx.Result(200, reply)
	}
}

func _Ingress_DeleteIngress0_HTTP_Handler(srv IngressHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		stdCtx := kcontext.SetKHTTPContextWithContext(ctx, ctx)
		var in DeleteIngressRequest
		if err := ctx.Bind(&in); err != nil {
			return err
		}
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationIngressDeleteIngress)
		auditRule := audit.NewAudit(
			"Service",
			"删除ingress",
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
						Const: "Service",
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
			return srv.DeleteIngress(ctx, req.(*DeleteIngressRequest))
		})
		out, err := h(stdCtx, &in)
		if err != nil {
			return err
		}
		reply := out.(*DeleteIngressResponse)
		return ctx.Result(200, reply)
	}
}

func _Ingress_GetIngressDetail0_HTTP_Handler(srv IngressHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		stdCtx := kcontext.SetKHTTPContextWithContext(ctx, ctx)
		var in IngressDetailRequest
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationIngressGetIngressDetail)

		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.GetIngressDetail(ctx, req.(*IngressDetailRequest))
		})
		out, err := h(stdCtx, &in)
		if err != nil {
			return err
		}
		reply := out.(*IngressDetailResponse)
		return ctx.Result(200, reply)
	}
}

type IngressHTTPClient interface {
	ListIngress(ctx context.Context, req *ListIngressRequest, opts ...http.CallOption) (rsp *ListIngressResponse, err error)
	ListIngressHost(ctx context.Context, req *ListHostRequest, opts ...http.CallOption) (rsp *ListHostResponse, err error)
	CreateIngress(ctx context.Context, req *CreateIngressRequest, opts ...http.CallOption) (rsp *CreateIngressResponse, err error)
	UpdateIngress(ctx context.Context, req *CreateIngressRequest, opts ...http.CallOption) (rsp *CreateIngressResponse, err error)
	DeleteIngress(ctx context.Context, req *DeleteIngressRequest, opts ...http.CallOption) (rsp *DeleteIngressResponse, err error)
	GetIngressDetail(ctx context.Context, req *IngressDetailRequest, opts ...http.CallOption) (rsp *IngressDetailResponse, err error)
}

type IngressHTTPClientImpl struct {
	cc *http.Client
}

func NewIngressHTTPClient(client *http.Client) IngressHTTPClient {
	return &IngressHTTPClientImpl{client}
}

func (c *IngressHTTPClientImpl) ListIngress(ctx context.Context, in *ListIngressRequest, opts ...http.CallOption) (*ListIngressResponse, error) {
	var out ListIngressResponse
	pattern := "/api/v1/ingress/list"
	path := binding.EncodeURL(pattern, in, true)
	opts = append(opts, http.Operation(OperationIngressListIngress))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "GET", path, nil, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *IngressHTTPClientImpl) ListIngressHost(ctx context.Context, in *ListHostRequest, opts ...http.CallOption) (*ListHostResponse, error) {
	var out ListHostResponse
	pattern := "/api/v1/ingress/host/list"
	path := binding.EncodeURL(pattern, in, true)
	opts = append(opts, http.Operation(OperationIngressListIngressHost))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "GET", path, nil, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *IngressHTTPClientImpl) CreateIngress(ctx context.Context, in *CreateIngressRequest, opts ...http.CallOption) (*CreateIngressResponse, error) {
	var out CreateIngressResponse
	pattern := "/api/v1/ingress/create"
	path := binding.EncodeURL(pattern, in, false)
	opts = append(opts, http.Operation(OperationIngressCreateIngress))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "POST", path, in, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *IngressHTTPClientImpl) UpdateIngress(ctx context.Context, in *CreateIngressRequest, opts ...http.CallOption) (*CreateIngressResponse, error) {
	var out CreateIngressResponse
	pattern := "/api/v1/ingress/update"
	path := binding.EncodeURL(pattern, in, false)
	opts = append(opts, http.Operation(OperationIngressUpdateIngress))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "POST", path, in, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *IngressHTTPClientImpl) DeleteIngress(ctx context.Context, in *DeleteIngressRequest, opts ...http.CallOption) (*DeleteIngressResponse, error) {
	var out DeleteIngressResponse
	pattern := "/api/v1/ingress/delete"
	path := binding.EncodeURL(pattern, in, false)
	opts = append(opts, http.Operation(OperationIngressDeleteIngress))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "POST", path, in, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *IngressHTTPClientImpl) GetIngressDetail(ctx context.Context, in *IngressDetailRequest, opts ...http.CallOption) (*IngressDetailResponse, error) {
	var out IngressDetailResponse
	pattern := "/api/v1/ingress/detail"
	path := binding.EncodeURL(pattern, in, true)
	opts = append(opts, http.Operation(OperationIngressGetIngressDetail))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "GET", path, nil, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, nil
}
