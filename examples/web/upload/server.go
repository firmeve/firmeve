package main

import (
	"fmt"
	"github.com/firmeve/firmeve"
	"github.com/firmeve/firmeve/binding"
	"github.com/firmeve/firmeve/kernel"
	"github.com/firmeve/firmeve/kernel/contract"
	"github.com/firmeve/firmeve/render"
)

type App struct {
	kernel.BaseProvider
}

func (a *App) Name() string {
	return `app`
}

func (a *App) Register() {
}

func (a *App) Boot() {
	router := a.Firmeve.Get(`http.router`).(contract.HttpRouter)
	web := router.Group("")
	{
		web.GET("/", func(c contract.Context) {
			c.RenderWith(200, render.Html, render.Template{
				Name:   "index",
				Append: nil,
			})
			c.Next()
		})
		web.POST("/uploads", func(c contract.Context) {
			fmt.Println(c.Protocol().(contract.HttpProtocol).ContentType())
			type FileStruct struct {
				Text  string `form:"text"`
				Files binding.MultipartFiles
			}
			v := new(FileStruct)
			c.Bind(v)
			files, err := v.Files.Save(`file`, &contract.UploadOption{
				Path:    "./files",
				Grading: true,
				Rename:  true,
			})
			if err != nil {
				c.Error(400, err)
				return
			}
			fmt.Printf("%#v\n", files)
			c.RenderWith(200, render.Plain, files[0])
			c.Next()
		})
	}
}

func main() {
	firmeve.RunDefault(firmeve.WithConfigPath("./config.yaml"), firmeve.WithProviders(
		[]contract.Provider{
			new(App),
		},
	))
}
