// Code generated by protoc-gen-go-http. DO NOT EDIT.
// versions:
// - protoc-gen-go-http v2.7.3
// - protoc             v5.27.2
// source: pb/statefulset.v1.proto

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

const OperationStatefulSetListStatefulSet = "/statefulset.StatefulSet/ListStatefulSet"
const OperationStatefulSetCreateOrUpdateStatefulSetByYaml = "/statefulset.StatefulSet/CreateOrUpdateStatefulSetByYaml"
const OperationStatefulSetDeleteStatefulSet = "/statefulset.StatefulSet/DeleteStatefulSet"
const OperationStatefulSetRestartStatefulSet = "/statefulset.StatefulSet/RestartStatefulSet"
const OperationStatefulSetScaleStatefulSet = "/statefulset.StatefulSet/ScaleStatefulSet"
const OperationStatefulSetGetStatefulSetDetail = "/statefulset.StatefulSet/GetStatefulSetDetail"
const OperationStatefulSetGetStatefulSetRevisions = "/statefulset.StatefulSet/GetStatefulSetRevisions"
const OperationStatefulSetRollbackStatefulSet = "/statefulset.StatefulSet/RollbackStatefulSet"
const OperationStatefulSetUpdateStatefulSetUpdateStrategy = "/statefulset.StatefulSet/UpdateStatefulSetUpdateStrategy"

type StatefulSetHTTPServer interface {
	// ListStatefulSet查看-云原生管理-StatefulSet-列表
	ListStatefulSet(context.Context, *ListStatefulSetRequest) (*ListStatefulSetResponse, error)
	// CreateOrUpdateStatefulSetByYaml管理-云原生管理-StatefulSet-Yaml创建更新
	CreateOrUpdateStatefulSetByYaml(context.Context, *CreateOrUpdateStatefulSetByYamlRequest) (*CreateOrUpdateStatefulSetByYamlResponse, error)
	// DeleteStatefulSet管理-云原生管理-StatefulSet-删除
	DeleteStatefulSet(context.Context, *DeleteStatefulSetRequest) (*DeleteStatefulSetResponse, error)
	// RestartStatefulSet管理-云原生管理-StatefulSet-重启
	RestartStatefulSet(context.Context, *RestartStatefulSetRequest) (*RestartStatefulSetResponse, error)
	// ScaleStatefulSet管理-云原生管理-StatefulSet-伸缩
	ScaleStatefulSet(context.Context, *ScaleStatefulSetRequest) (*ScaleStatefulSetResponse, error)
	// GetStatefulSetDetail查看-云原生管理-StatefulSet-详情
	GetStatefulSetDetail(context.Context, *GetStatefulSetDetailRequest) (*GetStatefulSetDetailResponse, error)
	// GetStatefulSetRevisions查看-云原生管理-StatefulSet-历史版本
	GetStatefulSetRevisions(context.Context, *GetStatefulSetHistoryRequest) (*GetStatefulSetHistoryResponse, error)
	// RollbackStatefulSet管理-云原生管理-StatefulSet-回滚
	RollbackStatefulSet(context.Context, *RollbackStatefulSetRequest) (*RollbackStatefulSetResponse, error)
	// UpdateStatefulSetUpdateStrategy管理-云原生管理-StatefulSet-更新策略
	UpdateStatefulSetUpdateStrategy(context.Context, *UpdateStatefulSetUpdateStrategyRequest) (*UpdateStatefulSetUpdateStrategyResponse, error)
}

func NewStatefulSetHTTPServerMiddleware() middleware.Middleware {
	return selector.Server(
		selector.Server().Path(OperationStatefulSetListStatefulSet).Build(),
		selector.Server().Path(OperationStatefulSetCreateOrUpdateStatefulSetByYaml).Build(),
		selector.Server().Path(OperationStatefulSetDeleteStatefulSet).Build(),
		selector.Server().Path(OperationStatefulSetRestartStatefulSet).Build(),
		selector.Server().Path(OperationStatefulSetScaleStatefulSet).Build(),
		selector.Server().Path(OperationStatefulSetGetStatefulSetDetail).Build(),
		selector.Server().Path(OperationStatefulSetGetStatefulSetRevisions).Build(),
		selector.Server().Path(OperationStatefulSetRollbackStatefulSet).Build(),
		selector.Server().Path(OperationStatefulSetUpdateStatefulSetUpdateStrategy).Build(),
	).Path(
		OperationStatefulSetListStatefulSet,
		OperationStatefulSetCreateOrUpdateStatefulSetByYaml,
		OperationStatefulSetDeleteStatefulSet,
		OperationStatefulSetRestartStatefulSet,
		OperationStatefulSetScaleStatefulSet,
		OperationStatefulSetGetStatefulSetDetail,
		OperationStatefulSetGetStatefulSetRevisions,
		OperationStatefulSetRollbackStatefulSet,
		OperationStatefulSetUpdateStatefulSetUpdateStrategy,
	).Build()
}

