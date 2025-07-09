package data

import (
	"context"
	"strconv"
	"time"

	"codo-cnmp/internal/biz"
	"codo-cnmp/internal/model/dao"
	"codo-cnmp/internal/model/model/entity"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/os/gtime"
)

type AuditLogRepo struct {
	data *Data
	log  *log.Helper
}

func (x *AuditLogRepo) ListQueryCondition(ctx context.Context) (*biz.QueryCondition, error) {
	model := dao.AuditLog.Ctx(ctx)
	result := &biz.QueryCondition{
		Cluster:      make([]string, 0),
		Namespace:    make([]string, 0),
		Module:       make([]string, 0),
		ResourceType: make([]string, 0),
		Action:       make([]string, 0),
		Status:       make([]int, 0),
	}
	// 定义接收结果的结构
	type Record struct {
		FieldType string `json:"field_type"`
		Value     string `json:"value"`
	}

	var records []Record
	err := model.UnionAll(
		// cluster
		model.Fields(
			"'cluster' as field_type",
			dao.AuditLog.Columns().Cluster+" as value",
		).Where(
			dao.AuditLog.Columns().Cluster+" IS NOT NULL AND "+
				dao.AuditLog.Columns().Cluster+" <> ''",
		).Group(dao.AuditLog.Columns().Cluster),

		// namespace
		model.Fields(
			"'namespace' as field_type",
			dao.AuditLog.Columns().Namespace+" as value",
		).Where(
			dao.AuditLog.Columns().Namespace+" IS NOT NULL AND "+
				dao.AuditLog.Columns().Namespace+" <> ''",
		).Group(dao.AuditLog.Columns().Namespace),

		// module
		model.Fields(
			"'module' as field_type",
			dao.AuditLog.Columns().Module+" as value",
		).Where(
			dao.AuditLog.Columns().Module+" IS NOT NULL AND "+
				dao.AuditLog.Columns().Module+" <> ''",
		).Group(dao.AuditLog.Columns().Module),

		// resource_type
		model.Fields(
			"'resource_type' as field_type",
			dao.AuditLog.Columns().ResourceType+" as value",
		).Where(
			dao.AuditLog.Columns().ResourceType+" IS NOT NULL AND "+
				dao.AuditLog.Columns().ResourceType+" <> ''",
		).Group(dao.AuditLog.Columns().ResourceType),

		// action
		model.Fields(
			"'action' as field_type",
			dao.AuditLog.Columns().Action+" as value",
		).Where(
			dao.AuditLog.Columns().Action+" IS NOT NULL AND "+
				dao.AuditLog.Columns().Action+" <> ''",
		).Group(dao.AuditLog.Columns().Action),

		// status
		model.Fields(
			"'status' as field_type",
			dao.AuditLog.Columns().Status+" as value",
		).Where(
			dao.AuditLog.Columns().Status+" IS NOT NULL AND "+
				dao.AuditLog.Columns().Status+" <> ''",
		).Group(dao.AuditLog.Columns().Status),
	).Order("field_type, value").Scan(&records)
	if err != nil {
		return nil, err
	}

	// 将结果分类到对应的切片中
	for _, record := range records {
		switch record.FieldType {
		case "cluster":
			result.Cluster = append(result.Cluster, record.Value)
		case "namespace":
			result.Namespace = append(result.Namespace, record.Value)
		case "module":
			result.Module = append(result.Module, record.Value)
		case "resource_type":
			result.ResourceType = append(result.ResourceType, record.Value)
		case "action":
			result.Action = append(result.Action, record.Value)
		case "status":
			status, err := strconv.Atoi(record.Value)
			if err != nil {
				x.log.Warnf("无效的状态值: %v", err)
				continue
			}
			result.Status = append(result.Status, status)
		}
	}

	return result, nil
}

