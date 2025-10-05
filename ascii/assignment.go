package ascii

import (
	"fmt"
	"strings"

	"github.com/Chance093/grader/constants"
	"github.com/Chance093/grader/types"
)

func DisplayAssignmentGrades(assignments types.Assignments) {
	maxA, maxG, maxT := getAssignmentColumnLengths(assignments)
	topLine := getHorizontalBorderLine([]int{maxA, maxG, maxT}, "top")
	bottomLine := getHorizontalBorderLine([]int{maxA, maxG, maxT}, "bottom")
	headerLine := getHeaderLine([]int{maxA, maxG, maxT}, []string{"Assignment", "Grade", "Type"})
	borderLine := getHeaderBorderLine([]int{maxA, maxG, maxT})
	assignmentLines := getAssignmentLines(maxA, maxG, maxT, assignments)

	fmt.Println("")
	fmt.Println(topLine)
	fmt.Println(headerLine)
	fmt.Println(borderLine)

	for _, line := range assignmentLines {
		fmt.Println(line)
	}

	fmt.Println(bottomLine)
	fmt.Println("")
}

func getAssignmentColumnLengths(assignments types.Assignments) (int, int, int) {
	maxA := constants.INIT_MAX_ASSIGNMENT_COLUMN_LENGTH
	maxG := constants.INIT_MAX_GRADE_COLUMN_LENGTH
	maxT := constants.INIT_MAX_TYPE_COLUMN_LENGTH

	for _, assignment := range assignments {
		if len(assignment.Name) > maxA {
			maxA = len(assignment.Name)
		}

		if len(assignment.Grade)+1 > maxG {
			maxG = len(assignment.Grade) + 1
		}

		if len(assignment.Type) > maxT {
			maxT = len(assignment.Type)
		}
	}

	return maxA, maxG, maxT
}

func getAssignmentLines(maxA, maxG, maxT int, assignments types.Assignments) []string {
	var lines []string

	// handle case of no assignments
	if len(assignments) <= 0 {
		assignments = types.Assignments{
			{
				Name:  "No Assignments",
				Grade: " N/A",
				Type:  " N/A",
			},
		}
	}

	for _, assignment := range assignments {
		assignmentLine := constants.WHITE_SPACE + constants.VERTICAL_BORDER_CHAR + constants.WHITE_SPACE + assignment.Name
		for i := 0; i < maxA-len(assignment.Name); i++ {
			assignmentLine += constants.WHITE_SPACE
		}
		assignmentLine += constants.WHITE_SPACE + constants.VERTICAL_BORDER_CHAR + constants.WHITE_SPACE + assignment.Grade

		// add percent if there is a grade
		if strings.Contains(assignment.Grade, "N/A") {
			assignmentLine += constants.WHITE_SPACE
		} else {
			assignmentLine += "%"
		}

		for i := 0; i < maxG-len(assignment.Grade)-1; i++ {
			assignmentLine += constants.WHITE_SPACE
		}
		assignmentLine += constants.WHITE_SPACE + constants.VERTICAL_BORDER_CHAR + constants.WHITE_SPACE + assignment.Type

		for i := 0; i < maxT-len(assignment.Type); i++ {
			assignmentLine += constants.WHITE_SPACE
		}

		assignmentLine += constants.WHITE_SPACE + constants.VERTICAL_BORDER_CHAR

		lines = append(lines, assignmentLine)
	}

	return lines
}
