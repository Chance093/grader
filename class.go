package main

import (
	"fmt"
	"log"
	"strconv"

	"github.com/manifoldco/promptui"
)

type ClassAndGrade = struct {
	className string
	grade     float64
}

type ClassesAndGrades = []ClassAndGrade

func (cfg *apiConfig) viewOverallGrades() {
	raw, err := cfg.getClassesAndGradesFromDB()
	if err != nil {
		log.Fatalf("Error while getting classes and grades: %s", err.Error())
	}

	classesAndGrades := make(map[string][]float64)
	for _, dat := range raw {
		value, ok := classesAndGrades[dat.className]
		if ok {
			classesAndGrades[dat.className] = append(value, dat.grade)
		} else {
			classesAndGrades[dat.className] = []float64{dat.grade}
		}
	}

	calculated := make(map[string]string)
	for key, val := range classesAndGrades {
		var sum float64
		for _, grade := range val {
			sum += grade
		}

    grade := sum / float64(len(val))

		calculated[key] = strconv.FormatFloat(grade, 'f', 1, 64)
	}

  getClassGradesAscii(calculated)

  cfg.startUpQuestion()
}

func (cfg *apiConfig) getClassesAndGradesFromDB() (ClassesAndGrades, error) {
	query := `
  SELECT 
    classes.name AS class_name, 
    assignments.percentage AS assignment_grade 
  FROM classes
  INNER JOIN assignments
    ON assignments.class_id = classes.id;
  `

	rows, err := cfg.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var classesAndGrades ClassesAndGrades

	for rows.Next() {
		var className string
		var grade float64

		if err := rows.Scan(&className, &grade); err != nil {
			return nil, err
		}

		classesAndGrades = append(classesAndGrades, ClassAndGrade{
			className,
			grade,
		})
	}

	return classesAndGrades, nil
}

func (cfg *apiConfig) addClass() {
	prompt := promptui.Prompt{
		Label: "Enter class name",
	}

	className, err := prompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return
	}

	prompt2 := promptui.Prompt{
		Label: "Enter subject (e.g. Math)",
	}

	subject, err := prompt2.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return
	}

	if err := cfg.addClassToDB(className, subject); err != nil {
		log.Fatalf("Failed to add to db: %s", err.Error())
	}

	fmt.Printf("%s added!\n", className)

	cfg.startUpQuestion()
}

func (cfg *apiConfig) addClassToDB(className, subject string) error {
	const sqlInsertClassStatement = `
      INSERT INTO classes (name, subject)
    VALUES (?, ?);
    `

	if _, err := cfg.db.Exec(sqlInsertClassStatement, className, subject); err != nil {
		return err
	}

	return nil
}

func (cfg *apiConfig) selectClass() {
	classes, err := cfg.getAllClassesFromDB()
	if err != nil {
		log.Fatalf("Error while getting classes: %s", err.Error())
	}

	classes = append(classes, "Go Back")

	prompt := promptui.Select{
		Label: "Select a Class",
		Items: classes,
	}

	_, result, err := prompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return
	}

	if result == "Go Back" {
		cfg.startUpQuestion()
	}

	cfg.selectAssignmentOption(result)
}

type Classes []struct {
	id      int
	name    string
	subject string
}

func (cfg *apiConfig) getAllClassesFromDB() ([]string, error) {
	const sqlQueryClassesStatement = `SELECT name FROM classes;`

	rows, err := cfg.db.Query(sqlQueryClassesStatement)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var classes []string

	for rows.Next() {
		var name string

		if err := rows.Scan(&name); err != nil {
			return nil, err
		}

		classes = append(classes, name)
	}

	return classes, nil
}
