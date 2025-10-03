package main

import (
	"log"

	"github.com/Chance093/gradr/constants"
	"github.com/Chance093/gradr/prompt"
)

func (cfg *apiConfig) displayMainMenu() {
	chosenOption, err := prompt.List(
		constants.CHOOSE_AN_OPTION,
		[]string{
			constants.VIEW_OVERALL_GRADES,
			constants.SELECT_CLASS,
			constants.ADD_CLASS,
			constants.EDIT_CLASS,
			constants.DELETE_CLASS,
			constants.QUIT,
		},
	)
	if err != nil {
		log.Fatalf("Prompt failed %v\n", err)
	}

	switch chosenOption {
	case constants.VIEW_OVERALL_GRADES:
		cfg.viewOverallGrades()
	case constants.SELECT_CLASS:
		cfg.selectClass()
	case constants.ADD_CLASS:
		cfg.addClass()
	case constants.EDIT_CLASS:
		cfg.editClass()
	case constants.DELETE_CLASS:
		cfg.deleteClass()
	case constants.QUIT:
		quit()
	default:
		log.Fatalf("Prompt failed %v\n", err)
	}
}

func (cfg *apiConfig) displayClassMenu(className string) {
	result, err := prompt.List(constants.CHOOSE_AN_OPTION, []string{
		constants.VIEW_ASSIGNMENTS,
		constants.ADD_ASSIGNMENT,
		constants.EDIT_ASSIGNMENT,
		constants.DELETE_ASSIGNMENT,
		constants.GO_BACK,
		constants.MAIN_MENU,
		constants.QUIT,
	})
	if err != nil {
		log.Fatalf("Prompt failed %v\n", err)
	}

	switch result {
	case constants.VIEW_ASSIGNMENTS:
		cfg.viewAssignments(className)
	case constants.ADD_ASSIGNMENT:
		cfg.addAssignment(className)
	case constants.EDIT_ASSIGNMENT:
		cfg.editAssignment(className)
	case constants.DELETE_ASSIGNMENT:
		cfg.deleteAssignment(className)
	case constants.GO_BACK:
		cfg.selectClass()
	case constants.MAIN_MENU:
		cfg.displayMainMenu()
	case constants.QUIT:
		quit()
	default: // Handles cases not explicitly matched
		log.Fatalf("Prompt failed %v\n", err)
	}
}
