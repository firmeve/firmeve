package firmeve

import (
	"github.com/firmeve/firmeve/container"
	"github.com/firmeve/firmeve/server"
)

func Run()  {
	//boot
	container.GetFirmeve().Boot()
	//fmt.Printf("%#v",b.firmeve.GetContainer().Get(`server`))
	container.GetFirmeve().GetContainer().Get(`server`).(server.Server).Run()
}