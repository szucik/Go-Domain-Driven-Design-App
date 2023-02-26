package user

import (
	"errors"
	"github.com/szucik/trade-helper/portfolio"
)

type Aggregate struct {
	user       User
	portfolios []portfolio.Entity
}

func (a *Aggregate) User() User {
	user := a.user
	return user
}

func (a *Aggregate) Portfolios() []portfolio.Entity {
	return a.portfolios
}

func (a *Aggregate) AddPortfolio(entity portfolio.Entity) error {
	portfolio := entity.Portfolio()

	for _, item := range a.portfolios {
		if item.Portfolio().Name == portfolio.Name {
			return errors.New("this name is not available")
		}
	}

	a.portfolios = append(a.portfolios, entity)

	return nil
}

func (a *Aggregate) FindPortfolio(name string) (portfolio.Portfolio, error) {
	for _, p := range a.Portfolios() {
		pto := p.Portfolio()
		if pto.Name == name {
			return pto, nil
		}
	}

	return portfolio.Portfolio{}, errors.New("the specified portfolio does not exist")
}
