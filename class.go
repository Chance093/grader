package main

import (
	"fmt"
	"log"

	"github.com/manifoldco/promptui"
)

func (cfg *apiConfig) addClass() {
	prompt := promptui.Prompt{
		Label: "Enter class name",
	}

	className, err := prompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return
	}

	prompt2 := promptui.Prompt{
		Label: "Enter subject (e.g. Math)",
	}

	subject, err := prompt2.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return
	}

	if err := cfg.addClassToDB(className, subject); err != nil {
		log.Fatalf("Failed to add to db: %s", err.Error())
	}

	cfg.startUpQuestion()
}

func (cfg *apiConfig) addClassToDB(className, subject string) error {
	const sqlInsertClassStatement = `
      INSERT INTO classes (name, subject)
    VALUES (?, ?);
    `

	if _, err := cfg.db.Exec(sqlInsertClassStatement, className, subject); err != nil {
		return err
	}

	return nil
}

func (cfg *apiConfig) selectClass() {
	prompt := promptui.Select{
		Label: "Select a Class",
		Items: []string{"Algebra 1", "Calculus", "Go Back"},
	}

	_, result, err := prompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return
	}

	if result == "Go Back" {
		cfg.startUpQuestion()
	}

	cfg.selectAssignmentOption(result)
}
