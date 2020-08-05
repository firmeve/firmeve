package kernel

import (
	"fmt"
	"github.com/firmeve/firmeve/kernel/contract"
	reflect2 "github.com/firmeve/firmeve/support/reflect"
	"reflect"
)

// defer
func Recover(logger contract.Loggable, params ...interface{}) {
	RecoverCallback(logger, nil, params...)
}

// defer
func RecoverCallback(logger contract.Loggable, callback func(err interface{}, params ...interface{}), params ...interface{}) {
	if err := recover(); err != nil {
		var message string
		if v, ok := err.(error); ok {
			message = v.Error()
		} else if v, ok := err.(string); ok {
			message = v
		} else {
			message = `unknown`
		}

		// merge all record params
		newParams := []interface{}{message, "error", err}
		// get params reflect name
		for i := range params {
			newParams = append(newParams, reflect2.IndirectType(reflect.TypeOf(params[i])).Name(), fmt.Sprintf("%v", params[i]))
		}
		logger.Error(newParams...)

		if callback != nil {
			callback(err, params...)
		}
	}
}
