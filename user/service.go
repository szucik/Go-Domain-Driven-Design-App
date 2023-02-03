package user

import (
	"fmt"
	"log"
	"time"
)

type UsersService interface {
	SignUp(user User) (string, error)
	SignIn() error
	GetUser(userName string) (UserResponse, error)
	GetUsers() (UsersOut, error)
	Update() error
}
type Repository interface {
	// GetUser Dashboard(ctx context.Context) (Aggregate, error)
	// SignIn(ctx context.Context) (Aggregate, error)
	GetUser(userName string) (Aggregate, error)
	GetUsers() ([]Aggregate, error)
	SignUp(aggregate Aggregate) (string, error)
}

type Users struct {
	Logger       *log.Logger
	Database     Repository
	NewAggregate func(User) (Aggregate, error)
}

func (u Users) SignUp(user User) (string, error) {
	user.Created = time.Now()

	aggregate, err := u.NewAggregate(user)
	if err != nil {
		return "", err
	}

	id, err := u.Database.SignUp(aggregate)
	if err != nil {
		return "", fmt.Errorf("database.SignUp failed: %w", err)
	}

	return id, nil
}

type UsersOut struct {
	Users []UserResponse
}

func (u Users) GetUser(userName string) (UserResponse, error) {
	aggregate, err := u.Database.GetUser(userName)
	if err != nil {
		return UserResponse{}, fmt.Errorf("database.GetUser failed: %w", err)
	}

	return UserResponse{
		Username:  aggregate.User().Username,
		Email:     aggregate.User().Email,
		Portfolio: nil,
		Created:   aggregate.User().Created,
	}, nil
}

func (u Users) GetUsers() (UsersOut, error) {
	var response UsersOut
	users, err := u.Database.GetUsers()
	if err != nil {
		return UsersOut{}, fmt.Errorf("database.GetUsers failed: %w", err)
	}

	for _, user := range users {
		response.Users = append(response.Users, transformToUserResponse(user))
	}

	return response, nil
}

func transformToUserResponse(aggregate Aggregate) UserResponse {
	user := aggregate.User()

	return UserResponse{
		Username:  user.Username,
		Email:     user.Email,
		Created:   user.Created,
		Portfolio: nil,
	}
}

func (u Users) SignIn() error {
	// TODO implement me
	panic("implement me")
}

func (u Users) Update() error {
	// TODO implement me
	panic("implement me")
}
