package main

import (
	"context"
	"github.com/zlyjoker102/simple-rest-api/handlers"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {
	l := log.New(os.Stdout, "logger", log.LstdFlags)
	hh := handlers.NewArticle(l)
	//hu := handlers.NewUser(l)
	hc := handlers.NewCurrency(l)

	sm := http.NewServeMux()
	sm.Handle("/", hh)
	//sm.Handle("/user", hu)
	sm.Handle("/currency", hc)

	s := &http.Server{
		Addr:           ":8080",
		Handler:        sm,
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
	tc, _ := context.WithTimeout(context.Background(), 30*time.Second)
	s.Shutdown(tc)
}
