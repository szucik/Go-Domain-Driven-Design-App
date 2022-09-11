package dao

import (
	"database/sql"
	"encoding/json"
	"io"

	"github.com/go-playground/validator/v10"
)

type Dao struct {
	db *sql.DB
}

func NewDatabase(db *sql.DB) *Dao {
	return &Dao{db}
}

func (d *Dao) Validate(i interface{}) error {
	return validator.New().Struct(i)
}

func (d *Dao) ToJSON(i interface{}, w io.Writer) error {
	return json.NewEncoder(w).Encode(i)
}

func (d *Dao) FromJSON(i interface{}, r io.Reader) error {
	return json.NewDecoder(r).Decode(i)
}

type MySQLError struct {
	Number  uint16
	Message string
}
