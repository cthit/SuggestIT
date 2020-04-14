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

	//"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

const (
	AUTH_URL =  "https://gamma.chalmers.it/api/oauth/authorize"
	TOKEN_URL =  "https://gamma.chalmers.it/api/oauth/token"
)

var (
	client_id = os.Getenv("CLIENT_ID")
	auth_secret = os.Getenv("AUTH_SECRET")
	allowed_group = os.Getenv("ALLOWED_GROUP")
	basicAothToken = base64.StdEncoding.EncodeToString([]byte(client_id + ":" + auth_secret))
)

/*type User struct {
	Username string   `json:"cid"`
	Nick     string   `json:"nick"`
	Groups   []string `json:"groups"`
	jwt.StandardClaims
}*/

type TokenResponse struct {
	Access_token string `json:"access_token"`
	Token_type   string `json:"token_type"`
	Expires_in   string `json:"expires_in"`
	Scope        string `json:"scope"`
	Nick         string `json:"nick"`
	Uid          string `json:"uid"`
	ISS          string `json:"iss"`
	JTI          string `json:"jti"`
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
	gammaQuery := fmt.Sprintf("https://gamma.chalmers.it/api/auth/valid_token?token=%s", token)

	client := http.Client{}
	req, _ := http.NewRequest("GET",gammaQuery, nil)

	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return false
	}

	text, er := ioutil.ReadAll(resp.Body)
	//TODO: Check if the token is for this client
	//TODO: Check if the user is a member of PRIT
	if er != nil || string(text) != "true" {
		return false
	}

	return true
}

func getToken(grant string, redirect_uri string) (TokenResponse, error) {

	gammaQuery := fmt.Sprintf(TOKEN_URL +
		"?grant_type=authorization_code&client_id=%s&redirect_uri=%s&code=%s",
		client_id,
		redirect_uri,
		grant)
	client := http.Client{}
	req, _ := http.NewRequest("POST",gammaQuery, nil)

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

func contains(elements []string, key string) bool {
	for _, v := range elements {
		if v == key {
			return true
		}
	}

	return false
}
