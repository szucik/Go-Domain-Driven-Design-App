package user_test

import (
	"context"
	"log"
	"os"
	"testing"

	shopspring "github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/szucik/trade-helper/database/fake"
	"github.com/szucik/trade-helper/user"
)

var (
	ctx         = context.Background()
	database    = fake.NewDatabase()
	userService = user.Users{
		Logger:       log.New(os.Stdout, "logger", log.LstdFlags),
		Database:     &database,
		NewAggregate: user.User.NewAggregate,
	}
)

func TestUserService_SignUp(t *testing.T) {
	t.Run("should return userName when user creation is complete", func(t *testing.T) {
		// when
		username, err := userService.SignUp(ctx, user.User(testUser))
		assert.NoError(t, err)
		// then
		assert.Equal(t, instanceAggregate.User().Username, username)
	})

	t.Run("should return error when user with same email already exists", func(t *testing.T) {
		db := fake.NewDatabase()
		svc := user.Users{
			Logger:       log.New(os.Stdout, "logger", log.LstdFlags),
			Database:     &db,
			NewAggregate: user.User.NewAggregate,
		}
		_, err := svc.SignUp(ctx, user.User(testUser))
		require.NoError(t, err)

		_, err = svc.SignUp(ctx, user.User(testUser))
		require.Error(t, err)
	})
}

func TestUserService_GetUsers(t *testing.T) {
	t.Run("should return list of users", func(t *testing.T) {
		databaseWithThreeUserInstances(t)

		out, err := userService.GetUsers(ctx, user.PaginationIn{})
		require.NoError(t, err)
		assert.Len(t, out.Users, 3, "Three user instances")
		assert.Equal(t, 3, out.Total)
	})

	t.Run("should return paginated list of users", func(t *testing.T) {
		databaseWithThreeUserInstances(t)

		out, err := userService.GetUsers(ctx, user.PaginationIn{Page: 0, Limit: 2})
		require.NoError(t, err)
		assert.Len(t, out.Users, 2)
	})
}

func TestUserService_SignIn(t *testing.T) {
	db := fake.NewDatabase()
	svc := user.Users{
		Logger:       log.New(os.Stdout, "logger", log.LstdFlags),
		Database:     &db,
		NewAggregate: user.User.NewAggregate,
	}
	_, err := svc.SignUp(ctx, user.User(testUser))
	require.NoError(t, err)

	t.Run("should return username when credentials are valid", func(t *testing.T) {
		username, err := svc.SignIn(ctx, user.AuthCredentials{
			Email:    testUser.Email,
			Password: testUser.Password,
		})
		require.NoError(t, err)
		assert.Equal(t, testUser.Username, username)
	})

	t.Run("should return error when password is wrong", func(t *testing.T) {
		_, err := svc.SignIn(ctx, user.AuthCredentials{
			Email:    testUser.Email,
			Password: "wrongpassword",
		})
		require.Error(t, err)
	})

	t.Run("should return error when user does not exist", func(t *testing.T) {
		_, err := svc.SignIn(ctx, user.AuthCredentials{
			Email:    "nonexistent@test.com",
			Password: testUser.Password,
		})
		require.Error(t, err)
	})
}

func TestUserService_GetUserByEmail(t *testing.T) {
	db := fake.NewDatabase()
	svc := user.Users{
		Logger:       log.New(os.Stdout, "logger", log.LstdFlags),
		Database:     &db,
		NewAggregate: user.User.NewAggregate,
	}
	_, err := svc.SignUp(ctx, user.User(testUser))
	require.NoError(t, err)

	t.Run("should return user when email exists", func(t *testing.T) {
		response, err := svc.GetUserByEmail(ctx, testUser.Email)
		require.NoError(t, err)
		assert.Equal(t, testUser.Username, response.Username)
		assert.Equal(t, testUser.Email, response.Email)
	})

	t.Run("should return error when email does not exist", func(t *testing.T) {
		_, err := svc.GetUserByEmail(ctx, "nonexistent@test.com")
		require.Error(t, err)
	})
}

func TestUserService_GetUserByName(t *testing.T) {
	db := fake.NewDatabase()
	svc := user.Users{
		Logger:       log.New(os.Stdout, "logger", log.LstdFlags),
		Database:     &db,
		NewAggregate: user.User.NewAggregate,
	}
	_, err := svc.SignUp(ctx, user.User(testUser))
	require.NoError(t, err)

	t.Run("should return user when username exists", func(t *testing.T) {
		response, err := svc.GetUserByName(ctx, testUser.Username)
		require.NoError(t, err)
		assert.Equal(t, testUser.Username, response.Username)
		assert.Equal(t, testUser.Email, response.Email)
	})

	t.Run("should return error when username does not exist", func(t *testing.T) {
		_, err := svc.GetUserByName(ctx, "nonexistent")
		require.Error(t, err)
	})
}

