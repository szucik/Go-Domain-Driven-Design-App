package transaction

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

var emptyFieldErr = errors.New("this field must not be empty")

type Transaction struct {
	ID            uuid.UUID       `json:"id"`
	UserName      string          `json:"user_id" validate:"required"`
	PortfolioName string          `json:"portfolio_id" validate:"required"`
	Symbol        string          `json:"symbol" validate:"required"`
	Quantity      decimal.Decimal `json:"quantity" validate:"required"`
	Price         decimal.Decimal `json:"price" validate:"required"`
	Created       time.Time       `json:"created"`
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
	switch {
	case isEmpty(t.UserName):
		return fmt.Errorf("%w: %s", emptyFieldErr, "username")
	case isEmpty(t.PortfolioName):
		return fmt.Errorf("%w: %s", emptyFieldErr, "portfolio-name")
	case isEmpty(t.Symbol):
		return fmt.Errorf("%w: %s", emptyFieldErr, "symbol")
	}
	if t.Quantity.LessThanOrEqual(decimal.NewFromInt(0)) {
		return fmt.Errorf("%w: %s", emptyFieldErr, "quantity")
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
