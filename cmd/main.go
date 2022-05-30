package main

import (
	"context"
	"github.com/go-openapi/runtime/middleware"
	_ "github.com/go-sql-driver/mysql"
	gohandlers "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/szucik/go-simple-rest-api/internal/data"
	"github.com/szucik/go-simple-rest-api/internal/database"
	"github.com/szucik/go-simple-rest-api/internal/handlers"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {
	l := log.New(os.Stdout, "logger", log.LstdFlags)
	dc, err := database.Connect()
	if err != nil {
		panic(err.Error())
	}
	defer dc.Close()
	db := data.NewDatabase(dc)
	users := handlers.NewUsers(l, db)
	auth := handlers.NewAuth(l, db)

	sm := mux.NewRouter()

	//Auth
	authRouter := sm.Methods(http.MethodPost).Subrouter()
	authRouter.HandleFunc("/login", auth.Login)
	authRouter.Use(auth.MiddlewareLoginValid)

	//Users
	getRouter := sm.Methods(http.MethodGet).Subrouter()
	getRouter.HandleFunc("/users", users.GetUsers)
	getRouter.HandleFunc("/", users.Dashboard)
	getRouter.Use(auth.MiddlewareAuth)

	postRouter := sm.Methods(http.MethodPost).Subrouter()
	postRouter.HandleFunc("/users", users.AddUser)
	postRouter.Use(users.MiddlewareUserValid)
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
	l.Println("Recived terminate, graceful shutdown", sig)
	s.ListenAndServe()

	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	s.Shutdown(ctx)
}
