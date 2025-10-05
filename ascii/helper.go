package ascii

import "github.com/Chance093/grader/constants"

func getHorizontalBorderLine(ints []int, line string) string {
	var horizontalBorder string
	if line == "top" {
		horizontalBorder += constants.WHITE_SPACE + constants.TOP_LEFT_CHAR
	} else if line == "bottom" {
		horizontalBorder += constants.WHITE_SPACE + constants.BOTTOM_LEFT_CHAR
	}

	for i, int := range ints {
		for range int + 2 {
			horizontalBorder += constants.HORIZONTAL_BORDER_CHAR
		}

		if i != len(ints)-1 {
			if line == "top" {
				horizontalBorder += constants.HORIZONTAL_TOP_CHAR
			} else if line == "bottom" {
				horizontalBorder += constants.HORIZONTAL_BOTTOM_CHAR
			}
		}
	}

	if line == "top" {
		horizontalBorder += constants.TOP_RIGHT_CHAR
	} else if line == "bottom" {
		horizontalBorder += constants.BOTTOM_RIGHT_CHAR
	}

	return horizontalBorder
}

func getHeaderLine(ints []int, headers []string) string {
	var headerLine string
	for i, int := range ints {
		header := headers[i]
		headerLine += constants.WHITE_SPACE + constants.VERTICAL_BORDER_CHAR + constants.WHITE_SPACE + header
		for i := 0; i < int-len(header); i++ {
			headerLine += constants.WHITE_SPACE
		}
	}

	headerLine += constants.WHITE_SPACE + constants.VERTICAL_BORDER_CHAR

	return headerLine
}

func getHeaderBorderLine(ints []int) string {
	headerBorderLine := constants.WHITE_SPACE + constants.VERTICAL_RIGHT_CHAR
	for i, int := range ints {
		if i != 0 {
			headerBorderLine += constants.CROSS_CHAR
		}
		for i := 0; i < int+2; i++ {
			headerBorderLine += constants.HORIZONTAL_BORDER_CHAR
		}
	}

	headerBorderLine += constants.VERTICAL_LEFT_CHAR

	return headerBorderLine
}
