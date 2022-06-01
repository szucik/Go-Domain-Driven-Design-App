package data

import (
	"database/sql"
	"errors"
)

var (
	ErrInvalidLoginCred = errors.New("invalid login credentials")
)

type Auth struct {
	Username string `json:"username" sql:"username"`
	Email    string `json:"email" validate:"required" sql:"email"`
	Password string `json:"password" validate:"required" sql:"password"`
}

func (db *Database) Login(email string) (*Auth, error) {
	var auth Auth

	err := db.db.QueryRow("SELECT username, password FROM users where email = ?", email).Scan(&auth.Username, &auth.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrInvalidLoginCred
		}
		return nil, err

	}

	return &auth, nil
}
