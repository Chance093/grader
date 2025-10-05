package ascii

import (
	"fmt"

	"github.com/Chance093/gradr/constants"
	"github.com/Chance093/gradr/types"
)

func DisplayClassGrades(data map[string]string) {
	maxC, maxG := getClassColumnLengths(data)
	topAndBottomLine := getHorizontalBorderLine([]int{maxC, maxG})
	headerLine := getHeaderLine([]int{maxC, maxG}, []string{"Class", "Grade"})
	headerBorderLine := getHeaderBorderLine([]int{maxC, maxG})
	classLines := getClassLines(maxC, maxG, data)

	fmt.Println("")
	fmt.Println(topAndBottomLine)
	fmt.Println(headerLine)
	fmt.Println(headerBorderLine)

	for _, line := range classLines {
		fmt.Println(line)
	}

	fmt.Println(topAndBottomLine)
	fmt.Println("")
}

func getClassColumnLengths(data types.ClassAndGradeMap) (int, int) {
	maxC := constants.INIT_MAX_CLASS_COLUMN_LENGTH
	maxG := constants.INIT_MAX_GRADE_COLUMN_LENGTH

	for className, grade := range data {
		if len(className) > maxC {
			maxC = len(className)
		}

		if len(grade)+1 > maxG {
			maxG = len(grade) + 1
		}
	}

	return maxC, maxG
}

func getClassLines(maxC, maxG int, data types.ClassAndGradeMap) []string {
	var lines []string

	// handle case of no classes
	if len(data) <= 0 {
		data = map[string]string{
			"No Classes": " N/A",
		}
	}

	for className, grade := range data {
		classLine := constants.WHITE_SPACE + constants.VERTICAL_BORDER_CHAR + constants.WHITE_SPACE + className
		for i := 0; i < maxC-len(className); i++ {
			classLine += constants.WHITE_SPACE
		}
		classLine += constants.WHITE_SPACE + constants.VERTICAL_BORDER_CHAR + constants.WHITE_SPACE + grade

		// add percent if there is a grade
		if grade == " N/A" {
			classLine += constants.WHITE_SPACE
		} else {
			classLine += "%"
		}

		for i := 0; i < maxG-len(grade); i++ {
			classLine += constants.WHITE_SPACE
		}
		classLine += constants.VERTICAL_BORDER_CHAR

		lines = append(lines, classLine)
	}

	return lines
}
