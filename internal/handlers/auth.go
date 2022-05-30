package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/golang-jwt/jwt"
	th_data "github.com/szucik/go-simple-rest-api/internal/data"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"time"
)

//var ErrUserPass = errors.New("")

//Users is a http.Handler
type Auth struct {
	l  *log.Logger
	db *th_data.Database
}

//NewUsers creates a users handler with the given logger
func NewAuth(l *log.Logger, db *th_data.Database) *Auth {
	return &Auth{l, db}
}

type keyAuth struct {
}

type JwtErrorMessage struct {
	message string
}

var (
	accessKey  = []byte("my_secret_key")
	refreshKey = []byte("dupa")
	domain     = "tradehelper.io"
)

type customClaims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

type jwtDetails struct {
	userName  string
	secretKey interface{}
	expiresAt time.Time
}

func newJwtToken(claims jwtDetails) (string, error) {
	c := customClaims{
		claims.userName,
		jwt.StandardClaims{
			ExpiresAt: claims.expiresAt.Unix(),
			Issuer:    domain,
		},
	}

	newJwt := jwt.NewWithClaims(jwt.SigningMethodHS256, c)

	jwtString, err := newJwt.SignedString(claims.secretKey)
	if err != nil {
		return "", err
	}

	return jwtString, nil
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func (a *Auth) Login(rw http.ResponseWriter, r *http.Request) {
	auth := r.Context().Value(keyAuth{}).(th_data.Auth)
	usr, err := a.db.Login(auth.Email)
	if err != nil {
		rw.WriteHeader(http.StatusUnauthorized)
		msg := fmt.Sprintf("%v", err)
		json.NewEncoder(rw).Encode(map[string]string{"message": msg})
		return
	}

	match := CheckPasswordHash(auth.Password, usr.Password)
	if !match {
		http.Error(rw, "Unable to unmarshal json", http.StatusNotFound)
	}

	tClaims := jwtDetails{
		userName:  usr.Username,
		secretKey: accessKey,
		expiresAt: time.Now().Add(5 * time.Minute),
	}
	access, err := newJwtToken(tClaims)
	if err != nil {
		fmt.Errorf("JWT error: %w", err)
	}

	http.SetCookie(rw, &http.Cookie{
		Name:    "Authorization",
		Value:   access,
		Expires: tClaims.expiresAt,
	})

	tClaims.secretKey = refreshKey
	tClaims.expiresAt = time.Now().Add(1 * time.Duration(24*time.Hour))
	refresh, err := newJwtToken(tClaims)
	if err != nil {
		fmt.Errorf("JWT error: %w", err)
	}

	http.SetCookie(rw, &http.Cookie{
		Name:     "Refresh",
		Value:    refresh,
		SameSite: 2,
		HttpOnly: true,
		Expires:  tClaims.expiresAt,
		Secure:   true,
	})
}

func (a *Auth) MiddlewareLoginValid(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		usr := &th_data.Auth{}
		err := th_data.FromJSON(usr, r.Body)
		if err != nil {
			a.l.Println("[ERROR] deserializing user", err)
			http.Error(rw, "Unable to unmarshal json", http.StatusBadRequest)
			return
		}
		err = th_data.Validate(usr)
		if err != nil {
			a.l.Println("[ERROR] validate user", err)
			http.Error(rw, "Unable validate user", http.StatusBadRequest)
			return
		}

		ctx := context.WithValue(r.Context(), keyAuth{}, *usr)
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
				return
			}
			rw.WriteHeader(http.StatusBadRequest)
			return
		}

		tknStr := c.Value
		claims := &customClaims{}

		tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
			return accessKey, nil
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
