package result

import (
	"codo-cnmp/common/xerr"
	"encoding/json"
	"google.golang.org/protobuf/proto"
)

type ResponseSuccessBean struct {
	Code uint32          `json:"code"`
	Msg  string          `json:"message"`
	Data json.RawMessage `json:"data"`
}
type NullJson struct{}

func Success(data interface{}) *ResponseSuccessBean {
	bs, _ := marshalJSON(data)
	return &ResponseSuccessBean{
		Code: 0,
		Msg:  "OK",
		Data: bs,
	}
}

type ResponseErrorBean struct {
	Code   uint32 `json:"code"`
	Msg    string `json:"message"`
	Reason string `json:"reason"`
}

func Error(errCode xerr.ErrCode, errMsg, reason string) *ResponseErrorBean {
	return &ResponseErrorBean{
		Code:   uint32(errCode),
		Msg:    errMsg,
		Reason: reason,
	}
}

func marshalJSON(v interface{}) ([]byte, error) {
	// 成功返回
	switch m := v.(type) {
	case json.Marshaler:
		return m.MarshalJSON()
	// 暂时注释
	case proto.Message:
		return protoMarshalOptions.Marshal(m)
	default:
		return json.Marshal(m)
	}
}
