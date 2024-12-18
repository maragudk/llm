package eval

import (
	"os"
	"strings"
	"testing"
)

type helper interface {
	Helper()
}

type skipper interface {
	helper
	SkipNow()
}

// skipIfNotEvaluating skips the test if "go test" is not being run with "-test.run=TestEval".
// Returns whether the test was skipped.
func skipIfNotEvaluating(t skipper) {
	t.Helper()

	for _, arg := range os.Args {
		if strings.HasPrefix(arg, "-test.run=TestEval") {
			return
		}
	}

	t.SkipNow()
}

type runnerSkipper interface {
	skipper
	Run(name string, f func(t *testing.T)) bool
}

// Run an evaluation.
// Behaves similar to [testing.T.Run], except it skips the test if "go test" is not being run with "-test.run=TestEval".
// The evaluation function [f] is passed an [E] to help with scoring, logging, etc.
func Run(t runnerSkipper, name string, f func(e *E)) {
	t.Helper()

	t.Run(name, func(t *testing.T) {
		skipIfNotEvaluating(t)

		f(&E{T: t})
	})
}

type E struct {
	T *testing.T
}

// Score a [Sample] using a [Scorer] and return the [Result].
// This is just a convenience method to make it easier to swap out scorers.
func (e *E) Score(s Sample, scorer Scorer) Result {
	r := scorer(s)
	r.Score.IsValid()
	return r
}

// Log a [Result].
func (e *E) Log(r Result) {
	e.T.Helper()
	e.T.Logf("result=%+v", r)
}
