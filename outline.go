package uopt

import (
	"strings"

	"github.com/jordanocokoljic/uopt/internal/uopterr"
)

type CommandOutline struct {
	arguments []string
	variadic  bool

	optionCache   []optionCacheLine
	optionBinding map[string]int
}

type optionCacheLine struct {
	name    string
	capture bool
}

func (outline CommandOutline) ApplyTo(args []string) (Result, error) {
	var result Result

	if outline.optionCache != nil {
		result.Options = make(map[string]string)
	}

	if outline.arguments != nil {
		result.Arguments = make(map[string]string)
	}

	i := 0
	for ; i < len(args); i++ {
		if result.Options != nil {
			if strings.HasPrefix(args[i], "--") {
				index, ok := outline.optionBinding[args[i][2:]]
				if !ok {
					return Result{}, uopterr.UnrecognizedOption(args[i])
				}

				opt := outline.optionCache[index]

				var value string
				if opt.capture {
					if i+1 >= len(args) {
						return Result{}, uopterr.NoCaptureValue(args[i])
					}

					if strings.HasPrefix(args[i+1], "-") {
						return Result{}, uopterr.NoCaptureValue(args[i])
					}

					value = args[i+1]
					i++
				}

				result.Options[opt.name] = value
				continue
			}

			if strings.HasPrefix(args[i], "-") {
				for j := 1; j < len(args[i]); j++ {
					index, ok := outline.optionBinding[args[i][j:j+1]]
					if !ok {
						return Result{}, uopterr.UnrecognizedOption(args[i])
					}

					opt := outline.optionCache[index]

					if opt.capture {
						if j+1 < len(args[i]) {
							result.Options[opt.name] = args[i][j+1:]
							break
						}

						if j >= len(args[i])-1 && i+1 < len(args) {
							if strings.HasPrefix(args[i+1], "-") {
								return Result{}, uopterr.NoCaptureValue(args[i])
							}

							result.Options[opt.name] = args[i+1]
							i++
							break
						}

						return Result{}, uopterr.NoCaptureValue(args[i])
					}

					result.Options[opt.name] = ""
				}

				continue
			}
		}

		if result.Arguments != nil {
			if len(result.Arguments) < len(outline.arguments) {
				name := outline.arguments[len(result.Arguments)]
				result.Arguments[name] = args[i]
				continue
			}
		}

		break
	}

	if i < len(args) {
		if !outline.variadic {
			return Result{}, uopterr.UnrecognizedArgument(args[i])
		}

		result.Variadic = args[i:]
	}

	return result, nil
}

type Result struct {
	Options   map[string]string
	Arguments map[string]string
	Variadic  []string
}
