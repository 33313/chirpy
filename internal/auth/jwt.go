package auth

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var ErrNoAuthHeader = errors.New("Missing authorization header")
var ErrBadAuthHeader = errors.New("Malformed authorization header")
var ErrInvalidIssuer = errors.New("Invalid token issuer")

func CreateJWT(userID int, secret string) (string, error) {
	currentTime := time.Now().UTC()
	duration := time.Duration(1) * time.Hour
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    "chirpy",
		IssuedAt:  jwt.NewNumericDate(currentTime),
		ExpiresAt: jwt.NewNumericDate(currentTime.Add(duration)),
		Subject:   strconv.Itoa(userID),
	})
	return token.SignedString([]byte(secret))
}

func ValidateJWT(tokenIn, secret string) (string, error) {
	claims := jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(
		tokenIn,
		&claims,
		func(tkn *jwt.Token) (interface{}, error) { return []byte(secret), nil },
	)
	if err != nil {
		return "", err
	}

	userIDStr, err := token.Claims.GetSubject()
	if err != nil {
		return "", err
	}

	issuer, err := token.Claims.GetIssuer()
	if err != nil {
		return "", err
	}
	if issuer != string("chirpy") {
		return "", ErrInvalidIssuer
	}
	return userIDStr, nil
}

func GetBearerToken(header string) (string, error) {
	if header == "" {
		return "", ErrNoAuthHeader
	}
	split := strings.Split(header, " ")
	if len(split) != 2 || split[0] != "Bearer" {
		return "", ErrBadAuthHeader
	}
	return split[1], nil
}

func CreateRefreshToken() string {
	bytes := make([]byte, 32)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}
