package handlers

import (
	"context"
	"github.com/gorilla/mux"
	"github.com/zlyjoker102/simple-rest-api/data"
	"log"
	"net/http"
	"strconv"
)

type Users struct {
	l *log.Logger
}
type KeyUser struct{}

func NewUser(l *log.Logger) *Users {
	return &Users{l}
}

func (u *Users) AddUser(rw http.ResponseWriter, r *http.Request) {
	u.l.Println("Handle POST Users")
	user := r.Context().Value(KeyUser{}).(data.User)
	data.AddUser(&user)
}

func (u *Users) DeleteUser(rw http.ResponseWriter, r *http.Request)  {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(rw, "Unable to convert id", http.StatusBadRequest)
		return
	}
	err = data.DeleteUser(id)
	if err == data.ErrUserNotFound {
		http.Error(rw, "User not found", http.StatusNotFound)
		return
	}
}

func (u *Users) UpdateUser(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(rw, "Unable to convert id", http.StatusBadRequest)
		return
	}

	user := r.Context().Value(KeyUser{}).(data.User)

	err = data.UpdateUser(id, &user)
	if err == data.ErrUserNotFound {
		http.Error(rw, "User not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(rw, "User not found", http.StatusInternalServerError)
		return
	}
}

func (u *Users) GetUsers(rw http.ResponseWriter, r *http.Request) {
	u.l.Println("Handle GET Users")
	lu := data.GetUsers()
	err := lu.ToJSON(rw)
	if err != nil {
		http.Error(rw, "Unable to marshal json", http.StatusInternalServerError)
	}
}

func (u *Users) MiddlewareUserValid(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		usr := data.User{}
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

		ctx := context.WithValue(r.Context(), KeyUser{}, usr)

		r = r.WithContext(ctx)
		next.ServeHTTP(rw, r)
	})
}
