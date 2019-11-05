package strings

import (
	"regexp"
	"strings"
)

func UcFirst(str string) string {
	first := strings.ToUpper(str[0:1])
	return strings.Join([]string{first, str[1:]}, ``)
}

func UcWords(words []string) string {
	return strings.ReplaceAll(strings.Title(strings.Join(words,` `)),` `,``)
}

func SnakeCase(str string) string {
	return strings.ToLower(regexp.MustCompile("([a-z0-9])([A-Z])").ReplaceAllString(str, "${1}_${2}"))
}

