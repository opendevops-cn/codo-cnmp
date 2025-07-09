package service

import (
	"context"
	"fmt"

	"codo-cnmp/common/utils"
	"codo-cnmp/internal/biz"
	"codo-cnmp/pb"
)

type AuditLogService struct {
	pb.UnimplementedAuditLogServer
	uc     *biz.AuditLogUseCase
	userUC *biz.UserUseCase
}

func NewAuditLogService(uc *biz.AuditLogUseCase, userUC *biz.UserUseCase) *AuditLogService {
	return &AuditLogService{uc: uc, userUC: userUC}
}

func (x *AuditLogService) convertDO2DTO(log *biz.AuditLogItem) *pb.AuditLogItem {
	username := log.UserName
	user, _ := x.userUC.GetUserByUsername(context.Background(), log.UserName)
	if user != nil {
		username += fmt.Sprintf("(%s)", user.Nickname)
	}
	return &pb.AuditLogItem{
		Id:            log.Id,
		Username:      username,
		ClientIp:      log.ClientIP,
		Module:        log.Module,
		Cluster:       log.Cluster,
		Namespace:     log.Namespace,
		ResourceType:  log.ResourceType,
		ResourceName:  log.ResourceName,
		Action:        log.Action,
		RequestPath:   log.RequestPath,
		RequestBody:   log.RequestBody,
		ResponseBody:  log.ResponseBody,
		Status:        pb.OperationStatus(int32(log.Status)),
		Duration:      log.Duration,
		OperationTime: setTime(log.OperationTime),
		CreatedTime:   setTime(log.CreateTime),
	}
}

func (x *AuditLogService) ListAuditLog(ctx context.Context, req *pb.ListAuditLogRequest) (*pb.ListAuditLogResponse, error) {
	// 构建查询参数
	username := req.GetUsername()
	if username != "" {
		// 如果用户名不为空，则查询用户信息
		if utils.IsEnglish(username) {
			// 如果是英文用户名，则直接使用
			req.Username = username
		} else {
			// 如果是中文用户名，则查询用户信息
			user, err := x.userUC.GetUserByNickname(ctx, username)
			if err != nil {
				req.Username = username
			}
			if user != nil {
				req.Username = user.Username
			}
		}
	}
	params := &biz.ListAuditLogRequest{
		AuditLogCommonParams: biz.AuditLogCommonParams{
			UserName:     req.Username,
			ClientIP:     req.ClientIp,
			Cluster:      req.Cluster,
			Namespace:    req.Namespace,
			Module:       req.Module,
			ResourceType: req.ResourceType,
			ResourceName: req.ResourceName,
			Action:       req.Action,
			Status:       int(req.Status),
		},
		Page:     req.Page,
		PageSize: req.PageSize,
		ListAll:  utils.IntToBool(int(req.ListAll)),
	}

	// 处理可选参数
	if req.StartTime != nil {
		params.StartTime = *req.StartTime
	}
	if req.EndTime != nil {
		params.EndTime = *req.EndTime
	}

	auditLogs, count, err := x.uc.ListAuditLog(ctx, params)
	if err != nil {
		return nil, err
	}

	// 转换响应数据
	list := make([]*pb.AuditLogItem, 0, len(auditLogs))
	for _, log := range auditLogs {
		list = append(list, x.convertDO2DTO(log))
	}

	// 返回响应
	return &pb.ListAuditLogResponse{
		List:  list,
		Total: count,
	}, nil
}

func (x *AuditLogService) GetAuditLog(ctx context.Context, req *pb.GetAuditLogRequest) (*pb.GetAuditLogResponse, error) {
	auditLog, err := x.uc.GetAuditLog(ctx, uint32(req.Id))
	if err != nil {
		return nil, err
	}

	return &pb.GetAuditLogResponse{
		Detail: x.convertDO2DTO(auditLog),
	}, nil
}

func (x *AuditLogService) ListAuditLogQueryCondition(ctx context.Context, req *pb.AuditLogQueryConditionRequest) (*pb.AuditLogQueryConditionResponse, error) {
	queryConditions, err := x.uc.ListQueryCondition(ctx)
	if err != nil {
		return nil, err
	}

	status := make([]pb.OperationStatus, 0, len(queryConditions.Status))
	for _, condition := range queryConditions.Status {
		status = append(status, pb.OperationStatus(condition))
	}

	return &pb.AuditLogQueryConditionResponse{
		Cluster:      queryConditions.Cluster,
		Module:       queryConditions.Module,
		Action:       queryConditions.Action,
		Namespace:    queryConditions.Namespace,
		ResourceType: queryConditions.ResourceType,
		Status:       status,
	}, nil
}
