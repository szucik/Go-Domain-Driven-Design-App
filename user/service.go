package user

import (
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"golang.org/x/crypto/bcrypt"

	"github.com/szucik/trade-helper/portfolio"
	"github.com/szucik/trade-helper/transaction"
)

type UsersService interface {
	SignUp(user User) (string, error)
	SignIn() error
	GetUser(userName string) (UserResponse, error)
	GetUsers() (UsersOut, error)
	AddPortfolio(in PortfolioIn) (string, error)
	AddTransaction(in TransactionIn) (string, error)
}

type Repository interface {
	// GetUser Dashboard(ctx context.Context) (Aggregate, error)
	// SignIn(ctx context.Context) (Aggregate, error)
	GetUser(userName string) (Aggregate, error)
	GetUsers() ([]Aggregate, error)
	SignUp(aggregate Aggregate) (string, error)
	SaveAggregate(aggregate Aggregate) error
	AddTransaction(transaction transaction.ValueObject) (string, error)
}

type Users struct {
	Logger       *log.Logger
	Database     Repository
	NewAggregate func(User) (Aggregate, error)
}

type TransactionIn struct {
	UserName      string
	PortfolioName string
	Symbol        string
	Amount        float64
	Quantity      float64
}

func (u Users) AddTransaction(in TransactionIn) (string, error) {
	aggregate, err := u.Database.GetUser(in.UserName)
	if err != nil {
		return "", fmt.Errorf("service.AddTransaction: %w", err)
	}

	_, err = aggregate.FindPortfolio(in.PortfolioName)
	if err != nil {
		return "", fmt.Errorf("service.AddTransaction: %w", err)
	}

	t, err := transaction.Transaction{
		ID:            uuid.New(),
		UserName:      in.UserName,
		PortfolioName: in.PortfolioName,
		Symbol:        in.Symbol,
		Created:       time.Now(),
		// Todo change this ->
		Quantity: decimal.Decimal{},
		Price:    decimal.Decimal{},
	}.NewTransaction()
	if err != nil {
		return "", fmt.Errorf("NewTransaction error: %w", err)
	}

	id, err := u.Database.AddTransaction(t)

	if err != nil {
		return "", fmt.Errorf("Database.AddTransaction: %w", err)
	}
	return id, nil
}

func (u Users) SignUp(user User) (string, error) {
	user.Created = time.Now()
	hash, err := hashPassword(user.Password)
	if err != nil {
		return "", fmt.Errorf("service.SignUp failed: %w", err)
	}

	user.Password = hash

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

	user := aggregate.User()

	return UserResponse{
		Username:  user.Username,
		Email:     user.Email,
		Portfolio: aggregate.Portfolios(),
		Created:   user.Created,
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
		Portfolio: aggregate.Portfolios(),
	}
}

func (u Users) SignIn() error {
	// TODO implement me
	panic("implement me")
}

type PortfolioIn struct {
	UserName string
	Name     string
}

func (u Users) AddPortfolio(in PortfolioIn) (name string, _ error) {
	aggregate, err := u.Database.GetUser(in.UserName)
	if err != nil {
		return "", fmt.Errorf("database.GetUser failed: %w", err)
	}

	p := portfolio.Portfolio{
		Name:            in.Name,
		TotalBalance:    decimal.NewFromFloat(0),
		TotalCost:       decimal.NewFromFloat(0),
		TotalProfitLoss: decimal.NewFromInt(0),
		ProfitLossDay:   decimal.NewFromInt(0),
		Created:         time.Now(),
	}

	err = aggregate.AddPortfolio(p)
	if err != nil {
		return "", fmt.Errorf("aggregate.AddPortfolio failed: %w", err)
	}

	err = u.Database.SaveAggregate(aggregate)
	if err != nil {
		return "", fmt.Errorf("aggregate.SaveAggregate failed: %w", err)
	}

	return p.Name, nil
}

func hashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 15)
	return string(hash), err
}
