package handlers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/szucik/trade-helper/apperrors"
	"github.com/szucik/trade-helper/user"
	"github.com/szucik/trade-helper/web/handlers"
)

var testStore = sessions.NewCookieStore([]byte("test-secret"))

type fakeService struct {
	signUpFn         func(ctx context.Context, u user.User) (string, error)
	signInFn         func(ctx context.Context, c user.AuthCredentials) (string, error)
	getUsersFn       func(ctx context.Context, p user.PaginationIn) (user.UsersOut, error)
	getUserByNameFn  func(ctx context.Context, name string) (user.UserResponse, error)
	addPortfolioFn   func(ctx context.Context, in user.PortfolioIn) (string, error)
	addTransactionFn func(ctx context.Context, in user.TransactionIn) (string, error)
	getTransactionsFn func(ctx context.Context, username, portfolioName string) (user.TransactionsOut, error)
	updateUserFn     func(ctx context.Context, username string, in user.UpdateUserIn) (string, error)
	deleteUserFn     func(ctx context.Context, username string) error
}

func (f *fakeService) SignUp(ctx context.Context, u user.User) (string, error) {
	return f.signUpFn(ctx, u)
}
func (f *fakeService) SignIn(ctx context.Context, c user.AuthCredentials) (string, error) {
	return f.signInFn(ctx, c)
}
func (f *fakeService) GetUsers(ctx context.Context, p user.PaginationIn) (user.UsersOut, error) {
	return f.getUsersFn(ctx, p)
}
func (f *fakeService) GetUserByEmail(_ context.Context, _ string) (user.UserResponse, error) {
	return user.UserResponse{}, nil
}
func (f *fakeService) GetUserByName(ctx context.Context, name string) (user.UserResponse, error) {
	return f.getUserByNameFn(ctx, name)
}
func (f *fakeService) AddPortfolio(ctx context.Context, in user.PortfolioIn) (string, error) {
	return f.addPortfolioFn(ctx, in)
}
func (f *fakeService) AddTransaction(ctx context.Context, in user.TransactionIn) (string, error) {
	return f.addTransactionFn(ctx, in)
}
func (f *fakeService) GetTransactions(ctx context.Context, username, portfolioName string) (user.TransactionsOut, error) {
	return f.getTransactionsFn(ctx, username, portfolioName)
}
func (f *fakeService) UpdateUser(ctx context.Context, username string, in user.UpdateUserIn) (string, error) {
	return f.updateUserFn(ctx, username, in)
}
func (f *fakeService) DeleteUser(ctx context.Context, username string) error {
	return f.deleteUserFn(ctx, username)
}

func newRequest(t *testing.T, method, path string, body any) *http.Request {
	t.Helper()
	var buf bytes.Buffer
	if body != nil {
		require.NoError(t, json.NewEncoder(&buf).Encode(body))
	}
	r, err := http.NewRequest(method, path, &buf)
	require.NoError(t, err)
	return r
}

func withVars(r *http.Request, vars map[string]string) *http.Request {
	return mux.SetURLVars(r, vars)
}

// --- SignUp ---

func TestSignUp_ReturnsUsername_WhenInputIsValid(t *testing.T) {
	svc := &fakeService{
		signUpFn: func(_ context.Context, u user.User) (string, error) { return u.Username, nil },
	}
	r := newRequest(t, http.MethodPost, "/signup", user.User{
		Username: "alice", Email: "alice@test.com", Password: "secret123",
	})
	rw := httptest.NewRecorder()
	handlers.SignUp(svc)(rw, r)

	assert.Equal(t, http.StatusOK, rw.Code)
	assert.Contains(t, rw.Body.String(), "alice")
}

func TestSignUp_Returns400_WhenRequiredFieldsMissing(t *testing.T) {
	svc := &fakeService{}
	r := newRequest(t, http.MethodPost, "/signup", user.User{Username: "alice"})
	rw := httptest.NewRecorder()
	handlers.SignUp(svc)(rw, r)

	assert.Equal(t, http.StatusBadRequest, rw.Code)
}

func TestSignUp_Returns400_WhenBodyIsInvalid(t *testing.T) {
	svc := &fakeService{}
	r, _ := http.NewRequest(http.MethodPost, "/signup", bytes.NewBufferString("not-json"))
	rw := httptest.NewRecorder()
	handlers.SignUp(svc)(rw, r)

	assert.Equal(t, http.StatusBadRequest, rw.Code)
}

func TestSignUp_Returns409_WhenUserAlreadyExists(t *testing.T) {
	svc := &fakeService{
		signUpFn: func(_ context.Context, _ user.User) (string, error) {
			return "", apperrors.Error("user already exists", "DuplicateUser", http.StatusConflict)
		},
	}
	r := newRequest(t, http.MethodPost, "/signup", user.User{
		Username: "alice", Email: "alice@test.com", Password: "secret123",
	})
	rw := httptest.NewRecorder()
	handlers.SignUp(svc)(rw, r)

	assert.Equal(t, http.StatusConflict, rw.Code)
}

