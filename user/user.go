package user

import (
	"regexp"
	"time"

	"github.com/szucik/trade-helper/portfolio"
	"github.com/szucik/trade-helper/web"
)

type User struct {
	Username  string    `json:"username" validate:"required"`
	Email     string    `json:"email" validate:"required"`
	Password  string    `json:"password" validate:"required"`
	TokenHash string    `json:"token_hash"`
	Created   time.Time `json:"created"`
	Updated   time.Time `json:"updated"`
}

type UserResponse struct {
	Username  string                `json:"username" validate:"required"`
	Email     string                `json:"email" validate:"required"`
	Portfolio []portfolio.Portfolio `json:"portfolios" validate:"required"`
	Created   time.Time             `json:"created"`
}

func (u User) NewAggregate() (Aggregate, error) {
	switch {
	case !isLengthValid(u.Username, 2):
		return Aggregate{}, web.BadRequestError(
			"User name is to short",
			"UserParamsValidation",
		)

	case !isEmailValid(u.Email):
		return Aggregate{}, web.BadRequestError(
			"Invalid user email",
			"UserParamsValidation",
		)

	case !isLengthValid(u.Password, 8):
		return Aggregate{},
			web.BadRequestError(
				"Password is to short, it should be longer than 8 characters",
				"UserParamsValidation",
			)
	}

	return Aggregate{
		user: u,
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
