package data

import "testing"

func TestValidation(t *testing.T) {
	u := &User{
		Name:  "test",
		Email: "test@test.pl",
	}

	err := u.Validate()
	if err != nil {
		t.Fatal(err)
	}
}
