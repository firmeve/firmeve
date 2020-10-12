package database

import (
	"database/sql"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"time"
)

type (
	DB struct {
		config      *Configuration
		connections map[string]*connection
	}

	connection struct {
		DB    *gorm.DB
		SqlDB *sql.DB
	}

	Configuration struct {
		Default     string `json:"default" yaml:"default"`
		Connections map[string]struct {
			Host     string `json:"host" yaml:"host"`
			Database string `json:"database" yaml:"database"`
			Username string `json:"username" yaml:"username"`
			Password string `json:"password" yaml:"password"`
			Charset  string `json:"charset" yaml:"charset"`
		} `json:"connections" yaml:"connections"`
		Pool struct {
			MaxIdle       int `json:"max_idle" yaml:"max_idle"`
			MaxConnection int `json:"max_connection" yaml:"max_connection"`
			MaxLifetime   int `json:"max_lifetime" yaml:"max_lifetime"`
		}
		Debug     bool `json:"debug" yaml:"debug"`
		Migration struct {
			Path string `json:"path" yaml:"path"`
		} `json:"migration" yaml:"migration"`
	}
)

func New(config *Configuration) *DB {
	return &DB{
		config:      config,
		connections: make(map[string]*connection, 0),
	}
}

func (db *DB) ConnectionDB(conn string) *gorm.DB {
	return db.Connection(conn).DB
}

func (db *DB) ConnectionNewDB(conn string) *gorm.DB {
	return db.Connection(conn).DB.Session(&gorm.Session{WithConditions: false})
}

// 数据库连接
func (db *DB) Connection(conn string) *connection {
	if v, ok := db.connections[conn]; ok {
		return v
	}

	var (
		connectionConfig = db.config.Connections[conn]
		dsn              = fmt.Sprintf(`%s:%s@(%s)/%s?charset=%s&parseTime=True&loc=Local`, connectionConfig.Username, connectionConfig.Password, connectionConfig.Host, connectionConfig.Database, connectionConfig.Charset)
	)

	db2, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		SkipDefaultTransaction: true,
		//Logger: logger.New(infrastructure.Logger, logger.Config{
		//	SlowThreshold: 100 * time.Millisecond,
		//	Colorful:      true,
		//	LogLevel:      logger.Info,
		//}),
	})

	if err != nil {
		panic(err)
	}

	// connection poll
	sqlDB, err := db2.DB()
	if err != nil {
		panic(err)
	}

	sqlDB.SetMaxIdleConns(db.config.Pool.MaxIdle)
	sqlDB.SetMaxOpenConns(db.config.Pool.MaxConnection)
	sqlDB.SetConnMaxLifetime(time.Duration(db.config.Pool.MaxLifetime))

	currentConnection := &connection{
		DB:    db2,
		SqlDB: sqlDB,
	}
	db.connections[conn] = currentConnection
	return currentConnection
}

func (db *DB) Close(connection string) {
	db.connections[connection].SqlDB.Close()
}
