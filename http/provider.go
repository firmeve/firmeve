package http

import (
	"github.com/firmeve/firmeve/container"
	"github.com/firmeve/firmeve/kernel"
	"github.com/spf13/cobra"
)

type Provider struct {
	kernel.BaseProvider
}

func (p *Provider) Name() string {
	return `http`
}

func (p *Provider) Register() {
	p.Firmeve.Bind(`http.router`, New(p.Firmeve), container.WithShare(true))

	p.Firmeve.Get("command").(*cobra.Command).AddCommand(
		NewServer(p.Firmeve).Cmd(),
	)
}

func (p *Provider) Boot() {

}
