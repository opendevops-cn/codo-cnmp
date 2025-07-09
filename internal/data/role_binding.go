/* Package data Description: 角色绑定数据仓库.
角色-用户组-集群-命名空间多对多绑定关系, 每一条记录代表一个用户组在某个集群的某个命名空间下拥有某个角色的权限.
*/

package data

import (
	"codo-cnmp/internal/biz"
	"codo-cnmp/internal/model/dao"
	"codo-cnmp/internal/model/model/entity"
	"context"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

type RoleBindingRepo struct {
	data *Data
	log  *log.Helper
}

func (x *RoleBindingRepo) DeleteByNamespace(ctx context.Context, clusterID uint32, namespace string) error {
	model := dao.RoleBinding.Ctx(ctx)
	_, err := model.Where(dao.RoleBinding.Columns().ClusterId, clusterID).Where(dao.RoleBinding.Columns().Namespace, namespace).Delete()
	return err
}

func (x *RoleBindingRepo) DeleteByClusterId(ctx context.Context, clusterID uint32) error {
	model := dao.RoleBinding.Ctx(ctx)
	_, err := model.Where(dao.RoleBinding.Columns().ClusterId, clusterID).Delete()
	return err
}

func (x *RoleBindingRepo) BulkUpdate(ctx context.Context, data []*biz.RoleBindingItem) error {
	return x.BulkCreate(ctx, data)
}

func (x *RoleBindingRepo) BulkCreate(ctx context.Context, data []*biz.RoleBindingItem) error {
	db := dao.RoleBinding.Ctx(ctx)
	entities := make([]g.Map, 0)
	// 开启事务
	err := db.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		for _, item := range data {
			entities = append(entities, g.Map{
				dao.RoleBinding.Columns().UserGroupId: item.UserGroupID,
				dao.RoleBinding.Columns().ClusterId:   item.ClusterID,
				dao.RoleBinding.Columns().RoleId:      item.RoleID,
				dao.RoleBinding.Columns().Namespace:   item.Namespace,
			})
		}
		_, err := tx.Save(dao.RoleBinding.Table(), entities)
		return err
	})
	return err

}

/* BatchUpdateByUserGroupID 批量更新role binding
* @param ctx context.Context
* @param data []*biz.RoleBindingItem
* @return error
* @note 批量更新role binding
 */

func (x *RoleBindingRepo) BatchUpdateByUserGroupID(ctx context.Context, data []*biz.RoleBindingItem) error {
	db := dao.RoleBinding.Ctx(ctx)
	var items []*entity.RoleBinding
	// 开启事务
	err := db.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		for _, item := range data {
			items = append(items, x.convertDO2VO(item))
		}
		_, err := tx.Save(dao.RoleBinding.Table(), items)
		return err
	})
	return err
}

func (x *RoleBindingRepo) List(ctx context.Context, query *biz.ListRoleBindingRequest) ([]*biz.RoleBindingItem, error) {
	db := dao.RoleBinding.Ctx(ctx)
	db = x.Apply(db, x.WithUserGroupID(query.UserGroupID), x.WithClusterID(query.ClusterID), x.WithRoleID(query.RoleID), x.WithPage(query.Page, query.PageSize, query.ListAll))
	var RoleBindings []*entity.RoleBinding
	err := db.Scan(&RoleBindings)
	if err != nil {
		return nil, err
	}
	var result []*biz.RoleBindingItem
	for _, item := range RoleBindings {
		result = append(result, x.convertVO2DO(item))
	}
	return result, nil
}

