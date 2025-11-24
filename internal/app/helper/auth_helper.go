package helper

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var secretKey = []byte(os.Getenv("JWT_SECRET_KEY"))

type UserClaims struct {
	jwt.RegisteredClaims
	Sub  uint64 `json:"sub"`
	Role string `json:"role"`
}

func GetTokens(userId uint) (string, string, error) {
	accessToken, err := getAccessToken(userId)
	if err != nil {
		return "", "", err
	}
	refreshToken, err := getRefreshToken(userId)
	if err != nil {
		return "", "", err
	}
	return accessToken, refreshToken, nil
}

func getAccessToken(userId uint) (string, error) {

	claims := jwt.MapClaims{
		"exp": jwt.NewNumericDate(time.Now().Add(2 * time.Hour)),
		"sub": userId,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString(secretKey)
	return ss, err
}

func getRefreshToken(userId uint) (string, error) {

	claims := jwt.MapClaims{
		"exp": jwt.NewNumericDate(time.Now().Add(168 * time.Hour)),
		"sub": userId,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString(secretKey)
	return ss, err
}

func ParseToken(tokenString string) (*UserClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &UserClaims{}, func(token *jwt.Token) (any, error) {
		return secretKey, nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*UserClaims)
	if !ok || !token.Valid {
		return nil, errors.New("could not parse claims")
	}

	return claims, nil
}
