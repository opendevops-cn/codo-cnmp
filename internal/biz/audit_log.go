package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
)

type AuditLogCommonParams struct {
	// 操作人
	UserName string
	// 操作IP
	ClientIP string
	// 集群
	Cluster string
	// 命名空间
	Namespace string
	// 模块
	Module string
	// 对象类型
	ResourceType string
	// 对象
	ResourceName string
	// 操作类型
	Action string
	// 操作状态
	Status int
	// 开始时间
	StartTime string
	// 结束时间
	EndTime string
}

type AuditLogItem struct {
	// 操作人
	UserName string `json:"user_name"`
	// 操作IP
	ClientIP string `json:"client_ip"`
	// 集群
	Cluster string `json:"cluster"`
	// 命名空间
	Namespace string `json:"namespace"`
	// 模块
	Module string `json:"module"`
	// 对象类型
	ResourceType string `json:"resource_type"`
	// 对象
	ResourceName string `json:"resource_name"`
	// 操作类型
	Action string `json:"action"`
	// 操作状态
	Status int `json:"status"`
	// Duration
	Duration   string `json:"duration"`
	Id         uint32
	CreateTime string `json:"create_time"`
	// 操作时间
	OperationTime string `json:"operation_time"`
	// 请求路径
	RequestPath string `json:"request_path"`
	// 请求内容
	RequestBody string `json:"request_body"`
	// 响应内容
	ResponseBody string `json:"response_body"`
	// TraceID
	TraceID string `json:"trace_id"`
}

type ListAuditLogRequest struct {
	AuditLogCommonParams
	// 页码
	Page uint32
	// 每页数量
	PageSize uint32
	// 是否全部
	ListAll bool
}

type CreateAuditLogRequest struct {
	AuditLogItem
}

type QueryCondition struct {
	// 集群
	Cluster []string
	// 命名空间
	Namespace []string
	// 模块
	Module []string
	// 对象类型
	ResourceType []string
	// 操作
	Action []string
	// 操作状态
	Status []int
}

type IAuditLogRepo interface {
	List(ctx context.Context, req *ListAuditLogRequest) ([]*AuditLogItem, uint32, error)
	Create(ctx context.Context, req *CreateAuditLogRequest) error
	Get(ctx context.Context, id uint32) (*AuditLogItem, error)
	ListQueryCondition(ctx context.Context) (*QueryCondition, error)
}

type IAuditLogUseCase interface {
	// ListAuditLog 审计日志列表
	ListAuditLog(ctx context.Context, req *ListAuditLogRequest) ([]*AuditLogItem, uint32, error)
	// CreateAuditLog 创建审计日志
	CreateAuditLog(ctx context.Context, req *CreateAuditLogRequest) error
	// GetAuditLog 获取审计日志
	GetAuditLog(ctx context.Context, id uint32) (*AuditLogItem, error)
	// ListQueryCondition 审计日志查询条件
	ListQueryCondition(ctx context.Context) (*QueryCondition, error)
}

type AuditLogUseCase struct {
	repo IAuditLogRepo
	log  *log.Helper
}

func (x *AuditLogUseCase) ListQueryCondition(ctx context.Context) (*QueryCondition, error) {
	return x.repo.ListQueryCondition(ctx)
}

func (x *AuditLogUseCase) GetAuditLog(ctx context.Context, id uint32) (*AuditLogItem, error) {
	return x.repo.Get(ctx, id)
}

func (x *AuditLogUseCase) ListAuditLog(ctx context.Context, req *ListAuditLogRequest) ([]*AuditLogItem, uint32, error) {
	return x.repo.List(ctx, req)
}

func (x *AuditLogUseCase) CreateAuditLog(ctx context.Context, req *CreateAuditLogRequest) error {
	return x.repo.Create(ctx, req)
}

func NewAuditLogUseCase(repo IAuditLogRepo, logger log.Logger) *AuditLogUseCase {
	return &AuditLogUseCase{repo: repo, log: log.NewHelper(log.With(logger, "module", "biz/audit_log"))}
}

func NewIAuditLogUseCase(x *AuditLogUseCase) IAuditLogUseCase {
	return x
}
