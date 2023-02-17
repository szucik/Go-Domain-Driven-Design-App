package user

import (
	"github.com/szucik/trade-helper/portfolio"
	"github.com/szucik/trade-helper/transaction"
)

type Aggregate struct {
	user        User
	transaction *transaction.Transaction
	portfolio.Entity
}

func (a *Aggregate) User() User {
	user := a.user
	return user
}

func (a *Aggregate) Transaction() *transaction.Transaction {
	return a.transaction
}
