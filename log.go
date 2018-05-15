// Package log is a simple logging strategy to output basic JSON to the stdout/stderr
package log

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

var stdOut io.Writer = os.Stdout
var errOut io.Writer = os.Stderr

func SetOut(o io.Writer) {
	stdOut = o
}
func SetErr(o io.Writer) {
	errOut = o
}

// Log writes a JSON message to the configured standard output
func Log(name, description string, parts ...interface{}) {
	m := fmt.Sprintf(description, parts...)
	fmt.Fprintf(stdOut, `{"name": %q, "desc": %q, "timestamp": %q}`, name, m, time.Now().Format(time.RFC3339))
	os.Stdout.Write([]byte{'\n'})
}

// Error writes an error as JSON to the configured error output
func Error(err error) {
	if err != nil {
		fmt.Fprintf(errOut, `{"error": %q, "timestamp": %q}`, err.Error(), time.Now().Format(time.RFC3339))
		errOut.Write([]byte{'\n'})
	}
}

// Record measures the duration of an event when logged
type Record struct {
	start time.Time
	name  string
	out   io.Writer
	err   io.Writer
}

// Log prints a simple json structure to the stdout with duration information since the record was created
func (r *Record) Log(description string, parts ...interface{}) {
	d := fmt.Sprintf(description, parts...)
	dur := time.Since(r.start) / time.Millisecond
	fmt.Fprintf(r.out, `{"name": %q, "desc": %q, "duration": %d, "timestamp": %q}`, r.name, d, dur, time.Now().Format(time.RFC3339))
	r.out.Write([]byte{'\n'})
}

// NewRecord creates a Record with a name
func NewRecord(name string) *Record {
	return &Record{time.Now(), name, stdOut, errOut}
}

// L log request information and duration
func L(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cr := codeRecorder{w, http.StatusOK}
		rf := NewRecord("request_finished")
		re := NewRecord("request_error")
		h.ServeHTTP(&cr, r)
		desc := fmt.Sprintf("code: %d, path: %s, params: %s", cr.code, r.URL.Path, r.URL.RawQuery)
		if cr.code >= http.StatusBadRequest {
			re.Log(desc)
		} else {
			rf.Log(desc)
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
