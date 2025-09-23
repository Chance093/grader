package main

import (
	"fmt"

	"github.com/manifoldco/promptui"
)

func selectAssignmentOption(className string) {
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
		addAssignment(className)
	case "Edit assignment":
		editAssignment(className)
	case "Edit grade weights":
		editGradeWeights(className)
	case "Go back":
		selectClass()
	default: // Handles cases not explicitly matched
		fmt.Printf("Prompt failed %v\n", err)
		return
	}
}

func addAssignment(className string) {
	fmt.Printf("Adding assignment to %s\n", className)
}

func editAssignment(className string) {
	fmt.Printf("Editing assignment in %s\n", className)
}

func editGradeWeights(className string) {
	fmt.Printf("Editing grade weights in %s\n", className)
}
