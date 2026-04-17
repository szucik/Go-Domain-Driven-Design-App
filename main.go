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

	// Public routes (no auth required)
	postRouter := sm.Methods(http.MethodPost).Subrouter()
	postRouter.HandleFunc("/signup", handlers.SignUp(users))
	postRouter.HandleFunc("/signin", handlers.SignIn(users, store))

	// Protected routes (require auth)
	protected := sm.PathPrefix("/").Subrouter()
	protected.Use(web.MiddlewareIsAuth(store))

	protectedPost := protected.Methods(http.MethodPost).Subrouter()
	protectedPost.HandleFunc("/users/{username:[a-zA-Z0-9]+}/portfolio", handlers.AddPortfolio(users))
	protectedPost.HandleFunc("/users/{username:[a-zA-Z0-9]+}/portfolio/{name:[a-zA-Z0-9]+}/transactions", handlers.AddTransaction(users))

	protectedGet := protected.Methods(http.MethodGet).Subrouter()
	protectedGet.HandleFunc("/users", handlers.GetUsers(users))
	protectedGet.HandleFunc("/users/{username:[a-zA-Z0-9]+}", handlers.GetUser(users))
	protectedGet.HandleFunc("/users/{username:[a-zA-Z0-9]+}/portfolio/{name:[a-zA-Z0-9]+}/transactions", handlers.GetTransactions(users))

	protectedPut := protected.Methods(http.MethodPut).Subrouter()
	protectedPut.HandleFunc("/users/{username:[a-zA-Z0-9]+}", handlers.UpdateUser(users))

	protectedDelete := protected.Methods(http.MethodDelete).Subrouter()
	protectedDelete.HandleFunc("/users/{username:[a-zA-Z0-9]+}", handlers.DeleteUser(users))

	// Serve static files from frontend directory (no auth required)
	sm.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir("./frontend/"))))

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
