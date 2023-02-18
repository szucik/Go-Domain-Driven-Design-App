package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/szucik/trade-helper/database/fake"
	"github.com/szucik/trade-helper/web/handlers"

	"github.com/szucik/trade-helper/user"

	gohandlers "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func main() {
	logger := log.New(os.Stdout, "logger", log.LstdFlags)
	database := fake.NewDatabase()

	users := user.Users{
		Logger:       logger,
		Database:     &database,
		NewAggregate: user.User.NewAggregate,
	}

	sm := mux.NewRouter()

	// Post
	signUpRouter := sm.Methods(http.MethodPost).Subrouter()
	signUpRouter.HandleFunc("/signup", handlers.SignUp(users))
	signUpRouter.HandleFunc(
		"/users/{username:[a-z, A-Z, 0-9]+}/portfolio", handlers.AddPortfolio(users))
	signUpRouter.HandleFunc(
		"/users/{username:[a-z, A-Z, 0-9]+}/portfolio/{name:[a-z, A-Z, 0-9]+}/transactions",
		handlers.AddTransaction(users),
	)
	// signUpRouter.Use(users.MiddlewareUserValid)

	// Users
	getRouter := sm.Methods(http.MethodGet).Subrouter()
	getRouter.HandleFunc("/users", handlers.GetUsers(users))
	getRouter.HandleFunc("/users/{username:[a-z, A-Z, 0-9]+}", handlers.GetUser(users))

	// putRouter := sm.Methods(http.MethodPut).Subrouter()
	// putRouter.HandleFunc("/users/{id:[0-9]+}", users.UpdateUser)
	// putRouter.Use(users.MiddlewareUserValid)

	// deleteRouter := sm.Methods(http.MethodDelete).Subrouter()
	// deleteRouter.HandleFunc("/users/{id:[0-9]+}", users.DeleteUser)

	// CORS
	ch := gohandlers.CORS(gohandlers.AllowedOrigins([]string{"localhost:3000"}))

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
			logger.Fatal(err)
		}
	}()

	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt)
	signal.Notify(sigChan, os.Kill)

	sig := <-sigChan
	logger.Println("Received terminate, graceful shutdown", sig)
	s.ListenAndServe()

	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	s.Shutdown(ctx)
}
