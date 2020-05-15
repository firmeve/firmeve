package render

import (
	"fmt"
	"github.com/firmeve/firmeve/kernel/contract"
	"io/ioutil"
	"os"
)

type (
	file struct {
	}

	FileOption struct {
		Mime     string
		Path     string
		Filename string
	}
)

var (
	File = file{}
)

func (f file) Render(protocol contract.Protocol, status int, v interface{}) error {
	httpProtocol := protocol.(contract.HttpProtocol)
	option := v.(FileOption)

	file, err := os.Open(option.Path)
	if err != nil {
		httpProtocol.ResponseWriter().WriteHeader(404)
		return err
	}
	httpProtocol.ResponseWriter().WriteHeader(status)
	httpProtocol.SetHeader(`Content-Type`, option.Mime)
	httpProtocol.SetHeader("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", option.Filename))

	defer file.Close()
	c, _ := ioutil.ReadAll(file)
	_, err = httpProtocol.Write(c)
	return err
}
