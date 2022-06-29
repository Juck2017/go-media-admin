package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
)

func Struct2JsonString(o interface{}) string {
	b, err := json.Marshal(o)
	if err != nil {
		return fmt.Sprintf("%+v", o)
	}
	var out bytes.Buffer
	err = json.Indent(&out, b, "", "  ")
	if err != nil {
		return fmt.Sprintf("%+v", o)
	}
	return out.String()
}
