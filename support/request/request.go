package request

import (
	"github.com/go-playground/form/v4"
	"net/url"
)

var bind = form.NewDecoder()

// 请求数据绑定
func BindWith(v interface{}, values url.Values) error {
	return bind.Decode(v, values)
}
