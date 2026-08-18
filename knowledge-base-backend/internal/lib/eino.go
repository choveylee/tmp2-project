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
		errMsg := tlog.E(ctx).Msgf("Init eino (chunk size: %d) err (eino text chunk size invalid)",
			einoTextChunkSize)

		errx := terror.NewRawTerror(ctx, terror.ErrConfInvalid("eino text chunk size"), errMsg)

		return errx
	}

	return nil
}
