package migration

import (
	"github.com/firmeve/firmeve/config"
	"github.com/spf13/cobra"
)

type migration struct {
	config config.Configurator
}

func New(config config.Configurator) *migration {
	return &migration{
		config: config,
	}
}

func (m *migration) Cmd() *cobra.Command {

}
