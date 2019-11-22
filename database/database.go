package database

import (
	"strings"

	"github.com/firmeve/firmeve/config"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mssql"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

type (
	DB struct {
		config      config.Configurator
		db          *gorm.DB
		connections dbConnection
	}

	dbConnection map[string]*gorm.DB
)

func New(config config.Configurator) *DB {
	return &DB{
		config:      config,
		connections: make(dbConnection, 0),
	}
}

func (d *DB) ConnectionDefault() *gorm.DB {
	return d.Connection(d.config.GetString(`default`))
}

func (d *DB) Connection(driver string) *gorm.DB {
	if connection, ok := d.connections[driver]; ok {
		return connection
	}

	config := d.config.GetString(strings.Join([]string{`connections`, driver, `addr`}, `.`))
	db, err := gorm.Open(driver, config)

	if err != nil {
		panic(err)
	}

	d.connections[driver] = db

	return db
}

func (d *DB) CloseDefault() {
	d.Close(d.config.GetString(`default`))
}

func (d *DB) Close(driver string) {
	if connection, ok := d.connections[driver]; ok {
		err := connection.Close()
		if err != nil {
			panic(err)
		}
		delete(d.connections, driver)
	}
}
