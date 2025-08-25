package main

import (
	"database/sql"
	"embed"
	"fmt"
	"io/fs"

	"github.com/CTSDM/gator-go/internal/config"
	"github.com/pressly/goose/v3"
)

//go:embed sql/schema/*.sql
var migrationFS embed.FS

func getDatabase(cfg config.Config) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.DB_URL)
	if err != nil {
		return &sql.DB{}, err
	}
	err = db.Ping()
	if err != nil {
		return &sql.DB{}, err
	}
	return db, err
}

// If the tales do not exist it will create the tables
func createTables(db *sql.DB) error {
	_, err := db.Query("SELECT * FROM users;")
	if err == nil {
		return nil
	}
	// tables creation using the migrations with goose

	sub, err := fs.Sub(migrationFS, "sql/schema")
	if err != nil {
		return fmt.Errorf("fs.Sub: %v\n", err)
	}

	fmt.Println("Empty database. Creating the tables...")
	fmt.Println("--------------------------------------")
	goose.SetBaseFS(sub)

	if err := goose.Up(db, "."); err != nil {
		return fmt.Errorf("goose up: %v\n", err)
	}

	fmt.Println("--------------------------------------")
	fmt.Println("Success on creating the tables.")
	fmt.Printf("Starting gator-go CLI...\n\n")

	return nil
}
