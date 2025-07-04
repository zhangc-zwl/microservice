package microservice

import (
	"strings"
	"unicode"
)

func SubStringLast(str, subStr string) string {
	index := strings.Index(str, subStr)
	if index < 0 {
		return ""
	}
	return str[index+len(subStr):]
}

func isASCII(s string) bool {
	for i := 0; i < len(s); i++ {
		if s[i] > unicode.MaxASCII {
			return false
		}
	}
	return true
}
