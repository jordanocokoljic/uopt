package uopt_test

import (
	"reflect"
	"testing"

	"github.com/jordanocokoljic/uopt/v2"
)

func TestVisit(t *testing.T) {
	arguments := make([]string, 0)
	flags := make([]string, 0)
	options := make(map[string]string)

	visitor := testVisitor{
		Argument: func(s string) error {
			arguments = append(arguments, s)
			return nil
		},
		Flag: func(opt string) error {
			if 'A' <= opt[0] && opt[0] <= 'Z' {
				return uopt.IsOption
			}

			flags = append(flags, opt)
			return nil
		},
		Option: func(opt string, val string) error {
			options[opt] = val
			return nil
		},
	}

	args := []string{
		"-a",
		"-B", "something",
		"-",
		"-dbS", "5432",
		"-Jfile.txt",
		"--store",
		"--Remove-from", "temp",
		"in.sql",
		"--Notfound=error",
		"-xfGcapture-this",
		"--",
		"-Z",
		"--zero",
	}

	expectedArguments := []string{
		"-",
		"in.sql",
		"-Z",
		"--zero",
	}

	expectedFlags := []string{
		"a",
		"d",
		"b",
		"store",
		"x",
		"f",
	}

	expectedOptions := map[string]string{
		"B":           "something",
		"S":           "5432",
		"J":           "file.txt",
		"Remove-from": "temp",
		"Notfound":    "error",
		"G":           "capture-this",
	}

	uopt.Visit(visitor, args)

	if !reflect.DeepEqual(expectedArguments, arguments) {
		t.Errorf("got %v, want %v", arguments, expectedArguments)
	}

	if !reflect.DeepEqual(expectedFlags, flags) {
		t.Errorf("got %v, want %v", flags, expectedFlags)
	}

	if !reflect.DeepEqual(expectedOptions, options) {
		t.Errorf("got %v, want %v", options, expectedOptions)
	}
}

func TestVisit_Arguments(t *testing.T) {
	var collected []string
	visitor := testVisitor{
		Argument: func(argument string) error {
			collected = append(collected, argument)
			return nil
		},
	}

	args := []string{"abc", "def", "--", "-"}
	expected := []string{"abc", "def", "-"}

	uopt.Visit(visitor, args)

	if !reflect.DeepEqual(collected, expected) {
		t.Errorf("got %v, want %v", expected, collected)
	}
}

func TestVisit_Flags(t *testing.T) {
	var collected []string
	visitor := testVisitor{
		Flag: func(option string) error {
			collected = append(collected, option)
			return nil
		},
		Argument: func(_ string) error { return nil },
	}

	args := []string{
		"-abc",
		"-d",
		"--efg",
		"-zX1y-wv",
		"--",
		"-hij",
		"-k",
		"--lmn",
	}

	expected := []string{
		"a",
		"b",
		"c",
		"d",
		"efg",
		"z",
		"X",
		"y",
		"w",
		"v",
	}

	uopt.Visit(visitor, args)

	if !reflect.DeepEqual(collected, expected) {
		t.Errorf("got %v, want %v", collected, expected)
	}
}

func TestVisit_Options(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		expected map[string]string
	}{
		{
			name: "Basic",
			args: []string{
				"-afile.txt",
				"-z", "image.png",
				"--in", "oneway",
				"--out=another",
				"--blank=",
			},
			expected: map[string]string{
				"a":     "file.txt",
				"z":     "image.png",
				"in":    "oneway",
				"out":   "another",
				"blank": "",
			},
		},
		{
			name: "LongMissingValue",
			args: []string{"--first", "--last"},
			expected: map[string]string{
				"first": "",
				"last":  "",
			},
		},
		{
			name: "ShortMissingValue",
			args: []string{"-a", "-z"},
			expected: map[string]string{
				"a": "",
				"z": "",
			},
		},
		{
			name: "DashIsNotCaptured",
			args: []string{"--first", "-", "-a", "-"},
			expected: map[string]string{
				"first": "",
				"a":     "",
			},
		},
		{
			name: "ShortIgnoresDoubleDash",
			args: []string{"-a", "--"},
			expected: map[string]string{
				"a": "",
			},
		},
		{
			name: "LongIgnoresDoubleDash",
			args: []string{"--first", "--"},
			expected: map[string]string{
				"first": "",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			collected := make(map[string]string)
			visitor := testVisitor{
				Flag:     func(_ string) error { return uopt.IsOption },
				Argument: func(_ string) error { return nil },
				Option: func(option string, value string) error {
					collected[option] = value
					return nil
				},
			}

			uopt.Visit(visitor, test.args)

			if !reflect.DeepEqual(collected, test.expected) {
				t.Errorf("got %v want %v", collected, test.expected)
			}
		})
	}
}

