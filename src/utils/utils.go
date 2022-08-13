package utils

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
)

func Struct2String(object interface{}) string {
	result, err := json.Marshal(object)
	if err != nil {
		return ""
	}
	return string(result)
}

// GetMd5Value 根据传入的key生成md5值
func GetMd5Value(key string) string {
	m := md5.New()
	m.Write([]byte(key))
	result := hex.EncodeToString(m.Sum([]byte("")))
	return result
}
