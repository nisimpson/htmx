package models

import (
	"context"
	"time"
)

type Snippet struct {
	ID      string
	Title   string
	Content string
	Created time.Time
	Expires time.Time
}

type SnippetStore interface {
	CreateSnippet(ctx context.Context, data *Snippet) (string, error)
	GetSnippetWithID(ctx context.Context, id string) (*Snippet, error)
	GetSnippets(ctx context.Context) ([]*Snippet, error)
}

type SnippetModel struct {
	Store SnippetStore
}

func (s SnippetModel) Create(ctx context.Context, data *Snippet) (string, error) {
	return s.Store.CreateSnippet(ctx, data)
}

func (s SnippetModel) Fetch(ctx context.Context, id string) (*Snippet, error) {
	return s.Store.GetSnippetWithID(ctx, id)
}

func (s SnippetModel) FetchAll(ctx context.Context) ([]*Snippet, error) {
	return s.Store.GetSnippets(ctx)
}
