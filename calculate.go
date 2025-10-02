package main

import (
	"strconv"

	"github.com/Chance093/gradr/types"
)


func calculateGrades(raw ClassesAndGradesRaw) types.ClassAndGradeMap {
	classAndWeightGradesMap := getClassAndWeightGradesMap(raw)
	classAndGradeMap := getClassAndGradeMap(classAndWeightGradesMap)

	return classAndGradeMap
}

func getClassAndWeightGradesMap(raw ClassesAndGradesRaw) types.ClassAndWeightGradesMap {
	classesAndGrades := make(types.ClassAndWeightGradesMap)
	for _, row := range raw {
		weightGradesMap, classExists := classesAndGrades[row.className]
		if classExists {
			grades, weightExists := weightGradesMap[row.weight]
			if weightExists {
				// append grade to grades slice associated with weight
				classesAndGrades[row.className][row.weight] = append(grades, row.grade)
			} else {
				// init weight and associate grade
				classesAndGrades[row.className][row.weight] = []float64{row.grade}
			}
		} else { // else initialize className to classesAndGrades map
			classesAndGrades[row.className] = map[int][]float64{
				row.weight: {row.grade},
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
