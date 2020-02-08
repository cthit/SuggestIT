package main

/**
* Pit falls:
* 	PostgreSQL query has placeholders $1, $2, $3.. not ?, ?, ?..
*	json.Marshal requires the names of the variables in a struct begins with a capital letter
 */

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func handleRoot(w http.ResponseWriter, r *http.Request) {
	defer recoverBadRequest(w)

	switch r.Method {
	case http.MethodPost:
		handleInsert(w, r)
	case http.MethodGet:
		if r.URL.Query().Get("Id") != "" {
			auth(handleGetSuggestion)(w, r)
		} else {
			auth(handleGetSuggestions)(w, r)
		}
	default:
		err := fmt.Sprintf("Request method not supported: %s", r.Method)
		panic(err)
	}
}

func handleInsert(w http.ResponseWriter, r *http.Request) {
	defer recoverBadRequest(w)
	var s Suggestion
	body, _ := ioutil.ReadAll(r.Body)
	if err := json.Unmarshal(body, &s); err != nil {
		panic("Could not parse body")
	}
	if _, err := insertSuggestion(s); err != nil {
		panic(err)
	}

	w.WriteHeader(http.StatusCreated)
}

func handleGetSuggestions(w http.ResponseWriter, r *http.Request) {
	var res []Suggestion
	var err error

	if res, err = getSuggestions(); err != nil {
		log.Println("Unable to get suggestions")
		http.Error(w, "Somthing went wrong", http.StatusInternalServerError)
	}

	ss, _ := json.Marshal(res)
	fmt.Fprint(w, string(ss))
}

func handleGetSuggestion(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("Id")

	var res Suggestion
	var err error

	if res, err = getSuggestion(id); err != nil {
		log.Printf("Unable to get suggestion with id: %s\n", id)
		http.Error(w, "Somthing went wrong", http.StatusInternalServerError)
		return
	}

	s, _ := json.Marshal(res)
	fmt.Fprint(w, string(s))
}

func handleDelete(w http.ResponseWriter, r *http.Request) {
	defer recoverBadRequest(w)

	if r.Method == http.MethodDelete {
		handleDeleteSuggestion(w, r)

	} else if r.Method == http.MethodPut {
		handleDeleteSuggestions(w, r)

	} else {
		http.Error(w, fmt.Sprintf("Request method not supported: %s", r.Method), http.StatusBadRequest)
	}
}

func handleDeleteSuggestion(w http.ResponseWriter, r *http.Request) {
	if err := deleteSuggestion(r.URL.Query().Get("Id")); err != nil {
		panic(err)
	}

	w.WriteHeader(http.StatusNoContent)
}

func handleDeleteSuggestions(w http.ResponseWriter, r *http.Request) {
	var suggestionIds struct {
		Ids []string
	}

	body, _ := ioutil.ReadAll(r.Body)
	if err := json.Unmarshal(body, &suggestionIds); err != nil {
		panic(fmt.Sprintf("Unable to parse body: %s", err))
	}

	for _, id := range suggestionIds.Ids {
		err := deleteSuggestion(id)

		if err != nil {
			panic(err)
		}
	}

	w.WriteHeader(http.StatusNoContent)

}

func recoverBadRequest(w http.ResponseWriter) {
	if rec := recover(); rec != nil {
		log.Println(rec)
		http.Error(w, fmt.Sprintf("%s", rec), http.StatusBadRequest)
	}
}
