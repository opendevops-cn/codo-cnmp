// Code generated by protoc-gen-go-http. DO NOT EDIT.
// versions:
// - protoc-gen-go-http v2.7.3
// - protoc             v5.27.2
// source: pb/daemonset.v1.proto

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

const OperationDaemonSetListDaemonSet = "/daemonset.DaemonSet/ListDaemonSet"
const OperationDaemonSetCreateOrUpdateDaemonSetByYaml = "/daemonset.DaemonSet/CreateOrUpdateDaemonSetByYaml"
const OperationDaemonSetDeleteDaemonSet = "/daemonset.DaemonSet/DeleteDaemonSet"
const OperationDaemonSetRestartDaemonSet = "/daemonset.DaemonSet/RestartDaemonSet"
const OperationDaemonSetGetDaemonSetDetail = "/daemonset.DaemonSet/GetDaemonSetDetail"
const OperationDaemonSetGetDaemonSetRevisions = "/daemonset.DaemonSet/GetDaemonSetRevisions"
const OperationDaemonSetRollbackDaemonSet = "/daemonset.DaemonSet/RollbackDaemonSet"
const OperationDaemonSetUpdateStatefulSetUpdateStrategy = "/daemonset.DaemonSet/UpdateStatefulSetUpdateStrategy"

type DaemonSetHTTPServer interface {
	// ListDaemonSet查看-云原生管理-DaemonSet-列表
	ListDaemonSet(context.Context, *ListDaemonSetRequest) (*ListDaemonSetResponse, error)
	// CreateOrUpdateDaemonSetByYaml管理-云原生管理-DaemonSet-Yaml创建更新
	CreateOrUpdateDaemonSetByYaml(context.Context, *CreateOrUpdateDaemonSetByYamlRequest) (*CreateOrUpdateDaemonSetByYamlResponse, error)
	// DeleteDaemonSet管理-云原生管理-DaemonSet-删除
	DeleteDaemonSet(context.Context, *DeleteDaemonSetRequest) (*DeleteDaemonSetResponse, error)
	// RestartDaemonSet管理-云原生管理-DaemonSet-重启
	RestartDaemonSet(context.Context, *RestartDaemonSetRequest) (*RestartDaemonSetResponse, error)
	// GetDaemonSetDetail查看-云原生管理-DaemonSet-详情
	GetDaemonSetDetail(context.Context, *GetDaemonSetDetailRequest) (*GetDaemonSetDetailResponse, error)
	// GetDaemonSetRevisions查看-云原生管理-DaemonSet-历史版本
	GetDaemonSetRevisions(context.Context, *GetDaemonSetHistoryRequest) (*GetDaemonSetHistoryResponse, error)
	// RollbackDaemonSet管理-云原生管理-DaemonSet-回滚
	RollbackDaemonSet(context.Context, *RollbackDaemonSetRequest) (*RollbackDaemonSetResponse, error)
	// UpdateStatefulSetUpdateStrategy管理-云原生管理-DaemonSet-更新策略
	UpdateStatefulSetUpdateStrategy(context.Context, *UpdateDaemonSetUpdateStrategyRequest) (*UpdateDaemonSetUpdateStrategyResponse, error)
}

func NewDaemonSetHTTPServerMiddleware() middleware.Middleware {
	return selector.Server(
		selector.Server().Path(OperationDaemonSetListDaemonSet).Build(),
		selector.Server().Path(OperationDaemonSetCreateOrUpdateDaemonSetByYaml).Build(),
		selector.Server().Path(OperationDaemonSetDeleteDaemonSet).Build(),
		selector.Server().Path(OperationDaemonSetRestartDaemonSet).Build(),
		selector.Server().Path(OperationDaemonSetGetDaemonSetDetail).Build(),
		selector.Server().Path(OperationDaemonSetGetDaemonSetRevisions).Build(),
		selector.Server().Path(OperationDaemonSetRollbackDaemonSet).Build(),
		selector.Server().Path(OperationDaemonSetUpdateStatefulSetUpdateStrategy).Build(),
	).Path(
		OperationDaemonSetListDaemonSet,
		OperationDaemonSetCreateOrUpdateDaemonSetByYaml,
		OperationDaemonSetDeleteDaemonSet,
		OperationDaemonSetRestartDaemonSet,
		OperationDaemonSetGetDaemonSetDetail,
		OperationDaemonSetGetDaemonSetRevisions,
		OperationDaemonSetRollbackDaemonSet,
		OperationDaemonSetUpdateStatefulSetUpdateStrategy,
	).Build()
}

