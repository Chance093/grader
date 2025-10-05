package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/Chance093/grader/ascii"
	"github.com/Chance093/grader/constants"
	"github.com/Chance093/grader/prompt"
	"github.com/Chance093/grader/validation"
)

func (cfg *config) viewAssignments(className string) {
	assignments, err := cfg.db.GetClassAssignments(className)
	if err != nil {
		log.Fatalf("Error while getting class assignments: %s", err.Error())
	}

	ascii.DisplayAssignmentGrades(assignments)

	cfg.displayClassMenu(className)
}

func (cfg *config) addAssignment(className string) {
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

	if err := validation.ValidatePoints(totalPoints, correctPoints); err != nil {
		// check to see if error is conversion error or validation error
		if strings.Contains(err.Error(), "conversion") {
			log.Fatalf("Error while validating points: %s", err)
		}
		fmt.Println(err.Error())
		cfg.addAssignment(className)
		return // explicit return so when the callstack clears, it doesn't make a ton of bad assignments
	}

	if err := cfg.db.AddAssignment(assignmentName, assignmentType, className, totalPoints, correctPoints); err != nil {
		log.Fatalf("Error while adding assignment to db: %s", err.Error())
	}

	fmt.Printf("%s: %s has been added!\n", assignmentType, assignmentName)

	cfg.displayClassMenu(className)
}

func (cfg *config) editAssignment(className string) {
	assignments, err := cfg.db.GetAllClassAssignments(className)
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

func (cfg *config) editAssignmentName(assignment, className string) {
	newName, err := prompt.Input(constants.ENTER_ASSIGNMENT_NAME)
	if err != nil {
		log.Fatalf("Prompt failed %v\n", err)
	}

	if err := cfg.db.EditAssignmentName(assignment, newName, className); err != nil {
		log.Fatalf("Error while editing assignment name: %s", err)
	}

	fmt.Printf("Assignment name updated to: %s\n", newName)
}

func (cfg *config) editAssignmentGrade(assignment, className string) {
	totalPoints, err := prompt.Input(constants.ENTER_TOTAL_POINTS)
	if err != nil {
		log.Fatalf("Prompt failed %v\n", err)
	}
	correctPoints, err := prompt.Input(constants.ENTER_CORRECT_POINTS)
	if err != nil {
		log.Fatalf("Prompt failed %v\n", err)
	}

	if err := validation.ValidatePoints(totalPoints, correctPoints); err != nil {
		// check to see if error is conversion error or validation error
		if strings.Contains(err.Error(), "conversion") {
			log.Fatalf("Error while validating points: %s", err)
		}
		fmt.Println(err.Error())
		cfg.editAssignmentGrade(assignment, className)
		return // explicit return so when the callstack clears, it doesn't make a ton of bad assignments
	}

	if err := cfg.db.EditAssignmentGrade(assignment, className, totalPoints, correctPoints); err != nil {
		log.Fatalf("Error while editing assignment grade: %s", err)
	}

	fmt.Println("Assignment grade updated!")
}

func (cfg *config) editAssignmentType(assignment, className string) {
	assignmentType, err := prompt.List(constants.CHOOSE_ASSIGNMENT_TYPE, constants.ASSIGNMENT_TYPES)
	if err != nil {
		log.Fatalf("Prompt failed %v\n", err)
	}

	if err := cfg.db.EditAssignmentType(assignment, className, assignmentType); err != nil {
		log.Fatalf("Error while editing assignment type: %s", err)
	}

	fmt.Printf("Assignment type updated to: %s\n", assignmentType)
}

func (cfg *config) deleteAssignment(className string) {
	assignments, err := cfg.db.GetAllClassAssignments(className)
	if err != nil {
		log.Fatalf("Error while getting classes : %s", err.Error())
	}

	assignments = append(assignments, constants.GO_BACK)
	assignments = append(assignments, constants.MAIN_MENU)
	assignments = append(assignments, constants.QUIT)

	assignment, err := prompt.List(constants.SELECT_ASSIGNMENT_DELETE, assignments)
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
	default:
		cfg.db.DeleteAssignment(assignment, className)
	}

	fmt.Printf("Deleted assignment: %s!\n", assignment)

	cfg.displayClassMenu(className)
}
