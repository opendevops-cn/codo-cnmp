package data

import (
	"codo-cnmp/internal/biz"
	"codo-cnmp/internal/model/dao"
	"codo-cnmp/internal/model/model/entity"
	"codo-cnmp/pb"
	"context"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/gogf/gf/v2/database/gdb"
)

type UserFollowRepo struct {
	data *Data
	log  *log.Helper
}

func (x *UserFollowRepo) convertQuery(db *gdb.Model, query *biz.ListUserFollowRequest) *gdb.Model {
	if query.Page <= 0 {
		query.Page = 1
	}
	if query.PageSize <= 0 {
		query.PageSize = 10
	}
	if query.UserID > 0 {
		db = db.Where(dao.UserFollow.Columns().UserId, query.UserID)
	}
	if query.FollowType >= 0 {
		db = db.Where(dao.UserFollow.Columns().FollowType, query.FollowType)
	}
	if query.FollowValue != "" {
		db = db.Where(dao.UserFollow.Columns().FollowValue, query.FollowValue)
	}
	if query.Keyword != "" {
		db = db.WhereLike(dao.UserFollow.Columns().FollowValue, "%"+query.Keyword+"%")
	}

	if query.ListAll {
		return db
	}
	return db.Page(int(query.Page), int(query.PageSize))
}

func (x *UserFollowRepo) List(ctx context.Context, data *biz.ListUserFollowRequest) ([]*biz.UserFollowItem, uint32, error) {
	result := make([]*biz.UserFollowItem, 0)
	db := dao.UserFollow.Ctx(ctx)
	db = x.convertQuery(db, data)
	userFollows := make([]*entity.UserFollow, 0)
	err := db.Scan(&userFollows)
	if err != nil {
		return result, 0, err
	}
	count, err := db.Count()
	if err != nil {
		return result, 0, err
	}
	for _, item := range userFollows {
		userFollow := x.convertVO2DO(item)
		result = append(result, userFollow)
	}
	return result, uint32(count), nil
}

func (x *UserFollowRepo) convertDO2VO(data *biz.UserFollowRequest) *entity.UserFollow {
	vo := &entity.UserFollow{
		UserId:      uint64(data.UserID),
		FollowType:  int(data.FollowType),
		FollowValue: data.FollowValue,
		ClusterName: data.ClusterName,
	}
	return vo
}

func (x *UserFollowRepo) convertVO2DO(data *entity.UserFollow) *biz.UserFollowItem {
	vo := &biz.UserFollowItem{
		UserID:      uint32(data.UserId),
		FollowType:  pb.FollowType(data.FollowType),
		FollowValue: data.FollowValue,
		CreatedTime: data.CreatedAt.String(),
		ClusterName: data.ClusterName,
	}
	return vo
}

func (x *UserFollowRepo) Create(ctx context.Context, data *biz.UserFollowRequest) error {
	db := dao.UserFollow.Ctx(ctx)
	var (
		exist int
		err   error
	)
	query := db.Where(dao.UserFollow.Columns().UserId, data.UserID).
		Where(dao.UserFollow.Columns().FollowValue, data.FollowValue).
		Where(dao.UserFollow.Columns().FollowType, data.FollowType)
	if data.FollowType != pb.FollowType_Cluster {
		query = query.Where(dao.UserFollow.Columns().ClusterName, data.ClusterName)
	}
	exist, err = query.Count()
	if err != nil {
		return err
	}
	if exist > 0 {
		return fmt.Errorf("已关注,请勿重复关注")
	}
	vo := x.convertDO2VO(data)
	_, err = db.InsertAndGetId(vo)
	return err
}

func (x *UserFollowRepo) Delete(ctx context.Context, data *biz.DeleteUserFollowRequest) error {
	return dao.UserFollow.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		db := dao.UserFollow.Ctx(ctx).TX(tx)
		query := db.Where(dao.UserFollow.Columns().UserId, data.UserID).
			Where(dao.UserFollow.Columns().FollowValue, data.FollowValue).
			Where(dao.UserFollow.Columns().FollowType, data.FollowType)

		if data.FollowType != pb.FollowType_Cluster {
			query = query.Where(dao.UserFollow.Columns().ClusterName, data.ClusterName)
		}

		count, err := query.Count()
		if err != nil {
			return fmt.Errorf("查询关注状态失败: %w", err)
		}
		if count < 1 {
			return fmt.Errorf("未关注,无法取消关注")
		}

		_, err = query.Delete()
		if err != nil {
			return fmt.Errorf("取消关注失败: %w", err)
		}
		return nil
	})
}

func NewUserFollowRepo(data *Data, logger log.Logger) *UserFollowRepo {
	return &UserFollowRepo{
		data: data,
		log:  log.NewHelper(log.With(logger, "module", "data/user_follow")),
	}
}

func NewIUserFollowRepo(x *UserFollowRepo) biz.IUserFollowRepo {
	return x
}
