package main

type ErrorReporter[E error] interface {
	ReportError(E)
}
