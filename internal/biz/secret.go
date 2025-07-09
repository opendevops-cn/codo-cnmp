package biz

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"codo-cnmp/pb"
	"k8s.io/client-go/kubernetes"

	"codo-cnmp/common/utils"
	"github.com/go-kratos/kratos/v2/log"
	corev1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/yaml"
)

type SecretCommonParams struct {
	ClusterName string `json:"cluster_name"`
	Namespace   string `json:"namespace"`
	Name        string `json:"name"`
}

// ListSecretRequest Secret列表请求
type ListSecretRequest struct {
	ClusterName string `json:"cluster_name"`
	Namespace   string `json:"namespace"`
	Keyword     string `json:"keyword"`
	Page        uint32 `json:"page"`
	PageSize    uint32 `json:"page_size"`
	ListAll     bool   `json:"list_all"`
	SecretType  pb.SecretType
}

// GetSecretRequest Secret详情请求
type GetSecretRequest struct {
	SecretCommonParams
}

// DeleteSecretRequest Secret删除请求
type DeleteSecretRequest struct {
	SecretCommonParams
}

// CreateSecretRequest Secret创建请求
type CreateSecretRequest struct {
	SecretCommonParams
	Data        map[string]string `json:"data"`
	Type        pb.SecretType     `json:"type"`
	Labels      map[string]string `json:"labels"`
	Annotations map[string]string `json:"annotations"`
}

// CreateOrUpdateSecretByYamlRequest Secret创建或更新请求
type CreateOrUpdateSecretByYamlRequest struct {
	ClusterName string `json:"cluster_name"`
	Yaml        string `json:"yaml"`
}

type ISecretUseCase interface {
	ListSecret(ctx context.Context, req *ListSecretRequest) ([]*corev1.Secret, uint32, error)
	GetSecret(ctx context.Context, req *GetSecretRequest) (*corev1.Secret, error)
	DeleteSecret(ctx context.Context, req *DeleteSecretRequest) error
	CreateSecret(ctx context.Context, req *CreateSecretRequest) error
	UpdateSecret(ctx context.Context, req *CreateSecretRequest) error
	CreateOrUpdateSecretByYamlRequest(ctx context.Context, req *CreateOrUpdateSecretByYamlRequest) error
}

type SecretUseCase struct {
	cluster IClusterUseCase
	log     *log.Helper
}

func (x *SecretUseCase) ListSecret(ctx context.Context, req *ListSecretRequest) ([]*corev1.Secret, uint32, error) {
	clientSet, err := x.cluster.GetClientSetByClusterName(ctx, req.ClusterName)
	if err != nil {
		return nil, 0, err
	}
	if req.Page == 0 {
		req.Page = 1
	}
	if req.PageSize == 0 {
		req.PageSize = 10
	}

	var secretType string
	switch req.SecretType {
	case pb.SecretType_Unknown_SecretType:
		secretType = ""
	case pb.SecretType_Opaque:
		secretType = string(corev1.SecretTypeOpaque)
	case pb.SecretType_DockerConfigJson:
		secretType = string(corev1.SecretTypeDockerConfigJson)
	case pb.SecretType_DockerConfig:
		secretType = string(corev1.SecretTypeDockerConfigJson)
	case pb.SecretType_TLS:
		secretType = string(corev1.SecretTypeTLS)
	case pb.SecretType_BasicAuth:
		secretType = string(corev1.SecretTypeBasicAuth)
	case pb.SecretType_ServiceAccountToken:
		secretType = string(corev1.SecretTypeServiceAccountToken)
	case pb.SecretType_SSHAuth:
		secretType = string(corev1.SecretTypeSSHAuth)
	case pb.SecretType_BootstrapToken:
		secretType = string(corev1.SecretTypeBootstrapToken)
	}

	var (
		allFilteredSecrets = make([]*corev1.Secret, 0)
		continueToken      = ""
		limit              = int64(req.PageSize)
	)
	for {
		secretListOptions := metav1.ListOptions{
			Limit:    limit,
			Continue: continueToken,
		}
		secretList, err := clientSet.CoreV1().Secrets(req.Namespace).List(ctx, secretListOptions)
		if err != nil {
			x.log.WithContext(ctx).Errorf("获取Secret列表失败: %v", err)
			return nil, 0, fmt.Errorf("获取Secret列表失败: %w", err)
		}
		filteredSecretList := x.filterSecretByKeyword(secretList, req.Keyword, secretType)
		for _, configMap := range filteredSecretList.Items {
			allFilteredSecrets = append(allFilteredSecrets, &configMap)
		}
		if filteredSecretList.Continue == "" {
			break
		}
		continueToken = filteredSecretList.Continue
	}
	// 返回所有结果或者分页结果
	if req.ListAll {
		return allFilteredSecrets, uint32(len(allFilteredSecrets)), nil
	}
	if len(allFilteredSecrets) == 0 {
		return allFilteredSecrets, 0, nil
	}
	paginatedSecrets, total := utils.K8sPaginate(allFilteredSecrets, req.Page, req.PageSize)
	return paginatedSecrets, total, nil

}

