package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/golang-jwt/jwt"
	th_data "github.com/szucik/go-simple-rest-api/internal/data"
	"github.com/szucik/go-simple-rest-api/internal/utils"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"time"
)

//var ErrUserPass = errors.New("")

type Authenticate struct {
	l  *log.Logger
	db *th_data.Database
}

func NewAuth(l *log.Logger, db *th_data.Database) *Authenticate {
	return &Authenticate{l, db}
}

type keyAuth struct{}

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

func (a *Authenticate) SignIn(rw http.ResponseWriter, r *http.Request) {
	auth := r.Context().Value(keyAuth{}).(th_data.AuthCredentials)
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

func (a *Authenticate) SignUp(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	rc := r.Context().Value(UserKey{}).(th_data.User)

	hp, err := utils.HashPassword(rc.Password)
	if err != nil {
		fmt.Errorf("%s", ErrHash)
	}

	salt := utils.RandomString(15)

	user := &th_data.User{
		Username:  rc.Username,
		Email:     rc.Email,
		Password:  hp,
		TokenHash: salt,
	}

	_, err = a.db.AddUser(user)
	if err != nil {
		message := fmt.Sprintf("Error message: %v", err)
		//jM, _ := json.Marshal(message)
		a.l.Print(message)
		th_data.ToJSON(&GenericResponse{Status: false, Message: MsgUserAlreadyExists}, rw)
		return
	}
	a.l.Print("UserKey created successfully")

	//rw.WriteHeader(http.StatusCreated)
	th_data.ToJSON(&GenericResponse{Status: true, Message: "user created successfully"}, rw)
}
