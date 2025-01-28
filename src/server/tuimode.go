package main

import (
	"context"

	zl "github.com/rs/zerolog"

	"server/internal/cache"
	"server/internal/types"
)

func TuiMode(ctx context.Context, logger *zl.Logger, opts *cache.Cache[types.Options]) error {
	logger.Info().Str("Mode", "Interactive").Msg("Starting Server")
	return nil
}
