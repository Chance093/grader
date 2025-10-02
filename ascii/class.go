package ascii

import (
	"fmt"

	"github.com/Chance093/gradr/types"
)

/*
e.g.

	===============================
	| Class               | Grade |
	|---------------------|-------|
	| Algrebra            | 90.0% |
	| Calculus            | 90.0% |
	===============================
*/
func DisplayClassGrades(data map[string]string) {
	maxC, maxG := getMaxCharLengths(data)
	topAndBottomLine := getHorizontalBorder(maxC, maxG)
	headerLine := getHeaderLine(maxC, maxG)
	headerBorderLine := getHeaderBorderLine(maxC, maxG)
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

func getMaxCharLengths(data types.ClassAndGradeMap) (int, int) {
	maxC := INIT_MAX_CLASS_COLUMN_LENGTH
	maxG := INIT_MAX_GRADE_COLUMN_LENGTH

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

// e.g. ======================
func getHorizontalBorder(maxC, maxG int) string {
	totalChars := 2 + maxC + 3 + maxG + 2

	horizontalBorder := WHITE_SPACE
	for range totalChars {
		horizontalBorder += HORIZONTAL_BORDER_CHAR
	}

	return horizontalBorder
}

// e.g. | Class             | Grade |
func getHeaderLine(maxC, maxG int) string {
	headerLine := WHITE_SPACE + VERTICAL_BORDER_CHAR + WHITE_SPACE + "Class"
	for i := 0; i < maxC-5; i++ {
		headerLine += WHITE_SPACE
	}
	headerLine += WHITE_SPACE + VERTICAL_BORDER_CHAR + WHITE_SPACE + "Grade"

	for i := 0; i < maxG-5; i++ {
		headerLine += WHITE_SPACE
	}
	headerLine += WHITE_SPACE + VERTICAL_BORDER_CHAR

	return headerLine
}

// e.g. |----------------|------|
func getHeaderBorderLine(maxC, maxG int) string {
	headerBorderLine := WHITE_SPACE + VERTICAL_BORDER_CHAR
	for i := 0; i < maxC+2; i++ {
		headerBorderLine += HEADER_BORDER_CHAR
	}

	headerBorderLine += VERTICAL_BORDER_CHAR
	for i := 0; i < maxG+2; i++ {
		headerBorderLine += HEADER_BORDER_CHAR
	}
	headerBorderLine += VERTICAL_BORDER_CHAR

	return headerBorderLine
}

/*
e.g.
 | Algebra          | 90.0% |
 | Calculus         | 93.8% |
*/
func getClassLines(maxC, maxG int, data types.ClassAndGradeMap) []string {
	var lines []string

	if len(data) <= 0 {
		data = map[string]string{
			"No Classes": " N/A",
		}
	}

	for className, grade := range data {
		classLine := WHITE_SPACE + VERTICAL_BORDER_CHAR + WHITE_SPACE + className
		for i := 0; i < maxC-len(className); i++ {
			classLine += WHITE_SPACE
		}
		classLine += WHITE_SPACE + VERTICAL_BORDER_CHAR + WHITE_SPACE + grade

		if grade == " N/A" {
			classLine += WHITE_SPACE
		} else {
			classLine += "%"
		}

		for i := 0; i < maxG-len(grade); i++ {
			classLine += WHITE_SPACE
		}
		classLine += VERTICAL_BORDER_CHAR

		lines = append(lines, classLine)
	}

	return lines
}
