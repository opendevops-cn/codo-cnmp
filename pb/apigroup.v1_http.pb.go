// Code generated by protoc-gen-go-http. DO NOT EDIT.
// versions:
// - protoc-gen-go-http v2.7.3
// - protoc             v5.27.2
// source: pb/apigroup.v1.proto

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

const OperationAPIGroupListAPIGroup = "/apigroup.v1.APIGroup/ListAPIGroup"

type APIGroupHTTPServer interface {
	// ListAPIGroup查看-云原生管理-APIGroup-列表
	ListAPIGroup(context.Context, *ListAPIGroupRequest) (*ListAPIGroupResponse, error)
}

func NewAPIGroupHTTPServerMiddleware() middleware.Middleware {
	return selector.Server(
		selector.Server().Path(OperationAPIGroupListAPIGroup).Build(),
	).Path(
		OperationAPIGroupListAPIGroup,
	).Build()
}

func RegisterAPIGroupHTTPServer(s *http.Server, srv APIGroupHTTPServer) {
	r := s.Route("/")
	r.GET("/api/v1/apigroup/list", _APIGroup_ListAPIGroup0_HTTP_Handler(srv))
}

func GenerateAPIGroupHTTPServerRouteInfo() []route.Route {
	routes := make([]route.Route, 0, 1)
	routes = append(routes, route.Route{
		Method:  "GET",
		Path:    "/api/v1/apigroup/list",
		Comment: "查看-云原生管理-APIGroup-列表",
	})
	return routes
}

func _APIGroup_ListAPIGroup0_HTTP_Handler(srv APIGroupHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		stdCtx := kcontext.SetKHTTPContextWithContext(ctx, ctx)
		var in ListAPIGroupRequest
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationAPIGroupListAPIGroup)

		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.ListAPIGroup(ctx, req.(*ListAPIGroupRequest))
		})
		out, err := h(stdCtx, &in)
		if err != nil {
			return err
		}
		reply := out.(*ListAPIGroupResponse)
		return ctx.Result(200, reply)
	}
}

type APIGroupHTTPClient interface {
	ListAPIGroup(ctx context.Context, req *ListAPIGroupRequest, opts ...http.CallOption) (rsp *ListAPIGroupResponse, err error)
}

type APIGroupHTTPClientImpl struct {
	cc *http.Client
}

func NewAPIGroupHTTPClient(client *http.Client) APIGroupHTTPClient {
	return &APIGroupHTTPClientImpl{client}
}

func (c *APIGroupHTTPClientImpl) ListAPIGroup(ctx context.Context, in *ListAPIGroupRequest, opts ...http.CallOption) (*ListAPIGroupResponse, error) {
	var out ListAPIGroupResponse
	pattern := "/api/v1/apigroup/list"
	path := binding.EncodeURL(pattern, in, true)
	opts = append(opts, http.Operation(OperationAPIGroupListAPIGroup))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "GET", path, nil, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, nil
}
