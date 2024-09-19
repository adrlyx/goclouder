package helpers

import "strings"

func FindInString(str string, substr string) bool {
	if strings.Contains(substr, str) {
		return true
	} else {
		return false
	}
}
