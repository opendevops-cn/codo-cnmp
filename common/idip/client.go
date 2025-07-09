package idip

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/opendevops-cn/codo-golang-sdk/client/xhttp"
)

type IDIPBaseAPI struct {
	signer    *Signer
	cc        xhttp.IClient
	baseURL   string
	appId     string // 应用ID, 用于区分不同游戏
	gameAppId string // 每个idip对应的密钥
	secret    string
	log       *log.Helper
}

func NewIDIPBaseAPI(idip string, gameAppId string, appId string, secret string, timeout time.Duration, logger *log.Helper) (*IDIPBaseAPI, error) {
	if idip == "" || appId == "" || secret == "" || gameAppId == "" {
		return nil, fmt.Errorf("idip, appID, secret, gameAppId 不能为空")
	}
	if timeout <= 0 {
		timeout = time.Second * 10
	}
	client, err := xhttp.NewClient(xhttp.WithClientOptionsTimeout(timeout))
	if err != nil {
		panic(err)
	}
	signer := NewSigner(appId, secret)
	return &IDIPBaseAPI{
		signer:    signer,
		cc:        client,
		baseURL:   fmt.Sprintf("%s/idip/", idip),
		appId:     appId,
		gameAppId: gameAppId,
		secret:    secret,
		log:       logger,
	}, nil
}

func (x *IDIPBaseAPI) PrepareRequest(ctx context.Context, body *RequestBody) (*http.Request, error) {
	// 生成签名头
	headers := map[string]string{
		"Content-Type": "application/json",
		"Accept":       "application/json",
	}
	if x.signer != nil {
		signHeader := x.signer.GenSignHeader(body.String(), "", "v1")
		for k, v := range signHeader.HeaderToMap() {
			headers[k] = v
		}
	}

	// 序列化请求体
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("序列化请求体失败: %w", err)
	}

	// 创建请求
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, x.baseURL, bytes.NewReader(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("创建HTTP请求失败: %w", err)
	}

	// 设置请求头
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	// 记录请求信息
	x.log.WithContext(ctx).Infof("idip请求url: %s, headers: %v, body: %s", x.baseURL, headers, string(jsonBody))

	return req, nil
}

func (x *IDIPBaseAPI) HandleResponse(ctx context.Context, resp *http.Response) (*ResponseBody, error) {
	// 读取响应体
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应体失败: %w", err)
	}

	// 记录响应信息
	x.log.WithContext(ctx).Infof("idip响应状态码: %s, 响应体: %s", resp.Status, string(respBody))

	// 检查HTTP状态码
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP请求失败: %s", resp.Status)
	}

	// 解析响应体
	var response ResponseBody
	if err := json.Unmarshal(respBody, &response); err != nil {
		return nil, fmt.Errorf("解析响应体失败: %w", err)
	}

	// 检查业务状态码
	if response.Head.Errno != 0 {
		return &response, fmt.Errorf("IDIP业务错误: %s", response.Head.Errmsg)
	}

	return &response, nil
}

func (x *IDIPBaseAPI) SendRequest(ctx context.Context, body *RequestBody) (*ResponseBody, error) {
	// 1. 准备请求
	req, err := x.PrepareRequest(ctx, body)
	if err != nil {
		return nil, fmt.Errorf("准备请求失败: %w", err)
	}

	// 2. 发送请求
	resp, err := x.cc.Do(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("发送HTTP请求失败: %w", err)
	}
	defer resp.Body.Close()

	// 3. 处理响应
	return x.HandleResponse(ctx, resp)
}
