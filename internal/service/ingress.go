package service

import (
	"context"
	"fmt"
	"sort"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"codo-cnmp/common/utils"
	"codo-cnmp/internal/biz"
	"codo-cnmp/pb"
	networkingv1 "k8s.io/api/networking/v1"
	"sigs.k8s.io/yaml"
)

type IngressService struct {
	pb.UnimplementedIngressServer
	uc *biz.IngressUseCase
	uf *biz.UserFollowUseCase
}

func (x *IngressService) convertVO2DO(request *pb.CreateIngressRequest) *biz.CreateIngressRequest {
	ingressReq := &biz.CreateIngressRequest{
		IngressCommonParams: biz.IngressCommonParams{
			ClusterName: request.ClusterName,
			Namespace:   request.Namespace,
			Name:        request.Name,
		},
	}
	ingressSpec := &networkingv1.IngressSpec{}
	if request.Tls != nil {
		for _, v := range request.Tls {
			ingressSpec.TLS = append(ingressSpec.TLS, networkingv1.IngressTLS{
				Hosts:      v.Hosts,
				SecretName: v.SecretName,
			})

		}
	}
	if request.IngressRules != nil {
		for _, rule := range request.IngressRules {
			// 为每个规则创建一个IngressRule对象
			ingressRule := networkingv1.IngressRule{
				Host: rule.Host,
				IngressRuleValue: networkingv1.IngressRuleValue{
					HTTP: &networkingv1.HTTPIngressRuleValue{
						Paths: []networkingv1.HTTPIngressPath{},
					},
				},
			}

			// 添加该规则的所有路径
			for _, path := range rule.IngressRuleValue.Http.Paths {
				pathType := networkingv1.PathType(path.PathType)
				// 确定使用端口数字还是端口名称
				var backendPort networkingv1.ServiceBackendPort
				if path.Backend.Service.Port.Name != "" {
					// 如果端口名称不为空，则优先使用端口名称
					backendPort = networkingv1.ServiceBackendPort{
						Name: path.Backend.Service.Port.Name,
					}
				} else if path.Backend.Service.Port.Number != 0 {
					// 如果端口名称为空，且端口数字不为空，则使用端口数字
					backendPort = networkingv1.ServiceBackendPort{
						Number: int32(path.Backend.Service.Port.Number),
					}
				}

				ingressRule.IngressRuleValue.HTTP.Paths = append(ingressRule.IngressRuleValue.HTTP.Paths,
					networkingv1.HTTPIngressPath{
						Path:     path.Path,
						PathType: &pathType,
						Backend: networkingv1.IngressBackend{
							Service: &networkingv1.IngressServiceBackend{
								Name: path.Backend.Service.Name,
								Port: backendPort,
							},
						},
					})
			}

			// 将完整的规则添加到规则列表
			ingressSpec.Rules = append(ingressSpec.Rules, ingressRule)
		}
	}
	if request.IngressClassName != "" {
		ingressSpec.IngressClassName = &request.IngressClassName
	}
	ingressReq.Annotations = request.Annotations
	ingressReq.Labels = request.Labels
	ingressReq.IngressSpec = *ingressSpec
	return ingressReq

}

// CreateIngress implements pb.IngressHTTPServer.
func (x *IngressService) CreateIngress(ctx context.Context, request *pb.CreateIngressRequest) (*pb.CreateIngressResponse, error) {
	req := x.convertVO2DO(request)
	err := x.uc.CreateIngress(ctx, req)
	if err != nil {
		return nil, err
	}
	return &pb.CreateIngressResponse{}, nil
}

// UpdateIngress implements pb.IngressHTTPServer.
func (x *IngressService) UpdateIngress(ctx context.Context, request *pb.CreateIngressRequest) (*pb.CreateIngressResponse, error) {
	req := x.convertVO2DO(request)
	err := x.uc.UpdateIngress(ctx, req)
	if err != nil {
		return nil, err
	}
	return &pb.CreateIngressResponse{}, nil
}

func NewIngressService(uc *biz.IngressUseCase, uf *biz.UserFollowUseCase) *IngressService {
	return &IngressService{uc: uc, uf: uf}
}

