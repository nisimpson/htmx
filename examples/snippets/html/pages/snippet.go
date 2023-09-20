package pages

import "github.com/nisimpson/htmx/examples/snippets/pkg/models"

type SnippetPage struct {
	*models.Snippet
}

func (SnippetPage) TemplateName() string { return "snippet.tmpl" }
