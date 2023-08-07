package handler

import "strings"

func isValidNameArray(name string) (int, bool) {
	lastIndex := len(name) - 1
	if lastIndex == -1 || name[lastIndex] != ']' {
		return -1, false
	}
	open := strings.Index(name, "[")
	return open, open != -1 && name[lastIndex] == ']' && lastIndex-open == 1
}

func extractArrayName(name string) (string, bool) {
	i, valid := isValidNameArray(name)
	if !valid {
		return name, false
	}

	return name[:i], true
}
