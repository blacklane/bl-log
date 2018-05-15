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
	log.SetOut(out)

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
}

func TestLog(t *testing.T) {
	out := bytes.NewBufferString("")
	log.SetOut(out)

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
	log.SetErr(out)

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
	log.SetOut(out)
	log.SetErr(out)

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

	t.Run("description", func(t *testing.T) {
		out.Reset()
		h := log.L(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
		req, _ := http.NewRequest("GET", "/foo?bar=123", nil)
		h.ServeHTTP(new(mockResponseWriter), req)

		rec := struct {
			Desc string `json:"desc"`
		}{}
		json.Unmarshal(out.Bytes(), &rec)

		exp := "code: 200, path: /foo, params: bar=123"
		got := rec.Desc
		if exp != got {
			t.Errorf("Expected %q, got %q", exp, got)
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
			Name string `json:"name"`
			Desc string `json:"desc"`
		}{}
		json.Unmarshal(out.Bytes(), &rec)

		exp := "request_error"
		got := rec.Name
		if exp != got {
			t.Errorf("Expected %q, got %q", exp, got)
		}
		exp = "code: 400, path: /foo, params: bar=123"
		got = rec.Desc
		if exp != got {
			t.Errorf("Expected %q, got %q", exp, got)
		}
	})
}
