package firmeve

import (
	cache2 "github.com/firmeve/firmeve/cache"
	"testing"
)

//func TestFirmeve_Bind(t *testing.T) {
//
//	bak := demo.Test()
//
//	firmeve := NewFirmeve()
//
//	firmeve.Bind(`bak`,bak)
//
//	firmeve.Bind(`bak2`,bak2.Bak{})
//	firmeve.Bind(`bak3`,"string")
//
//}
//
func TestRelectTest(t *testing.T) {
	//fmt.Println(demo.Test2())
	//var x string
	//RelectTest(x)
	//RelectTest(demo.Test)
	RelectTest(cache2.NewRepository)
}
//
////func TestFirmeve_Bind_Resolve(t *testing.T) {
////	firmeve := NewFirmeve()
////	firmeve.Bind("bak",&bak.Bak{Title:"abc",At:20})
////
////	t.Logf("%#v",firmeve.Resolve("bak").(*bak.Bak))
////}