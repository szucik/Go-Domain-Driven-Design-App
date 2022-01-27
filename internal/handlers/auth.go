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

var jwtKey = []byte("my_secret_key")

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
	type MyCustomClaims struct {
		Foo string `json:"foo"`
		jwt.StandardClaims
	}
	// Create the Claims
	expirationTime := time.Now().Add(5 * time.Minute)
	claims := MyCustomClaims{
		"jwt",
		jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
			Issuer:    "test",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		// If there is an error in creating the JWT return an internal server error
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	http.SetCookie(rw, &http.Cookie{
		Name:    "Authorization",
		Value:   tokenString,
		Expires: expirationTime,
	})
}

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

func (a *Auth) MiddlewareLoginValid(next http.Handler) http.Handler {
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
		fmt.Println(ctx)
		r = r.WithContext(ctx)

		next.ServeHTTP(rw, r)
	})
}

func (a *Auth) MiddlewareAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		c, err := r.Cookie("Authorization")
		if err != nil {
			if err == http.ErrNoCookie {
				rw.WriteHeader(http.StatusUnauthorized)
				http.Redirect(rw, r, "/login", http.StatusFound)
				return
			}
			// For any other type of error, return a bad request status
			rw.WriteHeader(http.StatusBadRequest)
			//http.Redirect(rw, r, "/login", http.StatusFound) ??
			return
		}

		tknStr := c.Value
		claims := &Claims{}
		tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})
		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				rw.WriteHeader(http.StatusUnauthorized)
				return
			}
			rw.WriteHeader(http.StatusBadRequest)
			return
		}
		if !tkn.Valid {
			rw.WriteHeader(http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(rw, r)
	})
}
