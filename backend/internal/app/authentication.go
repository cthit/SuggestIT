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
	gamma_url       = os.Getenv("GAMMA_URL")
	mock_mode       = os.Getenv("MOCK_MODE") == "True"
	gamma_authority = os.Getenv("GAMMA_AUTHORITY")
	cookie_domain   = os.Getenv("COOKIE_DOMAIN")
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
		if err != nil || !HasAuthority(GetUser(token)) {
			c.SetCookie("suggestit", "", -1000, "/", cookie_domain, true, true)
			c.AbortWithError(http.StatusUnauthorized, errors.New("You are not P.R.I.T."))
			return
		}

		h(c)
	}
}

func HasAuthority(user User) bool {
	return mock_mode || contains(user.Authorities,
		func(e Authority) bool { return e.Authority == gamma_authority })
}

func GetUser(token string) User {
	gammaQuery := fmt.Sprintf("%s/api/users/me", gamma_url)

	client := http.Client{}
	req, _ := http.NewRequest("GET", gammaQuery, nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	resp, err := client.Do(req)
	if err != nil || resp.StatusCode < 200 || resp.StatusCode >= 300 {
		log.Println(err)
		return User{}
	}

	me := User{}
	text, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal(text, &me)

	return me
}

func getToken(grant string) (*oauth2.Token, error) {
	return client.Exchange(context.Background(), grant)
}

func contains(elements []Authority, is func(Authority) bool) bool {
	for _, v := range elements {
		if is(v) {
			return true
		}
	}

	return false
}
