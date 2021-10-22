package data

import "testing"

func TestValidation(t *testing.T) {
	u := &User{
		Name:    "scscs",
		Email:   "scss@wp.pl",
	}

	err := u.Validate()
	if err != nil {
		t.Fatal(err)
	}
}