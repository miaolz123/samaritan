package candyjs

import "strings"

func isExported(name string) bool {
	return len(name) > 0 && name[0] >= 'A' && name[0] <= 'Z'
}

func nameToJavaScript(name string) string {
	return name
}

func nameToGo(name string) []string {
	return []string{
		strings.Title(name),
		name,
	}
}
