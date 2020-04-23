package database

import (
	"errors"
	"github.com/fatih/color"
	config2 "github.com/firmeve/firmeve/config"
	"github.com/firmeve/firmeve/kernel/contract"
	"github.com/firmeve/firmeve/support/path"
	"github.com/firmeve/firmeve/support/slices"
	"github.com/firmeve/firmeve/support/strings"
	"github.com/golang-migrate/migrate"
	_ "github.com/golang-migrate/migrate/database/mysql"
	_ "github.com/golang-migrate/migrate/source/file"
	"github.com/spf13/cobra"
	"os"
	path3 "path"
	"path/filepath"
	"strconv"
	"time"
)

type (
	MigrateCommand struct {
	}
)

var (
	actions = []string{
		`create`, `up`, `down`, `step`, `rollback`, `force`, `version`, //`drop`,
	}
)

func (m *MigrateCommand) CobraCmd() *cobra.Command {
	command := new(cobra.Command)
	command.Use = "migrate"
	command.Long = `Usage: migrate OPTIONS COMMAND [arg...]
       migrate [ -version | -help ]
Options:
  -driver       Run migrations against this driver (driver://url)
  -path         Shorthand for -path=path (Only create command)
Commands:
  create        NAME Create a set of timestamped up/down migrations titled NAME
  step N        Migrate to version V
  rollback N    Migrate rollback to version V
  up            Apply all up migrations
  down          Apply all down migrations
  drop          Drop everything inside database
  force V       Set version V but don't run migration (ignores dirty state)
  version       Print current migration version
`
	command.Short = "Migrate files"
	command.Args = func(cmd *cobra.Command, args []string) error {
		if !slices.InString(actions, args[0]) {
			return errors.New("the first parameter must be [" + strings.Join(`|`, actions...) + "] range")
		}

		return nil
	}

	command.Flags().StringP("driver", "D", "mysql", "migration driver example: mysql")
	command.Flags().StringP("path", "", "", "migration file create path")

	return command
}

func (m *MigrateCommand) Run(root contract.BaseCommand, cmd *cobra.Command, args []string) {
	var (
		config = root.Resolve(`config`).(*config2.Config)
		logger = root.Resolve(`logger`).(contract.Loggable)
		path2  string
		driver string
	)

	driver = cmd.Flag(`driver`).Value.String()
	path2, err := m.targetDir(
		cmd.Flag(`path`).Value.String(),
		config.Item(`database`).GetString(`migration.path`),
	)

	if err != nil {
		logger.Error("migration path create error", "error", err)
		return
	}

	// action
	action := args[0]

	if action == `create` {
		timestamp := time.Now().Format("20060102150405")
		//timestamp := strconv.FormatInt(time.Now().Unix(), 10)
		upFile := path3.Join(path2, timestamp+"_"+args[1]+".up.sql")
		downFile := path3.Join(path2, timestamp+"_"+args[1]+".down.sql")

		f1, err := os.Create(upFile)
		if err != nil {
			logger.Error("create up file error", "error", err)
			return
		}
		defer f1.Close()

		f2, err := os.Create(downFile)

		if err != nil {
			logger.Error("create down file error", "error", err)
			return
		}
		defer f2.Close()
	} else {
		dbConfig := root.Resolve(`config`).(*config2.Config).Item(`database`)
		connection := dbConfig.GetString(`connections.` + driver + `.addr`)
		m2, err := migrate.New(
			`file://`+path2,
			driver+"://"+connection,
		)
		if err != nil {
			logger.Error("migration connection error", "error", err)
		}
		/*else if action == `drop` {
			err = m2.Drop()
		}*/
		if action == `up` {
			err = m2.Up()
		} else if action == `down` {
			err = m2.Down()
		} else if action == `step` {
			step, _ := strconv.Atoi(args[1])
			err = m2.Steps(step)
		} else if action == `rollback` { // rollback
			step, _ := strconv.Atoi(args[1])
			err = m2.Steps(-step)
		} else if action == `force` {
			version, _ := strconv.Atoi(args[1])
			err = m2.Force(version)
		} else { //version
			var version uint
			version, _, err = m2.Version()
			greenColor := color.New(color.FgRed)
			greenColor.Printf("current version %d\n", version)
		}
		if err != nil {
			logger.Error("migration error", "error", err)
			return
		}
	}

	greenColor := color.New(color.FgGreen)
	greenColor.Println(action + " migration successfully")
}

func (m *MigrateCommand) targetDir(currentPath, defaultPath string) (string, error) {
	if currentPath == `` {
		currentPath = defaultPath
	}

	if !path.Exists(currentPath) {
		err := os.MkdirAll(currentPath, 0755)
		if err != nil {
			return ``, nil
		}
		currentPath, _ = filepath.Abs(currentPath)
	}

	return currentPath, nil
}
