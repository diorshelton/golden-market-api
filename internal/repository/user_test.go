package repository

import (
	"fmt"
	"os"
	"testing"

	"github.com/diorshelton/golden-market-api/internal/database"
	"github.com/joho/godotenv"
)

func TestSpunUpDB(t *testing.T) {
	// Load .env file before checking environment variables
	_ = godotenv.Load("../../.env")

	// Skip test if TEMP_DB_URL is not set
	if os.Getenv("TEMP_DB_URL") == "" {
		t.Skip("TEMP_DB_URL not set, skipping database tests")
	}

	dbConnection, err := database.SetupTestDB()
	if err != nil {
		t.Fatalf("Failed to SetupTestUserDB %v", err)
	}
	defer dbConnection.Close()

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
