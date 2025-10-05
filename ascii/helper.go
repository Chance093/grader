package ascii

import "github.com/Chance093/gradr/constants"

func getHorizontalBorderLine(ints []int) string {
	totalChars := 4                   // outside chars plus gap "|  |"
	totalChars += (len(ints) - 1) * 3 // inside chars plus gap " | "
	for _, int := range ints {        // add all max column lengths
		totalChars += int
	}

	horizontalBorder := constants.WHITE_SPACE
	for range totalChars {
		horizontalBorder += constants.HORIZONTAL_BORDER_CHAR
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
	headerBorderLine := constants.WHITE_SPACE
	for _, int := range ints {
		headerBorderLine += constants.VERTICAL_BORDER_CHAR
		for i := 0; i < int+2; i++ {
			headerBorderLine += constants.HEADER_BORDER_CHAR
		}
	}

	headerBorderLine += constants.VERTICAL_BORDER_CHAR

	return headerBorderLine
}
