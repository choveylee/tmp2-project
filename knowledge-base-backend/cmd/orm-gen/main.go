// Command orm-gen generates GORM query code for the service schema.
package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/choveylee/tcfg"
	"github.com/choveylee/tdb"
	"github.com/choveylee/tlog"
	"gorm.io/gen"

	"dev.choveylee.top/knowledge-base-backend/internal/const"
)

func main() {
	ctx := context.Background()

	runMode := tcfg.DefaultString(tcfg.LocalKey("RUN_MODE"), constant.RunModeDebug)

	serverDsn := strings.TrimSpace(tcfg.DefaultString(fmt.Sprintf("%s::%s", runMode, tcfg.LocalKey("SERVER_MYSQL_DSN")), ""))
	if serverDsn == "" {
		tlog.E(ctx).Msgf("ORM code generation failed because configuration key %s is empty", "server mysql dsn")

		return
	}

	serverClient, err := tdb.NewMysqlClient(ctx, serverDsn)
	if err != nil {
		tlog.E(ctx).Err(err).Msg("ORM code generation failed while creating the MySQL client")

		return
	}

	generator := gen.NewGenerator(gen.Config{
		OutPath: "./model/query",
		Mode:    gen.WithoutContext | gen.WithDefaultQuery | gen.WithQueryInterface, // Generation mode.
	})

	sqlDB := serverClient.DB(ctx, runMode)

	generator.UseDB(sqlDB)

	// Generate the basic type-safe DAO API for the template table.
	generator.ApplyBasic(generator.GenerateModel("template"))

	// Generate the type-safe API with dynamic SQL declared on the query interface.
	generator.ApplyInterface(func(TemplateMethod) {}, generator.GenerateModel("template"))

	// Generate the code.
	generator.Execute()
}
