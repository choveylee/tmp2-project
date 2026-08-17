package lib

import (
	"context"

	"github.com/choveylee/tcfg"
	"github.com/choveylee/terror"
	"github.com/choveylee/tlog"
)

var einoTextChunkSize int

func initEino(ctx context.Context) *terror.Terror {
	einoTextChunkSize = tcfg.DefaultInt(tcfg.LocalKey("EINO_TEXT_CHUNK_SIZE"), 1800)
	if einoTextChunkSize <= 0 {
		errMsg := tlog.E(ctx).Msgf("Init eino (config key: %s, chunk size: %d) err (chunk size invalid)",
			"eino text chunk size",
			einoTextChunkSize)

		errx := terror.NewRawTerror(ctx, terror.ErrConfInvalid("eino text chunk size"), errMsg)

		return errx
	}

	return nil
}
