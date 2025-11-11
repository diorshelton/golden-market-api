package repository

import (
	"fmt"
	"testing"

	"github.com/diorshelton/golden-market-api/internal/database"
)

func TestSpunUpDB(t *testing.T) {
	dbConnection, err := database.SetupTestUserDB()
	if err != nil {
		t.Fatalf("Failed to SetupTestUserDB %v", err)
	}

	type Test struct {
		firstName    string
		lastName     string
		email        string
		username     string
		passwordHash string
	}
	var j = Test{
		firstName:    "Jake",
		lastName:     "The Dog",
		email:        "jdog@example.com",
		username:     "jdog",
		passwordHash: "password123",
	}

	userDb := NewUserRepository(dbConnection)

	jake, err := userDb.CreateUser(j.username, j.firstName, j.lastName, j.email, j.passwordHash)

	if err != nil {
		t.Errorf("This error occurred :%v", err)
	}

	if jake.Username == "" {
		t.Errorf("No Username")
	}
	t.Run("pull from database", func(t *testing.T) {
		_, err := userDb.GetUserByEmail(jake.Email)
		if err != nil {
			t.Errorf("an error occurred %v", err)
		}
	})

	t.Run("Get all users", func(t *testing.T) {
		users, err := userDb.GetAllUsers()
		if err != nil {
			t.Errorf("An error occurred %v", err)
		}

		fmt.Printf("Returned %v user(s)\n", len(users))

		if len(users) < 1 {
			t.Errorf("No users")
		}
	})
}

func TestCreateUser(t *testing.T) {
	// userDatabase := setUpTestPGX(t)

	// _, err := userDatabase.CreateUser(jake.username, jake.firstName, jake.lastName, jake.email, jake.passwordHash)
	// if err != nil {
	// 	t.Errorf("Failed to create user %v", err)
	// }

	// if user.FirstName != "jake" {
	// 	t.Error("expected ")
	// }

	// t.Run("fetch user", func(t *testing.T) {
	// 	userDatabase.GetUserByEmail(jake.email)
	// })
}
