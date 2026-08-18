// Package init performs process-level initialization required before the service starts.
package init

import (
	"context"
	"strings"
	"time"

	"github.com/choveylee/tcfg"
	"github.com/choveylee/tlog"
)

func init() {
	ctx := context.Background()

	timeLocation := strings.TrimSpace(tcfg.DefaultString(tcfg.LocalKey("TIME_LOCATION"), "Asia/Shanghai"))
	if timeLocation == "" {
		tlog.F(ctx).Msg("Process init failed (time location empty)")

		return
	}

	location, err := time.LoadLocation(timeLocation)
	if err != nil {
		tlog.F(ctx).Err(err).Msgf("Process init failed (time location: %s) err (load location %v)",
			timeLocation, err)

		return
	}

	time.Local = location
}
