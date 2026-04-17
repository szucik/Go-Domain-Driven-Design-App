package user

import (
	"fmt"
	"math/rand"
	"net/http"

	"github.com/shopspring/decimal"
	"github.com/szucik/trade-helper/apperrors"
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

func (a *Aggregate) hashPassword() error {
	hash, err := hashPassword(a.user.Password)
	if err != nil {
		return fmt.Errorf("hashPassword failed: %w", err)
	}

	a.user.Password = hash
	return nil
}

func (a *Aggregate) AddPortfolio(entity portfolio.Entity) error {
	portfolio := entity.Portfolio()

	for _, item := range a.portfolios {
		if item.Portfolio().Name == portfolio.Name {
			return apperrors.Error("portfolio name already taken", "Conflict", http.StatusConflict)
		}
	}

	a.portfolios = append(a.portfolios, entity)

	return nil
}

func (a *Aggregate) UpdatePortfolioTotalCost(name string, delta decimal.Decimal) error {
	for i, p := range a.portfolios {
		pto := p.Portfolio()
		if pto.Name == name {
			pto.TotalCost = pto.TotalCost.Add(delta)
			entity, err := pto.NewPortfolio()
			if err != nil {
				return fmt.Errorf("UpdatePortfolioTotalCost: %w", err)
			}
			a.portfolios[i] = entity
			return nil
		}
	}
	return apperrors.Error("portfolio not found", "NotFound", http.StatusNotFound)
}

func (a *Aggregate) FindPortfolio(name string) (portfolio.Portfolio, error) {
	for _, p := range a.Portfolios() {
		pto := p.Portfolio()
		if pto.Name == name {
			return pto, nil
		}
	}

	return portfolio.Portfolio{}, apperrors.Error("portfolio not found", "NotFound", http.StatusNotFound)
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func randStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}
