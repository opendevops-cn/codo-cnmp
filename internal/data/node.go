package data

import (
	"codo-cnmp/internal/biz"
	"codo-cnmp/internal/model/dao"
	"codo-cnmp/internal/model/model/entity"
	"context"
	"encoding/json"
	"github.com/ccheers/xpkg/generic/arrayx"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
	corev1 "k8s.io/api/core/v1"
)

type NodeRepo struct {
	data *Data
}

func NewNodeRepo(data *Data) *NodeRepo {
	return &NodeRepo{
		data: data,
	}
}

func NewINodeRepo(repo *NodeRepo) biz.NodeRepo {
	return repo
}

func (n *NodeRepo) WithClusterID(id uint32) Option {
	return func(db *gdb.Model) *gdb.Model {
		return db.Where(dao.Cluster.Columns().Id, id)
	}
}

func (n *NodeRepo) WithClusterName(name string) Option {
	return func(db *gdb.Model) *gdb.Model {
		return db.Where(dao.Cluster.Columns().Name, name)
	}
}

func (n *NodeRepo) GetNode(ctx context.Context, req *biz.ListNodeRequest) (*biz.NodeItem, error) {
	db := dao.Node.Ctx(ctx)
	if req.ID != 0 {
		db = db.Where(dao.Node.Columns().Id, req.ID)
	}
	if req.Name != "" {
		db = db.Where(dao.Node.Columns().Name, req.Name)
	}
	var node entity.Node
	err := db.Scan(&node)
	if err != nil {
		return nil, err
	}
	return n.convertVO2DO(&node), nil
}

func (n *NodeRepo) convertVO2DO(vo *entity.Node) *biz.NodeItem {
	var (
		conditions      []corev1.NodeCondition
		capacity        corev1.ResourceList
		allocatable     corev1.ResourceList
		addresses       []corev1.NodeAddress
		labels          map[string]string
		annotations     map[string]string
		roles           []string
		spec            corev1.NodeSpec
		nodeInfo        corev1.NodeSystemInfo
		nodeHealthState []biz.NodeHealthState
	)

	json.Unmarshal([]byte(vo.Conditions), &conditions)
	json.Unmarshal([]byte(vo.Capacity), &capacity)
	json.Unmarshal([]byte(vo.Allocatable), &allocatable)
	json.Unmarshal([]byte(vo.Addresses), &addresses)
	json.Unmarshal([]byte(vo.NodeInfo), &nodeInfo)
	json.Unmarshal([]byte(vo.Labels), &labels)
	json.Unmarshal([]byte(vo.Annotations), &annotations)
	json.Unmarshal([]byte(vo.Roles), &roles)
	json.Unmarshal([]byte(vo.Spec), &spec)
	json.Unmarshal([]byte(vo.HealthState), &nodeHealthState)
	return &biz.NodeItem{
		ID:                uint32(vo.Id),
		Name:              vo.Name,
		Conditions:        conditions,
		Capacity:          capacity,
		Allocatable:       allocatable,
		Addresses:         addresses,
		CreationTimestamp: vo.CreationTimestamp,
		CpuUsage:          float32(vo.CpuUsage),
		MemoryUsage:       float32(vo.MemoryUsage),
		Status:            biz.NodeState(vo.Status),
		Labels:            labels,
		Roles:             roles,
		Annotations:       annotations,
		NodeInfo:          nodeInfo,
		UID:               vo.Uid,
		Spec:              spec,
		ResourceVersion:   vo.ResourceVersion,
		HealthState:       nodeHealthState,
	}
}

func (n *NodeRepo) convertDO2VO(data *biz.CreateNodeRequest) *entity.Node {
	conditionsJson, _ := json.Marshal(data.NodeItem.Conditions)
	capacityJson, _ := json.Marshal(data.NodeItem.Capacity)
	allocatableJson, _ := json.Marshal(data.NodeItem.Allocatable)
	AddressesJson, _ := json.Marshal(data.NodeItem.Addresses)
	return &entity.Node{
		ClusterId:         uint64(int32(data.ClusterID)),
		Name:              data.NodeItem.Name,
		Conditions:        string(conditionsJson),
		Capacity:          string(capacityJson),
		Allocatable:       string(allocatableJson),
		Addresses:         string(AddressesJson),
		CreationTimestamp: data.NodeItem.CreationTimestamp,
		CpuUsage:          float64(data.NodeItem.CpuUsage),
		MemoryUsage:       float64(data.NodeItem.MemoryUsage),
	}
}

func (n *NodeRepo) convertQuery(db *gdb.Model, req *biz.ListNodeRequest) *gdb.Model {
	if req.ClusterID != 0 {
		db = db.Where(dao.Node.Columns().ClusterId, int32(req.ClusterID))
	}
	if req.Keyword != "" {
		db = db.WhereLike(dao.Cluster.Columns().Name, "%"+req.Keyword+"%")
	}
	if req.ListAll {
		return db
	}
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}
	return db.Page(int(req.Page), int(req.PageSize))
}

