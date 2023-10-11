package app

import (
	"context"

	"github.com/korovindenis/go-market/internal/port/webserver"
)

func Run(ctx context.Context, cfg any) error {
	return webserver.Run(ctx, cfg)
}
