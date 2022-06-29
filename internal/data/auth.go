package data

import (
	"database/sql"
	"errors"
)

var (
	ErrInvalidLoginCred = errors.New("invalid login credentials")
)

type AuthCredentials struct {
	Username string `json:"username" sql:"username"`
	Email    string `json:"email" validate:"required" sql:"email"`
	Password string `json:"password" validate:"required" sql:"password"`
}

func (db *Database) Login(email string) (*AuthCredentials, error) {
	var auth AuthCredentials

	err := db.db.QueryRow("SELECT username, password FROM users where email = ?", email).Scan(&auth.Username, &auth.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrInvalidLoginCred
		}
		return nil, err

	}

	return &auth, nil
}

func (db *Database) AddUser(u *User) (int64, error) {
	stmt, err := db.db.Prepare("INSERT INTO users (username, email, password, tokenhash) VALUES (?, ?, ?, ?)")
	if err != nil {
		panic(err.Error())
	}
	result, err := stmt.Exec(u.Username, u.Email, u.Password, u.TokenHash)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return id, nil
}
