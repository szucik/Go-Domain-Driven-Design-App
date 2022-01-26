package data

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/go-sql-driver/mysql"
	"io"
)

type Database struct {
	db *sql.DB
}

func NewDatabase(db *sql.DB) *Database {
	return &Database{db}
}

// ToJSON serializes the given interface into a string based JSON format
func ToJSON(i interface{}, w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(i)
}

// FromJSON deserializes the object from JSON string
// given in the io.Reader to the given interface
func FromJSON(i interface{}, r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(i)
}

type MySQLError struct {
	Number  uint16
	Message string
}

func mysqlErrorMessage(err error) error {
	errorNumber := err.(*mysql.MySQLError).Number
	switch errorNumber {
	case 1062:
		return fmt.Errorf("Duplicate entry")
	default:
		return fmt.Errorf("%v", err)
	}
}
