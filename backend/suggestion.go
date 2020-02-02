package main

import (
	"database/sql"
	"errors"
	"fmt"
)

type Suggestion struct {
	Id        string `json:"id"`
	Timestamp string `json:"timestamp"`
	Title     string `json:"title"`
	Text      string `json:"text"`
	Author    string `json:"author"`
}

func createSuggestionTable() (sql.Result, error) {
	db.Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp"`)
	return db.Exec(`
		CREATE TABLE IF NOT EXISTS suggestion (
		id UUID NOT NULL,
		timestamp TIMESTAMP NOT NULL,
		title TEXT NOT NULL,
		text TEXT NOT NULL,
		author TEXT,
		PRIMARY KEY (UUID)
		)`)
}

func insertSuggestion(s Suggestion) (sql.Result, error) {
	return db.Exec(`
	INSERT INTO suggestion 
	(id, timestamp, title, text, author) 
	VALUES(uuid_generate_v4(), now(), $1,$2,$3)`,
		s.Title, s.Text, s.Author)
}

func getSuggestions() ([]Suggestion, error) {
	rows, err := db.Query("SELECT * FROM suggestion")
	if err != nil {
		return nil, err
	}

	var suggestions []Suggestion
	for rows.Next() {
		var s Suggestion
		err = rows.Scan(
			&s.Id,
			&s.Timestamp,
			&s.Title,
			&s.Text,
			&s.Author)

		if err == nil {
			suggestions = append(suggestions, s)
		}
	}

	return suggestions, nil
}

func getSuggestion(id string) (Suggestion, error) {
	row := db.QueryRow("SELECT * FROM suggestion WHERE id=$1", id)

	var s Suggestion
	err := row.Scan(
		&s.Id,
		&s.Timestamp,
		&s.Title,
		&s.Text,
		&s.Author)

	if err != nil {
		return Suggestion{}, err
	}

	return s, nil
}

func deleteSuggestion(Id string) error {
	res, err := db.Exec("DELETE FROM suggestion WHERE Id=$1", Id)
	if err != nil {
		return err
	}
	if n, err := res.RowsAffected(); err != nil {
		return err
	} else if n != 1 {
		return errors.New(fmt.Sprintf("Could not find element with id %s", Id))
	}
	return nil
}
