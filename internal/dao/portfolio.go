package dao

import (
	"time"

	"github.com/shopspring/decimal"
)

// Portfolio represents user portfolio of coins
// swagger:model Portfolio
type Portfolio struct {
	// the id for this portfolio
	// required: true
	ID int `json:"id"`
	// the name for this portfolio
	// required: true
	// min length: 6
	Name string `json:"name,omitempty" validate:"required"`
	// the id of the owner of this portfolio
	// required: true
	UserId int `json:"user_id" validate:"required"`
	// the total value of the portfolio
	TotalBalance    decimal.Decimal `json:"total_balance"`
	TotalCost       decimal.Decimal `json:"total_cost"`
	TotalProfitLoss decimal.Decimal `json:"total_profit_loss"`
	ProfitLossDay   decimal.Decimal `json:"profit_loss_day"`
	Created         time.Time       `json:"created" sql:"created"`
}

func (d *Dao) Add(u *Portfolio) (int64, error) {
	stmt, err := d.db.Prepare("INSERT INTO portfolios (user_id, name) VALUES (?, ?)")
	if err != nil {
		panic(err.Error())
	}
	result, err := stmt.Exec(u.UserId, u.Name)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return id, nil
}
