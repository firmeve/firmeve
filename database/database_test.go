package database

import (
	"github.com/firmeve/firmeve/kernel/contract"
	testing2 "github.com/firmeve/firmeve/testing"
	"testing"

	"github.com/firmeve/firmeve/config"
	"github.com/stretchr/testify/assert"
)

var (
	app contract.Application
)

func TestMain(m *testing.M) {
	//set up
	app = testing2.ApplicationDefault(new(Provider))

	// testing
	m.Run()

	//teardown
}

//func TestNew(t *testing.T) {
//	db := New(config.New(path.RunRelative(configPath)).Item(`database`))
//	assert.NotPanics(t, func() {
//		db.ConnectionDefault()
//		db.CloseDefault()
//	})
//	assert.Equal(t, db.ConnectionDefault(), db.ConnectionDefault())
//}

func TestNew_Connection_Error(t *testing.T) {
	//config := config.New(path.RunRelative(configPath)).Item(`database`)
	app.Resolve(`config`).(*config.Config).Item(`database`).Set("error_connection.addr", "nothing")
	assert.Panics(t, func() {
		app.Resolve(`db`).(*DB).Connection(`error_connection`)
	})
}

func TestNew_Close_Error(t *testing.T) {
	db := app.Resolve(`db`).(*DB)
	assert.NotPanics(t, func() {
		db.ConnectionDefault()
		db.CloseDefault()
		db.CloseDefault()
	})
}

func TestDB_Provider(t *testing.T) {
	firmeve := testing2.ApplicationDefault(new(Provider))

	firmeve.Boot()
	assert.Equal(t, true, firmeve.HasProvider("db"))
	assert.Equal(t, true, firmeve.Has(`db`))
	assert.Equal(t, true, firmeve.Has(`db.connection`))

	provider := firmeve.Make(new(Provider)).(*Provider)
	assert.Equal(t, firmeve, provider.Firmeve)
}
