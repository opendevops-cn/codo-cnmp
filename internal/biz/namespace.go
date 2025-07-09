package biz

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"codo-cnmp/common/utils"
	"github.com/ccheers/xpkg/sync/errgroup"
	"github.com/ccheers/xpkg/xmsgbus"
	"github.com/go-kratos/kratos/v2/log"
	corev1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/yaml"
)

var (
	DefaultKubernetesNameSpaces = map[string]struct{}{
		"default":         {},
		"kube-system":     {},
		"kube-public":     {},
		"kube-node-lease": {},
	}
	DESCRIPTION = "description"
)

func IsDefaultNameSpaces(name string) bool {
	_, exists := DefaultKubernetesNameSpaces[name]
	return exists
}

type NamespaceCommonParams struct {
	ClusterName string
	Name        string
}

type NameSpaceItem struct {
	UID         string
	Name        string
	Description string
	CreateTime  string
	State       string
	IsDefault   bool // 是否Default
	Yaml        string
	Labels      map[string]string
	Annotations map[string]string
}

type ListNameSpaceRequest struct {
	ClusterName string
	Keyword     string
	Page        uint32
	PageSize    uint32
	ListAll     bool
}

type CreateNameSpaceRequest struct {
	ClusterName string
	Name        string
	Description string
	Labels      map[string]string
	Annotations map[string]string
}

type ListNameSpaceResponse struct {
	List  []*NameSpaceItem
	Total int32
}

type DeleteNameSpaceRequest struct {
	ClusterName string
	Name        string
}

type UpdateNameSpaceRequest struct {
	ClusterName string
	UID         string
	Name        string
	Description string
	Labels      map[string]string
	Annotations map[string]string
}

type GetNamespaceYamlRequest struct {
	NamespaceCommonParams
}

type GetNamespaceRequest struct {
	NamespaceCommonParams
}

type CreateNameSpaceByYamlRequest struct {
	ClusterName string
	Yaml        string
}

type NameSpaceDeletionEvent struct {
	ClusterId uint32
	NameSpace string
}

func (x *NameSpaceDeletionEvent) Topic() string {
	return "namespace.deletion"
}

type INameSpaceUseCases interface {
	ListNameSpace(ctx context.Context, req *ListNameSpaceRequest) ([]*corev1.Namespace, uint32, error)
	CreateNameSpace(ctx context.Context, req *CreateNameSpaceRequest) error
	UpdateNameSpace(ctx context.Context, req *UpdateNameSpaceRequest) error
	DeleteNameSpace(ctx context.Context, req *DeleteNameSpaceRequest) error
	CreateNameSpaceByYaml(ctx context.Context, req *CreateNameSpaceByYamlRequest) error
	UpdateNameSpaceByYaml(ctx context.Context, req *CreateNameSpaceByYamlRequest) error
	GetNameSpace(ctx context.Context, req *GetNamespaceRequest) (*corev1.Namespace, error)
}

type NameSpaceUseCase struct {
	cluster              IClusterUseCase
	roleBinding          IRoleBindingUseCase
	log                  *log.Helper
	msgbus               xmsgbus.IMsgBus
	tm                   xmsgbus.ITopicManager
	otelOptions          *xmsgbus.OTELOptions
	pubNameSpaceDeletion xmsgbus.IPublisher[*NameSpaceDeletionEvent]
}

func NewINameSpaceUseCase(x *NameSpaceUseCase) INameSpaceUseCases {
	return x
}

func NewNameSpaceUseCase(ctx context.Context, cluster IClusterUseCase, logger log.Logger, bus xmsgbus.IMsgBus, tm xmsgbus.ITopicManager, roleBinding IRoleBindingUseCase) (*NameSpaceUseCase, func()) {
	ns := &NameSpaceUseCase{
		cluster:              cluster,
		roleBinding:          roleBinding,
		msgbus:               bus,
		otelOptions:          xmsgbus.NewOTELOptions(),
		tm:                   tm,
		log:                  log.NewHelper(log.With(logger, "module", "biz/namespace")),
		pubNameSpaceDeletion: xmsgbus.NewPublisher[*NameSpaceDeletionEvent](bus, tm, xmsgbus.NewOTELOptions()),
	}
	return ns, ns.Init(ctx)
}

