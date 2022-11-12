package user

import (
	"context"
	"errors"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrHash              = errors.New("Problem with hashing your password")
	MsgUserAlreadyExists = "UserKey already exists with the given email"
)

type Database interface {
	GetUsers(ctx context.Context) (Aggregate, error)
	GetUser(ctx context.Context) (Aggregate, error)
	UpdateUser(ctx context.Context) (Aggregate, error)
	Dashboard(ctx context.Context) (Aggregate, error)
	SignIn(ctx context.Context) (Aggregate, error)
	SignUp(ctx context.Context) (Aggregate, error)
}

type Users struct {
	Logger       *log.Logger
	Database     Database
	NewAggregate func(User) (Aggregate, error)
}

type UserKey struct{}

type GenericResponse struct {
	Status  bool        `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"dao"`
}

func (u *Users) GetUsers(rw http.ResponseWriter, r *http.Request) {
	//u.logger.Println("Handle GET Users")
	//rw.Header().Add("Content-Type", "application/json")
	//
	//users, err := u.Database.GetUsers()
	//if err != nil {
	//	fmt.Errorf("%s", err)
	//}
	//
	//j, _ := json.Marshal(users)
	//
	//if err != nil {
	//	http.Error(rw, "Unable to marshal json", http.StatusInternalServerError)
	//}
	//rw.Write(j)
}

func (u *Users) GetUser(rw http.ResponseWriter, r *http.Request) {
	//u.logger.Println("Handle GET User")
	//rw.Header().Add("Content-Type", "application/json")
	//rc := r.Context().Value(UserKey{}).(dao.User)
	//user, err := u.Database.GetUser(rc.ID)
	//if err != nil {
	//	fmt.Errorf("%s", err)
	//}
	//
	//j, _ := json.Marshal(user)
	//
	//if err != nil {
	//	http.Error(rw, "Unable to marshal json", http.StatusInternalServerError)
	//}
	//rw.Write(j)
}
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
	//auth := r.Context().Value(UserKey{}).(dao.AuthCredentials)
	//usr, err := u.Database.Login(auth.Email)
	//if err != nil {
	//	rw.WriteHeader(http.StatusUnauthorized)
	//	msg := fmt.Sprintf("%v", err)
	//	json.NewEncoder(rw).Encode(map[string]string{"message": msg})
	//	return
	//}
	//
	//match := CheckPasswordHash(auth.Password, usr.Password)
	//if !match {
	//	http.Error(rw, "Unable to unmarshal json", http.StatusNotFound)
	//}
	//
	//tClaims := jwtDetails{
	//	userName:  usr.Username,
	//	secretKey: accessKey,
	//	expiresAt: time.Now().Add(5 * time.Minute),
	//}
	//access, err := newJwtToken(tClaims)
	//if err != nil {
	//	fmt.Errorf("JWT error: %w", err)
	//}
	//
	//http.SetCookie(rw, &http.Cookie{
	//	Name:    "Authorization",
	//	Value:   access,
	//	Expires: tClaims.expiresAt,
	//})
	//
	//tClaims.secretKey = refreshKey
	//tClaims.expiresAt = time.Now().Add(1 * time.Duration(24*time.Hour))
	//refresh, err := newJwtToken(tClaims)
	//if err != nil {
	//	fmt.Errorf("JWT error: %w", err)
	//}
	//
	//http.SetCookie(rw, &http.Cookie{
	//	Name:     "Refresh",
	//	Value:    refresh,
	//	SameSite: 2,
	//	HttpOnly: true,
	//	Expires:  tClaims.expiresAt,
	//	Secure:   true,
	//})
}

func (u *Users) SignUp(rw http.ResponseWriter, r *http.Request) {
	//rw.Header().Set("Content-Type", "application/json")
	//rc := r.Context().Value(UserKey{}).(dao.User)
	//
	//hp, err := HashPassword(rc.Password)
	//if err != nil {
	//	fmt.Errorf("%s", ErrHash)
	//}
	//
	//salt := RandomString(15)
	//
	//user := &dao.User{
	//	Username:  rc.Username,
	//	Email:     rc.Email,
	//	Password:  hp,
	//	TokenHash: salt,
	//}
	//
	//_, err = u.db.AddUser(user)
	//if err != nil {
	//	message := fmt.Sprintf("Error message: %v", err)
	//	u.logger.Print(message)
	//	u.db.ToJSON(&GenericResponse{Status: false, Message: MsgUserAlreadyExists}, rw)
	//	return
	//}
	//u.logger.Print("UserKey created successfully")
	//
	//u.db.ToJSON(&GenericResponse{Status: true, Message: "user created successfully"}, rw)
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
