package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/szucik/trade-helper/user"
)

func SignUp(ctx context.Context, service user.UsersService) func(http.ResponseWriter, *http.Request) {
	return func(rw http.ResponseWriter, r *http.Request) {
		var user user.User

		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusBadRequest)
			return
		}

		username, err := service.SignUp(ctx, user)

		if err != nil {
			//TODO - handle errors, or logger
			http.Error(rw, err.Error(), 400)
			return
		}

		writeSuccessMessage(rw, []byte(username))
	}
}

func SignIn(ctx context.Context, service user.UsersService) func(http.ResponseWriter, *http.Request) {
	return func(rw http.ResponseWriter, r *http.Request) {
		var credentials user.AuthCredentials

		err := json.NewDecoder(r.Body).Decode(&credentials)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusBadRequest)
			return
		}

		err = service.SignIn(ctx, credentials)
		if err != nil {
			http.Error(rw, err.Error(), 400)
			return
		}

		//writeSuccessMessage(rw, []byte(username))
	}
}

func GetUsers(ctx context.Context, service user.UsersService) func(http.ResponseWriter, *http.Request) {
	return func(rw http.ResponseWriter, r *http.Request) {
		users, err := service.GetUsers(ctx)
		if err != nil {
			http.Error(rw, err.Error(), 400)
			return
		}

		result, err := json.Marshal(users)
		if err != nil {
			fmt.Printf("could not marshal json: %s\n", err)
			return
		}

		writeSuccessMessage(rw, result)
	}
}

func GetUser(ctx context.Context, service user.UsersService) func(http.ResponseWriter, *http.Request) {
	return func(rw http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		user, err := service.GetUserByEmail(ctx, vars["email"])
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

		writeSuccessMessage(rw, result)
	}
}

func writeSuccessMessage(rw http.ResponseWriter, result []byte) {
	rw.WriteHeader(http.StatusOK)
	rw.Write(result)
}
