//go:generate mockgen -package mock -destination ../../testing/mock/mock_render.go github.com/firmeve/firmeve/kernel/contract Render
package contract

type (
	Render interface {
		Render(protocol Protocol, status int, v interface{}) error
	}
)
