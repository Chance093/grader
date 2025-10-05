package db

import (
	"fmt"

	"github.com/Chance093/grader/types"
)

func (db *DB) GetClassAssignments(className string) (types.Assignments, error) {
	const getClassAssignmentsQuery = `
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

	rows, err := db.Query(getClassAssignmentsQuery, className)
	if err != nil {
		return nil, fmt.Errorf("Error querying class assignments: %w", err)
	}
	defer rows.Close()

	var assignments types.Assignments

	for rows.Next() {
		var a types.Assignment

		if err := rows.Scan(&a.Name, &a.Grade, &a.Type); err != nil {
			return nil, fmt.Errorf("Error scanning row: %w", err)
		}

		assignments = append(assignments, a)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("Error during iteration: %w", err)
	}

	return assignments, nil
}

func (db *DB) AddAssignment(name, assignmentType, className, totalPoints, correctPoints string) error {
	const getClassAndTypeIDsQuery = `
  SELECT c.id, t.id
  FROM classes c
  CROSS JOIN assignment_types t
  WHERE c.name = ? AND t.name = ?;
  `

	var classID, typeID int
	if err := db.QueryRow(getClassAndTypeIDsQuery, className, assignmentType).Scan(&classID, &typeID); err != nil {
		return fmt.Errorf("Error querying id's and scanning row: %w", err)
	}

	const createClassAssignmentStatement = `
  INSERT INTO assignments (name, type_id, class_id, total, correct)
  VALUES (?, ?, ?, ?, ?);
  `

	if _, err := db.Exec(createClassAssignmentStatement, name, typeID, classID, totalPoints, correctPoints); err != nil {
		return fmt.Errorf("Error creating class assignment: %w", err)
	}

	return nil
}

func (db *DB) GetAllClassAssignments(className string) ([]string, error) {
	const getAllClassAssignmentsQuery = `
  SELECT name FROM assignments WHERE class_id=(
    SELECT id FROM classes WHERE name = ?
  );
  `

	rows, err := db.Query(getAllClassAssignmentsQuery, className)
	if err != nil {
    return nil, fmt.Errorf("Error querying all class assignments: %w", err)
	}
	defer rows.Close()

	var assignments []string

	for rows.Next() {
		var name string

		if err := rows.Scan(&name); err != nil {
      return nil, fmt.Errorf("Error scanning row: %w", err)
		}

		assignments = append(assignments, name)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("Error during iteration: %w", err)
	}

	return assignments, nil
}

func (db *DB) EditAssignmentName(oldName, newName, className string) error {
	const updateAssignmentNameStatement = `
  UPDATE assignments 
  SET name = ?
  WHERE name = ? AND class_id = (
    SELECT id FROM classes WHERE name = ?
  );
  `

	if _, err := db.Exec(updateAssignmentNameStatement, newName, oldName, className); err != nil {
    return fmt.Errorf("Error updating assignment name: %w", err)
	}

	return nil
}

func (db *DB) EditAssignmentGrade(assignment, className, total, correct string) error {
	const updateAssignmentGradeStatement = `
  UPDATE assignments 
  SET correct = ?, total = ? 
  WHERE name = ? AND class_id = (
    SELECT id FROM classes WHERE name = ?
  );
  `

	if _, err := db.Exec(updateAssignmentGradeStatement, correct, total, assignment, className); err != nil {
    return fmt.Errorf("Error updating assignment grade: %w", err)
	}

	return nil
}

func (db *DB) EditAssignmentType(assignment, className, assignmentType string) error {
	const updateAssignmentTypeStatement = `
  UPDATE assignments 
  SET type_id = ? 
  WHERE name = ? AND class_id = (
    SELECT id FROM classes WHERE name = ?
  );
  `

	typeMap := map[string]int{
		"Test":     1,
		"Quiz":     2,
		"Homework": 3,
	}

	if _, err := db.Exec(updateAssignmentTypeStatement, typeMap[assignmentType], assignment, className); err != nil {
    return fmt.Errorf("Error updating assignment type: %w", err)
	}

	return nil
}

func (db *DB) DeleteAssignment(assignmentName, className string) error {
	const deleteAssignmentStatement = `
  DELETE FROM assignments 
  WHERE name = ? AND class_id = (
    SELECT id FROM classes WHERE name = ?
  );
  `

	if _, err := db.Exec(deleteAssignmentStatement, assignmentName, className); err != nil {
    return fmt.Errorf("Error deleting assignment: %w", err)
	}

	return nil
}
