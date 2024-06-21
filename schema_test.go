package uopt_test

import (
	"errors"
	"testing"

	"github.com/jordanocokoljic/uopt"
	"github.com/jordanocokoljic/uopt/internal/uopterr"
)

func TestCommandSchema_Validate(t *testing.T) {
	tests := []struct {
		name   string
		schema uopt.CommandSchema
		err    error
	}{
		{
			name: "DuplicateArgument",
			schema: uopt.CommandSchema{
				Arguments: []string{
					"arg1",
					"arg2",
					"arg1",
				},
			},
			err: uopterr.DuplicateName("arg1"),
		},
		{
			name: "DuplicateOption",
			schema: uopt.CommandSchema{
				Options: []uopt.OptionSchema{
					{
						Name:  "opt1",
						Short: "-a",
					},
					{
						Name:  "opt2",
						Short: "-b",
					},
					{
						Name:  "opt1",
						Short: "-c",
					},
				},
			},
			err: uopterr.DuplicateName("opt1"),
		},
		{
			name: "OptionTooLong",
			schema: uopt.CommandSchema{
				Options: []uopt.OptionSchema{
					{
						Name:  "opt1",
						Short: "-ab",
					},
				},
			},
			err: uopterr.InvalidShortFlag("-ab"),
		},
		{
			name: "ShortMissingHyphen",
			schema: uopt.CommandSchema{
				Options: []uopt.OptionSchema{
					{
						Name:  "opt1",
						Short: "ab",
					},
				},
			},
			err: uopterr.InvalidShortFlag("ab"),
		},
		{
			name: "ShortNonLetter",
			schema: uopt.CommandSchema{
				Options: []uopt.OptionSchema{
					{
						Name:  "opt1",
						Short: "-0",
					},
				},
			},
			err: uopterr.InvalidShortFlag("-0"),
		},
		{
			name: "LongTooShort",
			schema: uopt.CommandSchema{
				Options: []uopt.OptionSchema{
					{
						Name: "opt1",
						Long: "--",
					},
				},
			},
			err: uopterr.InvalidLongFlag("--"),
		},
		{
			name: "LongBadPrefix",
			schema: uopt.CommandSchema{
				Options: []uopt.OptionSchema{
					{
						Name: "opt1",
						Long: "long",
					},
				},
			},
			err: uopterr.InvalidLongFlag("long"),
		},
		{
			name: "LongBadStart",
			schema: uopt.CommandSchema{
				Options: []uopt.OptionSchema{
					{
						Name: "opt1",
						Long: "--0no",
					},
				},
			},
			err: uopterr.InvalidLongFlag("--0no"),
		},
		{
			name: "LongHasSpace",
			schema: uopt.CommandSchema{
				Options: []uopt.OptionSchema{
					{
						Name: "opt1",
						Long: "--long name",
					},
				},
			},
			err: uopterr.InvalidLongFlag("--long name"),
		},
		{
			name: "NoFlags",
			schema: uopt.CommandSchema{
				Options: []uopt.OptionSchema{
					{
						Name: "opt1",
					},
				},
			},
			err: uopterr.NoFlag("opt1"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.schema.Validate()
			if err == nil {
				t.Fatal("expected error, got nil")
			}

			if !errors.Is(err, tt.err) {
				t.Errorf("expected: %s\ngot: %s", tt.err, err)
			}
		})
	}
}
