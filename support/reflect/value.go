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
// If is zero value well convert base type zero , not nil
func CallFieldValue(reflectValue reflect.Value, name string) interface{} {
	fieldValue := reflect.Indirect(reflectValue).FieldByName(name)
	return InterfaceValue(reflect.TypeOf(fieldValue), fieldValue)
}

// Call original Field value
// .IsZero() method
func CallOriginalFieldValue(reflectValue reflect.Value, name string) reflect.Value {
	return reflect.Indirect(reflectValue).FieldByName(name)
}
