package main

import "fmt"

func getMaxCharLengths(data map[string]string) (int, int) {
	maxClassCharLength := 5
	maxGradeCharLength := 5

	for className, grade := range data {
		if len(className) > maxClassCharLength {
			maxClassCharLength = len(className)
		}

		if len(grade)+1 > maxGradeCharLength {
			maxGradeCharLength = len(grade) + 1
		}
	}

	return maxClassCharLength, maxGradeCharLength
}

func getClassGradesAscii(data map[string]string) {
	maxC, maxG := getMaxCharLengths(data)
	topAndBottomLine := getTopAndBottomLine(maxC, maxG)
	headerLine := getHeaderLine(maxC, maxG)
	borderLine := getBorderLine(maxC, maxG)
	assignmentLines := getAssignmentLines(maxC, maxG, data)

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

func getAssignmentLines(maxC, maxG int, data map[string]string) []string {
	var lines []string

	for className, grade := range data {
		assignmentLine := fmt.Sprintf(" | %s", className)
		for i := 0; i < maxC-len(className); i++ {
			assignmentLine += " "
		}
		assignmentLine += fmt.Sprintf(" | %s", grade) + "%"

		for i := 0; i < maxG-len(grade); i++ {
			assignmentLine += " "
		}
		assignmentLine += "|"

		lines = append(lines, assignmentLine)
	}

	return lines
}