// 检查关键词是否在 Pod 名称、状态、命名空间、镜像中
func filterNamespaceByKeyword(namespaces *corev1.NamespaceList, keyword string) []*corev1.Namespace {
	var result []*corev1.Namespace
	// 编译正则表达式，(?i)表示忽略大小写
	pattern := "(?i).*" + regexp.QuoteMeta(keyword) + ".*"
	for _, ns := range namespaces.Items {
		description := ns.GetAnnotations()[DESCRIPTION]
		if utils.MatchString(pattern, ns.Name) || (description != "" && utils.MatchString(pattern, description)) {
			result = append(result, &ns)
		}
	}
	return result
}

func (x *NameSpaceUseCase) ListNameSpace(ctx context.Context, req *ListNameSpaceRequest) ([]*corev1.Namespace, uint32, error) {
	clientSet, err := x.cluster.GetClientSetByClusterName(ctx, req.ClusterName)
	if err != nil {
		return nil, 0, fmt.Errorf("创建k8s client失败: %v", err)
	}

	if req.Page == 0 {
		req.Page = 1
	}

	if req.PageSize == 0 {
		req.PageSize = 10
	}

	var (
		continueToken         string
		limit                 = int64(req.PageSize)
		allFilteredNamespaces []*corev1.Namespace
	)

	for {
		namespaceListOptions := metav1.ListOptions{
			Limit:    limit,
			Continue: continueToken,
		}
		namespaces, err := clientSet.CoreV1().Namespaces().List(ctx, namespaceListOptions)
		if err != nil {
			x.log.Warnf("查询namespaces失败: %v", err)
			return nil, 0, fmt.Errorf("查询namespaces失败: %w", err)
		}
		filteredNamespaces := filterNamespaceByKeyword(namespaces, req.Keyword)
		allFilteredNamespaces = append(allFilteredNamespaces, filteredNamespaces...)

		if namespaces.Continue == "" {
			break
		}
		// 更新 continueToken，继续获取下一页
		continueToken = namespaces.Continue
	}

	// 如果 ListAll 为 true，直接返回所有结果
	if req.ListAll {
		return allFilteredNamespaces, uint32(len(allFilteredNamespaces)), nil
	}
	if len(allFilteredNamespaces) == 0 {
		return nil, 0, nil
	}
	paginatedNameSpaces, total := utils.K8sPaginate(allFilteredNamespaces, req.Page, req.PageSize)
	return paginatedNameSpaces, total, nil
}

// CreateNameSpace 创建NameSpace
func (x *NameSpaceUseCase) CreateNameSpace(ctx context.Context, req *CreateNameSpaceRequest) error {
	clientSet, err := x.cluster.GetClientSetByClusterName(ctx, req.ClusterName)
	if err != nil {
		return fmt.Errorf("创建k8s client失败: %w", err)
	}
	_, err = clientSet.CoreV1().Namespaces().Get(ctx, req.Name, metav1.GetOptions{})
	if err == nil {
		// 如果没有错误，说明命名空间已存在
		return fmt.Errorf("namespace[%s]已存在", req.Name)
	}
	// 处理错误情况
	if !k8serrors.IsNotFound(err) {
		// 如果不是 NotFound 错误，说明是其他错误
		if k8serrors.IsForbidden(err) {
			return fmt.Errorf("权限不足: %w", err)
		}
		return fmt.Errorf("查询namespace失败: %w", err)
	}
	nameSpace := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: req.Name,
			Annotations: map[string]string{
				DESCRIPTION: req.Description,
			},
		},
	}
	if req.Labels != nil {
		nameSpace.SetLabels(req.Labels)
	}
	if req.Annotations != nil {
		for k, v := range req.Annotations {
			nameSpace.Annotations[k] = v
		}
	}
	_, err = clientSet.CoreV1().Namespaces().Create(ctx, nameSpace, metav1.CreateOptions{DryRun: []string{metav1.DryRunAll}})
	if err != nil {
		return fmt.Errorf("DryRun参数校验失败: %w", err)
	}
	_, err = clientSet.CoreV1().Namespaces().Create(ctx, nameSpace, metav1.CreateOptions{})
	if err != nil {
		x.log.Warnf("创建命名空间失败: %v", err)
		return fmt.Errorf("创建命名空间失败: %w", err)
	}
	return nil
}

