package user

import (
	"errors"
	"log"
)

var (
	ErrHash              = errors.New("Problem with hashing your password")
	MsgUserAlreadyExists = "UserKey already exists with the given email"
)

type UsersService interface {
	SignUp(user User) error
	SignIn() error
	GetUser(userId string) error
	GetUsers() ([]Aggregate, error)
	Update() error
}

type Repository interface {
	// GetUsers UpdateUser(ctx context.Context) (Aggregate, error)
	//Dashboard(ctx context.Context) (Aggregate, error)
	//SignIn(ctx context.Context) (Aggregate, error)
	GetUsers() ([]Aggregate, error)
	SignUp(aggregate Aggregate) error
}

type Users struct {
	Logger       *log.Logger
	Database     Repository
	NewAggregate func(User) (Aggregate, error)
}

func (u Users) SignUp(user User) error {
	aggregate, _ := u.NewAggregate(user)
	err := u.Database.SignUp(aggregate)
	if err != nil {
		return err
	}
	return nil
}

func (u Users) GetUsers() ([]Aggregate, error) {
	var (
		users []Aggregate
		err   error
	)
	users, err = u.Database.GetUsers()
	if err != nil {
		return users, err
	}
	return users, nil
}

func (u Users) SignIn() error {
	//TODO implement me
	panic("implement me")
}

func (u Users) GetUser(userId string) error {
	//TODO implement me
	panic("implement me")
}

func (u Users) Update() error {
	//TODO implement me
	panic("implement me")
}
