package usecase

import (
	"context"

	"github.com/orewaee/nanolink/internal/core/domain"
	"github.com/orewaee/nanolink/internal/core/driven/repo"
	"github.com/orewaee/nanolink/internal/core/driving"
)

type LinkService struct {
	linkRepo repo.LinkRepo
}

func (s *LinkService) AddLink(ctx context.Context, id string, location string) (domain.Link, error) {
	return s.linkRepo.GetLinkById(ctx, id)
}

func (s *LinkService) GetLinkById(ctx context.Context, id string) (domain.Link, error) {
	return s.linkRepo.GetLinkById(ctx, id)
}

func (s *LinkService) RemoveLinkById(ctx context.Context, id string) error {
	return s.linkRepo.RemoveLinkById(ctx, id)
}

func NewLinkService(linkRepo repo.LinkRepo) driving.LinkApi {
	return &LinkService{linkRepo}
}
