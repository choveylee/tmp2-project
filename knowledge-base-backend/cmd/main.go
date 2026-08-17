// Command knowledge-base-backend starts the generated HTTP service.
package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/choveylee/tcfg"
	"github.com/choveylee/terror"
	"github.com/choveylee/thttp"
	"github.com/choveylee/tlog"
	"github.com/choveylee/tserver"
	"github.com/choveylee/tutil"
	"github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	_ "dev.choveylee.top/knowledge-base-backend/cmd/init"

	"dev.choveylee.top/knowledge-base-backend/internal/const"
	"dev.choveylee.top/knowledge-base-backend/internal/cron"
	"dev.choveylee.top/knowledge-base-backend/internal/lib"
	"dev.choveylee.top/knowledge-base-backend/internal/model"
	"dev.choveylee.top/knowledge-base-backend/internal/router"
	"dev.choveylee.top/knowledge-base-backend/internal/service"
)

func main() {
	ctx := context.Background()

	// Run migrations first so the database schema is ready before startup.
	errx := runMigrate(ctx)
	if errx != nil {
		tlog.E(ctx).Err(errx).Msgf("Main err (run migrate %v)",
			errx)

		return
	}

	// Initialize external dependency clients such as object storage.
	errx = lib.InitLib(ctx)
	if errx != nil {
		tlog.E(ctx).Err(errx).Msgf("Main err (init lib %v)",
			errx)

		return
	}

	// Initialize model-layer dependencies such as database and cache clients.
	errx = model.InitModel(ctx)
	if errx != nil {
		tlog.E(ctx).Err(errx).Msgf("Main err (init model %v)",
			errx)

		return
	}

	// Initialize the business service layer.
	errx = service.InitService(ctx)
	if errx != nil {
		tlog.E(ctx).Err(errx).Msgf("Main err (init service %v)",
			errx)

		return
	}

	// Start the resident ingest worker before the compensation cron is registered.
	errx = service.StartIngestWorker(ctx)
	if errx != nil {
		tlog.E(ctx).Err(errx).Msgf("Main err (start ingest worker %v)",
			errx)

		return
	}

	// Initialize and start low-frequency compensation cron jobs.
	errx = crontab.InitCron(ctx)
	if errx != nil {
		tlog.E(ctx).Err(errx).Msgf("Main err (init cron %v)",
			errx)

		return
	}

	errx = crontab.StartCron(ctx)
	if errx != nil {
		tlog.E(ctx).Err(errx).Msgf("Main err (start cron %v)",
			errx)

		return
	}

	httpPort := tcfg.DefaultInt(tcfg.LocalKey("HTTP_PORT"), 8080)

	go func() {
		if err := waitForTcpDial(ctx, httpPort, 30*time.Second); err != nil {
			tlog.W(ctx).Err(err).Msgf("Main (http port: %d) err (wait for tcp dial %v)",
				httpPort, err)

			return
		}

		errx := pingServer(ctx, httpPort)
		if errx != nil {
			tlog.W(ctx).Err(errx).Msgf("Main (http port: %d) err (ping server %v)",
				httpPort, errx)
		} else {
			tlog.I(ctx).Msgf("Main (http port: %d) info (http server started)",
				httpPort)
		}
	}()

	if err := tserver.StartHttpServer(ctx, router.NewRouter(ctx), httpPort); err != nil {
		tlog.F(ctx).Err(err).Msgf("Main (http port: %d) err (start http server %v)",
			httpPort, err)
	}
}