// --- SignIn ---

func TestSignIn_ReturnsUsername_WhenCredentialsAreValid(t *testing.T) {
	svc := &fakeService{
		signInFn: func(_ context.Context, _ user.AuthCredentials) (string, error) { return "alice", nil },
	}
	r := newRequest(t, http.MethodPost, "/signin", user.AuthCredentials{
		Email: "alice@test.com", Password: "secret123",
	})
	rw := httptest.NewRecorder()
	handlers.SignIn(svc, testStore)(rw, r)

	assert.Equal(t, http.StatusOK, rw.Code)
	assert.Contains(t, rw.Body.String(), "alice")
}

func TestSignIn_Returns400_WhenCredentialsAreWrong(t *testing.T) {
	svc := &fakeService{
		signInFn: func(_ context.Context, _ user.AuthCredentials) (string, error) {
			return "", errors.New("invalid password")
		},
	}
	r := newRequest(t, http.MethodPost, "/signin", user.AuthCredentials{
		Email: "alice@test.com", Password: "wrong",
	})
	rw := httptest.NewRecorder()
	handlers.SignIn(svc, testStore)(rw, r)

	assert.Equal(t, http.StatusInternalServerError, rw.Code)
}

func TestSignIn_Returns400_WhenRequiredFieldsMissing(t *testing.T) {
	svc := &fakeService{}
	r := newRequest(t, http.MethodPost, "/signin", user.AuthCredentials{Email: "alice@test.com"})
	rw := httptest.NewRecorder()
	handlers.SignIn(svc, testStore)(rw, r)

	assert.Equal(t, http.StatusBadRequest, rw.Code)
}

// --- GetUsers ---

func TestGetUsers_ReturnsUserList(t *testing.T) {
	svc := &fakeService{
		getUsersFn: func(_ context.Context, _ user.PaginationIn) (user.UsersOut, error) {
			return user.UsersOut{
				Users: []user.UserResponse{{Username: "alice"}, {Username: "bob"}},
				Total: 2,
			}, nil
		},
	}
	r, _ := http.NewRequest(http.MethodGet, "/users", nil)
	rw := httptest.NewRecorder()
	handlers.GetUsers(svc)(rw, r)

	assert.Equal(t, http.StatusOK, rw.Code)
	var out user.UsersOut
	require.NoError(t, json.NewDecoder(rw.Body).Decode(&out))
	assert.Len(t, out.Users, 2)
}

func TestGetUsers_PassesPaginationParams(t *testing.T) {
	var capturedPagination user.PaginationIn
	svc := &fakeService{
		getUsersFn: func(_ context.Context, p user.PaginationIn) (user.UsersOut, error) {
			capturedPagination = p
			return user.UsersOut{}, nil
		},
	}
	r, _ := http.NewRequest(http.MethodGet, "/users?page=2&limit=10", nil)
	rw := httptest.NewRecorder()
	handlers.GetUsers(svc)(rw, r)

	assert.Equal(t, 2, capturedPagination.Page)
	assert.Equal(t, 10, capturedPagination.Limit)
}

// --- GetUser ---

func TestGetUser_ReturnsUser_WhenUsernameExists(t *testing.T) {
	svc := &fakeService{
		getUserByNameFn: func(_ context.Context, name string) (user.UserResponse, error) {
			return user.UserResponse{Username: name, Email: "alice@test.com"}, nil
		},
	}
	r := withVars(newRequest(t, http.MethodGet, "/users/alice", nil), map[string]string{"username": "alice"})
	rw := httptest.NewRecorder()
	handlers.GetUser(svc)(rw, r)

	assert.Equal(t, http.StatusOK, rw.Code)
	var out user.UserResponse
	require.NoError(t, json.NewDecoder(rw.Body).Decode(&out))
	assert.Equal(t, "alice", out.Username)
}

func TestGetUser_Returns404_WhenUsernameNotFound(t *testing.T) {
	svc := &fakeService{
		getUserByNameFn: func(_ context.Context, _ string) (user.UserResponse, error) {
			return user.UserResponse{}, apperrors.Error("user not found", "UserNotFound", http.StatusNotFound)
		},
	}
	r := withVars(newRequest(t, http.MethodGet, "/users/ghost", nil), map[string]string{"username": "ghost"})
	rw := httptest.NewRecorder()
	handlers.GetUser(svc)(rw, r)

	assert.Equal(t, http.StatusNotFound, rw.Code)
}

// --- AddPortfolio ---

func TestAddPortfolio_ReturnsName_WhenInputIsValid(t *testing.T) {
	svc := &fakeService{
		addPortfolioFn: func(_ context.Context, in user.PortfolioIn) (string, error) { return in.Name, nil },
	}
	r := withVars(
		newRequest(t, http.MethodPost, "/users/alice/portfolio", map[string]string{"Name": "tech"}),
		map[string]string{"username": "alice"},
	)
	rw := httptest.NewRecorder()
	handlers.AddPortfolio(svc)(rw, r)

	assert.Equal(t, http.StatusOK, rw.Code)
}

