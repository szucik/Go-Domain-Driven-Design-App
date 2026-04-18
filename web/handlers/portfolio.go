package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/szucik/trade-helper/apperrors"
	"github.com/szucik/trade-helper/user"
)

func AddPortfolio(service user.UsersService) func(http.ResponseWriter, *http.Request) {
	return func(rw http.ResponseWriter, r *http.Request) {
		var p user.PortfolioIn

		if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
			writeError(rw, apperrors.Error("invalid request body", "BadRequest", http.StatusBadRequest))
			return
		}

		if p.Name == "" {
			writeError(rw, apperrors.Error("portfolio name is required", "ValidationError", http.StatusBadRequest))
			return
		}

		p.UserName = mux.Vars(r)["username"]
		name, err := service.AddPortfolio(r.Context(), p)
		if err != nil {
			writeError(rw, err)
			return
		}

		writeJSON(rw, http.StatusCreated, map[string]string{"name": name})
	}
}

func AddTransaction(service user.UsersService) func(http.ResponseWriter, *http.Request) {
	return func(rw http.ResponseWriter, r *http.Request) {
		var p user.TransactionIn
		vars := mux.Vars(r)

		if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
			writeError(rw, apperrors.Error("invalid request body", "BadRequest", http.StatusBadRequest))
			return
		}

		if p.Symbol == "" || p.Amount == "" || p.Quantity == "" {
			writeError(rw, apperrors.Error("symbol, amount and quantity are required", "ValidationError", http.StatusBadRequest))
			return
		}

		p.UserName = vars["username"]
		p.PortfolioName = vars["name"]

		result, err := service.AddTransaction(r.Context(), p)
		if err != nil {
			writeError(rw, err)
			return
		}

		writeJSON(rw, http.StatusCreated, map[string]string{"id": result})
	}
}

func GetTransactions(service user.UsersService) func(http.ResponseWriter, *http.Request) {
	return func(rw http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		username := vars["username"]
		portfolioName := vars["name"]

		out, err := service.GetTransactions(r.Context(), username, portfolioName)
		if err != nil {
			writeError(rw, err)
			return
		}

		writeJSON(rw, http.StatusOK, out)
	}
}
