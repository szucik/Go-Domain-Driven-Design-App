package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/szucik/trade-helper/user"
)

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

		writeSuccesMessage(rw, []byte(username))
	}
}

func GetUsers(service user.UsersService) func(http.ResponseWriter, *http.Request) {
	return func(rw http.ResponseWriter, r *http.Request) {
		users, err := service.GetUsers()
		if err != nil {
			http.Error(rw, err.Error(), 400)
			return
		}

		result, err := json.Marshal(users)
		if err != nil {
			fmt.Printf("could not marshal json: %s\n", err)
			return
		}

		writeSuccesMessage(rw, result)
	}
}

func GetUser(service user.UsersService) func(http.ResponseWriter, *http.Request) {
	return func(rw http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		user, err := service.GetUser(vars["username"])
		if err != nil {
			fmt.Errorf("%s", err.Error())
			return
		}

		rw.WriteHeader(http.StatusOK)
		result, err := json.Marshal(user)
		if err != nil {
			fmt.Printf("could not marshal json: %s\n", err)
			return
		}

		writeSuccesMessage(rw, result)
	}
}

func writeSuccesMessage(rw http.ResponseWriter, result []byte) {
	rw.WriteHeader(http.StatusOK)
	rw.Write(result)
}
