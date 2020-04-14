package database

import (
	testing2 "github.com/firmeve/firmeve/testing"
	"testing"

	"github.com/firmeve/firmeve/config"
	"github.com/firmeve/firmeve/support/path"
	"github.com/stretchr/testify/assert"
)

var (
	configPath = "../testdata/config/config.yaml"
)

func TestNew(t *testing.T) {
	db := New(config.New(path.RunRelative(configPath)).Item(`database`))
	assert.NotPanics(t, func() {
		db.ConnectionDefault()
		db.CloseDefault()
	})
	assert.Equal(t, db.ConnectionDefault(), db.ConnectionDefault())
}

func TestNew_Connection_Error(t *testing.T) {
	config := config.New(path.RunRelative(configPath)).Item(`database`)
	config.Set("error_connection.addr", "nothing")
	db := New(config)
	assert.Panics(t, func() {
		db.Connection(`error_connection`)
	})
}

func TestNew_Close_Error(t *testing.T) {
	db := New(config.New(path.RunRelative(configPath)).Item(`database`))
	assert.NotPanics(t, func() {
		db.ConnectionDefault()
		db.CloseDefault()
		db.CloseDefault()
	})
}

func TestDB_Provider(t *testing.T) {
	firmeve := testing2.TestingModeFirmeve()
	firmeve.Register(new(Provider), true)

	firmeve.Boot()
	assert.Equal(t, true, firmeve.HasProvider("db"))
	assert.Equal(t, true, firmeve.Has(`db`))
	assert.Equal(t, true, firmeve.Has(`db.connection`))

	provider := firmeve.Make(new(Provider)).(*Provider)
	assert.Equal(t, firmeve, provider.Firmeve)
}
