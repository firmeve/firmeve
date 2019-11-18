package database

import (
	"github.com/firmeve/firmeve/config"
	"github.com/firmeve/firmeve/support/path"
	"github.com/stretchr/testify/assert"
	"testing"

	firmeve2 "github.com/firmeve/firmeve"
)

//func TestNew(t *testing.T) {
//	db := New(config.New(path.RunRelative("../testdata/config")))
//	assert.NotPanics(t, func() {
//		db.ConnectionDefault()
//		db.CloseDefault()
//	})
//	assert.Equal(t, db.ConnectionDefault(), db.ConnectionDefault())
//}
//
//func TestNew_Connection_Error(t *testing.T) {
//	config := config.New(path.RunRelative("../testdata/config"))
//	config.Item("database").Set("error_connection.addr", "nothing")
//	db := New(config)
//	assert.Panics(t, func() {
//		db.Connection(`error_connection`)
//	})
//}
//
//func TestNew_Close_Error(t *testing.T) {
//	db := New(config.New(path.RunRelative("../testdata/config")))
//	assert.NotPanics(t, func() {
//		db.ConnectionDefault()
//		db.CloseDefault()
//		db.CloseDefault()
//	})
//}

func TestDB_Provider(t *testing.T) {
	firmeve := firmeve2.New()
	firmeve.Bind(`config`, config.New(path.RunRelative("../testdata/config")))
	firmeve.Register(firmeve.Make(new(Provider)).(firmeve2.Provider))

	//z := &Provider{
	//	firmeve2.BaseFirmeve{
	//		Firmeve: firmeve,
	//	},
	//}
	//z2 := cache.Provider{
	//	firmeve2.BaseFirmeve{
	//		Firmeve: firmeve,
	//	},
	//}
	//fmt.Printf("%#v\n", z)
	//fmt.Printf("%#v", z2)

	firmeve.Boot()
	assert.Equal(t, true, firmeve.HasProvider("db"))
	assert.Equal(t, true, firmeve.Has(`db`))
	assert.Equal(t, true, firmeve.Has(`db.connection`))

	provider := firmeve.Make(new(Provider)).(*Provider)
	assert.Equal(t, firmeve, provider.Firmeve)
}
