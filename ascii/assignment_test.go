package ascii

import (
	"io"
	"os"
	"strings"
	"testing"

	"github.com/Chance093/gradr/constants"
	"github.com/Chance093/gradr/types"
)

func init() {
	// Sub constants for predictable layout
	constants.WHITE_SPACE = " "
	constants.VERTICAL_BORDER_CHAR = "|"
	constants.HORIZONTAL_BORDER_CHAR = "-"
	constants.INIT_MAX_ASSIGNMENT_COLUMN_LENGTH = 5
	constants.INIT_MAX_GRADE_COLUMN_LENGTH = 5
	constants.INIT_MAX_TYPE_COLUMN_LENGTH = 5
}

func TestGetAssignmentColumnLengths(t *testing.T) {
	assignments := types.Assignments{
		{Name: "Homework1", Grade: "95.5", Type: "Quiz"},
		{Name: "H1", Grade: "83.0", Type: "Homework"},
		{Name: "FinalProject", Grade: "100.0", Type: "Project"},
	}

	gotA, gotG, gotT := getAssignmentColumnLengths(assignments)

	if gotA != len("FinalProject") {
		t.Errorf("expected maxA=%d, got %d", len("FinalProject"), gotA)
	}
	if gotG != len("100.0")+1 {
		t.Errorf("expected maxG=%d, got %d", len("100.0")+1, gotG)
	}
	if gotT != len("Homework") {
		t.Errorf("expected maxT=%d, got %d", len("Homework"), gotT)
	}
}

func TestGetAssignmentColumnLengths_Empty(t *testing.T) {
	var assignments types.Assignments
	gotA, gotG, gotT := getAssignmentColumnLengths(assignments)

	if gotA != constants.INIT_MAX_ASSIGNMENT_COLUMN_LENGTH ||
		gotG != constants.INIT_MAX_GRADE_COLUMN_LENGTH ||
		gotT != constants.INIT_MAX_TYPE_COLUMN_LENGTH {
		t.Errorf(
			"expected initial constants (%d,%d,%d), got (%d,%d,%d)",
			constants.INIT_MAX_ASSIGNMENT_COLUMN_LENGTH,
			constants.INIT_MAX_GRADE_COLUMN_LENGTH,
			constants.INIT_MAX_TYPE_COLUMN_LENGTH,
			gotA, gotG, gotT,
		)
	}
}

func TestGetAssignmentLines_WithAssignments(t *testing.T) {
	assignments := types.Assignments{
		{Name: "HW1", Grade: "92.0", Type: "Quiz"},
		{Name: "Essay", Grade: "N/A", Type: "Exam"},
	}

	lines := getAssignmentLines(6, 5, 6, assignments)

	if len(lines) != 2 {
		t.Fatalf("expected 2 lines, got %d", len(lines))
	}

	for _, line := range lines {
		if !strings.Contains(line, constants.VERTICAL_BORDER_CHAR) {
			t.Errorf("missing border character in line: %q", line)
		}
		if strings.Contains(line, "92.0") && !strings.Contains(line, "%") {
			t.Errorf("expected percent symbol for numeric grade: %q", line)
		}
		if strings.Contains(line, "N/A") && strings.Contains(line, "%") {
			t.Errorf("should not append percent for 'N/A': %q", line)
		}
	}
}

func TestGetAssignmentLines_Empty(t *testing.T) {
	lines := getAssignmentLines(5, 5, 5, types.Assignments{})

	if len(lines) != 1 {
		t.Fatalf("expected 1 fallback line, got %d", len(lines))
	}
	got := lines[0]
	if !strings.Contains(got, "No Assignments") {
		t.Errorf("expected 'No Assignments' fallback, got %q", got)
	}
	if !strings.Contains(got, "N/A") {
		t.Errorf("expected 'N/A' placeholders, got %q", got)
	}
	if !strings.Contains(got, "|") {
		t.Errorf("expected border characters, got %q", got)
	}
}

func TestDisplayAssignmentGrades(t *testing.T) {
	assignments := types.Assignments{
		{Name: "Lab1", Grade: "90.0", Type: "Lab"},
		{Name: "Midterm", Grade: "85.5", Type: "Exam"},
	}

	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	DisplayAssignmentGrades(assignments)

	w.Close()
	out, _ := io.ReadAll(r)
	os.Stdout = old

	output := string(out)

	if !strings.Contains(output, "Assignment") || !strings.Contains(output, "Grade") || !strings.Contains(output, "Type") {
		t.Errorf("header cells missing: %s", output)
	}
	if !strings.Contains(output, "Lab1") || !strings.Contains(output, "Midterm") {
		t.Errorf("expected assignment rows not found: %s", output)
	}
	if !strings.Contains(output, constants.HORIZONTAL_BORDER_CHAR) {
		t.Errorf("expected border lines, got none: %s", output)
	}
	if !strings.Contains(output, "%") {
		t.Errorf("expected percent signs after grades, got none: %s", output)
	}
}
