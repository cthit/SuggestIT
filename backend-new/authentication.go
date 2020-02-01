package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

var (
	prit_password = os.Getenv("PRIT_PASSWORD")
	prit_key      = os.Getenv("PRIT_AUTH_KEY")
)

func auth(h func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if key := r.Header.Get("Authorization"); key != prit_key {
			http.Error(w, "You are not P.R.I.T.", http.StatusUnauthorized)
			return
		}

		h(w, r)
	}
}

func authHandler(w http.ResponseWriter, r *http.Request) {
	defer recoverBadRequest(w)

	if key := r.Header.Get("Authorization"); r.Method == http.MethodGet && key == prit_key {
		fmt.Fprintf(w, "You are P.R.I.T.")
		return
	}

	if r.Method != http.MethodPut {
		http.Error(w, fmt.Sprintf("Request method not supported: %s", r.Method), http.StatusBadRequest)
		return
	}

	var pass struct {
		Password string
	}

	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		panic(err)
	}

	if err := json.Unmarshal(body, &pass); err != nil {
		panic(err)
	}

	if pass.Password != prit_password {
		http.Error(w, "You are not P.R.I.T.", http.StatusUnauthorized)
		return
	}

	fmt.Fprintln(w, fmt.Sprintf(`{"key":"%s"}`, prit_key))
}