func (x *SecretUseCase) filterSecretByKeyword(secrets *corev1.SecretList, keyword, secretType string) *corev1.SecretList {
	if keyword == "" && secretType == "" {
		return secrets
	}
	results := &corev1.SecretList{}
	pattern := "(?i).*" + regexp.QuoteMeta(keyword) + ".*"
	for _, secret := range secrets.Items {
		if secretType == "" {
			if utils.MatchString(pattern, secret.Name) {
				results.Items = append(results.Items, secret)
			}
		} else {
			if utils.MatchString(pattern, secret.Name) && secret.Type == corev1.SecretType(secretType) {

			}
		}

	}
	return results
}

func (x *SecretUseCase) GetSecret(ctx context.Context, req *GetSecretRequest) (*corev1.Secret, error) {
	clientSet, err := x.cluster.GetClientSetByClusterName(ctx, req.ClusterName)
	if err != nil {
		return nil, err
	}
	secret, err := clientSet.CoreV1().Secrets(req.Namespace).Get(ctx, req.Name, metav1.GetOptions{})
	if err != nil {
		x.log.WithContext(ctx).Errorf("获取Secret失败: %v", err)
		return nil, fmt.Errorf("获取Secret失败: %w", err)
	}
	return secret, nil
}

func (x *SecretUseCase) DeleteSecret(ctx context.Context, req *DeleteSecretRequest) error {
	clientSet, err := x.cluster.GetClientSetByClusterName(ctx, req.ClusterName)
	if err != nil {
		return err
	}
	err = clientSet.CoreV1().Secrets(req.Namespace).Delete(ctx, req.Name, metav1.DeleteOptions{})
	if err != nil {
		x.log.WithContext(ctx).Errorf("删除Secret失败: %v", err)
		return fmt.Errorf("删除Secret失败: %w", err)
	}
	return nil
}

// createDockerConfigJson 创建 docker config json
func createDockerConfigJson(data map[string]string) (string, error) {
	// 获取认证信息
	registry := data["registry"]
	username := data["username"]
	password := data["password"]

	// 验证必要参数
	if registry == "" || username == "" || password == "" {
		return "", fmt.Errorf("registry, username and password are required")
	}

	// 生成 base64 编码的认证信息
	auth := base64.StdEncoding.EncodeToString([]byte(username + ":" + password))

	// 构建 docker config 结构
	dockerConfig := struct {
		Auths map[string]struct {
			Username string `json:"username"`
			Password string `json:"password"`
			Auth     string `json:"auth"`
		} `json:"auths"`
	}{
		Auths: map[string]struct {
			Username string `json:"username"`
			Password string `json:"password"`
			Auth     string `json:"auth"`
		}{
			registry: {
				Username: username,
				Password: password,
				Auth:     auth,
			},
		},
	}

	// 转换为 JSON
	configBytes, err := json.Marshal(dockerConfig)
	if err != nil {
		return "", fmt.Errorf("marshal docker config: %w", err)
	}

	return string(configBytes), nil
}

