package main

import (
	"fmt"
	"log"
	"strconv"

	"github.com/manifoldco/promptui"
)

func (cfg *apiConfig) selectAssignmentOption(className string) {
	prompt := promptui.Select{
		Label: "Choose an option",
		Items: []string{"View assignments", "Add assignment", "Edit assignment", "Delete assignment", "Go back"},
	}

	_, result, err := prompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return
	}

	switch result {
	case "View assignments":
		cfg.viewAssignments(className)
	case "Add assignment":
		cfg.addAssignment(className)
	case "Edit assignment":
		cfg.editAssignment(className)
	case "Delete assignment":
		cfg.deleteAssignment(className)
	case "Go back":
		cfg.selectClass()
	default: // Handles cases not explicitly matched
		fmt.Printf("Prompt failed %v\n", err)
		return
	}
}

func (cfg *apiConfig) viewAssignments(className string) {
	raw, err := cfg.getClassAssignmentsFromDB(className)
	if err != nil {
		log.Fatalf("Error while getting class assignments: %s", err.Error())
	}

	getAssignmentGradesAscii(raw)

	cfg.selectAssignmentOption(className)
}

type AssignmentsRaw struct {
	assignment     string
	grade          string
	assignmentType string
}

func (cfg *apiConfig) getClassAssignmentsFromDB(className string) ([]AssignmentsRaw, error) {
	const getClassAssignmentsStatement = `
  SELECT assignments.name, 
    assignments.percentage AS grade, 
    assignment_types.name AS type 
  FROM assignments
  INNER JOIN assignment_types
    ON assignments.type_id = assignment_types.id
  WHERE assignments.class_id = (
    SELECT id FROM classes WHERE name = ?
  );
  `

	rows, err := cfg.db.Query(getClassAssignmentsStatement, className)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var assignments []AssignmentsRaw

	for rows.Next() {
		var assignment string
		var grade float64
		var assignmentType string

		if err := rows.Scan(&assignment, &grade, &assignmentType); err != nil {
			return nil, err
		}

		assignments = append(assignments, AssignmentsRaw{
			assignment:     assignment,
			grade:          strconv.FormatFloat(grade, 'f', 1, 64),
			assignmentType: assignmentType,
		})
	}
	return assignments, nil
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

	cfg.validatePoints(totalPoints, correctPoints, className)

	if err := cfg.addAssignmentToDB(assignmentName, assignmentType, className, totalPoints, correctPoints); err != nil {
		log.Fatalf("Error while adding assignment to db: %s", err.Error())
	}

	fmt.Printf("%s: %s has been added!\n", assignmentType, assignmentName)

	cfg.selectAssignmentOption(className)
}

func (cfg *apiConfig) validatePoints(totalPoints, correctPoints, className string) {
	totalPointsInt, err := strconv.Atoi(totalPoints)
	if err != nil {
		log.Fatalf("Failed to convert string to int: %s", err.Error())
	}

	correctPointsInt, err := strconv.Atoi(correctPoints)
	if err != nil {
		log.Fatalf("Failed to convert string to int: %s", err.Error())
	}

	if totalPointsInt < correctPointsInt {
		fmt.Println("Total points must be greater than or equal to correct points")
		cfg.addAssignment(className)
	}
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

func (cfg *apiConfig) editAssignment(className string) {
	assignments, err := cfg.getAllClassAssignmentsFromDB(className)
	if err != nil {
		log.Fatalf("Error while getting class assignments: %s", err.Error())
	}

	assignments = append(assignments, "Go Back")

	prompt := promptui.Select{
		Label: "Select an assignment to edit",
		Items: assignments,
	}

	_, assignment, err := prompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return
	}

	if assignment == "Go Back" {
		cfg.selectAssignmentOption(className)
	}

	editPrompt := promptui.Select{
		Label: "Choose an option to edit",
		Items: []string{"Name", "Grade", "Assignment Type"},
	}

	_, result, err := editPrompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return
	}

	switch result {
	case "Name":
		cfg.editAssignmentName(assignment, className)
	case "Grade":
		cfg.editAssignmentGrade(assignment, className)
	case "Assignment Type":
		cfg.editAssignmentType(assignment, className)
	default: // Handles cases not explicitly matched
		fmt.Printf("Prompt failed %v\n", err)
		return
	}

	cfg.selectAssignmentOption(className)
}

func (cfg *apiConfig) getAllClassAssignmentsFromDB(className string) ([]string, error) {
	// take the class name and find the class id in the db
	const sqlQueryClassIdStatement = `SELECT id FROM classes WHERE name=?`
	var classID int
	if err := cfg.db.QueryRow(sqlQueryClassIdStatement, className).Scan(&classID); err != nil {
		log.Fatal(err)
	}

	const sqlQueryAssignmentsStatement = `SELECT name FROM assignments WHERE class_id=?;`

	rows, err := cfg.db.Query(sqlQueryAssignmentsStatement, classID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var assignments []string

	for rows.Next() {
		var name string

		if err := rows.Scan(&name); err != nil {
			return nil, err
		}

		assignments = append(assignments, name)
	}

	return assignments, nil
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

	assignments = append(assignments, "Go Back")

	prompt := promptui.Select{
		Label: "Select an assignment to delete",
		Items: assignments,
	}

	_, result, err := prompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return
	}

	if result == "Go Back" {
		cfg.startUpQuestion()
	}

	cfg.deleteAssignmentFromDB(result, className)

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