// DeleteNameSpace 删除NameSpace
func (x *NameSpaceUseCase) DeleteNameSpace(ctx context.Context, req *DeleteNameSpaceRequest) error {
	clientSet, err := x.cluster.GetClientSetByClusterName(ctx, req.ClusterName)
	if err != nil {
		return err
	}
	clusterItem, err := x.cluster.GetClusterByName(ctx, req.ClusterName)
	if err != nil {
		return err
	}
	err = clientSet.CoreV1().Namespaces().Delete(ctx, req.Name, metav1.DeleteOptions{})
	if err != nil {
		x.log.WithContext(ctx).Errorf("删除命名空间失败: %v", err)
		return fmt.Errorf("删除命名空间失败: %w", err)
	}
	err = x.pubNameSpaceDeletion.Publish(ctx, &NameSpaceDeletionEvent{
		ClusterId: clusterItem.ID,
		NameSpace: req.Name,
	})
	return err
}

func (x *NameSpaceUseCase) Init(ctx context.Context) context.CancelFunc {
	ctx, cancel := context.WithCancel(ctx)
	eg := errgroup.WithCancel(ctx)

	eg.Go(func(ctx context.Context) error {
		sub := xmsgbus.NewSubscriber((*NameSpaceDeletionEvent)(nil).Topic(), "NameSpaceUseCase",
			x.msgbus, x.otelOptions, x.tm,
			xmsgbus.WithHandleFunc(func(ctx context.Context, msg *NameSpaceDeletionEvent) error {
				// handle event
				return x.HandleNameSpaceDeletion(ctx, msg.ClusterId, msg.NameSpace)
			}),
		)
		defer sub.Close(ctx)
		for {
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
				if err := sub.Handle(ctx); err != nil {
					if errors.Is(err, context.DeadlineExceeded) {
						continue
					}
					x.log.Errorf("subcribe error: %v", err)
				}
			}
		}
	})
	return func() {
		cancel()
		_ = eg.Wait()
	}
}

// UpdateNameSpace 更新NameSpace
func (x *NameSpaceUseCase) UpdateNameSpace(ctx context.Context, req *UpdateNameSpaceRequest) error {
	clientSet, err := x.cluster.GetClientSetByClusterName(ctx, req.ClusterName)
	if err != nil {
		return err
	}

	// 获取集群命名空间
	namespaces, err := clientSet.CoreV1().Namespaces().Get(ctx, req.Name, metav1.GetOptions{})
	if k8serrors.IsNotFound(err) {
		return fmt.Errorf("命名空间不存在")
	}
	if err != nil {
		x.log.WithContext(ctx).Errorf("获取命名空间失败: %v", err)
		return fmt.Errorf("获取命名空间失败: %w", err)
	}
	if req.Labels != nil {
		namespaces.SetLabels(req.Labels)
	}
	annotations := make(map[string]string)
	if req.Annotations != nil {
		for k, v := range req.Annotations {
			annotations[k] = v
		}
	}
	if req.Description != "" {
		annotations[DESCRIPTION] = req.Description
	}
	namespaces.SetAnnotations(annotations)
	_, err = clientSet.CoreV1().Namespaces().Update(ctx, namespaces, metav1.UpdateOptions{DryRun: []string{metav1.DryRunAll}})
	if err != nil {
		return fmt.Errorf("DryRun参数校验失败: %w", err)
	}
	// 更新命名空间
	_, err = clientSet.CoreV1().Namespaces().Update(ctx, namespaces, metav1.UpdateOptions{})
	if err != nil {
		x.log.WithContext(ctx).Errorf("更新命名空间失败: %v", err)
		return fmt.Errorf("更新命名空间失败: %w", err)
	}
	return nil
}

