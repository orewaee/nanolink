package repo

import (
	"context"

	"github.com/orewaee/nanolink/internal/core/domain"
)

type LinkRepo interface {
	LinkReader
	LinkWriter
}

type LinkReader interface {
	// May return domain.ErrLinkNotFound
	GetLinkById(ctx context.Context, id string) (domain.Link, error)
}

type LinkWriter interface {
	// May return domain.ErrLinkAlreadyExists
	AddLink(ctx context.Context, link domain.Link) error

	// May return domain.ErrLinkNotFound
	RemoveLinkById(ctx context.Context, id string) error
}
