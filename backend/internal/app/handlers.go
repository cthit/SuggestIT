package app

/**
* Pit falls:
* 	PostgreSQL query has placeholders $1, $2, $3.. not ?, ?, ?..
*	json.Marshal requires the names of the variables in a struct begins with a capital letter
 */

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

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

func GetClientId(c *gin.Context) {
	var id struct {
		Client_id string `json:"client_id"`
	}
	id.Client_id = os.Getenv("CLIENT_ID")

	c.JSON(http.StatusOK, id)
}

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

func HandleGetSuggestions(c *gin.Context) {
	res, _ := getSuggestions()
	c.JSON(http.StatusOK, res)
}

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

func HandleDeleteSuggestion(c *gin.Context) {
	if err := deleteSuggestion(c.Query("Id")); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
	}

	c.Status(http.StatusNoContent)
}

func HandleDeleteSuggestions(c *gin.Context) {
	var suggestionIds struct {
		Ids []string
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

func HandleGetToken(c *gin.Context) {
	code := c.Query("code")
	redirect_uri := c.Query("redirect_uri")
	authKey := base64.StdEncoding.EncodeToString([]byte(
		fmt.Sprintf("%s:%s", os.Getenv("CLIENT_ID"), os.Getenv("AUTH_SECRET"))))
	gammaQuery := fmt.Sprintf("https://gamma.chalmers.it/api/oauth/token" +
		"?grant_type=authorization_code&client_id=%s&redirect_uri=%s&code=%s",
		os.Getenv("CLIENT_ID"),
		redirect_uri,
		code)
	client := http.Client{}
	req, _ := http.NewRequest("POST",gammaQuery, nil)

	req.Header.Set("Authorization", fmt.Sprintf("Basic %s", authKey))

	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return
	}

	text, er := ioutil.ReadAll(resp.Body)
	if er != nil {
		log.Println("Unable to read body")
		log.Println(er)
		c.AbortWithError(http.StatusInternalServerError, er)
		return
	}

	body := TokenResponse{}
	json.Unmarshal(text, &body)
	c.JSON(http.StatusOK, body)
}
