// Code generated by protoc-gen-go-http. DO NOT EDIT.
// versions:
// - protoc-gen-go-http v2.7.3
// - protoc             v5.27.2
// source: pb/configmap.v1.proto

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

const OperationConfigMapListConfigMap = "/configmap.v1.ConfigMap/ListConfigMap"
const OperationConfigMapCreateOrUpdateConfigMapByYaml = "/configmap.v1.ConfigMap/CreateOrUpdateConfigMapByYaml"
const OperationConfigMapDeleteConfigMap = "/configmap.v1.ConfigMap/DeleteConfigMap"
const OperationConfigMapGetConfigMapDetail = "/configmap.v1.ConfigMap/GetConfigMapDetail"
const OperationConfigMapCreateConfigMap = "/configmap.v1.ConfigMap/CreateConfigMap"
const OperationConfigMapUpdateConfigMap = "/configmap.v1.ConfigMap/UpdateConfigMap"

type ConfigMapHTTPServer interface {
	// ListConfigMap查看-云原生管理-ConfigMap-列表
	ListConfigMap(context.Context, *ListConfigMapsRequest) (*ListConfigMapsResponse, error)
	// CreateOrUpdateConfigMapByYaml管理-云原生管理-ConfigMap-Yaml创建更新
	CreateOrUpdateConfigMapByYaml(context.Context, *CreateOrUpdateConfigMapByYamlRequest) (*CreateOrUpdateConfigMapByYamlResponse, error)
	// DeleteConfigMap管理-云原生管理-ConfigMap-删除
	DeleteConfigMap(context.Context, *DeleteConfigMapRequest) (*DeleteConfigMapResponse, error)
	// GetConfigMapDetail查看-云原生管理-ConfigMap-详情
	GetConfigMapDetail(context.Context, *ConfigMapDetailRequest) (*ConfigMapDetailResponse, error)
	// CreateConfigMap管理-云原生管理-ConfigMap-创建
	CreateConfigMap(context.Context, *CreateConfigMapRequest) (*CreateConfigMapResponse, error)
	// UpdateConfigMap管理-云原生管理-ConfigMap-更新
	UpdateConfigMap(context.Context, *UpdateConfigMapRequest) (*UpdateConfigMapResponse, error)
}

func NewConfigMapHTTPServerMiddleware() middleware.Middleware {
	return selector.Server(
		selector.Server().Path(OperationConfigMapListConfigMap).Build(),
		selector.Server().Path(OperationConfigMapCreateOrUpdateConfigMapByYaml).Build(),
		selector.Server().Path(OperationConfigMapDeleteConfigMap).Build(),
		selector.Server().Path(OperationConfigMapGetConfigMapDetail).Build(),
		selector.Server().Path(OperationConfigMapCreateConfigMap).Build(),
		selector.Server().Path(OperationConfigMapUpdateConfigMap).Build(),
	).Path(
		OperationConfigMapListConfigMap,
		OperationConfigMapCreateOrUpdateConfigMapByYaml,
		OperationConfigMapDeleteConfigMap,
		OperationConfigMapGetConfigMapDetail,
		OperationConfigMapCreateConfigMap,
		OperationConfigMapUpdateConfigMap,
	).Build()
}

func RegisterConfigMapHTTPServer(s *http.Server, srv ConfigMapHTTPServer) {
	r := s.Route("/")
	r.GET("/api/v1/configmap/list", _ConfigMap_ListConfigMap0_HTTP_Handler(srv))
	r.POST("/api/v1/configmap/create_or_update_by_yaml", _ConfigMap_CreateOrUpdateConfigMapByYaml0_HTTP_Handler(srv))
	r.POST("/api/v1/configmap/delete", _ConfigMap_DeleteConfigMap0_HTTP_Handler(srv))
	r.GET("/api/v1/configmap/detail", _ConfigMap_GetConfigMapDetail0_HTTP_Handler(srv))
	r.POST("/api/v1/configmap/create", _ConfigMap_CreateConfigMap0_HTTP_Handler(srv))
	r.POST("/api/v1/configmap/update", _ConfigMap_UpdateConfigMap0_HTTP_Handler(srv))
}

