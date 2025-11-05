package error_reporters

import (
	"fmt"
	"io"
	"os"
)

type writeReporter struct {
	writable io.Writer
}

func (wr *writeReporter) ReportError(err error) {
	fmt.Fprintln(wr.writable, err.Error())
}

type StdoutReporter struct {
	wr writeReporter
}

type StderrReporter struct {
	wr writeReporter
}

func NewStdoutReporter() *StdoutReporter {
	return &StdoutReporter{
		wr: writeReporter{
			writable: os.Stdout,
		},
	}
}

func (e *StdoutReporter) ReportError(err error) {
	e.wr.ReportError(err)
}

func NewStderrReporter() *StderrReporter {
	return &StderrReporter{
		wr: writeReporter{
			writable: os.Stderr,
		},
	}
}

func (e *StderrReporter) ReportError(err error) {
	e.wr.ReportError(err)
}
