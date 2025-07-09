package data

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"reflect"

	"codo-cnmp/common/utils"
	"codo-cnmp/internal/biz"
	"codo-cnmp/internal/conf"
	"codo-cnmp/internal/model/dao"
	"codo-cnmp/internal/model/model/entity"
	"codo-cnmp/pb"
	"github.com/ccheers/xpkg/generic/arrayx"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

type ClusterRepo struct {
	data *Data
	log  *log.Helper
	bc   *conf.Bootstrap
}

func (x *ClusterRepo) UpdateClusterBasicInfo(ctx context.Context, req *biz.UpdateClusterBasicRequest) error {
	db := dao.Cluster.Ctx(ctx)
	//vo := x.convertDO2VO(data)
	updates := make(g.Map)
	var HealthState string
	if req.HealthState != nil && len(req.HealthState) > 0 {
		HealthStateBytes, err := json.Marshal(req.HealthState)
		if err != nil {
			x.log.Errorf("序列化HealthState失败: %v", err)
			return nil
		}
		HealthState = string(HealthStateBytes)
	}
	if HealthState != "" {
		updates[dao.Cluster.Columns().HealthState] = req.HealthState
	}
	if req.BuildDate != "" {
		updates[dao.Cluster.Columns().BuildDate] = req.BuildDate
	}
	if req.ServerVersion != "" {
		updates[dao.Cluster.Columns().ServerVersion] = req.ServerVersion
	}
	if req.CpuUsage != 0.0 {
		updates[dao.Cluster.Columns().CpuUsage] = req.CpuUsage
	}
	if req.Platform != "" {
		updates[dao.Cluster.Columns().Platform] = req.Platform
	}
	if req.CpuTotal != 0.0 {
		updates[dao.Cluster.Columns().CpuTotal] = req.CpuTotal
	}
	if req.MemoryTotal != 0.0 {
		updates[dao.Cluster.Columns().MemoryTotal] = req.MemoryTotal
	}
	if req.MemoryUsage != 0.0 {
		updates[dao.Cluster.Columns().MemoryUsage] = req.MemoryUsage
	}
	if req.NodeCount != 0 {
		updates[dao.Cluster.Columns().NodeCount] = req.NodeCount
	}
	if &req.ClusterState != nil {
		updates[dao.Cluster.Columns().ClusterState] = req.ClusterState
	}

	if len(updates) == 0 {
		return nil
	}
	_, err := db.Where(dao.Cluster.Columns().Id, req.Id).Update(updates)
	return err
}

func (x *ClusterRepo) UpdateClusterState(ctx context.Context, req *biz.UpdateClusterStateRequest) error {
	var HealthState string
	if req.HealthState != nil && len(req.HealthState) > 0 {
		HealthStateBytes, err := json.Marshal(req.HealthState)
		if err != nil {
			x.log.Errorf("序列化HealthState失败: %v", err)
			return nil
		}
		HealthState = string(HealthStateBytes)
	}
	db := dao.Cluster.Ctx(ctx)
	_, err := db.Where(dao.Cluster.Columns().Id, req.Id).Update(g.Map{
		dao.Cluster.Columns().ClusterState: int(req.ClusterState),
		dao.Cluster.Columns().HealthState:  HealthState,
	})
	if err != nil {
		return fmt.Errorf("更新集群状态失败: %v", err)
	}
	return nil
}

type Option func(*gdb.Model) *gdb.Model

func (x *ClusterRepo) WithClusterID(id uint32) Option {
	return func(db *gdb.Model) *gdb.Model {
		return db.Where(dao.Cluster.Columns().Id, id)
	}
}

func (x *ClusterRepo) WithClusterName(name string) Option {
	return func(db *gdb.Model) *gdb.Model {
		return db.Where(dao.Cluster.Columns().Name, name)
	}
}

