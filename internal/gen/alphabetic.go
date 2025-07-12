package gen

import (
	"context"
	"fmt"

	"github.com/orewaee/nanolink/internal/core/driven"
)

type AlphabeticIdProvider struct{}

func NewAlphabeticIdProvider() driven.IdProvider {
	return &AlphabeticIdProvider{}
}

const (
	a = 97
	z = 122
)

func (provider *AlphabeticIdProvider) GenerateId(ctx context.Context, len int) (string, error) {
	if err := ctx.Err(); err != nil {
		return "", fmt.Errorf("operation canceled: %w", err)
	}

	id := make([]byte, len)
	for i := 0; i < len; i++ {
		select {
		case <-ctx.Done():
			return "", ctx.Err()
		default:
			id[i] = byte(minMaxIntN(a, z))
		}
	}

	return string(id), nil
}
