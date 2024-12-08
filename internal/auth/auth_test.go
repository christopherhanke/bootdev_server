package auth

import (
	"log"
	"net/http/httptest"
	"strings"
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
		return
	}
	token, err := MakeJWT(uuid.New(), "Test", expiresIn)
	if err != nil {
		t.Errorf("failed to create token: %s", err)
		return
	}
	log.Print(token)
}

func TestValidateJWT(t *testing.T) {
	id := uuid.New()
	tokenSecret := "Test"
	expiresIn, err := time.ParseDuration("1m")
	if err != nil {
		t.Errorf("failed to create expire duration")
		return
	}
	tokenString, err := MakeJWT(id, tokenSecret, expiresIn)
	if err != nil {
		t.Errorf("failed to create token: %s", err)
		return
	}
	valID, err := ValidateJWT(tokenString, tokenSecret)
	if err != nil {
		t.Errorf("Validation failed: %s", err)
		return
	}
	if valID != id {
		t.Error("Validation failed: ID failed")
		return
	}
}

func TestValidateJWTDiffTime(t *testing.T) {
	id := uuid.New()
	tokenSecret := "Test"
	expiresIn, err := time.ParseDuration("1ms")
	if err != nil {
		t.Errorf("failed to create expire duration")
		return
	}
	tokenString, err := MakeJWT(id, tokenSecret, expiresIn)
	if err != nil {
		t.Errorf("failed to create token: %s", err)
		return
	}
	expired, err := time.ParseDuration("2ms")
	if err != nil {
		t.Errorf("failed to create duration")
		return
	}
	time.Sleep(expired)
	valID, err := ValidateJWT(tokenString, tokenSecret)
	if valID != uuid.Nil {
		t.Errorf("token validation didn't catch expire")
		return
	}
	log.Printf("Validation failed: %s", err)
}

func TestValidateJWTDiffSecret(t *testing.T) {
	id := uuid.New()
	tokenSecret := "Test"
	expiresIn, err := time.ParseDuration("1m")
	if err != nil {
		t.Errorf("failed to create expire duration")
		return
	}
	tokenString, err := MakeJWT(id, tokenSecret, expiresIn)
	if err != nil {
		t.Errorf("failed to create token: %s", err)
		return
	}
	valID, err := ValidateJWT(tokenString, "1234")
	if valID != uuid.Nil {
		t.Errorf("token validation didn't catch secret")
		return
	}
	log.Printf("Validation failed: %s", err)
}

func TestGetBearerToken(t *testing.T) {
	tests := []struct {
		input, expected string
	}{
		{"Bearer Test", "Test"},
		{"Bearer Hash256-1", "Hash256-1"},
	}

	for _, test := range tests {
		req := httptest.NewRequest("GET", "/api/chirps", strings.NewReader(""))
		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("Authorization", test.input)

		stringToken, err := GetBearerToken(req.Header)
		if err != nil {
			t.Errorf("Get bearer token failed: %s", err)
			return
		}
		if stringToken != test.expected {
			t.Errorf("String token didn't match\nExpected: %s\nActual: %s", test.expected, stringToken)
		}
	}

}

func TestGetBearerTokenInvalid(t *testing.T) {
	tests := []struct {
		input, expected string
	}{
		{"missing", ""},
		{"", ""},
		{"bear this header", ""},
	}

	for _, test := range tests {
		req := httptest.NewRequest("GET", "/api/chirps", strings.NewReader(""))
		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("Authorization", test.input)

		stringToken, err := GetBearerToken(req.Header)
		if err == nil {
			t.Errorf("GetBearerToken didn't fail although header: %s\nactual token: %s", test.input, stringToken)
		}
	}
}
