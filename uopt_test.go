package uopt_test

import (
	"github.com/jordanocokoljic/uopt"
	"testing"
)

var benchResult uopt.Result
var benchError error

var benchSchema = uopt.CommandSchema{
	Options: []uopt.OptionSchema{
		{
			Name:  "help",
			Short: "-h",
			Long:  "--help",
		},
		{
			Name: "version",
			Long: "--version",
		},
		{
			Name:  "verbose",
			Short: "-v",
			Long:  "--verbose",
		},
		{
			Name:  "stream",
			Short: "-s",
			Long:  "--stream",
		},
		{
			Name: "debug",
			Long: "--debug",
		},
	},
	Arguments: []string{
		"src_file",
		"out_file",
		"log_file",
		"debug_file",
	},
}

var benchArguments = []string{
	"--verbose",
	"-s",
	"--debug",
	"in.file",
	"out.file",
	"log.file",
	"debug_file",
}

func Benchmark_ValidateBuildApply(b *testing.B) {
	// Historically this produces about
	// 1400~1500 ns/op, 1056 B/op, 7 allocs/op

	for i := 0; i < b.N; i++ {
		benchError = benchSchema.Validate()
		outline := benchSchema.Build()
		benchResult, benchError = outline.ApplyTo(benchArguments)
	}
}

func Benchmark_BuildApply(b *testing.B) {
	// Historically this produces about
	// 800~900 ns/op, 1056 B/op, 7 allocs/op

	for i := 0; i < b.N; i++ {
		outline := benchSchema.Build()
		benchResult, benchError = outline.ApplyTo(benchArguments)
	}
}

func Benchmark_Apply(b *testing.B) {
	// Historically this produces about
	// 400~500 ns/op, 672 B/op, 4 allocs/op

	outline := benchSchema.Build()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		benchResult, benchError = outline.ApplyTo(benchArguments)
	}
}
