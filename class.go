package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/Chance093/gradr/ascii"
	"github.com/Chance093/gradr/calculate"
	"github.com/Chance093/gradr/constants"
	"github.com/Chance093/gradr/prompt"
	"github.com/Chance093/gradr/validation"
)

func (cfg *config) viewOverallGrades() {
	raw, err := cfg.db.GetClassesAndGrades()
	if err != nil {
		log.Fatalf("Error while getting classes and grades: %s", err.Error())
	}

	calculated := calculate.CalculateGrades(raw)

	// for the case of classes with no assignments
	if len(calculated) == 0 {
		classes, err := cfg.db.GetAllClasses()
		if err != nil {
			log.Fatalf("Error while getting classes: %s", err.Error())
		}

		for _, class := range classes {
			calculated[class] = " N/A"
		}
	}

	ascii.DisplayClassGrades(calculated)

	cfg.displayMainMenu()
}

func (cfg *config) selectClass() {
	classes, err := cfg.db.GetAllClasses()
	if err != nil {
		log.Fatalf("Error while getting classes: %s", err.Error())
	}

	classes = append(classes, constants.MAIN_MENU)

	result, err := prompt.List(constants.SELECT_CLASS, classes)
	if err != nil {
		log.Fatalf("Prompt failed %v\n", err)
	}

	if result == constants.MAIN_MENU {
		cfg.displayMainMenu()
	}

	cfg.displayClassMenu(result)
}

func (cfg *config) addClass() {
	className, err := prompt.Input(constants.ENTER_CLASS_NAME)
	if err != nil {
		log.Fatalf("Prompt failed %v\n", err)
	}

	subject, err := prompt.Input(constants.ENTER_SUBJECT)
	if err != nil {
		log.Fatalf("Prompt failed %v\n", err)
	}

	if err := cfg.db.AddClass(className, subject); err != nil {
		log.Fatalf("Error while adding class: %s", err.Error())
	}

	fmt.Printf("%s added!\n", className)

	cfg.displayMainMenu()
}

func (cfg *config) editClass() {
	classes, err := cfg.db.GetAllClasses()
	if err != nil {
		log.Fatalf("Error while getting classes : %s", err.Error())
	}

	classes = append(classes, constants.MAIN_MENU)
	classes = append(classes, constants.QUIT)

	class, err := prompt.List(constants.SELECT_CLASS_EDIT, classes)
	if err != nil {
		log.Fatalf("Prompt failed %v\n", err)
	}

	switch class {
	case constants.MAIN_MENU:
		cfg.displayMainMenu()
	case constants.QUIT:
		quit()
	}

	result, err := prompt.List(constants.CHOOSE_OPTION_EDIT, constants.EDIT_CLASS_OPTS)
	if err != nil {
		log.Fatalf("Prompt failed %v\n", err)
	}

	switch result {
	case constants.EDIT_CLASS_OPTS[0]:
		cfg.editClassName(result)
	case constants.EDIT_CLASS_OPTS[1]:
		cfg.editClassWeights(result)
	default:
		log.Fatalf("Prompt failed %v\n", err)
	}
}

func (cfg *config) editClassName(oldClassName string) {
	newClassName, err := prompt.Input(constants.ENTER_CLASS_NAME)
	if err != nil {
		log.Fatalf("Prompt failed %v\n", err)
	}

	cfg.db.EditClassName(oldClassName, newClassName)

	fmt.Printf("Changed class name '%s' to '%s'!\n", oldClassName, newClassName)

	cfg.displayMainMenu()
}

func (cfg *config) editClassWeights(className string) {
	typeMap := map[int]string{
		1: "Test",
		2: "Quiz",
		3: "Homework",
	}

	currentWeights, err := cfg.db.GetClassWeights(className)
	if err != nil {
		log.Fatalf("Error while getting class weights: %s", err)
	}

	formattedStr := fmt.Sprintf(" (Current: %s-%d, %s-%d, %s-%d)", typeMap[currentWeights[0].Type_id], currentWeights[0].Weight, typeMap[currentWeights[1].Type_id], currentWeights[1].Weight, typeMap[currentWeights[2].Type_id], currentWeights[2].Weight)
	testWeight, err := prompt.Input(constants.ENTER_TEST_WEIGHT + formattedStr)
	if err != nil {
		log.Fatalf("Prompt failed %v\n", err)
	}
	quizWeight, err := prompt.Input(constants.ENTER_QUIZ_WEIGHT + formattedStr)
	if err != nil {
		log.Fatalf("Prompt failed %v\n", err)
	}
	homeworkWeight, err := prompt.Input(constants.ENTER_HOMEWORK_WEIGHT + formattedStr)
	if err != nil {
		log.Fatalf("Prompt failed %v\n", err)
	}

	if err := validation.ValidateWeights([]string{testWeight, quizWeight, homeworkWeight}); err != nil {
    // check to see if error is conversion error or validation error
		if strings.Contains(err.Error(), "conversion") {
			log.Fatalf("Error while validating weights: %s", err)
		}
		fmt.Println(err.Error())
		cfg.editClassWeights(className)
		return // explicit return so when the callstack clears, it doesn't make a ton of bad assignments
	}

	if err := cfg.db.UpdateClassWeights(className, testWeight, quizWeight, homeworkWeight); err != nil {
		log.Fatalf("Error while updating class weights: %s", err)
	}

	fmt.Println("Successfully upgraded class weights!")

	cfg.displayMainMenu()
}

func (cfg *config) deleteClass() {
	classes, err := cfg.db.GetAllClasses()
	if err != nil {
		log.Fatalf("Error while getting classes : %s", err.Error())
	}

	classes = append(classes, constants.MAIN_MENU)
	classes = append(classes, constants.QUIT)

	class, err := prompt.List(constants.SELECT_CLASS_DELETE, classes)
	if err != nil {
		log.Fatalf("Prompt failed %v\n", err)
	}

	switch class {
	case constants.MAIN_MENU:
		cfg.displayMainMenu()
	case constants.QUIT:
		quit()
	default:
		cfg.db.DeleteClass(class)
	}

	fmt.Printf("Deleted %s!\n", class)

	cfg.displayMainMenu()
}
