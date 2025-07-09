package data

import (
	"codo-cnmp/internal/biz"
	"codo-cnmp/pb"
	"context"
	"github.com/go-kratos/kratos/v2/log"
)

type nameSpaceRepo struct {
	data *Data
	log  *log.Helper
}

func (n *nameSpaceRepo) ListNameSpace(ctx context.Context, req *biz.ListNameSpaceRequest) (resp *pb.ListNameSpaceResponse, err error) {
	//TODO implement me
	panic("implement me")
}

func (n *nameSpaceRepo) CreateNameSpace(ctx context.Context, req *biz.NameSpaceItem) error {
	//TODO implement me
	panic("implement me")
}

func (n *nameSpaceRepo) DeleteNameSpace(ctx context.Context, req *biz.DeleteNameSpaceRequest) error {
	// TODO implement me
	panic("implement me")
}

func (n *nameSpaceRepo) UpdateNameSpace(ctx context.Context, req *biz.UpdateNameSpaceRequest) error {
	// TODO implement me
	panic("implement me")
}

func (n *nameSpaceRepo) GetNameSpaceYaml(ctx context.Context, req *biz.GetNamespaceYamlRequest) (string, error) {
	// TODO implement me
	panic("implement me")
}

func (n *nameSpaceRepo) CreateNameSpaceByYaml(ctx context.Context, req *biz.CreateNameSpaceRequest) error {
	// TODO implement me
	panic("implement me")
}

func (n *nameSpaceRepo) UpdateNameSpaceByYaml(ctx context.Context, req *biz.CreateNameSpaceByYamlRequest) error {
	// TODO implement me
	panic("implement me")
}

//func NewNameSpaceRepo(data *Data, logger log.Logger) biz.NameSpaceRepo {
//	return &nameSpaceRepo{
//		data: data,
//		log:  log.NewHelper(logger),
//	}
//}
