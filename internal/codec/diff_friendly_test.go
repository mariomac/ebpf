package codec

import (
	"bytes"
	"slices"
	"testing"
)

func TestConvertToDiffFriendly(t *testing.T) {
	tests := []struct {
		name     string
		previous []byte
		current  []byte
		expected []byte
	}{
		{
			name:     "empty inputs",
			previous: []byte{},
			current:  []byte{},
			expected: []byte{},
		}, {
			name:     "nil inputs",
			previous: nil,
			current:  nil,
			expected: []byte{},
		}, {
			name:     "empty previous, non-empty current",
			previous: []byte{},
			current:  []byte("0123456789012345678901234567890123456789"),
			expected: []byte(`01234567890123456789012345678901
23456789`),
		}, {
			name:     "non-empty previous, empty current",
			previous: []byte("123456789"),
			current:  []byte{},
			expected: []byte{},
		}, {
			name: "identical content",
			previous: []byte(`01234567890123456789012345678901
23456789012345678901234567890123`),
			current: []byte(`01234567890123456789012345678901
23456789012345678901234567890123`),
			expected: []byte(`01234567890123456789012345678901
23456789012345678901234567890123`),
		}, {
			name: "completely different content",
			previous: []byte(`01234567890123456789012345678901
23456789012345678901234567890123`),
			current: []byte(`aabbccddeeffaabbccddeeffaabbccdd
ddaafffffffffffffffffaafafafafaf`),
			expected: []byte(`aabbccddeeffaabbccddeeffaabbccdd
ddaafffffffffffffffffaafafafafaf`),
		}, {
			name: "partial overlap at beginning",
			previous: []byte(`01234567890123456789012345678901
23456789012345678901234567890123`),
			current: []byte(`01234567890123456789012345678901
ddaafffffffffffffffffaafafafafaf`),
			expected: []byte(`01234567890123456789012345678901
ddaafffffffffffffffffaafafafafaf`),
		}, {
			name: "partial overlap at end",
			previous: []byte(`01234567890123456789012345678901
23456789012345678901234567890123`),
			current: []byte(`ddaafffffffffffffffffaafafafafaf
23456789012345678901234567890123`),
			expected: []byte(`ddaafffffffffffffffffaafafafafaf
23456789012345678901234567890123`),
		}, {
			name: "current is subset of previous",
			previous: []byte(`01234567890123456789012345678901
ddaafffffffffffffffffaafafafafaf
23456789012345678901234567890123`),
			current:  []byte(`ddaafffffffffffffffffaafafafafaf`),
			expected: []byte(`ddaafffffffffffffffffaafafafafaf`),
		}, {
			name:     "previous is subset of current",
			previous: []byte(`ddaafffffffffffffffffaafafafafaf`),
			current: []byte(`01234567890123456789012345678901
ddaafffffffffffffffffaafafafafaf
23456789012345678901234567890123`),
			expected: []byte(`01234567890123456789012345678901
ddaafffffffffffffffffaafafafafaf
23456789012345678901234567890123`),
		}, {
			name: "insertion and shifting in the beginning",
			previous: []byte(`01234567890123456789012345678901
ddaafffffffffffffffffaafafafafaf
23456789012345678901234567890123`),
			current: []byte(`aabb0123456789012345678901234567
8901ddaafffffffffffffffffaafafaf
afaf2345678901234567890123456789
0123`),
			expected: []byte(`aabb
01234567890123456789012345678901
ddaafffffffffffffffffaafafafafaf
23456789012345678901234567890123`),
		}, {
			name: "insertion and shifting in the middle",
			previous: []byte(`01234567890123456789012345678901
ddaafffffffffffffffffaafafafafaf
23456789012345678901234567890123`),
			current: []byte(`01234567890123456FFFF78901234567
8901ddaafffffffffffffffffaafafaf
afaf2345678901234567890123456789
0123`),
			expected: []byte(`01234567890123456FFFF78901234567
8901
ddaafffffffffffffffffaafafafafaf
23456789012345678901234567890123`),
		},
		{
			name: "insertion and shifting in the middle with empty lines",
			previous: []byte(`01234567890123456789012345678901


ddaafffffffffffffffffaafafafafaf

23456789012345678901234567890123`),
			current: []byte(`01234567890123456FFFF78901234567
8901ddaafffffffffffffffffaafafaf
afaf2345678901234567890123456789
0123`),
			expected: []byte(`01234567890123456FFFF78901234567
8901
ddaafffffffffffffffffaafafafafaf
23456789012345678901234567890123`),
		}, {
			name: "multiple combinations",
			previous: []byte(`0123456789
0123456789
0123456789
01ddaaffff
ffffffffff
fffaafafaf
afaf234567
8901234567
8901234567
890123`),
			current: []byte(`0123456789
0123456FFF
F789012345
678901ddaa
ffffffffff
fffffffaaf
afafafaf23
45678dd901
2345678901
2345678901
23`),
			expected: []byte(`0123456789
0123456FFFF789
0123456789
01ddaaffff
ffffffffff
fffaafafaf
afaf234567
89dd01234567
8901234567
890123`),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ConvertToDiffFriendly(tt.previous, tt.current, 32)
			checkEqualIgnoreSpaces(t, tt.current, result)
			if !bytes.Equal(result, tt.expected) {
				t.Errorf("ConvertToDiffFriendly():\n%v\nwant:\n%v", string(result), string(tt.expected))
			}
		})
	}
}

func checkEqualIgnoreSpaces(t *testing.T, a, b []byte) {
	t.Helper()
	ra, rb := removeNonHex(slices.Clone(a)), removeNonHex(slices.Clone(b))
	if !bytes.Equal(ra, rb) {
		t.Errorf("bytes are not equal:\n%v\n%v", string(ra), string(rb))
	}
}
