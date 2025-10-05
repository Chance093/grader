package main

import (
	"log"

	"github.com/Chance093/gradr/internal/db"
	_ "github.com/mattn/go-sqlite3"
)

type config struct {
	db *db.DB
}

func main() {
	db, err := db.NewDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	cfg := config{
		db: db,
	}

	cfg.displayMainMenu()
}
