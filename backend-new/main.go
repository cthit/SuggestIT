package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/lib/pq"
)

var (
	db_host     = os.Getenv("SUGGESTIT_POSTGRES_HOST")
	db_port     = os.Getenv("SUGGESTIT_POSTGRES_PORT")
	db_user     = os.Getenv("SUGGESTIT_POSTGRES_USER")
	db_password = os.Getenv("SUGGESTIT_POSTGRES_PASSWORD")
	db_dbname   = os.Getenv("SUGGESTIT_POSTGRES_DB")
)

var db *sql.DB

func main() {
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		db_host, db_port, db_user, db_password, db_dbname)

	time.Sleep(4 * time.Second)

	var err error
	db, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	createSuggestionTable()

	http.HandleFunc("/api/", corse(handleRoot))
	http.HandleFunc("/api/delete", corse(auth(handleDelete)))
	http.HandleFunc("/api/authenticate", corse(authHandler))
	log.Fatal(http.ListenAndServe(":5001", nil))
}

func corse(h func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodOptions {
			h(w, r)
			return
		}

		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "*")
		w.Header().Set("Access-Control-Allow-Headers", "*")
		w.WriteHeader(http.StatusOK)
	}
}
