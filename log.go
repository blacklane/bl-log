// Package log is a simple logging strategy to output basic JSON to the stdout/stderr
package log

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

// Out is the default output used to log messages
var Out io.Writer = os.Stdout

// Err is the writer used to log errors with Error
var Err io.Writer = os.Stderr

// Log writes a JSON message to the configured standard output
func Log(name, desc string) {
	fmt.Fprintf(Out, `{"name": %q, "desc": %q, "timestamp": %q}
`, name, desc, formattedNow())
}

// Error writes an error as JSON to the configured error output
func Error(err error) {
	if err != nil {
		fmt.Fprintf(Err, `{"error": %q, "timestamp": %q}
`, err.Error(), formattedNow())
	}
}

// Record measures the duration of an event when logged
type Record struct {
	start time.Time
	name  string
}

// Log prints a simple json structure to the stdout with duration information since the record was created
func (r *Record) Log(desc string) {
	Duration(r.name, msSince(r.start), desc)
}

// NewRecord creates a Record with a name
func NewRecord(name string) *Record {
	return &Record{time.Now(), name}
}

func Duration(name string, dur time.Duration, desc string) {
	fmt.Fprintf(Out, `{"name": %q, "desc": %q, "duration": %d, "timestamp": %q}
`, name, desc, dur, formattedNow())
}

// Response logs relevant information about the request/response
func Response(name string, res *http.Response, duration time.Duration) {
	req := res.Request
	fmt.Fprintf(Out, `{"name": %q, "code": %d, "uri": %q, "params": %q, "duration": %d, "timestamp": %q}
`, name, res.StatusCode, req.URL.Path, req.URL.RawQuery, duration, formattedNow())
}

// L log request information and duration
func L(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cr := codeRecorder{w, http.StatusOK}
		start := time.Now()
		h.ServeHTTP(&cr, r)
		res := &http.Response{StatusCode: cr.code, Request: r}
		if res.StatusCode >= http.StatusBadRequest {
			Response("request_error", res, msSince(start))
		} else {
			Response("request_finished", res, msSince(start))
		}
	})
}

type codeRecorder struct {
	rw   http.ResponseWriter
	code int
}

func (cr *codeRecorder) Header() http.Header {
	return cr.rw.Header()
}
func (cr *codeRecorder) Write(p []byte) (int, error) {
	return cr.rw.Write(p)
}
func (cr *codeRecorder) WriteHeader(code int) {
	cr.code = code
	cr.rw.WriteHeader(code)
}

// Silence sets the normal and error outputs to the Noop writer so nothing gets logged, useful for use in tests.
func Silence() {
	Out = Noop
	Err = Noop
}

// Noop is an io.Writer which does nothing
var Noop = &dummyWriter{}

type dummyWriter struct{}

func (w *dummyWriter) Write(p []byte) (int, error) {
	return len(p), nil
}

func msSince(t time.Time) time.Duration {
	return time.Since(t) / time.Millisecond
}

func formattedNow() string {
	return time.Now().Format(time.RFC3339)
}
