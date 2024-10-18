//go:build integration

package mongo_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/szucik/trade-helper/database/mongo"
	mongotest "github.com/szucik/trade-helper/database/mongo/internal/test"
)

const databaseName = "test"

var setup *mongotest.Mongo

func TestMain(m *testing.M) {
	var code int
	ctx, cancel := context.WithCancel(context.Background())

	defer func() {
		cancel()
		os.Exit(code)
	}()

	var err error
	setup, err = mongotest.RunMongoDB(ctx)
	if err != nil {
		fmt.Println(err)

		code = 1

		return
	}

	code = m.Run()

	if err := setup.Cleanup(); err != nil {
		fmt.Println(err)

		code = 2

		return
	}
}

func TestDatabase(t *testing.T) {
	test.Database(t, newDB)
}

// To jest do poprawy
func newDB(t *testing.T) test.DB {
	ctx := context.Background()

	db := setup.Client.Database(databaseName)
	err := db.Drop(ctx)
	require.NoError(t, err)

	mongoDB, err := mongo.NewDatabase(ctx)
	require.NoError(t, err)

	return mongoDB
}
