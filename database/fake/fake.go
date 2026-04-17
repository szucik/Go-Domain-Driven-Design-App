package fake

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"sync"

	"github.com/google/uuid"

	"github.com/szucik/trade-helper/apperrors"
	"github.com/szucik/trade-helper/transaction"
	"github.com/szucik/trade-helper/user"
)

type userKey struct {
	name string
}

type MemoryRepository struct {
	users        map[userKey]user.Aggregate
	transactions map[string]map[string]map[uuid.UUID]transaction.Transaction
	sync.Mutex
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

	key := userKey{name: aggregate.User().Email}

	if _, ok := mr.users[key]; ok {
		return "", apperrors.ErrorResponse{
			Code:    http.StatusConflict,
			Message: "user already exists",
			Type:    "DuplicateUser",
		}
	}

	mr.users[key] = aggregate
	return aggregate.User().Username, nil
}

func (mr MemoryRepository) GetUserByEmail(_ context.Context, email string) (user.Aggregate, error) {
	mr.Lock()
	defer mr.Unlock()

	if a, ok := mr.users[userKey{name: email}]; ok {
		return a, nil
	}

	return user.Aggregate{}, apperrors.ErrorResponse{
		Code:    http.StatusNotFound,
		Message: "user not found",
		Type:    "UserNotFound",
	}
}

func (mr MemoryRepository) GetUserByName(_ context.Context, userName string) (user.Aggregate, error) {
	mr.Lock()
	defer mr.Unlock()

	for _, a := range mr.users {
		if a.User().Username == userName {
			return a, nil
		}
	}

	return user.Aggregate{}, apperrors.ErrorResponse{
		Code:    http.StatusNotFound,
		Message: "user not found",
		Type:    "UserNotFound",
	}
}

func (mr MemoryRepository) GetUsers(_ context.Context, p user.PaginationIn) ([]user.Aggregate, error) {
	mr.Lock()
	defer mr.Unlock()

	all := make([]user.Aggregate, 0, len(mr.users))
	for _, a := range mr.users {
		all = append(all, a)
	}

	if p.Limit <= 0 {
		return all, nil
	}

	start := p.Page * p.Limit
	if start >= len(all) {
		return nil, nil
	}
	end := start + p.Limit
	if end > len(all) {
		end = len(all)
	}
	return all[start:end], nil
}

func (mr MemoryRepository) SaveAggregate(_ context.Context, aggregate user.Aggregate) error {
	mr.Lock()
	defer mr.Unlock()

	key := userKey{name: aggregate.User().Email}
	if _, ok := mr.users[key]; !ok {
		return errors.New("user does not exist")
	}
	mr.users[key] = aggregate
	return nil
}

func (mr MemoryRepository) UpdateUser(_ context.Context, currentUsername string, aggregate user.Aggregate) error {
	mr.Lock()
	defer mr.Unlock()

	var oldKey userKey
	var found bool
	for k, a := range mr.users {
		if a.User().Username == currentUsername {
			oldKey = k
			found = true
			break
		}
	}
	if !found {
		return apperrors.ErrorResponse{
			Code:    http.StatusNotFound,
			Message: "user not found",
			Type:    "UserNotFound",
		}
	}

	delete(mr.users, oldKey)
	mr.users[userKey{name: aggregate.User().Email}] = aggregate
	return nil
}

func (mr MemoryRepository) DeleteUser(_ context.Context, username string) error {
	mr.Lock()
	defer mr.Unlock()

	for k, a := range mr.users {
		if a.User().Username == username {
			delete(mr.users, k)
			return nil
		}
	}

	return apperrors.ErrorResponse{
		Code:    http.StatusNotFound,
		Message: "user not found",
		Type:    "UserNotFound",
	}
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

func (mr MemoryRepository) GetTransactions(_ context.Context, username, portfolioName string) ([]transaction.ValueObject, error) {
	mr.Lock()
	defer mr.Unlock()

	portfolios, ok := mr.transactions[username]
	if !ok {
		return nil, nil
	}
	txns, ok := portfolios[portfolioName]
	if !ok {
		return nil, nil
	}

	var result []transaction.ValueObject
	for _, t := range txns {
		vo, err := t.NewTransaction()
		if err != nil {
			return nil, err
		}
		result = append(result, vo)
	}
	return result, nil
}
