// Package model initializes storage clients and repository-layer dependencies.
package model

import (
	"context"

	"github.com/choveylee/terror"
	"github.com/choveylee/tlog"

	"dev.choveylee.top/knowledge-base-backend/internal/model/mysql"
	"dev.choveylee.top/knowledge-base-backend/internal/model/redis"
)

// InitModel initializes all persistence-layer dependencies.
func InitModel(ctx context.Context) *terror.Terror {
	errx := dbmodel.InitMysqlModel(ctx)
	if errx != nil {
		errMsg := tlog.E(ctx).Err(errx).Msgf("Init model err (init mysql dependencies %v)",
			errx)
		errx.AttachErrMsg(errMsg)

		return errx
	}

	errx = redmodel.InitRedisModel(ctx)
	if errx != nil {
		errMsg := tlog.E(ctx).Err(errx).Msgf("Init model err (init redis dependencies %v)",
			errx)
		errx.AttachErrMsg(errMsg)

		return errx
	}

	return nil
}