func TestVisit_HandlesArgumentsThatLookLikeOptions(t *testing.T) {
	var collected []string
	visitor := testVisitor{
		Flag:   func(_ string) error { return nil },
		Option: func(_ string, _ string) error { return nil },
		Argument: func(arg string) error {
			collected = append(collected, arg)
			return nil
		},
	}

	arguments := []string{"-1", "-.", "--1", "--."}
	expected := []string{"-1", "-.", "--1", "--."}

	uopt.Visit(visitor, arguments)

	if !reflect.DeepEqual(collected, expected) {
		t.Errorf("got %v, want %v", collected, expected)
	}
}

func TestVisit_TerminatesOnHalt(t *testing.T) {
	tests := []struct {
		name      string
		arguments []string
		index     int
	}{
		{
			name:      "ForArgument",
			arguments: []string{"halt", "--hello=world"},
		},
		{
			name:      "ForAggregateFlag",
			arguments: []string{"-gh", "--hello=world"},
			index:     1,
		},
		{
			name:      "ForStandaloneShortFlag",
			arguments: []string{"-h", "--hello=world"},
		},
		{
			name:      "ForStandaloneLongFlag",
			arguments: []string{"--halt", "--hello=world"},
		},
		{
			name:      "ForAggregateShortOption",
			arguments: []string{"-chalt", "--hello=world"},
		},
		{
			name:      "ForSeparateShortOption",
			arguments: []string{"-c", "halt", "--hello=world"},
		},
		{
			name:      "ForSeparateLongOption",
			arguments: []string{"--cmd", "halt", "--hello=world"},
		},
		{
			name:      "ForEqualsLongOption",
			arguments: []string{"--cmd=halt", "--hello=world"},
		},
		{
			name:      "ForWhenShortOptionIsIncomplete",
			arguments: []string{"-b", "-c", "--hello=world"},
		},
		{
			name:      "ForWhenLongOptionIsIncomplete",
			arguments: []string{"--blank", "-c", "--hello=world"},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var count int

			// This visitor will halt on any of the following conditions:
			//
			// The argument 'halt' is provided.
			// The flags -h or --halt are provided
			// The option -c or --cmd are provided with the value 'halt'
			// The option -b or --blank are provided and an empty string is the value
			visitor := testVisitor{
				Argument: func(a string) error {
					count++

					if a == "halt" {
						return uopt.Halt
					}

					return nil
				},
				Flag: func(f string) error {
					if f == "h" || f == "halt" {
						count++
						return uopt.Halt
					}

					if f == "c" || f == "cmd" || f == "b" || f == "blank" {
						return uopt.IsOption
					}

					count++
					return nil
				},
				Option: func(o, v string) error {
					count++

					if (o == "c" || o == "cmd") && v == "halt" {
						return uopt.Halt
					}

					if (o == "b" || o == "blank") && v == "" {
						return uopt.Halt
					}

					return nil
				},
			}

			err := uopt.Visit(visitor, test.arguments)
			if err != nil {
				t.Fatalf("visit returned an error: %s", err.Error())
			}

			if count-1 != test.index {
				t.Errorf(
					"argument count mismatch: %d vs %d",
					count, test.index,
				)
			}
		})
	}
}

type testVisitor struct {
	Flag     func(f string) error
	Option   func(o string, v string) error
	Argument func(a string) error
}

func (t testVisitor) VisitFlag(f string) error      { return t.Flag(f) }
func (t testVisitor) VisitOption(o, v string) error { return t.Option(o, v) }
func (t testVisitor) VisitArgument(a string) error  { return t.Argument(a) }
