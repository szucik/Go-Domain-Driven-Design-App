package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	gohandlers "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"

	"github.com/szucik/trade-helper/database/mongo"
	"github.com/szucik/trade-helper/user"
	"github.com/szucik/trade-helper/web"
	"github.com/szucik/trade-helper/web/handlers"
)

func main() {
	connectCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	logger := log.New(os.Stdout, "logger ", log.LstdFlags)
	database, err := mongo.NewDatabase(connectCtx)
	if err != nil {
		panic(err)
	}

	sessionSecret := os.Getenv("SESSION_SECRET")
	if sessionSecret == "" {
		log.Fatal("SESSION_SECRET environment variable is required")
	}
	store := sessions.NewCookieStore([]byte(sessionSecret))

	users := user.Users{
		Logger:       logger,
		Database:     &database,
		NewAggregate: user.User.NewAggregate,
	}

	sm := mux.NewRouter()
	sm.Use(web.MiddlewareIsAuth(store))

	postRouter := sm.Methods(http.MethodPost).Subrouter()
	postRouter.HandleFunc("/signup", handlers.SignUp(users))
	postRouter.HandleFunc("/signin", handlers.SignIn(users, store))
	postRouter.HandleFunc("/users/{username:[a-zA-Z0-9]+}/portfolio", handlers.AddPortfolio(users))
	postRouter.HandleFunc("/users/{username:[a-zA-Z0-9]+}/portfolio/{name:[a-zA-Z0-9]+}/transactions", handlers.AddTransaction(users))

	getRouter := sm.Methods(http.MethodGet).Subrouter()
	getRouter.HandleFunc("/users", handlers.GetUsers(users))
	getRouter.HandleFunc("/users/{username:[a-zA-Z0-9]+}", handlers.GetUser(users))
	getRouter.HandleFunc("/users/{username:[a-zA-Z0-9]+}/portfolio/{name:[a-zA-Z0-9]+}/transactions", handlers.GetTransactions(users))

	putRouter := sm.Methods(http.MethodPut).Subrouter()
	putRouter.HandleFunc("/users/{username:[a-zA-Z0-9]+}", handlers.UpdateUser(users))

	deleteRouter := sm.Methods(http.MethodDelete).Subrouter()
	deleteRouter.HandleFunc("/users/{username:[a-zA-Z0-9]+}", handlers.DeleteUser(users))

	ch := gohandlers.CORS(gohandlers.AllowedOrigins([]string{"localhost:3000"}))

	s := &http.Server{
		Addr:           ":9090",
		Handler:        ch(sm),
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	go func() {
		log.Println("Starting server on port 9090")
		if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	sig := <-c
	log.Println("Got signal:", sig)

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()
	s.Shutdown(shutdownCtx)
}
