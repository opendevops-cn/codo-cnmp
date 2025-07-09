package service

import (
	"codo-cnmp/internal/biz"
	"codo-cnmp/pb"
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/openkruise/kruise-api/apps/v1alpha1"
	"strings"
)

type CRRService struct {
	pb.UnimplementedCRRServer
	uc    *biz.CRRUseCase
	redis *redis.Client
}

func (x *CRRService) BatchCreateCrr(ctx context.Context, request *pb.BatchCreateCRRRequest) (*pb.BatchCreateCRRResponse, error) {
	containerItems := make([]biz.ContainerItems, 0)
	for _, container := range request.ContainerItems {
		containerItems = append(containerItems, biz.ContainerItems{
			PodName:    container.PodName,
			Containers: container.GetContainerNames(),
		})
	}
	names, err := x.uc.BatchCreateCRR(ctx, &biz.BatchCreateCRRRequest{
		CRRCommonParams: biz.CRRCommonParams{
			ClusterName: request.ClusterName,
			Namespace:   request.Namespace,
		},
		ContainerItems: containerItems,
	})
	if err != nil {
		return nil, err
	}
	list := make([]*pb.ResponseContainerItem, 0)
	for _, name := range names {
		podName := strings.Split(name, ".")[1]
		list = append(list, &pb.ResponseContainerItem{Name: name, PodName: podName})
	}
	return &pb.BatchCreateCRRResponse{List: list}, nil
}

func (x *CRRService) convertCRR(crr *v1alpha1.ContainerRecreateRequest) *pb.GetCRRDetailResponse {
	var successCount uint32
	for _, state := range crr.Status.ContainerRecreateStates {
		if state.Phase == v1alpha1.ContainerRecreateRequestSucceeded {
			successCount++
		}
	}
	return &pb.GetCRRDetailResponse{
		PodName:      strings.Split(crr.Name, ".")[1],
		TotalCount:   uint32(len(crr.Status.ContainerRecreateStates)),
		SuccessCount: successCount,
	}
}

func (x *CRRService) GetCRRDetail(ctx context.Context, request *pb.GetCRRDetailRequest) (*pb.GetCRRDetailResponse, error) {
	crr, err := x.uc.GetCRR(ctx, &biz.GetCRRRequest{
		CRRCommonParams: biz.CRRCommonParams{
			ClusterName: request.ClusterName,
			Namespace:   request.Namespace,
		},
		CRRName: request.Name,
	})
	if err != nil {
		return nil, err
	}
	return x.convertCRR(crr), nil
}

func (x *CRRService) BatchQueryCRR(ctx context.Context, request *pb.GetBatchCRRDetailRequest) (*pb.GetBatchCRRDetailResponse, error) {
	crrs, err := x.uc.BatchQueryCRR(ctx, &biz.BatchQueryCRRDetailRequest{
		CRRCommonParams: biz.CRRCommonParams{
			ClusterName: request.ClusterName,
			Namespace:   request.Namespace,
		},
		CRRNames: request.GetNames(),
	})
	if err != nil {
		return nil, err
	}
	list := make([]*pb.BatchCRRDetail, 0)
	for _, crr := range crrs {
		var successCount uint32
		for _, state := range crr.Status.ContainerRecreateStates {
			if state.Phase == v1alpha1.ContainerRecreateRequestSucceeded {
				successCount++
			}
		}
		cachedCRR := fmt.Sprintf("codo:cnmp:%s:crr:%s", request.Namespace, crr.Name)
		totalCount, _ := x.redis.Get(ctx, cachedCRR).Int()
		list = append(list, &pb.BatchCRRDetail{
			Name:         crr.Name,
			PodName:      strings.Split(crr.Name, ".")[1],
			TotalCount:   uint32(totalCount),
			SuccessCount: successCount,
		})
	}
	return &pb.GetBatchCRRDetailResponse{List: list}, nil
}

func (x *CRRService) CreateCrr(ctx context.Context, request *pb.CreateCRRRequest) (*pb.CreateCRRResponse, error) {
	name, err := x.uc.CreateCRR(ctx, &biz.CreateCRRRequest{
		CRRCommonParams: biz.CRRCommonParams{
			ClusterName: request.ClusterName,
			Namespace:   request.Namespace,
		},
		PodName:    request.PodName,
		Containers: request.GetContainerNames(),
	})
	if err != nil {
		return nil, err
	}
	return &pb.CreateCRRResponse{Item: &pb.ResponseContainerItem{Name: name, PodName: request.PodName}}, nil
}

func NewCRRService(uc *biz.CRRUseCase, redis *redis.Client) *CRRService {
	return &CRRService{uc: uc, redis: redis}
}