func (x *SecretUseCase) CreateSecret(ctx context.Context, req *CreateSecretRequest) error {
	clientSet, err := x.cluster.GetClientSetByClusterName(ctx, req.ClusterName)
	if err != nil {
		return err
	}
	data := make(map[string][]byte)
	var secretType corev1.SecretType
	switch req.Type {
	case pb.SecretType_Opaque:
		secretType = corev1.SecretTypeOpaque
		if req.Data != nil {
			for k, v := range req.Data {
				data[k] = []byte(v)
			}
		}
	case pb.SecretType_DockerConfigJson:
		secretType = corev1.SecretTypeDockerConfigJson
		if req.Data != nil {
			// 创建 docker config json
			dockerConfigJson, err := createDockerConfigJson(req.Data)
			if err != nil {
				return fmt.Errorf("创建docker config json失败: %w", err)
			}

			// 更新 req.Data
			req.Data = map[string]string{
				".dockerconfigjson": dockerConfigJson,
			}
		}
	case pb.SecretType_TLS:
		secretType = corev1.SecretTypeTLS
		data["tls.crt"] = []byte(req.Data["tls_crt"])
		data["tls.key"] = []byte(req.Data["tls_key"])
	default:
		return fmt.Errorf("不支持的Secret类型: %d", req.Type)
	}
	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:        req.Name,
			Namespace:   req.Namespace,
			Labels:      req.Labels,
			Annotations: req.Annotations,
		},
		Data: data,
		Type: secretType,
	}
	_, err = clientSet.CoreV1().Secrets(req.Namespace).Create(ctx, secret, metav1.CreateOptions{})
	if err != nil {
		x.log.WithContext(ctx).Errorf("创建Secret失败: %v", err)
		return fmt.Errorf("创建Secret失败: %w", err)
	}
	return nil
}

func (x *SecretUseCase) UpdateSecret(ctx context.Context, req *CreateSecretRequest) error {
	clientSet, err := x.cluster.GetClientSetByClusterName(ctx, req.ClusterName)
	if err != nil {
		return err
	}
	data := make(map[string][]byte)
	var secretType corev1.SecretType
	switch req.Type {
	case pb.SecretType_Opaque:
		secretType = corev1.SecretTypeOpaque
		if req.Data != nil {
			for k, v := range req.Data {
				data[k] = []byte(v)
			}
		}
	case pb.SecretType_DockerConfigJson:
		secretType = corev1.SecretTypeDockerConfigJson
		if req.Data != nil {
			// 创建 docker config json
			dockerConfigJson, err := createDockerConfigJson(req.Data)
			if err != nil {
				return fmt.Errorf("创建docker config json失败: %w", err)
			}

			// 更新 req.Data
			req.Data = map[string]string{
				corev1.DockerConfigJsonKey: dockerConfigJson,
			}
		}
	case pb.SecretType_TLS:
		secretType = corev1.SecretTypeTLS
		data["tls.crt"] = []byte(req.Data["tls_crt"])
		data["tls.key"] = []byte(req.Data["tls_key"])
	default:
		return fmt.Errorf("不支持的Secret类型: %d", req.Type)
	}
	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:        req.Name,
			Namespace:   req.Namespace,
			Labels:      req.Labels,
			Annotations: req.Annotations,
		},
		Data: data,
		Type: secretType,
	}
	_, err = clientSet.CoreV1().Secrets(req.Namespace).Update(ctx, secret, metav1.UpdateOptions{})
	if err != nil {
		x.log.WithContext(ctx).Errorf("更新Secret失败: %v", err)
		return fmt.Errorf("更新Secret失败: %w", err)
	}
	return nil
}

func (x *SecretUseCase) CreateOrUpdateSecretByYamlRequest(ctx context.Context, req *CreateOrUpdateSecretByYamlRequest) error {
	clientSet, err := x.cluster.GetClientSetByClusterName(ctx, req.ClusterName)
	if err != nil {
		return err
	}
	secret := &corev1.Secret{}
	decoder := yaml.NewYAMLOrJSONDecoder(strings.NewReader(req.Yaml), 4096)
	if err := decoder.Decode(&secret); err != nil {
		x.log.WithContext(ctx).Errorf("解析YAML失败, 请检查格式是否正确: %v", err)
		return fmt.Errorf("解析YAML失败, 请检查格式是否正确: %v", err)
	}
	_, err = clientSet.CoreV1().Secrets(secret.Namespace).Get(ctx, secret.Name, metav1.GetOptions{})
	if err != nil {
		if k8serrors.IsNotFound(err) {
			// 如果不存在则创建
			_, err = clientSet.CoreV1().Secrets(secret.Namespace).Create(ctx, secret, metav1.CreateOptions{})
			if err != nil {
				x.log.WithContext(ctx).Errorf("创建ConfigMap失败: %v", err)
				return fmt.Errorf("创建ConfigMap失败: %w", err)
			}
			return nil
		}
		x.log.WithContext(ctx).Errorf("获取Secret失败: %v", err)
		return fmt.Errorf("获取Secret失败: %w", err)
	}
	// 如果存在则更新
	_, err = clientSet.CoreV1().Secrets(secret.Namespace).Update(ctx, secret, metav1.UpdateOptions{})
	if err != nil {
		x.log.WithContext(ctx).Errorf("更新Secret失败: %v", err)
		return fmt.Errorf("更新Secret失败: %w", err)
	}
	return nil
}