func (x *ClusterRepo) GetCluster(ctx context.Context, options ...Option) (*biz.ClusterItem, error) {
	db := dao.Cluster.Ctx(ctx)
	for _, option := range options {
		db = option(db)
	}
	var cluster entity.Cluster
	err := db.Scan(&cluster)
	if err != nil {
		return nil, err
	}
	return x.convertVO2DO(&cluster), nil
}

func NewClusterRepo(data *Data, logger log.Logger, bc *conf.Bootstrap) *ClusterRepo {
	return &ClusterRepo{
		data: data,
		log:  log.NewHelper(logger),
		bc:   bc,
	}
}

func NewIClusterRepo(repo *ClusterRepo) biz.ClusterRepo {
	return repo
}

// EncryptK8sConfig 加密k8s配置
func (x *ClusterRepo) EncryptK8sConfig(jsonStr string) (string, error) {
	// 编码为Base64
	base64KubeConfig := base64.StdEncoding.EncodeToString([]byte(jsonStr))

	// 对称加密
	key := x.bc.APP.GetSECRET()
	if key == "" {
		return "", fmt.Errorf("配置文件中secret key为空")
	}
	encrypted, err := utils.AESEncrypt(base64KubeConfig, key)
	if err != nil {
		return "", fmt.Errorf("加密失败: %v", err)
	}
	return encrypted, nil
}

// DecryptK8sConfig 解密k8s配置
func (x *ClusterRepo) DecryptK8sConfig(jsonStr string) (string, error) {
	// 解密
	key := x.bc.APP.GetSECRET()
	if key == "" {
		return "", fmt.Errorf("配置文件中secret key为空")
	}
	decrypted, err := utils.AESDecrypt(jsonStr, key)
	if err != nil {
		return "", fmt.Errorf("解密失败: %v", err)
	}

	// 解码为Base64
	base64KubeConfig, err := base64.StdEncoding.DecodeString(decrypted)
	if err != nil {
		return "", fmt.Errorf("解码失败: %v", err)
	}
	return string(base64KubeConfig), nil
}

// convertQuery convert query to ORM query
func convertQuery(db *gdb.Model, query *biz.QueryClusterReq) *gdb.Model {
	if query.ID > 0 {
		db = db.Where(dao.Cluster.Columns().Id, query.ID)
	}
	if query.Keyword != "" {
		db = db.WhereLike(dao.Cluster.Columns().Name, "%"+query.Keyword+"%").WhereOrLike(dao.Cluster.Columns().Description, "%"+query.Keyword+"%")
	}
	if query.Page <= 0 {
		query.Page = 1
	}
	if query.PageSize <= 0 {
		query.PageSize = 10
	}
	if query.ListAll {
		return db
	}
	return db.Page(int(query.Page), int(query.PageSize))
}

// convertCountQuery convert query to count query
func convertCountQuery(db *gdb.Model, query *biz.QueryClusterReq) *gdb.Model {
	if query.ID > 0 {
		db = db.Where(dao.Cluster.Columns().Id, query.ID)
	}
	if query.Keyword != "" {
		db = db.WhereLike(dao.Cluster.Columns().Name, "%"+query.Keyword+"%").WhereLike(dao.Cluster.Columns().Description, "%"+query.Keyword+"%")
	}
	return db.Page(int(query.Page), int(query.PageSize))
}

