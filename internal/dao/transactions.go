package dao

import (
	"time"

	"github.com/shopspring/decimal"
)

// Transaction represents a single transaction in the portfolio
// swagger:model Transaction
type Transaction struct {
	// the id for this transaction
	// required: true
	ID int `json:"id"`
	// the id for this transaction
	// required: true
	UserId int `json:"user_id" validate:"required"`
	// portfolio id  where the transaction will be stored
	// required: true
	PortfolioId int `json:"portfolio_id" validate:"required"`
	// cryptocurrency short name
	// required: true
	// min length 2
	Symbol string `json:"symbol" validate:"required"`
	// the quantity of cryptocurrency purchased
	// required: true
	Quantity decimal.Decimal `json:"quantity" validate:"required"`
	// the price of the purchased cryptocurrency
	// required: true
	Price   decimal.Decimal `json:"price" validate:"required"`
	Created time.Time       `json:"created" sql:"created"`
}

func (d *Dao) AddTransaction(u *Transaction) (int64, error) {
	stmt, err := d.db.Prepare("INSERT INTO transactions (user_id, symbol, price, quantity) VALUES (?, ?, ?, ?)")
	if err != nil {
		panic(err.Error())
	}
	result, err := stmt.Exec(u.UserId, u.Symbol, u.Price, u.Quantity)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return id, nil
}
