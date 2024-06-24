package uopt

import "strings"

// A Visitor is used to handle the "events" emitted by the Visit function.
type Visitor interface {
	// VisitFlag is called for any argument that appears to be a UNIX flag.
	//The flag argument will be the argument, excepting the leading hyphen(s).
	// Because the Visit function cannot distinguish between arguments intended
	// to be flags or options, it is up to VisitFlag to do so.
	// Returning true hints to Visit that the argument was a flag, returning
	// false indicates that there is a value that should be captured through
	// VisitOption.
	VisitFlag(flag string) (wasFlag bool)

	// VisitOption is called for any argument that VisitFlag returned false for.
	// The option argument, like VisitFlag will be the argument omitting any
	// prefixed hyphen(s). The value will be captured from the arguments given
	// to Visit, see it for more details.
	VisitOption(option string, value string)

	// VisitArgument will be called for any value that isn't a flag, or isn't
	// captured by an option. The exception to this is '--' which is used to
	// indicate that option parsing should stop.
	VisitArgument(argument string)
}

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
func Visit(visitor Visitor, arguments []string) {
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

				if visitor.VisitFlag(opt) {
					continue
				}

				if idx < len(arg) {
					visitor.VisitOption(opt, arg[idx+1:])
					continue
				}

				var value string
				if !last && isOptionValue(arguments[i+1]) {
					value = arguments[i+1]
					i++
				}

				visitor.VisitOption(opt, value)
				continue
			}

			if isShortOption(arg) {
				for j := 1; j < len(arg); j++ {
					opt := arg[j : j+1]

					if !isAlpha(opt[0]) {
						continue
					}

					if visitor.VisitFlag(opt) {
						continue
					}

					if j >= len(arg)-1 {
						if !last && isOptionValue(arguments[i+1]) {
							visitor.VisitOption(opt, arguments[i+1])
							i++
							continue
						}

						visitor.VisitOption(opt, "")
						continue
					}

					visitor.VisitOption(opt, arg[j+1:])
					break
				}

				continue
			}
		}

		visitor.VisitArgument(arg)
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
