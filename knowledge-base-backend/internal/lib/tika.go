package lib

import (
	"context"
	"strings"
	"time"

	"github.com/choveylee/tcfg"
	"github.com/choveylee/terror"
	"github.com/choveylee/thttp"
	"github.com/choveylee/tlog"
)

var (
	tikaEndpoint string

	tikaClient *thttp.HttpClient
)

func initTika(ctx context.Context) *terror.Terror {
	tikaEndpoint = strings.TrimRight(strings.TrimSpace(tcfg.DefaultString(tcfg.LocalKey("TIKA_ENDPOINT"), "http://127.0.0.1:9998")), "/")
	if tikaEndpoint == "" {
		errMsg := tlog.E(ctx).Msg("Init tika err (tika endpoint empty)")

		errx := terror.NewRawTerror(ctx, terror.ErrConfInvalid("tika endpoint"), errMsg)

		return errx
	}

	tikaClient = thttp.NewHttpClient().WithTimeout(2 * time.Minute)
	if tikaClient == nil {
		errMsg := tlog.E(ctx).Msg("Init tika err (tika client nil)")

		errx := terror.NewRawTerror(ctx, terror.ErrConfInvalid("tika client"), errMsg)

		return errx
	}

	return nil
}
