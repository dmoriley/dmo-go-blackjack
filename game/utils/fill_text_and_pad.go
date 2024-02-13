package utils

import (
	"bytes"
	"fmt"
)

// Fill supplied buffer with text
func FillTextAndPad(
	buffer *bytes.Buffer,
	length int,
	fillerChar byte,
	borderChar byte,
	insertText string,
	alignment string,
) {
	if len(insertText) != 0 {
		insertText = fmt.Sprintf(" %s ", insertText)
	}
	fillerLength := length - len(insertText) - 2 // minus 2 for the border chars
	strFillerChar := string(fillerChar)
	strBorderChar := string(borderChar)

	// start border char
	buffer.WriteString(strBorderChar)
	switch alignment {
	case "left":
		buffer.WriteString(insertText)
		for i := 0; i < fillerLength; i++ {
			buffer.WriteString(strFillerChar)
		}

	case "middle", "": // default to middle when alignment is blank
		halfLength := fillerLength / 2

		for i := 0; i < halfLength; i++ {
			buffer.WriteString(strFillerChar)
		}
		buffer.WriteString(insertText)

		if halfLength*2 != fillerLength {
			// need to increment for the 1 lost on odd number division
			halfLength++
		}

		for i := 0; i < halfLength; i++ {
			buffer.WriteString(strFillerChar)
		}
	case "right":
		for i := 0; i < fillerLength; i++ {
			buffer.WriteString(strFillerChar)
		}
		buffer.WriteString(insertText)
	}

	// end border char added
	buffer.WriteString(strBorderChar)
	buffer.WriteString("\n")
}

// Print text to fmt buffer
func PrintTextAndPad(
	length int,
	fillerChar byte,
	borderChar byte,
	insertText string,
	alignment string,
	finishWithNewLine bool,
) {
	if len(insertText) != 0 {
		insertText = fmt.Sprintf(" %s ", insertText)
	}
	fillerLength := length - len(insertText) - 2 // minus 2 for the border chars
	strFillerChar := string(fillerChar)
	strBorderChar := string(borderChar)

	// start border char
	fmt.Print(strBorderChar)
	switch alignment {
	case "left":
		fmt.Print(insertText)
		for i := 0; i < fillerLength; i++ {
			fmt.Print(strFillerChar)
		}

	case "middle", "": // default to middle when alignment is blank
		halfLength := fillerLength / 2

		for i := 0; i < halfLength; i++ {
			fmt.Print(strFillerChar)
		}
		fmt.Print(insertText)

		if halfLength*2 != fillerLength {
			// need to increment for the 1 lost on odd number division
			halfLength++
		}

		for i := 0; i < halfLength; i++ {
			fmt.Print(strFillerChar)
		}
	case "right":
		for i := 0; i < fillerLength; i++ {
			fmt.Print(strFillerChar)
		}
		fmt.Print(insertText)
	}

	// end border char added
	fmt.Print(strBorderChar)
	if finishWithNewLine {
		fmt.Print("\n")
	}
}
