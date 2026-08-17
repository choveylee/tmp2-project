// Package crontab loads configuration for scheduled jobs and registers them with the cron runtime.
package crontab

import (
	"context"
	"strings"
	"time"

	"github.com/choveylee/tcfg"
	"github.com/choveylee/tcron"
	"github.com/choveylee/terror"
	"github.com/choveylee/tlog"

	redmodel "dev.choveylee.top/knowledge-base-backend/internal/model/redis"
)

var (
	testSyncCron string

	ingestJobCron       string
	ingestJobWorkerName string
	ingestJobBatchLimit int
)

// InitCron loads cron configuration.
func InitCron(ctx context.Context) *terror.Terror {
	testSyncCron = strings.TrimSpace(tcfg.DefaultString(tcfg.LocalKey("TEST_SYNC_CRON"), ""))

	ingestJobCron = strings.TrimSpace(tcfg.DefaultString(tcfg.LocalKey("INGEST_JOB_CRON"), "0 */5 * * * *"))

	ingestJobWorkerName = strings.TrimSpace(tcfg.DefaultString(tcfg.LocalKey("INGEST_JOB_WORKER_NAME"), "knowledge-base-ingest-worker"))
	if ingestJobWorkerName == "" {
		errMsg := tlog.E(ctx).Msgf("Init cron (config key: %s) err (ingest job worker name empty)",
			"ingest job worker name")

		errx := terror.NewRawTerror(ctx, terror.ErrConfInvalid("ingest job worker name"), errMsg)

		return errx
	}

	ingestJobBatchLimit = tcfg.DefaultInt(tcfg.LocalKey("INGEST_JOB_BATCH_LIMIT"), 5)
	if ingestJobBatchLimit <= 0 {
		errMsg := tlog.E(ctx).Msgf("Init cron (config key: %s, ingest job batch limit: %d) err (ingest job batch limit invalid)",
			"ingest job batch limit", ingestJobBatchLimit)

		errx := terror.NewRawTerror(ctx, terror.ErrConfInvalid("ingest job batch limit"), errMsg)

		return errx
	}

	return nil
}

// StartCron registers the configured cron jobs.
func StartCron(ctx context.Context) *terror.Terror {
	cronRedisClient := redmodel.GetCronRedisClient()

	if testSyncCron != "" {
		_, err := tcron.RegisterSingletonCron(testSyncCron, runTestSync, cronRedisClient, 10*time.Minute)
		if err != nil {
			errMsg := tlog.E(ctx).Err(err).Msgf("cron job registration failed for schedule %q",
				testSyncCron)

			errx := terror.NewRawTerror(ctx, err, errMsg)

			return errx
		}
	}

	if ingestJobCron != "" {
		_, err := tcron.RegisterSingletonCron(ingestJobCron, runIngestJob, cronRedisClient, 30*time.Minute)
		if err != nil {
			errMsg := tlog.E(ctx).Err(err).Msgf("cron job registration failed for ingest job schedule %q",
				ingestJobCron)

			errx := terror.NewRawTerror(ctx, err, errMsg)

			return errx
		}
	}

	return nil
}
