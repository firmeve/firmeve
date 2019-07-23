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

func TestCacheProvider(t *testing.T)  {
	f.Bind("firmeve.provider",new(firmeve.FirmeveServiceProvider))
	serviceProvider := new(cache.CacheServiceProvider)
	fmt.Printf("%#v\n", f.Resolve(serviceProvider).(*cache.CacheServiceProvider).Provider)
	f.Resolve(serviceProvider).(*cache.CacheServiceProvider).Register()

	assert.IsType(t,cache.NewManager(f.Get("config").(*config.Config)),f.Get("cache").(*cache.Manager))
}