package main

import (
	"fmt"
	"log"
	"strconv"

	"github.com/Chance093/gradr/ascii"
	"github.com/Chance093/gradr/constants"
	"github.com/Chance093/gradr/prompt"
	"github.com/manifoldco/promptui"
)

type GradeAndWeight = struct {
	grade  float64
	weight int
}
type GradesAndWeights = []GradeAndWeight

func (cfg *apiConfig) viewOverallGrades() {
	raw, err := cfg.getClassesAndGradesFromDB()
	if err != nil {
		log.Fatalf("Error while getting classes and grades: %s", err.Error())
	}

	calculated := calculateGrades(raw)

	// for the case of classes with no assignments
	if len(calculated) == 0 {
		classes, err := cfg.getAllClassesFromDB()
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

func (cfg *apiConfig) selectClass() {
	classes, err := cfg.getAllClassesFromDB()
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

func (cfg *apiConfig) addClass() {
	className, err := prompt.Input("Enter class name")
	if err != nil {
		log.Fatalf("Prompt failed %v\n", err)
	}

	subject, err := prompt.Input("Enter subject (e.g. Math)")
	if err != nil {
		log.Fatalf("Prompt failed %v\n", err)
	}

	if err := cfg.addClassToDB(className, subject); err != nil {
		log.Fatalf("Failed to add to db: %s", err.Error())
	}

	fmt.Printf("%s added!\n", className)

	cfg.displayMainMenu()
}

type Classes []struct {
	id      int
	name    string
	subject string
}

func (cfg *apiConfig) deleteClass() {
	classes, err := cfg.getAllClassesFromDB()
	if err != nil {
		log.Fatalf("Error while getting classes : %s", err.Error())
	}

	classes = append(classes, "Main Menu")

	prompt := promptui.Select{
		Label: "Select a class to delete",
		Items: classes,
	}

	_, result, err := prompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return
	}

	if result == "Main Menu" {
		cfg.displayMainMenu()
	}

	cfg.deleteClassFromDB(result)

	fmt.Printf("Deleted %s!\n", result)

	cfg.displayMainMenu()
}

func (cfg *apiConfig) deleteClassFromDB(className string) {
	// take the class name and find the class id in the db
	const sqlDeleteClassStatement = `DELETE FROM classes WHERE name=?`
	if _, err := cfg.db.Exec(sqlDeleteClassStatement, className); err != nil {
		log.Fatal(err)
	}
}

func (cfg *apiConfig) editClass() {
	classes, err := cfg.getAllClassesFromDB()
	if err != nil {
		log.Fatalf("Error while getting classes : %s", err.Error())
	}

	classes = append(classes, "Main Menu")

	prompt := promptui.Select{
		Label: "Select a class to edit",
		Items: classes,
	}

	_, result, err := prompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return
	}

	if result == "Main Menu" {
		cfg.displayMainMenu()
	}

	prompt2 := promptui.Select{
		Label: "Choose an option",
		Items: []string{"Change class name", "Change grade weights", "Go back", "Main Menu"},
	}

	_, result2, err := prompt2.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return
	}

	switch result2 {
	case "Change class name":
		cfg.editClassName(result)
	case "Change grade weights":
		cfg.editClassWeights(result)
	case "Go back":
		cfg.editClass()
	case "Main Menu":
		cfg.displayMainMenu()
	default: // Handles cases not explicitly matched
		fmt.Printf("Prompt failed %v\n", err)
		return
	}
}

func (cfg *apiConfig) editClassName(oldClassName string) {
	classNamePrompt := promptui.Prompt{
		Label: "Enter new class name",
	}
	newClassName, err := classNamePrompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return
	}

	cfg.editClassNameInDB(oldClassName, newClassName)

	fmt.Printf("Changed class name '%s' to '%s'!\n", oldClassName, newClassName)

	cfg.displayMainMenu()
}

func (cfg *apiConfig) editClassNameInDB(oldClassName, newClassName string) {
	// take the class name and find the class id in the db
	const sqlUpdateClassStatement = `UPDATE classes SET name = ? WHERE name = ?`
	if _, err := cfg.db.Exec(sqlUpdateClassStatement, newClassName, oldClassName); err != nil {
		log.Fatal(err)
	}
}

