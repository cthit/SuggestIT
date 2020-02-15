package main

import (
	"net/http"
	"os"

	"github.com/dgrijalva/jwt-go"
)

var (
	auth_secret = []byte(os.Getenv("AUTH_SECRET"))
)

type User struct {
	Username string   `json:"cid"`
	Nick     string   `json:"nick"`
	Groups   []string `json:"groups"`
	jwt.StandardClaims
}

func auth(h func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if token := r.Header.Get("Authorization"); !ValidUser(token) {
			http.Error(w, "You are not P.R.I.T.", http.StatusUnauthorized)
			return
		}

		h(w, r)
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

	return contains(claims.Groups, "digit")
}

func contains(elements []string, key string) bool {
	for _, v := range elements {
		if v == key {
			return true
		}
	}

	return false
}
