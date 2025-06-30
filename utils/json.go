//revive:disable:var-naming
package utils

import (
	"encoding/json"
)

func ToJSONString(data any) string {
	b, err := json.Marshal(data)
	if err != nil {
		return ""
	}
	return string(b)
}
