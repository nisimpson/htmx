package main

import (
	"net/http"

	"github.com/nisimpson/htmx"
	"github.com/nisimpson/htmx/examples/snippets/html/pages"
)

func (s *SnippetBox) home(w *htmx.ResponseWriter, r *htmx.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r.Request)
		return
	}

	snippets, err := s.FetchAll(r.Context())
	if err != nil {
		s.serverError(w, err)
		return
	}

	page := pages.HomePage{}
	page.Snippets = snippets
	s.render(w, r, page)
}
