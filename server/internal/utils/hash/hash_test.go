package hash

import (
	"errors"
	"testing"
)

func TestHashAndCompareString(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		compare  string
		expected bool
	}{
		{
			name:     "Simple test",
			input:    "Tickr is awesome",
			compare:  "Tickr is awesome",
			expected: true,
		},
		{
			name:     "Different string test",
			input:    "Tickr is awesome",
			compare:  "Tickr_is_awesome",
			expected: false,
		},
		{
			name:     "Empty test",
			input:    "",
			compare:  "",
			expected: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			hashedStr, err := String(tc.input)
			if errors.Is(err, ErrStringTooLong) {

			}
			if err != nil {
				t.Fatalf("failed to hash string: %v", err)
			}
			if hashedStr == "" {
				t.Fatalf("expected token string not to be empty")
			}

			match := CheckString(tc.compare, hashedStr)
			if match != tc.expected {
				t.Errorf("expected match to be %v got %v", tc.expected, match)
			}
		})
	}

}

func TestHashLongString(t *testing.T) {
	tooLongStr := "Tickr is awesome, multi-tenant, concurrent movie ticketing platform with strict role-based access control"

	_, err := String(tooLongStr)
	if !errors.Is(err, ErrStringTooLong) {
		t.Errorf("expected error %v, got %v", ErrStringTooLong, err)
	}
}
