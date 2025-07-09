package service

import (
	"codo-cnmp/common/utils"
	"codo-cnmp/internal/biz"
	"codo-cnmp/pb"
	"context"
	"time"
)

type RoleService struct {
	pb.UnimplementedRoleServiceServer
	role *biz.RoleUseCase
}

func NewRoleService(role *biz.RoleUseCase) *RoleService {
	return &RoleService{role: role}
}

func setTime(datetime string) uint64 {
	localTime, err := utils.Datetime2time(datetime, time.DateTime, "Asia/Shanghai")
	if err == nil && !localTime.IsZero() {
		return uint64(localTime.UnixNano() / 1e6)
	}
	return 0
}

func (s *RoleService) CreateRole(ctx context.Context, req *pb.CreateRoleRequest) (*pb.CreateRoleResponse, error) {
	err := s.role.CreateRole(ctx, &biz.RoleItem{
		Name:        req.Name,
		Description: req.Description,
		RoleType:    req.RoleType,
		ISDefault:   false,
		YamlStr:     req.Yaml,
	})
	if err != nil {
		return nil, err
	}
	return &pb.CreateRoleResponse{
		Success: err == nil,
	}, nil
}

func (s *RoleService) DeleteRole(ctx context.Context, req *pb.DeleteRoleRequest) (*pb.DeleteRoleResponse, error) {
	err := s.role.DeleteRole(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &pb.DeleteRoleResponse{
		Success: err == nil,
	}, nil
}

func (s *RoleService) UpdateRole(ctx context.Context, req *pb.UpdateRoleRequest) (*pb.UpdateRoleResponse, error) {
	err := s.role.UpdateRole(ctx, &biz.RoleItem{
		ID:          req.Id,
		Name:        req.Name,
		Description: req.Description,
		RoleType:    req.RoleType,
		YamlStr:     req.Yaml,
	})
	if err != nil {
		return nil, err
	}
	return &pb.UpdateRoleResponse{
		Success: err == nil,
	}, nil
}

func (s *RoleService) ListRoles(ctx context.Context, req *pb.ListRoleRequest) (*pb.ListRoleResponse, error) {
	roles, total, err := s.role.ListRoles(ctx, &biz.ListRoleRequest{
		Page:     req.Page,
		PageSize: req.PageSize,
		Keyword:  req.Keyword,
		ListAll:  utils.IntToBool(int(req.ListAll)),
	})
	if err != nil {
		return nil, err
	}
	list := make([]*pb.RoleItem, 0)
	for _, role := range roles {
		var (
			createTime uint64
			updateTime uint64
		)
		createTime = setTime(role.CreateTime)
		updateTime = setTime(role.UpdateTime)
		list = append(list, &pb.RoleItem{
			Id:          role.ID,
			Name:        role.Name,
			Description: role.Description,
			RoleType:    role.RoleType,
			CreateTime:  createTime,
			UpdateTime:  updateTime,
			Yaml:        role.YamlStr,
			IsDefault:   role.ISDefault,
			UpdateBy:    role.UpdateBy,
		})
	}
	return &pb.ListRoleResponse{
		List:  list,
		Total: total,
	}, nil
}

func (s *RoleService) ListRoleBinding(ctx context.Context, req *pb.ListRoleBindingRequest) (*pb.ListRoleBindingResponse, error) {
	roles, total, err := s.role.ListRoleBindings(ctx, &biz.ListRoleBindingRequest{
		Page:     req.Page,
		PageSize: req.PageSize,
		ListAll:  utils.IntToBool(int(req.ListAll)),
		RoleBindingCommonParams: biz.RoleBindingCommonParams{
			RoleID: req.RoleId,
		},
	})
	if err != nil {
		return nil, err
	}
	list := make([]*pb.RoleBindingItem, 0)
	for _, role := range roles {
		list = append(list, &pb.RoleBindingItem{
			Namespace:   role.Namespace,
			UserGroupId: role.UserGroupID,
			ClusterId:   role.ClusterID,
			RoleId:      role.RoleID,
		})
	}
	return &pb.ListRoleBindingResponse{
		List:  list,
		Total: total,
	}, nil
}

func (s *RoleService) UpdateRoleBinding(ctx context.Context, req *pb.UpdateRoleBindingRequest) (*pb.UpdateRoleBindingResponse, error) {
	data := make([]*biz.RoleBindingItem, 0)
	if len(req.Bindings) == 0 {
		data = append(data, &biz.RoleBindingItem{
			RoleBindingCommonParams: biz.RoleBindingCommonParams{
				RoleID: req.RoleId,
			},
		})

	} else {
		for _, item := range req.Bindings {
			data = append(data, &biz.RoleBindingItem{
				RoleBindingCommonParams: biz.RoleBindingCommonParams{
					RoleID:      req.RoleId,
					ClusterID:   item.ClusterId,
					UserGroupID: item.UserGroupId,
				},
				Namespace: item.Namespace,
			})
		}
	}
	err := s.role.UpdateRoleBinding(ctx, data)
	if err != nil {
		return nil, err
	}
	return &pb.UpdateRoleBindingResponse{
		Success: err == nil,
	}, nil
}

func (s *RoleService) GetRoleDetail(ctx context.Context, req *pb.GetRoleDetailRequest) (*pb.GetRoleDetailResponse, error) {
	role, err := s.role.GetRole(ctx, req.RoleId)
	if err != nil {
		return nil, err
	}
	return &pb.GetRoleDetailResponse{
		Detail: &pb.RoleItem{
			Id:          role.ID,
			Name:        role.Name,
			Description: role.Description,
			RoleType:    role.RoleType,
			Yaml:        role.YamlStr,
			CreateTime:  setTime(role.CreateTime),
			UpdateTime:  setTime(role.UpdateTime),
			IsDefault:   role.ISDefault,
			UpdateBy:    role.UpdateBy,
		},
	}, nil

}
