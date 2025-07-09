package service

import (
	"codo-cnmp/common/utils"
	"codo-cnmp/internal/biz"
	pb "codo-cnmp/pb"
	"context"
	corev1 "k8s.io/api/core/v1"
)

type EventService struct {
	pb.UnimplementedEventServer
	uc *biz.EventUseCase
}

func NewEventService(uc *biz.EventUseCase) *EventService {
	return &EventService{
		uc: uc,
	}
}

func (s *EventService) convertEvent2DTO(event *corev1.Event) *pb.EventItem {
	var lastTimestamp, firstTimestamp *uint64
	if !event.LastTimestamp.IsZero() {
		tmp := uint64(event.LastTimestamp.UnixNano() / 1e6)
		lastTimestamp = &tmp
	}
	if !event.FirstTimestamp.IsZero() {
		tmp := uint64(event.FirstTimestamp.UnixNano() / 1e6)
		firstTimestamp = &tmp
	}
	return &pb.EventItem{
		FirstTimeStamp:     firstTimestamp,
		LastTimeStamp:      lastTimestamp,
		Reason:             event.Reason,
		Message:            event.Message,
		Type:               event.Type,
		Count:              uint32(event.Count),
		InvolvedObjectKind: event.InvolvedObject.Kind,
		InvolvedObjectName: event.InvolvedObject.Name,
	}
}

func (s *EventService) ListEvent(ctx context.Context, req *pb.ListEventRequest) (*pb.ListEventResponse, error) {
	events, total, err := s.uc.ListEvent(ctx, &biz.ListEventRequest{
		EventCommonParams: biz.EventCommonParams{
			ClusterName:    req.ClusterName,
			Namespace:      req.Namespace,
			ControllerName: req.ControllerName,
			ControllerType: req.ControllerType.String(),
		},
		Page:     req.Page,
		PageSize: req.PageSize,
		Keyword:  req.Keyword,
		ListAll:  utils.IntToBool(int(req.ListAll)),
	})
	if err != nil {
		return nil, err
	}
	list := make([]*pb.EventItem, 0, len(events))
	for _, event := range events {
		list = append(list, s.convertEvent2DTO(event))
	}
	return &pb.ListEventResponse{
		List:  list,
		Total: total,
	}, nil
}