func (x *SecretUseCase) GetSecretReferences(ctx context.Context, req *SecretCommonParams) (int32, []map[string]string, error) {
	clientSet, err := x.cluster.GetClientSetByClusterName(ctx, req.ClusterName)
	if err != nil {
		return 0, nil, err
	}

	refMap := make(map[string]struct{})
	continueToken := ""

	// 分页获取所有 pods
	for {
		pods, err := clientSet.CoreV1().Pods(req.Namespace).List(ctx, metav1.ListOptions{
			Limit:    100,
			Continue: continueToken,
		})
		if err != nil {
			return 0, nil, fmt.Errorf("获取Pod列表失败: %w", err)
		}

		// 处理当前页的 pods
		for _, pod := range pods.Items {
			if !isPodReferencingSecret(&pod, req.Name) {
				continue
			}
			topOwner, err := x.findTopOwner(ctx, clientSet, pod.Namespace, pod.OwnerReferences)
			if err != nil {
				return 0, nil, err
			}
			if topOwner != nil {
				refMap[fmt.Sprintf("%s/%s", topOwner.Kind, topOwner.Name)] = struct{}{}
			} else {
				refMap[fmt.Sprintf("Pod/%s", pod.Name)] = struct{}{}
			}
		}

		continueToken = pods.Continue
		if continueToken == "" {
			break
		}
	}

	// 转换结果格式
	result := make([]map[string]string, 0, len(refMap))
	for ref := range refMap {
		kind, name, _ := strings.Cut(ref, "/")
		result = append(result, map[string]string{
			"kind": kind,
			"name": name,
		})
	}

	return int32(len(refMap)), result, nil
}

// isPodReferencingSecret 检查 Pod 是否引用了指定的 Secret
func isPodReferencingSecret(pod *corev1.Pod, secretName string) bool {
	// 检查镜像拉取密钥
	for _, pullSecret := range pod.Spec.ImagePullSecrets {
		if pullSecret.Name == secretName {
			return true
		}
	}
	// 检查卷挂载
	for _, volume := range pod.Spec.Volumes {
		if volume.Secret != nil && volume.Secret.SecretName == secretName {
			return true
		}
		// 检查 projected 卷
		if volume.Projected != nil {
			for _, source := range volume.Projected.Sources {
				if source.Secret != nil && source.Secret.Name == secretName {
					return true
				}
			}
		}
	}
	// 检查所有容器
	for _, container := range pod.Spec.Containers {
		// 检查 envFrom
		for _, envFrom := range container.EnvFrom {
			if envFrom.SecretRef != nil && envFrom.SecretRef.Name == secretName {
				return true
			}
		}
		// 检查 env
		for _, env := range container.Env {
			if env.ValueFrom != nil &&
				env.ValueFrom.SecretKeyRef != nil &&
				env.ValueFrom.SecretKeyRef.Name == secretName {
				return true
			}
		}
	}

	return false
}

// findTopOwner 递归查找最顶层的所有者
func (x *SecretUseCase) findTopOwner(ctx context.Context, clientSet *kubernetes.Clientset, namespace string, owners []metav1.OwnerReference) (*metav1.OwnerReference, error) {
	if len(owners) == 0 {
		return nil, nil
	}
	owner := owners[0]
	switch owner.Kind {
	case "ReplicaSet":
		rs, err := clientSet.AppsV1().ReplicaSets(namespace).Get(ctx, owner.Name, metav1.GetOptions{})
		if err != nil {
			return nil, err
		}
		if len(rs.OwnerReferences) > 0 {
			return x.findTopOwner(ctx, clientSet, namespace, rs.OwnerReferences)
		}
	case "StatefulSet", "DaemonSet", "Deployment", "Job", "CronJob", "CloneSet", "GameServerSet":
		return &owner, nil
	}

	return &owner, nil
}

func NewSecretUseCase(cluster IClusterUseCase, logger log.Logger) *SecretUseCase {
	return &SecretUseCase{
		cluster: cluster,
		log:     log.NewHelper(log.With(logger, "module", "biz/secret")),
	}
}

func NewISecretUseCase(x *SecretUseCase) ISecretUseCase {
	return x
}
