package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"

	"github.com/szucik/trade-helper/apperrors"
	"github.com/szucik/trade-helper/user"
)

func SignUp(service user.UsersService) func(http.ResponseWriter, *http.Request) {
	return func(rw http.ResponseWriter, r *http.Request) {
		var u user.User
		if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
			writeError(rw, apperrors.Error(err.Error(), "BadRequest", http.StatusBadRequest))
			return
		}

		if u.Username == "" || u.Email == "" || u.Password == "" {
			writeError(rw, apperrors.Error("username, email and password are required", "ValidationError", http.StatusBadRequest))
			return
		}

		username, err := service.SignUp(r.Context(), u)
		if err != nil {
			writeError(rw, err)
			return
		}

		writeSuccessMessage(rw, []byte(username))
	}
}

func SignIn(service user.UsersService, store *sessions.CookieStore) func(http.ResponseWriter, *http.Request) {
	return func(rw http.ResponseWriter, r *http.Request) {
		var credentials user.AuthCredentials
		if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
			writeError(rw, apperrors.Error(err.Error(), "BadRequest", http.StatusBadRequest))
			return
		}

		if credentials.Email == "" || credentials.Password == "" {
			writeError(rw, apperrors.Error("email and password are required", "ValidationError", http.StatusBadRequest))
			return
		}

		username, err := service.SignIn(r.Context(), credentials)
		if err != nil {
			writeError(rw, err)
			return
		}

		session, _ := store.Get(r, "X-Auth")
		session.Options = &sessions.Options{
			Path:     "/",
			MaxAge:   86400 * 7,
			HttpOnly: true,
		}
		session.Values["username"] = username

		if err = session.Save(r, rw); err != nil {
			writeError(rw, apperrors.Error(err.Error(), "SessionError", http.StatusInternalServerError))
			return
		}

		writeSuccessMessage(rw, []byte(username))
	}
}

func GetUsers(service user.UsersService) func(http.ResponseWriter, *http.Request) {
	return func(rw http.ResponseWriter, r *http.Request) {
		p := user.PaginationIn{
			Page:  queryInt(r, "page", 0),
			Limit: queryInt(r, "limit", 0),
		}

		out, err := service.GetUsers(r.Context(), p)
		if err != nil {
			writeError(rw, err)
			return
		}

		result, err := json.Marshal(out)
		if err != nil {
			writeError(rw, apperrors.Error(err.Error(), "MarshalError", http.StatusInternalServerError))
			return
		}

		writeSuccessMessage(rw, result)
	}
}

func GetUser(service user.UsersService) func(http.ResponseWriter, *http.Request) {
	return func(rw http.ResponseWriter, r *http.Request) {
		u, err := service.GetUserByName(r.Context(), mux.Vars(r)["username"])
		if err != nil {
			writeError(rw, err)
			return
		}

		result, err := json.Marshal(u)
		if err != nil {
			writeError(rw, apperrors.Error(err.Error(), "MarshalError", http.StatusInternalServerError))
			return
		}

		writeSuccessMessage(rw, result)
	}
}

func UpdateUser(service user.UsersService) func(http.ResponseWriter, *http.Request) {
	return func(rw http.ResponseWriter, r *http.Request) {
		var in user.UpdateUserIn
		if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
			writeError(rw, apperrors.Error(err.Error(), "BadRequest", http.StatusBadRequest))
			return
		}

		if in.Username == "" && in.Email == "" && in.Password == "" {
			writeError(rw, apperrors.Error("at least one field is required", "ValidationError", http.StatusBadRequest))
			return
		}

		username := mux.Vars(r)["username"]
		updated, err := service.UpdateUser(r.Context(), username, in)
		if err != nil {
			writeError(rw, err)
			return
		}

		writeSuccessMessage(rw, []byte(updated))
	}
}

func DeleteUser(service user.UsersService) func(http.ResponseWriter, *http.Request) {
	return func(rw http.ResponseWriter, r *http.Request) {
		username := mux.Vars(r)["username"]
		if err := service.DeleteUser(r.Context(), username); err != nil {
			writeError(rw, err)
			return
		}

		rw.WriteHeader(http.StatusNoContent)
	}
}

func writeError(rw http.ResponseWriter, err error) {
	var appErr apperrors.ErrorResponse
	if errors.As(err, &appErr) {
		appErr.JSONError(rw)
		return
	}
	http.Error(rw, err.Error(), http.StatusInternalServerError)
}

func writeSuccessMessage(rw http.ResponseWriter, result []byte) {
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	rw.Write(result)
}

func queryInt(r *http.Request, key string, def int) int {
	s := r.URL.Query().Get(key)
	if s == "" {
		return def
	}
	v, err := strconv.Atoi(s)
	if err != nil {
		return def
	}
	return v
}
