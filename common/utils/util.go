package utils

import (
	"bytes"
	"codo-cnmp/common/consts"
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/apimachinery/pkg/runtime/serializer/json"
	"k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/client-go/kubernetes/scheme"
	"net"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// UTCToLocal UTC时间转本地时间
func UTCToLocal(utcTime string) (string, error) {
	// 加载本地时区
	local, err := time.LoadLocation("Local")
	if err != nil {
		return "", fmt.Errorf("加载本地时区失败: %w", err)
	}

	// 解析UTC时间
	localTime, err := time.ParseInLocation(time.RFC3339, utcTime, local)
	if err != nil {
		return "", fmt.Errorf("解析UTC时间失败: %w", err)
	}

	// 返回本地时间字符串
	return localTime.Format(time.DateTime), nil
}

// Datetime2time 解析时间字符串并转为本地时间
func Datetime2time(datetime string, format string, location string) (time.Time, error) {
	// 加载本地时区
	if location == "" {
		location = "Asia/Shanghai"
	}
	local, err := time.LoadLocation(location)
	if err != nil {
		return time.Time{}, fmt.Errorf("加载本地时区失败: %w", err)
	}

	// 解析传入的时间字符串和格式
	localTime, err := time.ParseInLocation(format, datetime, local)
	if err != nil {
		return time.Time{}, fmt.Errorf("解析datetime时间失败: %w", err)
	}

	// 返回本地时间
	return localTime, nil
}

// MatchString 正则匹配字符串
func MatchString(pattern, str string) bool {
	if pattern == "" || str == "" {
		return true
	}
	matched, err := regexp.MatchString(pattern, str)
	if err != nil {
		return false
	}
	return matched
}

// MatchLabels 检查标签是否包含关键词
func MatchLabels(pattern string, labels map[string]string) bool {
	if labels == nil {
		return false
	}
	for key, value := range labels {
		if MatchString(pattern, key) || MatchString(pattern, value) {
			return true
		}
	}
	return false
}

// MatchContainerImages 检查容器镜像是否包含关键词
func MatchContainerImages(pattern string, containers []corev1.Container) bool {
	if containers == nil {
		return false
	}
	for _, container := range containers {
		if MatchString(pattern, container.Image) {
			return true
		}
	}
	return false
}

// MatchWorkloadName 检查工作负载名称是否包含关键词
func MatchWorkloadName(pattern string, ownerRefs []metav1.OwnerReference) bool {
	if ownerRefs == nil {
		return false
	}
	for _, ownerRef := range ownerRefs {
		if MatchString(pattern, ownerRef.Name) {
			return true
		}
	}
	return false
}

// MatchSelector 检查标签选择器是否包含关键词
func MatchSelector(pattern string, selector *metav1.LabelSelector) bool {
	if selector == nil {
		return false
	}
	for key, value := range selector.MatchLabels {
		if MatchString(pattern, key) || MatchString(pattern, value) {
			return true
		}
	}
	return false
}

// Base64Decode base64 解码
func Base64Decode(base64String string) ([]byte, error) {
	// 解码 Base64 字符串
	decodedBytes, err := base64.StdEncoding.DecodeString(base64String)
	if err != nil {
		return nil, fmt.Errorf("解码 Base64 字符串失败: %w", err)
	}
	// 将字节转换为字符串并打印输出
	return decodedBytes, nil
}

// AESEncrypt AES加密
func AESEncrypt(plainText string, key string) (string, error) {
	// 转换密钥和明文为字节数组
	keyBytes := []byte(key)
	plainTextBytes := []byte(plainText)

	// 创建AES加密块
	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		return "", fmt.Errorf("创建AES加密块失败: %w", err)
	}

	// 填充明文
	plainTextBytes = PKCS7Padding(plainTextBytes, block.BlockSize())

	// 初始化向量和blockSize大小保持一致
	iv := make([]byte, block.BlockSize())
	if _, err := rand.Read(iv); err != nil {
		return "", fmt.Errorf("初始化向量失败: %w", err)
	}

	// 加密密文
	cipherText := make([]byte, block.BlockSize()+len(plainTextBytes))
	copy(cipherText[:block.BlockSize()], iv)
	cipher.NewCBCEncrypter(block, iv).CryptBlocks(cipherText[block.BlockSize():], plainTextBytes)

	// 转换密文为Base64字符串
	return base64.StdEncoding.EncodeToString(cipherText), nil
}