func (n *NodeRepo) ListNodes(ctx context.Context, req *biz.ListNodeRequest) ([]*biz.NodeItem, error) {
	db := dao.Node.Ctx(ctx)
	db = n.convertQuery(db, req)
	var nodes []*entity.Node
	err := db.Scan(&nodes)
	if err != nil {
		return nil, err
	}
	return arrayx.Map(nodes, func(t *entity.Node) *biz.NodeItem {
		return n.convertVO2DO(t)
	}), nil
}

func (n *NodeRepo) CountNode(ctx context.Context, query *biz.ListNodeRequest) (uint32, error) {
	db := dao.Node.Ctx(ctx)
	db = n.convertQuery(db, query)
	count, err := db.Count()
	if err != nil {
		return 0, err
	}
	return uint32(count), nil
}

func (n *NodeRepo) CreateNode(ctx context.Context, req *biz.CreateNodeRequest) (uint32, error) {
	db := dao.Node.Ctx(ctx)
	vo := n.convertDO2VO(req)
	id, err := db.InsertAndGetId(vo)
	if err != nil {
		return 0, err
	}
	return uint32(id), nil
}

func isNodeSpecEmpty(spec corev1.NodeSpec) bool {
	return spec.PodCIDR == "" &&
		(spec.PodCIDRs == nil || len(spec.PodCIDRs) == 0) &&
		(spec.Taints == nil || len(spec.Taints) == 0) &&
		spec.ProviderID == "" && !spec.Unschedulable
}

func (n *NodeRepo) DeleteNode(ctx context.Context, id uint32) error {
	//TODO implement me
	panic("implement me")
}

func (n *NodeRepo) UpdateNode(ctx context.Context, req *biz.NodeItem) error {
	db := dao.Node.Ctx(ctx)
	updates := make(g.Map)
	if req.Conditions != nil {
		conditionsJson, _ := json.Marshal(req.Conditions)
		updates[dao.Node.Columns().Conditions] = string(conditionsJson)
	}
	if req.Capacity != nil {
		capacityJson, _ := json.Marshal(req.Capacity)
		updates[dao.Node.Columns().Capacity] = string(capacityJson)
	}
	if req.Allocatable != nil {
		allocatableJson, _ := json.Marshal(req.Allocatable)
		updates[dao.Node.Columns().Allocatable] = string(allocatableJson)
	}
	// todo 更新逻辑优化
	if req.CpuUsage != 0.0 {
		updates[dao.Node.Columns().CpuUsage] = req.CpuUsage
	}
	if req.MemoryUsage != 0.0 {
		updates[dao.Node.Columns().MemoryUsage] = req.MemoryUsage
	}
	if req.Status != 0 {
		updates[dao.Node.Columns().Status] = req.Status
	}
	if req.Labels != nil {
		updates[dao.Node.Columns().Labels] = req.Labels
	}
	if req.Annotations != nil {
		updates[dao.Node.Columns().Annotations] = req.Annotations
	}
	if req.Roles != nil {
		updates[dao.Node.Columns().Roles] = req.Roles
	}
	if req.UID != "" {
		updates[dao.Node.Columns().Uid] = req.UID
	}
	if !isNodeSpecEmpty(req.Spec) {
		specJson, _ := json.Marshal(req.Spec)
		updates[dao.Node.Columns().Spec] = string(specJson)
	}
	if req.ResourceVersion != "" {
		updates[dao.Node.Columns().ResourceVersion] = req.ResourceVersion
	}
	if req.NodeInfo.SystemUUID != "" {
		nodeInfoJson, _ := json.Marshal(req.NodeInfo)
		updates[dao.Node.Columns().NodeInfo] = string(nodeInfoJson)
	}
	if req.HealthState != nil {
		healthStateJson, _ := json.Marshal(req.HealthState)
		updates[dao.Node.Columns().HealthState] = string(healthStateJson)
	}
	_, err := db.Where(dao.Node.Columns().Name, req.Name).Where(dao.Node.Columns().ClusterId, int32(req.ClusterID)).Update(updates)
	return err
}

func (n *NodeRepo) ExistNode(ctx context.Context, name string, clusterID uint32) (bool, error) {
	db := dao.Node.Ctx(ctx)
	count, err := db.Where(dao.Node.Columns().Name, name).Where(dao.Node.Columns().ClusterId, int32(clusterID)).Count()
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (n *NodeRepo) CreateOrUpdateNode(ctx context.Context, req *biz.NodeItem) (uint32, error) {
	exist, err := n.ExistNode(ctx, req.Name, req.ClusterID)
	if err != nil {
		return 0, err
	}
	if exist {
		return req.ID, n.UpdateNode(ctx, req)
	}
	return n.CreateNode(ctx, &biz.CreateNodeRequest{
		ClusterID: req.ClusterID,
		NodeItem: biz.NodeItem{
			Name:              req.Name,
			Conditions:        req.Conditions,
			Capacity:          req.Capacity,
			Allocatable:       req.Allocatable,
			Addresses:         req.Addresses,
			CreationTimestamp: req.CreationTimestamp,
			CpuUsage:          req.CpuUsage,
			MemoryUsage:       req.MemoryUsage,
			Status:            req.Status,
			Labels:            req.Labels,
			NodeInfo:          req.NodeInfo,
			UID:               req.UID,
			ResourceVersion:   req.ResourceVersion,
			Spec:              req.Spec,
			Annotations:       req.Annotations,
			HealthState:       req.HealthState,
		},
	})
}