func TestUserService_AddPortfolio(t *testing.T) {
	db := fake.NewDatabase()
	svc := user.Users{
		Logger:       log.New(os.Stdout, "logger", log.LstdFlags),
		Database:     &db,
		NewAggregate: user.User.NewAggregate,
	}
	_, err := svc.SignUp(ctx, user.User(testUser))
	require.NoError(t, err)

	t.Run("should return portfolio name when portfolio is created", func(t *testing.T) {
		name, err := svc.AddPortfolio(ctx, user.PortfolioIn{
			UserName: testUser.Username,
			Name:     "myportfolio",
		})
		require.NoError(t, err)
		assert.Equal(t, "myportfolio", name)
	})

	t.Run("should return error when portfolio name is invalid", func(t *testing.T) {
		_, err := svc.AddPortfolio(ctx, user.PortfolioIn{
			UserName: testUser.Username,
			Name:     "invalid name!",
		})
		require.Error(t, err)
	})

	t.Run("should return error when user does not exist", func(t *testing.T) {
		_, err := svc.AddPortfolio(ctx, user.PortfolioIn{
			UserName: "nonexistent",
			Name:     "portfolio",
		})
		require.Error(t, err)
	})

	t.Run("should return error when portfolio with same name already exists", func(t *testing.T) {
		_, err := svc.AddPortfolio(ctx, user.PortfolioIn{
			UserName: testUser.Username,
			Name:     "myportfolio",
		})
		require.Error(t, err)
	})
}

func TestUserService_AddTransaction(t *testing.T) {
	db := fake.NewDatabase()
	svc := user.Users{
		Logger:       log.New(os.Stdout, "logger", log.LstdFlags),
		Database:     &db,
		NewAggregate: user.User.NewAggregate,
	}
	_, err := svc.SignUp(ctx, user.User(testUser))
	require.NoError(t, err)

	_, err = svc.AddPortfolio(ctx, user.PortfolioIn{
		UserName: testUser.Username,
		Name:     "myportfolio",
	})
	require.NoError(t, err)

	t.Run("should return id when transaction is added", func(t *testing.T) {
		id, err := svc.AddTransaction(ctx, user.TransactionIn{
			UserName:      testUser.Username,
			PortfolioName: "myportfolio",
			Symbol:        "AAPL",
			Amount:        "150.00",
			Quantity:      "2",
		})
		require.NoError(t, err)
		assert.NotEmpty(t, id)
	})

	t.Run("should return error when portfolio does not exist", func(t *testing.T) {
		_, err := svc.AddTransaction(ctx, user.TransactionIn{
			UserName:      testUser.Username,
			PortfolioName: "nonexistent",
			Symbol:        "AAPL",
			Amount:        "150.00",
			Quantity:      "2",
		})
		require.Error(t, err)
	})

	t.Run("should return error when symbol is empty", func(t *testing.T) {
		_, err := svc.AddTransaction(ctx, user.TransactionIn{
			UserName:      testUser.Username,
			PortfolioName: "myportfolio",
			Symbol:        "",
			Amount:        "150.00",
			Quantity:      "2",
		})
		require.Error(t, err)
	})

	t.Run("should return error when quantity is zero", func(t *testing.T) {
		_, err := svc.AddTransaction(ctx, user.TransactionIn{
			UserName:      testUser.Username,
			PortfolioName: "myportfolio",
			Symbol:        "AAPL",
			Amount:        "150.00",
			Quantity:      "0",
		})
		require.Error(t, err)
	})

	t.Run("should update portfolio TotalCost after transaction", func(t *testing.T) {
		db := fake.NewDatabase()
		svc := user.Users{
			Logger:       log.New(os.Stdout, "logger", log.LstdFlags),
			Database:     &db,
			NewAggregate: user.User.NewAggregate,
		}
		_, err := svc.SignUp(ctx, user.User(testUser))
		require.NoError(t, err)
		_, err = svc.AddPortfolio(ctx, user.PortfolioIn{UserName: testUser.Username, Name: "tech"})
		require.NoError(t, err)

		_, err = svc.AddTransaction(ctx, user.TransactionIn{
			UserName:      testUser.Username,
			PortfolioName: "tech",
			Symbol:        "AAPL",
			Amount:        "150.00",
			Quantity:      "2",
		})
		require.NoError(t, err)

		response, err := svc.GetUserByName(ctx, testUser.Username)
		require.NoError(t, err)
		var found bool
		for _, p := range response.Portfolio {
			if p.Name == "tech" {
				found = true
				assert.True(t, p.TotalCost.Equal(shopspring.NewFromFloat(300.00)))
			}
		}
		require.True(t, found, "portfolio 'tech' not found in response")
	})
}