// ListByUserGroupID 根据用户组ID查询角色绑定关系 聚合查询
func (x *RoleBindingRepo) ListByUserGroupID(ctx context.Context, query *biz.ListRoleBindingRequest) ([]*biz.ListGrantedUserGroupResponseItem, uint32, error) {
	results := make([]*biz.ListGrantedUserGroupResponseItem, 0)
	db := dao.RoleBinding.Ctx(ctx)
	//  使用原生 SQL 查询
	db = db.Raw("SELECT user_group_id, " +
		"MAX(updated_at) AS update_time, " +
		"COUNT(DISTINCT cluster_id) AS cluster_count," +
		" COUNT(DISTINCT role_id) AS role_count " +
		"FROM role_binding " +
		"WHERE deleted_at IS NULL" +
		" GROUP BY user_group_id " +
		"ORDER BY user_group_id ASC")

	total, err := db.Count()
	if err != nil {
		return results, 0, err
	}
	if !query.ListAll {
		if query.Page <= 0 {
			query.Page = 1
		}
		if query.PageSize <= 0 {
			query.PageSize = 10
		}
		db = db.Page(int(query.Page), int(query.PageSize))
	}
	err = db.Scan(&results)
	if err != nil {
		return results, 0, err
	}
	return results, uint32(total), nil
}

// ListByRoleID 根据角色ID查询角色绑定关系 聚合查询
func (x *RoleBindingRepo) ListByRoleID(ctx context.Context, query *biz.ListRoleBindingRequest) ([]*biz.RoleBindingItem, uint32, error) {
	results := make([]*biz.RoleBindingItem, 0)
	db := dao.RoleBinding.Ctx(ctx)
	db = db.Where(dao.RoleBinding.Columns().RoleId, query.RoleID).OrderAsc(dao.RoleBinding.Columns().UserGroupId)
	total, err := db.Count()
	if err != nil {
		return results, 0, err
	}
	if !query.ListAll {
		if query.Page <= 0 {
			query.Page = 1
		}
		if query.PageSize <= 0 {
			query.PageSize = 10
		}
		db = db.Page(int(query.Page), int(query.PageSize))
	}
	err = db.Scan(&results)
	if err != nil {
		return results, 0, err
	}
	return results, uint32(total), nil
}

func (x *RoleBindingRepo) Count(ctx context.Context, query *biz.ListRoleBindingRequest) (uint32, error) {
	db := dao.RoleBinding.Ctx(ctx)
	db = x.Apply(db, x.WithUserGroupID(query.UserGroupID), x.WithClusterID(query.ClusterID), x.WithRoleID(query.RoleID))
	count, err := db.Count()
	if err != nil {
		return 0, err
	}
	return uint32(count), nil
}

