package app

import (
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

func contains(elements []string, key string) bool {
	for _, v := range elements {
		if v == key {
			return true
		}
	}

	return false
}
