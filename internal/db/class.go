package db

import (
	"fmt"

	"github.com/Chance093/grader/types"
)

func (db *DB) GetClassesAndGrades() (types.ClassesAndGradesRaw, error) {
	const getClassesAndGradesQuery = `
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

	rows, err := db.Query(getClassesAndGradesQuery)
	if err != nil {
		return nil, fmt.Errorf("Error querying classes and grades: %w", err)
	}
	defer rows.Close()

	var classesAndGrades types.ClassesAndGradesRaw

	for rows.Next() {
		var cg types.ClassAndGradeRaw

		if err := rows.Scan(&cg.ClassName, &cg.Grade, &cg.Weight); err != nil {
			return nil, fmt.Errorf("Error scanning row: %w", err)
		}

		classesAndGrades = append(classesAndGrades, cg)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("Error during iteration: %w", err)
	}

	return classesAndGrades, nil
}

func (db *DB) GetAllClasses() ([]string, error) {
	const getAllClassesQuery = `SELECT name FROM classes;`

	rows, err := db.Query(getAllClassesQuery)
	if err != nil {
		return nil, fmt.Errorf("Error querying classes: %w", err)
	}
	defer rows.Close()

	var classes []string

	for rows.Next() {
		var name string

		if err := rows.Scan(&name); err != nil {
			return nil, fmt.Errorf("Error scanning row: %w", err)
		}

		classes = append(classes, name)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("Error during iteration: %w", err)
	}

	return classes, nil
}

func (db *DB) AddClass(className, subject string) error {
	const addClassStatement = `
  INSERT INTO classes (name, subject)
  VALUES (?, ?);
  `

	res, err := db.Exec(addClassStatement, className, subject)
	if err != nil {
		return fmt.Errorf("Error inserting class: %w", err)
	}

	classId, err := res.LastInsertId()
	if err != nil {
		return fmt.Errorf("Error getting insert id: %w", err)
	}

	// Test-50% -- Quiz-20% -- Homework-30%
	const createDefaultWeightsStatement = `
  INSERT INTO assignment_weights (weight, type_id, class_id)
  VALUES (50, 1, ?), (20, 2, ?), (30, 3, ?);
  `

	if _, err := db.Exec(createDefaultWeightsStatement, classId, classId, classId); err != nil {
		return fmt.Errorf("Error inserting default weights: %w", err)
	}

	return nil
}

func (db *DB) DeleteClass(className string) error {
	const deleteClassStatement = `DELETE FROM classes WHERE name=?`

	if _, err := db.Exec(deleteClassStatement, className); err != nil {
		return fmt.Errorf("Error deleting class: %w", err)
	}

	return nil
}

func (db *DB) EditClassName(oldClassName, newClassName string) error {
	const updateClassNameStatement = `UPDATE classes SET name = ? WHERE name = ?`

	if _, err := db.Exec(updateClassNameStatement, newClassName, oldClassName); err != nil {
		return fmt.Errorf("Error updating class name: %w", err)
	}

	return nil
}

func (db *DB) GetClassWeights(className string) ([]types.AssignmentWeight, error) {
	const getClassWeightsQuery = `
  SELECT weight, type_id FROM assignment_weights
  WHERE class_id = (SELECT id FROM classes WHERE name = ?);
  `

	rows, err := db.Query(getClassWeightsQuery, className)
	if err != nil {
		return nil, fmt.Errorf("Error getting class weights: %w", err)
	}
	defer rows.Close()

	var weights []types.AssignmentWeight

	for rows.Next() {
		var aw types.AssignmentWeight

		if err := rows.Scan(&aw.Weight, &aw.Type_id); err != nil {
			return nil, fmt.Errorf("Error scanning row: %w", err)
		}

		weights = append(weights, aw)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("Error during iteration: %w", err)
	}

	return weights, nil
}

func (db *DB) UpdateClassWeights(className, test, quiz, homework string) error {
	const getClassIDQuery = `SELECT id FROM classes WHERE name = ?`

	var classID int
	if err := db.QueryRow(getClassIDQuery, className).Scan(&classID); err != nil {
		return fmt.Errorf("Error querying class id and scanning row: %w", err)
	}

	const updateClassWeightsStatement = `
  UPDATE assignment_weights SET weight = ? WHERE class_id = ? AND type_id = 1;
  UPDATE assignment_weights SET weight = ? WHERE class_id = ? AND type_id = 2;
  UPDATE assignment_weights SET weight = ? WHERE class_id = ? AND type_id = 3;
  `

	if _, err := db.Exec(updateClassWeightsStatement, test, classID, quiz, classID, homework, classID); err != nil {
		return fmt.Errorf("Error updating class weights: %w", err)
	}

	return nil
}
