package microservice

import "strings"

func SubStringLast(str, subStr string) string {
	index := strings.Index(str, subStr)
	if index < 0 {
		return ""
	}
	return str[index+len(subStr):]
}
