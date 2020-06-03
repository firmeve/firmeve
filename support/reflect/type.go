package reflect

import (
	"reflect"
)

func IndirectType(reflectType reflect.Type) reflect.Type {
	kind := reflectType.Kind()

	if kind == reflect.Ptr {
		return reflectType.Elem()
	}

	return reflectType
}

func KindElemType(reflectType reflect.Type) reflect.Kind {
	return IndirectType(reflectType).Kind()
}

type CallInParameterFunc func(i int, param reflect.Type) interface{}

// Call in func params
// It panics if the type's Kind is not Func.
func CallInParameterType(reflectType reflect.Type, paramFunc CallInParameterFunc) []interface{} {
	// make size = reflectType.NumIn()
	results := make([]interface{}, reflectType.NumIn())
	for i := 0; i < reflectType.NumIn(); i++ {
		values := paramFunc(i, reflectType.In(i))
		results[i] = values
		//results = append(results, values)
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
