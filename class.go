package main

import (
	"fmt"

	"github.com/manifoldco/promptui"
)

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

  if result == "Go Back" {
    startUpQuestion()
  }

  selectAssignmentOption(result)
}
