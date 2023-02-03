package user_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/szucik/trade-helper/user"
	"github.com/szucik/trade-helper/user/internal/test"
)

var (
	testUser = test.User{
		Username: "username",
		Email:    "email@test.com",
		Password: "12345678",
	}
)

func TestUser_NewAggregate(t *testing.T) {
	t.Run("should return an error when ", func(t *testing.T) {
		testCases := map[string]struct {
			user test.User
		}{
			"address e-mail has incorrect format": {
				user: testUser.WithEmail("com.invalid-email@test"),
			},
			"e-mail is shorter than 6 characters": {
				user: testUser.WithEmail("e@p.l"),
			},
			"password has less than 2 characters": {
				user: testUser.WithPassword("1234567"),
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

	t.Run("should create new user aggregate when all validations are passed", func(t *testing.T) {
		// when
		aggregate, err := user.User(testUser).NewAggregate()
		require.NoError(t, err)

		// then
		testUser.Password = aggregate.User().Password
		assert.Equal(t, user.User(testUser), aggregate.User())
	})
}