// Exists 检查是否存在
func (x *RoleBindingRepo) Exists(ctx context.Context, userGroupId uint32, clusterId uint32, roleId uint32, namespace string) (bool, error) {
	db := dao.RoleBinding.Ctx(ctx)
	count, err := db.Where(dao.RoleBinding.Columns().UserGroupId, userGroupId).
		Where(dao.RoleBinding.Columns().ClusterId, clusterId).
		Where(dao.RoleBinding.Columns().RoleId, roleId).
		Where(dao.RoleBinding.Columns().Namespace, namespace).Count()
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (x *RoleBindingRepo) Create(ctx context.Context, item *biz.RoleBindingItem) error {
	db := dao.RoleBinding.Ctx(ctx)
	_, err := db.Insert(x.convertDO2VO(item))
	return err
}

func (x *RoleBindingRepo) Update(ctx context.Context, item *biz.RoleBindingItem) error {
	db := dao.RoleBinding.Ctx(ctx)
	_, err := db.Where(dao.RoleBinding.Columns().UserGroupId, item.UserGroupID).Update(x.convertDO2VO(item))
	return err
}

func (x *RoleBindingRepo) DeleteByUserGroupID(ctx context.Context, userGroupID uint32) error {
	db := dao.RoleBinding.Ctx(ctx)
	_, err := db.Where(dao.RoleBinding.Columns().UserGroupId, userGroupID).Delete()
	return err
}

func NewRoleBindingRepo(data *Data, logger log.Logger) *RoleBindingRepo {
	return &RoleBindingRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func NewIRoleBindingRepoRepo(x *RoleBindingRepo) biz.IRoleBindingRepo {
	return x
}

type UserGroupClusterRefQueryOption func(*gdb.Model) *gdb.Model

func (x *RoleBindingRepo) WithUserGroupID(userGroupId uint32) UserGroupClusterRefQueryOption {
	return func(db *gdb.Model) *gdb.Model {
		if userGroupId == 0 {
			return db
		}
		return db.Where(dao.RoleBinding.Columns().UserGroupId, userGroupId)
	}
}

func (x *RoleBindingRepo) WithClusterID(clusterId uint32) UserGroupClusterRefQueryOption {
	return func(db *gdb.Model) *gdb.Model {
		if clusterId == 0 {
			return db
		}
		return db.Where(dao.RoleBinding.Columns().ClusterId, clusterId)
	}
}

func (x *RoleBindingRepo) WithRoleID(roleId uint32) UserGroupClusterRefQueryOption {
	return func(db *gdb.Model) *gdb.Model {
		if roleId == 0 {
			return db
		}
		return db.Where(dao.RoleBinding.Columns().RoleId, roleId)
	}
}

func (x *RoleBindingRepo) WithNamespace(namespace string) UserGroupClusterRefQueryOption {
	return func(db *gdb.Model) *gdb.Model {
		if namespace == "" {
			return db
		}
		return db.Where(dao.RoleBinding.Columns().Namespace, namespace)
	}
}

func (x *RoleBindingRepo) WithPage(page, limit uint32, listAll bool) UserGroupClusterRefQueryOption {
	return func(db *gdb.Model) *gdb.Model {
		if listAll {
			return db
		}
		if page <= 0 {
			page = 1
		}
		if limit <= 0 {
			limit = 10
		}
		return db.Page(int(page), int(limit))
	}
}

func (x *RoleBindingRepo) Apply(db *gdb.Model, options ...UserGroupClusterRefQueryOption) *gdb.Model {
	for _, option := range options {
		db = option(db)
	}
	return db
}

func (x *RoleBindingRepo) convertDO2VO(data *biz.RoleBindingItem) *entity.RoleBinding {
	vo := &entity.RoleBinding{
		UserGroupId: uint64(data.UserGroupID),
		ClusterId:   uint64(data.ClusterID),
		RoleId:      uint64(data.RoleID),
		Namespace:   data.Namespace,
	}
	return vo
}

func (x *RoleBindingRepo) convertVO2DO(data *entity.RoleBinding) *biz.RoleBindingItem {
	vo := &biz.RoleBindingItem{
		RoleBindingCommonParams: biz.RoleBindingCommonParams{
			UserGroupID: uint32(data.UserGroupId),
			ClusterID:   uint32(data.ClusterId),
			RoleID:      uint32(data.RoleId),
		},
		Namespace: data.Namespace,
	}
	return vo
}

// calculateDiff 计算差异
func (x *RoleBindingRepo) calculateDiff(existingData []*entity.RoleBinding, newData []*biz.RoleBindingItem) (toAdd []*biz.RoleBindingItem, toRemove []*entity.RoleBinding) {
	existingMap := make(map[string]*entity.RoleBinding)
	newMap := make(map[string]*biz.RoleBindingItem)

	// 将现有数据放入map，方便比较
	for _, item := range existingData {
		key := x.generateKey(uint32(item.UserGroupId), uint32(item.ClusterId), uint32(item.RoleId), item.Namespace)
		existingMap[key] = item
	}

	// 将新数据放入map，方便比较
	for _, item := range newData {
		key := x.generateKey(item.UserGroupID, item.ClusterID, item.RoleID, item.Namespace)
		newMap[key] = item
	}

	// 找到需要删除的现有数据（新数据中没有）
	for key, item := range existingMap {
		if _, exists := newMap[key]; !exists {
			toRemove = append(toRemove, item)
		}
	}

	// 找到需要新增的新数据（现有数据中没有）
	for key, item := range newMap {
		if _, exists := existingMap[key]; !exists {
			toAdd = append(toAdd, item)
		}
	}

	return toAdd, toRemove
}

// generateKey 生成唯一键，用于比较差异
func (x *RoleBindingRepo) generateKey(userGroupID, clusterID, roleID uint32, namespace string) string {
	return fmt.Sprintf("%d-%d-%d-%s", userGroupID, clusterID, roleID, namespace)
}

// ManageRoleBindingByUserGroupID  管理角色绑定关系
func (x *RoleBindingRepo) ManageRoleBindingByUserGroupID(ctx context.Context, userGroupID uint32, data []*biz.RoleBindingItem) error {
	db := dao.RoleBinding.Ctx(ctx)
	// 开启事务
	err := db.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		// 删除现有记录
		_, err := tx.Delete(dao.RoleBinding.Table(), g.Map{
			dao.RoleBinding.Columns().UserGroupId: userGroupID,
		})
		if err != nil {
			return err
		}
		// 插入新记录
		entities := make([]g.Map, 0)
		uniqueEntities := make(map[string]struct{}) // 用于去重
		for _, item := range data {
			uniqueKey := fmt.Sprintf("%d-%d-%d-%s", item.UserGroupID, item.ClusterID, item.RoleID, item.Namespace)
			if _, exists := uniqueEntities[uniqueKey]; !exists {
				uniqueEntities[uniqueKey] = struct{}{} // 标记为已出现
				entities = append(entities, g.Map{
					dao.RoleBinding.Columns().UserGroupId: item.UserGroupID,
					dao.RoleBinding.Columns().ClusterId:   item.ClusterID,
					dao.RoleBinding.Columns().RoleId:      item.RoleID,
					dao.RoleBinding.Columns().Namespace:   item.Namespace,
				})
			}
		}
		_, err = tx.Save(dao.RoleBinding.Table(), entities)
		if err != nil {
			return err
		}

		return nil
	})

	return err
}

