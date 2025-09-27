package main

import "fmt"

var data = []struct {
	className string
	grade     string
}{
	{
		className: "Beginning and Intermediate Algebra",
		grade:     "97.3",
	},
	{
		className: "Geometry",
		grade:     "98.52",
	},
	{
		className: "Calculus",
		grade:     "92.5",
	},
	{
		className: "Linear Algrebra",
		grade:     "90.88",
	},
}

func getMaxCharLengths() (int, int) {
	maxClassCharLength := 5
	maxGradeCharLength := 5

	for _, dat := range data {
		if len(dat.className) > maxClassCharLength {
			maxClassCharLength = len(dat.className)
		}

		if len(dat.grade)+1 > maxGradeCharLength {
			maxGradeCharLength = len(dat.grade) + 1
		}
	}

	return maxClassCharLength, maxGradeCharLength
}

func getClassGradesAscii() {
	maxC, maxG := getMaxCharLengths()
	topAndBottomLine := getTopAndBottomLine(maxC, maxG)
	headerLine := getHeaderLine(maxC, maxG)
	borderLine := getBorderLine(maxC, maxG)
	assignmentLines := getAssignmentLines(maxC, maxG)

	fmt.Println("")
	fmt.Println(topAndBottomLine)
	fmt.Println(headerLine)
	fmt.Println(borderLine)

	for _, line := range assignmentLines {
		fmt.Println(line)
	}

	fmt.Println(topAndBottomLine)
	fmt.Println("")
}

func getTopAndBottomLine(maxC, maxG int) string {
	total := 2 + maxC + 3 + maxG + 2

	topAndBottomLine := " "
	for range total {
		topAndBottomLine += "="
	}

	return topAndBottomLine
}

func getHeaderLine(maxC, maxG int) string {
	headerLine := " | Class"
	for i := 0; i < maxC-5; i++ {
		headerLine += " "
	}
	headerLine += " | Grade"

	for i := 0; i < maxG-5; i++ {
		headerLine += " "
	}
	headerLine += " |"

	return headerLine
}

func getBorderLine(maxC, maxG int) string {
	borderLine := " |"
	for i := 0; i < maxC+2; i++ {
		borderLine += "-"
	}

	borderLine += "|"
	for i := 0; i < maxG+2; i++ {
		borderLine += "-"
	}
	borderLine += "|"

	return borderLine
}

func getAssignmentLines(maxC, maxG int) []string {
	var lines []string

	for _, dat := range data {
		assignmentLine := fmt.Sprintf(" | %s", dat.className)
		for i := 0; i < maxC-len(dat.className); i++ {
			assignmentLine += " "
		}
		assignmentLine += fmt.Sprintf(" | %s", dat.grade) + "%"

		for i := 0; i < maxG-len(dat.grade); i++ {
			assignmentLine += " "
		}
		assignmentLine += "|"

		lines = append(lines, assignmentLine)
	}

	return lines
}
