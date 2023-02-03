package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/szucik/trade-helper/user"
)

type ResponseMessage struct {
	code    int
	message string
}

func SignUp(service user.UsersService) func(http.ResponseWriter, *http.Request) {
	return func(rw http.ResponseWriter, r *http.Request) {
		var user user.User

		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusBadRequest)
			return
		}

		username, err := service.SignUp(user)

		if err != nil {
			http.Error(rw, err.Error(), 400)
			return
		}

		rw.WriteHeader(http.StatusOK)
		rw.Write([]byte(username))
		return
	}
}

func GetUsers(service user.UsersService) func(http.ResponseWriter, *http.Request) {
	return func(rw http.ResponseWriter, r *http.Request) {
		users, err := service.GetUsers()
		if err != nil {
			rw.WriteHeader(400)
			return
		}

		result, err := json.Marshal(users)
		if err != nil {
			fmt.Printf("could not marshal json: %s\n", err)
			return
		}

		rw.WriteHeader(http.StatusOK)
		rw.Write(result)
		return
	}
}

func GetUser(service user.UsersService) func(http.ResponseWriter, *http.Request) {
	return func(rw http.ResponseWriter, r *http.Request) {

		vars := mux.Vars(r)
		user, err := service.GetUser(vars["username"])
		if err != nil {
			http.Error(rw, err.Error(), 400)
			return
		}

		rw.WriteHeader(http.StatusOK)
		result, err := json.Marshal(user)
		if err != nil {
			fmt.Printf("could not marshal json: %s\n", err)
			return
		}

		rw.WriteHeader(http.StatusOK)
		rw.Write(result)
		return
	}
}
