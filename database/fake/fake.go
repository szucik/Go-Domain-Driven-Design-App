package fake

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/szucik/go-simple-rest-api/user"
	"sync"
)

var (
	// ErrCustomerNotFound is returned when a customer is not found.
	ErrCustomerNotFound = errors.New("the customer was not found in the repository")
	// ErrFailedToAddCustomer is returned when the customer could not be added to the repository.
	ErrFailedToAddCustomer = errors.New("failed to add the customer to the repository")
	// ErrUpdateCustomer is returned when the customer could not be updated in the repository.
	ErrUpdateCustomer = errors.New("failed to update the customer in the repository")
)

type MemoryRepository struct {
	user  map[uuid.UUID]user.Aggregate
	mutex *sync.Mutex
}

func NewDatabase() MemoryRepository {
	return MemoryRepository{
		user:  map[uuid.UUID]user.Aggregate{},
		mutex: &sync.Mutex{},
	}
}

func (mr MemoryRepository) GetUsers() ([]user.Aggregate, error) {
	var users []user.Aggregate
	if len(mr.user) > 0 {
		for _, aggregate := range mr.user {
			users = append(users, aggregate)
		}
	}

	return users, nil
}

func (mr MemoryRepository) GetUser(ctx context.Context) (user.Aggregate, error) {
	//TODO implement me
	panic("implement me")
}

func (mr MemoryRepository) UpdateUser(ctx context.Context) (user.Aggregate, error) {
	//TODO implement me
	panic("implement me")
}

func (mr MemoryRepository) Dashboard(ctx context.Context) (user.Aggregate, error) {
	//TODO implement me
	panic("implement me")
}

func (mr MemoryRepository) SignIn(ctx context.Context) (user.Aggregate, error) {
	//TODO implement me
	panic("implement me")
}

func (mr MemoryRepository) SignUp(aggregate user.Aggregate) error {
	if mr.user == nil {
		mr.mutex.Lock()
		mr.user = make(map[uuid.UUID]user.Aggregate)
		mr.mutex.Unlock()
	}

	if _, ok := mr.user[aggregate.GetID()]; ok {
		return fmt.Errorf("customer already exists: %w", ErrFailedToAddCustomer)
	}

	mr.mutex.Lock()
	mr.user[aggregate.GetID()] = aggregate
	mr.mutex.Unlock()
	fmt.Println(aggregate.GetID())
	return nil
}

//func (d UserRepository) DeletePortfolio(_ context.Context) (portfolio.Aggregate, error) {
//	d.mutex.Lock()
//	defer d.mutex.Unlock()
//
//	aggregate := portfolio.Aggregate{}
//
//	return aggregate, nil
//}
//
//func (d UserRepository) UpdatePortfolio(_ context.Context) (portfolio.Aggregate, error) {
//	d.mutex.Lock()
//	defer d.mutex.Unlock()
//
//	aggregate := portfolio.Aggregate{}
//
//	return aggregate, nil
//}
//func (d UserRepository) SaveTargetGroup(_ context.Context, aggregate targetgroup.Aggregate) error {
//	d.mutex.Lock()
//	defer d.mutex.Unlock()
//
//	version := aggregate.Version()
//	newTargetGroup := version == 0
//
//	document, found := d.targetGroupsByARN[aggregate.TargetGroup().ARN]
//
//	if !newTargetGroup {
//		if !found {
//			return fmt.Errorf("document update failed: previous version of the document not found")
//		}
//
//		if document.Version != version {
//			return fmt.Errorf("document update failed: invalid version")
//		}
//	}
//
//	if newTargetGroup && found {
//		return fmt.Errorf("document insert failed: duplicate target group")
//	}
//
//	if err := aggregate.SetVersion(version + 1); err != nil {
//		return err
//	}
//
//	arn := aggregate.TargetGroup().ARN
//	d.targetGroupsByARN[arn] = tgdocument.New(aggregate)
//
//	return nil
//}
