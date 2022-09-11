package dao

import (
	"database/sql"
	"errors"
	"time"
)

var (
	ErrInvalidLoginCred = errors.New("invalid login credentials")
)

type AuthCredentials struct {
	Username string `json:"username" sql:"username"`
	Email    string `json:"email" validate:"required" sql:"email"`
	Password string `json:"password" validate:"required" sql:"password"`
}

func (d *Dao) Login(email string) (*AuthCredentials, error) {
	var auth AuthCredentials

	err := d.db.QueryRow("SELECT username, password FROM users where email = ?", email).Scan(&auth.Username, &auth.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrInvalidLoginCred
		}
		return nil, err

	}

	return &auth, nil
}

func (d *Dao) AddUser(u *User) (int64, error) {
	stmt, err := d.db.Prepare("INSERT INTO users (username, email, password, tokenhash) VALUES (?, ?, ?, ?)")
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

// User represents the user for this application
// swagger:model User
type User struct {
	// the id for this user
	// required: true
	// min: 1
	ID int `json:"id"`
	// the name for this user
	// required: true
	// min length: 3
	Username string `json:"username" validate:"required"`
	// the email address for this user
	// required: true
	// example: user@provider.net
	Email     string    `json:"email" validate:"required" sql:"email"`
	Password  string    `json:"password" validate:"required" sql:"password"`
	TokenHash string    `json:"token_hash" sql:"token_hash"`
	Created   time.Time `json:"created" sql:"created"`
	Updated   time.Time `json:"updated" sql:"updated"`
}

type Users []*User

func (d *Dao) GetUsers() (*Users, error) {
	selDB, err := d.db.Query("SELECT * FROM users ORDER BY id DESC")
	if err != nil {
		panic(err.Error())
	}
	user := &User{}
	users := Users{}
	for selDB.Next() {
		var id int
		var created, updated []uint8
		var email, username, password, tokenhash string
		err = selDB.Scan(&id, &username, &email, &password, &tokenhash, &created, &updated)
		if err != nil {
			panic(err.Error())
		}
		user.ID = id
		user.Username = username
		user.Email = email
		user.Password = password
		users = append(users, user)
	}

	return &users, nil
}

func (d *Dao) GetUser(id int) (*User, error) {
	var user User
	err := d.db.QueryRow("SELECT username, email FROM users where id = ?", id).Scan(&user.Username, &user.Email)
	if err != nil {
		panic(err.Error())
	}
	return &user, nil
}

func UpdateUser(id int, u *User) error {
	//_, pos, err := findUser(id)
	//if err != nil {
	//	return err
	//}
	//u.ID = id
	//usersList[pos] = u
	return nil
}

func DeleteUser(id int) error {
	//_, pos, err := findUser(id)
	//if err != nil {
	//	return err
	//}
	//usersList = append(usersList[0:pos], usersList[pos+1:]...)
	return nil
}
