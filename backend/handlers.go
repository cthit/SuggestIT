package main

/**
* Pit falls:
* 	PostgreSQL query has placeholders $1, $2, $3.. not ?, ?, ?..
*	json.Marshal requires the names of the variables in a struct begins with a capital letter
 */

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func handleInsert(c *gin.Context) {

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

func handleGetSuggestions(c *gin.Context) {
	res, _ := getSuggestions()
	c.JSON(http.StatusOK, res)
}

func handleGetSuggestion(c *gin.Context) {
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

func handleDeleteSuggestion(c *gin.Context) {
	if err := deleteSuggestion(c.Query("Id")); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
	}

	c.Status(http.StatusNoContent)
}

func handleDeleteSuggestions(c *gin.Context) {
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
