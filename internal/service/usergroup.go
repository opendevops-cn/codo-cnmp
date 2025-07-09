package service

import (
	"codo-cnmp/common/utils"
	"codo-cnmp/internal/biz"
	pb "codo-cnmp/pb"
	"context"
	"strconv"
)

// UserGroupService implements the Cluster interface.
type UserGroupService struct {
	pb.UnimplementedUserGroupServiceServer
	ug *biz.UserGroupUseCase
}

// NewUserGroupService creates a new UserGroupService object.
func NewUserGroupService(ug *biz.UserGroupUseCase) *UserGroupService {
	return &UserGroupService{ug: ug}
}

// GrantUserGroup 用户组授权
func (x *UserGroupService) GrantUserGroup(ctx context.Context, req *pb.GrantUserGroupRequest) (*pb.GrantUserGroupResponse, error) {
	var roleDetail []biz.RoleDetail
	for _, role := range req.Roles {
		roleDetail = append(roleDetail, biz.RoleDetail{
			RoleID:    role.RoleId,
			ClusterID: role.ClusterId,
			Namespace: role.Namespace,
		})
	}
	err := x.ug.CreateGrantedUserGroup(ctx, &biz.CreateGrantedUserGroupRequest{
		UserGroupIDS: req.UserGroupIds,
		RoleDetail:   roleDetail,
	})
	if err != nil {
		return nil, err
	}
	return &pb.GrantUserGroupResponse{
		Success: err == nil,
	}, nil
}

// ListGrantedUserGroup 用户组授权列表
func (x *UserGroupService) ListGrantedUserGroup(ctx context.Context, req *pb.ListGrantedUserGroupRequest) (*pb.ListGrantedUserGroupResponse, error) {
	userGroupList, total, err := x.ug.ListGrantedUserGroups(ctx, &biz.ListGrantedUserGroupRequest{
		Page:     req.Page,
		PageSize: req.PageSize,
		Keyword:  req.Keyword,
		ListAll:  utils.IntToBool(int(req.ListAll)),
	})
	if err != nil {
		return nil, err
	}
	respList := make([]*pb.GrantedUserGroup, 0)
	for _, userGroup := range userGroupList {
		respList = append(respList, &pb.GrantedUserGroup{
			UserGroupId:         userGroup.UserGroupID,
			UserGroupName:       userGroup.Name,
			GrantedClusterCount: userGroup.GrantedClusterCount,
			GrantedRoleCount:    userGroup.GrantedRoleCount,
			UpdateTime:          setTime(userGroup.UpdateTime),
		})
	}
	return &pb.ListGrantedUserGroupResponse{
		List:  respList,
		Total: total,
	}, nil
}

// DeleteGrantedUserGroup 删除用户组授权
func (x *UserGroupService) DeleteGrantedUserGroup(ctx context.Context, req *pb.DeleteGrantedUserGroupRequest) (*pb.DeleteGrantedUserGroupResponse, error) {
	err := x.ug.DeleteGrantedUserGroup(ctx, req.UserGroupId)
	if err != nil {
		return nil, err
	}
	return &pb.DeleteGrantedUserGroupResponse{
		Success: err == nil,
	}, nil
}

// UpdateGrantedUserGroup 更新用户组授权
func (x *UserGroupService) UpdateGrantedUserGroup(ctx context.Context, req *pb.UpdateGrantedUserGroupRequest) (*pb.UpdateGrantedUserGroupResponse, error) {
	var roleDetail []biz.RoleDetail
	for _, role := range req.Roles {
		roleDetail = append(roleDetail, biz.RoleDetail{
			RoleID:    role.RoleId,
			ClusterID: role.ClusterId,
			Namespace: role.Namespace,
		})
	}
	err := x.ug.UpdateGrantedUserGroup(ctx, &biz.CreateGrantedUserGroupRequest{
		UserGroupIDS: req.UserGroupIds,
		RoleDetail:   roleDetail,
	})
	if err != nil {
		return nil, err
	}
	return &pb.UpdateGrantedUserGroupResponse{
		Success: err == nil,
	}, nil
}

