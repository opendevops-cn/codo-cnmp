package data

import (
	"codo-cnmp/common/utils"
	"codo-cnmp/internal/biz"
	"codo-cnmp/internal/model/dao"
	"codo-cnmp/internal/model/model/entity"
	"codo-cnmp/pb"
	"context"
	"fmt"
	"github.com/ccheers/xpkg/generic/arrayx"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

type RoleRepo struct {
	data *Data
	log  *log.Helper
}

func (x *RoleRepo) GetRoleByName(ctx context.Context, roleName string) (*biz.RoleItem, error) {
	db := dao.Role.Ctx(ctx)
	var role entity.Role
	err := db.Where(dao.Role.Columns().Name, roleName).Scan(&role)
	if err != nil {
		return nil, err
	}
	return x.convertVO2DO(&role), nil
}

func (x *RoleRepo) ExistRoleByName(ctx context.Context, roleName string) (bool, error) {
	db := dao.Role.Ctx(ctx)
	count, err := db.Where(dao.Role.Columns().Name, roleName).Count()
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (x *RoleRepo) GetRoleByID(ctx context.Context, roleID uint32) (*biz.RoleItem, error) {
	db := dao.Role.Ctx(ctx)
	var role entity.Role
	err := db.Where(dao.Role.Columns().Id, roleID).Scan(&role)
	if err != nil {
		return nil, err
	}
	return x.convertVO2DO(&role), nil
}

func NewRoleRepo(data *Data, logger log.Logger) *RoleRepo {
	return &RoleRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func NewIGrantedRoleRepo(repo *RoleRepo) biz.IRoleRepo {
	return repo
}

func (x *RoleRepo) convertQuery(db *gdb.Model, query *biz.ListRoleRequest) *gdb.Model {
	if query.Keyword != "" {
		db = db.WhereLike(dao.Role.Columns().Name, "%"+query.Keyword+"%")
	}
	if query.ListAll {
		return db
	}
	if query.Page <= 0 {
		query.Page = 1
	}
	if query.PageSize <= 0 {
		query.PageSize = 10
	}
	return db.Page(int(query.Page), int(query.PageSize))

}

func (x *RoleRepo) convertDO2VO(data *biz.RoleItem) *entity.Role {
	vo := &entity.Role{
		Id:          uint64(data.ID),
		Name:        data.Name,
		Description: data.Description,
		RoleType:    int(data.RoleType),
		IsDefault:   utils.BoolToInt(data.ISDefault),
		Yaml:        data.YamlStr,
		UpdateBy:    data.UpdateBy,
	}
	return vo
}

// Create 创建角色
func (x *RoleRepo) Create(ctx context.Context, role *biz.RoleItem) error {
	db := dao.Role.Ctx(ctx)
	exist, err := db.Where(dao.Role.Columns().Name, role.Name).Count()
	if err != nil {
		return err
	}
	if exist > 0 {
		return fmt.Errorf("角色名称已存在")
	}
	vo := x.convertDO2VO(role)
	_, err = db.InsertAndGetId(vo)
	return err
}

// List 获取用户组授权列表
func (x *RoleRepo) List(ctx context.Context, query *biz.ListRoleRequest) ([]*biz.RoleItem, error) {
	db := dao.Role.Ctx(ctx)
	db = x.convertQuery(db, query)
	var roles []*entity.Role
	err := db.Scan(&roles)
	if err != nil {
		return nil, err
	}
	return arrayx.Map(roles, func(t *entity.Role) *biz.RoleItem {
		return x.convertVO2DO(t)
	}), nil
}

// Delete 删除用户组授权
func (x *RoleRepo) Delete(ctx context.Context, roleID uint32) error {
	db := dao.Role.Ctx(ctx)
	_, err := db.Where(dao.Role.Columns().Id, roleID).Delete()
	return err
}

// Update 更新用户组授权
func (x *RoleRepo) Update(ctx context.Context, data *biz.RoleItem) error {
	db := dao.Role.Ctx(ctx)
	vo := x.convertDO2VO(data)
	updates := make(g.Map)
	if vo.Name != "" {
		updates[dao.Role.Columns().Name] = vo.Name
	}
	if vo.Description != "" {
		updates[dao.Role.Columns().Description] = vo.Description
	}
	if vo.RoleType >= 0 {
		updates[dao.Role.Columns().RoleType] = vo.RoleType
	}
	if vo.Yaml != "" {
		updates[dao.Role.Columns().Yaml] = vo.Yaml
	}
	if vo.IsDefault != 0 {
		updates[dao.Role.Columns().IsDefault] = vo.IsDefault
	}
	_, err := db.Where(dao.Role.Columns().Id, data.ID).Update(updates)
	if err != nil {
		return err
	}
	return nil
}

func (x *RoleRepo) convertVO2DO(t *entity.Role) *biz.RoleItem {
	return &biz.RoleItem{
		ID:          uint32(t.Id),
		Name:        t.Name,
		Description: t.Description,
		RoleType:    pb.RoleType(t.RoleType),
		ISDefault:   utils.IntToBool(t.IsDefault),
		YamlStr:     t.Yaml,
		CreateTime:  t.CreatedAt.String(),
		UpdateTime:  t.UpdatedAt.String(),
		UpdateBy:    t.UpdateBy,
	}

}

// Count 获取用户组授权列表数量
func (x *RoleRepo) Count(ctx context.Context, query *biz.ListRoleRequest) (uint32, error) {
	db := dao.Role.Ctx(ctx)
	db = x.convertQuery(db, query)
	count, err := db.Count()
	if err != nil {
		return 0, err
	}
	return uint32(count), nil
}