// AESDecrypt AES解密
func AESDecrypt(cipherText string, key string) (string, error) {
	// 转为字节数组
	keyBytes := []byte(key)
	cipherTextBytes, err := base64.StdEncoding.DecodeString(cipherText)
	if err != nil {
		return "", fmt.Errorf("转换为字节数组失败: %w", err)
	}

	// 创建AES加密块
	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		return "", fmt.Errorf("创建AES加密块失败: %w", err)
	}
	blockSize := block.BlockSize()

	// 初始化向量和blockSize大小保持一致
	iv := cipherTextBytes[:blockSize]
	cipherTextBytes = cipherTextBytes[blockSize:]

	// CBC解密模式
	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(cipherTextBytes, cipherTextBytes)

	// 解除填充
	plainTextBytes, err := PKCS7UnPadding(cipherTextBytes)
	if err != nil {
		return "", fmt.Errorf("解除填充失败: %w", err)
	}

	return string(plainTextBytes), nil
}

// PKCS7Padding 对数据进行 PKCS#7 填充
func PKCS7Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	text := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, text...)
}

// PKCS7UnPadding 去除 PKCS#7 填充
func PKCS7UnPadding(origData []byte) ([]byte, error) {
	if len(origData) == 0 {
		return nil, fmt.Errorf("数据bytes为空")
	}
	unPadding := int(origData[len(origData)-1])
	if unPadding > len(origData) || unPadding > aes.BlockSize {
		return nil, fmt.Errorf("填充不合法")
	}
	return origData[:(len(origData) - unPadding)], nil
}

// ConvertMemoryToGiB 转换内存单位为 GiB
func ConvertMemoryToGiB(m resource.Quantity) float64 {
	// 转换为 GiB
	gbMemory := float64(m.Value()) / 1024 / 1024 / 1024
	// 保留两位小数
	formattedMemory, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", gbMemory), 64)
	return formattedMemory
}

// ConvertCPUToCores 转换 CPU 单位为 cores
func ConvertCPUToCores(cpu resource.Quantity) float64 {
	cores := float64(cpu.Value()) / 1000.0 // 将 CPU 转换为cores
	formattedCPU, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", cores), 64)
	return formattedCPU
}

// ConvertResourceTOYaml 将资源对象转换为 YAML 格式字符串
func ConvertResourceTOYaml(obj runtime.Object) (string, error) {
	encoder := json.NewSerializerWithOptions(
		json.DefaultMetaFactory, nil, nil,
		json.SerializerOptions{Yaml: true, Pretty: true, Strict: true},
	)
	var jsonBuffer strings.Builder
	if err := encoder.Encode(obj, &jsonBuffer); err != nil {
		return "", fmt.Errorf("序列化对象为 YAML 失败: %v", err)
	}
	return jsonBuffer.String(), nil
}

// IntToBool int转bool
func IntToBool(i int) bool {
	return i != 0
}

// BoolToInt bool转int
func BoolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

// Contains 检查字符串切片中是否包含指定字符串
func Contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// ParseK8sClusterRoleYAML 解析K8s ClusterRole YAML
func ParseK8sClusterRoleYAML(yamlStr string) (*rbacv1.ClusterRole, error) {
	// 解析YAML
	clusterRole := &rbacv1.ClusterRole{}
	decoder := yaml.NewYAMLOrJSONDecoder(strings.NewReader(yamlStr), 4096)
	if err := decoder.Decode(&clusterRole); err != nil {
		return nil, err
	}
	return clusterRole, nil
}

// K8sResourceValidator K8s资源验证器
type K8sResourceValidator struct {
	codecs serializer.CodecFactory
	scheme *runtime.Scheme
}

// NewK8sResourceValidator 创建新的资源验证器
func NewK8sResourceValidator() *K8sResourceValidator {
	// 使用默认的scheme
	s := runtime.NewScheme()
	scheme.AddToScheme(s)

	return &K8sResourceValidator{
		codecs: serializer.NewCodecFactory(s),
		scheme: s,
	}
}

// ValidateYAML 验证YAML格式的K8s资源
func (x *K8sResourceValidator) ValidateYAML(yamlContent string) error {
	if strings.TrimSpace(yamlContent) == "" {
		return fmt.Errorf("empty YAML content")
	}

	// 将YAML解析为非结构化对象
	_, gvk, err := x.decodeYAML(yamlContent)
	if err != nil {
		return fmt.Errorf("failed to decode YAML: %w", err)
	}

	// 验证资源类型
	if err := x.validateResourceType(gvk); err != nil {
		return err
	}

	return nil
}

// validateResourceType 验证资源类型是否支持
func (x *K8sResourceValidator) validateResourceType(gvk *schema.GroupVersionKind) error {
	if gvk == nil {
		return fmt.Errorf("GroupVersionKind is nil")
	}

	// 检查是否注册了该类型
	_, err := x.scheme.New(*gvk)
	if err != nil {
		return fmt.Errorf("unsupported resource type %s: %w", gvk.String(), err)
	}

	return nil
}