// GetGrantedUserGroupDetail 获取用户组授权详情
func (x *UserGroupService) GetGrantedUserGroupDetail(ctx context.Context, req *pb.GetGrantedUserGroupDetailRequest) (*pb.GetGrantedUserGroupDetailResponse, error) {
	ug, err := x.ug.GetGrantedUserGroup(ctx, &biz.ListRoleBindingRequest{
		RoleBindingCommonParams: biz.RoleBindingCommonParams{
			UserGroupID: req.UserGroupId,
		},
		Page:     req.Page,
		PageSize: req.PageSize,
		ListAll:  utils.IntToBool(int(req.ListAll)),
	})
	if err != nil {
		return nil, err
	}
	roles := make([]*pb.RoleDetail, 0)
	for _, role := range ug.RoleDetail {
		roles = append(roles, &pb.RoleDetail{
			RoleId:      role.RoleID,
			ClusterId:   role.ClusterID,
			Namespace:   role.Namespace,
			ClusterName: "", //TODO implement me
			RoleName:    "", //TODO implement me
		})
	}
	return &pb.GetGrantedUserGroupDetailResponse{
		UserGroup: &pb.UserGroup{
			UserGroupId:   ug.UserGroupID,
			UserGroupName: ug.Name,
		},
		RoleDetail: roles,
	}, nil
}

// ListUserGroup 用户组列表
func (x *UserGroupService) ListUserGroup(ctx context.Context, req *pb.ListUserGroupRequest) (*pb.ListUserGroupResponse, error) {
	userGroupList, total, err := x.ug.ListUserGroups(ctx, &biz.ListUserGroupRequest{
		Page:     req.Page,
		PageSize: req.PageSize,
		Keyword:  req.Keyword,
		ListAll:  utils.IntToBool(int(req.ListAll)),
	})
	if err != nil {
		return nil, err
	}
	respList := make([]*pb.UserGroup, 0)
	for _, userGroup := range userGroupList {
		respList = append(respList, &pb.UserGroup{
			UserGroupId:   uint32(userGroup.ID),
			UserGroupName: userGroup.RoleName,
		})
	}
	return &pb.ListUserGroupResponse{
		List:  respList,
		Total: total,
	}, nil
}

// ListUserGroupUsers 用户组用户列表
func (x *UserGroupService) ListUserGroupUsers(ctx context.Context, req *pb.ListUserGroupUsersRequest) (*pb.ListUserGroupUsersResponse, error) {
	userList, total, err := x.ug.GetUsersByUserGroupID(ctx, req.UserGroupId)
	if err != nil {
		return nil, err
	}
	respList := make([]*pb.UserGroupUser, 0)
	for _, user := range userList {
		respList = append(respList, &pb.UserGroupUser{
			UserId:   user.UserID,
			UserName: user.Username,
			Email:    user.Email,
			NickName: user.Nickname,
			Source:   user.Source},
		)
	}
	return &pb.ListUserGroupUsersResponse{
		List:  respList,
		Total: total,
	}, nil
}

// ListUser 用户列表
func (x *UserGroupService) ListUser(ctx context.Context, req *pb.ListUserRequest) (*pb.ListUserResponse, error) {
	userList, total, err := x.ug.ListUsers(ctx, &biz.ListUserRequest{
		Page:     req.Page,
		PageSize: req.PageSize,
		Keyword:  req.Keyword,
		ListAll:  utils.IntToBool(int(req.ListAll)),
	})
	if err != nil {
		return nil, err
	}
	respList := make([]*pb.User, 0)
	for _, user := range userList {
		respList = append(respList, &pb.User{
			Username: user.Username,
			UserId:   strconv.Itoa(int(user.UserID)),
			Nickname: user.Nickname,
		},
		)
	}
	return &pb.ListUserResponse{
		List:  respList,
		Total: total,
	}, nil
}
