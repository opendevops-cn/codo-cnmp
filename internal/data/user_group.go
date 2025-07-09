package data

import (
	"codo-cnmp/common/consts"
	"codo-cnmp/internal/biz"
	"codo-cnmp/internal/dep"
	"codo-cnmp/internal/model/dao"
	"codo-cnmp/internal/model/model/entity"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-redis/redis/v8"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
	"io"
	"net/http"
)

type UserGroupRepoV2 struct {
	data  *Data
	apiGw *dep.CODOAPIGateway
	log   *log.Helper
}

func (x *UserGroupRepoV2) ListUsers(ctx context.Context, req *biz.ListUserRequest) ([]*biz.User, uint32, error) {
	users := make([]*biz.User, 0)
	page := req.Page
	if page == 0 {
		page = 1
	}
	pageSize := req.PageSize
	if pageSize == 0 {
		pageSize = 10
	}
	if req.ListAll {
		page = 1
		pageSize = 300
	}

	url := fmt.Sprintf("/api/p/v4/user/?page_number=%d&page_size=%d&order_by=id&order=ascend&searchVal=%s",
		page, pageSize, req.Keyword)
	response, err := x.apiGw.SendRequest(ctx, "GET", url, nil, nil)
	if err != nil {
		x.log.WithContext(ctx).Errorf("获取用户列表失败: %v", err)
		return users, 0, err
	}
	if response.StatusCode != http.StatusOK {
		x.log.WithContext(ctx).Errorf("获取用户列表失败, 状态码异常: %v", response.Status)
		return users, 0, errors.New(response.Status)
	}
	data, err := io.ReadAll(response.Body)
	defer response.Body.Close()
	if err != nil {
		x.log.WithContext(ctx).Errorf("读取响应体失败: %v", err)
		return users, 0, err
	}
	type UserResponse struct {
		Code  int        `json:"code"`
		Count int        `json:"count"`
		Data  []biz.User `json:"data"`
		Msg   string     `json:"msg"`
	}
	var userResp UserResponse
	if err := json.Unmarshal(data, &userResp); err != nil {
		x.log.WithContext(ctx).Errorf("解析响应体失败: %v", err)
	}
	if userResp.Code != 0 {
		return users, 0, errors.New(userResp.Msg)
	}

	for _, user := range userResp.Data {
		users = append(users, &user)
	}
	return users, uint32(userResp.Count), nil

}

func (x *UserGroupRepoV2) GetUsersByUserGroupID(ctx context.Context, groupID uint32) ([]*biz.GroupUser, uint32, error) {
	users := make([]*biz.GroupUser, 0)
	response, err := x.apiGw.SendRequest(ctx, "GET", fmt.Sprintf("/api/p/v4/role_user/?role_id=%d", groupID), nil, nil)
	if err != nil {
		x.log.WithContext(ctx).Errorf("获取用户列表失败: %v", err)
		return users, 0, err
	}
	if response.StatusCode != http.StatusOK {
		x.log.WithContext(ctx).Errorf("获取用户列表失败, 状态码异常: %v", response.Status)
		return users, 0, err
	}
	data, err := io.ReadAll(response.Body)
	defer response.Body.Close()
	if err != nil {
		x.log.WithContext(ctx).Errorf("读取响应体失败: %v", err)
		return users, 0, err
	}
	type UserResponse struct {
		Code  int             `json:"code"`
		Count int             `json:"count"`
		Data  []biz.GroupUser `json:"data"`
		Msg   string          `json:"msg"`
	}
	var userResp UserResponse
	err = json.Unmarshal(data, &userResp)
	if err != nil {
		x.log.WithContext(ctx).Errorf("解析响应体失败: %v", err)
	}
	if userResp.Code != 0 {
		x.log.WithContext(ctx).Errorf("获取用户列表失败: %v", userResp.Msg)
		return users, 0, fmt.Errorf(userResp.Msg)
	}
	for _, user := range userResp.Data {
		users = append(users, &user)
	}
	return users, uint32(userResp.Count), nil

}

func (x *UserGroupRepoV2) convertDO2VO(data *biz.UserGroupV2) *entity.UserGroup {
	return &entity.UserGroup{
		Name:        data.Name,
		UserGroupId: data.UserGroupID,
	}
}

func (x *UserGroupRepoV2) convertVO2DO(entity *entity.UserGroup) *biz.UserGroupV2 {
	return &biz.UserGroupV2{
		Name:        entity.Name,
		UserGroupID: entity.UserGroupId,
	}
}

