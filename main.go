package main

import (
	"fmt"

	"github.com/manifoldco/promptui"
)

func main() {
	prompt := promptui.Select{
		Label: "Choose an option",
		Items: []string{"Add a Class", "Select a Class", "View Overall Grades", "Quit"},
	}

	_, result, err := prompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return
	}

	switch result {
	case "Add a Class":
		addClass()
	case "Select a Class":
		selectClass()
	case "View Overall Grades":
		viewOverallGrades()
	case "Quit":
		quit()
	default: // Handles cases not explicitly matched
		fmt.Printf("Prompt failed %v\n", err)
		return
	}
}

func addClass() {
}

func selectClass() {
	prompt := promptui.Select{
		Label: "Select a Class",
		Items: []string{"Algebra 1", "Calculus", "Go Back"},
	}

	_, result, err := prompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return
	}

	fmt.Printf("You chose %s\n", result)
}

func viewOverallGrades() {
}

func quit() {
}
