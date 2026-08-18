package crontab

import (
	"context"

	"github.com/choveylee/tlog"

	"dev.choveylee.top/knowledge-base-backend/internal/service"
)

func runIngestJob(ctx context.Context, params ...any) error {
	errx := service.RunPendingIngestJob(ctx, ingestJobWorkerName, ingestJobBatchLimit)
	if errx != nil {
		errMsg := tlog.E(ctx).Err(errx).Msgf("Run ingest job (worker name: %s, batch limit: %d) err (run pending ingest job %v)",
			ingestJobWorkerName, ingestJobBatchLimit, errx)
		errx.AttachErrMsg(errMsg)

		return errx
	}

	return nil
}
