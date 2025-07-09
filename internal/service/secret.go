package service

import (
	"codo-cnmp/common/utils"
	"codo-cnmp/internal/biz"
	"codo-cnmp/pb"
	"context"
	"encoding/json"
	"fmt"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sort"
	"time"
)

type SecretService struct {
	pb.UnimplementedSecretServer
	uc *biz.SecretUseCase
	uf *biz.UserFollowUseCase
}

func NewSecretService(uc *biz.SecretUseCase, uf *biz.UserFollowUseCase) *SecretService {
	return &SecretService{
		uc: uc,
		uf: uf,
	}
}

// getUserFollowMap 获取用户关注的CloneSet列表
func (x *SecretService) getUserFollowMap(ctx context.Context, userID uint32) (map[string]bool, error) {
	userFollows, _, err := x.uf.ListUserFollow(ctx, &biz.ListUserFollowRequest{
		UserFollowCommonParams: biz.UserFollowCommonParams{
			UserID:     userID,
			FollowType: pb.FollowType_Secret,
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
func (x *SecretService) setFollowedStatus(ctx context.Context, clusterName string, items []*pb.SecretItem) error {
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

// parseDockerConfigJson 解析 docker config json
func parseDockerConfigJson(configJson string) (registry, username, password string, err error) {
	var config struct {
		Auths map[string]struct {
			Username string `json:"username"`
			Password string `json:"password"`
			Auth     string `json:"auth"`
		} `json:"auths"`
	}

	if err := json.Unmarshal([]byte(configJson), &config); err != nil {
		return "", "", "", fmt.Errorf("unmarshal docker config: %w", err)
	}

	// 获取第一个 registry 的信息
	for reg, auth := range config.Auths {
		return reg, auth.Username, auth.Password, nil
	}

	return "", "", "", fmt.Errorf("no registry found in docker config")
}

func (x *SecretService) convertDO2DTO(secret *corev1.Secret) *pb.SecretItem {
	var (
		createTime, updateTime       time.Time
		yamlStr                      string
		secretType                   pb.SecretType
		tlsCert                      string
		tlsKey                       string
		registry, username, password string
	)
	createTime = secret.CreationTimestamp.Time
	for _, managedField := range secret.ManagedFields {
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
	if secret.Kind == "" {
		secret.Kind = "Secret"
	}
	if secret.APIVersion == "" {
		secret.APIVersion = "v1"
	}
	yamlStr, err := utils.ConvertResourceTOYaml(secret)
	if err != nil {
		yamlStr = ""
	}

	switch secret.Type {
	case corev1.SecretTypeOpaque:
		secretType = pb.SecretType_Opaque
	case corev1.SecretTypeDockerConfigJson:
		secretType = pb.SecretType_DockerConfigJson
		// 解析 .dockerconfigjson
		dockerConfigJson := string(secret.Data[corev1.DockerConfigJsonKey])
		if dockerConfigJson == "" {
			registry = ""
			username = ""
			password = ""
		} else {
			registry, username, password, err = parseDockerConfigJson(dockerConfigJson)
			if err != nil {
				registry = ""
				username = ""
				password = ""
			}
		}
	case corev1.SecretTypeTLS:
		secretType = pb.SecretType_TLS
		tlsCert = string(secret.Data["tls.crt"])
		tlsKey = string(secret.Data["tls.key"])
	default:
		secretType = pb.SecretType_Unknown_SecretType
	}

	return &pb.SecretItem{
		Name:        secret.Name,
		Labels:      secret.Labels,
		Annotations: secret.Annotations,
		CreateTime:  uint64(createTime.UnixNano() / 1e6),
		UpdateTime:  uint64(updateTime.UnixNano() / 1e6),
		Data:        secret.Data,
		RefCount:    0, // todo
		Yaml:        yamlStr,
		Type:        secretType,
		TlsCrt:      tlsCert,
		TlsKey:      tlsKey,
		Registry:    registry,
		Username:    username,
		Password:    password,
	}
}

func (x *SecretService) ListSecret(ctx context.Context, req *pb.ListSecretsRequest) (*pb.ListSecretsResponse, error) {
	secrets, total, err := x.uc.ListSecret(ctx, &biz.ListSecretRequest{
		ClusterName: req.ClusterName,
		Namespace:   req.Namespace,
		Keyword:     req.Keyword,
		Page:        req.Page,
		PageSize:    req.PageSize,
		ListAll:     utils.IntToBool(int(req.ListAll)),
		SecretType:  req.Type,
	})
	if err != nil {
		return nil, err
	}

	list := make([]*pb.SecretItem, 0, len(secrets))
	for _, secret := range secrets {
		dto := x.convertDO2DTO(secret)
		refCount, refs, err := x.uc.GetSecretReferences(ctx, &biz.SecretCommonParams{
			ClusterName: req.ClusterName,
			Namespace:   secret.Namespace,
			Name:        secret.Name,
		})
		if err != nil {
			continue
		}
		for _, ref := range refs {
			dto.Refs = append(dto.Refs, &pb.SecretReference{
				Kind: ref["kind"],
				Name: ref["name"],
			})

		}
		dto.RefCount = uint32(refCount)
		list = append(list, dto)
	}
	if err := x.setFollowedStatus(ctx, req.ClusterName, list); err != nil {
		return nil, err
	}
	sort.SliceStable(list, func(i, j int) bool {
		return list[i].IsFollowed && !list[j].IsFollowed
	})

	return &pb.ListSecretsResponse{
		List:  list,
		Total: total,
	}, nil
}

func (x *SecretService) CreateOrUpdateSecretByYaml(ctx context.Context, req *pb.CreateOrUpdateSecretByYamlRequest) (*pb.CreateOrUpdateSecretByYamlResponse, error) {
	err := x.uc.CreateOrUpdateSecretByYamlRequest(ctx, &biz.CreateOrUpdateSecretByYamlRequest{
		ClusterName: req.ClusterName,
		Yaml:        req.Yaml,
	})
	if err != nil {
		return nil, err
	}
	return &pb.CreateOrUpdateSecretByYamlResponse{}, nil
}

func (x *SecretService) CreateSecret(ctx context.Context, req *pb.CreateSecretRequest) (*pb.CreateSecretResponse, error) {
	err := x.uc.CreateSecret(ctx, &biz.CreateSecretRequest{
		SecretCommonParams: biz.SecretCommonParams{
			ClusterName: req.ClusterName,
			Name:        req.Name,
			Namespace:   req.Namespace,
		},
		Labels:      req.Labels,
		Annotations: req.Annotations,
		Data:        req.Data,
	})
	if err != nil {
		return nil, err
	}
	return &pb.CreateSecretResponse{}, nil
}

func (x *SecretService) UpdateSecret(ctx context.Context, req *pb.UpdateSecretRequest) (*pb.UpdateSecretResponse, error) {
	err := x.uc.UpdateSecret(ctx, &biz.CreateSecretRequest{
		SecretCommonParams: biz.SecretCommonParams{
			ClusterName: req.ClusterName,
			Name:        req.Name,
			Namespace:   req.Namespace,
		},
		Labels:      req.Labels,
		Annotations: req.Annotations,
		Data:        req.Data,
	})
	if err != nil {
		return nil, err
	}
	return &pb.UpdateSecretResponse{}, nil
}

func (x *SecretService) DeleteSecret(ctx context.Context, req *pb.DeleteSecretRequest) (*pb.DeleteSecretResponse, error) {
	err := x.uc.DeleteSecret(ctx, &biz.DeleteSecretRequest{
		SecretCommonParams: biz.SecretCommonParams{
			ClusterName: req.ClusterName,
			Name:        req.Name,
			Namespace:   req.Namespace,
		},
	})
	if err != nil {
		return nil, err
	}
	return &pb.DeleteSecretResponse{}, nil
}

func (x *SecretService) GetSecretDetail(ctx context.Context, req *pb.SecretDetailRequest) (*pb.SecretDetailResponse, error) {
	secret, err := x.uc.GetSecret(ctx, &biz.GetSecretRequest{
		SecretCommonParams: biz.SecretCommonParams{
			ClusterName: req.ClusterName,
			Namespace:   req.Namespace,
			Name:        req.Name,
		},
	})
	if err != nil {
		return nil, err
	}
	count, _, err := x.uc.GetSecretReferences(ctx, &biz.SecretCommonParams{
		ClusterName: req.ClusterName,
		Namespace:   req.Namespace,
		Name:        req.Name,
	})
	if err != nil {
		return nil, err
	}
	detail := x.convertDO2DTO(secret)
	detail.RefCount = uint32(count)
	return &pb.SecretDetailResponse{
		Detail: detail,
	}, nil
}
