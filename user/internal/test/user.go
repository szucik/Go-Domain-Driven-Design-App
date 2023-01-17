package test

import "github.com/szucik/trade-helper/user"

type FakeUser user.User

func (u FakeUser) WithEmail(email string) FakeUser {
	u.Email = email
	return u
}

func (u FakeUser) WithName(userName string) FakeUser {
	u.Username = userName
	return u
}

func (u FakeUser) WithPassword(pass string) FakeUser {
	u.Password = pass
	return u
}
