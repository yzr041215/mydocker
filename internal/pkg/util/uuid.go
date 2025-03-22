package util

import (
	"crypto/rand"
	"encoding/hex"
	"strings"
)

func GenerateUUID() string {
	randomBytes := make([]byte, 32)
	_, err := rand.Read(randomBytes)
	if err != nil {
		panic(err)
	}

	// 将随机字节转换为十六进制字符串
	randomHex := hex.EncodeToString(randomBytes)
	//fmt.Println(randomHex)
	return randomHex
}
func GenerateIinkUUID() string {
	randomBytes := make([]byte, 27)
	_, err := rand.Read(randomBytes)
	if err != nil {
		panic(err)
	}

	// 将随机字节转换为十六进制字符串
	randomHex := hex.EncodeToString(randomBytes)
	//fmt.Println(randomHex)
	//

	return strings.ToUpper(randomHex)[0:26]
}
