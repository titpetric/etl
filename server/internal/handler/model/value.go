package model

import "fmt"

// DBValue converts any database value into a string representation.
func DBValue(in any) string {
	if v, ok := in.([]byte); ok {
		return string(v)
	}
	if in == nil {
		return ""
	}
	return fmt.Sprint(in)
}
