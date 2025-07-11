package driving

import (
	"context"

	"github.com/orewaee/nanolink/internal/core/domain"
)

type LinkApi interface {
	GetLinkById(ctx context.Context, id string) (domain.Link, error)

	AddLink(ctx context.Context, id, location string) (domain.Link, error)

	RemoveLinkById(ctx context.Context, id string) error
}
