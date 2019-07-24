package provider

import (
	"fmt"
	"github.com/firmeve/firmeve"
	"github.com/firmeve/firmeve/cache"
	"github.com/firmeve/firmeve/config"
	"github.com/stretchr/testify/assert"
	"testing"
)

var f = firmeve.NewFirmeve("../")

func TestCacheProvider(t *testing.T) {
	serviceProvider := new(cache.ServiceProvider)
	fmt.Printf("%#v\n", f.Resolve(serviceProvider).(*cache.ServiceProvider))
	f.Resolve(serviceProvider).(*cache.ServiceProvider).Register()

	assert.IsType(t, cache.NewManager(f.Get("config").(*config.Config)), f.Get("cache").(*cache.Manager))
}

//func TestHttpProvider(t *testing.T)  {
//	f.Bind("firmeve.provider",new(firmeve.FirmeveServiceProvider))
//	provider := f.Resolve(new(http.HttpServiceProvider)).(*http.HttpServiceProvider)
//	provider.Register()
//	//fmt.Printf("%#v",provider)
//	f.Get("http.server").(*gin.Engine).Run(":22122")
//}



//func TestFirmeve_Run(t *testing.T) {
//	serviceProvider := f.Resolve(new(server.ServiceProvider)).(*server.ServiceProvider)
//	httpServiceProvider := f.Resolve(new(http.ServiceProvider)).(*http.ServiceProvider)
//
//	f.Register(`server`, serviceProvider)
//	f.Register(`http`, httpServiceProvider)
//
//	//firmeve.Run(f)
//	//f.Get(`server`).(server.Server).Run()
//	f.Boot()
//	f.Get(`server`).(server.Server).Run()
//}
