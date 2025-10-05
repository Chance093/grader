package db

import (
	"log"
	"strconv"

	"github.com/Chance093/gradr/types"
)

func (db *DB) GetClassAssignments(className string) (types.Assignments, error) {
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

	rows, err := db.Query(getClassAssignmentsStatement, className)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var assignments types.Assignments

	for rows.Next() {
		var Name string
		var Grade float64
		var Type string

		if err := rows.Scan(&Name, &Grade, &Type); err != nil {
			return nil, err
		}

		assignments = append(assignments, types.Assignment{
			Name:  Name,
			Grade: strconv.FormatFloat(Grade, 'f', 1, 64),
			Type:  Type,
		})
	}
	return assignments, nil
}

func (db *DB) AddAssignment(name, assignmentType, className, totalPoints, correctPoints string) error {
	// take the class name and type name, and find their id's in the db
	const sqlQueryIdsStatement = `
    SELECT c.id, t.id
    FROM classes c
    CROSS JOIN assignment_types t
    WHERE c.name = ? AND t.name = ?;
    `
	var classID, typeID int
	if err := db.QueryRow(sqlQueryIdsStatement, className, assignmentType).Scan(&classID, &typeID); err != nil {
		log.Fatal(err)
	}

	// create assignment that is associated with class
	const sqlInsertAssignmentStatement = `
      INSERT INTO assignments (name, type_id, class_id, total, correct)
    VALUES (?, ?, ?, ?, ?);
    `
	if _, err := db.Exec(sqlInsertAssignmentStatement, name, typeID, classID, totalPoints, correctPoints); err != nil {
		return err
	}

	return nil
}

func (db *DB) GetAllClassAssignments(className string) ([]string, error) {
	const sqlQueryAssignmentsStatement = `SELECT name FROM assignments WHERE class_id=(
      SELECT id FROM classes WHERE name = ?
    );`

	rows, err := db.Query(sqlQueryAssignmentsStatement, className)
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

func (db *DB) EditAssignmentName(oldName, newName, className string) error {
	const sqlUpdateAssignmentNameStatement = `
      UPDATE assignments SET name = ? WHERE name = ? AND 
    class_id = (SELECT id FROM classes WHERE name = ?);
    `

	if _, err := db.Exec(sqlUpdateAssignmentNameStatement, newName, oldName, className); err != nil {
		return err
	}

	return nil
}

func (db *DB) EditAssignmentGrade(assignment, className, total, correct string) error {
	const sqlUpdateAssignmentGradeStatement = `
      UPDATE assignments SET correct = ?, total = ? WHERE name = ? AND 
    class_id = (SELECT id FROM classes WHERE name = ?);
    `

	if _, err := db.Exec(sqlUpdateAssignmentGradeStatement, correct, total, assignment, className); err != nil {
		return err
	}

	return nil
}

func (db *DB) EditAssignmentType(assignment, className, assignmentType string) error {
	const sqlUpdateAssignmentNameStatement = `
      UPDATE assignments SET type_id = ? WHERE name = ? AND 
    class_id = (SELECT id FROM classes WHERE name = ?);
    `

	typeMap := map[string]int{
		"Test":     1,
		"Quiz":     2,
		"Homework": 3,
	}

	if _, err := db.Exec(sqlUpdateAssignmentNameStatement, typeMap[assignmentType], assignment, className); err != nil {
		return err
	}

	return nil
}

func (db *DB) DeleteAssignment(assignmentName, className string) error {
	const sqlDeleteAssignmentStatement = `
      DELETE FROM assignments WHERE name = ? AND class_id = 
    (SELECT id FROM classes WHERE name = ?);
    `

	if _, err := db.Exec(sqlDeleteAssignmentStatement, assignmentName, className); err != nil {
		return err
	}

	return nil
}
