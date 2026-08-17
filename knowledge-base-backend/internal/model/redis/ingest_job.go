package redmodel

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/choveylee/terror"
	"github.com/choveylee/tlog"
	"github.com/redis/go-redis/v9"

	constant "dev.choveylee.top/knowledge-base-backend/internal/const"
)

func PushParseIngestJob(ctx context.Context, jobId string) *terror.Terror {
	jobId = strings.TrimSpace(jobId)
	if jobId == "" {
		errMsg := tlog.E(ctx).Msg("Push parse ingest job err (job id empty)")

		errx := terror.NewTerror(ctx, terror.ErrParamInvalid("job_id"), constant.ErrorCodeRequestParamInvalid, errMsg)

		return errx
	}

	if serverClient == nil {
		errMsg := tlog.E(ctx).Msgf("Push parse ingest job (job id: %s) err (redis client nil)",
			jobId)

		errx := terror.NewTerror(ctx, terror.ErrConfInvalid("redis client"), constant.ErrorCodeRedisServerAbnormal, errMsg)

		return errx
	}

	// Redis queue wakes workers quickly; MySQL ingest_jobs remains the source of truth.
	err := serverClient.Client().RPush(ctx, constant.RedisParseIngestJobQueueKey, jobId).Err()
	if err != nil {
		errMsg := tlog.E(ctx).Err(err).Msgf("Push parse ingest job (job id: %s) err (redis rpush %v)",
			jobId, err)

		errx := terror.NewTerror(ctx, err, constant.ErrorCodeRedisServerAbnormal, errMsg)

		return errx
	}

	return nil
}

func BlockPopParseIngestJobId(ctx context.Context, timeout time.Duration) (string, *terror.Terror) {
	if timeout <= 0 {
		timeout = 5 * time.Second
	}

	if serverClient == nil {
		errMsg := tlog.E(ctx).Msgf("Block pop parse ingest job id (timeout: %s) err (redis client nil)",
			timeout)

		errx := terror.NewTerror(ctx, terror.ErrConfInvalid("redis client"), constant.ErrorCodeRedisServerAbnormal, errMsg)

		return "", errx
	}

	jobItems, err := serverClient.Client().BLPop(ctx, timeout, constant.RedisParseIngestJobQueueKey).Result()
	if errors.Is(err, redis.Nil) || errors.Is(err, context.Canceled) {
		return "", nil
	}
	if err != nil {
		errMsg := tlog.E(ctx).Err(err).Msgf("Block pop parse ingest job id (timeout: %s) err (redis blpop %v)",
			timeout, err)

		errx := terror.NewTerror(ctx, err, constant.ErrorCodeRedisServerAbnormal, errMsg)

		return "", errx
	}

	if len(jobItems) < 2 {
		errMsg := tlog.E(ctx).Msgf("Block pop parse ingest job id (timeout: %s, item count: %d) err (redis blpop result invalid)",
			timeout, len(jobItems))

		errx := terror.NewTerror(ctx, terror.ErrParamInvalid("redis blpop result"), constant.ErrorCodeRedisServerAbnormal, errMsg)

		return "", errx
	}

	jobId := strings.TrimSpace(jobItems[1])
	if jobId == "" {
		return "", nil
	}

	return jobId, nil
}
