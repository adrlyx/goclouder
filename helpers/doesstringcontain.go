package helpers

import "strings"

func DoesStringContain(str string, substr string) bool {
	if strings.Contains(substr, str) {
		return true
	} else {
		return false
	}
}
