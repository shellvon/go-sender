package utils

import (
	"encoding/json"
)

// ToJSONString 将 map[string]string 转换为 JSON 字符串
// 如果 map 为空或序列化失败，返回空字符串
func ToJSONString(m map[string]string) string {
	if len(m) == 0 {
		return ""
	}

	b, err := json.Marshal(m)
	if err != nil {
		return ""
	}

	return string(b)
}