func GenerateConfigMapHTTPServerRouteInfo() []route.Route {
	routes := make([]route.Route, 0, 6)
	routes = append(routes, route.Route{
		Method:  "GET",
		Path:    "/api/v1/configmap/list",
		Comment: "查看-云原生管理-ConfigMap-列表",
	})
	routes = append(routes, route.Route{
		Method:  "POST",
		Path:    "/api/v1/configmap/create_or_update_by_yaml",
		Comment: "管理-云原生管理-ConfigMap-Yaml创建更新",
	})
	routes = append(routes, route.Route{
		Method:  "POST",
		Path:    "/api/v1/configmap/delete",
		Comment: "管理-云原生管理-ConfigMap-删除",
	})
	routes = append(routes, route.Route{
		Method:  "GET",
		Path:    "/api/v1/configmap/detail",
		Comment: "查看-云原生管理-ConfigMap-详情",
	})
	routes = append(routes, route.Route{
		Method:  "POST",
		Path:    "/api/v1/configmap/create",
		Comment: "管理-云原生管理-ConfigMap-创建",
	})
	routes = append(routes, route.Route{
		Method:  "POST",
		Path:    "/api/v1/configmap/update",
		Comment: "管理-云原生管理-ConfigMap-更新",
	})
	return routes
}

func _ConfigMap_ListConfigMap0_HTTP_Handler(srv ConfigMapHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		stdCtx := kcontext.SetKHTTPContextWithContext(ctx, ctx)
		var in ListConfigMapsRequest
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationConfigMapListConfigMap)

		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.ListConfigMap(ctx, req.(*ListConfigMapsRequest))
		})
		out, err := h(stdCtx, &in)
		if err != nil {
			return err
		}
		reply := out.(*ListConfigMapsResponse)
		return ctx.Result(200, reply)
	}
}

func _ConfigMap_CreateOrUpdateConfigMapByYaml0_HTTP_Handler(srv ConfigMapHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		stdCtx := kcontext.SetKHTTPContextWithContext(ctx, ctx)
		var in CreateOrUpdateConfigMapByYamlRequest
		if err := ctx.Bind(&in); err != nil {
			return err
		}
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationConfigMapCreateOrUpdateConfigMapByYaml)
		auditRule := audit.NewAudit(
			"configmap",
			"Yaml创建更新",
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
						Const: "configmap",
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
			return srv.CreateOrUpdateConfigMapByYaml(ctx, req.(*CreateOrUpdateConfigMapByYamlRequest))
		})
		out, err := h(stdCtx, &in)
		if err != nil {
			return err
		}
		reply := out.(*CreateOrUpdateConfigMapByYamlResponse)
		return ctx.Result(200, reply)
	}
}

func _ConfigMap_DeleteConfigMap0_HTTP_Handler(srv ConfigMapHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		stdCtx := kcontext.SetKHTTPContextWithContext(ctx, ctx)
		var in DeleteConfigMapRequest
		if err := ctx.Bind(&in); err != nil {
			return err
		}
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationConfigMapDeleteConfigMap)
		auditRule := audit.NewAudit(
			"configmap",
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
						Const: "configmap",
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
			return srv.DeleteConfigMap(ctx, req.(*DeleteConfigMapRequest))
		})
		out, err := h(stdCtx, &in)
		if err != nil {
			return err
		}
		reply := out.(*DeleteConfigMapResponse)
		return ctx.Result(200, reply)
	}
}

func _ConfigMap_GetConfigMapDetail0_HTTP_Handler(srv ConfigMapHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		stdCtx := kcontext.SetKHTTPContextWithContext(ctx, ctx)
		var in ConfigMapDetailRequest
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationConfigMapGetConfigMapDetail)

		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.GetConfigMapDetail(ctx, req.(*ConfigMapDetailRequest))
		})
		out, err := h(stdCtx, &in)
		if err != nil {
			return err
		}
		reply := out.(*ConfigMapDetailResponse)
		return ctx.Result(200, reply)
	}
}

func _ConfigMap_CreateConfigMap0_HTTP_Handler(srv ConfigMapHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		stdCtx := kcontext.SetKHTTPContextWithContext(ctx, ctx)
		var in CreateConfigMapRequest
		if err := ctx.Bind(&in); err != nil {
			return err
		}
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationConfigMapCreateConfigMap)
		auditRule := audit.NewAudit(
			"configmap",
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
						Const: "configmap",
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
			return srv.CreateConfigMap(ctx, req.(*CreateConfigMapRequest))
		})
		out, err := h(stdCtx, &in)
		if err != nil {
			return err
		}
		reply := out.(*CreateConfigMapResponse)
		return ctx.Result(200, reply)
	}
}

