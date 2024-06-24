package uopt

import "strings"

type Visitor interface {
	VisitArgument(argument string)
	VisitFlag(flag string) (wasFlag bool)
	VisitOption(option string, value string)
}

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
	return strings.HasPrefix(arg, "--") && arg != "--"
}

func isShortOption(arg string) bool {
	return strings.HasPrefix(arg, "-") && arg != "-"
}

func isOptionValue(arg string) bool {
	return arg != "-" && !isLongOption(arg) && !isShortOption(arg)
}

func isAlpha(str byte) bool {
	return ('a' <= str && str <= 'z') || ('A' <= str && str <= 'Z')
}