func (x *ClusterRepo) convertVO2DO(cluster *entity.Cluster) *biz.ClusterItem {
	var (
		ImportDetail biz.ImportDetail
		healthState  []biz.HealthState
		ops          []string
		links        []*pb.Link
	)
	// 解密k8s配置
	decrypted, err := x.DecryptK8sConfig(cluster.ImportDetail)
	if err != nil {
		x.log.Errorf("解密k8s配置失败: %v", err)
		return nil
	}
	err = json.Unmarshal([]byte(decrypted), &ImportDetail)
	if err != nil {
		x.log.Errorf("解码ImportDetail失败: %v", err)
		return nil
	}
	if cluster.HealthState != "" {
		err = json.Unmarshal([]byte(cluster.HealthState), &healthState)
		if err != nil {
			x.log.Errorf("解码HealthState失败: %v", err)
			return nil
		}
	}

	if cluster.Ops != "" {
		err = json.Unmarshal([]byte(cluster.Ops), &ops)
		if err != nil {
			x.log.Errorf("解码Ops失败: %v", err)
			return nil
		}
	}
	if cluster.Links != "" {
		err = json.Unmarshal([]byte(cluster.Links), &links)
		if err != nil {
			x.log.Errorf("解码Links失败: %v", err)
			return nil
		}
	}

	return &biz.ClusterItem{
		ID:            uint32(cluster.Id),
		Name:          cluster.Name,
		Description:   cluster.Description,
		ClusterState:  biz.ClusterState(cluster.ClusterState),
		HealthState:   healthState,
		BuildDate:     cluster.BuildDate,
		ServerVersion: cluster.ServerVersion,
		CpuUsage:      float32(cluster.CpuUsage),
		Platform:      cluster.Platform,
		CpuTotal:      float32(cluster.CpuTotal),
		MemoryTotal:   float32(cluster.MemoryTotal),
		MemoryUsage:   float32(cluster.MemoryUsage),
		ImportType:    biz.ImportType(cluster.ImportType),
		ImportDetail:  ImportDetail,
		NodeCount:     int(cluster.NodeCount),
		UID:           cluster.Uid,
		IDIP:          cluster.Idip,
		AppId:         cluster.AppId,
		AppSecret:     cluster.AppSecret,
		Ops:           ops,
		DstAgentId:    uint32(cluster.DstAgentId),
		ConnectType:   pb.ConnectType(cluster.ConnetType),
		MeshAddr:      cluster.MeshAddr,
		Links:         links,
	}
}

// 更新JSON序列化字段
func updateJSONField[T any](updates *g.Map, column string, value T, allowEmpty bool) {
	// 检查值是否有效
	if !reflect.ValueOf(value).IsValid() {
		return
	}

	reflectValue := reflect.ValueOf(value)
	reflectType := reflect.TypeOf(value)

	if reflectType.Kind() == reflect.Slice && reflectValue.Len() == 0 && !allowEmpty {
		return
	}
	// 特殊处理切片类型
	if reflectType.Kind() == reflect.Slice {
		// 如果切片长度为0,设置为null
		if reflectValue.Len() == 0 {
			(*updates)[column] = nil
			return
		}
		// 非空切片正常序列化
		jsonBytes, err := json.Marshal(value)
		if err == nil {
			(*updates)[column] = string(jsonBytes)
		}
	} else if !reflectValue.IsZero() {
		// 非切片类型，且不为零值
		jsonBytes, err := json.Marshal(value)
		if err == nil && len(jsonBytes) > 0 && string(jsonBytes) != "null" {
			(*updates)[column] = string(jsonBytes)
		}
	}
}

// 更新数值字段，仅当新值不为零时
func updateIfNonZero[T comparable](updates *g.Map, column string, value T) {
	var zero T
	if value != zero {
		(*updates)[column] = value
	}
}

// 更新字符串字段，仅当新值不为空且与旧值不同时
func updateIfChanged(updates *g.Map, column string, newValue, oldValue string, allowEmpty bool) {
	if allowEmpty {
		if newValue != oldValue {
			(*updates)[column] = newValue
		}
	}
	if newValue != "" && newValue != oldValue {
		(*updates)[column] = newValue
	}
}

