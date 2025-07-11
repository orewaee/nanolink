package delivery

import "context"

type Controller interface {
	Run() error
	Shutdown(ctx context.Context) error
}
