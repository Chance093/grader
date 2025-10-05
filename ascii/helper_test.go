package ascii

import (
	"strings"
	"testing"

	"github.com/Chance093/gradr/constants"
)

func init() {
	constants.WHITE_SPACE = " "
	constants.HORIZONTAL_BORDER_CHAR = "="
	constants.HEADER_BORDER_CHAR = "-"
	constants.VERTICAL_BORDER_CHAR = "|"
}

func TestGetHorizontalBorderLine(t *testing.T) {
	tests := []struct {
		name      string
		ints      []int
		wantCount int
	}{
		{
			name:      "Single column width 5",
			ints:      []int{5},
			wantCount: 4 + 5, // 4 for "|  |" + width
		},
		{
			name: "Three columns width 3,4,2",
			ints: []int{3, 4, 2},
			// 4 start + (len-1)*3 = 4 + (3-1)*3 = 10 + total widths = 3+4+2=9 => 19
			wantCount: 19,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getHorizontalBorderLine(tt.ints)
			if !strings.HasPrefix(got, constants.WHITE_SPACE) {
				t.Errorf("border should begin with a space, got %q", got)
			}
			borderRunes := strings.TrimPrefix(got, constants.WHITE_SPACE)
			if len(borderRunes) != tt.wantCount {
				t.Errorf("got %d border chars, want %d", len(borderRunes), tt.wantCount)
			}
			for _, ch := range borderRunes {
				if string(ch) != constants.HORIZONTAL_BORDER_CHAR {
					t.Errorf("expected only horizontal border chars, got %q", got)
				}
			}
		})
	}
}

func TestGetHeaderLine(t *testing.T) {
	tests := []struct {
		name     string
		ints     []int
		headers  []string
		expected string
	}{
		{
			name:     "Single column exact fit",
			ints:     []int{4},
			headers:  []string{"Test"},
			expected: " " + "|" + " " + "Test" + " " + "|",
		},
		{
			name:    "Two columns different lengths",
			ints:    []int{5, 3},
			headers: []string{"Num", "ID"},
			// | Num..| ID.| (padding '.' for clarity)
			expected: " | Num  | ID |",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getHeaderLine(tt.ints, tt.headers)

			// Trim any excessive whitespace for matching (since output always begins with a space)
			if !strings.Contains(got, constants.VERTICAL_BORDER_CHAR) {
				t.Errorf("missing vertical borders in output, got %q", got)
			}

			for _, header := range tt.headers {
				if !strings.Contains(got, header) {
					t.Errorf("missing header %q in output %q", header, got)
				}
			}
		})
	}
}

func TestGetHeaderBorderLine(t *testing.T) {
	tests := []struct {
		name        string
		ints        []int
		expectedRun string
	}{
		{
			name:        "Single column width 4",
			ints:        []int{4},
			expectedRun: " " + "|" + strings.Repeat(constants.HEADER_BORDER_CHAR, 6) + "|",
		},
		{
			name:        "Two columns width 2, 3",
			ints:        []int{2, 3},
			expectedRun: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getHeaderBorderLine(tt.ints)
			if !strings.HasPrefix(got, constants.WHITE_SPACE) {
				t.Errorf("header border should start with a space, got %q", got)
			}
			if !strings.HasSuffix(got, constants.VERTICAL_BORDER_CHAR) {
				t.Errorf("header border should end with a vertical border char, got %q", got)
			}
			for _, ch := range strings.TrimSpace(got) {
				if !(string(ch) == constants.VERTICAL_BORDER_CHAR ||
					string(ch) == constants.HEADER_BORDER_CHAR) {
					t.Errorf("unexpected character %q in header border", ch)
				}
			}
		})
	}
}
