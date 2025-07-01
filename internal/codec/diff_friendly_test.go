package codec

import (
	"strings"
	"testing"
)

func TestConvertToDiffFriendly(t *testing.T) {
	tests := []struct {
		skipReason string
		name       string
		previous   string
		current    string
		expected   string
	}{
		{
			name:     "empty inputs",
			previous: "",
			current:  "",
			expected: "",
		}, {
			name:     "empty previous, non-empty current",
			previous: "",
			current:  "0123456789012345678901234567890123456789",
			expected: `01234567890123456789012345678901
23456789`,
		}, {
			name:     "non-empty previous, empty current",
			previous: "123456789",
			current:  "",
			expected: "",
		}, {
			name: "identical content",
			previous: `01234567890123456789012345678901
23456789012345678901234567890123`,
			current: `01234567890123456789012345678901
23456789012345678901234567890123`,
			expected: `01234567890123456789012345678901
23456789012345678901234567890123`,
		}, {
			name: "completely different content",
			previous: `01234567890123456789012345678901
23456789012345678901234567890123`,
			current: `aabbccddeeffaabbccddeeffaabbccdd
ddaafffffffffffffffffaafafafafaf`,
			expected: `aabbccddeeffaabbccddeeffaabbccdd
ddaafffffffffffffffffaafafafafaf`,
		}, {
			name: "partial overlap at beginning",
			previous: `01234567890123456789012345678901
23456789012345678901234567890123`,
			current: `01234567890123456789012345678901
ddaafffffffffffffffffaafafafafaf`,
			expected: `01234567890123456789012345678901
ddaafffffffffffffffffaafafafafaf`,
		}, {
			name: "partial overlap at end",
			previous: `01234567890123456789012345678901
23456789012345678901234567890123`,
			current: `ddaafffffffffffffffffaafafafafaf
23456789012345678901234567890123`,
			expected: `ddaafffffffffffffffffaafafafafaf
23456789012345678901234567890123`,
		}, {
			name: "current is subset of previous",
			previous: `01234567890123456789012345678901
ddaafffffffffffffffffaafafafafaf
23456789012345678901234567890123`,
			current:  `ddaafffffffffffffffffaafafafafaf`,
			expected: `ddaafffffffffffffffffaafafafafaf`,
		}, {
			name:     "previous is subset of current",
			previous: `ddaafffffffffffffffffaafafafafaf`,
			current: `01234567890123456789012345678901
ddaafffffffffffffffffaafafafafaf
23456789012345678901234567890123`,
			expected: `01234567890123456789012345678901
ddaafffffffffffffffffaafafafafaf
23456789012345678901234567890123`,
		}, {
			name: "insertion and shifting in the beginning",
			previous: `01234567890123456789012345678901
ddaafffffffffffffffffaafafafafaf
23456789012345678901234567890123`,
			current: `aabb0123456789012345678901234567
8901ddaafffffffffffffffffaafafaf
afaf2345678901234567890123456789
0123`,
			expected: `aabb
01234567890123456789012345678901
ddaafffffffffffffffffaafafafafaf
23456789012345678901234567890123`,
		}, {
			name: "insertion and shifting in the middle",
			previous: `01234567890123456789012345678901
ddaafffffffffffffffffaafafafafaf
23456789012345678901234567890123`,
			current: `01234567890123456FFFF78901234567
8901ddaafffffffffffffffffaafafaf
afaf2345678901234567890123456789
0123`,
			expected: `01234567890123456FFFF78901234567
8901
ddaafffffffffffffffffaafafafafaf
23456789012345678901234567890123`,
		},
		{
			name: "insertion and shifting in the middle with empty lines",
			previous: `01234567890123456789012345678901


ddaafffffffffffffffffaafafafafaf

23456789012345678901234567890123`,
			current: `01234567890123456FFFF78901234567
8901ddaafffffffffffffffffaafafaf
afaf2345678901234567890123456789
0123`,
			expected: `01234567890123456FFFF78901234567
8901
ddaafffffffffffffffffaafafafafaf
23456789012345678901234567890123`,
		}, {
			skipReason: "this initial greedy line breaker algorithm is sub-optimum. TODO: provide a DP approach",
			name:       "multiple combinations",
			previous: `0123456789
0123456789
0123456789
01ddaaffff
ffffffffff
fffaafafaf
afaf234567
8901234567
8901234567
890123`,
			current: `0123456789
0123456FFF
F789012345
678901ddaa
ffffffffff
fffffffaaf
afafafaf23
45678dd901
2345678901
2345678901
23`,
			expected: `0123456789
0123456FFFF789
0123456789
01ddaaffff
ffffffffff
fffaafafaf
afaf234567
89dd01234567
8901234567
890123`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.skipReason != "" {
				t.Skip(tt.skipReason)
			}
			result := ConvertToDiffFriendly(tt.previous, tt.current, 32)
			// check that tt.Current and result are equivalent (no hex digits added or removed)
			resRaw := strings.ReplaceAll(result, "\n", "")
			curRaw := strings.ReplaceAll(tt.current, "\n", "")
			if resRaw != curRaw {
				t.Errorf("generated result is not binary-equal to input:\nwant: %sgot:  %s", curRaw, resRaw)
			}
			if tt.expected != result {
				t.Errorf("ConvertToDiffFriendly():\n      %v\nwant: %v", result, tt.expected)
			}
		})
	}
}
