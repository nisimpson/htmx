package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/nisimpson/htmx"
	"github.com/nisimpson/htmx/examples/snippets/html/components"
	"github.com/nisimpson/htmx/examples/snippets/pkg/models"
)

func (s *SnippetBox) snippets(w *htmx.ResponseWriter, r *htmx.Request) {
	if r.Method == http.MethodGet {
		s.pollSnippets(w, r)
		return
	} else if r.Method == http.MethodPost {
		s.createSnippet(w, r)
		return
	}
	w.Header().Set("Allow", strings.Join([]string{http.MethodGet, http.MethodPost}, ","))
	http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
}

func (s *SnippetBox) pollSnippets(w *htmx.ResponseWriter, r *htmx.Request) {
	snippets, err := s.SnippetModel.FetchAll(r.Context())
	if err != nil {
		s.serverError(w, err)
		return
	}

	s.render(w, r, components.SnippetsList{Snippets: snippets})
}

func (s *SnippetBox) createSnippet(w *htmx.ResponseWriter, r *htmx.Request) {
	id, err := s.Store.CreateSnippet(r.Context(), &models.Snippet{
		Title:   "O snail",
		Content: "O snail\nClimb Mount Fuji,\nBut slowly, slow-ly!\n\n- Kobayashi Issa",
	})

	if err != nil {
		s.serverError(w, err)
		return
	}

	if r.IsHTMXRequest() {
		// The snippets polling mechanism will fetch the new snippet, so just return 201
		w.WriteHeader(http.StatusCreated)
		return
	}

	// redirect to the newly created snippet.
	http.Redirect(w, r.Request, fmt.Sprintf("/snippet?id=%s", id), http.StatusSeeOther)
}
