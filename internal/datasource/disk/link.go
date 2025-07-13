package disk

import (
	"context"
	"fmt"
	"os"

	"github.com/goccy/go-yaml"
	"github.com/orewaee/nanolink/internal/core/domain"
	"github.com/orewaee/nanolink/internal/core/driven/repo"
	"github.com/orewaee/nanolink/internal/model"
)

type YamlLinkRepo struct {
	Dir string
}

func (r *YamlLinkRepo) AddLink(ctx context.Context, link domain.Link) error {
	path := fmt.Sprintf("%s/%s.yaml", r.Dir, link.Id)
	info, err := os.Stat(path)
	if err == nil && info.IsDir() {
		return domain.ErrInvalidLink
	}

	if err == nil {
		return domain.ErrLinkAlreadyExists
	}

	if !os.IsNotExist(err) {
		return err
	}

	data, err := yaml.Marshal(&model.Link{
		Id:       link.Id,
		Location: link.Location,
	})
	if err != nil {
		return err
	}

	file, err := os.Create(path)
	if err != nil {
		return nil
	}

	file.Write(data)
	return file.Close()
}

func (r *YamlLinkRepo) GetLinkById(ctx context.Context, id string) (domain.Link, error) {
	link := domain.Link{Id: id}
	path := fmt.Sprintf("%s/%s.yaml", r.Dir, id)
	info, err := os.Stat(path)
	if err != nil && os.IsNotExist(err) {
		return link, domain.ErrLinkNotFound
	}

	if err != nil {
		return link, err
	}

	if info.IsDir() {
		return link, domain.ErrInvalidLink
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return link, err
	}

	if err := yaml.Unmarshal(data, &link); err != nil {
		return link, err
	}

	return link, nil
}

func (r *YamlLinkRepo) RemoveLinkById(ctx context.Context, id string) error {
	path := fmt.Sprintf("%s/%s.yaml", r.Dir, id)
	info, err := os.Stat(path)
	if err != nil && os.IsNotExist(err) {
		return domain.ErrLinkNotFound
	}

	if err != nil {
		return err
	}

	if info.IsDir() {
		return domain.ErrInvalidLink
	}

	return os.Remove(path)
}

func NewLinkRepo() repo.LinkRepo {
	return &YamlLinkRepo{}
}
