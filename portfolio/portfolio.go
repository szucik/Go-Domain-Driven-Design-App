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
	portfolios []Portfolio
}

func (e *Entity) Portfolios() []Portfolio {
	return e.portfolios
}

func (e *Entity) FindPortfolio(name string) (Portfolio, error) {
	for _, p := range e.Portfolios() {
		if p.Name == name {
			return p, nil
		}
	}

	return Portfolio{}, errors.New("the specified portfolio does not exist")
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
