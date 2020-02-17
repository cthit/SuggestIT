package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

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

	log.Println("Starting")
	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowOrigins:     []string{"https://suggestit.chalmers.it", "http://localhost:3000"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		AllowCredentials: true}))

	router.DELETE("/api/delete", auth(handleDeleteSuggestion))
	router.PUT("/api/delete", auth(handleDeleteSuggestions))
	router.GET("/api/authenticate", auth(func(c *gin.Context) {}))
	router.GET("/api/", auth(handleGetSuggestions))
	router.POST("/api/", handleInsert)

	router.Run(":3001")
}
