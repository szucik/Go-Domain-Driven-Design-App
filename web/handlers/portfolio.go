package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/szucik/trade-helper/user"
)

func AddPortfolio(service user.UsersService) func(http.ResponseWriter, *http.Request) {
	return func(rw http.ResponseWriter, r *http.Request) {
		var p user.PortfolioIn
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

func AddTransaction(service user.UsersService) func(http.ResponseWriter, *http.Request) {
	return func(rw http.ResponseWriter, r *http.Request) {
		var p user.TransactionIn
		vars := mux.Vars(r)

		err := json.NewDecoder(r.Body).Decode(&p)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusBadRequest)
			return
		}

		p.UserName = vars["username"]
		p.PortfolioName = vars["name"]

		result, err := service.AddTransaction(p)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusBadRequest)
			return
		}

		writeSuccessMessage(rw, []byte(result))
	}
}
