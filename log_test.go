package log_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"testing"
	"time"

	log "github.com/blacklane/bl-log"
)

func TestRecord(t *testing.T) {
	out := bytes.NewBufferString("")
	log.Out = out

	t.Run("duration in milliseconds", func(t *testing.T) {
		out.Reset()
		r := log.NewRecord("test")
		time.Sleep(8 * time.Millisecond)
		r.Log("foo")

		rec := struct {
			Duration int `json:"duration"`
		}{}
		json.Unmarshal(out.Bytes(), &rec)

		exp := 8
		got := rec.Duration
		if exp != got {
			t.Errorf("Expected %q, got %q", exp, got)
		}
	})

	t.Run("timestamp in ISO8601", func(t *testing.T) {
		out.Reset()
		r := log.NewRecord("test")
		r.Log("foo")

		rec := struct {
			Time string `json:"timestamp"`
		}{}
		json.Unmarshal(out.Bytes(), &rec)

		_, err := time.Parse(time.RFC3339, rec.Time)
		if err != nil {
			t.Error(err)
		}
	})

	t.Run("response", func(t *testing.T) {
		out.Reset()
		r := log.NewRecord("test_response")

		req, _ := http.NewRequest("GET", "/foo", nil)
		r.Response(&http.Response{StatusCode: 200, Request: req})

		rec := struct {
			Code int    `json:"code"`
			URI  string `json:"uri"`
		}{}
		json.Unmarshal(out.Bytes(), &rec)

		code, path := 200, "/foo"
		if rec.Code != code || rec.URI != path {
			t.Errorf("Expected {%d %s}, got %v", code, path, rec)
		}
	})
}

func TestLog(t *testing.T) {
	out := bytes.NewBufferString("")
	log.Out = out

	t.Run("timestamp in ISO8601", func(t *testing.T) {
		out.Reset()
		log.Log("foo", "bar")

		rec := struct {
			Time string `json:"timestamp"`
		}{}
		json.Unmarshal(out.Bytes(), &rec)

		_, err := time.Parse(time.RFC3339, rec.Time)
		if err != nil {
			t.Error(err)
		}
	})
}
func TestError(t *testing.T) {
	out := bytes.NewBufferString("")
	log.Err = out

	t.Run("timestamp in ISO8601", func(t *testing.T) {
		out.Reset()
		log.Error(errors.New("foo"))

		rec := struct {
			Time string `json:"timestamp"`
		}{}
		json.Unmarshal(out.Bytes(), &rec)

		_, err := time.Parse(time.RFC3339, rec.Time)
		if err != nil {
			t.Error(err)
		}
	})
}

type mockResponseWriter struct{}

func (m *mockResponseWriter) Header() (h http.Header) {
	return http.Header{}
}
func (m *mockResponseWriter) Write(p []byte) (n int, err error) {
	return len(p), nil
}
func (m *mockResponseWriter) WriteHeader(int) {}

func TestRequest(t *testing.T) {
	out := bytes.NewBufferString("")
	log.Out = out
	log.Err = out

	t.Run("event name", func(t *testing.T) {
		out.Reset()
		h := log.L(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
		req, _ := http.NewRequest("GET", "/foo?bar=123", nil)
		h.ServeHTTP(new(mockResponseWriter), req)

		rec := struct {
			Name string `json:"name"`
		}{}
		json.Unmarshal(out.Bytes(), &rec)

		exp := "request_finished"
		got := rec.Name
		if exp != got {
			t.Errorf("Expected %q, got %q", exp, got)
		}
	})

	t.Run("request data", func(t *testing.T) {
		out.Reset()
		h := log.L(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
		req, _ := http.NewRequest("GET", "/foo?bar=123", nil)
		h.ServeHTTP(new(mockResponseWriter), req)

		rec := struct {
			Code   int    `json:"code"`
			URI    string `json:"uri"`
			Params string `json:"params"`
		}{}
		json.Unmarshal(out.Bytes(), &rec)

		code, URI, params := 200, "/foo", "bar=123"
		if rec.Code != code || rec.URI != URI || rec.Params != params {
			t.Errorf("Expected %d %s %s, got %v", code, URI, params, rec)
		}
	})

	t.Run("errored requests", func(t *testing.T) {
		out.Reset()
		h := log.L(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusBadRequest)
			w.Header().Set("X-Foo", "foo")
			w.Write([]byte("Hello"))
		}))
		req, _ := http.NewRequest("GET", "/foo?bar=123", nil)
		h.ServeHTTP(new(mockResponseWriter), req)

		rec := struct {
			Name   string `json:"name"`
			Code   int    `json:"code"`
			URI    string `json:"uri"`
			Params string `json:"params"`
		}{}
		json.Unmarshal(out.Bytes(), &rec)

		exp := "request_error"
		got := rec.Name
		if exp != got {
			t.Errorf("Expected %q, got %q", exp, got)
		}
		code, URI, params := 400, "/foo", "bar=123"
		if rec.Code != code || rec.URI != URI || rec.Params != params {
			t.Errorf("Expected %d %s %s, got %v", code, URI, params, rec)
		}
	})
}

func TestSilence(t *testing.T) {
	out := bytes.NewBufferString("")
	err := bytes.NewBufferString("")
	log.Out = out
	log.Err = err

	log.Silence()

	log.Log("test_out", "foo")
	log.Error(errors.New("foo"))

	if out.String() != "" {
		t.Errorf("Out should be empty, got: %s", out)
	}
	if err.String() != "" {
		t.Errorf("Err should be empty, got: %s", err)
	}
}
