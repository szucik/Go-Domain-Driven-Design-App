package user

import (
	"regexp"
	"time"

	"github.com/shopspring/decimal"
	"github.com/szucik/trade-helper/apperrors"
	"github.com/szucik/trade-helper/portfolio"
	"github.com/szucik/trade-helper/transaction"
)

type User struct {
	Username  string
	Email     string
	Password  string
	TokenHash string
	Created   time.Time
	Updated   time.Time
}

type AuthCredentials struct {
	Email    string `json:"email"    validate:"required"`
	Password string `json:"password" validate:"required"`
}

type UpdateUserIn struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserResponse struct {
	Username  string                `json:"username"`
	Email     string                `json:"email"`
	Portfolio []portfolio.Portfolio `json:"portfolios"`
	Created   time.Time             `json:"created"`
}

type TransactionResponse struct {
	ID            string                   `json:"id"`
	Symbol        string                   `json:"symbol"`
	Type          transaction.TransactionType `json:"type"`
	Quantity      decimal.Decimal          `json:"quantity"`
	Price         decimal.Decimal          `json:"price"`
	Created       time.Time                `json:"created"`
}

type TransactionsOut struct {
	Transactions []TransactionResponse `json:"transactions"`
}

type PaginationIn struct {
	Page  int
	Limit int
}

func (u User) NewAggregate() (Aggregate, error) {
	switch {
	case !isLengthValid(u.Username, 2):
		return Aggregate{}, apperrors.Error(
			"User name is to short",
			"UserParamsValidation",
			400,
		)

	case !isEmailValid(u.Email):
		return Aggregate{}, apperrors.Error(
			"Invalid user email",
			"UserParamsValidation",
			400,
		)

	case len(u.Password) < 8:
		return Aggregate{},
			apperrors.Error(
				"Password is to short, it should be longer than 8 characters",
				"UserParamsValidation",
				400,
			)
	}

	return Aggregate{
		user: u,
	}, nil
}

func isLengthValid(value string, length int) bool {
	return len(value) >= length
}

func isEmailValid(email string) bool {
	var emailRegex = regexp.MustCompile("(?:[a-z0-9!#$%&'*+/=?^_`{|}~-]+(?:\\.[a-z0-9!#$%&'*+/=?^_`{|}~-]+)*|\"" +
		"(?:[\\x01-\\x08\\x0b\\x0c\\x0e-\\x1f\\x21\\x23-\\x5b\\x5d-\\x7f]|\\\\[\\x01-\\x09\\x0b\\x0c\\x0e-\\x7f])*\")" +
		"@(?:(?:[a-z0-9](?:[a-z0-9-]*[a-z0-9])?\\.)+[a-z0-9](?:[a-z0-9-]*[a-z0-9])?|\\[(?:(?:(2(5[0-5]|[0-4][0-9])|" +
		"1[0-9][0-9]|[1-9]?[0-9]))\\.){3}(?:(2(5[0-5]|[0-4][0-9])|1[0-9][0-9]|[1-9]?[0-9])|[a-z0-9-]*[a-z0-9]:" +
		"(?:[\\x01-\\x08\\x0b\\x0c\\x0e-\\x1f\\x21-\\x5a\\x53-\\x7f]|\\\\[\\x01-\\x09\\x0b\\x0c\\x0e-\\x7f])+)\\])")
	return isLengthValid(email, 6) && emailRegex.MatchString(email)
}
