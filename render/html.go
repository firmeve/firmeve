package render

import (
	"fmt"
	"github.com/firmeve/firmeve/kernel/contract"
	"github.com/spf13/viper"
	template2 "html/template"
	path2 "path"
	"strings"
)

type (
	html struct {
	}

	Template struct {
		Name   string
		Data   interface{}
		Append []string //附加模板，多个模板关联，后期测试是否需要
	}
)

var (
	Html     = html{}
	basePath string
	suffix   string
)

func (html) Render(protocol contract.Protocol, status int, v interface{}) error {
	if p, ok := protocol.(contract.HttpProtocol); ok {
		p.ResponseWriter().WriteHeader(status)
		p.SetHeader(`Content-Type`, `text/html`)

		// template parse
		if tmpl, ok := v.(Template); ok {
			// Get base views config
			if basePath == "" || suffix == "" {
				config := p.Application().Resolve(`config`).(contract.Configuration).Get("view").(*viper.Viper)
				basePath = config.GetString("path")
				suffix = config.GetString("suffix")
			}

			// template name to conver path
			path := strings.ReplaceAll(tmpl.Name, ".", "/")
			fullPath := path2.Join(basePath, path+suffix)
			// @todo mutil files parse
			t, err := template2.ParseFiles(fullPath)
			if err != nil {
				return err
			}
			return t.Execute(p.ResponseWriter(), tmpl.Data)
		}
	}

	// except tmpl parse
	_, err := protocol.Write([]byte(fmt.Sprintf("%v", v)))

	return err
}
