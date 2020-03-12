package app

import (
	"errors"
	"net/http"
	"os"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

var (
	auth_secret = []byte(os.Getenv("AUTH_SECRET"))
	allowed_group = os.Getenv("ALLOWED_GROUP")
)

type User struct {
	Username string   `json:"cid"`
	Nick     string   `json:"nick"`
	Groups   []string `json:"groups"`
	jwt.StandardClaims
}

func Auth(h func(*gin.Context)) func(*gin.Context) {
	return func(c *gin.Context) {
		if token := c.GetHeader("Authorization"); !ValidUser(token) {
			c.AbortWithError(http.StatusUnauthorized, errors.New("You are not P.R.I.T."))
			return
		}

		h(c)
	}
}

func ValidUser(token string) bool {
	keyFunc := func(t *jwt.Token) (interface{}, error) { return auth_secret, nil }
	parsedToken, err := jwt.ParseWithClaims(token, &User{}, keyFunc)
	if err != nil {
		return false
	}

	claims, ok := parsedToken.Claims.(*User)
	if !ok || !parsedToken.Valid {
		return false
	}

	return contains(claims.Groups, allowed_group)
}

func contains(elements []string, key string) bool {
	for _, v := range elements {
		if v == key {
			return true
		}
	}

	return false
}
