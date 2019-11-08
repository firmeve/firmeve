package reflect

import (
	"fmt"
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
	kind := reflectValue.Kind()
	if kind != reflect.Slice && kind != reflect.Array {
		panic(`only support slice or array type`)
	}

	newInterfaces := make([]interface{}, reflectValue.Len())
	for i := 0; i < reflectValue.Len(); i++ {
		newInterfaces[i] = reflectValue.Index(i).Interface()
	}

	return newInterfaces
}

//
//func KindValue(reflectValue reflect.Value) reflect.Kind {
//	return reflectValue.Kind()
//}

func CallFuncValue(reflectValue reflect.Value, params ...interface{}) []interface{} {
	newParams := make([]reflect.Value, 0)
	for _, param := range params {
		newParams = append(newParams, reflect.ValueOf(param))
	}

	results := make([]interface{}, 0)
	for _, value := range reflectValue.Call(newParams) {
		current := InterfaceValue(reflect.TypeOf(value), value)
		if err, ok := current.(error); ok && err != nil {
			panic(fmt.Errorf("call func execute result error. %w", err))
		}

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

//// It panics if the type's Kind is not Struct.
//func StructFieldsValue(reflectValue reflect.Value) map[string]reflect.StructField {
//	reflectType = IndirectType(reflectType)
//
//	fields := make(map[string]reflect.StructField, 0)
//	for i := 0; i < reflectValue.NumField(); i++ {
//		reflectField := reflectValue.Field(i)
//		fields[reflectField.Name] = reflectField
//	}
//
//	return fields
//}
