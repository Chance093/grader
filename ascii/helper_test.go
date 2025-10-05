package ascii

import (
	"strings"
	"testing"

	"github.com/Chance093/gradr/constants"
)

func init() {
	// Sub the constants for predictable test results
	constants.WHITE_SPACE = " "
	constants.HORIZONTAL_BORDER_CHAR = "-"
	constants.TOP_LEFT_CHAR = "+"
	constants.TOP_RIGHT_CHAR = "+"
	constants.BOTTOM_LEFT_CHAR = "+"
	constants.BOTTOM_RIGHT_CHAR = "+"
	constants.HORIZONTAL_TOP_CHAR = "^"
	constants.HORIZONTAL_BOTTOM_CHAR = "v"
	constants.VERTICAL_BORDER_CHAR = "|"
	constants.VERTICAL_LEFT_CHAR = "<"
	constants.VERTICAL_RIGHT_CHAR = ">"
	constants.CROSS_CHAR = "*"
}

func TestGetHorizontalBorderLine_Top(t *testing.T) {
	widths := []int{3, 5}
	got := getHorizontalBorderLine(widths, "top")

	// Expected pattern: " +---^-----+"
	wantPrefix := " " + constants.TOP_LEFT_CHAR
	wantSuffix := constants.TOP_RIGHT_CHAR

	if !strings.HasPrefix(got, wantPrefix) {
		t.Errorf("top line should start with %q, got %q", wantPrefix, got)
	}
	if !strings.HasSuffix(got, wantSuffix) {
		t.Errorf("top line should end with %q, got %q", wantSuffix, got)
	}
	if !strings.Contains(got, constants.HORIZONTAL_TOP_CHAR) {
		t.Errorf("top line missing HORIZONTAL_TOP_CHAR (%q): %q", constants.HORIZONTAL_TOP_CHAR, got)
	}
	if strings.Contains(got, constants.HORIZONTAL_BOTTOM_CHAR) {
		t.Errorf("top line should not contain bottom connector (%q): %q", constants.HORIZONTAL_BOTTOM_CHAR, got)
	}
	if len(got) == 0 {
		t.Fatal("top line should not be empty")
	}
}

func TestGetHorizontalBorderLine_Bottom(t *testing.T) {
	widths := []int{2, 3}
	got := getHorizontalBorderLine(widths, "bottom")

	// Expected: " +--v---+"
	if !strings.HasPrefix(got, " "+constants.BOTTOM_LEFT_CHAR) {
		t.Errorf("bottom line should start with %q, got %q", constants.BOTTOM_LEFT_CHAR, got)
	}
	if !strings.HasSuffix(got, constants.BOTTOM_RIGHT_CHAR) {
		t.Errorf("bottom line should end with %q, got %q", constants.BOTTOM_RIGHT_CHAR, got)
	}
	if !strings.Contains(got, constants.HORIZONTAL_BOTTOM_CHAR) {
		t.Errorf("bottom line missing bottom connector (%q): %q", constants.HORIZONTAL_BOTTOM_CHAR, got)
	}
	if strings.Contains(got, constants.HORIZONTAL_TOP_CHAR) {
		t.Errorf("bottom line should not contain top connector (%q): %q", constants.HORIZONTAL_TOP_CHAR, got)
	}
}

func TestGetHeaderLine(t *testing.T) {
	ints := []int{4, 3}
	headers := []string{"Name", "Age"}
	got := getHeaderLine(ints, headers)

	for _, h := range headers {
		if !strings.Contains(got, h) {
			t.Errorf("expected header %q to appear in line: %q", h, got)
		}
	}
	if strings.Count(got, constants.VERTICAL_BORDER_CHAR) < 2 {
		t.Errorf("expected at least two vertical borders, got %q", got)
	}
}

func TestGetHeaderBorderLine(t *testing.T) {
	ints := []int{3, 5}
	got := getHeaderBorderLine(ints)

	// It should start with space + right vertical char
	if !strings.HasPrefix(got, " "+constants.VERTICAL_RIGHT_CHAR) {
		t.Errorf("expected prefix with vertical-right %q, got %q", constants.VERTICAL_RIGHT_CHAR, got)
	}

	// It should contain cross-join chars between columns
	if !strings.Contains(got, constants.CROSS_CHAR) {
		t.Errorf("expected to contain cross chars (%q), got %q", constants.CROSS_CHAR, got)
	}

	// It should end with vertical-left char
	if !strings.HasSuffix(got, constants.VERTICAL_LEFT_CHAR) {
		t.Errorf("expected suffix with vertical-left %q, got %q", constants.VERTICAL_LEFT_CHAR, got)
	}

	// Horizontal fill check
	for _, ch := range strings.TrimSpace(got) {
		if !(string(ch) == constants.VERTICAL_RIGHT_CHAR ||
			string(ch) == constants.VERTICAL_LEFT_CHAR ||
			string(ch) == constants.CROSS_CHAR ||
			string(ch) == constants.HORIZONTAL_BORDER_CHAR) {
			t.Errorf("unexpected character %q in header border line: %q", string(ch), got)
			break
		}
	}
}