func RegisterStatefulSetHTTPServer(s *http.Server, srv StatefulSetHTTPServer) {
	r := s.Route("/")
	r.GET("/api/v1/statefulset/list", _StatefulSet_ListStatefulSet0_HTTP_Handler(srv))
	r.POST("/api/v1/statefulset/create_or_update_by_yaml", _StatefulSet_CreateOrUpdateStatefulSetByYaml0_HTTP_Handler(srv))
	r.POST("/api/v1/statefulset/delete", _StatefulSet_DeleteStatefulSet0_HTTP_Handler(srv))
	r.POST("/api/v1/statefulset/restart", _StatefulSet_RestartStatefulSet0_HTTP_Handler(srv))
	r.POST("/api/v1/statefulset/scale", _StatefulSet_ScaleStatefulSet0_HTTP_Handler(srv))
	r.GET("/api/v1/statefulset/detail", _StatefulSet_GetStatefulSetDetail0_HTTP_Handler(srv))
	r.GET("/api/v1/statefulset/revisions", _StatefulSet_GetStatefulSetRevisions0_HTTP_Handler(srv))
	r.POST("/api/v1/statefulset/rollback", _StatefulSet_RollbackStatefulSet0_HTTP_Handler(srv))
	r.POST("/api/v1/statefulset/upgrade_strategy/update", _StatefulSet_UpdateStatefulSetUpdateStrategy1_HTTP_Handler(srv))
}

func GenerateStatefulSetHTTPServerRouteInfo() []route.Route {
	routes := make([]route.Route, 0, 9)
	routes = append(routes, route.Route{
		Method:  "GET",
		Path:    "/api/v1/statefulset/list",
		Comment: "查看-云原生管理-StatefulSet-列表",
	})
	routes = append(routes, route.Route{
		Method:  "POST",
		Path:    "/api/v1/statefulset/create_or_update_by_yaml",
		Comment: "管理-云原生管理-StatefulSet-Yaml创建更新",
	})
	routes = append(routes, route.Route{
		Method:  "POST",
		Path:    "/api/v1/statefulset/delete",
		Comment: "管理-云原生管理-StatefulSet-删除",
	})
	routes = append(routes, route.Route{
		Method:  "POST",
		Path:    "/api/v1/statefulset/restart",
		Comment: "管理-云原生管理-StatefulSet-重启",
	})
	routes = append(routes, route.Route{
		Method:  "POST",
		Path:    "/api/v1/statefulset/scale",
		Comment: "管理-云原生管理-StatefulSet-伸缩",
	})
	routes = append(routes, route.Route{
		Method:  "GET",
		Path:    "/api/v1/statefulset/detail",
		Comment: "查看-云原生管理-StatefulSet-详情",
	})
	routes = append(routes, route.Route{
		Method:  "GET",
		Path:    "/api/v1/statefulset/revisions",
		Comment: "查看-云原生管理-StatefulSet-历史版本",
	})
	routes = append(routes, route.Route{
		Method:  "POST",
		Path:    "/api/v1/statefulset/rollback",
		Comment: "管理-云原生管理-StatefulSet-回滚",
	})
	routes = append(routes, route.Route{
		Method:  "POST",
		Path:    "/api/v1/statefulset/upgrade_strategy/update",
		Comment: "管理-云原生管理-StatefulSet-更新策略",
	})
	return routes
}

