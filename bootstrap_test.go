package firmeve

import (
	"github.com/firmeve/firmeve/container"
	"github.com/firmeve/firmeve/server"
	"testing"
)
import _ "net/http/pprof"

var basePath = "./testdata"

func TestFirmeve_Run(t *testing.T) {
	// boot
	container.GetFirmeve().Boot()
	//fmt.Printf("%#v",b.firmeve.GetContainer().Get(`server`))
	container.GetFirmeve().GetContainer().Get(`server`).(server.Server).Run()
}


