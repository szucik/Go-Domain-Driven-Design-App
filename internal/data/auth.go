package data

import (
	"database/sql"
	"errors"
)

type Auth struct {
	// the email address for this user
	// required: true
	// example: user@provider.net
	Email    string `json:"email" validate:"required" sql:"email"`
	Password string `json:"password" validate:"required" sql:"password"`
}

func (db *Database) Login(email string) (*Auth, error) {
	var auth Auth
	// Execute the query
	err := db.db.QueryRow("SELECT email, password FROM users where email = ?", email).Scan(&auth.Email, &auth.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("Invalid user email.\r\n")
		}

		return nil, err
	}
	//if err != nil {
	//	panic(err.Error()) // proper error handling instead of panic in your app
	//}
	return &auth, nil
}
