package database

import (
	"testing"

	firmeve2 "github.com/firmeve/firmeve"

	"github.com/firmeve/firmeve/config"
	"github.com/firmeve/firmeve/support/path"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	db := New(config.New(path.RunRelative("../testdata/config")))
	assert.NotPanics(t, func() {
		db.ConnectionDefault()
		db.CloseDefault()
	})
	assert.Equal(t, db.ConnectionDefault(), db.ConnectionDefault())
}

func TestNew_Connection_Error(t *testing.T) {
	config := config.New(path.RunRelative("../testdata/config"))
	config.Item("database").Set("error_connection.addr", "nothing")
	db := New(config)
	assert.Panics(t, func() {
		db.Connection(`error_connection`)
	})
}

func TestNew_Close_Error(t *testing.T) {
	db := New(config.New(path.RunRelative("../testdata/config")))
	assert.NotPanics(t, func() {
		db.ConnectionDefault()
		db.CloseDefault()
		db.CloseDefault()
	})
}

func TestDB_Provider(t *testing.T) {
	firmeve := firmeve2.Instance()
	firmeve.Boot()
	assert.Equal(t, true, firmeve.HasProvider("db"))
	assert.Equal(t, true, firmeve.Has(`db`))
	assert.Equal(t, true, firmeve.Has(`db.connection`))

	provider := firmeve.Resolve(new(Provider)).(*Provider)
	assert.Equal(t, firmeve, provider.Firmeve)
}
