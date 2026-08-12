// SPDX-License-Identifier: 0BSD

package server

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestHomeIsDynamicEscapedAndNotCached(t *testing.T) {
	base := time.Date(2026, time.August, 11, 14, 30, 0, 0, time.UTC)
	clockCalls := 0
	handler := NewWithClock(func() time.Time {
		result := base.Add(time.Duration(clockCalls) * time.Second)
		clockCalls++
		return result
	})

	first := request(t, handler, http.MethodGet, `/?name=%3Cscript%3Ealert%281%29%3C%2Fscript%3E`)
	second := request(t, handler, http.MethodGet, "/?name=friend")

	if first.Code != http.StatusOK || second.Code != http.StatusOK {
		t.Fatalf("unexpected statuses: first=%d second=%d", first.Code, second.Code)
	}
	if first.Header().Get("Cache-Control") != "no-store" {
		t.Fatalf("Cache-Control = %q", first.Header().Get("Cache-Control"))
	}
	if timing := first.Header().Get("Server-Timing"); !strings.HasPrefix(timing, "sando;dur=") || !strings.Contains(timing, "buffered component render") {
		t.Fatalf("Server-Timing = %q", timing)
	}
	if strings.Contains(first.Body.String(), "<script>") || !strings.Contains(first.Body.String(), "&lt;script&gt;") {
		t.Fatalf("visitor was not safely escaped: %s", first.Body.String())
	}
	if !strings.Contains(first.Body.String(), "2026-08-11T14:30:00Z") || !strings.Contains(first.Body.String(), "#1") {
		t.Fatalf("first response is missing its dynamic proof: %s", first.Body.String())
	}
	if !strings.Contains(second.Body.String(), "2026-08-11T14:30:01Z") || !strings.Contains(second.Body.String(), "#2") {
		t.Fatalf("second response is missing fresh dynamic proof: %s", second.Body.String())
	}
	if first.Body.String() == second.Body.String() {
		t.Fatal("separate requests produced identical bodies")
	}
}

func TestHeadHasGetHeadersAndNoBody(t *testing.T) {
	handler := NewWithClock(func() time.Time {
		return time.Date(2026, time.August, 11, 14, 30, 0, 0, time.UTC)
	})

	response := request(t, handler, http.MethodHead, "/")
	if response.Code != http.StatusOK {
		t.Fatalf("status = %d", response.Code)
	}
	if response.Body.Len() != 0 {
		t.Fatalf("HEAD wrote %d body bytes", response.Body.Len())
	}
	length, err := strconv.Atoi(response.Header().Get("Content-Length"))
	if err != nil || length <= 0 {
		t.Fatalf("Content-Length = %q, err = %v", response.Header().Get("Content-Length"), err)
	}
	if response.Header().Get("Cache-Control") != "no-store" || response.Header().Get("Server-Timing") == "" {
		t.Fatalf("HEAD omitted dynamic response headers: %v", response.Header())
	}
}

func TestRouterOwnsMethodAndPathPolicy(t *testing.T) {
	handler := New()

	if response := request(t, handler, http.MethodPost, "/"); response.Code != http.StatusMethodNotAllowed {
		t.Fatalf("POST / status = %d", response.Code)
	}
	if response := request(t, handler, http.MethodGet, "/missing"); response.Code != http.StatusNotFound {
		t.Fatalf("GET /missing status = %d", response.Code)
	}
}

func request(t *testing.T, handler http.Handler, method, target string) *httptest.ResponseRecorder {
	t.Helper()
	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, httptest.NewRequest(method, target, nil))
	return recorder
}

func TestGetContentLengthMatchesBufferedBody(t *testing.T) {
	handler := New()
	response := request(t, handler, http.MethodGet, "/")
	want, err := strconv.Atoi(response.Header().Get("Content-Length"))
	if err != nil {
		t.Fatal(err)
	}
	body, err := io.ReadAll(response.Result().Body)
	if err != nil {
		t.Fatal(err)
	}
	if len(body) != want {
		t.Fatalf("body length = %d, Content-Length = %d", len(body), want)
	}
}
