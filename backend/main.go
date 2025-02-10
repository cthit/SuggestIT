package main

import (
	"database/sql"
	"log"
	"os"

	"chalmers.it/suggestit/internal/app"
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
	if err := app.OpenDB(); err != nil {
		panic(err)
	}
	defer app.CloseDB()

	app.CreateSuggestionTable()

	log.Println("Starting")
	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowOrigins:     []string{"https://suggestit.chalmers.it", "http://localhost:3000"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		AllowCredentials: true}))

	r.POST("/api/auth/withCode", app.HandleAuthenticationWithCode)
	r.GET("/api/checkLogin", app.HandleCheckLogin)
	r.GET("/api/login", app.HandleLogin)
	r.POST("/api/logout", app.HandleLogout)

	r.DELETE("/api/delete", app.Auth(app.HandleDeleteSuggestion))
	r.PUT("/api/delete", app.Auth(app.HandleDeleteSuggestions))
	r.GET("/api/suggestion", app.Auth(app.HandleGetSuggestion))
	r.GET("/api/", app.Auth(app.HandleGetSuggestions))
	r.POST("/api/", app.HandleInsert)

	r.Run(":3001")
}
