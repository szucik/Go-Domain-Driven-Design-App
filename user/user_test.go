package user_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/szucik/trade-helper/user"
	"github.com/szucik/trade-helper/user/internal/test"
)

var fakeUser = test.FakeUser{
	Username: "username",
	Email:    "email@test.com",
	Password: "12345678",
}

func TestUser_NewAggregate(t *testing.T) {
	t.Run("should return an error when ", func(t *testing.T) {
		testCases := map[string]struct {
			user test.FakeUser
		}{
			"address email has incorrect format": {
				user: fakeUser.WithEmail("com.invalid-email@test"),
			},
			"e-mail is shorter than 6 characters": {
				user: fakeUser.WithEmail("e@p.l"),
			},
			"password has less than 2 characters": {
				user: fakeUser.WithPassword("1234567"),
			},
			"name has less than 2 characters": {
				user: fakeUser.WithName("u"),
			},
		}

		for name, testCase := range testCases {
			t.Run(name, func(t *testing.T) {
				// when
				_, err := user.User(testCase.user).NewAggregate()
				// then
				require.Error(t, err)
			})
		}
	})

	t.Run("should create new user aggregate", func(t *testing.T) {
		// when
		aggregate, err := user.User(fakeUser).NewAggregate()
		require.NoError(t, err)
		// then
		assert.Equal(t, user.User(fakeUser), aggregate.User())
	})
}
