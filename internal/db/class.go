package db

import (
	"log"

	"github.com/Chance093/gradr/types"
)

func (db *DB) GetClassesAndGrades() (types.ClassesAndGradesRaw, error) {
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

func (db *DB) GetAllClasses() ([]string, error) {
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

func (db *DB) AddClass(className, subject string) error {
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

func (db *DB) DeleteClass(className string) {
	// take the class name and find the class id in the db
	const sqlDeleteClassStatement = `DELETE FROM classes WHERE name=?`
	if _, err := db.Exec(sqlDeleteClassStatement, className); err != nil {
		log.Fatal(err)
	}
}

func (db *DB) EditClassName(oldClassName, newClassName string) {
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