func _StatefulSet_ListStatefulSet0_HTTP_Handler(srv StatefulSetHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		stdCtx := kcontext.SetKHTTPContextWithContext(ctx, ctx)
		var in ListStatefulSetRequest
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationStatefulSetListStatefulSet)

		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.ListStatefulSet(ctx, req.(*ListStatefulSetRequest))
		})
		out, err := h(stdCtx, &in)
		if err != nil {
			return err
		}
		reply := out.(*ListStatefulSetResponse)
		return ctx.Result(200, reply)
	}
}

func _StatefulSet_CreateOrUpdateStatefulSetByYaml0_HTTP_Handler(srv StatefulSetHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		stdCtx := kcontext.SetKHTTPContextWithContext(ctx, ctx)
		var in CreateOrUpdateStatefulSetByYamlRequest
		if err := ctx.Bind(&in); err != nil {
			return err
		}
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationStatefulSetCreateOrUpdateStatefulSetByYaml)
		auditRule := audit.NewAudit(
			"statefulset",
			"yaml创建更新",
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
						Const: "statefulset",
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
			return srv.CreateOrUpdateStatefulSetByYaml(ctx, req.(*CreateOrUpdateStatefulSetByYamlRequest))
		})
		out, err := h(stdCtx, &in)
		if err != nil {
			return err
		}
		reply := out.(*CreateOrUpdateStatefulSetByYamlResponse)
		return ctx.Result(200, reply)
	}
}

func _StatefulSet_DeleteStatefulSet0_HTTP_Handler(srv StatefulSetHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		stdCtx := kcontext.SetKHTTPContextWithContext(ctx, ctx)
		var in DeleteStatefulSetRequest
		if err := ctx.Bind(&in); err != nil {
			return err
		}
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationStatefulSetDeleteStatefulSet)
		auditRule := audit.NewAudit(
			"statefulset",
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
						Const: "statefulset",
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
			return srv.DeleteStatefulSet(ctx, req.(*DeleteStatefulSetRequest))
		})
		out, err := h(stdCtx, &in)
		if err != nil {
			return err
		}
		reply := out.(*DeleteStatefulSetResponse)
		return ctx.Result(200, reply)
	}
}

func _StatefulSet_RestartStatefulSet0_HTTP_Handler(srv StatefulSetHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		stdCtx := kcontext.SetKHTTPContextWithContext(ctx, ctx)
		var in RestartStatefulSetRequest
		if err := ctx.Bind(&in); err != nil {
			return err
		}
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationStatefulSetRestartStatefulSet)
		auditRule := audit.NewAudit(
			"statefulset",
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
						Const: "statefulset",
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
			return srv.RestartStatefulSet(ctx, req.(*RestartStatefulSetRequest))
		})
		out, err := h(stdCtx, &in)
		if err != nil {
			return err
		}
		reply := out.(*RestartStatefulSetResponse)
		return ctx.Result(200, reply)
	}
}

func _StatefulSet_ScaleStatefulSet0_HTTP_Handler(srv StatefulSetHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		stdCtx := kcontext.SetKHTTPContextWithContext(ctx, ctx)
		var in ScaleStatefulSetRequest
		if err := ctx.Bind(&in); err != nil {
			return err
		}
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationStatefulSetScaleStatefulSet)
		auditRule := audit.NewAudit(
			"statefulset",
			"伸缩",
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
						Const: "statefulset",
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
			return srv.ScaleStatefulSet(ctx, req.(*ScaleStatefulSetRequest))
		})
		out, err := h(stdCtx, &in)
		if err != nil {
			return err
		}
		reply := out.(*ScaleStatefulSetResponse)
		return ctx.Result(200, reply)
	}
}

func _StatefulSet_GetStatefulSetDetail0_HTTP_Handler(srv StatefulSetHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		stdCtx := kcontext.SetKHTTPContextWithContext(ctx, ctx)
		var in GetStatefulSetDetailRequest
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationStatefulSetGetStatefulSetDetail)

		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.GetStatefulSetDetail(ctx, req.(*GetStatefulSetDetailRequest))
		})
		out, err := h(stdCtx, &in)
		if err != nil {
			return err
		}
		reply := out.(*GetStatefulSetDetailResponse)
		return ctx.Result(200, reply)
	}
}

