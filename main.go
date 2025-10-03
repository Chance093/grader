package main

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

type apiConfig struct {
	db *sql.DB
}

func main() {
	db, err := initDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	cfg := apiConfig{
		db: db,
	}

	cfg.displayMainMenu()
}

func quit() {
	os.Exit(0)
}
