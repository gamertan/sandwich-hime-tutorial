// SPDX-License-Identifier: 0BSD

// Package server owns the example application's HTTP policy and request data.
package server

import (
	"bytes"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync/atomic"
	"time"
	"unicode/utf8"

	"gamertan.com/sandwich-hime/sando"
	"gitea.speelman.ca/gamertan/sandwich-hime-tutorial/internal/views"
)

const maxVisitorRunes = 80

// New returns the complete example application using the system clock.
func New() http.Handler {
	return NewWithClock(time.Now)
}

// NewWithClock returns the application with an injectable request clock.
// It is exported from this internal package so tests can make time exact.
func NewWithClock(now func() time.Time) http.Handler {
	if now == nil {
		now = time.Now
	}

	application := &application{now: now}
	mux := http.NewServeMux()
	mux.HandleFunc("GET /{$}", application.home)
	return mux
}

type application struct {
	now      func() time.Time
	requests atomic.Uint64
}

func (a *application) home(w http.ResponseWriter, r *http.Request) {
	renderedAt := a.now().UTC()
	requestNumber := a.requests.Add(1)

	body := views.Home(views.HomeView{
		Visitor: normalizeVisitor(r.URL.Query().Get("name")),
		Trails:  trailsForRequest(),
	})
	page := views.Layout(views.LayoutView{
		Title:         "A small Sandwich Hime site",
		Body:          body,
		RenderedAtUTC: renderedAt.Format(time.RFC3339Nano),
		RequestNumber: requestNumber,
	})

	var output bytes.Buffer
	renderStarted := time.Now()
	if err := sando.Render(r.Context(), &output, page); err != nil {
		w.Header().Set("Cache-Control", "no-store")
		http.Error(w, "could not render page", http.StatusInternalServerError)
		return
	}
	renderDuration := time.Since(renderStarted)

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("Content-Length", strconv.Itoa(output.Len()))
	w.Header().Set("Content-Security-Policy", "default-src 'none'; style-src 'unsafe-inline'; base-uri 'none'; form-action 'self'; frame-ancestors 'none'")
	w.Header().Set("Referrer-Policy", "no-referrer")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("Server-Timing", fmt.Sprintf(`sando;dur=%.3f;desc="buffered component render"`, float64(renderDuration)/float64(time.Millisecond)))
	w.WriteHeader(http.StatusOK)
	if r.Method == http.MethodHead {
		return
	}
	_, _ = w.Write(output.Bytes())
}

func trailsForRequest() []views.Trail {
	return []views.Trail{
		{
			Label:       "Read the official tutorial",
			Description: "Walk through the source one typed component at a time.",
			URL:         "https://sandwichhime.com/docs/tutorial/",
		},
		{
			Label:       "Study the language and safety boundary",
			Description: "See where contextual escaping succeeds and where compilation deliberately stops.",
			URL:         "https://sandwichhime.com/docs/",
		},
		{
			Label:       "Render another visitor",
			Description: "This local link sends different untrusted request data through the same compiled templates.",
			URL:         "/?name=friend",
		},
	}
}

func normalizeVisitor(value string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return "traveler"
	}
	if utf8.RuneCountInString(value) <= maxVisitorRunes {
		return value
	}
	runes := []rune(value)
	return string(runes[:maxVisitorRunes])
}
