package app

import (
	"database/sql"
	"fmt"
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

func OpenDB() error {
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		db_host, db_port, db_user, db_password, db_dbname)

	time.Sleep(4 * time.Second)

	var err error
	db, err = sql.Open("postgres", psqlInfo)
	return err
}

func CloseDB() {
	db.Close()
}