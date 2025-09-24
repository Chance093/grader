package main

import (
	"fmt"
	"log"

	"github.com/manifoldco/promptui"
)

func (cfg *apiConfig) selectAssignmentOption(className string) {
	prompt := promptui.Select{
		Label: "Choose an option",
		Items: []string{"Add assignment", "Edit assignment", "Edit grade weights", "Go back"},
	}

	_, result, err := prompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return
	}

	switch result {
	case "Add assignment":
		cfg.addAssignment(className)
	case "Edit assignment":
		editAssignment(className)
	case "Edit grade weights":
		editGradeWeights(className)
	case "Go back":
		cfg.selectClass()
	default: // Handles cases not explicitly matched
		fmt.Printf("Prompt failed %v\n", err)
		return
	}
}

func (cfg *apiConfig) addAssignment(className string) {
	assignmentPrompt := promptui.Prompt{
		Label: "Enter assignment name",
	}
	assignmentName, err := assignmentPrompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return
	}

	assignmentTypePrompt := promptui.Select{
		Label: "Choose an assignment type",
		Items: []string{"Test", "Quiz", "Homework"},
	}
	_, assignmentType, err := assignmentTypePrompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return
	}

	totalPointsPrompt := promptui.Prompt{
		Label: "Enter the total amount of points possible",
	}
	totalPoints, err := totalPointsPrompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return
	}

	correctPointsPrompt := promptui.Prompt{
		Label: "Enter the total amount of points you got correct",
	}
	correctPoints, err := correctPointsPrompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return
	}

	if err := cfg.addAssignmentToDB(assignmentName, assignmentType, className, totalPoints, correctPoints); err != nil {
		log.Fatalf("Error while adding assignment to db: %s", err.Error())
	}

	fmt.Printf("%s: %s has been added!\n", assignmentType, assignmentName)

	cfg.selectAssignmentOption(className)
}

func (cfg *apiConfig) addAssignmentToDB(name, assignmentType, className, totalPoints, correctPoints string) error {
	// take the class name and type name, and find their id's in the db
	const sqlQueryIdsStatement = `
    SELECT c.id, t.id
    FROM classes c
    CROSS JOIN assignment_types t
    WHERE c.name = ? AND t.name = ?;
    `
	var classID, typeID int
	if err := cfg.db.QueryRow(sqlQueryIdsStatement, className, assignmentType).Scan(&classID, &typeID); err != nil {
		log.Fatal(err)
	}

	// create assignment that is associated with class
	const sqlInsertAssignmentStatement = `
      INSERT INTO assignments (name, type_id, class_id, total, correct)
    VALUES (?, ?, ?, ?, ?);
    `
	if _, err := cfg.db.Exec(sqlInsertAssignmentStatement, name, typeID, classID, totalPoints, correctPoints); err != nil {
		return err
	}

	return nil
}

func editAssignment(className string) {
	fmt.Printf("Editing assignment in %s\n", className)
}

func editGradeWeights(className string) {
	fmt.Printf("Editing grade weights in %s\n", className)
}
