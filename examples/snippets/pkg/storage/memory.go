package storage

import (
	"context"
	"strconv"
	"time"

	"github.com/nisimpson/htmx/examples/snippets"
	"github.com/nisimpson/htmx/examples/snippets/pkg/models"
)

type MemoryStorage struct {
	counter  int
	snippets map[string]*models.Snippet
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{snippets: make(map[string]*models.Snippet)}
}

func (m *MemoryStorage) CreateSnippet(ctx context.Context, data *models.Snippet) (string, error) {
	m.counter++
	data.ID = strconv.Itoa(m.counter)
	data.Created = time.Now()
	data.Expires = data.Created.Add(24 * time.Hour)
	m.snippets[data.ID] = data
	return data.ID, nil
}

func (m MemoryStorage) GetSnippetWithID(ctx context.Context, id string) (*models.Snippet, error) {
	if item, ok := m.snippets[id]; !ok {
		return nil, snippets.ErrItemNotFound
	} else {
		return item, nil
	}
}

func (m MemoryStorage) GetSnippets(ctx context.Context) ([]*models.Snippet, error) {
	snippets := make([]*models.Snippet, 0, len(m.snippets))
	for _, value := range m.snippets {
		snippets = append(snippets, value)
	}
	return snippets, nil
}
