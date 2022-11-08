package gorm_masking

import (
	"encoding/base64"
	"github.com/duke-git/lancet/v2/cryptor"
	"strings"
)

func generate(len int) string {
	if len == 0 {
		return ""
	}
	str := strings.Builder{}
	for i := 0; i < len; i++ {
		str.WriteString("*")
	}
	return str.String()
}

func encrypt(str, key string) string {
	return base64.StdEncoding.EncodeToString(cryptor.AesEcbEncrypt([]byte(str), []byte(key)))
}

func decryptValue(str, key string) string {
	rs, _ := base64.StdEncoding.DecodeString(str)
	v := cryptor.AesEcbDecrypt(rs, []byte(key))
	return string(v)
}
