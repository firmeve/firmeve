package reflect

import (
	"reflect"
)

func IndirectType(reflectType reflect.Type) reflect.Type {
	kind := reflectType.Kind()

	// type's Kind is Array, Chan, Map, Ptr, or Slice.
	if kind == reflect.Array || kind == reflect.Ptr || kind == reflect.Chan || kind == reflect.Map || kind == reflect.Slice {
		return reflectType.Elem()
	}

	return reflectType
}

func KindType(reflectType reflect.Type) reflect.Kind {
	return reflectType.Kind()
}

//func KindType(object interface{}) reflect.Kind {
//	if t, ok := object.(reflect.Type); ok {
//		return t.Kind()
//	}
//
//	return reflect.TypeOf(object).Kind()
//}

type CallParameterFunc func(i int, param reflect.Type) interface{}

// Call func params
// It panics if the type's Kind is not Func.
func CallParameterType(reflectType reflect.Type, paramFunc CallParameterFunc) []interface{} {

	results := make([]interface{}, 0)
	for i := 0; i < reflectType.NumIn(); i++ {
		results = append(results, paramFunc(i, reflectType.In(i)))
	}

	return results
}

type CallFieldFunc func(i int, field reflect.StructField) interface{}

// Call struct fields
// It panics if the type's Kind is not Struct.
func CallFieldType(reflectType reflect.Type, fieldFunc CallFieldFunc) []interface{} {

	results := make([]interface{}, 0)
	for i := 0; i < reflectType.NumField(); i++ {
		results = append(results, fieldFunc(i, reflectType.Field(i)))
	}

	return results
}
