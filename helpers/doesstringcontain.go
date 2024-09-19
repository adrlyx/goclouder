package helpers

import "strings"

// Rename FindInString + file name
func DoesStringContain(str string, substr string) bool {
	if strings.Contains(substr, str) {
		return true
	} else {
		return false
	}
}
