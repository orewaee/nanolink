package driven

import "context"

type IdProvider interface {
	GenerateId(ctx context.Context, len int) (string, error)
}