func runMigrate(ctx context.Context) *terror.Terror {
	runMode := tcfg.DefaultString(tcfg.LocalKey("RUN_MODE"), constant.RunModeDebug)

	serverDsn := strings.TrimSpace(tcfg.DefaultString(fmt.Sprintf("%s::%s", runMode, tcfg.LocalKey("SERVER_MYSQL_DSN")), ""))
	if serverDsn == "" {
		errMsg := tlog.E(ctx).Msgf("Run migrate (run mode: %s, config key: %s) err (server mysql dsn empty)",
			runMode, "server mysql dsn")

		errx := terror.NewRawTerror(ctx, terror.ErrConfInvalid("server mysql dsn"), errMsg)

		return errx
	}

	serverDBName := ""
	parsedCfg, err := mysql.ParseDSN(serverDsn)
	if err == nil {
		serverDBName = parsedCfg.DBName
	}

	client, err := migrate.New("file://migration", "mysql://"+tutil.MysqlDsnEncode(serverDsn))
	if err != nil {
		initialClientErr := err

		serverCfg, err := mysql.ParseDSN(serverDsn)
		if err != nil {
			errMsg := tlog.E(ctx).Err(err).Msgf("Run migrate (run mode: %s, migration source: %s, initial client err: %v) err (parse dsn %v)",
				runMode, "file://migration", initialClientErr, err)

			errx := terror.NewRawTerror(ctx, err, errMsg)

			return errx
		}

		dbName := serverCfg.DBName

		tlog.W(ctx).Err(initialClientErr).Msgf("Run migrate (run mode: %s, database: %s, migration source: %s) info (migrate new failed, try db create database %v)",
			runMode, dbName, "file://migration", initialClientErr)

		serverCfg.DBName = ""
		tmpDsn := serverCfg.FormatDSN()

		db, err := sql.Open("mysql", tmpDsn)
		if err != nil {
			errMsg := tlog.E(ctx).Err(err).Msgf("Run migrate (run mode: %s, database: %s, migration source: %s, initial client err: %v) err (db open mysql %v)",
				runMode, dbName, "file://migration", initialClientErr, err)

			errx := terror.NewRawTerror(ctx, err, errMsg)

			return errx
		}

		defer db.Close()

		_, err = db.Exec("CREATE DATABASE IF NOT EXISTS " + "`" + dbName + "` DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci")
		if err != nil {
			errMsg := tlog.E(ctx).Err(err).Msgf("Run migrate (run mode: %s, database: %s, migration source: %s, initial client err: %v) err (db exec create database %v)",
				runMode, dbName, "file://migration", initialClientErr, err)

			errx := terror.NewRawTerror(ctx, err, errMsg)

			return errx
		}

		client, err = migrate.New("file://migration", "mysql://"+tutil.MysqlDsnEncode(serverDsn))
		if err != nil {
			errMsg := tlog.E(ctx).Err(err).Msgf("Run migrate (run mode: %s, database: %s, migration source: %s, initial client err: %v) err (db migrate new %v)",
				runMode, dbName, "file://migration", initialClientErr, err)

			errx := terror.NewRawTerror(ctx, err, errMsg)

			return errx
		}
	}

	defer func() {
		srcErr, dbErr := client.Close()
		if srcErr != nil || dbErr != nil {
			tlog.W(ctx).Msgf("Run migrate (run mode: %s, database: %s) err (db migrate close source %v database %v)",
				runMode, serverDBName, srcErr, dbErr)
		}
	}()

	err = client.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		if serverDBName != "" {
			errMsg := tlog.E(ctx).Err(err).Msgf("Run migrate (run mode: %s, database: %s, migration source: %s) err (db migrate up %v)",
				runMode, serverDBName, "file://migration", err)

			errx := terror.NewRawTerror(ctx, err, errMsg)

			return errx
		}

		errMsg := tlog.E(ctx).Err(err).Msgf("Run migrate (run mode: %s, migration source: %s) err (db migrate up %v)",
			runMode, "file://migration", err)

		errx := terror.NewRawTerror(ctx, err, errMsg)

		return errx
	}

	return nil
}

func waitForTcpDial(ctx context.Context, port int, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)

	address := fmt.Sprintf("127.0.0.1:%d", port)

	dialer := net.Dialer{Timeout: 150 * time.Millisecond}

	for time.Now().Before(deadline) {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		conn, err := dialer.DialContext(ctx, "tcp", address)
		if err == nil {
			_ = conn.Close()

			return nil
		}

		time.Sleep(50 * time.Millisecond)
	}

	return fmt.Errorf("timed out waiting for the TCP listener at %s", address)
}

// Keep the local health probe URL aligned with HTTP_PORT.
func resolvePingBaseUrl(httpPort int) string {
	serverPingHost := strings.TrimSpace(tcfg.DefaultString(tcfg.LocalKey("SERVER_PING_HOST"), ""))
	if serverPingHost == "" {
		return fmt.Sprintf("http://127.0.0.1:%d", httpPort)
	}

	parsedUrl, err := url.Parse(serverPingHost)
	if err != nil {
		return fmt.Sprintf("http://127.0.0.1:%d", httpPort)
	}

	if parsedUrl.Scheme == "" {
		parsedUrl.Scheme = "http"
	}

	switch parsedUrl.Hostname() {
	case "127.0.0.1", "localhost", "::1":
		parsedUrl.Host = net.JoinHostPort(parsedUrl.Hostname(), strconv.Itoa(httpPort))
	}

	return strings.TrimSuffix(parsedUrl.String(), "/")
}

func pingServer(ctx context.Context, httpPort int) *terror.Terror {
	pingCount := tcfg.DefaultInt(tcfg.LocalKey("SERVER_PING_COUNT"), 3)

	baseUrl := resolvePingBaseUrl(httpPort)
	pingUrl := baseUrl + "/healthz"

	for i := range pingCount {
		if i > 0 {
			time.Sleep(time.Second)
		}

		response, err := thttp.Get(ctx, pingUrl, nil, nil)
		if err != nil {
			continue
		}

		statusCode, _, err := response.ToString()
		if err != nil || statusCode != http.StatusOK {
			continue
		}

		return nil
	}

	err := terror.ErrSvcExecute("http server")

	errMsg := tlog.E(ctx).Err(err).Msgf("Ping server (http port: %d, ping count: %d, ping url: %s) err (get %v)",
		httpPort, pingCount, pingUrl, err)

	errx := terror.NewRawTerror(ctx, err, errMsg)

	return errx
}
