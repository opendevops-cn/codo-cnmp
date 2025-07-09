package service

import (
	"codo-cnmp/common/utils"
	"codo-cnmp/internal/biz"
	"codo-cnmp/pb"
	"context"
)

type UserFollowService struct {
	pb.UnimplementedUserFollowServer
	uc *biz.UserFollowUseCase
}

func NewUserFollowService(uc *biz.UserFollowUseCase) *UserFollowService {
	return &UserFollowService{uc: uc}
}

func (x *UserFollowService) CreateUserFollow(ctx context.Context, req *pb.UserFollowRequest) (*pb.UserFollowResponse, error) {
	userId, err := utils.GetUserIDFromCtx(ctx)
	if err != nil {
		return nil, err
	}
	err = x.uc.CreateUserFollow(ctx, &biz.UserFollowRequest{
		UserFollowCommonParams: biz.UserFollowCommonParams{
			UserID:      userId,
			FollowType:  req.FollowType,
			FollowValue: req.FollowValue,
			ClusterName: req.ClusterName,
		},
	})
	if err != nil {
		return nil, err
	}
	return &pb.UserFollowResponse{}, nil
}

func (x *UserFollowService) DeleteUserFollow(ctx context.Context, req *pb.DeleteUserFollowRequest) (*pb.DeleteUserFollowResponse, error) {
	userId, err := utils.GetUserIDFromCtx(ctx)
	if err != nil {
		return nil, err
	}
	err = x.uc.DeleteUserFollow(ctx, &biz.DeleteUserFollowRequest{
		UserFollowCommonParams: biz.UserFollowCommonParams{
			UserID:      userId,
			FollowType:  req.FollowType,
			FollowValue: req.FollowValue,
			ClusterName: req.ClusterName,
		},
	})
	if err != nil {
		return nil, err
	}
	return &pb.DeleteUserFollowResponse{}, nil
}

func (x *UserFollowService) ListUserFollow(ctx context.Context, req *pb.ListUserFollowRequest) (*pb.ListUserFollowResponse, error) {
	userId, err := utils.GetUserIDFromCtx(ctx)
	if err != nil {
		return nil, err
	}
	userFlows, count, err := x.uc.ListUserFollow(ctx, &biz.ListUserFollowRequest{
		UserFollowCommonParams: biz.UserFollowCommonParams{
			UserID:      userId,
			FollowType:  req.FollowType,
			FollowValue: req.FollowValue,
		},
		Page:     req.Page,
		PageSize: req.PageSize,
		Keyword:  req.Keyword,
		ListAll:  utils.IntToBool(int(req.ListAll)),
	})
	if err != nil {
		return nil, err
	}
	list := make([]*pb.FollowItem, 0, len(userFlows))
	for _, userFlow := range userFlows {
		list = append(list, &pb.FollowItem{
			Id:          userFlow.ID,
			FollowType:  userFlow.FollowType,
			FollowValue: userFlow.FollowValue,
			CreateTime:  userFlow.CreatedTime,
			ClusterName: userFlow.ClusterName,
		})
	}
	return &pb.ListUserFollowResponse{
		List:  list,
		Total: count,
	}, nil
}
