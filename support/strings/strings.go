package strings

import (
	"html"
	"math/rand"
	"regexp"
	"strings"
	"time"
)

const charset = `abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789`

func UcFirst(str string) string {
	first := strings.ToUpper(str[0:1])
	return strings.Join([]string{first, str[1:]}, ``)
}

func UcWords(words ...string) string {
	return strings.ReplaceAll(strings.Title(Join(` `, words...)), ` `, ``)
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

func RandWithCharset(length int, charset string) string {
	var (
		seededRand = rand.New(rand.NewSource(time.Now().UnixNano()))
		b = make([]byte,length)
	)

	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}

	return string(b)
}

func Rand(length int) string {
	return RandWithCharset(length, charset)
}