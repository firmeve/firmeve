package reflect

import (
	"fmt"
	"reflect"
)

func CanSetValue(reflectValue reflect.Value) bool {
	return reflectValue.CanSet()
}

func InterfaceValue(reflectValue reflect.Value) interface{} {
	//@todo 这里可能会有问题
	//@todo 需要仔细研究CanAddr
	if reflectValue.CanAddr() {//&& KindValue(reflectValue) <= reflect.Complex128
		return reflectValue.Addr().Interface()
	}

	return reflectValue.Interface()
}

func KindValue(reflectValue reflect.Value) reflect.Kind {
	return reflectValue.Kind()
}

func CallFuncValue(reflectValue reflect.Value, params ...interface{}) []interface{} {
	newParams := make([]reflect.Value, 0)
	for _, param := range params {
		newParams = append(newParams, reflect.ValueOf(param))
	}

	results := make([]interface{}, 0)
	for _, value := range reflectValue.Call(newParams) {
		current := InterfaceValue(value)
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
	// 如果父级是指针，结构体字段则也是支持CanAddr()
	// 所以需要再包装
	var fieldValue reflect.Value
	if KindValue(reflectValue) == reflect.Ptr {
		return reflect.Indirect(reflectValue).FieldByName(name).Interface()
	} else {
		fieldValue = reflect.Indirect(reflectValue).FieldByName(name)
	}

	return InterfaceValue(fieldValue)
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
