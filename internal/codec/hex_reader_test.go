package codec

import (
	"bytes"
	"errors"
	"io"
	"strings"
	"testing"
)

func TestNewHexReader(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []byte
		wantErr  bool
	}{
		{
			name:     "simple hex string",
			input:    "0102030405",
			expected: []byte{0x01, 0x02, 0x03, 0x04, 0x05},
			wantErr:  false,
		},
		{
			name:     "hex with line breaks",
			input:    "010203040506\naabbccddEEFF\n1255",
			expected: []byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff, 0x12, 0x55},
			wantErr:  false,
		},
		{
			name:     "hex with various whitespace",
			input:    "01 02\t03\r\n04 05",
			expected: []byte{0x01, 0x02, 0x03, 0x04, 0x05},
			wantErr:  false,
		},
		{
			name:     "mixed case hex",
			input:    "0aAbBcCdDeEfFf",
			expected: []byte{0x0a, 0xab, 0xbc, 0xcd, 0xde, 0xef, 0xff},
			wantErr:  false,
		},
		{
			name:     "empty input",
			input:    "",
			expected: []byte{},
			wantErr:  false,
		},
		{
			name:     "only whitespace",
			input:    "   \n\t\r\n  ",
			expected: []byte{},
			wantErr:  false,
		},
		{
			name:     "hex with non-hex characters ignored",
			input:    "01g02h03",
			expected: []byte{0x01, 0x02, 0x03},
			wantErr:  false,
		},
		{
			name:     "single hex digit (odd)",
			input:    "0",
			expected: []byte{},
			wantErr:  false, // Should handle gracefully since we process pairs
		},
		{
			name:     "large hex string",
			input:    strings.Repeat("0123456789abcdef", 1000), // 16KB of hex chars
			expected: bytes.Repeat([]byte{0x01, 0x23, 0x45, 0x67, 0x89, 0xab, 0xcd, 0xef}, 1000),
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hexReader, err := NewHexReader([]byte(tt.input))

			if tt.wantErr {
				if !errors.Is(err, io.EOF) {
					t.Errorf("NewHexReader() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}

			// Read all data from the hex reader
			result := make([]byte, len(tt.expected))
			n, err := hexReader.ReadAt(result, 0)

			if err != nil && !errors.Is(err, io.EOF) {
				t.Errorf("ReadAt() error = %v", err)
				return
			}

			if n != len(tt.expected) {
				t.Errorf("ReadAt() read %d bytes, expected %d", n, len(tt.expected))
				return
			}

			if !bytes.Equal(result, tt.expected) {
				t.Errorf("ReadAt() result = %v, expected %v", result, tt.expected)
			}
		})
	}
}

func TestHexReader_ReadAt(t *testing.T) {
	// Create a hex reader with known data
	input := "0123456789abcdef"
	expected := []byte{0x01, 0x23, 0x45, 0x67, 0x89, 0xab, 0xcd, 0xef}

	hexReader, err := NewHexReader([]byte(input))
	if err != nil {
		t.Fatalf("NewHexReader() error = %v", err)
	}

	tests := []struct {
		name     string
		offset   int64
		length   int
		expected []byte
		wantErr  bool
	}{
		{
			name:     "read from beginning",
			offset:   0,
			length:   4,
			expected: []byte{0x01, 0x23, 0x45, 0x67},
			wantErr:  false,
		},
		{
			name:     "read from middle",
			offset:   2,
			length:   3,
			expected: []byte{0x45, 0x67, 0x89},
			wantErr:  false,
		},
		{
			name:     "read single byte",
			offset:   5,
			length:   1,
			expected: []byte{0xab},
			wantErr:  false,
		},
		{
			name:     "read all data",
			offset:   0,
			length:   8,
			expected: expected,
			wantErr:  false,
		},
		{
			name:     "read beyond end",
			offset:   10,
			length:   2,
			expected: []byte{},
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buffer := make([]byte, tt.length)
			n, err := hexReader.ReadAt(buffer, tt.offset)

			if (err != nil) != tt.wantErr {
				t.Errorf("ReadAt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				return
			}

			result := buffer[:n]
			if !bytes.Equal(result, tt.expected) {
				t.Errorf("ReadAt() result = %v, expected %v", result, tt.expected)
			}
		})
	}
}

func TestHexReader_ExampleFromSpec(t *testing.T) {
	// Test the exact example from the user's specification
	input := `010203040506
aabbccddEEFF
1255`

	expected := []byte{
		0x01, 0x02, 0x03, 0x04, 0x05, 0x06, // 010203040506
		0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff, // aabbccddEEFF
		0x12, 0x55, // 1255
	}

	hexReader, err := NewHexReader([]byte(input))
	if err != nil {
		t.Fatalf("NewHexReader() error = %v", err)
	}

	result := make([]byte, len(expected))
	n, err := hexReader.ReadAt(result, 0)
	if err != nil {
		t.Fatalf("ReadAt() error = %v", err)
	}

	if n != len(expected) {
		t.Errorf("ReadAt() read %d bytes, expected %d", n, len(expected))
	}

	if !bytes.Equal(result, expected) {
		t.Errorf("ReadAt() result = %v, expected %v", result, expected)
	}
}
