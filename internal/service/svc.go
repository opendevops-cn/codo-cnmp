package service

import (
	"context"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"codo-cnmp/common/utils"
	"codo-cnmp/internal/biz"
	"codo-cnmp/pb"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

type SvcService struct {
	pb.UnimplementedSVCServer
	uc *biz.SvcUseCase
	uf *biz.UserFollowUseCase
}

func (x *SvcService) PrepareData(ctx context.Context, request *pb.CreateSvcRequest) (*biz.CreateSvcRequest, error) {
	req := &biz.CreateSvcRequest{
		SvcCommonParams: biz.SvcCommonParams{
			ClusterName: request.ClusterName,
			Namespace:   request.Namespace,
			Name:        request.Name,
		},
		SvcType: request.SvcType,
	}
	var ports []corev1.ServicePort
	if request.Ports != nil {
		for _, port := range request.Ports {
			servicePort := corev1.ServicePort{}

			if port.Name != nil {
				servicePort.Name = *port.Name
			}

			if port.Port != nil {
				servicePort.Port = *port.Port
			}

			if port.TargetPort != nil {
				if targetPortInt, err := strconv.Atoi(*port.TargetPort); err == nil {
					servicePort.TargetPort = intstr.FromInt32(int32(targetPortInt))
				} else {
					servicePort.TargetPort = intstr.FromString(*port.TargetPort)
				}
			}

			if port.Protocol != nil {
				servicePort.Protocol = corev1.Protocol(*port.Protocol)
			}

			if port.AppProtocol != nil {
				servicePort.AppProtocol = port.AppProtocol
			}

			if port.NodePort != nil {
				servicePort.NodePort = *port.NodePort
			}
			ports = append(ports, servicePort)
		}
		req.Ports = ports
	}
	if len(request.ExternalIps) > 0 {
		req.ExternalIPs = request.ExternalIps
	}
	if request.PublishNotReadyAddresses != nil {
		req.PublishNotReadyAddresses = *request.PublishNotReadyAddresses
	}
	if request.SessionAffinity != nil && *request.SessionAffinity != pb.SessionAffinity_SESSION_AFFINITY_UNSPECIFIED {
		req.SessionAffinity = request.SessionAffinity.String()
	}
	if request.SessionAffinitySeconds != nil {
		req.SessionAffinitySeconds = *request.SessionAffinitySeconds
	}
	if request.Headless != nil {
		req.Headless = *request.Headless
	}
	if request.ExternalName != nil {
		req.ExternalName = *request.ExternalName
	}
	if request.Labels != nil {
		req.Labels = request.Labels
	}
	if request.Annotations != nil {
		req.Annotations = request.Annotations
	}
	if request.Selector != nil {
		req.Selector = request.Selector
	}
	return req, nil
}

func (x *SvcService) CreateSvc(ctx context.Context, request *pb.CreateSvcRequest) (*pb.CreateSvcResponse, error) {
	req, _ := x.PrepareData(ctx, request)
	err := x.uc.CreateSvc(ctx, req)
	if err != nil {
		return nil, err
	}
	return &pb.CreateSvcResponse{}, nil
}

func (x *SvcService) UpdateSvc(ctx context.Context, request *pb.CreateSvcRequest) (*pb.CreateSvcResponse, error) {
	req, _ := x.PrepareData(ctx, request)
	err := x.uc.UpdateSvc(ctx, req)
	if err != nil {
		return nil, err
	}
	return &pb.CreateSvcResponse{}, nil
}

// getUserFollowMap 获取用户关注的Service列表
func (x *SvcService) getUserFollowMap(ctx context.Context, userID uint32) (map[string]bool, error) {
	userFollows, _, err := x.uf.ListUserFollow(ctx, &biz.ListUserFollowRequest{
		UserFollowCommonParams: biz.UserFollowCommonParams{
			UserID:     userID,
			FollowType: pb.FollowType_Svc,
		},
		ListAll: true,
	})
	if err != nil {
		return nil, err
	}

	followMap := make(map[string]bool)
	for _, follow := range userFollows {
		followKey := fmt.Sprintf("%s.%s", follow.ClusterName, follow.FollowValue)
		followMap[followKey] = true
	}
	return followMap, nil
}

// setFollowedStatus 设置关注状态
func (x *SvcService) setFollowedStatus(ctx context.Context, clusterName string, items []*pb.SvcItem) error {
	userID, err := utils.GetUserIDFromCtx(ctx)
	if err != nil {
		return err
	}
	followMap, err := x.getUserFollowMap(ctx, userID)
	if err != nil {
		return err
	}
	for _, item := range items {
		followKey := fmt.Sprintf("%s.%s", clusterName, item.Name)
		item.IsFollowed = followMap[followKey]
	}
	return nil
}

func (x *SvcService) convertDO2DTO(svc *corev1.Service) *pb.SvcItem {
	var (
		svcType pb.SvcType
	)
	switch svc.Spec.Type {
	case corev1.ServiceTypeClusterIP:
		svcType = pb.SvcType_ClusterIP
	case corev1.ServiceTypeLoadBalancer:
		svcType = pb.SvcType_LoadBalancer
	case corev1.ServiceTypeNodePort:
		svcType = pb.SvcType_NodePort
	case corev1.ServiceTypeExternalName:
		svcType = pb.SvcType_ExternalName
	default:
		svcType = pb.SvcType_SVC_TYPE_UNSPECIFIED
	}

	var vip []string
	if svc.Status.LoadBalancer.Ingress != nil {
		if len(svc.Status.LoadBalancer.Ingress) > 0 {
			for _, i := range svc.Status.LoadBalancer.Ingress {
				if i.IP != "" {
					vip = append(vip, i.IP)
				}
				if i.Hostname != "" {
					vip = append(vip, i.Hostname)
				}
			}
		}
	}
	var (
		createTime, updateTime time.Time
	)
	createTime = svc.CreationTimestamp.Time
	for _, managedField := range svc.ManagedFields {
		if managedField.Operation == metav1.ManagedFieldsOperationUpdate {
			if updateTime.IsZero() || managedField.Time.Time.After(updateTime) {
				updateTime = managedField.Time.Time
			}
		}
	}

	// 如果没有更新记录，则更新时间等于创建时间
	if updateTime.IsZero() {
		updateTime = createTime
	}

	var ports []*pb.ServicePort
	for _, port := range svc.Spec.Ports {
		targetPort := port.TargetPort.String()
		protocol := string(port.Protocol)
		servicePort := &pb.ServicePort{
			Protocol:    &protocol,
			AppProtocol: port.AppProtocol,
			TargetPort:  &targetPort,
		}
		if port.NodePort != 0 {
			servicePort.NodePort = &port.NodePort
		}
		if port.Port != 0 {
			servicePort.Port = &port.Port
		}
		if port.Name != "" {
			servicePort.Name = &port.Name
		}

		ports = append(ports, servicePort)
	}

	if svc.Kind == "" {
		svc.Kind = "Service"
	}
	if svc.APIVersion == "" {
		svc.APIVersion = "v1"
	}
	yamlStr, _ := utils.ConvertResourceTOYaml(svc)

	var sessionAffinity pb.SessionAffinity
	switch svc.Spec.SessionAffinity {
	case corev1.ServiceAffinityClientIP:
		sessionAffinity = pb.SessionAffinity_ClientIP
	case corev1.ServiceAffinityNone:
		sessionAffinity = pb.SessionAffinity_None
	default:
		sessionAffinity = pb.SessionAffinity_SESSION_AFFINITY_UNSPECIFIED
	}

	var headless bool
	if svc.Spec.ClusterIP == "None" {
		headless = true
	}
	var sessionAffinitySeconds int32
	if svc.Spec.SessionAffinityConfig != nil && svc.Spec.SessionAffinityConfig.ClientIP != nil {
		sessionAffinitySeconds = *svc.Spec.SessionAffinityConfig.ClientIP.TimeoutSeconds
	}
	uint32SessionAffinitySeconds := uint32(sessionAffinitySeconds)

	publishNotReadyAddresses := svc.Spec.PublishNotReadyAddresses

	return &pb.SvcItem{
		Name:                     svc.Name,
		Namespace:                svc.Namespace,
		SvcType:                  svcType,
		Ports:                    ports,
		ClusterIp:                svc.Spec.ClusterIP,
		CreateTime:               uint64(createTime.UnixNano() / 1e6),
		UpdateTime:               uint64(updateTime.UnixNano() / 1e6),
		Yaml:                     yamlStr,
		Vip:                      strings.Join(vip, ","),
		Labels:                   svc.Labels,
		Annotations:              svc.Annotations,
		Selector:                 svc.Spec.Selector,
		SessionAffinity:          sessionAffinity,
		ExternalName:             svc.Spec.ExternalName,
		Headless:                 headless,
		ExternalIps:              svc.Spec.ExternalIPs,
		SessionAffinitySeconds:   &uint32SessionAffinitySeconds,
		PublishNotReadyAddresses: &publishNotReadyAddresses,
	}
}
func (x *SvcService) ListSvc(ctx context.Context, request *pb.ListSvcRequest) (*pb.ListSvcResponse, error) {
	svcs, total, err := x.uc.ListSvc(ctx, &biz.ListSvcRequest{
		ClusterName: request.ClusterName,
		Keyword:     request.Keyword,
		ListAll:     utils.IntToBool(int(request.ListAll)),
		Namespace:   request.Namespace,
		Page:        request.Page,
		PageSize:    request.PageSize,
		SvcType:     request.GetSvcType(),
	})
	if err != nil {
		return nil, err
	}
	list := make([]*pb.SvcItem, 0, len(svcs))
	for _, svc := range svcs {
		dto := x.convertDO2DTO(svc)
		refCount, refs, err := x.uc.GetSvcReferences(ctx, svc, &biz.SvcCommonParams{
			ClusterName: request.ClusterName,
			Namespace:   svc.Namespace,
			Name:        svc.Name,
		})
		if err != nil {
			continue
		}
		for _, ref := range refs {
			dto.Refs = append(dto.Refs, &pb.SvcReference{
				Kind: ref["kind"],
				Name: ref["name"],
			})

		}
		dto.RefCount = uint32(refCount)
		list = append(list, dto)
	}
	if err := x.setFollowedStatus(ctx, request.ClusterName, list); err != nil {
		return nil, err
	}
	sort.SliceStable(list, func(i, j int) bool {
		return list[i].IsFollowed && !list[j].IsFollowed
	})
	return &pb.ListSvcResponse{
		List:  list,
		Total: total,
	}, nil
}

func (x *SvcService) DeleteSvc(ctx context.Context, request *pb.DeleteSvcRequest) (*pb.DeleteSvcResponse, error) {
	err := x.uc.DeleteSvc(ctx, &biz.DeleteSvcRequest{
		SvcCommonParams: biz.SvcCommonParams{
			ClusterName: request.ClusterName,
			Name:        request.Name,
			Namespace:   request.Namespace,
		},
	})
	if err != nil {
		return nil, err
	}
	return &pb.DeleteSvcResponse{}, nil
}

func (x *SvcService) GetSvcDetail(ctx context.Context, request *pb.SvcDetailRequest) (*pb.SvcDetailResponse, error) {
	svc, err := x.uc.GetSvc(ctx, &biz.GetSvcRequest{
		SvcCommonParams: biz.SvcCommonParams{
			ClusterName: request.ClusterName,
			Name:        request.Name,
			Namespace:   request.Namespace,
		},
	})
	if err != nil {
		return nil, err
	}
	detail := x.convertDO2DTO(svc)
	refCount, refs, err := x.uc.GetSvcReferences(ctx, svc, &biz.SvcCommonParams{
		ClusterName: request.ClusterName,
		Namespace:   svc.Namespace,
		Name:        svc.Name,
	})
	if err != nil {
		return nil, err
	}
	for _, ref := range refs {
		detail.Refs = append(detail.Refs, &pb.SvcReference{
			Kind: ref["kind"],
			Name: ref["name"],
		})

	}
	detail.RefCount = uint32(refCount)
	return &pb.SvcDetailResponse{
		Detail: x.convertDO2DTO(svc),
	}, nil
}

func NewSvcService(uc *biz.SvcUseCase, uf *biz.UserFollowUseCase) *SvcService {
	return &SvcService{
		uc: uc,
		uf: uf,
	}
}