func _StatefulSet_GetStatefulSetRevisions0_HTTP_Handler(srv StatefulSetHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		stdCtx := kcontext.SetKHTTPContextWithContext(ctx, ctx)
		var in GetStatefulSetHistoryRequest
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationStatefulSetGetStatefulSetRevisions)

		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.GetStatefulSetRevisions(ctx, req.(*GetStatefulSetHistoryRequest))
		})
		out, err := h(stdCtx, &in)
		if err != nil {
			return err
		}
		reply := out.(*GetStatefulSetHistoryResponse)
		return ctx.Result(200, reply)
	}
}

func _StatefulSet_RollbackStatefulSet0_HTTP_Handler(srv StatefulSetHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		stdCtx := kcontext.SetKHTTPContextWithContext(ctx, ctx)
		var in RollbackStatefulSetRequest
		if err := ctx.Bind(&in); err != nil {
			return err
		}
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationStatefulSetRollbackStatefulSet)
		auditRule := audit.NewAudit(
			"statefulset",
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
						Const: "statefulset",
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
			return srv.RollbackStatefulSet(ctx, req.(*RollbackStatefulSetRequest))
		})
		out, err := h(stdCtx, &in)
		if err != nil {
			return err
		}
		reply := out.(*RollbackStatefulSetResponse)
		return ctx.Result(200, reply)
	}
}

func _StatefulSet_UpdateStatefulSetUpdateStrategy1_HTTP_Handler(srv StatefulSetHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		stdCtx := kcontext.SetKHTTPContextWithContext(ctx, ctx)
		var in UpdateStatefulSetUpdateStrategyRequest
		if err := ctx.Bind(&in); err != nil {
			return err
		}
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationStatefulSetUpdateStatefulSetUpdateStrategy)
		auditRule := audit.NewAudit(
			"statefulset",
			"修改更新策略",
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
						Const: "statefulset",
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
			return srv.UpdateStatefulSetUpdateStrategy(ctx, req.(*UpdateStatefulSetUpdateStrategyRequest))
		})
		out, err := h(stdCtx, &in)
		if err != nil {
			return err
		}
		reply := out.(*UpdateStatefulSetUpdateStrategyResponse)
		return ctx.Result(200, reply)
	}
}

type StatefulSetHTTPClient interface {
	ListStatefulSet(ctx context.Context, req *ListStatefulSetRequest, opts ...http.CallOption) (rsp *ListStatefulSetResponse, err error)
	CreateOrUpdateStatefulSetByYaml(ctx context.Context, req *CreateOrUpdateStatefulSetByYamlRequest, opts ...http.CallOption) (rsp *CreateOrUpdateStatefulSetByYamlResponse, err error)
	DeleteStatefulSet(ctx context.Context, req *DeleteStatefulSetRequest, opts ...http.CallOption) (rsp *DeleteStatefulSetResponse, err error)
	RestartStatefulSet(ctx context.Context, req *RestartStatefulSetRequest, opts ...http.CallOption) (rsp *RestartStatefulSetResponse, err error)
	ScaleStatefulSet(ctx context.Context, req *ScaleStatefulSetRequest, opts ...http.CallOption) (rsp *ScaleStatefulSetResponse, err error)
	GetStatefulSetDetail(ctx context.Context, req *GetStatefulSetDetailRequest, opts ...http.CallOption) (rsp *GetStatefulSetDetailResponse, err error)
	GetStatefulSetRevisions(ctx context.Context, req *GetStatefulSetHistoryRequest, opts ...http.CallOption) (rsp *GetStatefulSetHistoryResponse, err error)
	RollbackStatefulSet(ctx context.Context, req *RollbackStatefulSetRequest, opts ...http.CallOption) (rsp *RollbackStatefulSetResponse, err error)
	UpdateStatefulSetUpdateStrategy(ctx context.Context, req *UpdateStatefulSetUpdateStrategyRequest, opts ...http.CallOption) (rsp *UpdateStatefulSetUpdateStrategyResponse, err error)
}

type StatefulSetHTTPClientImpl struct {
	cc *http.Client
}

func NewStatefulSetHTTPClient(client *http.Client) StatefulSetHTTPClient {
	return &StatefulSetHTTPClientImpl{client}
}

