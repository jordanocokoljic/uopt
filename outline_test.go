package uopt_test

import (
	"errors"
	"reflect"
	"testing"

	"github.com/jordanocokoljic/uopt"
	"github.com/jordanocokoljic/uopt/internal/uopterr"
)

func TestCommandOutline_ApplyTo(t *testing.T) {
	tests := []struct {
		name      string
		arguments []string
		schema    uopt.CommandSchema
		result    uopt.Result
		err       error
	}{
		{
			name:      "LongOption",
			arguments: []string{"--help"},
			schema: uopt.CommandSchema{
				Options: []uopt.OptionSchema{
					{
						Name: "help",
						Long: "--help",
					},
				},
			},
			result: uopt.Result{
				Options: map[string]string{
					"help": "",
				},
			},
		},
		{
			name:      "UnrecognizedLongOption",
			arguments: []string{"--not-help"},
			schema: uopt.CommandSchema{
				Options: []uopt.OptionSchema{
					{
						Name: "help",
						Long: "--help",
					},
				},
			},
			err: uopterr.UnrecognizedOption("--not-help"),
		},
		{
			name:      "ShortOption",
			arguments: []string{"-h"},
			schema: uopt.CommandSchema{
				Options: []uopt.OptionSchema{
					{
						Name:  "help",
						Short: "-h",
					},
				},
			},
			result: uopt.Result{
				Options: map[string]string{
					"help": "",
				},
			},
		},
		{
			name:      "UnrecognizedShortOption",
			arguments: []string{"-g"},
			schema: uopt.CommandSchema{
				Options: []uopt.OptionSchema{
					{
						Name:  "help",
						Short: "-h",
					},
				},
			},
			err: uopterr.UnrecognizedOption("-g"),
		},
		{
			name:      "SingleArgument",
			arguments: []string{"jordan"},
			schema: uopt.CommandSchema{
				Arguments: []string{
					"name",
				},
			},
			result: uopt.Result{
				Arguments: map[string]string{
					"name": "jordan",
				},
			},
		},
		{
			name:      "NonVariadicExtraArguments",
			arguments: []string{"jordan", "uopt"},
			schema: uopt.CommandSchema{
				Arguments: []string{
					"name",
				},
			},
			err: uopterr.UnrecognizedArgument("uopt"),
		},
		{
			name:      "VariadicExtraArguments",
			arguments: []string{"jordan", "uopt", "golang"},
			schema: uopt.CommandSchema{
				Arguments: []string{
					"name",
				},
				Variadic: true,
			},
			result: uopt.Result{
				Arguments: map[string]string{
					"name": "jordan",
				},
				Variadic: []string{
					"uopt",
					"golang",
				},
			},
		},
		{
			name:      "CombinedShortOptions",
			arguments: []string{"-abc"},
			schema: uopt.CommandSchema{
				Options: []uopt.OptionSchema{
					{
						Name:  "a",
						Short: "-a",
					},
					{
						Name:  "b",
						Short: "-b",
					},
					{
						Name:  "c",
						Short: "-c",
					},
				},
			},
			result: uopt.Result{
				Options: map[string]string{
					"a": "",
					"b": "",
					"c": "",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.schema.Validate()
			if err != nil {
				t.Fatalf("schema validation failed: %s", err)
			}

			result, err := tt.schema.Build().ApplyTo(tt.arguments)
			if !errors.Is(err, tt.err) {
				t.Fatalf("want error %v, got %v", tt.err, err)
			}

			if !reflect.DeepEqual(result, tt.result) {
				t.Fatalf(
					"result does not match:\n\thave %#v\n\twant: %#v",
					result, tt.result)
			}
		})
	}
}
