package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/szucik/trade-helper/user"
)

type GenericResponse struct {
	Status  bool        `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"dao"`
}

func SignUp(service user.UsersService) func(http.ResponseWriter, *http.Request) {
	return func(rw http.ResponseWriter, r *http.Request) {
		var user user.User
		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusBadRequest)
			return
		}

		id, err := service.SignUp(user)

		if err != nil {
			http.Error(rw, err.Error(), 400)
			return
		}
		rw.WriteHeader(200)
		fmt.Printf("json data: %s\n", id)
		rw.Write([]byte("test"))
		return
	}
}

func GetUsers(service user.UsersService) func(http.ResponseWriter, *http.Request) {
	return func(rw http.ResponseWriter, r *http.Request) {
		users, err := service.GetUsers()
		if err != nil {
			rw.WriteHeader(400)
		}

		jsonData, err := json.Marshal(users)
		if err != nil {
			fmt.Printf("could not marshal json: %s\n", err)
			return
		}
		fmt.Printf("%s", jsonData)

		rw.Write(jsonData)
	}
}