func (x *NameSpaceUseCase) CreateNameSpaceByYaml(ctx context.Context, req *CreateNameSpaceByYamlRequest) error {
	clientSet, err := x.cluster.GetClientSetByClusterName(ctx, req.ClusterName)
	if err != nil {
		return err
	}
	var namespace corev1.Namespace
	decoder := yaml.NewYAMLOrJSONDecoder(strings.NewReader(req.Yaml), 4096)
	if err := decoder.Decode(&namespace); err != nil {
		return fmt.Errorf("解析YAML失败, 请检查格式是否正确: %v", err)
	}
	_, err = clientSet.CoreV1().Namespaces().Create(ctx, &namespace, metav1.CreateOptions{})
	if err != nil {
		if k8serrors.IsAlreadyExists(err) {
			return fmt.Errorf("命名空间[%s]已存在", namespace.Name)
		} else if k8serrors.IsForbidden(err) {
			return fmt.Errorf("权限不足: %w", err)
		}
		x.log.WithContext(ctx).Errorf("创建命名空间失败: %v", err)
		return fmt.Errorf("创建命名空间失败: %w", err)
	}
	return nil
}

func (x *NameSpaceUseCase) UpdateNameSpaceByYaml(ctx context.Context, req *CreateNameSpaceByYamlRequest) error {
	clientSet, err := x.cluster.GetClientSetByClusterName(ctx, req.ClusterName)
	if err != nil {
		return err
	}
	var namespace corev1.Namespace
	decoder := yaml.NewYAMLOrJSONDecoder(strings.NewReader(req.Yaml), 4096)
	if err := decoder.Decode(&namespace); err != nil {
		return fmt.Errorf("解析YAML失败, 请检查格式是否正确: %v", err)
	}

	// 获取当前的命名空间对象
	currentNamespace, err := clientSet.CoreV1().Namespaces().Get(ctx, namespace.Name, metav1.GetOptions{})
	if err != nil {
		x.log.WithContext(ctx).Errorf("获取当前命名空间失败: %v", err)
		if k8serrors.IsNotFound(err) {
			return fmt.Errorf("命名空间不存在")
		} else if k8serrors.IsForbidden(err) {
			return fmt.Errorf("权限不足: %w", err)
		}
		return fmt.Errorf("获取当前命名空间失败: %w", err)
	}

	namespace.ResourceVersion = currentNamespace.ResourceVersion
	_, err = clientSet.CoreV1().Namespaces().Update(ctx, &namespace, metav1.UpdateOptions{})
	if err != nil {
		x.log.WithContext(ctx).Errorf("更新命名空间失败: %v", err)
		return fmt.Errorf("更新命名空间失败: %w", err)
	}
	return nil
}

func (x *NameSpaceUseCase) GetNameSpace(ctx context.Context, req *GetNamespaceRequest) (*corev1.Namespace, error) {
	clientSet, err := x.cluster.GetClientSetByClusterName(ctx, req.ClusterName)
	if err != nil {
		return nil, err
	}
	// 获取集群命名空间
	namespace, err := clientSet.CoreV1().Namespaces().Get(ctx, req.Name, metav1.GetOptions{})
	if err != nil {
		x.log.WithContext(ctx).Errorf("查询命名空间失败: %v", err)
		return nil, fmt.Errorf("查询命名空间失败: %w", err)
	}
	return namespace, nil
}

func (x *NameSpaceUseCase) HandleNameSpaceDeletion(ctx context.Context, ClusterId uint32, namespace string) error {
	return x.roleBinding.DeleteByNamespace(ctx, ClusterId, namespace)

}
