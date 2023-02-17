package fake

import (
	"context"
	"errors"
	"net/http"
	"sync"

	"github.com/google/uuid"

	"github.com/szucik/trade-helper/user"
	"github.com/szucik/trade-helper/web"
)

type MemoryRepository struct {
	user map[string]user.Aggregate
	sync.Mutex
}

func NewDatabase() MemoryRepository {
	return MemoryRepository{
		user: map[string]user.Aggregate{},
	}
}

func (mr MemoryRepository) SignUp(aggregate user.Aggregate) (string, error) {
	mr.Lock()
	defer mr.Unlock()

	if mr.user == nil {
		mr.user = make(map[string]user.Aggregate)
	}

	if _, ok := mr.user[aggregate.User().Username]; ok {
		return "", web.ErrorResponse{
			TraceId: uuid.New().String(),
			Code:    http.StatusBadRequest,
			Message: "user already exists",
			Type:    "DuplicateUser",
		}
	}

	mr.user[aggregate.User().Username] = aggregate

	return aggregate.User().Username, nil
}

func (mr MemoryRepository) GetUsers() ([]user.Aggregate, error) {
	mr.Lock()
	defer mr.Unlock()

	var users []user.Aggregate
	if len(mr.user) > 0 {
		for _, aggregate := range mr.user {
			users = append(users, aggregate)
		}
	}

	return users, nil
}

func (mr MemoryRepository) GetUser(username string) (user.Aggregate, error) {
	mr.Lock()
	defer mr.Unlock()

	if user, exist := mr.user[username]; exist {
		return user, nil
	}

	return user.Aggregate{}, web.ErrorResponse{
		TraceId: uuid.New().String(),
		Code:    http.StatusNotFound,
		Message: "user not found",
		Type:    "UserNotFound",
	}
}

func (mr MemoryRepository) UpdateUser(ctx context.Context) (user.Aggregate, error) {
	// TODO implement me
	panic("implement me")
}

func (mr MemoryRepository) Dashboard(ctx context.Context) (user.Aggregate, error) {
	// TODO implement me
	panic("implement me")
}

func (mr MemoryRepository) SaveAggregate(aggregate user.Aggregate) error {
	mr.Lock()
	defer mr.Unlock()

	val, exist := mr.user[aggregate.User().Username]
	if !exist {
		return errors.New("user dont exist")
	}

	mr.user[val.User().Username] = aggregate

	return nil
}
