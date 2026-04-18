package user

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/sessions"
	"github.com/shopspring/decimal"
	"golang.org/x/crypto/bcrypt"

	"github.com/szucik/trade-helper/apperrors"
	"github.com/szucik/trade-helper/portfolio"
	"github.com/szucik/trade-helper/transaction"
)

type UsersService interface {
	SignUp(ctx context.Context, user User) (string, error)
	SignIn(ctx context.Context, credentials AuthCredentials) (username string, err error)
	GetUserByEmail(ctx context.Context, email string) (UserResponse, error)
	GetUserByName(ctx context.Context, name string) (UserResponse, error)
	GetUsers(ctx context.Context, p PaginationIn) (UsersOut, error)
	GetTransactions(ctx context.Context, username, portfolioName string) (TransactionsOut, error)
	AddPortfolio(ctx context.Context, in PortfolioIn) (string, error)
	AddTransaction(ctx context.Context, in TransactionIn) (string, error)
	UpdateUser(ctx context.Context, currentUsername string, in UpdateUserIn) (string, error)
	DeleteUser(ctx context.Context, username string) error
}

type Repository interface {
	GetUserByEmail(ctx context.Context, email string) (Aggregate, error)
	GetUserByName(ctx context.Context, userName string) (Aggregate, error)
	GetUsers(ctx context.Context, p PaginationIn) ([]Aggregate, error)
	SignUp(ctx context.Context, aggregate Aggregate) (string, error)
	SaveAggregate(ctx context.Context, aggregate Aggregate) error
	AddTransaction(ctx context.Context, t transaction.ValueObject) (string, error)
	GetTransactions(ctx context.Context, username, portfolioName string) ([]transaction.ValueObject, error)
	UpdateUser(ctx context.Context, currentUsername string, aggregate Aggregate) error
	DeleteUser(ctx context.Context, username string) error
}

type Users struct {
	Logger       *log.Logger
	Database     Repository
	Sessions     sessions.CookieStore
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
	Type          string `json:"type"`
}

type UsersOut struct {
	Users []UserResponse `json:"users"`
	Total int            `json:"total"`
	Page  int            `json:"page"`
	Limit int            `json:"limit"`
}

func (u Users) SignUp(ctx context.Context, user User) (string, error) {
	user.Created = time.Now()
	user.TokenHash = randStringBytes(15)

	aggregate, err := u.NewAggregate(user)
	if err != nil {
		return "", err
	}

	if err = aggregate.hashPassword(); err != nil {
		return "", err
	}

	id, err := u.Database.SignUp(ctx, aggregate)
	if err != nil {
		u.Logger.Printf("SignUp failed for %s: %v", user.Email, err)
		return "", fmt.Errorf("database.SignUp failed: %w", err)
	}

	u.Logger.Printf("SignUp: new user %s", user.Username)
	return id, nil
}

func (u Users) GetUserByName(ctx context.Context, username string) (UserResponse, error) {
	aggregate, err := u.Database.GetUserByName(ctx, username)
	if err != nil {
		return UserResponse{}, fmt.Errorf("database.GetUserByName failed: %w", err)
	}
	return transformToUserResponse(aggregate), nil
}

func (u Users) GetUserByEmail(ctx context.Context, email string) (UserResponse, error) {
	aggregate, err := u.Database.GetUserByEmail(ctx, email)
	if err != nil {
		return UserResponse{}, fmt.Errorf("database.GetUserByEmail failed: %w", err)
	}
	return transformToUserResponse(aggregate), nil
}

func (u Users) GetUsers(ctx context.Context, p PaginationIn) (UsersOut, error) {
	aggregates, err := u.Database.GetUsers(ctx, p)
	if err != nil {
		return UsersOut{}, fmt.Errorf("database.GetUsers failed: %w", err)
	}

	out := UsersOut{Page: p.Page, Limit: p.Limit, Total: len(aggregates)}
	for _, a := range aggregates {
		out.Users = append(out.Users, transformToUserResponse(a))
	}
	return out, nil
}

func (u Users) GetTransactions(ctx context.Context, username, portfolioName string) (TransactionsOut, error) {
	vos, err := u.Database.GetTransactions(ctx, username, portfolioName)
	if err != nil {
		return TransactionsOut{}, fmt.Errorf("database.GetTransactions failed: %w", err)
	}

	var out TransactionsOut
	for _, vo := range vos {
		t := vo.Transaction()
		out.Transactions = append(out.Transactions, TransactionResponse{
			ID:       t.ID.String(),
			Symbol:   t.Symbol,
			Type:     t.Type,
			Quantity: t.Quantity,
			Price:    t.Price,
			Created:  t.Created,
		})
	}
	return out, nil
}