func (x *AuditLogRepo) convertQuery(ctx context.Context, query *biz.ListAuditLogRequest) *gdb.Model {
	model := dao.AuditLog.Ctx(ctx)
	if query.UserName != "" {
		model = model.Where(dao.AuditLog.Columns().Username, query.UserName)
	}
	if query.ClientIP != "" {
		model = model.Where(dao.AuditLog.Columns().ClientIp, query.ClientIP)
	}
	if query.Cluster != "" {
		model = model.Where(dao.AuditLog.Columns().Cluster, query.Cluster)
	}
	if query.Namespace != "" {
		model = model.Where(dao.AuditLog.Columns().Namespace, query.Namespace)
	}
	if query.Module != "" {
		model = model.Where(dao.AuditLog.Columns().Module, query.Module)
	}
	if query.ResourceType != "" {
		model = model.Where(dao.AuditLog.Columns().ResourceType, query.ResourceType)
	}
	if query.ResourceName != "" {
		model = model.Where(dao.AuditLog.Columns().ResourceName, query.ResourceName)
	}
	if query.Action != "" {
		model = model.Where(dao.AuditLog.Columns().Action, query.Action)
	}
	if query.Status > 0 {
		model = model.Where(dao.AuditLog.Columns().Status, query.Status)
	}
	if query.StartTime != "" {
		model = model.WhereGTE(dao.AuditLog.Columns().CreatedAt, query.StartTime)
	}
	if query.EndTime != "" {
		model = model.WhereLTE(dao.AuditLog.Columns().CreatedAt, query.EndTime)
	}
	// 默认按照创建时间倒序排列
	model = model.OrderDesc(dao.AuditLog.Columns().CreatedAt)
	if query.ListAll {
		return model
	}
	if query.Page == 0 {
		query.Page = 1
	}
	if query.PageSize == 0 {
		query.PageSize = 10
	}
	return model.Page(int(query.Page), int(query.PageSize))
}

func (x *AuditLogRepo) convertDO2VO(data *biz.CreateAuditLogRequest) *entity.AuditLog {
	operationTime, err := gtime.StrToTime(data.OperationTime, time.DateTime)
	if err != nil {
		operationTime = gtime.Now()
	}
	return &entity.AuditLog{
		Username:      data.UserName,
		ClientIp:      data.ClientIP,
		Cluster:       data.Cluster,
		Namespace:     data.Namespace,
		Module:        data.Module,
		ResourceType:  data.ResourceType,
		ResourceName:  data.ResourceName,
		Action:        data.Action,
		Status:        data.Status,
		Duration:      data.Duration,
		RequestPath:   data.RequestPath,
		RequestBody:   data.RequestBody,
		ResponseBody:  data.ResponseBody,
		OperationTime: operationTime,
		TraceId:       data.TraceID,
	}
}

func (x *AuditLogRepo) convertVO2DO(data *entity.AuditLog) *biz.AuditLogItem {
	return &biz.AuditLogItem{
		UserName:      data.Username,
		ClientIP:      data.ClientIp,
		Cluster:       data.Cluster,
		Namespace:     data.Namespace,
		Module:        data.Module,
		ResourceType:  data.ResourceType,
		ResourceName:  data.ResourceName,
		Action:        data.Action,
		Status:        data.Status,
		Duration:      data.Duration,
		Id:            uint32(data.Id),
		CreateTime:    data.CreatedAt.String(),
		OperationTime: data.OperationTime.String(),
		RequestPath:   data.RequestPath,
		RequestBody:   data.RequestBody,
		ResponseBody:  data.ResponseBody,
		TraceID:       data.TraceId,
	}
}

func (x *AuditLogRepo) Get(ctx context.Context, id uint32) (*biz.AuditLogItem, error) {
	db := dao.AuditLog.Ctx(ctx)
	var auditLog entity.AuditLog
	err := db.Where(dao.AuditLog.Columns().Id, id).Scan(&auditLog)
	if err != nil {
		return nil, err
	}
	return x.convertVO2DO(&auditLog), nil
}

func (x *AuditLogRepo) Count(ctx context.Context, req *biz.ListAuditLogRequest) (uint32, error) {
	db := x.convertQuery(ctx, req)
	count, err := db.Count()
	if err != nil {
		return 0, err
	}
	return uint32(count), nil
}

func (x *AuditLogRepo) List(ctx context.Context, req *biz.ListAuditLogRequest) ([]*biz.AuditLogItem, uint32, error) {
	result := make([]*biz.AuditLogItem, 0)
	db := dao.AuditLog.Ctx(ctx)
	db = x.convertQuery(ctx, req)
	auditLogs := make([]*entity.AuditLog, 0)
	err := db.Scan(&auditLogs)
	if err != nil {
		return result, 0, err
	}
	count, err := db.Count()
	if err != nil {
		return result, 0, err
	}
	for _, auditLog := range auditLogs {
		item := x.convertVO2DO(auditLog)
		result = append(result, item)
	}
	return result, uint32(count), nil
}

func (x *AuditLogRepo) Create(ctx context.Context, req *biz.CreateAuditLogRequest) error {
	db := dao.AuditLog.Ctx(ctx)
	vo := x.convertDO2VO(req)
	_, err := db.Insert(vo)
	return err
}

func NewAuditLogRepo(data *Data, logger log.Logger) *AuditLogRepo {
	return &AuditLogRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func NewIAuditLogRepo(repo *AuditLogRepo) biz.IAuditLogRepo {
	return repo
}
