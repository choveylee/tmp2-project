// Package redmodel configures the Redis client shared by the service.
package redmodel

import (
	"context"
	"fmt"

	"github.com/choveylee/tcfg"
	"github.com/choveylee/tdb"
	"github.com/choveylee/terror"
	"github.com/choveylee/tlog"
	"github.com/redis/go-redis/v9"

	"dev.choveylee.top/knowledge-base-backend/internal/const"
)

var (
	runMode string

	serverClient *tdb.RedisClient
)

// GetCronRedisClient returns the Redis client used for cron coordination.
func GetCronRedisClient() *redis.Client {
	return serverClient.Client()
}

// InitRedisModel initializes the shared Redis client.
func InitRedisModel(ctx context.Context) *terror.Terror {
	runMode = tcfg.DefaultString(tcfg.LocalKey("RUN_MODE"), constant.RunModeDebug)

	serverAddress := tcfg.DefaultString(fmt.Sprintf("%s::%s", runMode, tcfg.LocalKey("SERVER_REDIS_ADDRESS")), "")
	if serverAddress == "" {
		errMsg := tlog.E(ctx).Msg("Init redis err (server redis address empty)")

		errx := terror.NewRawTerror(ctx, terror.ErrConfInvalid("server redis address"), errMsg)

		return errx
	}

	serverPassword := tcfg.DefaultString(fmt.Sprintf("%s::%s", runMode, tcfg.LocalKey("SERVER_REDIS_PASSWORD")), "")

	serverPoolSize := tcfg.DefaultInt(tcfg.LocalKey("SERVER_REDIS_POOLSIZE"), 100)
	if serverPoolSize <= 0 {
		errMsg := tlog.E(ctx).Msgf("Init redis (pool size: %d) err (server redis pool size invalid)",
			serverPoolSize)

		errx := terror.NewRawTerror(ctx, terror.ErrConfInvalid("server redis pool size"), errMsg)

		return errx
	}

	var err error

	serverClient, err = tdb.NewRedisClient(ctx, serverAddress, serverPassword, serverPoolSize)
	if err != nil {
		errMsg := tlog.E(ctx).Err(err).Msgf("Init redis (server address: %s, pool size: %d) err (new redis client %v)",
			serverAddress, serverPoolSize, err)

		errx := terror.NewRawTerror(ctx, err, errMsg)

		return errx
	}

	if serverClient == nil {
		errMsg := tlog.E(ctx).Msgf("Init redis (server address: %s, pool size: %d) err (redis client invalid)",
			serverAddress, serverPoolSize)

		errx := terror.NewRawTerror(ctx, terror.ErrConfInvalid("redis client"), errMsg)

		return errx
	}

	return nil
}
