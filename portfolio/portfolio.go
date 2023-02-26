package portfolio

import (
	"errors"
	"regexp"
	"strings"
	"time"

	"github.com/shopspring/decimal"
)

type Portfolio struct {
	Name            string          `json:"name,omitempty" validate:"required"`
	TotalBalance    decimal.Decimal `json:"total_balance"`
	TotalCost       decimal.Decimal `json:"total_cost"`
	TotalProfitLoss decimal.Decimal `json:"total_profit_loss"`
	ProfitLossDay   decimal.Decimal `json:"profit_loss_day"`
	Created         time.Time       `json:"created"`
}

// Entity is an object with lifecycle (can be mutated by adding, removing and renaming). It can be embedded into other
// entities to avoid repetition.
type Entity struct {
	portfolio Portfolio
}

func (p Portfolio) NewPortfolio() (Entity, error) {
	err := validateName(p.Name)
	if err != nil {
		return Entity{}, err
	}

	return Entity{
		portfolio: Portfolio{
			Name:            p.Name,
			TotalBalance:    p.TotalBalance,
			TotalCost:       p.TotalCost,
			TotalProfitLoss: p.TotalProfitLoss,
			ProfitLossDay:   p.ProfitLossDay,
			Created:         p.Created,
		},
	}, nil
}

func (e *Entity) Portfolio() Portfolio {
	return e.portfolio
}

func validateName(name string) error {
	regex := regexp.MustCompile(`^[a-zA-Z0-9]+$`)
	name = strings.Trim(name, " ")

	if name == "" {
		return errors.New("the portfolio name must not be empty")
	}

	if !regex.MatchString(name) {
		return errors.New("the portfolio name can contain numbers, lowercase letters, uppercase letters")
	}

	return nil
}
