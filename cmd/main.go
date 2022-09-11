package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/go-openapi/runtime/middleware"

	_ "github.com/go-sql-driver/mysql"

	gohandlers "github.com/gorilla/handlers"
	"github.com/gorilla/mux"

	"github.com/szucik/go-simple-rest-api/internal/dao"
	"github.com/szucik/go-simple-rest-api/internal/database"
	"github.com/szucik/go-simple-rest-api/internal/handlers"
)

func main() {
	l := log.New(os.Stdout, "logger", log.LstdFlags)

	dc, err := database.Connection()
	if err != nil {
		panic(err.Error())
	}
	defer dc.Close()

	db := dao.NewDatabase(dc)
	users := handlers.NewUsers(l, db)
	transactions := handlers.NewTransactions(l, db)
	portfolio := handlers.NewPortfolios(l, db)
	sm := mux.NewRouter()

	// SignUp
	signUpRouter := sm.Methods(http.MethodPost).Subrouter()
	signUpRouter.HandleFunc("/signup", users.SignUp)
	signUpRouter.Use(users.MiddlewareUserValid)

	// SignIn
	signInRouter := sm.Methods(http.MethodPost).Subrouter()
	signInRouter.HandleFunc("/signin", users.SignIn)
	signInRouter.Use(users.MiddlewareLoginValid)

	// Users
	getRouter := sm.Methods(http.MethodGet).Subrouter()
	getRouter.HandleFunc("/users", users.GetUsers)
	getRouter.HandleFunc("/users/{id:[0-9]+}", users.GetUsers)
	getRouter.HandleFunc("/", users.Dashboard)
	getRouter.Use(users.MiddlewareIsAuth)

	// Transactions
	coinsRouter := sm.Methods(http.MethodPost).Subrouter()
	coinsRouter.HandleFunc("/users/transactions", transactions.AddTransaction)
	coinsRouter.Use(transactions.MiddlewareTransactionValid)

	// Portfolio
	portfolioRouter := sm.Methods(http.MethodPost).Subrouter()
	portfolioRouter.HandleFunc("/users/portfolios", portfolio.AddPortfolio)
	portfolioRouter.Use(portfolio.MiddlewarePortfolioValid)

	putRouter := sm.Methods(http.MethodPut).Subrouter()
	putRouter.HandleFunc("/users/{id:[0-9]+}", users.UpdateUser)
	putRouter.Use(users.MiddlewareUserValid)

	deleteRouter := sm.Methods(http.MethodDelete).Subrouter()
	deleteRouter.HandleFunc("/users/{id:[0-9]+}", users.DeleteUser)

	// handler for documentation
	opts := middleware.RedocOpts{SpecURL: "/swagger.yaml"}
	sh := middleware.Redoc(opts, nil)

	// CORS
	ch := gohandlers.CORS(gohandlers.AllowedOrigins([]string{"localhost:3000"}))
	getRouter.Handle("/docs", sh)
	getRouter.Handle("/swagger.yaml", http.FileServer(http.Dir("./")))
	s := &http.Server{
		Addr:           ":9090",
		Handler:        ch(sm),
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	go func() {
		err := s.ListenAndServe()
		if err != nil {
			l.Fatal(err)
		}
	}()

	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt)
	signal.Notify(sigChan, os.Kill)

	sig := <-sigChan
	l.Println("Received terminate, graceful shutdown", sig)
	s.ListenAndServe()

	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	s.Shutdown(ctx)
}
