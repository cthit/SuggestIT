package app

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
)

var (
	gamma_url     = os.Getenv("GAMMA_URL")
	mock_mode     = os.Getenv("MOCK_MODE") == "True"
	allowed_group = os.Getenv("ALLOWED_GROUP")
)

var client = oauth2.Config{
	ClientID:     os.Getenv("CLIENT_ID"),
	ClientSecret: os.Getenv("AUTH_SECRET"),
	Endpoint: oauth2.Endpoint{
		AuthURL:   fmt.Sprintf("%s/api/oauth/authorize", os.Getenv("REDIRECT_GAMMA_URL")),
		TokenURL:  fmt.Sprintf("%s/api/oauth/token", gamma_url),
		AuthStyle: 0,
	},
	RedirectURL: os.Getenv("CALLBACK_URL"),
	Scopes:      nil,
}

func Auth(h func(*gin.Context)) func(*gin.Context) {
	return func(c *gin.Context) {
		token, err := c.Cookie("suggestit")
		if err != nil || !ValidUser(token) {
			c.SetCookie("suggestit", "", -1000, "/", c.Request.Host, true, true)
			c.AbortWithError(http.StatusUnauthorized, errors.New("You are not P.R.I.T."))
			return
		}

		h(c)
	}
}

func MemberOfGroup(user User, group string) bool {
	return mock_mode || contains(user.Groups,
		func(e Group) bool { return e.Active && e.SuperGroup.Name == group })
}

func ValidUser(token string) bool {
	gammaQuery := fmt.Sprintf("%s/api/users/me", gamma_url)

	client := http.Client{}
	req, _ := http.NewRequest("GET", gammaQuery, nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	resp, err := client.Do(req)
	if err != nil || resp.StatusCode < 200 || resp.StatusCode >= 300 {
		log.Println(err)
		return false
	}

	me := User{}
	text, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal(text, &me)

	return MemberOfGroup(me, "prit")
}

func getToken(grant string) (*oauth2.Token, error) {
	return client.Exchange(context.Background(), grant)
}

func contains(elements []Group, is func(Group) bool) bool {
	for _, v := range elements {
		if is(v) {
			return true
		}
	}

	return false
}
