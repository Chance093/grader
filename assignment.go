package main

import (
	"errors"
	"fmt"
	"log"
	"strconv"

	"github.com/Chance093/gradr/ascii"
	"github.com/Chance093/gradr/constants"
	"github.com/Chance093/gradr/prompt"
	"github.com/manifoldco/promptui"
)

func (cfg *apiConfig) viewAssignments(className string) {
	assignments, err := cfg.getClassAssignmentsFromDB(className)
	if err != nil {
		log.Fatalf("Error while getting class assignments: %s", err.Error())
	}

	ascii.DisplayAssignmentGrades(assignments)

	cfg.displayClassMenu(className)
}

func (cfg *apiConfig) addAssignment(className string) {
	// capture inputs
	assignmentName, err := prompt.Input(constants.ENTER_ASSIGNMENT_NAME)
	if err != nil {
		log.Fatalf("Prompt failed %v\n", err)
	}
	assignmentType, err := prompt.List(constants.CHOOSE_ASSIGNMENT_TYPE, constants.ASSIGNMENT_TYPES)
	if err != nil {
		log.Fatalf("Prompt failed %v\n", err)
	}
	totalPoints, err := prompt.Input(constants.ENTER_TOTAL_POINTS)
	if err != nil {
		log.Fatalf("Prompt failed %v\n", err)
	}
	correctPoints, err := prompt.Input(constants.ENTER_CORRECT_POINTS)
	if err != nil {
		log.Fatalf("Prompt failed %v\n", err)
	}

	if err := validatePoints(totalPoints, correctPoints, className); err != nil {
		fmt.Println(err.Error())
		cfg.addAssignment(className)
		return // explicit return so when the callstack clears, it doesn't make a ton of bad assignments
	}

	if err := cfg.addAssignmentToDB(assignmentName, assignmentType, className, totalPoints, correctPoints); err != nil {
		log.Fatalf("Error while adding assignment to db: %s", err.Error())
	}

	fmt.Printf("%s: %s has been added!\n", assignmentType, assignmentName)

	cfg.displayClassMenu(className)
}

func validatePoints(totalPoints, correctPoints, className string) error {
	totalPointsInt, err := strconv.Atoi(totalPoints)
	if err != nil {
		log.Fatalf("Failed to convert string to int: %s", err.Error())
	}

	correctPointsInt, err := strconv.Atoi(correctPoints)
	if err != nil {
		log.Fatalf("Failed to convert string to int: %s", err.Error())
	}

	if totalPointsInt < correctPointsInt {
		return errors.New("Total points must be greater than or equal to correct points")
	}

	return nil
}

func (cfg *apiConfig) editAssignment(className string) {
	assignments, err := cfg.getAllClassAssignmentsFromDB(className)
	if err != nil {
		log.Fatalf("Error while getting class assignments: %s", err.Error())
	}

	assignments = append(assignments, constants.GO_BACK)
	assignments = append(assignments, constants.MAIN_MENU)
	assignments = append(assignments, constants.QUIT)

	assignment, err := prompt.List(constants.SELECT_ASSIGNMENT_EDIT, assignments)
	if err != nil {
		log.Fatalf("Prompt failed %v\n", err)
	}

	switch assignment {
	case constants.GO_BACK:
		cfg.displayClassMenu(className)
	case constants.MAIN_MENU:
		cfg.displayMainMenu()
	case constants.QUIT:
		quit()
	}

  result, err := prompt.List(constants.CHOOSE_OPTION_EDIT, constants.EDIT_ASSIGNMENT_OPTS)
	if err != nil {
		log.Fatalf("Prompt failed %v\n", err)
	}

	switch result {
	case constants.EDIT_ASSIGNMENT_OPTS[0]:
		cfg.editAssignmentName(assignment, className)
	case constants.EDIT_ASSIGNMENT_OPTS[1]:
		cfg.editAssignmentGrade(assignment, className)
	case constants.EDIT_ASSIGNMENT_OPTS[2]:
		cfg.editAssignmentType(assignment, className)
	default: // Handles cases not explicitly matched
		log.Fatalf("Prompt failed %v\n", err)
	}

	cfg.displayClassMenu(className)
}

