package data

type User struct {
	name, surname, email string
}

var users = []User{
	{
		name:    "Janusz",
		surname: "Koalski",
		email:   "janusz@wp.pl",
	},
	{
		name:    "Tomasz",
		surname: "Jakut",
		email:   "tj@gmail.pl",
	},
}
