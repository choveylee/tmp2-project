// Package dbmodel configures the MySQL client shared by the service.
package dbmodel

import (
	"context"
	"fmt"

	"github.com/choveylee/tcfg"
	"github.com/choveylee/tdb"
	"github.com/choveylee/terror"
	"github.com/choveylee/tlog"
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"gorm.io/gorm"

	"dev.choveylee.top/knowledge-base-backend/internal/const"
)

var (
	runMode string

	serverClient *tdb.MysqlClient
)

// InitMysqlModel initializes the shared MySQL client.
func InitMysqlModel(ctx context.Context) *terror.Terror {
	runMode = tcfg.DefaultString(tcfg.LocalKey("RUN_MODE"), constant.RunModeDebug)

	serverDsn := tcfg.DefaultString(fmt.Sprintf("%s::%s", runMode, tcfg.LocalKey("SERVER_MYSQL_DSN")), "")
	if serverDsn == "" {
		errMsg := tlog.E(ctx).Msg("Init mysql err (server mysql dsn empty)")

		errx := terror.NewRawTerror(ctx, terror.ErrConfInvalid("server mysql dsn"), errMsg)

		return errx
	}

	var err error

	if runMode == constant.RunModeDebug {
		serverClient, err = tdb.NewMysqlClientWithLog(ctx, serverDsn)
	} else {
		serverClient, err = tdb.NewMysqlClient(ctx, serverDsn)
	}
	if err != nil {
		errMsg := tlog.E(ctx).Err(err).Msg("Init mysql err (new mysql client)")

		errx := terror.NewRawTerror(ctx, err, errMsg)

		return errx
	}
	if serverClient == nil {
		errMsg := tlog.E(ctx).Msg("Init mysql err (mysql client invalid)")

		errx := terror.NewRawTerror(ctx, terror.ErrConfInvalid("mysql client"), errMsg)

		return errx
	}

	maxIdleConns := tcfg.DefaultInt(tcfg.LocalKey("MYSQL_MAX_IDLE_CONNS"), 10)

	err = serverClient.SetMaxIdleConns(maxIdleConns)
	if err != nil {
		errMsg := tlog.E(ctx).Err(err).Msgf("Init mysql (max idle conns: %d) err (set max idle conns %v)",
			maxIdleConns, err)

		errx := terror.NewRawTerror(ctx, err, errMsg)

		return errx
	}

	maxOpenConns := tcfg.DefaultInt(tcfg.LocalKey("MYSQL_MAX_OPEN_CONNS"), 100)

	err = serverClient.SetMaxOpenConns(maxOpenConns)
	if err != nil {
		errMsg := tlog.E(ctx).Err(err).Msgf("Init mysql (max open conns: %d) err (set max open conns %v)",
			maxOpenConns, err)

		errx := terror.NewRawTerror(ctx, err, errMsg)

		return errx
	}

	return nil
}

// Tx returns the request-scoped GORM database handle.
func Tx(ctx context.Context) *gorm.DB {
	return serverClient.Tx(ctx, runMode)
}

// DB returns the request-scoped GORM database handle.
func DB(ctx context.Context) *gorm.DB {
	return serverClient.DB(ctx, runMode)
}