func (cfg *apiConfig) editAssignmentName(assignment, className string) {
	prompt := promptui.Prompt{
		Label: "Enter new assignment name",
	}
	newName, err := prompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return
	}

	if err := cfg.editAssignmentNameInDB(assignment, newName, className); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Assignment name updated to: %s\n", newName)
}

func (cfg *apiConfig) editAssignmentNameInDB(oldName, newName, className string) error {
	const sqlUpdateAssignmentNameStatement = `
      UPDATE assignments SET name = ? WHERE name = ? AND 
    class_id = (SELECT id FROM classes WHERE name = ?);
    `

	if _, err := cfg.db.Exec(sqlUpdateAssignmentNameStatement, newName, oldName, className); err != nil {
		return err
	}

	return nil
}

func (cfg *apiConfig) editAssignmentGrade(assignment, className string) {
	totalPrompt := promptui.Prompt{
		Label: "Enter the total amount of points possible",
	}
	total, err := totalPrompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return
	}

	correctPrompt := promptui.Prompt{
		Label: "Enter the total amount of points you got correct",
	}
	correct, err := correctPrompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return
	}

	totalInt, err := strconv.Atoi(total)
	if err != nil {
		log.Fatal(err)
	}

	correctInt, err := strconv.Atoi(correct)
	if err != nil {
		log.Fatal(err)
	}

	if err := cfg.editAssignmentGradeInDB(assignment, className, totalInt, correctInt); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Assignment grade updated!")
}

func (cfg *apiConfig) editAssignmentGradeInDB(assignment, className string, total, correct int) error {
	const sqlUpdateAssignmentGradeStatement = `
      UPDATE assignments SET correct = ?, total = ? WHERE name = ? AND 
    class_id = (SELECT id FROM classes WHERE name = ?);
    `

	if _, err := cfg.db.Exec(sqlUpdateAssignmentGradeStatement, correct, total, assignment, className); err != nil {
		return err
	}

	return nil
}

func (cfg *apiConfig) editAssignmentType(assignment, className string) {
	assignmentTypePrompt := promptui.Select{
		Label: "Choose an assignment type",
		Items: []string{"Test", "Quiz", "Homework"},
	}
	_, assignmentType, err := assignmentTypePrompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return
	}

	if err := cfg.editAssignmentTypeInDB(assignment, className, assignmentType); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Assignment type updated to: %s\n", assignmentType)
}

func (cfg *apiConfig) editAssignmentTypeInDB(assignment, className, assignmentType string) error {
	const sqlUpdateAssignmentNameStatement = `
      UPDATE assignments SET type_id = ? WHERE name = ? AND 
    class_id = (SELECT id FROM classes WHERE name = ?);
    `

	typeMap := map[string]int{
		"Test":     1,
		"Quiz":     2,
		"Homework": 3,
	}

	if _, err := cfg.db.Exec(sqlUpdateAssignmentNameStatement, typeMap[assignmentType], assignment, className); err != nil {
		return err
	}

	return nil
}

func (cfg *apiConfig) deleteAssignment(className string) {
	assignments, err := cfg.getAllClassAssignmentsFromDB(className)
	if err != nil {
		log.Fatalf("Error while getting classes : %s", err.Error())
	}

	assignments = append(assignments, "Go back")
	assignments = append(assignments, "Main Menu")

	prompt := promptui.Select{
		Label: "Select an assignment to delete",
		Items: assignments,
	}

	_, result, err := prompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return
	}

	switch result {
	case "Go back":
		cfg.selectAssignmentOption(className)
	case "Main Menu":
		cfg.displayMainMenu()
	default:
		cfg.deleteAssignmentFromDB(result, className)
	}

	fmt.Printf("Deleted assignment: %s!\n", result)

	cfg.selectAssignmentOption(className)
}

func (cfg *apiConfig) deleteAssignmentFromDB(assignmentName, className string) error {
	const sqlDeleteAssignmentStatement = `
      DELETE FROM assignments WHERE name = ? AND class_id = 
    (SELECT id FROM classes WHERE name = ?);
    `

	if _, err := cfg.db.Exec(sqlDeleteAssignmentStatement, assignmentName, className); err != nil {
		return err
	}

	return nil
}