func (u Users) AddPortfolio(ctx context.Context, in PortfolioIn) (name string, _ error) {
	aggregate, err := u.Database.GetUserByName(ctx, in.UserName)
	if err != nil {
		return "", fmt.Errorf("database.GetUserByName failed: %w", err)
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

	if err = aggregate.AddPortfolio(entity); err != nil {
		return "", fmt.Errorf("aggregate.AddPortfolio failed: %w", err)
	}

	if err = u.Database.SaveAggregate(ctx, aggregate); err != nil {
		u.Logger.Printf("AddPortfolio SaveAggregate failed for %s: %v", in.UserName, err)
		return "", fmt.Errorf("aggregate.SaveAggregate failed: %w", err)
	}

	u.Logger.Printf("AddPortfolio: user %s added portfolio %s", in.UserName, in.Name)
	return entity.Portfolio().Name, nil
}

func (u Users) AddTransaction(ctx context.Context, in TransactionIn) (string, error) {
	aggregate, err := u.Database.GetUserByName(ctx, in.UserName)
	if err != nil {
		return "", fmt.Errorf("service.AddTransaction: %w", err)
	}

	if _, err = aggregate.FindPortfolio(in.PortfolioName); err != nil {
		return "", fmt.Errorf("service.AddTransaction: %w", err)
	}

	quantity, err := decimal.NewFromString(in.Quantity)
	if err != nil {
		return "", apperrors.Error("invalid quantity format", "BadRequest", http.StatusBadRequest)
	}

	price, err := decimal.NewFromString(in.Amount)
	if err != nil {
		return "", apperrors.Error("invalid amount format", "BadRequest", http.StatusBadRequest)
	}

	txType := parseTransactionType(in.Type)

	t, err := transaction.Transaction{
		ID:            uuid.New(),
		UserName:      in.UserName,
		PortfolioName: in.PortfolioName,
		Symbol:        in.Symbol,
		Type:          txType,
		Created:       time.Now(),
		Quantity:      quantity,
		Price:         price,
	}.NewTransaction()
	if err != nil {
		return "", fmt.Errorf("NewTransaction error: %w", err)
	}

	id, err := u.Database.AddTransaction(ctx, t)
	if err != nil {
		u.Logger.Printf("AddTransaction DB failed for %s/%s: %v", in.UserName, in.PortfolioName, err)
		return "", fmt.Errorf("Database.AddTransaction: %w", err)
	}

	costDelta := quantity.Mul(price)
	if txType == transaction.Sell {
		costDelta = costDelta.Neg()
	}

	if err := aggregate.UpdatePortfolioTotalCost(in.PortfolioName, costDelta); err != nil {
		return "", fmt.Errorf("service.AddTransaction: %w", err)
	}

	if err := u.Database.SaveAggregate(ctx, aggregate); err != nil {
		u.Logger.Printf("AddTransaction SaveAggregate failed for %s: %v", in.UserName, err)
		return "", fmt.Errorf("service.AddTransaction SaveAggregate: %w", err)
	}

	u.Logger.Printf("AddTransaction: %s %s %s x%s", txType, in.Symbol, in.PortfolioName, in.Quantity)
	return id, nil
}

func (u Users) SignIn(ctx context.Context, auth AuthCredentials) (string, error) {
	const invalidCredentials = "invalid credentials"

	aggregate, err := u.Database.GetUserByEmail(ctx, auth.Email)
	if err != nil {
		if !apperrors.IsNotFound(err) {
			u.Logger.Printf("SignIn: unexpected database error: %v", err)
		}
		return "", apperrors.Error(invalidCredentials, "Unauthorized", http.StatusUnauthorized)
	}

	usr := aggregate.User()
	if err = compareHashAndPassword(usr.Password, auth.Password); err != nil {
		return "", apperrors.Error(invalidCredentials, "Unauthorized", http.StatusUnauthorized)
	}

	u.Logger.Printf("SignIn: user %s authenticated", usr.Username)
	return usr.Username, nil
}

func (u Users) UpdateUser(ctx context.Context, currentUsername string, in UpdateUserIn) (string, error) {
	aggregate, err := u.Database.GetUserByName(ctx, currentUsername)
	if err != nil {
		return "", fmt.Errorf("service.UpdateUser: %w", err)
	}

	current := aggregate.User()
	if in.Username != "" {
		current.Username = in.Username
	}
	if in.Email != "" {
		current.Email = in.Email
	}
	if in.Password != "" {
		current.Password = in.Password
	}
	current.Updated = time.Now()

	updated, err := u.NewAggregate(current)
	if err != nil {
		return "", fmt.Errorf("service.UpdateUser: %w", err)
	}

	if in.Password != "" {
		if err := updated.hashPassword(); err != nil {
			return "", fmt.Errorf("service.UpdateUser: %w", err)
		}
	}

	if err := u.Database.UpdateUser(ctx, currentUsername, updated); err != nil {
		u.Logger.Printf("UpdateUser failed for %s: %v", currentUsername, err)
		return "", fmt.Errorf("service.UpdateUser: %w", err)
	}

	u.Logger.Printf("UpdateUser: %s updated", currentUsername)
	return updated.User().Username, nil
}

func (u Users) DeleteUser(ctx context.Context, username string) error {
	if err := u.Database.DeleteUser(ctx, username); err != nil {
		u.Logger.Printf("DeleteUser failed for %s: %v", username, err)
		return fmt.Errorf("service.DeleteUser: %w", err)
	}
	u.Logger.Printf("DeleteUser: user %s deleted", username)
	return nil
}

func parseTransactionType(s string) transaction.TransactionType {
	if s == "sell" {
		return transaction.Sell
	}
	return transaction.Buy
}

func hashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 15)
	return string(hash), err
}

func compareHashAndPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func transformToUserResponse(aggregate Aggregate) UserResponse {
	usr := aggregate.User()

	var portfolios []portfolio.Portfolio
	for _, entity := range aggregate.Portfolios() {
		portfolios = append(portfolios, entity.Portfolio())
	}

	return UserResponse{
		Username:  usr.Username,
		Email:     usr.Email,
		Created:   usr.Created,
		Portfolio: portfolios,
	}
}