func RegisterDaemonSetHTTPServer(s *http.Server, srv DaemonSetHTTPServer) {
	r := s.Route("/")
	r.GET("/api/v1/daemonset/list", _DaemonSet_ListDaemonSet0_HTTP_Handler(srv))
	r.POST("/api/v1/daemonset/create_or_update_by_yaml", _DaemonSet_CreateOrUpdateDaemonSetByYaml0_HTTP_Handler(srv))
	r.POST("/api/v1/daemonset/delete", _DaemonSet_DeleteDaemonSet0_HTTP_Handler(srv))
	r.POST("/api/v1/daemonset/restart", _DaemonSet_RestartDaemonSet0_HTTP_Handler(srv))
	r.GET("/api/v1/daemonset/detail", _DaemonSet_GetDaemonSetDetail0_HTTP_Handler(srv))
	r.GET("/api/v1/daemonset/revisions", _DaemonSet_GetDaemonSetRevisions0_HTTP_Handler(srv))
	r.POST("/api/v1/daemonset/rollback", _DaemonSet_RollbackDaemonSet0_HTTP_Handler(srv))
	r.POST("/api/v1/daemonset/upgrade_strategy/update", _DaemonSet_UpdateStatefulSetUpdateStrategy0_HTTP_Handler(srv))
}

func GenerateDaemonSetHTTPServerRouteInfo() []route.Route {
	routes := make([]route.Route, 0, 8)
	routes = append(routes, route.Route{
		Method:  "GET",
		Path:    "/api/v1/daemonset/list",
		Comment: "查看-云原生管理-DaemonSet-列表",
	})
	routes = append(routes, route.Route{
		Method:  "POST",
		Path:    "/api/v1/daemonset/create_or_update_by_yaml",
		Comment: "管理-云原生管理-DaemonSet-Yaml创建更新",
	})
	routes = append(routes, route.Route{
		Method:  "POST",
		Path:    "/api/v1/daemonset/delete",
		Comment: "管理-云原生管理-DaemonSet-删除",
	})
	routes = append(routes, route.Route{
		Method:  "POST",
		Path:    "/api/v1/daemonset/restart",
		Comment: "管理-云原生管理-DaemonSet-重启",
	})
	routes = append(routes, route.Route{
		Method:  "GET",
		Path:    "/api/v1/daemonset/detail",
		Comment: "查看-云原生管理-DaemonSet-详情",
	})
	routes = append(routes, route.Route{
		Method:  "GET",
		Path:    "/api/v1/daemonset/revisions",
		Comment: "查看-云原生管理-DaemonSet-历史版本",
	})
	routes = append(routes, route.Route{
		Method:  "POST",
		Path:    "/api/v1/daemonset/rollback",
		Comment: "管理-云原生管理-DaemonSet-回滚",
	})
	routes = append(routes, route.Route{
		Method:  "POST",
		Path:    "/api/v1/daemonset/upgrade_strategy/update",
		Comment: "管理-云原生管理-DaemonSet-更新策略",
	})
	return routes
}

func _DaemonSet_ListDaemonSet0_HTTP_Handler(srv DaemonSetHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		stdCtx := kcontext.SetKHTTPContextWithContext(ctx, ctx)
		var in ListDaemonSetRequest
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationDaemonSetListDaemonSet)

		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.ListDaemonSet(ctx, req.(*ListDaemonSetRequest))
		})
		out, err := h(stdCtx, &in)
		if err != nil {
			return err
		}
		reply := out.(*ListDaemonSetResponse)
		return ctx.Result(200, reply)
	}
}

func _DaemonSet_CreateOrUpdateDaemonSetByYaml0_HTTP_Handler(srv DaemonSetHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		stdCtx := kcontext.SetKHTTPContextWithContext(ctx, ctx)
		var in CreateOrUpdateDaemonSetByYamlRequest
		if err := ctx.Bind(&in); err != nil {
			return err
		}
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationDaemonSetCreateOrUpdateDaemonSetByYaml)
		auditRule := audit.NewAudit(
			"daemonset",
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
						Const: "daemonset",
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
			return srv.CreateOrUpdateDaemonSetByYaml(ctx, req.(*CreateOrUpdateDaemonSetByYamlRequest))
		})
		out, err := h(stdCtx, &in)
		if err != nil {
			return err
		}
		reply := out.(*CreateOrUpdateDaemonSetByYamlResponse)
		return ctx.Result(200, reply)
	}
}

func _DaemonSet_DeleteDaemonSet0_HTTP_Handler(srv DaemonSetHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		stdCtx := kcontext.SetKHTTPContextWithContext(ctx, ctx)
		var in DeleteDaemonSetRequest
		if err := ctx.Bind(&in); err != nil {
			return err
		}
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationDaemonSetDeleteDaemonSet)
		auditRule := audit.NewAudit(
			"daemonset",
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
						Const: "daemonset",
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
			return srv.DeleteDaemonSet(ctx, req.(*DeleteDaemonSetRequest))
		})
		out, err := h(stdCtx, &in)
		if err != nil {
			return err
		}
		reply := out.(*DeleteDaemonSetResponse)
		return ctx.Result(200, reply)
	}
}

