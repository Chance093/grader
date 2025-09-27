package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/manifoldco/promptui"
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

	cfg.startUpQuestion()
}

func (cfg *apiConfig) startUpQuestion() {
	prompt := promptui.Select{
		Label: "Choose an option",
		Items: []string{"View Overall Grades", "Add a Class", "Select a Class", "Quit"},
	}

	_, result, err := prompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return
	}

	switch result {
	case "Add a Class":
		cfg.addClass()
	case "Select a Class":
		cfg.selectClass()
	case "View Overall Grades":
	  cfg.viewOverallGrades()
	case "Quit":
		quit()
	default: // Handles cases not explicitly matched
		fmt.Printf("Prompt failed %v\n", err)
		return
	}
}

func quit() {
	os.Exit(0)
}
