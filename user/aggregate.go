package user

import (
	"github.com/szucik/trade-helper/portfolio"
)

type Aggregate struct {
	user User
	portfolio.Entity
}

func (a *Aggregate) User() User {
	user := a.user
	return user
}
