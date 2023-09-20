package pages

import "github.com/nisimpson/htmx/examples/snippets/html/components"

type HomePage struct {
	components.SnippetsList
}

func (HomePage) TemplateName() string { return "home.tmpl" }
