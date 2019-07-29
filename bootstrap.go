package firmeve

//import (
//	"github.com/firmeve/firmeve/cache"
//	"github.com/firmeve/firmeve/container"
//	"github.com/firmeve/firmeve/server"
//	"github.com/firmeve/firmeve/server/http"
//	"sync"
//)
//
//var (
//	bootstrap *Bootstrap
//	once      sync.Once
//)
//
//type Bootstrap struct {
//	firmeve *container.Firmeve
//}
//
//func NewBootstrap(basePath string) *Bootstrap {
//	if bootstrap != nil {
//		return bootstrap
//	}
//	once.Do(func() {
//		bootstrap = &Bootstrap{
//			firmeve: container.NewFirmeve(basePath),
//		}
//	})
//	return bootstrap
//}
//
//func (b *Bootstrap) RegisterService() *Bootstrap {
//	b.firmeve.Register(`cache`, b.firmeve.GetContainer().Resolve(new(cache.ServiceProvider)).(*cache.ServiceProvider))
//	b.firmeve.Register(`server`, b.firmeve.GetContainer().Resolve(new(server.ServiceProvider)).(*server.ServiceProvider))
//	b.firmeve.Register(`http`, b.firmeve.GetContainer().Resolve(new(http.ServiceProvider)).(*http.ServiceProvider))
//	return b
//}
//
//func (b *Bootstrap) Run() {
//	// boot
//	b.firmeve.Boot()
//	//fmt.Printf("%#v",b.firmeve.GetContainer().Get(`server`))
//	b.firmeve.GetContainer().Get(`server`).(server.Server).Run()
//}
