package main

import (
	"fmt"
	"os"

	"github.com/manifoldco/promptui"
)

func main() {
	startUpQuestion()
}

func viewOverallGrades() {
}

func startUpQuestion() {
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

func quit() {
	fmt.Println("Quitting...")
	os.Exit(0)
}
