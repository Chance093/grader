package main

import (
	"log"
	"strconv"

	"github.com/Chance093/gradr/types"
)

type ClassAndGradeRaw = struct {
	className string
	grade     float64
	weight    int
}

type ClassesAndGradesRaw = []ClassAndGradeRaw

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

func (cfg *apiConfig) getClassAssignmentsFromDB(className string) (types.Assignments, error) {
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

func (cfg *apiConfig) getAllClassAssignmentsFromDB(className string) ([]string, error) {
	const sqlQueryAssignmentsStatement = `SELECT name FROM assignments WHERE class_id=(
      SELECT id FROM classes WHERE name = ?
    );`

	rows, err := cfg.db.Query(sqlQueryAssignmentsStatement, className)
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
