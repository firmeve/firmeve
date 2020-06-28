package render

import (
	"fmt"
	"github.com/firmeve/firmeve/kernel/contract"
)

type (
	errorRender struct {
		Message string                 `json:"message,omitempty"`
		Stack   string                 `json:"stack,omitempty"`
		Meta    map[string]interface{} `json:"meta,omitempty"`
	}
)

var Error = errorRender{}

func (e errorRender) Render(protocol contract.Protocol, status int, v interface{}) error {
	switch v := v.(type) {
	case contract.Error:
		e.Message = v.String()
		e.Meta = v.Meta()
	case error:
		e.Message = v.Error()
	default:
		e.Message = fmt.Sprintf("%v", v)
	}

	// append stack
	if v, ok := v.(contract.ErrorStack); ok {
		e.Stack = v.StackString()
	}

	return Render(protocol, status, e)
}

func (e errorRender) String() string {
	return e.Message
}
