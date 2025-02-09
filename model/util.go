package model

import (
	"fmt"
)

func dbValue(in any) string {
	if v, ok := in.([]byte); ok {
		return string(v)
	}
	if in == nil {
		return ""
	}
	return fmt.Sprint(in)
}
