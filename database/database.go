package database

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"

	"github.com/go-playground/validator/v10"
	"github.com/szucik/go-simple-rest-api/app"
)

var db *sql.DB
var err error

func Validate(i interface{}) error {
	return validator.New().Struct(i)
}

func ToJSON(i interface{}, w io.Writer) error {
	return json.NewEncoder(w).Encode(i)
}

func FromJSON(i interface{}, r io.Reader) error {
	return json.NewDecoder(r).Decode(i)
}

func ConnectWithMysqlDb(config app.Configuration) (*sql.DB, error) {
	if err != nil {
		return nil, err
	}
	dataSourceName := fmt.Sprintf("%s:%s@/%s", config.Database.User, config.Database.Password, config.Database.DbName)
	db, err = sql.Open(config.Database.Driver, dataSourceName)
	if err != nil {
		return nil, err
	}
	return db, nil
}
