package portfolio

import (
	"errors"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"

	"github.com/szucik/trade-helper/transaction"
)

type Portfolio struct {
	ID              uuid.UUID                 `json:"id"`
	Name            string                    `json:"name,omitempty" validate:"required"`
	TotalBalance    decimal.Decimal           `json:"total_balance"`
	TotalCost       decimal.Decimal           `json:"total_cost"`
	TotalProfitLoss decimal.Decimal           `json:"total_profit_loss"`
	ProfitLossDay   decimal.Decimal           `json:"profit_loss_day"`
	Transaction     []transaction.Transaction `json:"transactions"`
	Created         time.Time                 `json:"created"`
}

// Entity is an object with lifecycle (can be mutated by adding, removing and renaming). It can be embedded into other
// entities to avoid repetition.
type Entity struct {
	portfolios []Portfolio
}

func (e *Entity) Portfolios() []Portfolio {
	return e.portfolios
}

func (e *Entity) AddPortfolio(p Portfolio) error {
	err := validateName(p.Name)
	if err != nil {
		return err
	}

	for _, item := range e.Portfolios() {
		if item.Name == p.Name {
			return errors.New("this name is not available")
		}
	}

	e.portfolios = append(e.portfolios, p)

	return nil
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
