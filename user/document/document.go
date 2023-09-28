package document

import (
	"github.com/szucik/trade-helper/user"
	"time"
)

type User struct {
	Username  string    `bson:"username,omitempty"`
	Email     string    `bson:"email,omitempty"`
	Password  string    `bson:"password,omitempty"`
	TokenHash string    `bson:"token_hash,omitempty"`
	Created   time.Time `bson:"created,omitempty"`
	Updated   time.Time `bson:"updated,omitempty"`
}

func NewDocument(aggregate user.Aggregate) User {
	user := aggregate.User()
	return User{
		Username:  user.Username,
		Email:     user.Email,
		Password:  user.Password,
		TokenHash: user.TokenHash,
		Created:   user.Created,
		Updated:   user.Updated,
	}
}

func NewAggregate(aggregate user.Aggregate) User {
	user := aggregate.User()
	return User{
		Username:  user.Username,
		Email:     user.Email,
		Password:  user.Password,
		TokenHash: user.TokenHash,
		Created:   user.Created,
		Updated:   user.Updated,
	}
}