func _DaemonSet_RestartDaemonSet0_HTTP_Handler(srv DaemonSetHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		stdCtx := kcontext.SetKHTTPContextWithContext(ctx, ctx)
		var in RestartDaemonSetRequest
		if err := ctx.Bind(&in); err != nil {
			return err
		}
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationDaemonSetRestartDaemonSet)
		auditRule := audit.NewAudit(
			"daemonset",
			"重启",
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
						Const: "daemonset",
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
			return srv.RestartDaemonSet(ctx, req.(*RestartDaemonSetRequest))
		})
		out, err := h(stdCtx, &in)
		if err != nil {
			return err
		}
		reply := out.(*RestartDaemonSetResponse)
		return ctx.Result(200, reply)
	}
}

func _DaemonSet_GetDaemonSetDetail0_HTTP_Handler(srv DaemonSetHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		stdCtx := kcontext.SetKHTTPContextWithContext(ctx, ctx)
		var in GetDaemonSetDetailRequest
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationDaemonSetGetDaemonSetDetail)

		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.GetDaemonSetDetail(ctx, req.(*GetDaemonSetDetailRequest))
		})
		out, err := h(stdCtx, &in)
		if err != nil {
			return err
		}
		reply := out.(*GetDaemonSetDetailResponse)
		return ctx.Result(200, reply)
	}
}

func _DaemonSet_GetDaemonSetRevisions0_HTTP_Handler(srv DaemonSetHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		stdCtx := kcontext.SetKHTTPContextWithContext(ctx, ctx)
		var in GetDaemonSetHistoryRequest
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationDaemonSetGetDaemonSetRevisions)

		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.GetDaemonSetRevisions(ctx, req.(*GetDaemonSetHistoryRequest))
		})
		out, err := h(stdCtx, &in)
		if err != nil {
			return err
		}
		reply := out.(*GetDaemonSetHistoryResponse)
		return ctx.Result(200, reply)
	}
}

func _DaemonSet_RollbackDaemonSet0_HTTP_Handler(srv DaemonSetHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		stdCtx := kcontext.SetKHTTPContextWithContext(ctx, ctx)
		var in RollbackDaemonSetRequest
		if err := ctx.Bind(&in); err != nil {
			return err
		}
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationDaemonSetRollbackDaemonSet)
		auditRule := audit.NewAudit(
			"daemonset",
			"回滚",
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
						Const: "daemonset",
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
			return srv.RollbackDaemonSet(ctx, req.(*RollbackDaemonSetRequest))
		})
		out, err := h(stdCtx, &in)
		if err != nil {
			return err
		}
		reply := out.(*RollbackDaemonSetResponse)
		return ctx.Result(200, reply)
	}
}

func _DaemonSet_UpdateStatefulSetUpdateStrategy0_HTTP_Handler(srv DaemonSetHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		stdCtx := kcontext.SetKHTTPContextWithContext(ctx, ctx)
		var in UpdateDaemonSetUpdateStrategyRequest
		if err := ctx.Bind(&in); err != nil {
			return err
		}
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationDaemonSetUpdateStatefulSetUpdateStrategy)
		auditRule := audit.NewAudit(
			"daemonset",
			"更新策略",
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
						Const: "daemonset",
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
			return srv.UpdateStatefulSetUpdateStrategy(ctx, req.(*UpdateDaemonSetUpdateStrategyRequest))
		})
		out, err := h(stdCtx, &in)
		if err != nil {
			return err
		}
		reply := out.(*UpdateDaemonSetUpdateStrategyResponse)
		return ctx.Result(200, reply)
	}
}

type DaemonSetHTTPClient interface {
	ListDaemonSet(ctx context.Context, req *ListDaemonSetRequest, opts ...http.CallOption) (rsp *ListDaemonSetResponse, err error)
	CreateOrUpdateDaemonSetByYaml(ctx context.Context, req *CreateOrUpdateDaemonSetByYamlRequest, opts ...http.CallOption) (rsp *CreateOrUpdateDaemonSetByYamlResponse, err error)
	DeleteDaemonSet(ctx context.Context, req *DeleteDaemonSetRequest, opts ...http.CallOption) (rsp *DeleteDaemonSetResponse, err error)
	RestartDaemonSet(ctx context.Context, req *RestartDaemonSetRequest, opts ...http.CallOption) (rsp *RestartDaemonSetResponse, err error)
	GetDaemonSetDetail(ctx context.Context, req *GetDaemonSetDetailRequest, opts ...http.CallOption) (rsp *GetDaemonSetDetailResponse, err error)
	GetDaemonSetRevisions(ctx context.Context, req *GetDaemonSetHistoryRequest, opts ...http.CallOption) (rsp *GetDaemonSetHistoryResponse, err error)
	RollbackDaemonSet(ctx context.Context, req *RollbackDaemonSetRequest, opts ...http.CallOption) (rsp *RollbackDaemonSetResponse, err error)
	UpdateStatefulSetUpdateStrategy(ctx context.Context, req *UpdateDaemonSetUpdateStrategyRequest, opts ...http.CallOption) (rsp *UpdateDaemonSetUpdateStrategyResponse, err error)
}

