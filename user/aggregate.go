package user

type Aggregate struct {
	user User
}

func (a Aggregate) User() User {
	return a.user
}
