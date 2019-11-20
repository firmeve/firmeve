package strings

import (
	"html"
	"regexp"
	"strings"
)

func UcFirst(str string) string {
	first := strings.ToUpper(str[0:1])
	return strings.Join([]string{first, str[1:]}, ``)
}

func UcWords(words []string) string {
	return strings.ReplaceAll(strings.Title(strings.Join(words, ` `)), ` `, ``)
}

func SnakeCase(str string) string {
	return strings.ToLower(regexp.MustCompile("([a-z0-9])([A-Z])").ReplaceAllString(str, "${1}_${2}"))
}

func Join(sep string, str ...string) string {
	return strings.Join(str, sep)
}

func HTMLEntity(str string) string {
	return html.EscapeString(str)
}
