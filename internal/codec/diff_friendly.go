package codec

import (
	"bytes"
	"errors"
	"fmt"
	"io"
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

func ConvertToDiffFriendly(previous, current []byte, lineLen int) []byte {
	// Remove all line breaks to get raw hexadecimal content
	currentRaw := removeNonHex(current)

	// Get lines from previous string
	previousLines := bytes.Split(previous, []byte{'\n'})

	// Remove empty lines
	var nonEmptyPreviousLines [][]byte
	for _, line := range previousLines {
		if len(line) > 0 {
			nonEmptyPreviousLines = append(nonEmptyPreviousLines, line)
		}
	}

	var resultLines [][]byte
	remainingCurrent := currentRaw

	// Try to match lines from previous in order
	for _, previousLine := range nonEmptyPreviousLines {
		if len(previousLine) <= lineLen && bytes.HasPrefix(remainingCurrent, previousLine) {
			// Found a matching line at the beginning
			resultLines = append(resultLines, previousLine)
			remainingCurrent = remainingCurrent[len(previousLine):]
		} else {
			// Try to find the line anywhere in the remaining current content
			index := bytes.Index(remainingCurrent, previousLine)
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
	if len(remainingCurrent) > 0 {
		resultLines = append(resultLines, breakIntoLines(remainingCurrent, lineLen)...)
	}

	return bytes.Join(resultLines, []byte{'\n'})
}

func removeNonHex(current []byte) []byte {
	wi := 0
	for _, b := range current {
		if (b >= '0' && b <= '9') ||
			(b >= 'A' && b <= 'F') ||
			(b >= 'a' && b <= 'f') {
			current[wi] = b
			wi++
		}
	}
	return current[:wi]
}

// breakIntoLines splits a string into lines of maximum lineLen characters
func breakIntoLines(content []byte, lineLen int) [][]byte {
	var lines [][]byte
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
