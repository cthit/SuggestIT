package app

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
)

var (
	gamma_url     = os.Getenv("GAMMA_URL")
	mock_mode     = os.Getenv("MOCK_MODE") == "True"
	cookie_domain = os.Getenv("COOKIE_DOMAIN")
)

var oauthConfig = &oauth2.Config{
	ClientID:     os.Getenv("OAUTH_CLIENT_ID"),
	ClientSecret: os.Getenv("OAUTH_CLIENT_SECRET"),
	RedirectURL:  os.Getenv("OAUTH_CALLBACK_URL"),
	Scopes:       []string{"openid"},
	Endpoint: oauth2.Endpoint{
		AuthURL:  "https://auth.chalmers.it/oauth2/authorize",
		TokenURL: "https://auth.chalmers.it/oauth2/token",
	},
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
	return true
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
	text, _ := io.ReadAll(resp.Body)
	json.Unmarshal(text, &me)

	return me
}

func getToken(grant string) (*oauth2.Token, error) {
	return oauthConfig.Exchange(context.Background(), grant)
}

func contains(elements []Authority, is func(Authority) bool) bool {
	for _, v := range elements {
		if is(v) {
			return true
		}
	}

	return false
}
