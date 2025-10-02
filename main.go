package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/Chance093/gradr/prompt"
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

func (cfg *apiConfig) displayMainMenu() {
	chosenOption, err := prompt.List(
		CHOOSE_AN_OPTION,
		[]string{VIEW_OVERALL_GRADES, SELECT_CLASS, ADD_CLASS, EDIT_CLASS, DELETE_CLASS, QUIT},
	)
	if err != nil {
		log.Fatalf("Prompt failed %v\n", err)
	}

	switch chosenOption {
	case VIEW_OVERALL_GRADES:
		cfg.viewOverallGrades()
	case SELECT_CLASS:
		cfg.selectClass()
	case ADD_CLASS:
		cfg.addClass()
	case EDIT_CLASS:
		cfg.editClass()
	case DELETE_CLASS:
		cfg.deleteClass()
	case QUIT:
		quit()
	default:
		log.Fatalf("Prompt failed %v\n", err)
	}
}

func quit() {
	os.Exit(0)
}