func TestUserService_GetTransactions(t *testing.T) {
	db := fake.NewDatabase()
	svc := user.Users{
		Logger:       log.New(os.Stdout, "logger", log.LstdFlags),
		Database:     &db,
		NewAggregate: user.User.NewAggregate,
	}
	_, err := svc.SignUp(ctx, user.User(testUser))
	require.NoError(t, err)
	_, err = svc.AddPortfolio(ctx, user.PortfolioIn{UserName: testUser.Username, Name: "tech"})
	require.NoError(t, err)
	_, err = svc.AddTransaction(ctx, user.TransactionIn{
		UserName: testUser.Username, PortfolioName: "tech",
		Symbol: "AAPL", Amount: "150.00", Quantity: "2", Type: "buy",
	})
	require.NoError(t, err)

	t.Run("should return transactions for portfolio", func(t *testing.T) {
		out, err := svc.GetTransactions(ctx, testUser.Username, "tech")
		require.NoError(t, err)
		assert.Len(t, out.Transactions, 1)
		assert.Equal(t, "AAPL", out.Transactions[0].Symbol)
	})

	t.Run("should return empty list when portfolio has no transactions", func(t *testing.T) {
		_, err = svc.AddPortfolio(ctx, user.PortfolioIn{UserName: testUser.Username, Name: "empty"})
		require.NoError(t, err)
		out, err := svc.GetTransactions(ctx, testUser.Username, "empty")
		require.NoError(t, err)
		assert.Empty(t, out.Transactions)
	})
}

func TestUserService_AddTransaction_BuyIncreasesTotalCost(t *testing.T) {
	db := fake.NewDatabase()
	svc := user.Users{
		Logger:       log.New(os.Stdout, "logger", log.LstdFlags),
		Database:     &db,
		NewAggregate: user.User.NewAggregate,
	}
	_, err := svc.SignUp(ctx, user.User(testUser))
	require.NoError(t, err)
	_, err = svc.AddPortfolio(ctx, user.PortfolioIn{UserName: testUser.Username, Name: "tech"})
	require.NoError(t, err)

	_, err = svc.AddTransaction(ctx, user.TransactionIn{
		UserName: testUser.Username, PortfolioName: "tech",
		Symbol: "AAPL", Amount: "100.00", Quantity: "3", Type: "buy",
	})
	require.NoError(t, err)

	_, err = svc.AddTransaction(ctx, user.TransactionIn{
		UserName: testUser.Username, PortfolioName: "tech",
		Symbol: "AAPL", Amount: "100.00", Quantity: "1", Type: "sell",
	})
	require.NoError(t, err)

	response, err := svc.GetUserByName(ctx, testUser.Username)
	require.NoError(t, err)
	for _, p := range response.Portfolio {
		if p.Name == "tech" {
			// buy 3x100 = 300, sell 1x100 = -100 => 200
			assert.True(t, p.TotalCost.Equal(shopspring.NewFromFloat(200.00)))
		}
	}
}

func TestUserService_UpdateUser(t *testing.T) {
	db := fake.NewDatabase()
	svc := user.Users{
		Logger:       log.New(os.Stdout, "logger", log.LstdFlags),
		Database:     &db,
		NewAggregate: user.User.NewAggregate,
	}
	_, err := svc.SignUp(ctx, user.User(testUser))
	require.NoError(t, err)

	t.Run("should update username", func(t *testing.T) {
		newName, err := svc.UpdateUser(ctx, testUser.Username, user.UpdateUserIn{
			Username: "newusername",
		})
		require.NoError(t, err)
		assert.Equal(t, "newusername", newName)

		_, err = svc.GetUserByName(ctx, "newusername")
		require.NoError(t, err)
	})

	t.Run("should update email", func(t *testing.T) {
		_, err := svc.UpdateUser(ctx, "newusername", user.UpdateUserIn{
			Email: "new@test.com",
		})
		require.NoError(t, err)

		response, err := svc.GetUserByName(ctx, "newusername")
		require.NoError(t, err)
		assert.Equal(t, "new@test.com", response.Email)
	})

	t.Run("should return error when user does not exist", func(t *testing.T) {
		_, err := svc.UpdateUser(ctx, "nonexistent", user.UpdateUserIn{Username: "x"})
		require.Error(t, err)
	})

	t.Run("should return error when new username is too short", func(t *testing.T) {
		_, err := svc.UpdateUser(ctx, "newusername", user.UpdateUserIn{Username: "x"})
		require.Error(t, err)
	})
}

func TestUserService_DeleteUser(t *testing.T) {
	db := fake.NewDatabase()
	svc := user.Users{
		Logger:       log.New(os.Stdout, "logger", log.LstdFlags),
		Database:     &db,
		NewAggregate: user.User.NewAggregate,
	}
	_, err := svc.SignUp(ctx, user.User(testUser))
	require.NoError(t, err)

	t.Run("should delete existing user", func(t *testing.T) {
		err := svc.DeleteUser(ctx, testUser.Username)
		require.NoError(t, err)

		_, err = svc.GetUserByName(ctx, testUser.Username)
		require.Error(t, err)
	})

	t.Run("should return error when user does not exist", func(t *testing.T) {
		err := svc.DeleteUser(ctx, "nonexistent")
		require.Error(t, err)
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
		_, err = database.SignUp(ctx, aggregate)
		require.NoError(t, err)
	}
}
