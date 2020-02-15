package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/rs/cors"

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

	mux := http.NewServeMux()
	c := cors.New(cors.Options{
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedOrigins:   []string{"https://suggestit.chalmers.it", "http://localhost:3000"},
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
		AllowCredentials: true})

	mux.HandleFunc("/api/", handleRoot)
	mux.HandleFunc("/api/delete", auth(handleDelete))
	mux.HandleFunc("/api/authenticate", auth(func(w http.ResponseWriter, r *http.Request) { fmt.Fprintln(w, "OK") }))

	handler := c.Handler(mux)
	log.Fatal(http.ListenAndServe(":3001", handler))
}
