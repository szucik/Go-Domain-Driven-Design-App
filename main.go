package main

import (
	"context"
	"github.com/szucik/trade-helper/app"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/szucik/trade-helper/database/mongo"

	"github.com/szucik/trade-helper/web/handlers"

	"github.com/szucik/trade-helper/user"

	gohandlers "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func main() {
	ctx, stop := context.WithTimeout(context.Background(), 30*time.Second)
	defer stop()

	logger := log.New(os.Stdout, "logger", log.LstdFlags)
	database, err := mongo.NewDatabase(ctx)
	if err != nil {
		panic(err)
	}

	users := user.Users{
		Logger:       logger,
		Database:     &database,
		NewAggregate: user.User.NewAggregate,
	}

	sm := mux.NewRouter()
	sm.Use(app.MiddlewareIsAuth)
	// Post
	postRouter := sm.Methods(http.MethodPost).Subrouter()
	postRouter.HandleFunc("/signup", handlers.SignUp(ctx, users))
	postRouter.HandleFunc("/signin", handlers.SignIn(ctx, users))
	postRouter.HandleFunc(
		"/users/{username:[a-z, A-Z, 0-9]+}/portfolio", handlers.AddPortfolio(ctx, users))
	postRouter.HandleFunc(
		"/users/{username:[a-z, A-Z, 0-9]+}/portfolio/{name:[a-z, A-Z, 0-9]+}/transactions",
		handlers.AddTransaction(ctx, users),
	)
	// postRouter.Use(users.MiddlewareUserValid)

	// Users
	getRouter := sm.Methods(http.MethodGet).Subrouter()
	getRouter.HandleFunc("/users", handlers.GetUsers(ctx, users))
	getRouter.HandleFunc("/users/{username:[a-z, A-Z, 0-9]+}", handlers.GetUser(ctx, users))
	//getRouter.HandleFunc("/users/{email:[a-z, A-Z, 0-9]+}", handlers.GetUser(ctx, users))

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
		log.Println("Starting server on port 9090")

		err = s.ListenAndServe()
		if err != nil {
			log.Fatal(err)
			os.Exit(1)
		}
	}()

	// trap sigterm or interupt and gracefully shutdown the server
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, os.Kill)

	// Block until a signal is received.
	sig := <-c
	log.Println("Got signal:", sig)

	// gracefully shutdown the server, waiting max 30 seconds for current operations to complete
	s.Shutdown(ctx)
}
