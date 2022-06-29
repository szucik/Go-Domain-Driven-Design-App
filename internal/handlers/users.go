// Package classification of UserKey API
//
// Documentation for UserKey API
//
//	Schemes: http
//	BasePath: /
//	Version: 1.0.0
//
//	Consumes:
//	- application/json
//
//	Produces:
//	- application/json
//
// swagger:meta
package handlers

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/szucik/go-simple-rest-api/internal/data"
	"log"
	"net/http"
)

var (
	ErrHash = errors.New("Problem with hashing your password")

	MsgUserAlreadyExists  = "UserKey already exists with the given email"
	MsgUserNotFound       = "No user account exists with given email. Please sign in first"
	MsgUserCreationFailed = "Unable to create user.Please try again later"
)

type UserKey struct{}

//A list of users returns in the response
// swagger:response usersResponse
type usersResponseWrapper struct {
	//All users in the system
	//in: body
	Body []data.User
}

type GenericResponse struct {
	Status  bool        `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// swagger:response noContent
type usersNoContent struct {
}

//swagger:parameters deleteUser
type usersIDParameterWrapper struct {
	//The id of user delete from database
	//	in:path
	//	required: true
	ID int `json:"id"`
}

//Users is a http.Handler
type Users struct {
	l  *log.Logger
	db *data.Database
}

func NewUsers(l *log.Logger, db *data.Database) *Users {
	return &Users{l, db}
}

// swagger:route GET /users listUsers
// Return a list of users from the database
// responses:
//	200: usersResponse

// GetUsers returns Users from data store
func (u *Users) GetUsers(rw http.ResponseWriter, r *http.Request) {
	u.l.Println("Handle GET Users")
	rw.Header().Add("Content-Type", "application/json")

	users, err := u.db.GetUsers()
	if err != nil {
		fmt.Errorf("%s", err)
	}

	j, _ := json.Marshal(users)

	if err != nil {
		http.Error(rw, "Unable to marshal json", http.StatusInternalServerError)
	}
	rw.Write(j)
}

// swagger:route DELETE /users/{id} users deleteUser
// Return a list of users from the database
// responses:
//	201: noContent

// DeleteUser deletes a user from DB

func (u *Users) DeleteUser(rw http.ResponseWriter, r *http.Request) {}

func (u *Users) UpdateUser(rw http.ResponseWriter, r *http.Request) {}

func (u *Users) Dashboard(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	//TODO Remember about Authorization cookies in client
	rw.Write([]byte("Welcom in dashboard"))
}
