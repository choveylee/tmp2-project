// Package service implements service-layer logic exposed through HTTP handlers and scheduled jobs.
package service

import (
	"context"
	"strings"
	"time"

	"github.com/choveylee/tcfg"
	"github.com/choveylee/terror"
	"github.com/choveylee/tlog"
)

var (
	ingestWorkerEnabled      bool
	ingestJobWorkerName      string
	ingestWorkerBlockTimeout time.Duration
	ingestWorkerErrorSleep   time.Duration
)

// InitService initializes service-layer dependencies.
func InitService(ctx context.Context) *terror.Terror {
	ingestWorkerEnabled = tcfg.DefaultBool(tcfg.LocalKey("INGEST_JOB_WORKER_ENABLED"), true)

	ingestJobWorkerName = strings.TrimSpace(tcfg.DefaultString(tcfg.LocalKey("INGEST_JOB_WORKER_NAME"), "knowledge-base-ingest-worker"))
	if ingestJobWorkerName == "" {
		errMsg := tlog.E(ctx).Msg("Init service err (ingest job worker name empty)")

		errx := terror.NewRawTerror(ctx, terror.ErrConfInvalid("ingest job worker name"), errMsg)

		return errx
	}

	blockSeconds := tcfg.DefaultInt(tcfg.LocalKey("INGEST_JOB_WORKER_BLOCK_SECONDS"), 5)
	if blockSeconds <= 0 {
		errMsg := tlog.E(ctx).Msgf("Init service (block seconds: %d) err (ingest job worker block seconds invalid)",
			blockSeconds)

		errx := terror.NewRawTerror(ctx, terror.ErrConfInvalid("ingest job worker block seconds"), errMsg)

		return errx
	}

	errorSeconds := tcfg.DefaultInt(tcfg.LocalKey("INGEST_JOB_WORKER_ERROR_SECONDS"), 3)
	if errorSeconds <= 0 {
		errMsg := tlog.E(ctx).Msgf("Init service (error seconds: %d) err (ingest job worker error seconds invalid)",
			errorSeconds)

		errx := terror.NewRawTerror(ctx, terror.ErrConfInvalid("ingest job worker error seconds"), errMsg)

		return errx
	}

	ingestWorkerBlockTimeout = time.Duration(blockSeconds) * time.Second
	ingestWorkerErrorSleep = time.Duration(errorSeconds) * time.Second

	return nil
}
