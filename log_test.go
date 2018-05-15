package log_test

import (
	"bytes"
	"encoding/json"
	"errors"
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
