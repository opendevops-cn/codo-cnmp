package data

import (
	"codo-cnmp/internal/biz"
	"codo-cnmp/internal/model/dao"
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

type UserUserGroupRelRepo struct {
	data *Data
	log  *log.Helper
}

func (x *UserRepo) BulkInsertOrUpdateUserUserGroupRel(ctx context.Context, data []*biz.UserUserGroupRel) (bool, error) {
	model := dao.UserUserGroupRel.Ctx(ctx)
	entities := make([]g.Map, 0)
	// 开启事务
	err := model.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		for _, item := range data {
			entities = append(entities, g.Map{
				dao.UserUserGroupRel.Columns().UserId:      item.UserID,
				dao.UserUserGroupRel.Columns().UserGroupId: item.UserGroupID,
			})
		}
		_, err := tx.Save(dao.UserUserGroupRel.Table(), entities)
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
