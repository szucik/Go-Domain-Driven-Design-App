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

func TestPortfolio_NewPortfolio(t *testing.T) {
	t.Run("should return an error when portfolio name ", func(t *testing.T) {
		testCases := map[string]struct {
			p Portfolio
		}{
			"is empty": {
				p: portfolio.WithName(""),
			},
			"has unauthorized signs": {
				p: portfolio.WithName("name123$!-+ "),
			},
			"contains only spaces": {
				p: portfolio.WithName("   "),
			},
		}

		for name, testCase := range testCases {
			t.Run(name, func(t *testing.T) {
				// when
				_, err := pto.Portfolio(testCase.p).NewPortfolio()
				// then
				require.Error(t, err)
			})
		}
	})

	t.Run("should create portfolio when name is valid", func(t *testing.T) {
		entity, err := pto.Portfolio(portfolio).NewPortfolio()
		require.NoError(t, err)
		assert.Equal(t, portfolio.Name, entity.Portfolio().Name)
	})
}
