package fake

import (
	"context"
	"sync"

	"github.com/szucik/go-simple-rest-api/portfolio"
)

type Database struct {
	mutex *sync.Mutex
}

func NewDatabase() Database {
	return Database{
		mutex: &sync.Mutex{},
	}
}

func (d Database) AddPortfolio(ctx context.Context) (portfolio.Aggregate, error) {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	aggregate := portfolio.Aggregate{}

	return aggregate, nil
}

//
//func (d Database) DeletePortfolio(_ context.Context) (portfolio.Aggregate, error) {
//	d.mutex.Lock()
//	defer d.mutex.Unlock()
//
//	aggregate := portfolio.Aggregate{}
//
//	return aggregate, nil
//}
//
//func (d Database) UpdatePortfolio(_ context.Context) (portfolio.Aggregate, error) {
//	d.mutex.Lock()
//	defer d.mutex.Unlock()
//
//	aggregate := portfolio.Aggregate{}
//
//	return aggregate, nil
//}
