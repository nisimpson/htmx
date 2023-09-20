package main

import (
	"errors"
	"net/http"

	"github.com/nisimpson/htmx"
	"github.com/nisimpson/htmx/examples/snippets"
	"github.com/nisimpson/htmx/examples/snippets/html/pages"
)

func (s *SnippetBox) viewSnippet(w *htmx.ResponseWriter, r *htmx.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "missing id query", http.StatusBadRequest)
		return
	}

	item, err := s.Store.GetSnippetWithID(r.Context(), id)
	if errors.Is(err, snippets.ErrItemNotFound) {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	view := pages.SnippetPage{Snippet: item}
	s.render(w, view)
}
