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
	name string
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

				result.Options[outline.optionCache[index].name] = ""
				continue
			}

			if strings.HasPrefix(args[i], "-") {
				for j := 1; j < len(args[i]); j++ {
					index, ok := outline.optionBinding[args[i][j:j+1]]
					if !ok {
						return Result{}, uopterr.UnrecognizedOption(args[i])
					}

					result.Options[outline.optionCache[index].name] = ""
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
