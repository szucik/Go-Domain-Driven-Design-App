package app

import (
	"context"
	"fmt"

	"github.com/szucik/go-simple-rest-api/database/fake"
	"github.com/szucik/go-simple-rest-api/portfolio"
)

//func Serve(ctx context.Context, config Config) error {

func Serve(ctx context.Context) error {
	database := fake.NewDatabase()
	p := portfolio.Portfolios{Database: database, NewAggregate: portfolio.Portfolio.NewAggregate}

	fmt.Println(p)
	return nil
}
