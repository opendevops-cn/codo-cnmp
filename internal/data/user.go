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

type UserRepo struct {
	data *Data
	log  *log.Helper
}

func (x *UserRepo) CreateUser(ctx context.Context, data *biz.RoleUser) (bool, error) {
	model := dao.User.Ctx(ctx)
	exist, err := x.ExistUser(ctx, data.UserID)
	if err != nil {
		return false, err
	}
	if !exist {
		return false, fmt.Errorf("用户不存在")
	}
	vo := x.convertDO2VO(data)
	_, err = model.InsertAndGetId(vo)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (x *UserRepo) ExistUser(ctx context.Context, userID uint32) (bool, error) {
	model := dao.User.Ctx(ctx)
	count, err := model.Where(dao.User.Columns().UserId, userID).Count()
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (x *UserRepo) DeleteUser(ctx context.Context, userID uint32) (bool, error) {
	model := dao.User.Ctx(ctx)
	exist, err := x.ExistUser(ctx, userID)
	if err != nil {
		return false, err
	}
	if !exist {
		return false, fmt.Errorf("用户不存在")
	}
	_, err = model.Where(dao.User.Columns().UserId, userID).Delete()
	if err != nil {
		return false, err
	}
	return true, nil
}

func (x *UserRepo) UpdateUser(ctx context.Context, data *biz.RoleUser) (bool, error) {
	model := dao.User.Ctx(ctx)
	exist, err := x.ExistUser(ctx, data.UserID)
	if err != nil {
		return false, err
	}
	if !exist {
		return false, fmt.Errorf("用户不存在")
	}
	vo := x.convertDO2VO(data)
	_, err = model.Where(dao.User.Columns().UserId, data.UserID).Update(vo)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (x *UserRepo) GetUserByUserID(ctx context.Context, userID uint32) (*biz.RoleUser, error) {
	model := dao.User.Ctx(ctx)
	exist, err := x.ExistUser(ctx, userID)
	if err != nil {
		return nil, err
	}
	if !exist {
		return nil, fmt.Errorf("用户不存在")
	}
	var user entity.User
	err = model.Where(dao.User.Columns().UserId, userID).Scan(&user)
	if err != nil {
		return nil, err
	}
	return x.convertVO2DO(&user), nil
}

func (x *UserRepo) convertDO2VO(data *biz.RoleUser) *entity.User {
	return &entity.User{
		Username: data.Username,
		Nickname: data.Nickname,
		UserId:   uint64(data.UserID),
	}
}

func (x *UserRepo) convertVO2DO(e *entity.User) *biz.RoleUser {
	return &biz.RoleUser{
		Username: e.Username,
		Nickname: e.Nickname,
		UserID:   uint32(e.UserId),
	}
}

func NewUserRepo(data *Data, logger log.Logger) *UserRepo {
	return &UserRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func NewIUserRepo(repo *UserRepo) biz.IUserRepo {
	return repo
}

func (x *UserRepo) BulkInsertOrUpdateUser(ctx context.Context, data []*biz.RoleUser) (bool, error) {
	model := dao.User.Ctx(ctx)
	entities := make([]g.Map, 0)
	// 开启事务
	err := model.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		for _, item := range data {
			entities = append(entities, g.Map{
				dao.User.Columns().UserId:   item.UserID,
				dao.User.Columns().Username: item.Username,
				dao.User.Columns().Nickname: item.Nickname,
				dao.User.Columns().Email:    item.Email,
				dao.User.Columns().Source:   item.Source,
			})
		}
		_, err := tx.Save(dao.User.Table(), entities)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return false, err
	}
	return true, nil
}
