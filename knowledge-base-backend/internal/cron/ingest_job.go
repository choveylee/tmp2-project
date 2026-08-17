package crontab

import (
	"context"

	"dev.choveylee.top/knowledge-base-backend/internal/service"
)

func runIngestJob(ctx context.Context, params ...any) error {
	errx := service.RunPendingIngestJob(ctx, ingestJobWorkerName, ingestJobBatchLimit)
	if errx != nil {
		return errx
	}

	return nil
}