type DaemonSetHTTPClientImpl struct {
	cc *http.Client
}

func NewDaemonSetHTTPClient(client *http.Client) DaemonSetHTTPClient {
	return &DaemonSetHTTPClientImpl{client}
}

func (c *DaemonSetHTTPClientImpl) ListDaemonSet(ctx context.Context, in *ListDaemonSetRequest, opts ...http.CallOption) (*ListDaemonSetResponse, error) {
	var out ListDaemonSetResponse
	pattern := "/api/v1/daemonset/list"
	path := binding.EncodeURL(pattern, in, true)
	opts = append(opts, http.Operation(OperationDaemonSetListDaemonSet))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "GET", path, nil, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *DaemonSetHTTPClientImpl) CreateOrUpdateDaemonSetByYaml(ctx context.Context, in *CreateOrUpdateDaemonSetByYamlRequest, opts ...http.CallOption) (*CreateOrUpdateDaemonSetByYamlResponse, error) {
	var out CreateOrUpdateDaemonSetByYamlResponse
	pattern := "/api/v1/daemonset/create_or_update_by_yaml"
	path := binding.EncodeURL(pattern, in, false)
	opts = append(opts, http.Operation(OperationDaemonSetCreateOrUpdateDaemonSetByYaml))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "POST", path, in, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *DaemonSetHTTPClientImpl) DeleteDaemonSet(ctx context.Context, in *DeleteDaemonSetRequest, opts ...http.CallOption) (*DeleteDaemonSetResponse, error) {
	var out DeleteDaemonSetResponse
	pattern := "/api/v1/daemonset/delete"
	path := binding.EncodeURL(pattern, in, false)
	opts = append(opts, http.Operation(OperationDaemonSetDeleteDaemonSet))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "POST", path, in, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *DaemonSetHTTPClientImpl) RestartDaemonSet(ctx context.Context, in *RestartDaemonSetRequest, opts ...http.CallOption) (*RestartDaemonSetResponse, error) {
	var out RestartDaemonSetResponse
	pattern := "/api/v1/daemonset/restart"
	path := binding.EncodeURL(pattern, in, false)
	opts = append(opts, http.Operation(OperationDaemonSetRestartDaemonSet))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "POST", path, in, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *DaemonSetHTTPClientImpl) GetDaemonSetDetail(ctx context.Context, in *GetDaemonSetDetailRequest, opts ...http.CallOption) (*GetDaemonSetDetailResponse, error) {
	var out GetDaemonSetDetailResponse
	pattern := "/api/v1/daemonset/detail"
	path := binding.EncodeURL(pattern, in, true)
	opts = append(opts, http.Operation(OperationDaemonSetGetDaemonSetDetail))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "GET", path, nil, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *DaemonSetHTTPClientImpl) GetDaemonSetRevisions(ctx context.Context, in *GetDaemonSetHistoryRequest, opts ...http.CallOption) (*GetDaemonSetHistoryResponse, error) {
	var out GetDaemonSetHistoryResponse
	pattern := "/api/v1/daemonset/revisions"
	path := binding.EncodeURL(pattern, in, true)
	opts = append(opts, http.Operation(OperationDaemonSetGetDaemonSetRevisions))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "GET", path, nil, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *DaemonSetHTTPClientImpl) RollbackDaemonSet(ctx context.Context, in *RollbackDaemonSetRequest, opts ...http.CallOption) (*RollbackDaemonSetResponse, error) {
	var out RollbackDaemonSetResponse
	pattern := "/api/v1/daemonset/rollback"
	path := binding.EncodeURL(pattern, in, false)
	opts = append(opts, http.Operation(OperationDaemonSetRollbackDaemonSet))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "POST", path, in, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *DaemonSetHTTPClientImpl) UpdateStatefulSetUpdateStrategy(ctx context.Context, in *UpdateDaemonSetUpdateStrategyRequest, opts ...http.CallOption) (*UpdateDaemonSetUpdateStrategyResponse, error) {
	var out UpdateDaemonSetUpdateStrategyResponse
	pattern := "/api/v1/daemonset/upgrade_strategy/update"
	path := binding.EncodeURL(pattern, in, false)
	opts = append(opts, http.Operation(OperationDaemonSetUpdateStatefulSetUpdateStrategy))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "POST", path, in, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, nil
}
