package db

import (
	"log"
	"strconv"

	"github.com/Chance093/gradr/types"
)

// TODO: Make this file its own db package

func (db *DB) GetClassesAndGradesFromDB() (types.ClassesAndGradesRaw, error) {
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

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var classesAndGrades types.ClassesAndGradesRaw

	for rows.Next() {
		var ClassName string
		var Grade float64
		var Weight int

		if err := rows.Scan(&ClassName, &Grade, &Weight); err != nil {
			return nil, err
		}

		classesAndGrades = append(classesAndGrades, types.ClassAndGradeRaw{
			ClassName,
			Grade,
			Weight,
		})
	}

	return classesAndGrades, nil
}

func (db *DB) GetAllClassesFromDB() ([]string, error) {
	const sqlQueryClassesStatement = `SELECT name FROM classes;`

	rows, err := db.Query(sqlQueryClassesStatement)
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

func (db *DB) AddClassToDB(className, subject string) error {
	const sqlInsertClassStatement = `
      INSERT INTO classes (name, subject)
    VALUES (?, ?);
    `

	res, err := db.Exec(sqlInsertClassStatement, className, subject)
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

	if _, err := db.Exec(sqlCreateWeightsStatement, classId, classId, classId); err != nil {
		return err
	}

	return nil
}

func (db *DB) GetClassAssignmentsFromDB(className string) (types.Assignments, error) {
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

func (db *DB) AddAssignmentToDB(name, assignmentType, className, totalPoints, correctPoints string) error {
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

func (db *DB) GetAllClassAssignmentsFromDB(className string) ([]string, error) {
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

func (db *DB) EditAssignmentNameInDB(oldName, newName, className string) error {
	const sqlUpdateAssignmentNameStatement = `
      UPDATE assignments SET name = ? WHERE name = ? AND 
    class_id = (SELECT id FROM classes WHERE name = ?);
    `

	if _, err := db.Exec(sqlUpdateAssignmentNameStatement, newName, oldName, className); err != nil {
		return err
	}

	return nil
}

func (db *DB) EditAssignmentGradeInDB(assignment, className, total, correct string) error {
	const sqlUpdateAssignmentGradeStatement = `
      UPDATE assignments SET correct = ?, total = ? WHERE name = ? AND 
    class_id = (SELECT id FROM classes WHERE name = ?);
    `

	if _, err := db.Exec(sqlUpdateAssignmentGradeStatement, correct, total, assignment, className); err != nil {
		return err
	}

	return nil
}

func (db *DB) EditAssignmentTypeInDB(assignment, className, assignmentType string) error {
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

func (db *DB) DeleteAssignmentFromDB(assignmentName, className string) error {
	const sqlDeleteAssignmentStatement = `
      DELETE FROM assignments WHERE name = ? AND class_id = 
    (SELECT id FROM classes WHERE name = ?);
    `

	if _, err := db.Exec(sqlDeleteAssignmentStatement, assignmentName, className); err != nil {
		return err
	}

	return nil
}

func (db *DB) DeleteClassFromDB(className string) {
	// take the class name and find the class id in the db
	const sqlDeleteClassStatement = `DELETE FROM classes WHERE name=?`
	if _, err := db.Exec(sqlDeleteClassStatement, className); err != nil {
		log.Fatal(err)
	}
}

func (db *DB) EditClassNameInDB(oldClassName, newClassName string) {
	// take the class name and find the class id in the db
	const sqlUpdateClassStatement = `UPDATE classes SET name = ? WHERE name = ?`
	if _, err := db.Exec(sqlUpdateClassStatement, newClassName, oldClassName); err != nil {
		log.Fatal(err)
	}
}

func (db *DB) GetClassWeights(className string) ([]types.AssignmentWeight, error) {
	const sqlGetClassWeightsStatement = `SELECT weight, type_id FROM assignment_weights
    WHERE class_id = (SELECT id FROM classes WHERE name = ?);
    `

	rows, err := db.Query(sqlGetClassWeightsStatement, className)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var weights []types.AssignmentWeight

	for rows.Next() {
		var Weight int
		var Type_id int

		if err := rows.Scan(&Weight, &Type_id); err != nil {
			return nil, err
		}

		newWeight := types.AssignmentWeight{Weight, Type_id}

		weights = append(weights, newWeight)
	}

	return weights, nil
}

func (db *DB) UpdateClassWeights(className, test, quiz, homework string) error {
	const sqlGetClassIDStatement = `SELECT id FROM classes WHERE name = ?`
	var classID int
	if err := db.QueryRow(sqlGetClassIDStatement, className).Scan(&classID); err != nil {
		log.Fatal(err)
	}

	const sqlUpdateWeightsStatement = `
      UPDATE assignment_weights SET weight = ? WHERE class_id = ? AND type_id = 1;
      UPDATE assignment_weights SET weight = ? WHERE class_id = ? AND type_id = 2;
      UPDATE assignment_weights SET weight = ? WHERE class_id = ? AND type_id = 3;
    `

	if _, err := db.Exec(sqlUpdateWeightsStatement, test, classID, quiz, classID, homework, classID); err != nil {
		log.Fatal(err)
	}

	return nil
}