// decodeYAML 解码YAML内容
func (x *K8sResourceValidator) decodeYAML(content string) (runtime.Object, *schema.GroupVersionKind, error) {
	jsonData, err := yaml.ToJSON([]byte(content))
	if err != nil {
		return nil, nil, fmt.Errorf("failed to convert YAML to JSON: %w", err)
	}

	obj, gvk, err := x.codecs.UniversalDeserializer().Decode(jsonData, nil, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to decode resource: %w", err)
	}

	return obj, gvk, nil
}

// AggregateToACL 将 Kubernetes RBAC ClusterRole 规则转换为 ACL 格式
// clusterName: 集群名称
// namespace: 命名空间，"all" 将被转换为 "*"
// 返回: ACL 规则列表
func AggregateToACL(clusterRole *rbacv1.ClusterRole, clusterName string, namespace string) ([]string, error) {
	aclList := make([]string, 0)
	// 参数验证
	if clusterRole == nil {
		return aclList, fmt.Errorf("ClusterRole 为空")
	}
	if clusterName == "" {
		return aclList, fmt.Errorf("集群名称为空")
	}
	if namespace == "" {
		return aclList, fmt.Errorf("命名空间为空")
	}

	uniqueCombos := make(map[string]bool)

	// 处理资源规则
	for _, rule := range clusterRole.Rules {
		apiGroups := rule.APIGroups
		if len(apiGroups) == 0 {
			apiGroups = []string{"*"}
		}

		ns := namespace
		if ns == "all" {
			ns = "*"
		}

		for _, apiGroup := range apiGroups {
			if apiGroup == "" {
				apiGroup = "*"
			}
			for _, resourceStr := range rule.Resources {
				for _, verb := range rule.Verbs {
					aclRule := fmt.Sprintf("/%s/%s/%s/%s/%s",
						clusterName,
						ns,
						apiGroup,
						resourceStr,
						verb)

					if !uniqueCombos[aclRule] {
						aclList = append(aclList, aclRule)
						uniqueCombos[aclRule] = true
					}
				}
			}
		}
	}

	// 处理非资源 URL 规则
	for _, rule := range clusterRole.Rules {
		if len(rule.NonResourceURLs) > 0 {
			for _, nonResourceURL := range rule.NonResourceURLs {
				for _, verb := range rule.Verbs {
					aclRule := fmt.Sprintf("/%s/*/%s/%s",
						clusterName,
						nonResourceURL,
						verb)

					if !uniqueCombos[aclRule] {
						aclList = append(aclList, aclRule)
						uniqueCombos[aclRule] = true
					}
				}
			}
		}
	}
	return aclList, nil
}

func GetUserIDFromCtx(ctx context.Context) (uint32, error) {
	contextVal := ctx.Value(consts.ContextUserIDKey)
	if contextVal == nil {
		return 0, fmt.Errorf("用户未登录")
	}
	// 断言
	userId, ok := contextVal.(string)
	if !ok {
		return 0, fmt.Errorf("用户未登录")
	}
	// 转换为整数
	id, err := strconv.ParseUint(userId, 10, 32)
	if err != nil {
		return 0, fmt.Errorf("用户未登录")
	}
	return uint32(id), nil
}

func GetUserNameFromCtx(ctx context.Context) (string, error) {
	contextVal := ctx.Value(consts.ContextUserNameKey)
	if contextVal == nil {
		return "", fmt.Errorf("用户未登录")
	}
	// 断言
	userName, ok := contextVal.(string)
	if !ok {
		return "", fmt.Errorf("用户未登录")
	}
	return userName, nil
}

// CheckConnection 检查连接
func CheckConnection(domain string) (bool, error) {
	parsedURL, err := url.Parse(domain)
	if err != nil {
		return false, fmt.Errorf("解析域名失败: %w", err)
	}
	host := parsedURL.Host
	if !strings.Contains(host, ":") {
		switch parsedURL.Scheme {
		case "https":
			host += ":443"
		case "http":
			host += ":80"
		default:
			host += ":80"
		}
	}
	conn, err := net.DialTimeout("tcp", host, 2*time.Second)
	if err != nil {
		return false, fmt.Errorf("连接失败: %w", err)
	}
	defer conn.Close()
	return true, nil
}

// IsEnglish 判断是否为全英文字符串（仅限 A-Z 和 a-z）
func IsEnglish(s string) bool {
	if len(s) == 0 {
		return false
	}
	for _, r := range s {
		if !((r >= 'A' && r <= 'Z') || (r >= 'a' && r <= 'z')) {
			return false
		}
	}
	return true
}
