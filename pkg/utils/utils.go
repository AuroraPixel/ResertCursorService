package utils

import (
	"crypto/rand"
	"math/big"
)

// GenerateRandomString 生成指定长度的随机字符串
func GenerateRandomString(length int) string {
	const charset = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	b := make([]byte, length)
	for i := range b {
		randNum, _ := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		b[i] = charset[randNum.Int64()]
	}
	return string(b)
}
