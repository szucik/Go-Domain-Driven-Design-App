package test

import "github.com/szucik/trade-helper/user"

type User user.User

func (u User) WithEmail(email string) User {
	u.Email = email
	return u
}

func (u User) WithName(userName string) User {
	u.Username = userName
	return u
}

func (u User) WithPassword(pass string) User {
	u.Password = pass
	return u
}
