package app

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"math/rand/v2"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
)

var oauthState = randomString(16)

var (
	oauth_url     = os.Getenv("OAUTH_URL")
	cookie_domain = os.Getenv("COOKIE_DOMAIN")
)

var oauthConfig = &oauth2.Config{
	ClientID:     os.Getenv("OAUTH_CLIENT_ID"),
	ClientSecret: os.Getenv("OAUTH_CLIENT_SECRET"),
	RedirectURL:  os.Getenv("CALLBACK_URL"),
	Scopes:       []string{"openid", "profile"},
	Endpoint: oauth2.Endpoint{
		AuthURL:  fmt.Sprintf("%s/oauth2/authorize", oauth_url),
		TokenURL: fmt.Sprintf("%s/oauth2/token", oauth_url),
	},
}

func Auth(h func(*gin.Context)) func(*gin.Context) {
	return func(c *gin.Context) {
		token, err := c.Cookie("suggestit")
		if err != nil || !HasAuthority(GetUser(token)) {
			c.SetCookie("suggestit", "", -1000, "/", cookie_domain, true, true)
			c.AbortWithError(http.StatusUnauthorized, errors.New("you are not P.R.I.T."))
			return
		}

		h(c)
	}
}

func HasAuthority(user User) bool {
	return user.Cid != ""
}

func GetUser(token string) User {
	gammaQuery := fmt.Sprintf("%s/oauth2/userinfo", oauth_url)

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

func getToken(code string) (*oauth2.Token, error) {
	return oauthConfig.Exchange(context.Background(), code)
}

func randomString(n int) string {
	const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.IntN(len(letterBytes))]
	}
	return string(b)
}
