package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/szucik/trade-helper/database/fake"
	"github.com/szucik/trade-helper/web/handlers"

	"github.com/szucik/trade-helper/user"

	_ "github.com/go-sql-driver/mysql"

	gohandlers "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func main() {
	logger := log.New(os.Stdout, "logger", log.LstdFlags)

	// config, err := app.GetConfiguration()
	// if err != nil {
	//	panic("Loading config failed: " + err.Error())
	// }

	database := fake.NewDatabase()
	// portfolios := portfolio.Portfolios{
	//	Database:     database,
	//	NewAggregate: portfolio.Portfolio.NewAggregate,
	// }

	users := user.Users{
		Logger:       logger,
		Database:     &database,
		NewAggregate: user.User.NewAggregate,
	}

	sm := mux.NewRouter()

	fmt.Println(users)
	// SignUp
	signUpRouter := sm.Methods(http.MethodPost).Subrouter()
	signUpRouter.HandleFunc("/signup", handlers.SignUp(users))
	// signUpRouter.Use(users.MiddlewareUserValid)

	// SignIn
	// signInRouter := sm.Methods(http.MethodPost).Subrouter()
	// signInRouter.HandleFunc("/signin", users.SignIn)
	// signInRouter.Use(users.MiddlewareLoginValid)

	// Users
	getRouter := sm.Methods(http.MethodGet).Subrouter()
	getRouter.HandleFunc("/users", handlers.GetUsers(users))
	getRouter.HandleFunc("/users/{username:[a-z, A-Z, 0-9]+}", handlers.GetUser(users))
	// getRouter.HandleFunc("/", users.Dashboard)
	// getRouter.Use(users.MiddlewareIsAuth)

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
