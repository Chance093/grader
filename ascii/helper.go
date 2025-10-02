package ascii

func getHorizontalBorderLine(ints []int) string {
	totalChars := 4                   // outside chars plus gap "|  |"
	totalChars += (len(ints) - 1) * 3 // inside chars plus gap " | "
	for _, int := range ints {        // add all max column lengths
		totalChars += int
	}

	horizontalBorder := WHITE_SPACE
	for range totalChars {
		horizontalBorder += HORIZONTAL_BORDER_CHAR
	}

	return horizontalBorder
}

func getHeaderLine(ints []int, headers []string) string {
	var headerLine string
	for i, int := range ints {
		header := headers[i]
		headerLine += WHITE_SPACE + VERTICAL_BORDER_CHAR + WHITE_SPACE + header
		for i := 0; i < int-len(header); i++ {
			headerLine += WHITE_SPACE
		}
	}

	headerLine += WHITE_SPACE + VERTICAL_BORDER_CHAR

	return headerLine
}

func getHeaderBorderLine(ints []int) string {
	headerBorderLine := WHITE_SPACE
	for _, int := range ints {
		headerBorderLine += VERTICAL_BORDER_CHAR
		for i := 0; i < int+2; i++ {
			headerBorderLine += HEADER_BORDER_CHAR
		}
	}

	headerBorderLine += VERTICAL_BORDER_CHAR

	return headerBorderLine
}
