package data

import (
	"database/sql"
	"encoding/json"
	"github.com/go-playground/validator/v10"
	"io"
)

type Database struct {
	db *sql.DB
}

func NewDatabase(db *sql.DB) *Database {
	return &Database{db}
}

func Validate(i interface{}) error {
	v := validator.New()
	return v.Struct(i)
}

func ToJSON(i interface{}, w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(i)
}

func FromJSON(i interface{}, r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(i)
}

type MySQLError struct {
	Number  uint16
	Message string
}
