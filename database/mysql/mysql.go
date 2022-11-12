package mysql

import (
	"context"
	"database/sql"
	"errors"

	"github.com/szucik/go-simple-rest-api/portfolio"
	"github.com/szucik/go-simple-rest-api/transaction"
	"github.com/szucik/go-simple-rest-api/user"
)

type Database struct {
	database *sql.DB
}

func NewDatabase(db *sql.DB) Database {
	return Database{
		database: db,
	}
}

func (d Database) AddPortfolio(ctx context.Context) (portfolio.Aggregate, error) {
	//stmt, err := d.database.Prepare("INSERT INTO portfolio (user_id, name) VALUES (?, ?)")
	//if err != nil {
	//	panic(err.Error())
	//}
	//result, err := stmt.Exec(p.UserId, p.Name)
	//if err != nil {
	//	return portfolio.Aggregate{}, err
	//}
	//id, err := result.LastInsertId()
	//fmt.Println(id)
	//if err != nil {
	//	return portfolio.Aggregate{}, err
	//}
	return portfolio.Aggregate{}, nil
}

// Refactor code below
var (
	ErrInvalidLoginCred = errors.New("invalid login credentials")
)

type AuthCredentials struct {
	Username string `json:"username" sql:"username"`
	Email    string `json:"email" validate:"required" sql:"email"`
	Password string `json:"password" validate:"required" sql:"password"`
}

func (d Database) SignUp(ctx context.Context) (user.Aggregate, error) {
	//stmt, err := d.database.Prepare("INSERT INTO users (username, email, password, tokenhash) VALUES (?, ?, ?, ?)")
	//if err != nil {
	//	panic(err.Error())
	//}
	//result, err := stmt.Exec(u.Username, u.Email, u.Password, u.TokenHash)
	//if err != nil {
	//	return 0, err
	//}
	//id, err := result.LastInsertId()
	//if err != nil {
	//	return 0, err
	//}
	//return id, nil
	return user.Aggregate{}, nil
}
func (d Database) SignIn(ctx context.Context) (user.Aggregate, error) {
	//var auth AuthCredentials
	//
	//err := d.db.QueryRow("SELECT username, password FROM users where email = ?", email).Scan(&auth.Username, &auth.Password)
	//if err != nil {
	//	if err == sql.ErrNoRows {
	//		return nil, ErrInvalidLoginCred
	//	}
	//	return nil, err
	//
	//}
	//
	//return &auth, nil
	return user.Aggregate{}, nil
}
func (d Database) UpdateUser(ctx context.Context) (user.Aggregate, error) {
	return user.Aggregate{}, nil
}
func (d Database) Dashboard(ctx context.Context) (user.Aggregate, error) {
	return user.Aggregate{}, nil
}
func (d Database) GetUser(ctx context.Context) (user.Aggregate, error) {
	//var user User
	//err := d.db.QueryRow("SELECT username, email FROM users where id = ?", id).Scan(&user.Username, &user.Email)
	//if err != nil {
	//	panic(err.Error())
	//}
	//return &user, nil
	return user.Aggregate{}, nil
}
func (d Database) GetUsers(ctx context.Context) (user.Aggregate, error) {
	//selDB, err := d.db.Query("SELECT * FROM users ORDER BY id DESC")
	//if err != nil {
	//	panic(err.Error())
	//}
	//user := &User{}
	//users := Users{}
	//for selDB.Next() {
	//	var id int
	//	var created, updated []uint8
	//	var email, username, password, tokenhash string
	//	err = selDB.Scan(&id, &username, &email, &password, &tokenhash, &created, &updated)
	//	if err != nil {
	//		panic(err.Error())
	//	}
	//	user.ID = id
	//	user.Username = username
	//	user.Email = email
	//	user.Password = password
	//	users = append(users, user)
	//}
	//
	//return &users, nil
	return user.Aggregate{}, nil
}

func (d Database) AddTransaction(ctx context.Context) (transaction.Aggregate, error) {
	//stmt, err := d.Database.Prepare("INSERT INTO transactions (user_id, symbol, price, quantity) VALUES (?, ?, ?, ?)")
	//if err != nil {
	//	panic(err.Error())
	//}
	//result, err := stmt.Exec(u.UserId, u.Symbol, u.Price, u.Quantity)
	//if err != nil {
	//	return 0, err
	//}
	//id, err := result.LastInsertId()
	//if err != nil {
	//	return 0, err
	//}
	return transaction.Aggregate{}, nil
}
