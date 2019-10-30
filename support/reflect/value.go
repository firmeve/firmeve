package reflect

import (
	"fmt"
	"reflect"
)

func CanSetValue(reflectValue reflect.Value) bool {
	return reflectValue.CanSet()
}

func InterfaceValue(reflectValue reflect.Value) interface{} {
	if reflectValue.CanAddr() {
		return reflectValue.Addr().Interface()
	}

	return reflectValue.Interface()
}

func KindValue(reflectValue reflect.Value) reflect.Kind {
	//if v, ok := object.(reflect.Value); ok {
	//	return v.Kind()
	//}

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
