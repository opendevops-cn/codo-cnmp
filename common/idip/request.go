package idip

import (
	"encoding/json"
)

// HeadRequest 请求头部结构
type HeadRequest struct {
	GameAppid string `json:"game_appid"`
	BigArea   string `json:"big_area"`
	Cmd       int    `json:"cmd"`
}

// BodyRequest 请求体结构
type BodyRequest struct {
	Lock       bool   `json:"lock"`
	ServerName string `json:"server_name"`
}

// RequestBody 完整的请求结构
type RequestBody struct {
	Head HeadRequest `json:"head"`
	Body BodyRequest `json:"body"`
}

// HeadResponse 响应头部结构
type HeadResponse struct {
	Errno  int    `json:"errno"`
	Errmsg string `json:"errmsg"`
}

// BodyResponse 响应体结构
type BodyResponse struct {
	// 可以根据需要添加字段
}

// ResponseBody 完整的响应结构
type ResponseBody struct {
	Head HeadResponse `json:"head"`
	Body BodyResponse `json:"body"`
}

// RequestOption 请求选项函数类型
type RequestOption func(*RequestBody)

// WithGameAppid 设置 GameAppid
func WithGameAppid(gameAppid string) RequestOption {
	return func(r *RequestBody) {
		r.Head.GameAppid = gameAppid
	}
}

// WithBigArea 设置 BigArea
func WithBigArea(bigArea string) RequestOption {
	return func(r *RequestBody) {
		r.Head.BigArea = bigArea
	}
}

// WithCmd 设置 Cmd
func WithCmd(cmd int) RequestOption {
	return func(r *RequestBody) {
		r.Head.Cmd = cmd
	}
}

// WithLock 设置 Lock
func WithLock(lock bool) RequestOption {
	return func(r *RequestBody) {
		r.Body.Lock = lock
	}
}

// WithServerName 设置 ServerName
func WithServerName(serverName string) RequestOption {
	return func(r *RequestBody) {
		r.Body.ServerName = serverName
	}
}

// NewRequestBody 创建新的请求体
func NewRequestBody(opts ...RequestOption) *RequestBody {
	req := &RequestBody{
		Head: HeadRequest{},
		Body: BodyRequest{},
	}
	// 应用选项
	for _, opt := range opts {
		opt(req)
	}

	return req
}

// String 实现 Stringer 接口
func (r *RequestBody) String() string {
	b, _ := json.Marshal(r)
	return string(b)
}
