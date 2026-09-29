package middleware

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

// topLevelKeys returns the keys of a JSON object, duplicates included.
func topLevelKeys(t *testing.T, line string) []string {
	t.Helper()
	dec := json.NewDecoder(strings.NewReader(line))
	if _, err := dec.Token(); err != nil {
		t.Fatal(err)
	}
	var keys []string
	for dec.More() {
		key, err := dec.Token()
		if err != nil {
			t.Fatal(err)
		}
		keys = append(keys, key.(string))
		var value json.RawMessage
		if err := dec.Decode(&value); err != nil {
			t.Fatal(err)
		}
	}
	return keys
}

func TestNewLoggerDoesNotAccumulateAttributes(t *testing.T) {
	gin.SetMode(gin.TestMode)

	var buf bytes.Buffer
	logger := slog.New(slog.NewJSONHandler(&buf, nil))

	router := gin.New()
	router.Use(NewLogger(logger))
	router.GET("/ping", func(c *gin.Context) { c.Status(http.StatusOK) })

	for range 5 {
		router.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/ping", nil))
	}

	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) != 5 {
		t.Fatalf("expected 5 log lines, got %d", len(lines))
	}
	want := len(topLevelKeys(t, lines[0]))
	for i, line := range lines {
		keys := topLevelKeys(t, line)
		if len(keys) != want {
			t.Fatalf("line %d: expected %d keys, got %d: %s", i, want, len(keys), line)
		}
		seen := map[string]bool{}
		for _, k := range keys {
			if seen[k] {
				t.Fatalf("line %d: duplicated key %q: %s", i, k, line)
			}
			seen[k] = true
		}
	}
}
