package user

import (
	"regexp"
	"time"

	"github.com/google/uuid"

	"github.com/szucik/trade-helper/web"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"username" validate:"required"`
	Email     string    `json:"email" validate:"required"`
	Password  string    `json:"password" validate:"required"`
	TokenHash string    `json:"token_hash"`
	Created   time.Time `json:"created"`
	Updated   time.Time `json:"updated"`
}

func (u User) NewAggregate() (Aggregate, error) {
	switch {
	case !isLengthValid(u.Username, 2):
		return Aggregate{}, web.ErrorResponse{
			TraceId: "",
			Code:    400,
			Message: "user name is to short",
			Type:    "UserParamsValidation",
		}

	case !isEmailValid(u.Email):
		return Aggregate{}, web.ErrorResponse{
			TraceId: "",
			Code:    400,
			Message: "user email is to short",
			Type:    "UserParamsValidation",
		}

	case !isLengthValid(u.Password, 8):
		return Aggregate{}, web.ErrorResponse{
			TraceId: "",
			Code:    400,
			Message: "user password is to short",
			Type:    "UserParamsValidation",
		}
	}
	return Aggregate{
		user: u,
		// transaction: &transaction.Transaction{},
	}, nil
}

func isEmailValid(email string) bool {
	var emailRegex = regexp.MustCompile("(?:[a-z0-9!#$%&'*+/=?^_`{|}~-]+(?:\\.[a-z0-9!#$%&'*+/=?^_`{|}~-]+)*|\"" +
		"(?:[\\x01-\\x08\\x0b\\x0c\\x0e-\\x1f\\x21\\x23-\\x5b\\x5d-\\x7f]|\\\\[\\x01-\\x09\\x0b\\x0c\\x0e-\\x7f])*\")" +
		"@(?:(?:[a-z0-9](?:[a-z0-9-]*[a-z0-9])?\\.)+[a-z0-9](?:[a-z0-9-]*[a-z0-9])?|\\[(?:(?:(2(5[0-5]|[0-4][0-9])|" +
		"1[0-9][0-9]|[1-9]?[0-9]))\\.){3}(?:(2(5[0-5]|[0-4][0-9])|1[0-9][0-9]|[1-9]?[0-9])|[a-z0-9-]*[a-z0-9]:" +
		"(?:[\\x01-\\x08\\x0b\\x0c\\x0e-\\x1f\\x21-\\x5a\\x53-\\x7f]|\\\\[\\x01-\\x09\\x0b\\x0c\\x0e-\\x7f])+)\\])")
	return isLengthValid(email, 6) && emailRegex.MatchString(email)
}

func isLengthValid(value string, length int) bool {
	if len(value) < length {
		return false
	}
	return true
}
