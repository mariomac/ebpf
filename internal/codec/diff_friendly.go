package codec

import (
	"errors"
	"fmt"
	"io"
	"strings"
)

func HexString(src io.Reader, dst io.Writer) error {
	buf := make([]byte, 256)
	for {
		n, err := src.Read(buf)
		if err != nil {
			if errors.Is(err, io.EOF) {
				return nil
			}
			return fmt.Errorf("reading input bytes: %w", err)
		}
		for i := 0; i < n; i++ {
			if _, err := fmt.Fprintf(dst, "%02x", buf[i]); err != nil {
				return fmt.Errorf("writing output HEX digit: %w", err)
			}
		}

	}
}

func ConvertToDiffFriendly(previous, current string, lineLen int) string {
	// Remove all line breaks to get raw hexadecimal content
	currentRaw := strings.ReplaceAll(current, "\n", "")

	// Get lines from previous string
	previousLines := strings.Split(previous, "\n")

	// Remove empty lines
	var nonEmptyPreviousLines []string
	for _, line := range previousLines {
		if line != "" {
			nonEmptyPreviousLines = append(nonEmptyPreviousLines, line)
		}
	}

	var resultLines []string
	remainingCurrent := currentRaw

	// Try to match lines from previous in order
	for _, previousLine := range nonEmptyPreviousLines {
		if len(previousLine) <= lineLen && strings.HasPrefix(remainingCurrent, previousLine) {
			// Found a matching line at the beginning
			resultLines = append(resultLines, previousLine)
			remainingCurrent = remainingCurrent[len(previousLine):]
		} else {
			// Try to find the line anywhere in the remaining current content
			index := strings.Index(remainingCurrent, previousLine)
			if index != -1 && len(previousLine) <= lineLen {
				// Add content before the match as separate lines
				if index > 0 {
					beforeMatch := remainingCurrent[:index]
					resultLines = append(resultLines, breakIntoLines(beforeMatch, lineLen)...)
				}
				// Add the matching line
				resultLines = append(resultLines, previousLine)
				remainingCurrent = remainingCurrent[index+len(previousLine):]
			}
		}
	}

	// Add any remaining content as lines
	if remainingCurrent != "" {
		resultLines = append(resultLines, breakIntoLines(remainingCurrent, lineLen)...)
	}

	return strings.Join(resultLines, "\n")
}

// breakIntoLines splits a string into lines of maximum lineLen characters
func breakIntoLines(content string, lineLen int) []string {
	var lines []string
	for len(content) > 0 {
		if len(content) <= lineLen {
			lines = append(lines, content)
			break
		}
		lines = append(lines, content[:lineLen])
		content = content[lineLen:]
	}
	return lines
}
