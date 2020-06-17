package firmeve

import (
	"github.com/firmeve/firmeve/kernel/contract"
	"github.com/stretchr/testify/assert"
	"testing"
)

const configPath = "./testdata/config/config.testing.yaml"

func TestRunDefault(t *testing.T) {
	//assert.Nil(t, RunDefault(WithConfigPath(configPath)))
	RunDefault(WithConfigPath(configPath))
}

func TestRunWithSupportFunc(t *testing.T) {
	assert.Nil(t, RunWithSupportFunc(func(application contract.Application) {

	}, WithConfigPath(configPath)))
}
