package user_test

import (
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/szucik/trade-helper/clock"
	pto "github.com/szucik/trade-helper/portfolio"
	"github.com/szucik/trade-helper/user"
	"testing"
)

var (
	portfolio = pto.Portfolio{
		Name:            "name",
		TotalBalance:    decimal.Decimal{},
		TotalCost:       decimal.Decimal{},
		TotalProfitLoss: decimal.Decimal{},
		ProfitLossDay:   decimal.Decimal{},
		Created:         clock.FakeTime(),
	}

	instanceAggregate = func() user.Aggregate {
		aggregate, err := user.User(testUser).NewAggregate()
		if err != nil {
			panic(err)
		}

		return aggregate
	}()

	instancePortfolio = func() pto.Entity {
		newPortfolio, err := portfolio.NewPortfolio()
		if err != nil {
			panic(err)
		}
		return newPortfolio
	}()
)

func TestPortfolio_AddPortfolio(t *testing.T) {
	t.Run("should create new Portfolio entity", func(t *testing.T) {
		// when
		err := instanceAggregate.AddPortfolio(instancePortfolio)
		require.NoError(t, err)
		// then
		name := instanceAggregate.Portfolios()[0].Portfolio().Name
		assert.Equal(t, name, portfolio.Name)
	})
}

func TestPortfolio_FindPortfolio(t *testing.T) {
	portfolioInstances(t)

	t.Run("should return an error when portfolio does not exist ", func(t *testing.T) {
		_, err := instanceAggregate.FindPortfolio("nonexistent")
		require.Error(t, err)
	})

	t.Run("should find portfolio ", func(t *testing.T) {
		p, err := instanceAggregate.FindPortfolio("name1")
		require.NoError(t, err)
		assert.Equal(t, "name1", p.Name)
	})
}

func portfolioInstances(t *testing.T) {
	p1, _ := pto.Portfolio{
		Name:    "name1",
		Created: clock.FakeTime(),
	}.NewPortfolio()
	p2, _ := pto.Portfolio{
		Name:    "name2",
		Created: clock.FakeTime(),
	}.NewPortfolio()
	// when
	err := instanceAggregate.AddPortfolio(p1)
	require.NoError(t, err)
	err = instanceAggregate.AddPortfolio(p2)
	require.NoError(t, err)
}
