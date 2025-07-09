package idip

import (
	"crypto/md5"
	"fmt"
	"math/rand"
	"strings"
	"time"
)

// Signer 签名类，用于生成签名字符串
type Signer struct {
	secret    string
	appId     string
	timestamp int64
	rnd       int
}

// NewSigner 创建新的签名器
func NewSigner(appId, secret string) *Signer {
	return &Signer{
		secret:    secret,
		appId:     appId,
		timestamp: getCurrentTimestamp(),
		rnd:       generateRandInt(),
	}
}

// getCurrentTimestamp 获取当前时间戳
func getCurrentTimestamp() int64 {
	return time.Now().Unix()
}

// generateRandInt 生成随机数
func generateRandInt() int {
	return rand.Intn(1000)
}

// GetRnd 获取随机数
func (s *Signer) GetRnd() int {
	return s.rnd
}

// GetTimestamp 获取时间戳
func (s *Signer) GetTimestamp() int64 {
	return s.timestamp
}

// Sign 生成签名字符串
func (s *Signer) Sign(body, tag string) string {
	// 构建签名字符串
	signStr := fmt.Sprintf("%s%s%d%d%s", tag, body, s.timestamp, s.rnd, s.secret)

	// 计算MD5
	hash := md5.Sum([]byte(signStr))
	md5Str := fmt.Sprintf("%x", hash)

	return md5Str
}

// SignHeader 表示签名头部信息
type SignHeader struct {
	AppId   string `json:"cbb-sign-appid"`
	Time    string `json:"cbb-sign-time"`
	Rnd     string `json:"cbb-sign-rnd"`
	Sign    string `json:"cbb-sign"`
	Version string `json:"cbb-sign-version"`
}

// GenSignHeader 生成签名header
func (s *Signer) GenSignHeader(body, tag string, version string) SignHeader {
	if version == "" {
		version = "v1"
	}
	signStr := s.Sign(body, tag)

	return SignHeader{
		AppId:   s.appId,
		Time:    fmt.Sprintf("%d", s.timestamp),
		Rnd:     fmt.Sprintf("%d", s.rnd),
		Sign:    strings.ToUpper(strings.ReplaceAll(signStr, "-", "")),
		Version: version,
	}
}

// HeaderToMap 将签名头部信息转换为map
func (h SignHeader) HeaderToMap() map[string]string {
	return map[string]string{
		"cbb-sign-appid":   h.AppId,
		"cbb-sign-time":    h.Time,
		"cbb-sign-rnd":     h.Rnd,
		"cbb-sign":         h.Sign,
		"cbb-sign-version": h.Version,
	}
}