func (x *ClusterRepo) UpdateClusterV2(ctx context.Context, data *biz.ClusterItem) error {
	currentCluster, err := x.GetClusterByID(ctx, data.ID)
	if err != nil {
		return err
	}

	// 创建更新字段映射
	updates := make(g.Map)

	// 新cluster的名称不能与其他cluster的名称重复
	if data.Name != "" && data.Name != currentCluster.Name {
		exist, err := x.ExistCluster(ctx, data.Name)
		if err != nil {
			return err
		}
		if exist {
			return fmt.Errorf("集群名称已存在")
		}
	}

	// 处理常规字段
	updateIfChanged(&updates, dao.Cluster.Columns().Name, data.Name, currentCluster.Name, false)
	updateIfChanged(&updates, dao.Cluster.Columns().Description, data.Description, currentCluster.Description, true)
	updateIfChanged(&updates, dao.Cluster.Columns().Idip, data.IDIP, currentCluster.IDIP, true)
	updateIfChanged(&updates, dao.Cluster.Columns().AppId, data.AppId, currentCluster.AppId, true)
	updateIfChanged(&updates, dao.Cluster.Columns().MeshAddr, data.MeshAddr, currentCluster.MeshAddr, true)
	updateIfNonZero(&updates, dao.Cluster.Columns().ConnetType, int(data.ConnectType))
	updateIfNonZero(&updates, dao.Cluster.Columns().DstAgentId, int(data.DstAgentId))
	updateIfNonZero(&updates, dao.Cluster.Columns().ConnetType, int(data.ConnectType))

	// 处理特殊字段: ImportDetail
	if data.ImportDetail.KubeConfig != "" && data.ImportDetail.ApiServer != "" {
		importDetail, err := json.Marshal(data.ImportDetail)
		if err == nil {
			encryptedString, err := x.EncryptK8sConfig(string(importDetail))
			if err == nil {
				updates[dao.Cluster.Columns().ImportDetail] = encryptedString
			}
		}
	}

	// 处理特殊字段: AppSecret
	if data.AppSecret != "" && data.AppSecret != currentCluster.AppSecret {
		key := x.bc.APP.GetSECRET()
		if key != "" {
			encryptedAppSecret, err := utils.AESEncrypt(data.AppSecret, key)
			if err == nil {
				updates[dao.Cluster.Columns().AppSecret] = encryptedAppSecret
			}
		}
	}

	if data.Ops != nil {
		updateJSONField(&updates, dao.Cluster.Columns().Ops, data.Ops, true)
	}
	if data.Links != nil {
		updateJSONField(&updates, dao.Cluster.Columns().Links, data.Links, true)
	}

	// 如果没有需要更新的字段，直接返回
	if len(updates) == 0 {
		return nil
	}

	// 执行更新操作
	db := dao.Cluster.Ctx(ctx)
	_, err = db.Where(dao.Cluster.Columns().Id, data.ID).Update(updates)
	return err
}

func (x *ClusterRepo) convertDO2VO(data *biz.ClusterItem) *entity.Cluster {
	importDetail, err := json.Marshal(data.ImportDetail)
	if err != nil {
		x.log.Errorf("序列化ImportDetail失败: %v", err)
		return nil
	}
	var HealthState string
	if data.HealthState != nil && len(data.HealthState) > 0 {
		HealthStateBytes, err := json.Marshal(data.HealthState)
		if err != nil {
			x.log.Errorf("序列化HealthState失败: %v", err)
			return nil
		}
		HealthState = string(HealthStateBytes)
	}
	var ops string
	if data.Ops != nil && len(data.Ops) > 0 {
		opsBytes, err := json.Marshal(data.Ops)
		if err != nil {
			x.log.Errorf("序列化Ops失败: %v", err)
			return nil
		}
		ops = string(opsBytes)
	}
	// 加密k8s配置
	var encryptedString string
	if string(importDetail) != "" {
		encryptedString, err = x.EncryptK8sConfig(string(importDetail))
		if err != nil {
			x.log.Errorf("加密k8s配置失败: %v", err)
			return nil
		}
	}
	var Links string
	if data.Links != nil && len(data.Links) > 0 {
		LinksBytes, err := json.Marshal(data.Links)
		if err != nil {
			x.log.Errorf("序列化Links失败: %v", err)
			return nil
		}
		Links = string(LinksBytes)
	}
	return &entity.Cluster{
		Id:            uint64(data.ID),
		Name:          data.Name,
		Description:   data.Description,
		ImportType:    int(data.ImportType),
		ImportDetail:  encryptedString,
		ClusterState:  int(data.ClusterState),
		HealthState:   HealthState,
		BuildDate:     data.BuildDate,
		ServerVersion: data.ServerVersion,
		CpuUsage:      float64(data.CpuUsage),
		Platform:      data.Platform,
		CpuTotal:      float64(data.CpuTotal),
		MemoryTotal:   float64(data.MemoryTotal),
		MemoryUsage:   float64(data.MemoryUsage),
		NodeCount:     uint(data.NodeCount),
		Uid:           data.UID,
		Idip:          data.IDIP,
		AppId:         data.AppId,
		Ops:           ops,
		DstAgentId:    int64(data.DstAgentId),
		ConnetType:    uint(data.ConnectType),
		MeshAddr:      data.MeshAddr,
		Links:         Links,
	}
}