func (c *StatefulSetHTTPClientImpl) ListStatefulSet(ctx context.Context, in *ListStatefulSetRequest, opts ...http.CallOption) (*ListStatefulSetResponse, error) {
	var out ListStatefulSetResponse
	pattern := "/api/v1/statefulset/list"
	path := binding.EncodeURL(pattern, in, true)
	opts = append(opts, http.Operation(OperationStatefulSetListStatefulSet))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "GET", path, nil, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *StatefulSetHTTPClientImpl) CreateOrUpdateStatefulSetByYaml(ctx context.Context, in *CreateOrUpdateStatefulSetByYamlRequest, opts ...http.CallOption) (*CreateOrUpdateStatefulSetByYamlResponse, error) {
	var out CreateOrUpdateStatefulSetByYamlResponse
	pattern := "/api/v1/statefulset/create_or_update_by_yaml"
	path := binding.EncodeURL(pattern, in, false)
	opts = append(opts, http.Operation(OperationStatefulSetCreateOrUpdateStatefulSetByYaml))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "POST", path, in, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *StatefulSetHTTPClientImpl) DeleteStatefulSet(ctx context.Context, in *DeleteStatefulSetRequest, opts ...http.CallOption) (*DeleteStatefulSetResponse, error) {
	var out DeleteStatefulSetResponse
	pattern := "/api/v1/statefulset/delete"
	path := binding.EncodeURL(pattern, in, false)
	opts = append(opts, http.Operation(OperationStatefulSetDeleteStatefulSet))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "POST", path, in, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *StatefulSetHTTPClientImpl) RestartStatefulSet(ctx context.Context, in *RestartStatefulSetRequest, opts ...http.CallOption) (*RestartStatefulSetResponse, error) {
	var out RestartStatefulSetResponse
	pattern := "/api/v1/statefulset/restart"
	path := binding.EncodeURL(pattern, in, false)
	opts = append(opts, http.Operation(OperationStatefulSetRestartStatefulSet))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "POST", path, in, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *StatefulSetHTTPClientImpl) ScaleStatefulSet(ctx context.Context, in *ScaleStatefulSetRequest, opts ...http.CallOption) (*ScaleStatefulSetResponse, error) {
	var out ScaleStatefulSetResponse
	pattern := "/api/v1/statefulset/scale"
	path := binding.EncodeURL(pattern, in, false)
	opts = append(opts, http.Operation(OperationStatefulSetScaleStatefulSet))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "POST", path, in, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *StatefulSetHTTPClientImpl) GetStatefulSetDetail(ctx context.Context, in *GetStatefulSetDetailRequest, opts ...http.CallOption) (*GetStatefulSetDetailResponse, error) {
	var out GetStatefulSetDetailResponse
	pattern := "/api/v1/statefulset/detail"
	path := binding.EncodeURL(pattern, in, true)
	opts = append(opts, http.Operation(OperationStatefulSetGetStatefulSetDetail))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "GET", path, nil, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *StatefulSetHTTPClientImpl) GetStatefulSetRevisions(ctx context.Context, in *GetStatefulSetHistoryRequest, opts ...http.CallOption) (*GetStatefulSetHistoryResponse, error) {
	var out GetStatefulSetHistoryResponse
	pattern := "/api/v1/statefulset/revisions"
	path := binding.EncodeURL(pattern, in, true)
	opts = append(opts, http.Operation(OperationStatefulSetGetStatefulSetRevisions))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "GET", path, nil, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *StatefulSetHTTPClientImpl) RollbackStatefulSet(ctx context.Context, in *RollbackStatefulSetRequest, opts ...http.CallOption) (*RollbackStatefulSetResponse, error) {
	var out RollbackStatefulSetResponse
	pattern := "/api/v1/statefulset/rollback"
	path := binding.EncodeURL(pattern, in, false)
	opts = append(opts, http.Operation(OperationStatefulSetRollbackStatefulSet))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "POST", path, in, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *StatefulSetHTTPClientImpl) UpdateStatefulSetUpdateStrategy(ctx context.Context, in *UpdateStatefulSetUpdateStrategyRequest, opts ...http.CallOption) (*UpdateStatefulSetUpdateStrategyResponse, error) {
	var out UpdateStatefulSetUpdateStrategyResponse
	pattern := "/api/v1/statefulset/upgrade_strategy/update"
	path := binding.EncodeURL(pattern, in, false)
	opts = append(opts, http.Operation(OperationStatefulSetUpdateStatefulSetUpdateStrategy))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "POST", path, in, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, nil
}
