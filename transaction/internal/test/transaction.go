package test

import (
	"github.com/shopspring/decimal"

	txn "github.com/szucik/trade-helper/transaction"
)

type Transaction txn.Transaction

func (t Transaction) WithUserName(name string) Transaction {
	t.UserName = name
	return t
}

func (t Transaction) WithPortfolioName(name string) Transaction {
	t.PortfolioName = name
	return t
}

func (t Transaction) WithSymbol(symbol string) Transaction {
	t.Symbol = symbol
	return t
}

func (t Transaction) WithQuantity(quantity decimal.Decimal) Transaction {
	t.Quantity = quantity
	return t
}
