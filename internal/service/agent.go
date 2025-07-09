package service

import (
	"codo-cnmp/common/utils"
	"codo-cnmp/internal/biz"
	"codo-cnmp/pb"
	"context"
)

type AgentService struct {
	uc *biz.AgentUseCase
}

func (x *AgentService) ListAgent(ctx context.Context, request *pb.ListAgentRequest) (*pb.ListAgentResponse, error) {
	agents, total, err := x.uc.ListAgent(ctx, &biz.ListAgentRequest{
		Keyword:  request.GetKeyword(),
		Page:     request.GetPage(),
		PageSize: request.GetPageSize(),
		ListAll:  utils.IntToBool(int(request.GetListAll())),
	})
	if err != nil {
		return nil, err
	}
	list := make([]*pb.AgentItem, 0, len(agents))
	for _, agent := range agents {
		item := x.convertDO2DTO(agent)
		list = append(list, item)
	}
	return &pb.ListAgentResponse{
		List:  list,
		Total: total,
	}, nil

}

func (x *AgentService) CreateAgent(ctx context.Context, request *pb.CreateAgentRequest) (*pb.CreateAgentResponse, error) {
	err := x.uc.CreateAgent(ctx, &biz.CreateAgentRequest{
		AgentCommonParams: biz.AgentCommonParams{
			Name:    request.GetName(),
			AgentId: request.GetAgentId(),
		},
	})
	if err != nil {
		return nil, err
	}
	return &pb.CreateAgentResponse{}, nil
}

func (x *AgentService) DeleteAgent(ctx context.Context, request *pb.DeleteAgentRequest) (*pb.DeleteAgentResponse, error) {
	err := x.uc.DeleteAgent(ctx, &biz.DeleteAgentRequest{Id: int(request.Id)})
	if err != nil {
		return nil, err
	}
	return &pb.DeleteAgentResponse{}, nil
}

func (x *AgentService) UpdateAgent(ctx context.Context, request *pb.UpdateAgentRequest) (*pb.UpdateAgentResponse, error) {
	err := x.uc.UpdateAgent(ctx, &biz.UpdateAgentRequest{
		Id: int(request.Id),
		AgentCommonParams: biz.AgentCommonParams{
			Name:    request.GetName(),
			AgentId: request.GetAgentId(),
		},
	})
	if err != nil {
		return nil, err
	}
	return &pb.UpdateAgentResponse{}, nil
}

func (x *AgentService) convertDO2DTO(agent *biz.AgentItem) *pb.AgentItem {
	return &pb.AgentItem{
		Id:      uint32(agent.Id),
		AgentId: agent.AgentId,
		Name:    agent.Name,
	}
}

func NewAgentService(uc *biz.AgentUseCase) *AgentService {
	return &AgentService{uc: uc}
}
