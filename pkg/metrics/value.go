package metrics

import "strings"

// IsExecutable checks if it is expression to evaluate
func IsExecutable(value string) bool {
	if strings.Contains(value, "{{") && strings.Contains(value, "}}") {
		return true
	}

	return false
}
