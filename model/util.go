package model

import (
	"unicode"
)

// toUnderScoreCase myName => my_name
func toUnderScoreCase(s string) string {
	runes := []rune(s)
	length := len(runes)
	out := []rune{}
	for i := 0; i < length; i++ {
		if i > 0 && unicode.IsUpper(runes[i]) && ((i+1 < length && unicode.IsLower(runes[i+1])) || unicode.IsLower(runes[i-1])) {
			out = append(out, '_')
		}
		out = append(out, unicode.ToLower(runes[i]))
	}
	return string(out)
}
