package app

/**
* Pit falls:
* 	PostgreSQL query has placeholders $1, $2, $3.. not ?, ?, ?..
*	json.Marshal requires the names of the variables in a struct begins with a capital letter
 */

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Redirects the user to gamma authentication page
func HandleLogin(c *gin.Context) {
	c.Redirect(http.StatusPermanentRedirect,
		fmt.Sprintf("%s?response_type=code&client_id=%s&redirect_uri=%s",
			client.Endpoint.AuthURL,
			client.ClientID,
			client.RedirectURL))
}

// Removes the authentication cookie and makes the user unauthenticated
func HandleLogout(c *gin.Context) {
	c.SetCookie("suggestit", "", -1000, "/", cookie_domain, true, true)
}

// Exchanges Gamma authentication code with access token
// The cookie will be sent along with all further requests
func HandleAuthenticationWithCode(c *gin.Context) {
	code := c.Query("code")

	token, err := getToken(code)
	if err != nil {
		log.Println(err)
		c.AbortWithError(http.StatusUnauthorized, err)
		return
	}

	user := GetUser(token.AccessToken)
	if !HasAuthority(user) {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	c.SetCookie("suggestit",
		token.AccessToken,
		24*60*60,
		"/",
		cookie_domain,
		true,
		true)
}

func HandleCheckLogin(c *gin.Context) {
	token, err := c.Cookie("suggestit")
	user := GetUser(token)
	if err != nil || !HasAuthority(user) {
		c.SetCookie("suggestit", "", -1000, "/", cookie_domain, true, true)
		c.AbortWithError(http.StatusUnauthorized, errors.New("You are not P.R.I.T."))
		return
	}

	c.JSON(http.StatusOK, user)
}

// Creates a new suggestion
func HandleInsert(c *gin.Context) {

	var s Suggestion
	body, _ := ioutil.ReadAll(c.Request.Body)
	if err := json.Unmarshal(body, &s); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	if _, err := insertSuggestion(s); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	c.Status(http.StatusCreated)
}

// Returns all suggestions
func HandleGetSuggestions(c *gin.Context) {
	res, _ := getSuggestions()
	c.JSON(http.StatusOK, res)
}

// Returns only one suggestion when "Id" is specified
func HandleGetSuggestion(c *gin.Context) {
	id := c.Query("Id")

	var res Suggestion
	var err error

	if res, err = getSuggestion(id); err != nil {
		log.Printf("Unable to get suggestion with id: %s\n", id)
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	c.JSON(http.StatusOK, res)
}

// Deletes one suggestion
func HandleDeleteSuggestion(c *gin.Context) {
	if err := deleteSuggestion(c.Query("Id")); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
	}

	c.Status(http.StatusNoContent)
}

// Deletes multiple suggestions
func HandleDeleteSuggestions(c *gin.Context) {
	var suggestionIds struct {
		Ids []string `json:"ids"`
	}

	body, _ := ioutil.ReadAll(c.Request.Body)
	if err := json.Unmarshal(body, &suggestionIds); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	for _, id := range suggestionIds.Ids {
		deleteSuggestion(id)
	}

	c.Status(http.StatusNoContent)
}