func TestAddPortfolio_Returns400_WhenNameMissing(t *testing.T) {
	svc := &fakeService{}
	r := withVars(
		newRequest(t, http.MethodPost, "/users/alice/portfolio", map[string]string{}),
		map[string]string{"username": "alice"},
	)
	rw := httptest.NewRecorder()
	handlers.AddPortfolio(svc)(rw, r)

	assert.Equal(t, http.StatusBadRequest, rw.Code)
}

// --- AddTransaction ---

func TestAddTransaction_ReturnsID_WhenInputIsValid(t *testing.T) {
	svc := &fakeService{
		addTransactionFn: func(_ context.Context, _ user.TransactionIn) (string, error) { return "txn-id-123", nil },
	}
	r := withVars(
		newRequest(t, http.MethodPost, "/users/alice/portfolio/tech/transactions", user.TransactionIn{
			Symbol: "AAPL", Amount: "150.00", Quantity: "2",
		}),
		map[string]string{"username": "alice", "name": "tech"},
	)
	rw := httptest.NewRecorder()
	handlers.AddTransaction(svc)(rw, r)

	assert.Equal(t, http.StatusOK, rw.Code)
	assert.Contains(t, rw.Body.String(), "txn-id-123")
}

func TestAddTransaction_Returns400_WhenRequiredFieldsMissing(t *testing.T) {
	svc := &fakeService{}
	r := withVars(
		newRequest(t, http.MethodPost, "/users/alice/portfolio/tech/transactions", user.TransactionIn{Symbol: "AAPL"}),
		map[string]string{"username": "alice", "name": "tech"},
	)
	rw := httptest.NewRecorder()
	handlers.AddTransaction(svc)(rw, r)

	assert.Equal(t, http.StatusBadRequest, rw.Code)
}

// --- GetTransactions ---

func TestGetTransactions_ReturnsList(t *testing.T) {
	svc := &fakeService{
		getTransactionsFn: func(_ context.Context, _, _ string) (user.TransactionsOut, error) {
			return user.TransactionsOut{Transactions: []user.TransactionResponse{{Symbol: "AAPL"}}}, nil
		},
	}
	r := withVars(
		newRequest(t, http.MethodGet, "/users/alice/portfolio/tech/transactions", nil),
		map[string]string{"username": "alice", "name": "tech"},
	)
	rw := httptest.NewRecorder()
	handlers.GetTransactions(svc)(rw, r)

	assert.Equal(t, http.StatusOK, rw.Code)
	var out user.TransactionsOut
	require.NoError(t, json.NewDecoder(rw.Body).Decode(&out))
	assert.Len(t, out.Transactions, 1)
}

// --- UpdateUser ---

func TestUpdateUser_ReturnsNewUsername_WhenInputIsValid(t *testing.T) {
	svc := &fakeService{
		updateUserFn: func(_ context.Context, _ string, in user.UpdateUserIn) (string, error) { return in.Username, nil },
	}
	r := withVars(
		newRequest(t, http.MethodPut, "/users/alice", user.UpdateUserIn{Username: "alice2"}),
		map[string]string{"username": "alice"},
	)
	rw := httptest.NewRecorder()
	handlers.UpdateUser(svc)(rw, r)

	assert.Equal(t, http.StatusOK, rw.Code)
	assert.Contains(t, rw.Body.String(), "alice2")
}

func TestUpdateUser_Returns400_WhenAllFieldsEmpty(t *testing.T) {
	svc := &fakeService{}
	r := withVars(
		newRequest(t, http.MethodPut, "/users/alice", user.UpdateUserIn{}),
		map[string]string{"username": "alice"},
	)
	rw := httptest.NewRecorder()
	handlers.UpdateUser(svc)(rw, r)

	assert.Equal(t, http.StatusBadRequest, rw.Code)
}

// --- DeleteUser ---

func TestDeleteUser_Returns204_WhenUserExists(t *testing.T) {
	svc := &fakeService{
		deleteUserFn: func(_ context.Context, _ string) error { return nil },
	}
	r := withVars(newRequest(t, http.MethodDelete, "/users/alice", nil), map[string]string{"username": "alice"})
	rw := httptest.NewRecorder()
	handlers.DeleteUser(svc)(rw, r)

	assert.Equal(t, http.StatusNoContent, rw.Code)
}

func TestDeleteUser_Returns404_WhenUserNotFound(t *testing.T) {
	svc := &fakeService{
		deleteUserFn: func(_ context.Context, _ string) error {
			return apperrors.Error("user not found", "UserNotFound", http.StatusNotFound)
		},
	}
	r := withVars(newRequest(t, http.MethodDelete, "/users/ghost", nil), map[string]string{"username": "ghost"})
	rw := httptest.NewRecorder()
	handlers.DeleteUser(svc)(rw, r)

	assert.Equal(t, http.StatusNotFound, rw.Code)
}