func (x *UserGroupRepoV2) ListUserGroups(ctx context.Context, req *biz.ListUserGroupRequest) ([]*biz.UserGroup, uint32, error) {
	res := make([]*biz.UserGroup, 0)
	data, err := x.data.redis.Get(ctx, consts.UserGroupCacheKey).Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		x.log.WithContext(ctx).Errorf("redis get error: %v", err)
		return res, 0, nil
	}

	if data != "" {
		// 从缓存中获取数据
		var userGroups []*biz.UserGroup
		err := json.Unmarshal([]byte(data), &userGroups)
		if err != nil {
			return res, 0, err
		}
		return userGroups, uint32(len(userGroups)), nil
	}
	return res, 0, nil
}

func (x *UserGroupRepoV2) CreateUserGroup(ctx context.Context, data *biz.UserGroupV2) (bool, error) {
	model := dao.UserGroup.Ctx(ctx)
	exists, err := x.ExistUserGroup(ctx, data.UserGroupID)
	if err != nil {
		return false, err
	}
	if !exists {
		return false, fmt.Errorf("用户组不存在")
	}
	vo := x.convertDO2VO(data)
	_, err = model.InsertAndGetId(vo)
	if err != nil {
		return false, err
	}
	return true, nil

}

func (x *UserGroupRepoV2) ExistUserGroup(ctx context.Context, userGroupID uint64) (bool, error) {
	model := dao.UserGroup.Ctx(ctx)
	count, err := model.Where(dao.UserGroup.Columns().UserGroupId, userGroupID).Count()
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (x *UserGroupRepoV2) UpdateUserGroup(ctx context.Context, data *biz.UserGroupV2) (bool, error) {
	model := dao.UserGroup.Ctx(ctx)
	exists, err := x.ExistUserGroup(ctx, data.UserGroupID)
	if err != nil {
		return false, err
	}
	if !exists {
		return false, fmt.Errorf("用户组不存在")
	}
	vo := x.convertDO2VO(data)
	_, err = model.Where(dao.UserGroup.Columns().UserGroupId, data.UserGroupID).Update(vo)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (x *UserGroupRepoV2) DeleteUserGroup(ctx context.Context, userGroupID uint64) (bool, error) {
	model := dao.UserGroup.Ctx(ctx)
	exists, err := x.ExistUserGroup(ctx, userGroupID)
	if err != nil {
		return false, err
	}
	if !exists {
		return false, fmt.Errorf("用户组不存在")
	}
	_, err = model.Where(dao.UserGroup.Columns().UserGroupId, userGroupID).Delete()
	if err != nil {
		return false, err
	}
	return true, nil
}

func (x *UserGroupRepoV2) GetUserGroupByUserGroupID(ctx context.Context, userGroupID uint64) (*biz.UserGroupV2, error) {
	model := dao.UserGroup.Ctx(ctx)
	exist, err := model.Where(dao.UserGroup.Columns().UserGroupId, userGroupID).Count()
	if err != nil {
		return nil, err
	}
	if exist == 0 {
		return nil, fmt.Errorf("用户组不存在")
	}
	var user entity.UserGroup
	err = model.Where(dao.UserGroup.Columns().UserGroupId, userGroupID).Scan(&user)
	if err != nil {
		return nil, err
	}
	return x.convertVO2DO(&user), nil
}

func (x *UserGroupRepoV2) BulkInsertOrUpdateUserGroup(ctx context.Context, data []*biz.UserGroupV2) (bool, error) {
	model := dao.UserGroup.Ctx(ctx)
	entities := make([]g.Map, 0)
	// 开启事务
	err := model.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		for _, item := range data {
			entities = append(entities, g.Map{
				dao.UserGroup.Columns().UserGroupId: item.UserGroupID,
				dao.UserGroup.Columns().Name:        item.Name,
			})
		}
		_, err := tx.Save(dao.UserGroup.Table(), entities)
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

func NewUserGroupRepoRepoV2(data *Data, logger log.Logger, apiGW *dep.CODOAPIGateway) *UserGroupRepoV2 {
	return &UserGroupRepoV2{
		data:  data,
		apiGw: apiGW,
		log:   log.NewHelper(logger),
	}
}

func NewIUserGroupRepoV2(repo *UserGroupRepoV2) biz.IUserGroupV2Repo {
	return repo
}
