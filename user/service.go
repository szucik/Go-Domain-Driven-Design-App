package user

import (
	"fmt"
	"github.com/google/uuid"
	"log"
	"math/rand"
	"time"

	"github.com/shopspring/decimal"
	"golang.org/x/crypto/bcrypt"

	"github.com/szucik/trade-helper/portfolio"
	"github.com/szucik/trade-helper/transaction"
)

type UsersService interface {
	SignUp(user User) (string, error)
	SignIn(credentials AuthCredentials) error
	GetUserByEmail(email string) (UserResponse, error)
	GetUsers() (UsersOut, error)
	AddPortfolio(in PortfolioIn) (string, error)
	AddTransaction(in TransactionIn) (string, error)
}

type Repository interface {
	// GetUserByEmail Dashboard(ctx context.Context) (Aggregate, error)
	// SignIn(ctx context.Context) (Aggregate, error)
	GetUserByEmail(userName string) (Aggregate, error)
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

type PortfolioIn struct {
	UserName string
	Name     string
}

type TransactionIn struct {
	UserName      string `json:"user_name"`
	PortfolioName string `json:"portfolio_name"`
	Symbol        string `json:"symbol"`
	Amount        string `json:"amount"`
	Quantity      string `json:"quantity"`
}

type UsersOut struct {
	Users []UserResponse
}

func (u Users) SignUp(user User) (string, error) {
	user.Created = time.Now()
	hash, err := hashPassword(user.Password)
	if err != nil {
		return "", fmt.Errorf("service.SignUp failed: %w", err)
	}

	user.Password = hash
	user.TokenHash = randStringBytes(15)

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

func (u Users) GetUserByEmail(email string) (UserResponse, error) {
	aggregate, err := u.Database.GetUserByEmail(email)
	if err != nil {
		return UserResponse{}, fmt.Errorf("database.GetUserByEmail failed: %w", err)
	}

	return transformToUserResponse(aggregate), nil
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

func (u Users) AddPortfolio(in PortfolioIn) (name string, _ error) {
	aggregate, err := u.Database.GetUserByEmail(in.UserName)
	if err != nil {
		return "", fmt.Errorf("database.GetUserByEmail failed: %w", err)
	}

	entity, err := portfolio.Portfolio{
		Name:            in.Name,
		TotalBalance:    decimal.NewFromFloat(0),
		TotalCost:       decimal.NewFromFloat(0),
		TotalProfitLoss: decimal.NewFromInt(0),
		ProfitLossDay:   decimal.NewFromInt(0),
		Created:         time.Now(),
	}.NewPortfolio()

	if err != nil {
		return "", fmt.Errorf("portfolio.NewPortfolio failed: %w", err)
	}

	err = aggregate.AddPortfolio(entity)
	if err != nil {
		return "", fmt.Errorf("aggregate.AddPortfolio failed: %w", err)
	}

	err = u.Database.SaveAggregate(aggregate)
	if err != nil {
		return "", fmt.Errorf("aggregate.SaveAggregate failed: %w", err)
	}

	return entity.Portfolio().Name, nil
}

func (u Users) AddTransaction(in TransactionIn) (string, error) {
	aggregate, err := u.Database.GetUserByEmail(in.UserName)
	if err != nil {
		return "", fmt.Errorf("service.AddTransaction: %w", err)
	}

	_, err = aggregate.FindPortfolio(in.PortfolioName)
	if err != nil {
		return "", fmt.Errorf("service.AddTransaction: %w", err)
	}

	quantity, err := decimal.NewFromString(in.Quantity)
	if err != nil {
		return "", fmt.Errorf("service.AddTransaction: %w", err)
	}

	price, err := decimal.NewFromString(in.Amount)
	if err != nil {
		return "", fmt.Errorf("service.AddTransaction: %w", err)
	}

	t, err := transaction.Transaction{
		ID:            uuid.New(),
		UserName:      in.UserName,
		PortfolioName: in.PortfolioName,
		Symbol:        in.Symbol,
		Created:       time.Now(),
		Quantity:      quantity,
		Price:         price,
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

func (u Users) SignIn(auth AuthCredentials) error {

	//u.Database.GetUserByEmail()

	err := comparePasswords(auth.Password, auth.Password)

	if err != nil {
		return err
	}
	return nil
}

func hashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 15)
	return string(hash), err
}

func comparePasswords(hashedPassword, providedPassword string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(providedPassword))
	if err != nil {
		return err
	}

	return nil
}

func transformToUserResponse(aggregate Aggregate) UserResponse {
	user := aggregate.User()

	var portfolios []portfolio.Portfolio
	for _, entity := range aggregate.Portfolios() {
		portfolios = append(portfolios, entity.Portfolio())
	}

	return UserResponse{
		Username:  user.Username,
		Email:     user.Email,
		Created:   user.Created,
		Portfolio: portfolios,
	}
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func randStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}
