//go:generate mockgen -package mock -destination ../../testing/mock/mock_logger.go github.com/firmeve/firmeve/kernel/contract Loggable
package contract

import "io"

type (
	Loggable interface {
		Writer(channel string) io.Writer
		With(...interface{}) Loggable
		Debug(...interface{})
		Info(...interface{})
		Warn(...interface{})
		Error(...interface{})
		Fatal(...interface{})
		Panic(...interface{})
		//Channel(stack string) Loggable
	}
)
