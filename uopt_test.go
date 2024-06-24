package uopt_test

import (
	"github.com/jordanocokoljic/uopt"
	"reflect"
	"testing"
)

func TestVisit(t *testing.T) {
	arguments := make([]string, 0)
	flags := make([]string, 0)
	options := make(map[string]string)

	visitor := Visitor{
		Argument: func(s string) {
			arguments = append(arguments, s)
		},
		Flag: func(opt string) bool {
			if 'A' <= opt[0] && opt[0] <= 'Z' {
				return false
			}

			flags = append(flags, opt)
			return true
		},
		Option: func(opt string, val string) {
			options[opt] = val
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
	visitor := Visitor{
		Argument: func(argument string) {
			collected = append(collected, argument)
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
	visitor := Visitor{
		Flag: func(option string) bool {
			collected = append(collected, option)
			return true
		},
		Argument: func(_ string) {},
	}

	args := []string{"-abc", "-d", "--efg", "--", "-hij", "-k", "--lmn"}
	expected := []string{"a", "b", "c", "d", "efg"}

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
			args: []string{"-afile.txt", "-z", "image.png", "--in", "oneway", "--out=another", "--blank="},
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
			visitor := Visitor{
				Flag:     func(_ string) bool { return false },
				Argument: func(_ string) {},
				Option: func(option string, value string) {
					collected[option] = value
				},
			}

			uopt.Visit(visitor, test.args)

			if !reflect.DeepEqual(collected, test.expected) {
				t.Errorf("got %v want %v", collected, test.expected)
			}
		})
	}
}

type Visitor struct {
	Argument func(string)
	Flag     func(string) bool
	Option   func(string, string)
}

func (v Visitor) VisitArgument(arg string)      { v.Argument(arg) }
func (v Visitor) VisitFlag(arg string) bool     { return v.Flag(arg) }
func (v Visitor) VisitOption(arg, value string) { v.Option(arg, value) }
