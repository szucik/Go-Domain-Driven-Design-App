package portfolio_test

import (
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/szucik/trade-helper/clock"
	pto "github.com/szucik/trade-helper/portfolio"
)

type Portfolio pto.Portfolio

var (
	entity    = pto.Entity{}
	portfolio = Portfolio{
		Name:            "name",
		TotalBalance:    decimal.Decimal{},
		TotalCost:       decimal.Decimal{},
		TotalProfitLoss: decimal.Decimal{},
		ProfitLossDay:   decimal.Decimal{},
		Created:         clock.FakeTime(),
	}
)

func (p Portfolio) WithName(name string) Portfolio {
	p.Name = name
	return p
}

func TestPortfolio_AddPortfolio(t *testing.T) {
	t.Run("should return an error when portfolio name", func(t *testing.T) {
		testCases := map[string]struct {
			p Portfolio
		}{
			"is empty": {
				p: Portfolio{
					Name: "",
				},
			},
			"has unauthorized signs": {
				p: Portfolio{
					Name: "name123$!-+ ",
				},
			},
		}

		for name, testCase := range testCases {
			t.Run(name, func(t *testing.T) {
				// when
				err := entity.AddPortfolio(pto.Portfolio(testCase.p))
				// then
				require.Error(t, err)
			})
		}
	})

	t.Run("should create new Portfolio", func(t *testing.T) {
		// when
		err := entity.AddPortfolio(pto.Portfolio(portfolio))
		require.NoError(t, err)
		// then
		assert.Equal(t, entity.Portfolios()[0].Name, portfolio.Name)
	})
}

func TestPortfolio_FindPortfolio(t *testing.T) {
	portfolioInstances(t)

	t.Run("should return an error when portfolio does not exist ", func(t *testing.T) {
		_, err := entity.FindPortfolio("nonexistent")
		require.Error(t, err)
	})

	t.Run("should find portfolio ", func(t *testing.T) {
		p, err := entity.FindPortfolio("name1")
		require.NoError(t, err)
		assert.Equal(t, "name1", p.Name)
	})
}

func portfolioInstances(t *testing.T) {
	entity = pto.Entity{}

	p1 := Portfolio{
		Name:            "name1",
		TotalBalance:    decimal.Decimal{},
		TotalCost:       decimal.Decimal{},
		TotalProfitLoss: decimal.Decimal{},
		ProfitLossDay:   decimal.Decimal{},
		Created:         clock.FakeTime(),
	}

	// when
	err := entity.AddPortfolio(pto.Portfolio(p1))
	require.NoError(t, err)
	err = entity.AddPortfolio(pto.Portfolio(p1.WithName("name2")))
	require.NoError(t, err)
}
