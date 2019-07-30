package firmeve

import (
	"github.com/firmeve/firmeve/container"
	"github.com/firmeve/firmeve/server"
	"github.com/firmeve/firmeve/server/http"
	"github.com/gin-gonic/gin"
)

func Run()  {
	//boot
	container.GetFirmeve().Boot()

	container.GetFirmeve().GetContainer().Get(`http.server`).(*http.Http).Server.GET(`/test`, func(context *gin.Context) {
		//fmt.Printf("%#v\n",server.NewContext())
		//c,_ := context2.WithCancel(context)
		z := &http.Context{context}

		z.String(200, `test`)
	})

	//fmt.Printf("%#v",b.firmeve.GetContainer().Get(`server`))
	container.GetFirmeve().GetContainer().Get(`server`).(server.Server).Run()
}