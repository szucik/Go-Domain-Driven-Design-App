package document

import (
	"github.com/szucik/trade-helper/portfolio"
	"github.com/szucik/trade-helper/user"
)

type Document struct {
	User       user.User
	Portfolios []portfolio.Entity
}

func NewDocument(aggregate user.Aggregate) Document {
	return Document{
		User:       aggregate.User(),
		Portfolios: aggregate.Portfolios(),
	}
}

func (d *Document) NewAggregate() (user.Aggregate, error) {
	aggregate, err := d.NewAggregate()
	if err != nil {
		return user.Aggregate{}, err
	}
	return aggregate, nil
}
