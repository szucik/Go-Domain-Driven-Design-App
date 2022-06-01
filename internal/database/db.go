package database

import (
	"database/sql"
	"flag"
	"fmt"
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
	return message + fmt.Sprintf(" Can also be configured by env variable %s.", env)
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

func Connect() (*sql.DB, error) {
	c := config()
	db, err = sql.Open(c.driver, c.dbParams)
	if err != nil {
		return nil, err
	}
	return db, nil
}
