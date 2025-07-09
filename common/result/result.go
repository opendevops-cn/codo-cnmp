package result

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"google.golang.org/protobuf/encoding/protojson"
	nethttp "net/http"
	"strings"

	"github.com/gorilla/websocket"
	pkgerr "github.com/pkg/errors"

	"codo-cnmp/common/xerr"
)

// KHttpResult http返回
func KHttpResult(writer nethttp.ResponseWriter, _ *nethttp.Request, i interface{}) error {

	r := Success(i)
	writer.WriteHeader(nethttp.StatusOK)
	return json.NewEncoder(writer).Encode(r)
}

// KHttpParseErr http参数解析错误返回
func KHttpParseErr(err error) xerr.ErrCode {
	// 错误返回
	newErr := errors.Unwrap(err)
	if newErr != nil {
		err = newErr
	}
	errcode := xerr.ServerCommonError
	if impl, ok := err.(interface{ ErrorName() string }); ok &&
		strings.HasSuffix(impl.ErrorName(), "ValidationError") {
		errMsg := fmt.Sprintf("%s：%s", xerr.MapErrMsg(xerr.RequestParamError), err.Error())
		err = xerr.NewErrCodeMsg(xerr.RequestParamError, errMsg)
	}
	causeErr := pkgerr.Cause(err) // err类型
	var e *xerr.CodeError
	if errors.As(causeErr, &e) { // 自定义错误类型
		// 自定义CodeError
		errcode = e.GetErrCode()
	}
	return errcode
}

// KHttpError http错误返回
func KHttpError(w nethttp.ResponseWriter, r *nethttp.Request, err error) {
	errCode := KHttpParseErr(err)
	errMsg := xerr.MapErrMsg(errCode)
	statusCode := nethttp.StatusBadRequest
	if canRet200(errCode) {
		statusCode = nethttp.StatusOK
	}
	if errCode == xerr.ErrNotAllowed {
		statusCode = nethttp.StatusForbidden
	}
	if errCode == xerr.ErrUnAuthorization {
		statusCode = nethttp.StatusUnauthorized
	}
	if errCode == xerr.ServerCommonError {
		statusCode = nethttp.StatusInternalServerError
	}
	if errCode == xerr.DryRunResourceError {
		causeErr := pkgerr.Cause(err) // err类型
		var e *xerr.CodeError
		if errors.As(causeErr, &e) {
			errMsg = e.GetErrMsg()
		}
	}
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(Error(errCode, errMsg, err.Error()))
}

func WebSocketResult(ctx context.Context, conn *websocket.Conn, resp interface{}, err error) {
	if err == nil {
		// 成功返回
		r := Success(resp)
		conn.WriteJSON(r)
		return
	}
	// 错误返回
	errcode := xerr.ServerCommonError
	if impl, ok := err.(interface{ ErrorName() string }); ok &&
		strings.HasSuffix(impl.ErrorName(), "ValidationError") {
		errMsg := fmt.Sprintf("%s：%s", xerr.MapErrMsg(xerr.RequestParamError), err.Error())
		err = pkgerr.Wrapf(xerr.NewErrCodeMsg(xerr.RequestParamError, errMsg), err.Error())
	}
	causeErr := pkgerr.Cause(err) // err类型
	var e *xerr.CodeError
	if errors.As(causeErr, &e) { // 自定义错误类型
		// 自定义CodeError
		errcode = e.GetErrCode()
	}

	//logx.WithContext(ctx).Errorf("【WS-ERR】 : %+v ", err)

	_ = conn.WriteJSON(Error(errcode, xerr.MapErrMsg(errcode), err.Error()))
}

var whiteList = map[xerr.ErrCode]struct{}{}

func canRet200(statusCode xerr.ErrCode) bool {
	_, ok := whiteList[statusCode]
	return ok
}

var protoMarshalOptions = protojson.MarshalOptions{
	UseProtoNames:     true,
	UseEnumNumbers:    false,
	EmitDefaultValues: true,
	//EmitUnpopulated: true,
}
