package user

import (
	"errors"
	"fmt"
	"github.com/szucik/trade-helper/portfolio"
	"math/rand"
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

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func randStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}
