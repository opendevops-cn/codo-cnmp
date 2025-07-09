package utils

import (
	"bytes"
	"encoding/base32"
	"fmt"
	"hash/crc32"
	"strings"
)

func UInt32ToBytes(u32 uint32) []byte {
	return []byte{byte(u32 >> 24), byte(u32 >> 16), byte(u32 >> 8), byte(u32)}

}

func PathEscape(token string) string {
	return strings.ReplaceAll(token, "/", "-")
}

func PathUnescape(token string) string {
	return strings.ReplaceAll(token, "-", "/")
}

func EncodeToken(bs []byte) string {
	for i := range bs {
		bs[i] = ^bs[i]
	}
	return base32.HexEncoding.WithPadding(base32.NoPadding).EncodeToString(bs)
}

func DecodeToken(token string) ([]byte, error) {
	bs, err := base32.HexEncoding.WithPadding(base32.NoPadding).DecodeString(token)
	if err != nil {
		return nil, err
	}
	for i := range bs {
		bs[i] = ^bs[i]
	}
	return bs, nil
}

func EncodeTokenEscape64(bs []byte) string {
	return PathEscape(EncodeToken(bs))
}

func DecodeTokenEscape64(token string) ([]byte, error) {
	token = PathUnescape(token)
	return DecodeToken(token)
}

func EncodeTokenWithChecksumEscape64(bs []byte, key []byte) string {
	u32 := crc32.ChecksumIEEE(append(append(bs, key...)))
	bs = append(bs, UInt32ToBytes(u32)...)
	return EncodeTokenEscape64(bs)
}

func BytesToUInt32(bs []byte) uint32 {
	if len(bs) < 4 {
		return 0
	}
	return uint32(bs[0])<<24 | uint32(bs[1])<<16 | uint32(bs[2])<<8 | uint32(bs[3])
}

func DecodeTokenWithCheckSumEscape64(token string, key []byte) ([]byte, error) {
	bs, err := DecodeTokenEscape64(token)
	if err != nil {
		return nil, err
	}
	if len(bs) < 8 {
		return nil, fmt.Errorf("%s, invalid token", token)
	}

	// 检查报文完整性
	got := BytesToUInt32(bs[len(bs)-4:])
	bs = bs[:len(bs)-4]
	need := crc32.ChecksumIEEE(append(bs, key...))
	if need != got {
		return nil, fmt.Errorf("%s, checksum error", token)
	}
	return bs, nil
}

func GenerateToken(userId uint32, clusterId uint32, signKey string) (string, error) {
	userID := UInt32ToBytes(userId)
	clusterID := UInt32ToBytes(clusterId)
	token := EncodeTokenWithChecksumEscape64(bytes.Join([][]byte{userID, clusterID}, nil), []byte(signKey))
	return token, nil
}
