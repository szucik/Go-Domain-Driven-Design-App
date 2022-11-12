package transaction

import (
	"log"
	"net/http"
)

type Database interface {
}

type Transactions struct {
	Logger       *log.Logger
	Database     Database
	NewAggregate func(Transaction) (Aggregate, error)
}

type TransactionKey struct{}

func (t *Transactions) AddTransaction(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	//rc := r.Context().Value(TransactionKey{}).(Transaction)

	//transaction := &Transaction{
	//
	//	UserId:      rc.UserId,
	//	PortfolioId: rc.PortfolioId,
	//	Symbol:      rc.Symbol,
	//	Quantity:    rc.Quantity,
	//	Price:       rc.Price,
	//	Created:     time.Time{},
	//}

	//_, err := t.db.AddTransaction(transaction)
	//if err != nil {
	//	message := fmt.Sprintf("Error message: %v", err)
	//	t.l.Print(message)
	//	t.db.ToJSON(&GenericResponse{Status: false, Message: err.Error()}, rw)
	//	return
	//}
	//t.l.Print("TransactionKey created successfully")
	//
	//t.db.ToJSON(&GenericResponse{Status: true, Message: "Transaction created successfully"}, rw)
}
