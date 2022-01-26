// Package classification of User API
//
// Documentation for User API
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
	"context"
	"encoding/json"
	"errors"
	"fmt"
	datastruct "github.com/szucik/go-simple-rest-api/internal/data"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
)

var ErrHash = errors.New("Problem with hashing your password")

//A list of users returns in the response
// swagger:response usersResponse
type usersResponseWrapper struct {
	//All users in the system
	//in: body
	Body []datastruct.User
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
	db *datastruct.Database
}
type KeyUser struct{}

//NewUsers creates a users handler with the given logger
func NewUsers(l *log.Logger, db *datastruct.Database) *Users {
	return &Users{l, db}
}

func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(hash), err
}

func (u *Users) AddUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	rc := r.Context().Value(KeyUser{}).(datastruct.User)

	hash, err := HashPassword(rc.Password)
	if err != nil {
		fmt.Errorf("%s", ErrHash)
	}

	user := &datastruct.User{
		Username:  rc.Username,
		Email:     rc.Email,
		Password:  hash,
		TokenHash: rc.TokenHash,
	}

	id, err := u.db.AddUser(user)
	if err != nil {
		message := fmt.Sprintf("Error message: %v", err)
		jM, _ := json.Marshal(message)
		w.Write(jM)

	}
	fmt.Print(id)
}

// swagger:route GET /users users listUsers
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

func (u *Users) DeleteUser(rw http.ResponseWriter, r *http.Request) {
	//vars := mux.Vars(r)
	//id, err := strconv.Atoi(vars["id"])
	//if err != nil {
	//	http.Error(rw, "Unable to convert id", http.StatusBadRequest)
	//	return
	//}
	//err = data.DeleteUser(id)
	//if err == data.ErrUserNotFound {
	//	http.Error(rw, "User not found", http.StatusNotFound)
	//	return
	//}
}

func (u *Users) UpdateUser(rw http.ResponseWriter, r *http.Request) {
	//vars := mux.Vars(r)
	//id, err := strconv.Atoi(vars["id"])
	//if err != nil {
	//	http.Error(rw, "Unable to convert id", http.StatusBadRequest)
	//	return
	//}
	//
	//user := r.Context().Value(KeyUser{}).(data.User)
	//
	//err = data.UpdateUser(id, &user)
	//if err == data.ErrUserNotFound {
	//	http.Error(rw, "User not found", http.StatusNotFound)
	//	return
	//}
	//if err != nil {
	//	http.Error(rw, "User not found", http.StatusInternalServerError)
	//	return
	//}
}

func (u *Users) MiddlewareUserValid(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		usr := &datastruct.User{}
		err := usr.FromJSON(r.Body)
		if err != nil {
			u.l.Println("[ERROR] deserializing user", err)
			http.Error(rw, "Unable to unmarshal json", http.StatusBadRequest)
			return
		}
		err = usr.Validate()
		if err != nil {
			u.l.Println("[ERROR] validate user", err)
			http.Error(rw, "Unable validate user", http.StatusBadRequest)
			return
		}

		ctx := context.WithValue(r.Context(), KeyUser{}, *usr)

		r = r.WithContext(ctx)
		next.ServeHTTP(rw, r)
	})
}
