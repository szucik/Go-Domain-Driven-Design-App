package handlers

import (
	"encoding/json"
	"github.com/szucik/go-simple-rest-api/user"
	"net/http"
)

type GenericResponse struct {
	Status  bool        `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"dao"`
}

func SignUp(service user.UsersService) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var user user.User
		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		err = service.SignUp(user)
		if err != nil {
			w.WriteHeader(400)
		}
	}
}

func GetUsers(service user.UsersService) func(http.ResponseWriter, *http.Request) {
	return func(rw http.ResponseWriter, r *http.Request) {
		users, err := service.GetUsers()
		if err != nil {
			rw.WriteHeader(400)
		}
		result, err := json.Marshal(users)

		if err != nil {
			http.Error(rw, "Unable to marshal json", http.StatusInternalServerError)
		}
		rw.Write(result)
	}
}