// ListClusters list clusters
func (x *ClusterRepo) ListClusters(ctx context.Context, query *biz.QueryClusterReq) (biz.ClusterItems, error) {
	db := dao.Cluster.Ctx(ctx)
	db = convertQuery(db, query)
	var clusters []*entity.Cluster
	err := db.Scan(&clusters)
	if err != nil {
		return nil, err
	}
	return arrayx.Map(clusters, func(t *entity.Cluster) *biz.ClusterItem {
		return x.convertVO2DO(t)
	}), nil
}

// CreateCluster create cluster
func (x *ClusterRepo) CreateCluster(ctx context.Context, data *biz.ClusterItem) (uint32, error) {
	db := dao.Cluster.Ctx(ctx)
	vo := x.convertDO2VO(data)
	appSecret := data.AppSecret
	if appSecret != "" {
		key := x.bc.APP.GetSECRET()
		if key == "" {
			return 0, fmt.Errorf("配置文件中secret key为空")
		}
		encryptedAppSecretString, err := utils.AESEncrypt(data.AppSecret, key)
		if err != nil {
			return 0, fmt.Errorf("加密AppSecret失败: %v", err)
		}
		vo.AppSecret = encryptedAppSecretString
	}
	id, err := db.InsertAndGetId(vo)
	if err != nil {
		return 0, err
	}
	return uint32(id), nil
}

