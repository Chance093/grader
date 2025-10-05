package ascii

import (
	"fmt"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/Chance093/gradr/constants"
	"github.com/Chance093/gradr/types"
)

func init() {
	// Sub constants for predictable tests
	constants.WHITE_SPACE = " "
	constants.VERTICAL_BORDER_CHAR = "|"
	constants.HORIZONTAL_BORDER_CHAR = "-"
	constants.INIT_MAX_CLASS_COLUMN_LENGTH = 5
	constants.INIT_MAX_GRADE_COLUMN_LENGTH = 5
}

func TestGetClassColumnLengths(t *testing.T) {
	tests := []struct {
		name  string
		input types.ClassAndGradeMap
		wantC int
		wantG int
	}{
		{
			name: "normal class data widths",
			input: types.ClassAndGradeMap{
				"Math":    "90.0",
				"History": "82.5",
			},
			wantC: 7, // length(History)
			wantG: 5, // len("90.0")+1 = 5
		},
		{
			name:  "empty map uses initial constants",
			input: types.ClassAndGradeMap{},
			wantC: constants.INIT_MAX_CLASS_COLUMN_LENGTH,
			wantG: constants.INIT_MAX_GRADE_COLUMN_LENGTH,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotC, gotG := getClassColumnLengths(tt.input)
			if gotC != tt.wantC || gotG != tt.wantG {
				t.Errorf("getClassColumnLengths() = (%d,%d), want (%d,%d)",
					gotC, gotG, tt.wantC, tt.wantG)
			}
		})
	}
}

func TestGetClassLines_WithData(t *testing.T) {
	data := types.ClassAndGradeMap{
		"Math":    "88.0",
		"Science": "93.5",
	}

	lines := getClassLines(7, 6, data)
	if len(lines) != 2 {
		t.Fatalf("expected 2 lines, got %d", len(lines))
	}
	for _, line := range lines {
		if !strings.Contains(line, constants.VERTICAL_BORDER_CHAR) {
			t.Errorf("line missing border char: %q", line)
		}
		if !strings.Contains(line, "%") {
			t.Errorf("line missing percent char: %q", line)
		}
	}
}

func TestGetClassLines_EmptyData(t *testing.T) {
	data := types.ClassAndGradeMap{}

	lines := getClassLines(5, 5, data)
	if len(lines) != 1 {
		t.Fatalf("expected single 'No Classes' line, got %d", len(lines))
	}
	got := lines[0]
	if !strings.Contains(got, "No Classes") {
		t.Errorf("expected 'No Classes' fallback, got %q", got)
	}
	if !strings.Contains(got, "N/A") {
		t.Errorf("expected 'N/A' grade fallback, got %q", got)
	}
	if !strings.Contains(got, "|") {
		t.Errorf("expected border characters, got %q", got)
	}
}

func TestDisplayClassGrades(t *testing.T) {
	data := map[string]string{
		"Math":    "95.5",
		"English": "89.0",
	}

	// Capture stdout
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	DisplayClassGrades(data)

	w.Close()
	out, _ := io.ReadAll(r)
	os.Stdout = old

	output := string(out)
	lines := strings.Split(strings.TrimSpace(output), "\n")
	if len(lines) < 5 {
		t.Errorf("expected at least 5 lines (borders, header, and classes), got %d", len(lines))
	}

	if !strings.Contains(output, "Class") || !strings.Contains(output, "Grade") {
		t.Errorf("missing headers in output: %q", output)
	}
	if !strings.Contains(output, "Math") {
		t.Errorf("missing class row 'Math' in output: %q", output)
	}
	if !strings.Contains(output, constants.HORIZONTAL_BORDER_CHAR) {
		t.Errorf("expected horizontal border chars, got %q", output)
	}

	// should print a top border and bottom border line twice
	borderCount := strings.Count(output, constants.HORIZONTAL_BORDER_CHAR)
	if borderCount == 0 {
		t.Errorf("expected border lines, got none")
	}

	fmt.Printf("Captured DisplayClassGrades output:\n%s\n", output)
}
