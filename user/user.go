package user

import (
	"errors"
	"github.com/google/uuid"
	"github.com/szucik/go-simple-rest-api/transaction"
	"time"
)

var (
	invalidUserErr = errors.New("invalid user params")
)

type User struct {
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"username" validate:"required"`
	Email     string    `json:"email" validate:"required" sql:"email"`
	Password  string    `json:"password" validate:"required" sql:"password"`
	TokenHash string    `json:"token_hash" sql:"token_hash"`
	Created   time.Time `json:"created" sql:"created"`
	Updated   time.Time `json:"updated" sql:"updated"`
}

func (u User) NewAggregate() (Aggregate, error) {
	//TODO Add Validation error

	if u.Username == "" || u.Email == "" || u.Password == "" {
		return Aggregate{}, invalidUserErr
	}
	return Aggregate{
		user:        &u,
		transaction: &transaction.Transaction{},
	}, nil
}
