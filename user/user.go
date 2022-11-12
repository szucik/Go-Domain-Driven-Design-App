package user

import "time"

type User struct {
	ID        int       `json:"id"`
	Username  string    `json:"username" validate:"required"`
	Email     string    `json:"email" validate:"required" sql:"email"`
	Password  string    `json:"password" validate:"required" sql:"password"`
	TokenHash string    `json:"token_hash" sql:"token_hash"`
	Created   time.Time `json:"created" sql:"created"`
	Updated   time.Time `json:"updated" sql:"updated"`
}

func (u User) NewAggregate() (Aggregate, error) {
	//TODO Add Validation error
	//var aggregate Aggregate
	//return aggregate, tradehelpererrors.ValidationError
	//if u.Username == "" {
	//	return Aggregate{}, nil
	//}
	return Aggregate{
		user: u,
	}, nil
}
