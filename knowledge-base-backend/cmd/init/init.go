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
		tlog.F(ctx).Msgf("process initialization failed because configuration key %s is empty", "time location")
	}

	location, err := time.LoadLocation(timeLocation)
	if err != nil {
		tlog.F(ctx).Err(err).Msgf("process initialization failed while loading time location %q", timeLocation)
	}

	time.Local = location
}
