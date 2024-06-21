package uopterr_test

import (
	"testing"

	"github.com/jordanocokoljic/uopt/internal/uopterr"
)

func TestErrorMessages(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected string
	}{
		{
			name:     "DuplicateName",
			err:      uopterr.DuplicateName("arg"),
			expected: "name was already registered: arg",
		},
		{
			name:     "DuplicateFlag",
			err:      uopterr.DuplicateFlag("--flag"),
			expected: "flag was already registered: --flag",
		},
		{
			name:     "InvalidShortFlag",
			err:      uopterr.InvalidShortFlag("-0"),
			expected: "flag must be a hyphen followed by 1 alphabetic character: -0",
		},
		{
			name:     "InvalidLongFlag",
			err:      uopterr.InvalidLongFlag("--0a"),
			expected: "flag must be two hyphens followed by an alphabetic character, then any number of alphanumeric characters: --0a",
		},
		{
			name:     "NoFlag",
			err:      uopterr.NoFlag("opt"),
			expected: "option must have a short or long flag: opt",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := tt.err.Error()

			if msg != tt.expected {
				t.Errorf("expected: %s\ngot: %s", tt.expected, msg)
			}
		})
	}
}
