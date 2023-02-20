package transaction_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/szucik/trade-helper/clock"
	txn "github.com/szucik/trade-helper/transaction"
	"github.com/szucik/trade-helper/transaction/internal/test"
)

var (
	fakeName = "fakeName"
	instance = test.Transaction{
		ID:            uuid.MustParse("69359037-9599-48e7-b8f2-48393c019135"),
		UserName:      fakeName,
		PortfolioName: fakeName,
		Symbol:        "USD",
		Quantity:      decimal.New(1, 0),
		Price:         decimal.New(0, 0),
		Created:       clock.FakeTime(),
	}
)

func TestTransaction_NewTransaction(t *testing.T) {
	t.Run("should return an error when", func(t *testing.T) {
		testCases := map[string]struct {
			instance test.Transaction
		}{
			"userName is empty": {
				instance: instance.WithUserName(""),
			},
			"portfolioName is empty": {
				instance: instance.WithPortfolioName(""),
			},
			"symbol is empty": {
				instance: instance.WithSymbol(""),
			},
			"quantity is less than or equal to zero": {
				instance: instance.WithQuantity(decimal.New(0, 0)),
			},
		}

		for name, testCase := range testCases {
			t.Run(name, func(t *testing.T) {
				// when
				_, err := txn.Transaction(testCase.instance).NewTransaction()
				// then
				require.Error(t, err)
			})
		}
	})

	t.Run("should create Transaction", func(t *testing.T) {
		// when
		actual, err := txn.Transaction(instance).NewTransaction()
		require.NoError(t, err)
		// then
		assert.Equal(t, instance.Created, actual.Transaction().Created)
	})
}