func (cfg *apiConfig) editClassWeights(className string) {
	typeMap := map[int]string{
		1: "Test",
		2: "Quiz",
		3: "Homework",
	}

	currentWeights, err := cfg.getClassWeights(className)
	if err != nil {
		log.Fatal(err)
	}

	testWeightPrompt := promptui.Prompt{
		Label: fmt.Sprintf("Enter new test weight (Current: %s-%d, %s-%d, %s-%d)", typeMap[currentWeights[0].type_id], currentWeights[0].weight, typeMap[currentWeights[1].type_id], currentWeights[1].weight, typeMap[currentWeights[2].type_id], currentWeights[2].weight),
	}
	testWeight, err := testWeightPrompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return
	}

	quizWeightPrompt := promptui.Prompt{
		Label: fmt.Sprintf("Enter new quiz weight (Current: %s-%d, %s-%d, %s-%d)", typeMap[currentWeights[0].type_id], currentWeights[0].weight, typeMap[currentWeights[1].type_id], currentWeights[1].weight, typeMap[currentWeights[2].type_id], currentWeights[2].weight),
	}
	quizWeight, err := quizWeightPrompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return
	}

	homeworkWeightPrompt := promptui.Prompt{
		Label: fmt.Sprintf("Enter new homework weight (Current: %s-%d, %s-%d, %s-%d)", typeMap[currentWeights[0].type_id], currentWeights[0].weight, typeMap[currentWeights[1].type_id], currentWeights[1].weight, typeMap[currentWeights[2].type_id], currentWeights[2].weight),
	}
	homeworkWeight, err := homeworkWeightPrompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return
	}

	testWeightInt, err := strconv.Atoi(testWeight)
	if err != nil {
		fmt.Printf("String to int conversion failed %v\n", err)
		return
	}

	quizWeightInt, err := strconv.Atoi(quizWeight)
	if err != nil {
		fmt.Printf("String to int conversion failed %v\n", err)
		return
	}

	homeworkWeightInt, err := strconv.Atoi(homeworkWeight)
	if err != nil {
		fmt.Printf("String to int conversion failed %v\n", err)
		return
	}

	if testWeightInt+quizWeightInt+homeworkWeightInt != 100 {
		fmt.Println("Test, Quiz, and Homework weights should add to 100. Please try again.")
		cfg.editClassWeights(className)
	}

	if err := cfg.updateClassWeights(className, testWeightInt, quizWeightInt, homeworkWeightInt); err != nil {
		log.Fatal(err)
	}

	cfg.displayMainMenu()
}

type AssignmentWeight = struct {
	weight  int
	type_id int
}

func (cfg *apiConfig) getClassWeights(className string) ([]AssignmentWeight, error) {
	const sqlGetClassWeightsStatement = `SELECT weight, type_id FROM assignment_weights
    WHERE class_id = (SELECT id FROM classes WHERE name = ?);
    `

	rows, err := cfg.db.Query(sqlGetClassWeightsStatement, className)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var weights []AssignmentWeight

	for rows.Next() {
		var weight int
		var type_id int

		if err := rows.Scan(&weight, &type_id); err != nil {
			return nil, err
		}

		newWeight := AssignmentWeight{weight, type_id}

		weights = append(weights, newWeight)
	}

	return weights, nil
}

func (cfg *apiConfig) updateClassWeights(className string, test, quiz, homework int) error {
	const sqlGetClassIDStatement = `SELECT id FROM classes WHERE name = ?`
	var classID int
	if err := cfg.db.QueryRow(sqlGetClassIDStatement, className).Scan(&classID); err != nil {
		log.Fatal(err)
	}

	const sqlUpdateWeightsStatement = `
      UPDATE assignment_weights SET weight = ? WHERE class_id = ? AND type_id = 1;
      UPDATE assignment_weights SET weight = ? WHERE class_id = ? AND type_id = 2;
      UPDATE assignment_weights SET weight = ? WHERE class_id = ? AND type_id = 3;
    `

	if _, err := cfg.db.Exec(sqlUpdateWeightsStatement, test, classID, quiz, classID, homework, classID); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Successfully upgraded class weights!")

	return nil
}
