package utils

import (
	"fmt"
	"os"
	"strconv"
	"testproject/task004/models"
	"time"

	"github.com/golang-jwt/jwt"
)

func GenerateToken(user models.User) (string, error) {
	claims := jwt.MapClaims{
		"id":       strconv.FormatUint(user.ID, 10),
		"username": user.Username,
		"exp":      time.Now().Add(time.Hour * 72).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}

func ValidateToken(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
}
