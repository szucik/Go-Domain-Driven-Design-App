package handlers

import (
	"fmt"
	"log"
	"net/http"
	"time"

	th_data "github.com/szucik/go-simple-rest-api/internal/dao"
)

type Transactions struct {
	l  *log.Logger
	db *th_data.Dao
}

func NewTransactions(l *log.Logger, db *th_data.Dao) *Transactions {
	return &Transactions{l, db}
}

type TransactionKey struct{}

func (t *Transactions) AddTransaction(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	rc := r.Context().Value(TransactionKey{}).(th_data.Transaction)

	transaction := &th_data.Transaction{

		UserId:      rc.UserId,
		PortfolioId: rc.PortfolioId,
		Symbol:      rc.Symbol,
		Quantity:    rc.Quantity,
		Price:       rc.Price,
		Created:     time.Time{},
	}

	_, err := t.db.AddTransaction(transaction)
	if err != nil {
		message := fmt.Sprintf("Error message: %v", err)
		t.l.Print(message)
		t.db.ToJSON(&GenericResponse{Status: false, Message: err.Error()}, rw)
		return
	}
	t.l.Print("TransactionKey created successfully")

	t.db.ToJSON(&GenericResponse{Status: true, Message: "Transaction created successfully"}, rw)
}
