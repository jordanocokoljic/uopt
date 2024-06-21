package uopt

import (
	"fmt"
	"slices"
	"strings"

	"github.com/jordanocokoljic/uopt/internal/uopterr"
)

type CommandSchema struct {
	Arguments []string
	Options   []OptionSchema
	Variadic  bool
}

type OptionSchema struct {
	Name    string
	Short   string
	Long    string
	Capture bool
}

func (schema *CommandSchema) Validate() error {
	if len(schema.Arguments) > 0 {
		for i, argument := range schema.Arguments {
			if slices.Contains(schema.Arguments[:i], argument) {
				return fmt.Errorf(
					"argument validation failed: %w",
					uopterr.DuplicateName(argument))
			}
		}
	}

	if len(schema.Options) > 0 {
		names := make(map[string]struct{})
		flags := make(map[string]struct{})

		for _, option := range schema.Options {
			if _, ok := names[option.Name]; ok {
				return fmt.Errorf(
					"option validation failed: %w",
					uopterr.DuplicateName(option.Name))
			}

			names[option.Name] = struct{}{}

			if option.Short != "" {
				short := option.Short

				if _, ok := flags[short]; ok {
					return fmt.Errorf(
						"option validation failed: %w",
						uopterr.DuplicateFlag(short))
				}

				if !validShortFlag(short) {
					return uopterr.InvalidShortFlag(short)
				}

				flags[short] = struct{}{}
			}

			if option.Long != "" {
				long := option.Long

				if _, ok := flags[long]; ok {
					return fmt.Errorf(
						"option validation failed: %w",
						uopterr.DuplicateFlag(long))
				}

				if !validLongFlag(long) {
					return fmt.Errorf(
						"option validation failed: %w",
						uopterr.InvalidLongFlag(long))
				}

				flags[long] = struct{}{}
			}

			if option.Short == "" && option.Long == "" {
				return uopterr.NoFlag(option.Name)
			}
		}
	}

	return nil
}

func (schema *CommandSchema) Build() CommandOutline {
	var outline CommandOutline

	if schema.Arguments != nil && len(schema.Arguments) > 0 {
		outline.arguments = schema.Arguments
	}

	if schema.Options != nil && len(schema.Options) > 0 {
		outline.optionCache = make([]optionCacheLine, len(schema.Options))
		outline.optionBinding = make(map[string]int)

		for i, option := range schema.Options {
			outline.optionCache[i] = optionCacheLine{
				name:    option.Name,
				capture: option.Capture,
			}

			if option.Short != "" {
				outline.optionBinding[option.Short[1:]] = i
			}

			if option.Long != "" {
				outline.optionBinding[option.Long[2:]] = i
			}
		}
	}

	outline.variadic = schema.Variadic

	return outline
}

func validShortFlag(short string) bool {
	return len(short) == 2 && short[0] == '-' && isLetter(short[1])
}

func validLongFlag(long string) bool {
	return len(long) >= 3 &&
		strings.HasPrefix(long, "--") &&
		isLetter(long[2]) &&
		!strings.Contains(long, " ")
}

func isLetter(b byte) bool {
	return ('a' <= b && b <= 'z') || ('A' <= b && b <= 'Z')
}
