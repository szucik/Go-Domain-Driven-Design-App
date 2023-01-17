package user

import (
	"errors"
	"log"

	"github.com/google/uuid"
)

var (
	ErrHash              = errors.New("Problem with hashing your password")
	MsgUserAlreadyExists = "UserKey already exists with the given email"
)

type Id uuid.UUID

type UsersService interface {
	SignUp(user User) (Id, error)
	SignIn() error
	GetUser(userId string) error

	// TODO GetUsers should return map of aggregates???
	GetUsers() ([]Aggregate, error)
	Update() error
}
type Repository interface {
	// GetUsers UpdateUser(ctx context.Context) (Aggregate, error)
	// Dashboard(ctx context.Context) (Aggregate, error)
	// SignIn(ctx context.Context) (Aggregate, error)
	GetUsers() ([]Aggregate, error)
	SignUp(aggregate Aggregate) (Id, error)
}

type Users struct {
	Logger       *log.Logger
	Database     Repository
	NewAggregate func(User) (Aggregate, error)
}

// func HashPassword(password string) (string, error) {
// 	hash, err := bcrypt.GenerateFromPassword([]byte(password), 15)
// 	return string(hash), err
// }

func (u Users) SignUp(user User) (Id, error) {

	// hash, err := HashPassword(user.Password)
	// if err != nil {
	// 	fmt.Errorf("%s", ErrHash)
	// }

	aggregate, err := u.NewAggregate(user)
	if err != nil {
		return Id{}, err
	}

	id, err := u.Database.SignUp(aggregate)
	if err != nil {
		return Id{}, err
	}

	return id, nil
}

func (u Users) GetUsers() ([]Aggregate, error) {
	users, err := u.Database.GetUsers()
	if err != nil {
		return users, err
	}

	return users, nil
}

func (u Users) SignIn() error {
	// TODO implement me
	panic("implement me")
}

func (u Users) GetUser(userId string) error {
	// TODO implement me
	panic("implement me")
}

func (u Users) Update() error {
	// TODO implement me
	panic("implement me")
}
