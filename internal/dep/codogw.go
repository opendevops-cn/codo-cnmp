package dep

import (
	"bytes"
	"codo-cnmp/internal/conf"
	"context"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

const (
	CODOAPIGatewayAuthKeyHeader = "auth_key"
)

type CODOAPIGateway struct {
	BaseURL string
	client  *http.Client
	cookies []*http.Cookie
}

func NewCODOAPIGateway(ctx context.Context, bc *conf.Bootstrap) *CODOAPIGateway {
	cfg := bc.TIANMEN
	// 从配置中获取 auth_key 并设置为 Cookie
	authCookie := &http.Cookie{
		Name:  CODOAPIGatewayAuthKeyHeader,
		Value: cfg.AUTH_KEY, // 从配置读取 auth_key
	}
	timeout := time.Duration(cfg.Timeout) * time.Second
	if timeout == 0 {
		timeout = 30 * time.Second // 默认超时时间为 30 秒
	}

	return &CODOAPIGateway{
		client: &http.Client{
			Timeout: timeout,
		},
		BaseURL: cfg.ADDR,
		cookies: []*http.Cookie{authCookie},
	}
}

// SetCookies 设置统一的 Cookie
func (x *CODOAPIGateway) SetCookies(cookies []*http.Cookie) {
	x.cookies = cookies
}

func (x *CODOAPIGateway) SendRequest(ctx context.Context, method string, endpoint string, body []byte, queryParams map[string]string) (*http.Response, error) {
	// 构建完整的 URL，包含 query 参数
	fullURL, err := x.buildURLWithParams(endpoint, queryParams)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(method, fmt.Sprintf("%s%s", x.BaseURL, fullURL), bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	// 设置默认 Header
	req.Header.Set("Content-Type", "application/json")

	// 自动添加统一设置的 Cookie
	for _, cookie := range x.cookies {
		req.AddCookie(cookie)
	}
	// 发送请求
	return x.client.Do(req)
}

func (x *CODOAPIGateway) buildURLWithParams(endpoint string, queryParams map[string]string) (string, error) {
	u, err := url.Parse(endpoint)
	if err != nil {
		return "", err
	}
	// 如果 endpoint 已经带有查询参数，保留它们
	q := u.Query()
	for key, value := range queryParams {
		q.Add(key, value)
	}
	u.RawQuery = q.Encode()

	return u.String(), nil
}
