package utils

import (
	"bytes"
	"net"
	"net/http"
	"time"
)

// DoRequest 封装HTTP请求
func DoRequest(method, url string, headers map[string]string, body []byte) (*http.Response, error) {
	req, err := http.NewRequest(method, url, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	// 设置请求头
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	httpClient := &http.Client{
		Timeout: time.Second * 10, // 请求超时
		Transport: &http.Transport{
			MaxIdleConns:        100,
			MaxIdleConnsPerHost: 100,
			IdleConnTimeout:     90 * time.Second,
			// 自定义DialContext，设置拨号超时时间和保持连接
			DialContext: (&net.Dialer{
				Timeout:   5 * time.Second,  // 拨号超时
				KeepAlive: 30 * time.Second, // 保持连接存活
			}).DialContext,
			TLSHandshakeTimeout: 10 * time.Second, // TLS 握手超时
		},
	}

	// 执行请求
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	// 检查响应并确保关闭body
	if resp.Body != nil {
		defer resp.Body.Close()
	}

	return resp, nil
}
