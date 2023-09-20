package components

import (
	"sort"
	"strconv"

	"github.com/nisimpson/htmx/examples/snippets/pkg/models"
)

type SnippetsList struct {
	Snippets []*models.Snippet
}

func (SnippetsList) TemplateFiles() []string {
	return []string{
		"./html/components/snippets_list.tmpl",
	}
}

func (SnippetsList) TemplateName() string {
	return "snippets_list"
}

func (s SnippetsList) TemplateData() any {
	// sort the snippets list
	sort.Slice(s.Snippets, func(i, j int) bool {
		idi, _ := strconv.Atoi(s.Snippets[i].ID)
		idj, _ := strconv.Atoi(s.Snippets[j].ID)
		return idi < idj
	})
	return s
}