// ExistCluster check cluster exist
func (x *ClusterRepo) ExistCluster(ctx context.Context, name string) (bool, error) {
	db := dao.Cluster.Ctx(ctx)
	count, err := db.Where(dao.Cluster.Columns().Name, name).Count()
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// CountCluster CountClusters count clusters
func (x *ClusterRepo) CountCluster(ctx context.Context, query *biz.QueryClusterReq) (uint32, error) {
	db := dao.Cluster.Ctx(ctx)
	db = convertCountQuery(db, query)
	count, err := db.Count()
	if err != nil {
		return 0, err
	}
	return uint32(count), nil
}

// DeleteCluster delete cluster
func (x *ClusterRepo) DeleteCluster(ctx context.Context, id uint32) error {
	db := dao.Cluster.Ctx(ctx)
	_, err := db.Where(dao.Cluster.Columns().Id, id).Delete()
	if err != nil {
		return err
	}
	return nil
}

// FetchAllClusters fetch all clusters
func (x *ClusterRepo) FetchAllClusters(ctx context.Context) (biz.ClusterItems, error) {
	db := dao.Cluster.Ctx(ctx)
	var clusters []*entity.Cluster
	err := db.Scan(&clusters)
	if err != nil {
		return nil, err
	}
	return arrayx.Map(clusters, func(t *entity.Cluster) *biz.ClusterItem {
		return x.convertVO2DO(t)
	}), nil
}

func (x *ClusterRepo) UpdateCluster(ctx context.Context, data *biz.ClusterItem) error {
	db := dao.Cluster.Ctx(ctx)
	vo := x.convertDO2VO(data)
	if data.ImportDetail.KubeConfig == "" || data.ImportDetail.Token == "" {
		vo.ImportDetail = ""
	}
	currenCluster, err := x.GetClusterByID(ctx, data.ID)
	if err != nil {
		return err
	}
	if data.AppSecret != "" && data.AppSecret != currenCluster.AppSecret {
		key := x.bc.APP.GetSECRET()
		if key == "" {
			return fmt.Errorf("配置文件中secret key为空")
		}
		encryptedAppSecretString, err := utils.AESEncrypt(data.AppSecret, key)
		if err != nil {
			return fmt.Errorf("加密AppSecret失败: %v", err)
		}
		vo.AppSecret = encryptedAppSecretString
	}
	updates := make(g.Map)
	if vo.AppId != "" {
		updates[dao.Cluster.Columns().AppId] = vo.AppId
	}
	if vo.AppSecret != "" {
		updates[dao.Cluster.Columns().AppSecret] = vo.AppSecret
	}
	if vo.Ops != "" {
		updates[dao.Cluster.Columns().Ops] = vo.Ops
	}
	if vo.Description != "" {
		updates[dao.Cluster.Columns().Description] = vo.Description
	}
	if vo.Idip != "" {
		updates[dao.Cluster.Columns().Idip] = vo.Idip
	}
	if vo.Name != "" {
		updates[dao.Cluster.Columns().Name] = vo.Name
	}
	if &vo.NodeState != nil {
		updates[dao.Cluster.Columns().NodeState] = vo.NodeState
	}
	if vo.HealthState != "" {
		updates[dao.Cluster.Columns().HealthState] = vo.HealthState
	}
	if vo.BuildDate != "" {
		updates[dao.Cluster.Columns().BuildDate] = vo.BuildDate
	}
	if vo.ServerVersion != "" {
		updates[dao.Cluster.Columns().ServerVersion] = vo.ServerVersion
	}
	if vo.CpuUsage != 0.0 {
		updates[dao.Cluster.Columns().CpuUsage] = vo.CpuUsage
	}
	if vo.Platform != "" {
		updates[dao.Cluster.Columns().Platform] = vo.Platform
	}
	if vo.CpuTotal != 0.0 {
		updates[dao.Cluster.Columns().CpuTotal] = vo.CpuTotal
	}
	if vo.MemoryTotal != 0.0 {
		updates[dao.Cluster.Columns().MemoryTotal] = vo.MemoryTotal
	}
	if vo.MemoryUsage != 0.0 {
		updates[dao.Cluster.Columns().MemoryUsage] = vo.MemoryUsage
	}
	if vo.NodeCount != 0 {
		updates[dao.Cluster.Columns().NodeCount] = vo.NodeCount
	}
	if &vo.ClusterState != nil {
		updates[dao.Cluster.Columns().ClusterState] = vo.ClusterState
	}
	if vo.ImportDetail != "" {
		updates[dao.Cluster.Columns().ImportDetail] = vo.ImportDetail
	}
	if vo.Links != "" {
		updates[dao.Cluster.Columns().Links] = vo.Links
	}

	if len(updates) == 0 {
		return nil
	}
	_, err = db.Where(dao.Cluster.Columns().Id, data.ID).Update(updates)
	return err
}

// GetClusterByID get cluster by id
func (x *ClusterRepo) GetClusterByID(ctx context.Context, id uint32) (*biz.ClusterItem, error) {
	db := dao.Cluster.Ctx(ctx)
	var cluster entity.Cluster
	err := db.Where(dao.Cluster.Columns().Id, id).Scan(&cluster)
	if err != nil {
		return nil, err
	}
	if &cluster == nil {
		return nil, fmt.Errorf("集群不存在")
	}
	return x.convertVO2DO(&cluster), nil
}

// GetClusterByName get cluster by name
func (x *ClusterRepo) GetClusterByName(ctx context.Context, name string) (*biz.ClusterItem, error) {
	db := dao.Cluster.Ctx(ctx)
	var cluster entity.Cluster
	err := db.Where(dao.Cluster.Columns().Name, name).Scan(&cluster)
	if err != nil {
		return nil, err
	}
	if &cluster == nil {
		return nil, fmt.Errorf("集群不存在")
	}
	return x.convertVO2DO(&cluster), nil
}
