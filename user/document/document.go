package document

import (
	"time"

	"github.com/shopspring/decimal"

	"github.com/szucik/trade-helper/portfolio"
	"github.com/szucik/trade-helper/user"
)

type PortfolioDocument struct {
	Name            string    `bson:"name"`
	TotalBalance    string    `bson:"total_balance"`
	TotalCost       string    `bson:"total_cost"`
	TotalProfitLoss string    `bson:"total_profit_loss"`
	ProfitLossDay   string    `bson:"profit_loss_day"`
	Created         time.Time `bson:"created"`
}

type Document struct {
	User       user.User           `bson:"user"`
	Portfolios []PortfolioDocument `bson:"portfolios"`
}

func NewDocument(aggregate user.Aggregate) Document {
	portfolios := make([]PortfolioDocument, 0, len(aggregate.Portfolios()))
	for _, e := range aggregate.Portfolios() {
		p := e.Portfolio()
		portfolios = append(portfolios, PortfolioDocument{
			Name:            p.Name,
			TotalBalance:    p.TotalBalance.String(),
			TotalCost:       p.TotalCost.String(),
			TotalProfitLoss: p.TotalProfitLoss.String(),
			ProfitLossDay:   p.ProfitLossDay.String(),
			Created:         p.Created,
		})
	}
	return Document{
		User:       aggregate.User(),
		Portfolios: portfolios,
	}
}

func (d *Document) NewAggregate() (user.Aggregate, error) {
	aggregate, err := d.User.NewAggregate()
	if err != nil {
		return user.Aggregate{}, err
	}
	for _, pd := range d.Portfolios {
		entity, err := pd.toEntity()
		if err != nil {
			return user.Aggregate{}, err
		}
		if err := aggregate.AddPortfolio(entity); err != nil {
			return user.Aggregate{}, err
		}
	}
	return aggregate, nil
}

func (pd PortfolioDocument) toEntity() (portfolio.Entity, error) {
	totalBalance, err := decimal.NewFromString(pd.TotalBalance)
	if err != nil {
		return portfolio.Entity{}, err
	}
	totalCost, err := decimal.NewFromString(pd.TotalCost)
	if err != nil {
		return portfolio.Entity{}, err
	}
	totalProfitLoss, err := decimal.NewFromString(pd.TotalProfitLoss)
	if err != nil {
		return portfolio.Entity{}, err
	}
	profitLossDay, err := decimal.NewFromString(pd.ProfitLossDay)
	if err != nil {
		return portfolio.Entity{}, err
	}

	return portfolio.Portfolio{
		Name:            pd.Name,
		TotalBalance:    totalBalance,
		TotalCost:       totalCost,
		TotalProfitLoss: totalProfitLoss,
		ProfitLossDay:   profitLossDay,
		Created:         pd.Created,
	}.NewPortfolio()
}