// getUserFollowMap 获取用户关注的Ingress列表
func (x *IngressService) getUserFollowMap(ctx context.Context, userID uint32) (map[string]bool, error) {
	userFollows, _, err := x.uf.ListUserFollow(ctx, &biz.ListUserFollowRequest{
		UserFollowCommonParams: biz.UserFollowCommonParams{
			UserID:     userID,
			FollowType: pb.FollowType_Ingress,
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
func (x *IngressService) setFollowedStatus(ctx context.Context, clusterName string, items []*pb.IngressItem) error {
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

func (x *IngressService) ListIngress(ctx context.Context, request *pb.ListIngressRequest) (*pb.ListIngressResponse, error) {
	ingresses, total, err := x.uc.ListIngress(ctx, &biz.ListIngressRequest{
		ClusterName: request.ClusterName,
		Keyword:     request.Keyword,
		ListAll:     utils.IntToBool(int(request.ListAll)),
		Namespace:   request.Namespace,
		Page:        request.Page,
		PageSize:    request.PageSize,
		Host:        request.Host,
	})
	if err != nil {
		return nil, err
	}
	list := make([]*pb.IngressItem, 0, len(ingresses))
	for _, ingress := range ingresses {
		list = append(list, x.convertDO2DTO(ingress))
	}
	if err := x.setFollowedStatus(ctx, request.ClusterName, list); err != nil {
		return nil, err
	}
	sort.SliceStable(list, func(i, j int) bool {
		return list[i].IsFollowed && !list[j].IsFollowed
	})
	return &pb.ListIngressResponse{
		List:  list,
		Total: total,
	}, nil
}

func (x *IngressService) DeleteIngress(ctx context.Context, request *pb.DeleteIngressRequest) (*pb.DeleteIngressResponse, error) {
	err := x.uc.DeleteIngress(ctx, &biz.DeleteIngressRequest{
		IngressCommonParams: biz.IngressCommonParams{
			ClusterName: request.ClusterName,
			Name:        request.Name,
			Namespace:   request.Namespace,
		},
	})
	if err != nil {
		return nil, err
	}
	return &pb.DeleteIngressResponse{}, nil
}

func (x *IngressService) GetIngressDetail(ctx context.Context, request *pb.IngressDetailRequest) (*pb.IngressDetailResponse, error) {
	ingress, err := x.uc.GetIngress(ctx, &biz.GetIngressRequest{
		IngressCommonParams: biz.IngressCommonParams{
			ClusterName: request.ClusterName,
			Name:        request.Name,
			Namespace:   request.Namespace,
		},
	})
	if err != nil {
		return nil, err
	}
	return &pb.IngressDetailResponse{
		Detail: x.convertDO2DTO(ingress),
	}, nil
}

func (x *IngressService) convertDO2DTO(ingress *networkingv1.Ingress) *pb.IngressItem {
	if ingress.APIVersion == "" {
		ingress.APIVersion = "networking.k8s.io/v1"
	}
	if ingress.Kind == "" {
		ingress.Kind = "Ingress"
	}
	yamlStr, _ := utils.ConvertResourceTOYaml(ingress)
	var ingressClassName string
	if ingress.Spec.IngressClassName != nil {
		ingressClassName = *ingress.Spec.IngressClassName
	}
	var vip string
	if len(ingress.Status.LoadBalancer.Ingress) > 0 {
		vip = ingress.Status.LoadBalancer.Ingress[0].IP
	}
	var ingressRulesStr string
	if ingress.Spec.Rules != nil {
		ingressRuleBytes, _ := yaml.Marshal(ingress.Spec.Rules)
		ingressRulesStr = string(ingressRuleBytes)
	}
	var defaultBackend *pb.IngressBackend
	if ingress.Spec.DefaultBackend != nil {
		defaultBackend = x.convertIngressBackendDO2DTO(*ingress.Spec.DefaultBackend)
	}

	var updateTime time.Time
	for _, managedField := range ingress.ManagedFields {
		if managedField.Operation == metav1.ManagedFieldsOperationUpdate {
			if updateTime.IsZero() || managedField.Time.Time.After(updateTime) {
				updateTime = managedField.Time.Time
			}
		}
	}

	var tls []*pb.IngressTLS
	if ingress.Spec.TLS != nil {
		tls = x.convertIngressTLSDO2DTO(ingress.Spec.TLS)
	}

	return &pb.IngressItem{
		Name:             ingress.Name,
		Namespace:        ingress.Namespace,
		IngressRules:     x.convertIngressRulesDO2DTO(ingress.Spec.Rules),
		CreateTime:       uint64(ingress.CreationTimestamp.UnixNano() / 1e6),
		IngressClassName: ingressClassName,
		Vip:              vip,
		Yaml:             yamlStr,
		IngressRulesStr:  ingressRulesStr,
		UpdateTime:       uint64(updateTime.UnixNano() / 1e6),
		Labels:           ingress.Labels,
		Annotations:      ingress.Annotations,
		DefaultBackend:   defaultBackend,
		Tls:              tls,
	}
}

func (x *IngressService) convertIngressTLSDO2DTO(ingressTLS []networkingv1.IngressTLS) []*pb.IngressTLS {
	var tlsList []*pb.IngressTLS
	for _, tls := range ingressTLS {
		tlsList = append(tlsList, &pb.IngressTLS{
			Hosts:      tls.Hosts,
			SecretName: tls.SecretName,
		})
	}
	return tlsList
}

func (x *IngressService) convertIngressRulesDO2DTO(ingressRules []networkingv1.IngressRule) []*pb.IngressRule {
	var ingressRulesDTO []*pb.IngressRule
	for _, ingressRule := range ingressRules {
		ingressRulesDTO = append(ingressRulesDTO, &pb.IngressRule{
			Host:             ingressRule.Host,
			IngressRuleValue: x.convertIngressRuleValueDO2DTO(ingressRule.IngressRuleValue),
		})
	}
	return ingressRulesDTO
}

func (x *IngressService) convertIngressRuleValueDO2DTO(ingressRuleValue networkingv1.IngressRuleValue) *pb.IngressRuleValue {
	if ingressRuleValue.HTTP != nil {
		return &pb.IngressRuleValue{
			Http: x.convertHTTPIngressRuleValueDO2DTO(*ingressRuleValue.HTTP),
		}
	}
	return &pb.IngressRuleValue{}
}

func (x *IngressService) convertHTTPIngressRuleValueDO2DTO(httpIngressRuleValue networkingv1.HTTPIngressRuleValue) *pb.HTTPIngressRuleValue {
	HTTPIngressPath := make([]*pb.HTTPIngressPath, len(httpIngressRuleValue.Paths))
	for i, ingressPath := range httpIngressRuleValue.Paths {
		var pathType string
		switch *ingressPath.PathType {
		case networkingv1.PathTypeExact:
			pathType = "Exact"
		case networkingv1.PathTypeImplementationSpecific:
			pathType = "ImplementationSpecific"
		case networkingv1.PathTypePrefix:
			pathType = "Prefix"
		default:
			pathType = "ImplementationSpecific"
		}
		HTTPIngressPath[i] = &pb.HTTPIngressPath{
			Path:     ingressPath.Path,
			PathType: pathType,
			Backend:  x.convertIngressBackendDO2DTO(ingressPath.Backend),
		}
	}
	return &pb.HTTPIngressRuleValue{
		Paths: HTTPIngressPath,
	}
}

func (x *IngressService) convertHTTPIngressPathDO2DTO(httpIngressPath networkingv1.HTTPIngressPath) *pb.HTTPIngressPath {
	return &pb.HTTPIngressPath{
		Backend: x.convertIngressBackendDO2DTO(httpIngressPath.Backend),
	}
}

func (x *IngressService) convertIngressBackendDO2DTO(backend networkingv1.IngressBackend) *pb.IngressBackend {
	return &pb.IngressBackend{
		Service: x.convertIngressServiceBackendDO2DTO(backend.Service),
	}
}

func (x *IngressService) convertIngressServiceBackendDO2DTO(service *networkingv1.IngressServiceBackend) *pb.IngressServiceBackend {
	if service == nil {
		return &pb.IngressServiceBackend{}
	}
	return &pb.IngressServiceBackend{
		Name: service.Name,
		Port: x.convertServiceBackendPortDO2DTO(service.Port),
	}
}

func (x *IngressService) convertServiceBackendPortDO2DTO(port networkingv1.ServiceBackendPort) *pb.ServiceBackendPort {
	return &pb.ServiceBackendPort{
		Name:   port.Name,
		Number: uint32(port.Number),
	}
}

func (x *IngressService) ListIngressHost(ctx context.Context, request *pb.ListHostRequest) (*pb.ListHostResponse, error) {
	domains, err := x.uc.ListIngressHost(ctx, &biz.ListIngressHostRequest{
		ClusterName: request.ClusterName,
		Namespace:   request.Namespace,
	})
	if err != nil {
		return nil, err
	}
	return &pb.ListHostResponse{
		List: domains,
	}, nil
}
