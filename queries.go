package main

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
