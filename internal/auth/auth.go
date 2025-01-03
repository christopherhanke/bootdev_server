package auth

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Couldn't hash password: %s", password)
		return "", err
	}
	return string(hashedPassword), nil
}

func CheckPasswordHash(password, hash string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		log.Printf("Password and hash don't align.")
		return err
	}
	return nil
}

// generate JWT token for user verification
func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    "chirpy",
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiresIn)),
		Subject:   userID.String(),
	})
	signedToken, err := token.SignedString([]byte(tokenSecret))
	if err != nil {
		log.Printf("Error signing token: %s", err)
		return "", err
	}
	return signedToken, nil
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(tokenSecret), nil
	})
	if err != nil {
		return uuid.Nil, err
	}
	subject, err := token.Claims.GetSubject()
	if err != nil {
		log.Printf("Error getting subject from token: %s", err)
		return uuid.Nil, err
	}
	id, err := uuid.Parse(subject)
	if err != nil {
		log.Printf("Error convert string to uuid: %s", err)
		return uuid.Nil, err
	}
	return id, nil
}

func GetBearerToken(headers http.Header) (string, error) {
	BearerToken := headers.Get("Authorization")
	if BearerToken == "" {
		return "", fmt.Errorf("no authorization header found")
	}
	token, ok := strings.CutPrefix(BearerToken, "Bearer ")
	if !ok {
		return "", fmt.Errorf("authorization header didn't fit")
	}

	return token, nil
}

// generate a Refresh Token string
func MakeRefreshToken() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

// get API Key from http header for authorization
func GetAPIKey(headers http.Header) (string, error) {
	APIKey := headers.Get("Authorization")
	if APIKey == "" {
		return "", fmt.Errorf("no authorization header found")
	}

	key, ok := strings.CutPrefix(APIKey, "ApiKey ")
	if !ok {
		return "", fmt.Errorf("authorization header did not fit")
	}
	return key, nil
}
