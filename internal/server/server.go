package server

import (
	"context"

	"github.com/go-openapi/runtime/middleware"
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

func Run() {
	l := log.New(os.Stdout, "logger", log.LstdFlags)

	dc, err := database.Connection()
	if err != nil {
		panic(err.Error())
	}
	defer dc.Close()

	db := data.NewDatabase(dc)
	users := handlers.NewUsers(l, db)
	auth := handlers.NewAuth(l, db)

	sm := mux.NewRouter()

	//Sign Up
	signUpRouter := sm.Methods(http.MethodPost).Subrouter()
	signUpRouter.HandleFunc("/signup", auth.SignUp)
	signUpRouter.Use(users.MiddlewareUserValid)

	//Sign In
	signInRouter := sm.Methods(http.MethodPost).Subrouter()
	signInRouter.HandleFunc("/signin", auth.SignIn)
	signInRouter.Use(auth.MiddlewareLoginValid)

	//Users
	getRouter := sm.Methods(http.MethodGet).Subrouter()
	getRouter.HandleFunc("/users", users.GetUsers)
	getRouter.HandleFunc("/", users.Dashboard)
	getRouter.Use(auth.MiddlewareIsAuth)

	//postRouter := sm.Methods(http.MethodPost).Subrouter()

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
