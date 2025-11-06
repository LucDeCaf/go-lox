package error_reporters

type ErrorReporter[E error] interface {
	ReportError(E)
}
