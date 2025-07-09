package data

import (
	"context"

	"codo-cnmp/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
)

type greeterRepo struct {
	data *Data
	log  *log.Helper
}

// NewGreeterRepo .
func NewGreeterRepo(data *Data, logger log.Logger) biz.GreeterRepo {
	return &greeterRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (r *greeterRepo) Save(ctx context.Context, g *biz.Greeter) (*biz.Greeter, error) {
	r.log.WithContext(ctx).Infof("Greeter: %v", g)

	result, err := r.data.db.Query(ctx, "select * from user where username like ? limit 1", g.Hello+"%")
	if err != nil {
		return nil, err
	}

	r.log.WithContext(ctx).Infof("result size: %v", result.Size())
	return g, nil
}

func (r *greeterRepo) Update(ctx context.Context, g *biz.Greeter) (*biz.Greeter, error) {
	r.log.WithContext(ctx).Infof("Greeter: %v", g)
	return g, nil
}

func (r *greeterRepo) FindByID(ctx context.Context, id int64) (*biz.Greeter, error) {
	r.log.WithContext(ctx).Infof("FindByID: %v", id)
	return nil, nil
}

func (r *greeterRepo) ListByHello(ctx context.Context, hello string) ([]*biz.Greeter, error) {
	r.log.WithContext(ctx).Infof("ListByHello: %v", hello)
	return nil, nil
}

func (r *greeterRepo) ListAll(ctx context.Context) ([]*biz.Greeter, error) {
	r.log.WithContext(ctx).Infof("ListAll")
	return nil, nil
}
