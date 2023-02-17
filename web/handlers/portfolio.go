package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/szucik/trade-helper/user"
)

func AddPortfolio(service user.UsersService) func(http.ResponseWriter, *http.Request) {
	return func(rw http.ResponseWriter, r *http.Request) {
		var p user.AddPortfolioIn
		username := mux.Vars(r)["username"]

		err := json.NewDecoder(r.Body).Decode(&p)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusBadRequest)
			return
		}

		p.UserName = username
		name, err := service.AddPortfolio(p)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusBadRequest)
			return
		}

		writeSuccessMessage(rw, []byte(name))
	}
}
