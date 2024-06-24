package uopt

import (
	"errors"
	"strings"
)

// A Visitor is used to handle the "events" emitted by the Visit function.
type Visitor interface {
	// VisitFlag is called for any argument that appears to be a UNIX flag.
	//The flag argument will be the argument, excepting the leading hyphen(s).
	// Because the Visit function cannot distinguish between arguments intended
	// to be flags or options, it is up to VisitFlag to do so.
	// Returning nil hints to Visit that the argument was a flag, returning
	// IsOption indicates that there is a value that should be captured through
	// VisitOption.
	VisitFlag(flag string) error

	// VisitOption is called for any argument where VisitFlag returned IsOption.
	// The option argument, like VisitFlag will be the argument omitting any
	// prefixed hyphen(s). The value will be captured from the arguments given
	// to Visit, see it for more details.
	VisitOption(option string, value string) error

	// VisitArgument will be called for any value that isn't a flag, or isn't
	// captured by an option. The exception to this is '--' which is used to
	// indicate that option parsing should stop.
	VisitArgument(argument string) error
}

const (
	// IsOption can be returned by calls to Visitor.VisitFlag to indicate that
	// the argument should instead be handled like an option.
	IsOption visitError = iota
)

// Visit will step through each of the arguments calling the appropriate
// methods on the Visitor.
// Visit attempts to adhere to the principle of least surprise encountering
// flags and options by following widely adopted UNIX conventions.
//
// When Visit encounters values that could be a UNIX option, it will do one of
// the following:
//
// If the option is long (prefixed with '--') it will check if an '=' is also
// present. If it is, Visitor.VisitFlag will be called containing all the
// characters after '--' and before '='. If it returns true, then characters
// after '=' are ignored and Visit will move onto the next argument. If it
// returns false, Visitor.VisitOption will be called with all the characters
// after the '='. If there is no '=' present, and Visitor.VisitFlag returns
// false, then Visitor.VisitOption will instead be called with the next
// argument as the value, or an empty string if there are no more.
//
// If the option is short (prefixed with '-') it will begin to iterate over
// each of the characters in the group, calling Visitor.VisitFlag for each one.
// If any call returns false, it will see if there are any characters left to
// iterate over. If there are, Visitor.VisitOption will be called with the
// remaining characters as the value. If not, then Visitor.VisitOption will be
// called with the next argument as the value, or an empty string if there are
// none.
//
// If while iterating over a group of short options, Visit encounters a
// non-alphabetic character, it will simply be ignored.
//
// Visit will only return an error if the Visitor returns an error that isn't
// IsOption.
func Visit(visitor Visitor, arguments []string) error {
	visitOption := true

	for i := 0; i < len(arguments); i++ {
		arg := arguments[i]
		last := i+1 == len(arguments)

		if visitOption {
			if arg == "--" {
				visitOption = false
				continue
			}

			if isLongOption(arg) {
				idx := strings.Index(arg, "=")
				if idx == -1 {
					idx = len(arg)
				}

				opt := arg[2:idx]

				err := visitor.VisitFlag(opt)
				if err == nil {
					continue
				}

				if !errors.Is(err, IsOption) {
					return err
				}

				if idx < len(arg) {
					err = visitor.VisitOption(opt, arg[idx+1:])
					if err != nil {
						return err
					}

					continue
				}

				var value string
				if !last && isOptionValue(arguments[i+1]) {
					value = arguments[i+1]
					i++
				}

				err = visitor.VisitOption(opt, value)
				if err != nil {
					return err
				}

				continue
			}

			if isShortOption(arg) {
				for j := 1; j < len(arg); j++ {
					opt := arg[j : j+1]

					if !isAlpha(opt[0]) {
						continue
					}

					err := visitor.VisitFlag(opt)
					if err == nil {
						continue
					}

					if !errors.Is(err, IsOption) {
						return err
					}

					if j >= len(arg)-1 {
						if !last && isOptionValue(arguments[i+1]) {
							err = visitor.VisitOption(opt, arguments[i+1])
							if err != nil {
								return err
							}

							i++
							continue
						}

						err = visitor.VisitOption(opt, "")
						if err != nil {
							return err
						}

						continue
					}

					err = visitor.VisitOption(opt, arg[j+1:])
					if err != nil {
						return err
					}

					break
				}

				continue
			}
		}

		err := visitor.VisitArgument(arg)
		if err != nil {
			return err
		}
	}

	return nil
}

type visitError int

func (e visitError) Error() string {
	switch {
	case errors.Is(e, IsOption):
		return "flag should have been interpreted as an option"
	default:
		panic("unsupported error code")
	}
}

func isLongOption(arg string) bool {
	return strings.HasPrefix(arg, "--") &&
		len(arg) > 2 && isAlpha(arg[2])
}

func isShortOption(arg string) bool {
	return strings.HasPrefix(arg, "-") &&
		len(arg) > 1 && isAlpha(arg[1])
}

func isOptionValue(arg string) bool {
	return arg != "-" && arg != "--" &&
		!isLongOption(arg) && !isShortOption(arg)
}

func isAlpha(str byte) bool {
	return ('a' <= str && str <= 'z') || ('A' <= str && str <= 'Z')
}
