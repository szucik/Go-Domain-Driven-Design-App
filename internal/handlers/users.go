// Package classification of UserKey API
//
// Documentation for UserKey API
//
//	Schemes: http
//	BasePath: /
//	Version: 1.0.0
//
//	Consumes:
//	- application/json
//
//	Produces:
//	- application/json
//
// swagger:meta
package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"

	th_data "github.com/szucik/go-simple-rest-api/internal/dao"
)

var (
	ErrHash = errors.New("Problem with hashing your password")

	MsgUserAlreadyExists = "UserKey already exists with the given email"
)

type UserKey struct{}

//A list of users returns in the response
// swagger:response usersResponse
type usersResponseWrapper struct {
	//All users in the system
	//in: body
	Body []th_data.User
}

type GenericResponse struct {
	Status  bool        `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"dao"`
}

// swagger:response noContent
type usersNoContent struct {
}

//swagger:parameters deleteUser
type usersIDParameterWrapper struct {
	//The id of user delete from database
	//	in:path
	//	required: true
	ID int `json:"id"`
}

//Users is a http.Handler
type Users struct {
	l  *log.Logger
	db *th_data.Dao
}

func NewUsers(l *log.Logger, db *th_data.Dao) *Users {
	return &Users{l, db}
}

// swagger:route GET /users listUsers
// Return a list of users from the database
// responses:
//	200: usersResponse

// GetUsers returns Users from dao store
func (u *Users) GetUsers(rw http.ResponseWriter, r *http.Request) {
	u.l.Println("Handle GET Users")
	rw.Header().Add("Content-Type", "application/json")

	users, err := u.db.GetUsers()
	if err != nil {
		fmt.Errorf("%s", err)
	}

	j, _ := json.Marshal(users)

	if err != nil {
		http.Error(rw, "Unable to marshal json", http.StatusInternalServerError)
	}
	rw.Write(j)
}

func (u *Users) GetUser(rw http.ResponseWriter, r *http.Request) {
	u.l.Println("Handle GET User")
	rw.Header().Add("Content-Type", "application/json")
	rc := r.Context().Value(UserKey{}).(th_data.User)
	user, err := u.db.GetUser(rc.ID)
	if err != nil {
		fmt.Errorf("%s", err)
	}

	j, _ := json.Marshal(user)

	if err != nil {
		http.Error(rw, "Unable to marshal json", http.StatusInternalServerError)
	}
	rw.Write(j)
}

// swagger:route DELETE /users/{id} users deleteUser
// Return a list of users from the database
// responses:
//	201: noContent

// DeleteUser deletes a user from DB

func (u *Users) DeleteUser(rw http.ResponseWriter, r *http.Request) {}

func (u *Users) UpdateUser(rw http.ResponseWriter, r *http.Request) {}

func (u *Users) Dashboard(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	//TODO Remember about Authorization cookies in client
	rw.Write([]byte("Welcom in dashboard"))
}

type JwtErrorMessage struct {
	message string
}

var (
	accessKey  = []byte("accessKey")
	refreshKey = []byte("refreshKey")
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

func (u *Users) SignIn(rw http.ResponseWriter, r *http.Request) {
	auth := r.Context().Value(UserKey{}).(th_data.AuthCredentials)
	usr, err := u.db.Login(auth.Email)
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

func (u *Users) SignUp(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	rc := r.Context().Value(UserKey{}).(th_data.User)

	hp, err := HashPassword(rc.Password)
	if err != nil {
		fmt.Errorf("%s", ErrHash)
	}

	salt := RandomString(15)

	user := &th_data.User{
		Username:  rc.Username,
		Email:     rc.Email,
		Password:  hp,
		TokenHash: salt,
	}

	_, err = u.db.AddUser(user)
	if err != nil {
		message := fmt.Sprintf("Error message: %v", err)
		u.l.Print(message)
		u.db.ToJSON(&GenericResponse{Status: false, Message: MsgUserAlreadyExists}, rw)
		return
	}
	u.l.Print("UserKey created successfully")

	u.db.ToJSON(&GenericResponse{Status: true, Message: "user created successfully"}, rw)
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func RandomString(n int) string {
	var output strings.Builder
	for i := 0; i < n; i++ {
		random := rand.Intn(len(letterBytes))
		randomChar := letterBytes[random]
		output.WriteString(string(randomChar))
	}
	return output.String()
}

func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 15)
	return string(hash), err
}
