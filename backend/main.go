package main

import (
	"chalmers.it/suggestit/internal/app"
	"database/sql"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"log"
	"os"

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
	if err := app.OpenDB(); err != nil {
		panic(err)
	}
	defer app.CloseDB()


	app.CreateSuggestionTable()

	log.Println("Starting")
	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowOrigins:     []string{"https://suggestit.chalmers.it", "http://localhost:3000"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		AllowCredentials: true}))

	router.DELETE("/api/delete", app.Auth(app.HandleDeleteSuggestion))
	router.PUT("/api/delete", app.Auth(app.HandleDeleteSuggestions))
	router.GET("/api/authenticate", app.Auth(func(c *gin.Context) {}))
	router.GET("/api/clientid", app.GetClientId)
	router.GET("/api/", app.Auth(app.HandleGetSuggestions))
	router.POST("/api/", app.HandleInsert)

	router.Run(":3001")
}
