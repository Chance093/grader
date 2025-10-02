package ascii

import (
	"fmt"

	"github.com/Chance093/gradr/types"
)

func DisplayAssignmentGrades(data []types.AssignmentsRaw) {
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

func getMaxCharAssLengths(data []types.AssignmentsRaw) (int, int, int) {
	maxAssignmentCharLength := 20
	maxGradeCharLength := 5
	maxAssignmentTypeCharLength := 5

	for _, dat := range data {
		if len(dat.Assignment) > maxAssignmentCharLength {
			maxAssignmentCharLength = len(dat.Assignment)
		}

		if len(dat.Grade)+1 > maxGradeCharLength {
			maxGradeCharLength = len(dat.Grade) + 1
		}

		if len(dat.AssignmentType) > maxAssignmentTypeCharLength {
			maxAssignmentTypeCharLength = len(dat.AssignmentType)
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

func getAssignmentLines(maxA, maxG, maxT int, data []types.AssignmentsRaw) []string {
	var lines []string

	if len(data) <= 0 {
		data = []types.AssignmentsRaw{
			{
				Assignment:     "No Assignments",
				Grade:          " N/A",
				AssignmentType: " N/A",
			},
		}
	}

	for _, dat := range data {
		assignmentLine := fmt.Sprintf(" | %s", dat.Assignment)
		for i := 0; i < maxA-len(dat.Assignment); i++ {
			assignmentLine += " "
		}

		assignmentLine += fmt.Sprintf(" | %s", dat.Grade)
		if dat.Assignment == "No Assignments" {
			assignmentLine += " "
		} else {
			assignmentLine += "%"
		}

		for i := 0; i < maxG-len(dat.Grade)-1; i++ {
			assignmentLine += " "
		}

		assignmentLine += fmt.Sprintf(" | %s", dat.AssignmentType)
		for i := 0; i < maxT-len(dat.AssignmentType); i++ {
			assignmentLine += " "
		}

		assignmentLine += " |"

		lines = append(lines, assignmentLine)
	}

	return lines
}
