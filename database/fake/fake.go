package fake

import (
	"context"
	"errors"
	"fmt"
	"github.com/szucik/trade-helper/apperrors"
	"net/http"
	"sync"

	"github.com/google/uuid"

	"github.com/szucik/trade-helper/transaction"
	"github.com/szucik/trade-helper/user"
)

type userKey struct {
	name string
}

type MemoryRepository struct {
	users map[userKey]user.Aggregate
	sync.Mutex
	transactions map[string]map[string]map[uuid.UUID]transaction.Transaction
}

func (mr MemoryRepository) GetUserByName(ctx context.Context, userName string) (user.Aggregate, error) {
	//TODO implement me
	panic("implement me")
}

func NewDatabase() MemoryRepository {
	return MemoryRepository{
		users:        map[userKey]user.Aggregate{},
		transactions: map[string]map[string]map[uuid.UUID]transaction.Transaction{},
	}
}

func (mr MemoryRepository) SignUp(_ context.Context, aggregate user.Aggregate) (string, error) {
	mr.Lock()
	defer mr.Unlock()

	if mr.users == nil {
		mr.users = make(map[userKey]user.Aggregate)
	}

	key := userKey{
		name: aggregate.User().Email,
	}

	if _, ok := mr.users[key]; ok {
		return "", apperrors.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "user already exists",
			Type:    "DuplicateUser",
		}
	}

	mr.users[key] = aggregate

	return aggregate.User().Username, nil
}

func (mr MemoryRepository) GetUsers(_ context.Context) ([]user.Aggregate, error) {
	mr.Lock()
	defer mr.Unlock()

	var users []user.Aggregate
	if len(mr.users) > 0 {
		for _, aggregate := range mr.users {
			users = append(users, aggregate)
		}
	}

	return users, nil
}

func (mr MemoryRepository) GetUserByEmail(_ context.Context, email string) (user.Aggregate, error) {
	mr.Lock()
	defer mr.Unlock()

	key := userKey{
		name: email,
	}

	if user, exist := mr.users[key]; exist {
		return user, nil
	}

	return user.Aggregate{}, apperrors.ErrorResponse{
		Code:    http.StatusNotFound,
		Message: "user not found",
		Type:    "UserNotFound",
	}
}

func (mr MemoryRepository) SaveAggregate(_ context.Context, aggregate user.Aggregate) error {
	mr.Lock()
	defer mr.Unlock()

	key := userKey{
		name: aggregate.User().Email,
	}

	_, exist := mr.users[key]
	if !exist {
		return errors.New("user dont exist")

	}

	mr.users[key] = aggregate

	return nil
}

func (mr MemoryRepository) AddTransaction(_ context.Context, vo transaction.ValueObject) (string, error) {
	mr.Lock()
	defer mr.Unlock()

	t := vo.Transaction()

	if mr.transactions[t.UserName] == nil {
		mr.transactions[t.UserName] = make(map[string]map[uuid.UUID]transaction.Transaction)
	}

	if mr.transactions[t.UserName][t.PortfolioName] == nil {
		mr.transactions[t.UserName][t.PortfolioName] = make(map[uuid.UUID]transaction.Transaction)
	}

	mr.transactions[t.UserName][t.PortfolioName][t.ID] = t
	return fmt.Sprintf("%s: %s", t.Symbol, t.Quantity), nil
}

func (mr MemoryRepository) UpdateUser(ctx context.Context) (user.Aggregate, error) {
	// TODO implement me
	panic("implement me")
}

func (mr MemoryRepository) Dashboard(ctx context.Context) (user.Aggregate, error) {
	// TODO implement me
	panic("implement me")
}
