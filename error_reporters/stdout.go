package error_reporters

import (
	"fmt"
	"os"
)

type StdoutReporter struct{}

func NewStdoutReporter() *StdoutReporter {
	return &StdoutReporter{}
}

func (e *StdoutReporter) ReportError(err error) {
	fmt.Fprintln(os.Stdout, err.Error())
}
