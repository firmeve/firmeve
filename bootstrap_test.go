package firmeve

import (
	"testing"
)
import _ "net/http/pprof"

var basePath = "./testdata"

func TestFirmeve_Run(t *testing.T) {
	NewBootstrap(basePath).RegisterService().Run()
}


