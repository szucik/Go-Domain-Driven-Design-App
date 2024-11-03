package test

import (
	"context"
	"github.com/szucik/trade-helper/transaction"
	"github.com/szucik/trade-helper/user"
	"testing"
)

type DB interface {
	GetUserByEmail(ctx context.Context, email string) (user.Aggregate, error)
	GetUserByName(ctx context.Context, userName string) (user.Aggregate, error)
	GetUsers(ctx context.Context) ([]user.Aggregate, error)
	SignUp(ctx context.Context, aggregate user.Aggregate) (string, error)
	SaveAggregate(ctx context.Context, aggregate user.Aggregate) error
	AddTransaction(ctx context.Context, transaction transaction.ValueObject) (string, error)
}

func Database(t *testing.T, newDB func(t *testing.T) DB) {
	ctx := context.Background()
	aggregate, _ := user.User{
		Username: "username",
		Email:    "user@email.com",
		Password: "SomePassword123",
	}.NewAggregate()

	t.Run("should store user", func(t *testing.T) {
		db := newDB(t)
		_, err := db.SignUp(ctx, aggregate)
		if err != nil {
			return
		}
	})

}
