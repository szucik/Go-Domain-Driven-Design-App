package database

import (
	"database/sql"
	"flag"
	"fmt"
	"github.com/szucik/go-simple-rest-api/internal/configuration"
)

var db *sql.DB
var err error

type dbConfig struct {
	driver, dbParams string
}

const (
	DbDriver   = "DB_DRIVER"
	DbUser     = "DB_USER"
	DbPassword = "DB_PASS"
	DbName     = "DB_NAME"
)

func appendEnvMsg(message, env string) string {
	return message + fmt.Sprintf("Can also be configured by env variable %s.", env)
}

func config() *dbConfig {
	driver := flag.String("driver", "", appendEnvMsg("required: database driver", DbDriver))
	user := flag.String("user", "", appendEnvMsg("required: database user", DbUser))
	password := flag.String(
		"password",
		"",
		appendEnvMsg("required: database password", DbPassword))
	dbName := flag.String("dbName", "", appendEnvMsg("required: database name", DbName))

	flag.Parse()

	dbParameters := fmt.Sprintf("%s:%s@/%s", *user, *password, *dbName)
	return &dbConfig{
		*driver, dbParameters,
	}
}

func Connection() (*sql.DB, error) {
	confEnv := config()

	conf, err := configuration.New()
	if err != nil {
		return nil, err
	}
	db, err = sql.Open(conf.DBDriver, confEnv.dbParams)
	if err != nil {
		return nil, err
	}
	return db, nil
}
