package data

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"codo-cnmp/common/consts"
	"codo-cnmp/internal/biz"
	"codo-cnmp/internal/dep"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-redis/redis/v8"
)

type MeshRepo struct {
	data  *Data
	log   *log.Helper
	apiGw *dep.CODOAPIGateway
}

func (x *MeshRepo) DeleteMesh(ctx context.Context, id string) error {
	panic("implement me")
}

// MeshResponse 结构体映射 JSON 响应
type MeshResponse struct {
	Status int      `json:"status"`
	Code   int      `json:"code"`
	Data   MeshData `json:"data"`
	Msg    string   `json:"msg"`
}

// MeshData 结构体映射 data 字段
type MeshData struct {
	Addr string `json:"addr"`
}

func (x *MeshRepo) CreateMesh(ctx context.Context, req *biz.MeshItem) (string, error) {
	bytesBody, err := json.Marshal(req)
	if err != nil {
		x.log.WithContext(ctx).Errorf("组网注册请求序列化失败: %v", err)
		return "", err
	}
	response, err := x.apiGw.SendRequest(ctx, http.MethodPost, "/api/agent/v1/manager/agent/mesh/register", bytesBody, nil)
	if err != nil {
		x.log.WithContext(ctx).Errorf("组网注册请求失败， 请求Body: %v, %v", string(bytesBody), err)
		return "", err
	}
	if response.StatusCode != http.StatusOK {
		x.log.WithContext(ctx).Errorf("组网注册请求失败, 状态码异常: %v", response.Status)
		return "", err
	}
	data, err := io.ReadAll(response.Body)
	defer response.Body.Close()
	if err != nil {
		x.log.WithContext(ctx).Errorf("读取响应体失败: %v", err)
		return "", err
	}
	var meshResp MeshResponse
	err = json.Unmarshal(data, &meshResp)
	if err != nil {
		x.log.WithContext(ctx).Errorf("响应体解析失败: %v", err)
		return "", err
	}
	if meshResp.Status != 0 {
		x.log.WithContext(ctx).Errorf("组网注册失败， 请求Body: %v, %s", string(bytesBody), meshResp.Msg)
		return "", fmt.Errorf("组网注册失败: %v", meshResp.Msg)
	}
	return meshResp.Data.Addr, nil
}

func (x *MeshRepo) ListMesh(ctx context.Context, req *biz.ListMeshRequest) ([]*biz.MeshItem, uint32, error) {
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}
	if req.ListAll {
		req.Page = 1
		req.PageSize = 300
	}
	res := make([]*biz.MeshItem, 0)
	cachedData, err := x.data.redis.Get(ctx, consts.MeshCacheKey).Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		x.log.WithContext(ctx).Errorf("redis get error: %v", err)
		return res, 0, nil
	}
	if cachedData != "" {
		err = json.Unmarshal([]byte(cachedData), &res)
		if err != nil {
			x.log.WithContext(ctx).Errorf("json unmarshal error: %v", err)
			return res, 0, nil
		}
		return res, uint32(len(res)), nil
	}
	return res, 0, nil
}

func NewMeshRepo(data *Data, logger log.Logger, apiGw *dep.CODOAPIGateway) *MeshRepo {
	return &MeshRepo{
		data:  data,
		apiGw: apiGw,
		log:   log.NewHelper(logger),
	}
}

func NewIMeshRepo(x *MeshRepo) biz.IMeshRepository {
	return x
}
