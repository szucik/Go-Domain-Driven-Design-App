package user

import (
	"github.com/google/uuid"
	"github.com/szucik/go-simple-rest-api/portfolio"
	"github.com/szucik/go-simple-rest-api/transaction"
)

type Aggregate struct {
	user        *User
	transaction *transaction.Transaction
	portfolio   *portfolio.Portfolio
}

func (a *Aggregate) User() *User {
	return a.user
}

func (a *Aggregate) Transaction() *transaction.Transaction {
	return a.transaction
}

func (a *Aggregate) Portfolio() *portfolio.Portfolio {
	return a.portfolio
}

func (a *Aggregate) GetID() uuid.UUID {
	return a.User().ID
}
