package auth

import (
	"testing"
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
	tests := []struct {
		input, expected string
	}{}

}
