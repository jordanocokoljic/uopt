uopt
====

uopt is a library for building command line parsers. It works by exposing a
visitor esque pattern that can be used to read and process arguments item
by item.


Install
-------

In order to get started, use `go` to install the library.

   go get github.com/jordanocokoljic/uopt/v2

Then write an implementation of the uopt.Visitor interface. Then provide the
Visitor implementation, and a slice of strings to uopt.Visit and uopt will
parse the provided arguments, calling the methods on the Visitor.
