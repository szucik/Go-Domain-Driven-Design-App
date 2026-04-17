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

type TransactionType int

const (
	Buy  TransactionType = iota
	Sell TransactionType = iota
)

func (t TransactionType) String() string {
	if t == Sell {
		return "sell"
	}
	return "buy"
}

type Transaction struct {
	ID            uuid.UUID       `json:"id"           bson:"id"`
	UserName      string          `json:"user_id"      bson:"username"`
	PortfolioName string          `json:"portfolio_id" bson:"portfolio_name"`
	Symbol        string          `json:"symbol"       bson:"symbol"`
	Type          TransactionType `json:"type"         bson:"type"`
	Quantity      decimal.Decimal `json:"quantity"     bson:"quantity"`
	Price         decimal.Decimal `json:"price"        bson:"price"`
	Created       time.Time       `json:"created"      bson:"created"`
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
			Type:          t.Type,
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
