package auth

import (
	"log"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestHashPassword(t *testing.T) {
	tests := []string{
		"Test123!",
		"123456",
		"Password1",
	}

	for i, test := range tests {
		_, err := HashPassword(test)
		if err != nil {
			t.Errorf("Error hashing password: %s", tests[i])
		}

	}
}

func TestCompareHashedPassowrd(t *testing.T) {
	passwords := []string{
		"Test123!",
		"123456",
		"Password1",
	}
	tests := [3]struct {
		input, expected string
	}{}
	for i, password := range passwords {
		tests[i].input = password
		hash, err := HashPassword(password)
		if err != nil {
			t.Errorf("Error hashing password: %s", password)
		}
		tests[i].expected = hash
	}

	for _, test := range tests {
		err := CheckPasswordHash(test.input, test.expected)
		if err != nil {
			t.Errorf("Error comparing hash: %s", test.input)
		}
	}
}

func TestMakeJWT(t *testing.T) {
	expiresIn, err := time.ParseDuration("1m")
	if err != nil {
		t.Errorf("failed to create expire duration")
	}
	token, err := MakeJWT(uuid.New(), "Test", expiresIn)
	if err != nil {
		t.Errorf("failed to create token: %s", err)
	}
	log.Print(token)
}
