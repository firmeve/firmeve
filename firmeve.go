package firmeve

import (
	"fmt"
	"go/importer"
	"reflect"
	"strings"
	"sync"
)

var (
	firmeve *Firmeve
	once    sync.Once
)

type Container interface {
	Get(id string) interface{}
	Has(id string) bool
	//Bind(id interface{})
}

type binding struct {
	shared bool
	instance interface{}
	object func() interface{}
}

type Firmeve struct {
	bindings map[string]*binding
}

// Create a new firmeve container
func NewFirmeve() *Firmeve {
	if firmeve != nil {
		return firmeve
	}

	once.Do(func() {
		firmeve = &Firmeve{
			bindings: make(map[string]*binding),
		}
	})

	return firmeve
}

func (f *Firmeve) Bind(object func() interface{},share bool) { //, value interface{}
	//f.bindings[id] = newBinding(false, value)
	reflectType := reflect.TypeOf(object)
	if reflectType.Kind() == reflect.Ptr {
		fullName := strings.Join([]string{reflectType.Elem().PkgPath(), reflectType.Elem().Name()}, ".")
		f.bindings[fullName] = newBinding(share, object)
	}

}

func (f *Firmeve) Resolve(id interface{}) interface{}{
	reflectType := reflect.TypeOf(id)
	if reflectType.Kind() == reflect.Func {

		//params1Type := reflectType.In(0)//.Elem().Name()
		//z := nil
		////params1Value := &params1Type//reflect.NewAt(params1Type,unsafe.Pointer(params1Type.))
		//fmt.Printf("%#v",reflect.ValueOf(z).Convert(params1Type))

		// 重新设置一个新的struct对象，不过是nil
		//newObjPtr := reflect.New(params1Type)
		//newObj := reflect.Indirect(newObjPtr)

		//importer.Default()
		//fmt.Println(params1Type.Elem().PkgPath())
		pkg,err := importer.Default().Import(`github.com/firmeve/firmeve/demo`)
		if err != nil{
			fmt.Println(err.Error())
		}

		for _, declName := range pkg.Scope().Names() {
			fmt.Println(declName)
		}


		//println(params1Type,reflect.ValueOf(params1Value).Interface().(*T1))




		//fullName := strings.Join([]string{reflectType.In(0).Elem().PkgPath(), reflectType.In(0).Elem().Name()}, ".")
		//fmt.Println(fullName)
		//fmt.Printf("%#v",reflectType.In(0).Field(0))
		//stk := reflectType.In(0)
		//stkValue := reflect.ValueOf(reflect.TypeOf(id).In(0))
		////reflect.New(stk)
		////fmt.Printf("%#v",stkValue.MethodByName("NewT1").Call(make([]reflect.Value,0))[0].Interface())
		//fmt.Printf("%#v\n",stk.Elem().PkgPath())
		////zx := reflect.ValueOf(NewT1())
		//stkValue.Elem().Set(reflect.New(stk))
		//reflect.NewAt(stk, unsafe.Pointer(stk.UnsafeAddr())).Elem()
		//fmt.Printf("%#v",stkValue.Interface())
		//fmt.Printf("%#v",stkValue.Call([]reflect.Value{reflect.ValueOf("NewT1")}))

		//return reflect.ValueOf(id).Call([]reflect.Value{reflect.ValueOf(f.bindings[fullName].object)})[0].Interface()
	}

	return  nil
}

// ---------------------------- binding ------------------------

func newBinding(shared bool, object interface{}) *binding {
	return &binding{
		shared: shared,
		object: object,
	}
}