func _ConfigMap_UpdateConfigMap0_HTTP_Handler(srv ConfigMapHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		stdCtx := kcontext.SetKHTTPContextWithContext(ctx, ctx)
		var in UpdateConfigMapRequest
		if err := ctx.Bind(&in); err != nil {
			return err
		}
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationConfigMapUpdateConfigMap)
		auditRule := audit.NewAudit(
			"configmap",
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
						Const: "configmap",
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
			return srv.UpdateConfigMap(ctx, req.(*UpdateConfigMapRequest))
		})
		out, err := h(stdCtx, &in)
		if err != nil {
			return err
		}
		reply := out.(*UpdateConfigMapResponse)
		return ctx.Result(200, reply)
	}
}

type ConfigMapHTTPClient interface {
	ListConfigMap(ctx context.Context, req *ListConfigMapsRequest, opts ...http.CallOption) (rsp *ListConfigMapsResponse, err error)
	CreateOrUpdateConfigMapByYaml(ctx context.Context, req *CreateOrUpdateConfigMapByYamlRequest, opts ...http.CallOption) (rsp *CreateOrUpdateConfigMapByYamlResponse, err error)
	DeleteConfigMap(ctx context.Context, req *DeleteConfigMapRequest, opts ...http.CallOption) (rsp *DeleteConfigMapResponse, err error)
	GetConfigMapDetail(ctx context.Context, req *ConfigMapDetailRequest, opts ...http.CallOption) (rsp *ConfigMapDetailResponse, err error)
	CreateConfigMap(ctx context.Context, req *CreateConfigMapRequest, opts ...http.CallOption) (rsp *CreateConfigMapResponse, err error)
	UpdateConfigMap(ctx context.Context, req *UpdateConfigMapRequest, opts ...http.CallOption) (rsp *UpdateConfigMapResponse, err error)
}

type ConfigMapHTTPClientImpl struct {
	cc *http.Client
}

func NewConfigMapHTTPClient(client *http.Client) ConfigMapHTTPClient {
	return &ConfigMapHTTPClientImpl{client}
}

func (c *ConfigMapHTTPClientImpl) ListConfigMap(ctx context.Context, in *ListConfigMapsRequest, opts ...http.CallOption) (*ListConfigMapsResponse, error) {
	var out ListConfigMapsResponse
	pattern := "/api/v1/configmap/list"
	path := binding.EncodeURL(pattern, in, true)
	opts = append(opts, http.Operation(OperationConfigMapListConfigMap))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "GET", path, nil, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *ConfigMapHTTPClientImpl) CreateOrUpdateConfigMapByYaml(ctx context.Context, in *CreateOrUpdateConfigMapByYamlRequest, opts ...http.CallOption) (*CreateOrUpdateConfigMapByYamlResponse, error) {
	var out CreateOrUpdateConfigMapByYamlResponse
	pattern := "/api/v1/configmap/create_or_update_by_yaml"
	path := binding.EncodeURL(pattern, in, false)
	opts = append(opts, http.Operation(OperationConfigMapCreateOrUpdateConfigMapByYaml))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "POST", path, in, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *ConfigMapHTTPClientImpl) DeleteConfigMap(ctx context.Context, in *DeleteConfigMapRequest, opts ...http.CallOption) (*DeleteConfigMapResponse, error) {
	var out DeleteConfigMapResponse
	pattern := "/api/v1/configmap/delete"
	path := binding.EncodeURL(pattern, in, false)
	opts = append(opts, http.Operation(OperationConfigMapDeleteConfigMap))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "POST", path, in, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *ConfigMapHTTPClientImpl) GetConfigMapDetail(ctx context.Context, in *ConfigMapDetailRequest, opts ...http.CallOption) (*ConfigMapDetailResponse, error) {
	var out ConfigMapDetailResponse
	pattern := "/api/v1/configmap/detail"
	path := binding.EncodeURL(pattern, in, true)
	opts = append(opts, http.Operation(OperationConfigMapGetConfigMapDetail))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "GET", path, nil, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *ConfigMapHTTPClientImpl) CreateConfigMap(ctx context.Context, in *CreateConfigMapRequest, opts ...http.CallOption) (*CreateConfigMapResponse, error) {
	var out CreateConfigMapResponse
	pattern := "/api/v1/configmap/create"
	path := binding.EncodeURL(pattern, in, false)
	opts = append(opts, http.Operation(OperationConfigMapCreateConfigMap))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "POST", path, in, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *ConfigMapHTTPClientImpl) UpdateConfigMap(ctx context.Context, in *UpdateConfigMapRequest, opts ...http.CallOption) (*UpdateConfigMapResponse, error) {
	var out UpdateConfigMapResponse
	pattern := "/api/v1/configmap/update"
	path := binding.EncodeURL(pattern, in, false)
	opts = append(opts, http.Operation(OperationConfigMapUpdateConfigMap))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "POST", path, in, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, nil
}
