package fake

import (
	"context"
	"errors"
	"sync"

	"github.com/google/uuid"

	"github.com/szucik/trade-helper/user"
	"github.com/szucik/trade-helper/web"
)

var (
	errUserNotFound    = errors.New("the user was not found in the repository")
	errFailedToAddUser = errors.New("failed to add the user to the repository")
	errUpdateUser      = errors.New("failed to update the user in the repository")
)

type MemoryRepository struct {
	user map[string]user.Aggregate
	sync.Mutex
}

func NewDatabase() MemoryRepository {
	return MemoryRepository{
		user: make(map[string]user.Aggregate),
	}
}

func (mr MemoryRepository) SignUp(aggregate user.Aggregate) (user.Id, error) {
	mr.Lock()
	defer mr.Unlock()

	if mr.user == nil {
		mr.user = make(map[string]user.Aggregate)
	}
	id := uuid.New()
	aggregate.SetId(id)

	if _, ok := mr.user[aggregate.User().Username]; ok {
		return user.Id{}, web.ErrorResponse{
			TraceId: id.String(),
			Code:    400,
			Message: "user already exists",
			Type:    "DuplicateUser",
		}

	}
	mr.user[aggregate.User().Username] = aggregate
	return user.Id(aggregate.GetID()), nil
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

func (mr MemoryRepository) GetUser(ctx context.Context) (user.Aggregate, error) {
	// TODO implement me
	panic("implement me")
}

func (mr MemoryRepository) UpdateUser(ctx context.Context) (user.Aggregate, error) {
	// TODO implement me
	panic("implement me")
}

func (mr MemoryRepository) Dashboard(ctx context.Context) (user.Aggregate, error) {
	// TODO implement me
	panic("implement me")
}

func (mr MemoryRepository) SignIn(ctx context.Context) (user.Aggregate, error) {
	// TODO implement me
	panic("implement me")
}