// ManageRoleBindingByRoleID   管理角色绑定关系
func (x *RoleBindingRepo) ManageRoleBindingByRoleID(ctx context.Context, roleID uint32, data []*biz.RoleBindingItem) error {
	db := dao.RoleBinding.Ctx(ctx)
	// 开启事务
	err := db.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		// 删除现有记录
		_, err := tx.Delete(dao.RoleBinding.Table(), g.Map{
			dao.RoleBinding.Columns().RoleId: roleID,
		})
		if err != nil {
			return err
		}

		// 插入新记录
		entities := make([]g.Map, 0)
		uniqueEntities := make(map[string]struct{}) // 用于去重
		for _, item := range data {
			if item.ClusterID == 0 && item.UserGroupID == 0 {
				continue
			}
			uniqueKey := fmt.Sprintf("%d-%d-%d-%s", item.UserGroupID, item.ClusterID, item.RoleID, item.Namespace)
			if _, exists := uniqueEntities[uniqueKey]; !exists {
				uniqueEntities[uniqueKey] = struct{}{} // 标记为已出现
				entities = append(entities, g.Map{
					dao.RoleBinding.Columns().UserGroupId: item.UserGroupID,
					dao.RoleBinding.Columns().ClusterId:   item.ClusterID,
					dao.RoleBinding.Columns().RoleId:      item.RoleID,
					dao.RoleBinding.Columns().Namespace:   item.Namespace,
				})
			}
		}
		if len(entities) == 0 {
			return nil
		}
		_, err = tx.Save(dao.RoleBinding.Table(), entities)
		if err != nil {
			return err
		}
		return nil
	})

	return err
}

// InsertOrUpdate 插入或更新
func (x *RoleBindingRepo) InsertOrUpdate(ctx context.Context, data []*biz.RoleBindingItem) error {
	db := dao.RoleBinding.Ctx(ctx)

	// 将业务数据转换为数据库实体
	var items []*entity.RoleBinding
	for _, item := range data {
		items = append(items, x.convertDO2VO(item))
	}
	_, err := db.Save(dao.RoleBinding.Table(), items)

	return err
}
