package middleware

import (
	"codo-cnmp/common/xerr"
	"context"
	kmiddleware "github.com/go-kratos/kratos/v2/middleware"
)

type ValidateMiddleware struct {
}

func NewValidateMiddleware() *ValidateMiddleware {
	return &ValidateMiddleware{}
}

func (x *ValidateMiddleware) Server() kmiddleware.Middleware {
	return func(handler kmiddleware.Handler) kmiddleware.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			// 如果未启用，则直接跳过
			if validate, ok := req.(interface{ Validate() error }); ok {
				err := validate.Validate()
				if err != nil {
					return nil, xerr.NewErrCodeMsg(xerr.RequestParamError, err.Error())
				}
			}

			// 继续执行后续的中间件和服务
			return handler(ctx, req)
		}
	}
}
