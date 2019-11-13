package reflect

import (
	"reflect"
)

func IndirectType(reflectType reflect.Type) reflect.Type {
	kind := reflectType.Kind()

	if kind == reflect.Ptr {
		return reflectType.Elem()
	}
	//// type's Kind is Array, Chan, Map, Ptr, or Slice.
	//if kind == reflect.Array || kind == reflect.Ptr || kind == reflect.Chan || kind == reflect.Map || kind == reflect.Slice {
	//	return reflectType.Elem()
	//}

	return reflectType
}

func KindElemType(reflectType reflect.Type) reflect.Kind {
	return IndirectType(reflectType).Kind()
}

//func KindType(object interface{}) reflect.Kind {
//	if t, ok := object.(reflect.Type); ok {
//		return t.Kind()
//	}
//
//	return reflect.TypeOf(object).Kind()
//}

type CallInParameterFunc func(i int, param reflect.Type) interface{}

// Call in func params
// It panics if the type's Kind is not Func.
func CallInParameterType(reflectType reflect.Type, paramFunc CallInParameterFunc) []interface{} {
	results := make([]interface{}, 0)
	for i := 0; i < reflectType.NumIn(); i++ {
		values := paramFunc(i, reflectType.In(i))
		//valueKind := reflect.TypeOf(values).Kind()
		//if valueKind == reflect.Slice {
		//
		//}
		//if len(values) > 1 {
		//	results = append(results, values...)
		//} else {
		results = append(results, values)
		//}
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
