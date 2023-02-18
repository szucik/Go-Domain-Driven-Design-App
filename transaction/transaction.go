package transaction

import (
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type Transaction struct {
	// the id for this transaction
	// required: true
	ID uuid.UUID `json:"id"`
	// the id for this transaction
	// required: true
	UserName string `json:"user_id" validate:"required"`
	// portfolio id  where the transaction will be stored
	// required: true
	PortfolioName string `json:"portfolio_id" validate:"required"`
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
	Created time.Time       `json:"created"`
}

type ValueObject struct {
	transaction Transaction
}

func (vo ValueObject) Transaction() Transaction {
	return vo.transaction
}

func (t Transaction) NewTransaction() (ValueObject, error) {
	err := t.validate()
	if err != nil {
		return ValueObject{}, err
	}

	return ValueObject{
		transaction: Transaction{
			ID:            t.ID,
			UserName:      t.UserName,
			PortfolioName: t.PortfolioName,
			Symbol:        t.Symbol,
			Quantity:      t.Quantity,
			Price:         t.Price,
			Created:       t.Created,
		},
	}, nil
}

func (t Transaction) validate() error {
	//  Todo more validation
	if isEmpty(t.UserName) || isEmpty(t.PortfolioName) || isEmpty(t.Symbol) {
		return errors.New("can't be empty")
	}

	return nil
}

func isEmpty(value string) bool {
	value = strings.Trim(value, " ")
	if value == "" {
		return true
	}

	return false
}
