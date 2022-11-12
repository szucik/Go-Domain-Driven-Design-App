package portfolio

import (
	"time"

	"github.com/shopspring/decimal"
)

type Portfolio struct {
	ID              int             `json:"id"`
	Name            string          `json:"name,omitempty" validate:"required"`
	UserId          int             `json:"user_id" validate:"required"`
	TotalBalance    decimal.Decimal `json:"total_balance"`
	TotalCost       decimal.Decimal `json:"total_cost"`
	TotalProfitLoss decimal.Decimal `json:"total_profit_loss"`
	ProfitLossDay   decimal.Decimal `json:"profit_loss_day"`
	Created         time.Time       `json:"created" sql:"created"`
}

func (p Portfolio) NewAggregate() (Aggregate, error) {
	return Aggregate{
		portfolio: p,
	}, nil
}
