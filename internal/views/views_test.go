// SPDX-License-Identifier: 0BSD

package views

import (
	"context"
	"strings"
	"testing"

	"gamertan.com/sandwich-hime/sando"
)

func TestHomeEscapesUntrustedTextAndNestsBadge(t *testing.T) {
	t.Parallel()

	const attack = `<script>alert("no")</script>`
	var output strings.Builder
	page := Home(HomeView{Visitor: attack})
	if err := sando.Render(context.Background(), &output, page); err != nil {
		t.Fatal(err)
	}

	rendered := output.String()
	if strings.Contains(rendered, attack) || strings.Contains(rendered, "<script>") {
		t.Fatalf("visitor became markup: %s", rendered)
	}
	if !strings.Contains(rendered, "&lt;script&gt;") {
		t.Fatalf("escaped visitor is missing: %s", rendered)
	}
	if !strings.Contains(rendered, `<span class="badge">rendered on this request</span>`) {
		t.Fatalf("nested Badge component is missing: %s", rendered)
	}
}

func TestHomeRejectsDangerousURL(t *testing.T) {
	t.Parallel()

	var output strings.Builder
	page := Home(HomeView{Trails: []Trail{
		{Label: "unsafe", URL: "javascript:alert(1)"},
	}})
	if err := sando.Render(context.Background(), &output, page); err == nil {
		t.Fatal("dangerous URL rendered without an error")
	}
}
