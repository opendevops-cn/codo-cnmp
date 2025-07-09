package xerr

import (
	"fmt"

	codes "google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

/**
常用通用固定错误
*/

type CodeError struct {
	errCode ErrCode
	errMsg  string
}

// GetErrCode 返回给前端的错误码
func (e *CodeError) GetErrCode() ErrCode {
	return e.errCode
}

// GetErrMsg 返回给前端显示端错误信息
func (e *CodeError) GetErrMsg() string {
	return e.errMsg
}

func (e *CodeError) Error() string {
	return fmt.Sprintf("ErrCode:%d，ErrMsg:%s", e.errCode, e.errMsg)
}

func NewErrCodeMsg(errCode ErrCode, errMsg string) *CodeError {
	return &CodeError{errCode: errCode, errMsg: errMsg}
}

func NewErrCode(errCode ErrCode) *CodeError {
	return &CodeError{errCode: errCode, errMsg: MapErrMsg(errCode)}
}

func NewErrMsg(errMsg string) *CodeError {
	return &CodeError{errCode: ServerCommonError, errMsg: errMsg}
}

func NewForbiddenErrMsg(errMsg string) *CodeError {
	return &CodeError{errCode: ErrNotAllowed, errMsg: errMsg}
}

func NewRpcErrCode(errCode ErrCode) error {
	return status.New(codes.Code(errCode), MapErrMsg(errCode)).Err()
}

func NewRpcErrMsg(errMsg string) error {
	return status.New(codes.Unknown, errMsg).Err()
}

func NewRpcErrCodeMsg(errCode ErrCode, errMsg string) error {
	return status.New(codes.Code(errCode), errMsg).Err()
}
