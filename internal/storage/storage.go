package storage

import (
	"context"
	"time"
)

type Storage interface {
	Save(ctx context.Context, p *Page) error
	Pick(ctx context.Context, userName string) (*Page, error)
	Remove(ctx context.Context, p *Page) error
	PickAll(ctx context.Context, userName string) ([]Page, error)
	SelectTags(ctx context.Context, userName string) ([]string, error)
	SelectByTag(ctx context.Context, tag string, userName string) ([]string, error)
	BatchUpdate(ctx context.Context, pages []Page) error
}

type Page struct {
	URL      string
	Tags     string
	UserName string
	Created  time.Time
}
