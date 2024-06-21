package uopt

type CommandOutline struct {
	arguments []string

	optionCache   []optionCacheLine
	optionBinding map[string]int
}

type optionCacheLine struct {
	name string
}

func (outline CommandOutline) ApplyTo(args []string) (Result, error) {
	return Result{}, nil
}

type Result struct{}
