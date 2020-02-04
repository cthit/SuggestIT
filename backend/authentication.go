package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"
	"github.com/dgrijalva/jwt-go"
)

var (
	prit_password = os.Getenv("PRIT_PASSWORD")
	prit_key      = os.Getenv("PRIT_AUTH_KEY")
	auth_secret   = []byte(os.Getenv("AUTH_SECRET"))
)

func auth(h func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if token := r.Header.Get("Authorization"); validUser(token) {
			http.Error(w, "You are not P.R.I.T.", http.StatusUnauthorized)
			return
		}

		h(w, r)
	}
}

func authHandler(w http.ResponseWriter, r *http.Request) {
	defer recoverBadRequest(w)

	if r.Method != http.MethodPut && r.Method != http.MethodGet {
		http.Error(w, fmt.Sprintf("Request method not supported: %s", r.Method), http.StatusBadRequest)
		return
	}

	if token := r.Header.Get("Authorization"); validUser(token) {
		fmt.Fprintf(w, "You are P.R.I.T.")
		return
	}

	var userInput struct {
		CID string
		Password string
	}

	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		panic(err)
	}

	if err := json.Unmarshal(body, &userInput); err != nil {
		panic(err)
	}

	var user User
	user, err = login_ldap(userInput.CID, userInput.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	if !user.MemberOfPRIT {
		http.Error(w, "You are not P.R.I.T.", http.StatusUnauthorized)
		return
	}

	var token string
	token, err = createToken(user)
	if err != nil {
		http.Error(w, "Unable to create token", http.StatusUnauthorized)
		return
	}

	fmt.Fprintln(w, fmt.Sprintf(`{"token":"%s"}`, token))
}

func createToken(user User) (string, error) {
	user.ExpiresAt = time.Now().Add(time.Hour * 24*30).Unix() 

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, user)
	return token.SignedString(auth_secret)
}

func validUser(token string) bool {
	keyFunc := func(t *jwt.Token) (interface{}, error) {return auth_secret, nil}
	parsedToken, err := jwt.ParseWithClaims(token, &User{}, keyFunc)
	if err != nil {
		return false
	}

	claims, ok := parsedToken.Claims.(*User)
	if !ok || !parsedToken.Valid {
		return false
	}

	return claims.MemberOfPRIT
}