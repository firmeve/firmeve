package main

import (
	"github.com/firmeve/firmeve"
	"github.com/firmeve/firmeve/database"
	"github.com/firmeve/firmeve/kernel/contract"
	// if is pgsql or sqlite
	//_ "github.com/golang-migrate/migrate/database/postgres"
	//_ "github.com/golang-migrate/migrate/database/sqlite3"
)

func main() {
	firmeve.RunDefault(firmeve.WithConfigPath(`../../config.yaml`), firmeve.WithProviders(
		[]contract.Provider{
			new(database.Provider),
		},
	), firmeve.WithCommands(
		[]contract.Command{
			new(database.MigrateCommand),
		},
	))
}
