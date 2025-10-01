package main

import "fmt"

func getMaxCharLengths(data map[string]string) (int, int) {
	maxClassCharLength := 20
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
	classLines := getClassLines(maxC, maxG, data)

	fmt.Println("")
	fmt.Println(topAndBottomLine)
	fmt.Println(headerLine)
	fmt.Println(borderLine)

	for _, line := range classLines {
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

func getClassLines(maxC, maxG int, data map[string]string) []string {
	var lines []string

	if len(data) <= 0 {
		data = map[string]string{
			"No Classes": " N/A",
		}
	}

	for className, grade := range data {
		classLine := fmt.Sprintf(" | %s", className)
		for i := 0; i < maxC-len(className); i++ {
			classLine += " "
		}
		classLine += fmt.Sprintf(" | %s", grade)

		if grade == " N/A" {
			classLine += " "
		} else {
			classLine += "%"
		}

		for i := 0; i < maxG-len(grade); i++ {
			classLine += " "
		}
		classLine += "|"

		lines = append(lines, classLine)
	}

	return lines
}

func getAssignmentGradesAscii(data []AssignmentsRaw) {
	maxA, maxG, maxT := getMaxCharAssLengths(data)
	topAndBottomLine := getTopAndBottomAssLine(maxA, maxG, maxT)
	headerLine := getHeaderAssLine(maxA, maxG, maxT)
	borderLine := getBorderAssLine(maxA, maxG, maxT)
	assignmentLines := getAssignmentLines(maxA, maxG, maxT, data)

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

func getMaxCharAssLengths(data []AssignmentsRaw) (int, int, int) {
	maxAssignmentCharLength := 20
	maxGradeCharLength := 5
	maxAssignmentTypeCharLength := 5

	for _, dat := range data {
		if len(dat.assignment) > maxAssignmentCharLength {
			maxAssignmentCharLength = len(dat.assignment)
		}

		if len(dat.grade)+1 > maxGradeCharLength {
			maxGradeCharLength = len(dat.grade) + 1
		}

		if len(dat.assignmentType) > maxAssignmentTypeCharLength {
			maxAssignmentTypeCharLength = len(dat.assignmentType)
		}
	}

	return maxAssignmentCharLength, maxGradeCharLength, maxAssignmentTypeCharLength
}

func getTopAndBottomAssLine(maxA, maxG, maxT int) string {
	total := 2 + maxA + 3 + maxG + 3 + maxT + 2

	topAndBottomLine := " "
	for range total {
		topAndBottomLine += "="
	}

	return topAndBottomLine
}

func getHeaderAssLine(maxA, maxG, maxT int) string {
	headerLine := " | Assignment"
	for i := 0; i < maxA-10; i++ {
		headerLine += " "
	}

	headerLine += " | Grade"
	for i := 0; i < maxG-5; i++ {
		headerLine += " "
	}

	headerLine += " | Type"
	for i := 0; i < maxT-4; i++ {
		headerLine += " "
	}
	headerLine += " |"

	return headerLine
}

func getBorderAssLine(maxA, maxG, maxT int) string {
	borderLine := " |"
	for i := 0; i < maxA+2; i++ {
		borderLine += "-"
	}

	borderLine += "|"
	for i := 0; i < maxG+2; i++ {
		borderLine += "-"
	}

	borderLine += "|"
	for i := 0; i < maxT+2; i++ {
		borderLine += "-"
	}
	borderLine += "|"

	return borderLine
}

func getAssignmentLines(maxA, maxG, maxT int, data []AssignmentsRaw) []string {
	var lines []string

	if len(data) <= 0 {
		data = []AssignmentsRaw{
			{
				assignment:     "No Assignments",
				grade:          " N/A",
				assignmentType: " N/A",
			},
		}
	}

	for _, dat := range data {
		assignmentLine := fmt.Sprintf(" | %s", dat.assignment)
		for i := 0; i < maxA-len(dat.assignment); i++ {
			assignmentLine += " "
		}

		assignmentLine += fmt.Sprintf(" | %s", dat.grade)
		if dat.assignment == "No Assignments" {
			assignmentLine += " "
		} else {
			assignmentLine += "%"
		}

		for i := 0; i < maxG-len(dat.grade)-1; i++ {
			assignmentLine += " "
		}

		assignmentLine += fmt.Sprintf(" | %s", dat.assignmentType)
		for i := 0; i < maxT-len(dat.assignmentType); i++ {
			assignmentLine += " "
		}

		assignmentLine += " |"

		lines = append(lines, assignmentLine)
	}

	return lines
}
