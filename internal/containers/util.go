package containers

import (
	"crypto/rand"
	"encoding/hex"
	"strings"
)

func GenContainerId() string {
	randomBytes := make([]byte, 27)
	_, err := rand.Read(randomBytes)
	if err != nil {
		panic(err)
	}

	// 将随机字节转换为十六进制字符串
	randomHex := hex.EncodeToString(randomBytes)
	//fmt.Println(randomHex)
	//

	return strings.ToUpper(randomHex)[0:10]
}
