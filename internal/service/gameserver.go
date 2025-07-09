package service

import (
	"codo-cnmp/common/utils"
	"codo-cnmp/internal/biz"
	"codo-cnmp/pb"
	"context"
)

type GameServerService struct {
	pb.UnimplementedGameServerServer
	uc *biz.GameServerUseCase
}

func NewGameServerService(uc *biz.GameServerUseCase) *GameServerService {
	return &GameServerService{
		uc: uc,
	}
}

func (x *GameServerService) ListGameServerType(ctx context.Context, req *pb.ListGameServerTypeRequest) (*pb.ListGameServerTypeResponse, error) {
	types, total, err := x.uc.ListGameServeType(ctx, &biz.ListGameServerTypeRequest{
		Page:        req.Page,
		PageSize:    req.PageSize,
		Keyword:     req.Keyword,
		ListAll:     utils.IntToBool(int(req.ListAll)),
		ClusterName: req.ClusterName,
		Namespace:   req.Namespace,
	})
	if err != nil {
		return nil, err
	}
	list := make([]*pb.GameServerType, 0)
	for _, t := range types {
		list = append(list, &pb.GameServerType{
			Name: t.Name,
		})
	}
	return &pb.ListGameServerTypeResponse{
		List:  list,
		Total: total,
	}, nil
}

func (x *GameServerService) convertGameServerToDTO(server *biz.GameServer) *pb.GameServerItem {
	return &pb.GameServerItem{
		ServerName:        server.ServerName,
		ServerVersion:     server.ServerVersion,
		Workload:          server.Workload,
		WorkloadType:      server.WorkloadType,
		EntityNum:         server.EntityNum,
		OnlineNum:         server.OnlineNum,
		EntityLockStatus:  pb.EntityLockStatus(server.EntityLockStatus),
		LbLockStatus:      pb.LBLockStatus(server.LbLockStatus),
		PodName:           server.Pod,
		CodeVersionGame:   server.CodeVersionGame,
		CodeVersionConfig: server.CodeVersionConfig,
		CodeVersionScript: server.CodeVersionScript,
		Id:                server.ID,
	}
}

func (x *GameServerService) ListGameServer(ctx context.Context, req *pb.ListGameServerRequest) (*pb.ListGameServerResponse, error) {
	params := &biz.ListGameServerRequest{
		ClusterName:      req.ClusterName,
		Namespace:        req.Namespace,
		Page:             req.Page,
		PageSize:         req.PageSize,
		Keyword:          req.Keyword,
		ListAll:          utils.IntToBool(int(req.ListAll)),
		LbLockStatus:     uint32(req.LbLockStatus),
		EntityLockStatus: uint32(req.EntityLockStatus),
		ServerType:       req.ServerType,
	}
	servers, total, err := x.uc.ListGameServer(ctx, params)
	if err != nil {
		return nil, err
	}
	list := make([]*pb.GameServerItem, 0)
	for _, t := range servers {
		list = append(list, x.convertGameServerToDTO(t))
	}
	return &pb.ListGameServerResponse{
		List:  list,
		Total: total,
	}, nil

}

func (x *GameServerService) ManageGameServerEntity(ctx context.Context, req *pb.ManageGameServerEntityRequest) (*pb.ManageGameServerEntityResponse, error) {
	result, err := x.uc.ManageEntity(ctx, &biz.ManageEntityRequest{
		ClusterName: req.ClusterName,
		Namespace:   req.Namespace,
		ServerName:  req.ServerName,
		Lock:        req.Lock,
	})
	if err != nil {
		return nil, err
	}
	var (
		head *pb.ManageEntityHead
		body *pb.ManageEntityBody
	)
	if result != nil && &result.Body != nil && &result.Head != nil {
		head = &pb.ManageEntityHead{
			Errmsg: result.Head.Errmsg,
			Errno:  uint32(result.Head.Errno),
		}
		body = &pb.ManageEntityBody{}
	}

	return &pb.ManageGameServerEntityResponse{
		Head: head,
		Body: body,
	}, nil
}

func (x *GameServerService) ManageGameServerLB(ctx context.Context, req *pb.ManageGameServerEntityRequest) (*pb.ManageGameServerEntityResponse, error) {
	result, err := x.uc.ManageLB(ctx, &biz.ManageLBRequest{
		ClusterName: req.ClusterName,
		Namespace:   req.Namespace,
		ServerName:  req.ServerName,
		Lock:        req.Lock,
	})
	if err != nil {
		return nil, err
	}
	var (
		head *pb.ManageEntityHead
		body *pb.ManageEntityBody
	)
	if result != nil && &result.Body != nil && &result.Head != nil {
		head = &pb.ManageEntityHead{
			Errmsg: result.Head.Errmsg,
			Errno:  uint32(result.Head.Errno),
		}
		body = &pb.ManageEntityBody{}
	}

	return &pb.ManageGameServerEntityResponse{
		Head: head,
		Body: body,
	}, nil
}

func (x *GameServerService) BatchManageGameServerEntity(ctx context.Context, req *pb.BatchManageGameServerEntityRequest) (*pb.BatchManageGameServerEntityResponse, error) {
	err := x.uc.BatchManageEntity(ctx, &biz.BatchManageEntityRequest{
		ClusterName: req.ClusterName,
		Namespace:   req.Namespace,
		ServerNames: req.ServerNames,
		Lock:        req.Lock,
	})
	if err != nil {
		return nil, err
	}
	return &pb.BatchManageGameServerEntityResponse{}, nil
}

func (x *GameServerService) BatchManageGameServerLB(ctx context.Context, req *pb.BatchManageGameServerEntityRequest) (*pb.BatchManageGameServerEntityResponse, error) {
	err := x.uc.BatchManageLB(ctx, &biz.BatchManageLBRequest{
		ClusterName: req.ClusterName,
		Namespace:   req.Namespace,
		ServerNames: req.ServerNames,
		Lock:        req.Lock,
	})
	if err != nil {
		return nil, err
	}
	return &pb.BatchManageGameServerEntityResponse{}, nil
}
