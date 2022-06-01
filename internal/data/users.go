package data

import (
	"encoding/json"
	"fmt"
	"github.com/go-playground/validator/v10"
	"io"
	"time"
)

var (
	ErrUserNotFound     = fmt.Errorf("user not found")
	ErrCannotCreateUser = fmt.Errorf("user cannot be created")
)

// User represents the user for this application
//
// swagger:model
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
	TokenHash string    `json:"tokenhash" sql:"tokenhash"`
	Created   time.Time `json:"created" sql:"created"`
	Updated   time.Time `json:"updated" sql:"updated"`
}

type Users []*User

func (u *User) Validate() error {
	v := validator.New()
	return v.Struct(u)
}

func (u *User) FromJSON(r io.Reader) error {
	e := json.NewDecoder(r)
	return e.Decode(u)
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

func (db *Database) GetUsers() (*Users, error) {
	selDB, err := db.db.Query("SELECT * FROM users ORDER BY id DESC")
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

//func (db *Database) GetUser(id int) (*User, error) {
//	//selDB, err := db.db.Query("SELECT * FROM Users ORDER BY id DESC")
//	//if err != nil {
//	//	panic(err.Error())
//	//}
//	return nil, nil
//}

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
