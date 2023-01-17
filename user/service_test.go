package user_test

import (
	"log"
	"os"
	"testing"

	"github.com/szucik/trade-helper/database/fake"
	"github.com/szucik/trade-helper/user"
)

var userService = user.Users{
	Logger:       log.New(os.Stdout, "logger", log.LstdFlags),
	Database:     fake.NewDatabase(),
	NewAggregate: user.User.NewAggregate,
}

func TestUserService_SignUp(t *testing.T) {
	// t.Run("should create new user", func(t *testing.T) {
	// 	_, err := userService.SignUp(user.User(fakeUser))
	// 	assert.Error(t, err)
	// })
}

func TestUsers_GetUsers(t *testing.T) {
	// t.Run("should return an error when ", func(t *testing.T) {
	// 	testCases := map[string]struct {
	// 		user test.FakeUser
	// 	}{}
	//
	// 	for name, testCase := range testCases {
	// 		t.Run(name, func(t *testing.T) {
	//
	// 		})
	// 	}
	// })

	t.Run("should return list of users", func(t *testing.T) {
		// aggregates, _ := userService.GetUsers()
		// // when
		// aggregate, err := user.user(fakeUser).NewAggregate()
		// // then
		// assert.Equal(t, user.user(fakeUser), aggregate.user())
	})

}
