package auth

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var ErrJWTInvalidIssuer = errors.New("Invalid JWT issuer")
var ErrNoAuthHeaderIncluded = errors.New("Authentication header not included in request")
var ErrAuthHeaderMalformed = errors.New("Malformed authorization header")

const (
	defaultJWTExpiresInSeconds = 24
	defaultJWTIssuer           = "chirpy"
)

func HashPassword(password string) (string, error) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(passwordHash), nil
}

func CheckHashPassword(password, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

func CreateJWT(userId int, tokenSecret string, expiresInSeconds int) (string, error) {
	expiresAt := defaultJWTExpiresInSeconds * time.Hour

	if expiresInSeconds > 0 {
		expiresAt = time.Duration(expiresInSeconds)
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    defaultJWTIssuer,
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(expiresAt)),
		Subject:   strconv.Itoa(userId),
	})
	return token.SignedString([]byte(tokenSecret))
}

func ValidateJWT(tokenString, tokenSecret string) (int, error) {
	claimsStruct := jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(
		tokenString,
		&claimsStruct,
		func(token *jwt.Token) (interface{}, error) { return []byte(tokenSecret), nil },
	)
	if err != nil {
		return 0, err
	}

	userIDString, err := token.Claims.GetSubject()
	if err != nil {
		return 0, err
	}

	issuer, err := token.Claims.GetIssuer()
	if err != nil {
		return 0, err
	}

	if issuer != defaultJWTIssuer {
		return 0, ErrJWTInvalidIssuer
	}

	userId, err := strconv.Atoi(userIDString)

	if err != nil {
		return 0, err
	}

	return userId, nil
}

func createRefreshToken() (string, error) {
	token := make([]byte, 32)
	_, err := rand.Read(token)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(token), err
}

func GetBearerToken(headers http.Header) (string, error) {
	authHeader := headers.Get("Authorization")
	if authHeader == "" {
		return "", ErrNoAuthHeaderIncluded
	}
	splitAuth := strings.Split(authHeader, " ")
	if len(splitAuth) < 2 || splitAuth[0] != "Bearer" {
		return "", ErrAuthHeaderMalformed
	}

	return splitAuth[1], nil
}
