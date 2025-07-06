package utils

import "strings"

func IsEmptyOrBlank(lines []string) bool {
	return len(lines) == 0 || (len(lines) == 1 && lines[0] == "")
}

func EndsWithoutNewline(lines []string) bool {
	return len(lines) > 0 && !strings.HasSuffix(lines[len(lines)-1], "\n")
}
