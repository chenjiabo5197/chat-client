package utils

import (
	"encoding/json"
)

func Struct2String(object interface{}) string {
	result, err := json.Marshal(object)
	if err != nil {
		return ""
	}
	return string(result)
}
