package app

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

var (
	client_id      = os.Getenv("CLIENT_ID")
	auth_secret    = os.Getenv("AUTH_SECRET")
	allowed_group  = os.Getenv("ALLOWED_GROUP")
	gamma_url      = os.Getenv("GAMMA_URL")
	redirect_url   = os.Getenv("REDIRECT_URL")
	basicAothToken = base64.StdEncoding.EncodeToString([]byte(client_id + ":" + auth_secret))
	mock_mode      = os.Getenv("MOCK_MODE") == "True"
	auth_url       = fmt.Sprintf("%s/api/oauth/authorize", gamma_url)
	token_url      = fmt.Sprintf("%s/api/oauth/token", gamma_url)
)

func Auth(h func(*gin.Context)) func(*gin.Context) {
	return func(c *gin.Context) {
		if token := c.GetHeader("Authorization"); !ValidUser(token) {
			c.AbortWithError(http.StatusUnauthorized, errors.New("You are not P.R.I.T."))
			return
		}

		h(c)
	}
}

func ValidToken(token string) bool {
	keyFunc := func(t *jwt.Token) (interface{}, error) { return []byte(""), nil }

	parsedToken, _ := jwt.ParseWithClaims(token, &TokenClaims{}, keyFunc)

	claims, ok := parsedToken.Claims.(*TokenClaims)
	if !ok {
		return false
	}

	return claims.ClientId == client_id
}

func MemberOfGroup(user User, group string) bool {
	if mock_mode {
		return true
	}
	return contains(user.Relationships,
		func(e Relationship) bool { return e.Group.Active && e.SuperGroup.Name == group })
}

func ValidUser(token string) bool {
	gammaQuery := fmt.Sprintf("%s/api/users/me", gamma_url)

	client := http.Client{}
	req, _ := http.NewRequest("GET", gammaQuery, nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return false
	}

	me := User{}
	text, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal(text, &me)

	return ValidToken(token) && MemberOfGroup(me, allowed_group)
}

func getToken(grant string, redirect_uri string) (TokenResponse, error) {

	gammaQuery := fmt.Sprintf(token_url+
		"?grant_type=authorization_code&client_id=%s&redirect_uri=%s&code=%s",
		client_id,
		redirect_uri,
		grant)
	client := http.Client{}
	req, _ := http.NewRequest("POST", gammaQuery, nil)

	req.Header.Set("Authorization", fmt.Sprintf("Basic %s", basicAothToken))

	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return TokenResponse{}, err
	}

	text, er := ioutil.ReadAll(resp.Body)
	if er != nil {
		log.Println("Unable to read body")
		log.Println(er)
		return TokenResponse{}, er
	}

	body := TokenResponse{}
	json.Unmarshal(text, &body)
	return body, nil
}

func contains(elements []Relationship, is func(Relationship) bool) bool {
	for _, v := range elements {
		if is(v) {
			return true
		}
	}

	return false
}
