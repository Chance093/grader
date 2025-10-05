package calculate

import (
	"strconv"

	"github.com/Chance093/grader/types"
)

func CalculateGrades(raw types.ClassesAndGradesRaw) types.ClassAndGradeMap {
	classAndWeightGradesMap := getClassAndWeightGradesMap(raw)
	classAndGradeMap := getClassAndGradeMap(classAndWeightGradesMap)

	return classAndGradeMap
}

func getClassAndWeightGradesMap(raw types.ClassesAndGradesRaw) types.ClassAndWeightGradesMap {
	classesAndGrades := make(types.ClassAndWeightGradesMap)
	for _, row := range raw {
		weightGradesMap, classExists := classesAndGrades[row.ClassName]
		if classExists {
			grades, weightExists := weightGradesMap[row.Weight]
			if weightExists {
				// append grade to grades slice associated with weight
				classesAndGrades[row.ClassName][row.Weight] = append(grades, row.Grade)
			} else {
				// init weight and associate grade
				classesAndGrades[row.ClassName][row.Weight] = []float64{row.Grade}
			}
		} else { // else initialize className to classesAndGrades map
			classesAndGrades[row.ClassName] = map[int][]float64{
				row.Weight: {row.Grade},
			}
		}
	}

	return classesAndGrades
}

func getClassAndGradeMap(classAndWeightGradesMap types.ClassAndWeightGradesMap) types.ClassAndGradeMap {
	classAndGradeMap := make(types.ClassAndGradeMap)
	for className, weightGradeMap := range classAndWeightGradesMap {
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

		classAndGradeMap[className] = strconv.FormatFloat(totalPercentage, 'f', 1, 64)
	}

	return classAndGradeMap
}
