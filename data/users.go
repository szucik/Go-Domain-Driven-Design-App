package data

import (
	"encoding/json"
	"fmt"
	"io"
	"time"
)
var ErrUserNotFound = fmt.Errorf("user not found")
type User struct {
	ID int `json:"id"`
	Name string `json:"name"`
	Surname string `json:"surname"`
	Email string `json:"email"`
	Created string `json:"-"`
	Updated string `json:"-"`
}

type Users []*User

func GetUsers() Users {
	return usersList
}

func (u * User) FromJSON(r io.Reader) error {
	e := json.NewDecoder(r)
	return e.Decode(u)
}

func (u * Users) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(u)
}

func AddUser(u *User) {
	u.ID = getNextID()
	usersList = append(usersList, u)
}

func getNextID() int {
	lu := usersList[len(usersList)-1]
	return lu.ID + 1
}

func UpdateUser(id int, u * User) error {
	_, pos, err := findUser(id)
	if err != nil {
		return err
	}
	u.ID = id
	usersList[pos] = u
	return nil
}


func findUser(id int) (*User, int, error) {
	for i, user := range usersList {
		if user.ID == id{
			return user, i, nil
		}
	}
	return nil, -1, ErrUserNotFound
}

var usersList = []*User{
	&User{
		ID: 1,
		Name:    "Janusz",
		Surname: "Koalski",
		Email:   "janusz@wp.pl",
		Created: time.Now().UTC().String(),
		Updated: time.Now().UTC().String(),

	},
	&User{
		ID: 2,
		Name:    "Tomasz",
		Surname: "Jakut",
		Email:   "tj@gmail.pl",
		Created: time.Now().UTC().String(),
		Updated: time.Now().UTC().String(),
	},
}
