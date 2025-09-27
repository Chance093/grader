package main

import (
	"fmt"
	"log"
	"strconv"

	"github.com/manifoldco/promptui"
)

type ClassAndGradeRaw = struct {
	className string
	grade     float64
	weight    int
}

type ClassesAndGradesRaw = []ClassAndGradeRaw

type GradeAndWeight = struct {
	grade  float64
	weight int
}
type GradesAndWeights = []GradeAndWeight

type ClassMap = map[string]map[int][]float64

func (cfg *apiConfig) viewOverallGrades() {
	raw, err := cfg.getClassesAndGradesFromDB()
	if err != nil {
		log.Fatalf("Error while getting classes and grades: %s", err.Error())
	}

	calculated := calculateGrades(raw)

	getClassGradesAscii(calculated)

	cfg.startUpQuestion()
}

func calculateGrades(raw ClassesAndGradesRaw) map[string]string {
	classesAndGrades := make(ClassMap)
	for _, dat := range raw {
		value, ok := classesAndGrades[dat.className]
		if ok {
			innerVal, innerOk := value[dat.weight]
			if innerOk {
				classesAndGrades[dat.className][dat.weight] = append(innerVal, dat.grade)
			} else {
				classesAndGrades[dat.className][dat.weight] = []float64{dat.grade}
			}
		} else {
			classesAndGrades[dat.className] = map[int][]float64{
				dat.weight: {dat.grade},
			}
		}
	}

	calculated := make(map[string]string)
	for className, weightGradeMap := range classesAndGrades {
		var totalPercentage float64

		var totalWeight int
		for weight := range weightGradeMap {
			totalWeight += weight
		}

		for weight, grades := range weightGradeMap {
			var sum float64
			for _, grade := range grades {
				sum += grade
			}

			total := sum / float64(len(grades))
			newWeight := float64(weight) / float64(totalWeight)
			percent := total * newWeight

			totalPercentage += percent
		}

		calculated[className] = strconv.FormatFloat(totalPercentage, 'f', 1, 64)
	}

	return calculated
}

func (cfg *apiConfig) getClassesAndGradesFromDB() (ClassesAndGradesRaw, error) {
	query := `
  SELECT 
    classes.name AS class_name, 
    assignments.percentage AS assignment_grade,
    assignment_weights.weight AS assignment_weight
  FROM classes
  INNER JOIN assignments
    ON assignments.class_id = classes.id
  INNER JOIN assignment_weights
    ON assignment_weights.class_id = classes.id 
    AND assignment_weights.type_id = assignments.type_id;
  `

	rows, err := cfg.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var classesAndGrades ClassesAndGradesRaw

	for rows.Next() {
		var className string
		var grade float64
		var weight int

		if err := rows.Scan(&className, &grade, &weight); err != nil {
			return nil, err
		}

		classesAndGrades = append(classesAndGrades, ClassAndGradeRaw{
			className,
			grade,
			weight,
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

	res, err := cfg.db.Exec(sqlInsertClassStatement, className, subject)
	if err != nil {
		return err
	}

	classId, err := res.LastInsertId()
	if err != nil {
		return err
	}

	const sqlCreateWeightsStatement = `
    INSERT INTO assignment_weights (weight, type_id, class_id)
    VALUES (50, 1, ?), (20, 2, ?), (30, 3, ?);
    `

	if _, err := cfg.db.Exec(sqlCreateWeightsStatement, classId, classId, classId); err != nil {
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

func (cfg *apiConfig) deleteClass() {
	classes, err := cfg.getAllClassesFromDB()
	if err != nil {
		log.Fatalf("Error while getting classes : %s", err.Error())
	}

	classes = append(classes, "Go Back")

	prompt := promptui.Select{
		Label: "Select a class to delete",
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

	cfg.deleteClassFromDB(result)

	fmt.Printf("Deleted %s!\n", result)

	cfg.startUpQuestion()
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

	classes = append(classes, "Go Back")

	prompt := promptui.Select{
		Label: "Select a class to edit",
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

	prompt2 := promptui.Select{
		Label: "Choose an option",
		Items: []string{"Change class name", "Change grade weights", "Go back"},
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

	cfg.startUpQuestion()
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

	cfg.startUpQuestion()
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
