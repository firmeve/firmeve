## 辅助函数库



### 字符串

`UcFirst` 将首字母大写
```go
UcFirst(`hello`)
// Hello
```
`UcWords` 转换每个单词为大写，并合并单词
```go
UcWords(`hello`, `world`)
// HelloWorld
```
`SnakeCase` 驼峰写法转下划线
```go
SnakeCase(`helloWorld`)
// hello_world
```
`Join` 连接多个字符串
```go
Join(`-`,`hello`,`world`)
// hello-world
```
`HTMLEntity` Html实体转换
```go
HTMLEntity(`&`)
// &amp;
```
`Rand` 随机获取N长度字符串，随机空间为a-z0-9
```go
Rand(10)
// adEasF321D
```
`RandWithCharset` 按指定字符串随机获取
```go
RandWithCharset(5,`abcdef123`)
// abde5
```



## Slice

```go
// 字符串切片去重
func UniqueString(us []string) []string

// int类型数值去重
func UniqueInt(us []int) []int

// 任意类型去重
func UniqueInterface(v interface{}) []interface{}

// 判断给定字符串是否在slice中
func InString(iss []string, v string) bool

// 判断给定Int是否在slice中
func InInt(iss []int, v int) bool 

// 判断任意类型值是否在slice中
func InInterface(iss interface{}, v interface{}) bool
```



### Range

```go
// 获取指定区间的随机数字
func RangeInt(min, max int) int
```



### Path

```go

// 获取当前正在执行的文件路径
func RunFile() string

// 获取当前正在执行的文件目录路径
func RunDir() string

// 获取当前目录的相对路径
func RunRelative(rpath string) string

// 判断文件是否存在
func Exists(path string) bool
```



### Hash

```go
// 生成一个sha1 string
func Sha1String(target string) string

// 返回一个sha1 byte
func Sha1(target string) []byte
```



### Reflect

反射分为反射类型和反射值。

#### Reflect Type

```go
package reflect

import (
	"reflect"
)

// 获取类型的IndirectType 可以自动获取指针类型
func IndirectType(reflectType reflect.Type) reflect.Type

// 获取指定类型的Kind，支持指针
func KindElemType(reflectType reflect.Type) reflect.Kind {
	return IndirectType(reflectType).Kind()
}

type CallInParameterFunc func(i int, param reflect.Type) interface{}

// Call in func params
// It panics if the type's Kind is not Func.
func CallInParameterType(reflectType reflect.Type, paramFunc CallInParameterFunc) []interface{} {
	results := make([]interface{}, 0)
	for i := 0; i < reflectType.NumIn(); i++ {
		values := paramFunc(i, reflectType.In(i))
		results = append(results, values)
	}

	return results
}

type CallFieldFunc func(i int, field reflect.StructField) interface{}

// Call struct fields
// It panics if the type's Kind is not Struct.
func CallFieldType(reflectType reflect.Type, fieldFunc CallFieldFunc) map[string]interface{} {
	reflectType = IndirectType(reflectType)

	results := make(map[string]interface{}, 0)
	for i := 0; i < reflectType.NumField(); i++ {
		reflectField := reflectType.Field(i)
		results[reflectField.Name] = fieldFunc(i, reflectField)
	}

	return results
}

type CallMethodFunc func(i int, method reflect.Method) interface{}

func CallMethodType(reflectType reflect.Type, methodFunc CallMethodFunc) map[string]interface{} {
	results := make(map[string]interface{}, 0)
	for i := 0; i < reflectType.NumMethod(); i++ {
		reflectMethod := reflectType.Method(i)
		results[reflectMethod.Name] = methodFunc(i, reflectMethod)
	}

	return results
}

func Methods(reflectType reflect.Type) map[string]reflect.Method {
	methods := make(map[string]reflect.Method, 0)
	for i := 0; i < reflectType.NumMethod(); i++ {
		method := reflectType.Method(i)
		methods[method.Name] = method
	}

	return methods
}

func MethodExists(reflectType reflect.Type, name string) bool {
	_, ok := reflectType.MethodByName(name)
	return ok
}

// It panics if the type's Kind is not Struct.
func StructFields(reflectType reflect.Type) map[string]reflect.StructField {
	reflectType = IndirectType(reflectType)
	fields := make(map[string]reflect.StructField, 0)
	for i := 0; i < reflectType.NumField(); i++ {
		reflectField := reflectType.Field(i)
		fields[reflectField.Name] = reflectField
	}

	return fields
}

```



#### Reflect Value

```go
package reflect

import (
	"reflect"
)

func CanSetValue(reflectValue reflect.Value) bool {
	return reflectValue.CanSet()
}

func InterfaceValue(reflectType reflect.Type, reflectValue reflect.Value) interface{} {
	if reflectType.Kind() == reflect.Ptr {
		return reflectValue.Addr().Interface()
	}

	return reflectValue.Interface()
}

func SliceInterface(reflectValue reflect.Value) []interface{} {
	newReflectValue := reflect.Indirect(reflectValue)
	kind := newReflectValue.Kind()
	if kind != reflect.Slice && kind != reflect.Array {
		panic(`only support slice or array type`)
	}

	newInterfaces := make([]interface{}, newReflectValue.Len())
	for i := 0; i < newReflectValue.Len(); i++ {
		newInterfaces[i] = newReflectValue.Index(i).Interface()
	}

	return newInterfaces
}

func CallFuncValue(reflectValue reflect.Value, params ...interface{}) []interface{} {
	newParams := make([]reflect.Value, 0)
	for _, param := range params {
		if v, ok := param.(reflect.Value); ok {
			newParams = append(newParams, v)
		} else {
			newParams = append(newParams, reflect.ValueOf(param))
		}
	}

	results := make([]interface{}, 0)
	for _, value := range reflectValue.Call(newParams) {
		current := InterfaceValue(reflect.TypeOf(value), value)

		results = append(results, current)
	}

	return results
}

func CallMethodValue(reflectValue reflect.Value, name string, params ...interface{}) []interface{} {
	return CallFuncValue(reflectValue.MethodByName(name), params...)
}

// FieldByName returns the struct field with the given name.
// It returns the zero Value if no field was found.
// It panics if v's Kind is not struct.
func CallFieldValue(reflectValue reflect.Value, name string) interface{} {
	fieldValue := reflect.Indirect(reflectValue).FieldByName(name)
	return InterfaceValue(reflect.TypeOf(fieldValue), fieldValue)
}

```





