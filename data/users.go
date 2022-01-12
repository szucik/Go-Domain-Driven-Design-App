package data

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/go-playground/validator/v10"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

var ErrUserNotFound = fmt.Errorf("user not found")

// User represents the user for this application
//
// swagger:model
type User struct {
	// the id for this user
	//
	// required: true
	// min: 1
	ID int `json:"id"`
	// the name for this user
	// required: true
	// min length: 3
	Name    string `json:"name" validate:"required"`
	Surname string `json:"surname"`
	// the email address for this user
	//
	// required: true
	// example: user@provider.net
	Email   string `json:"login" validate:"email"`
	Created string `json:"-"`
	Updated string `json:"-"`
}

type Users []*User

func (u *User) Validate() error {
	v := validator.New()
	return v.Struct(u)
}

func GetUsers() Users {
	return usersList
}

func (u *User) FromJSON(r io.Reader) error {
	e := json.NewDecoder(r)
	return e.Decode(u)
}

func (u *Users) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(u)
}

//func AddUser(u *User) {
//	u.ID = getNextID()
//	usersList = append(usersList, u)
//}

func AddUser(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	stmt, err := db.Prepare("INSERT INTO posts(title) VALUES(?)")
	if err != nil {
		panic(err.Error())
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err.Error())
	}
	keyVal := make(map[string]string)
	json.Unmarshal(body, &keyVal)
	title := keyVal["title"]
	_, err = stmt.Exec(title)
	if err != nil {
		panic(err.Error())
	}
	fmt.Fprintf(w, "New post was created")
}

func getNextID() int {
	lu := usersList[len(usersList)-1]
	return lu.ID + 1
}

func UpdateUser(id int, u *User) error {
	_, pos, err := findUser(id)
	if err != nil {
		return err
	}
	u.ID = id
	usersList[pos] = u
	return nil
}

func DeleteUser(id int) error {
	_, pos, err := findUser(id)
	if err != nil {
		return err
	}
	usersList = append(usersList[0:pos], usersList[pos+1:]...)
	return nil
}

func findUser(id int) (*User, int, error) {
	for i, user := range usersList {
		if user.ID == id {
			return user, i, nil
		}
	}
	return nil, -1, ErrUserNotFound
}

var usersList = []*User{
	&User{
		ID:      1,
		Name:    "Janusz",
		Surname: "Koalski",
		Email:   "janusz@wp.pl",
		Created: time.Now().UTC().String(),
		Updated: time.Now().UTC().String(),
	},
	&User{
		ID:      2,
		Name:    "Tomasz",
		Surname: "Jakut",
		Email:   "tj@gmail.pl",
		Created: time.Now().UTC().String(),
		Updated: time.Now().UTC().String(),
	},
}
