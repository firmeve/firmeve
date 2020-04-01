package grpc

import "github.com/firmeve/firmeve/kernel"

type Provider struct {
	kernel.BaseProvider
}

func (p Provider) Name() string {
	return `grpc`
}

func (p *Provider) Register() {
	p.Firmeve.Bind(`grpc.server.router`, NewRouter(p.Firmeve))
}

func (p Provider) Boot() {

}
