// SPDX-License-Identifier: 0BSD

package views

import "gamertan.com/sandwich-hime/sando"

// Trail is one server-provided destination rendered by Home.
type Trail struct {
	Label       string
	Description string
	URL         string
}

// HomeView is the complete typed input to the inner page component.
type HomeView struct {
	Visitor string
	Trails  []Trail
}

// LayoutView is the typed input to the full-document component.
type LayoutView struct {
	Title         string
	Body          sando.Component
	RenderedAtUTC string
	RequestNumber uint64
}
