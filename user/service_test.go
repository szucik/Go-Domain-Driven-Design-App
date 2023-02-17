package user_test

import (
	"log"
	"os"
	"testing"
	"time"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/szucik/trade-helper/database/fake"
	"github.com/szucik/trade-helper/portfolio"
	"github.com/szucik/trade-helper/user"
)

var (
	database    = fake.NewDatabase()
	userService = user.Users{
		Logger:       log.New(os.Stdout, "logger", log.LstdFlags),
		Database:     &database,
		NewAggregate: user.User.NewAggregate,
	}

	instanceAggregate = func() user.Aggregate {
		aggregate, err := user.User(testUser).NewAggregate()
		if err != nil {
			panic(err)
		}

		return aggregate
	}()
)

func TestUserService_SignUp(t *testing.T) {
	t.Run("should return userName when user creation is complete", func(t *testing.T) {
		// when
		username, err := userService.SignUp(user.User(testUser))
		assert.NoError(t, err)
		// then
		assert.Equal(t, instanceAggregate.User().Username, username)
	})
}

func TestUserService_GetUsers(t *testing.T) {
	t.Run("should return list of users", func(t *testing.T) {
		databaseWithThreeUserInstances(t)

		// when
		out, err := userService.GetUsers()
		require.NoError(t, err)
		// then
		assert.Len(t, out.Users, 3, "Three user instances")
	})
}

func TestUserService_AddPortfolio(t *testing.T) {
	// TODO
	t.Run("should return an error when user dont exist", func(t *testing.T) {
		p := portfolio.Portfolio{
			Name:            "testName",
			TotalBalance:    decimal.Decimal{},
			TotalCost:       decimal.Decimal{},
			TotalProfitLoss: decimal.Decimal{},
			ProfitLossDay:   decimal.Decimal{},
			Transaction:     nil,
			Created:         time.Time{},
		}
		_, err := database.AddPortfolio(p)
		assert.NoError(t, err)
	})
}

func databaseWithThreeUserInstances(t *testing.T) {
	database = fake.NewDatabase()

	users := []user.User{
		user.User(testUser.WithEmail("email1@test.test").
			WithName("name1")),
		user.User(testUser.WithEmail("email2@test.test").
			WithName("name2")),
		user.User(testUser.WithEmail("email3@test.test").
			WithName("name3")),
	}

	for _, user := range users {
		aggregate, err := user.NewAggregate()
		require.NoError(t, err)
		_, err = database.SignUp(aggregate)
		require.NoError(t, err)
	}
}
