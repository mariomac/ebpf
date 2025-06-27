package codec

import (
	"strings"
	"testing"
)

func TestHexWriterWrap_Write(t *testing.T) {
	tests := []struct {
		name  string
		input []byte
		want  string
	}{
		{
			name:  "empty",
			input: []byte{},
			want:  "",
		},
		{
			name:  "one byte",
			input: []byte{0xde},
			want:  "de",
		},
		{
			name:  "many bytes",
			input: []byte{0xde, 0xad, 0xbe, 0xef, 0x00, 0x11, 0x22, 0x33},
			want:  "deadbeef00112233",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var b strings.Builder
			hw := NewHexWriterWrapper(&b)

			n, err := hw.Write(tt.input)
			if err != nil {
				t.Fatal(err)
			}
			if n != len(tt.input) {
				t.Errorf("expected to write %d bytes, got %d", len(tt.input), n)
			}

			if b.String() != tt.want {
				t.Errorf("got %q, want %q", b.String(), tt.want)
			}
		})
	}
}
