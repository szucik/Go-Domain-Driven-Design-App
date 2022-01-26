package handlers

import (
	"context"
	"fmt"
	"github.com/golang-jwt/jwt"
	data2 "github.com/szucik/go-simple-rest-api/internal/data"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"time"
)

//Users is a http.Handler
type Auth struct {
	l  *log.Logger
	db *data2.Database
}

//NewUsers creates a users handler with the given logger
func NewAuth(l *log.Logger, db *data2.Database) *Auth {
	return &Auth{l, db}
}

type KeyAuth struct {
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

var hmacSampleSecret []byte

func (a *Auth) Login(rw http.ResponseWriter, r *http.Request) {
	auth := r.Context().Value(KeyAuth{}).(data2.Auth)
	usr, err := a.db.Login(auth.Email)
	if err != nil {
		fmt.Println(err)
	}
	match := CheckPasswordHash(auth.Password, usr.Password)
	if !match {
		http.Error(rw, "Unable to unmarshal json", http.StatusNotFound)
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"foo": "marcin-token",
		"nbf": time.Now().Unix(),
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString(hmacSampleSecret)
	fmt.Print(tokenString)
}

func (a *Auth) MiddlewareAuthValid(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		usr := &data2.Auth{}
		err := data2.FromJSON(usr, r.Body)
		if err != nil {
			a.l.Println("[ERROR] deserializing user", err)
			http.Error(rw, "Unable to unmarshal json", http.StatusBadRequest)
			return
		}
		err = data2.Validate(usr)
		if err != nil {
			a.l.Println("[ERROR] validate user", err)
			http.Error(rw, "Unable validate user", http.StatusBadRequest)
			return
		}

		ctx := context.WithValue(r.Context(), KeyAuth{}, *usr)

		r = r.WithContext(ctx)
		next.ServeHTTP(rw, r)
	})
}
