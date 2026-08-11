package service

import (
	"context"
	"fmt"
	"net/http"

	"github.com/choveylee/terror"
	"github.com/choveylee/tlog"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/load"
	"github.com/shirou/gopsutil/mem"

	"dev.choveylee.top/knowledge-base-backend/internal/const"
	"dev.choveylee.top/knowledge-base-backend/internal/data"
)

const (
	loadPerCoreCritical = 1.0
	loadPerCoreWarning  = 0.85
)

// CpuCheck returns the current CPU health status for the host.
func CpuCheck(ctx context.Context) (*data.CpuCheckRespData, *terror.Terror) {
	cores, err := cpu.Counts(false)
	if err != nil {
		errMsg := tlog.E(ctx).Err(err).Msgf("CPU check err (cpu counts %v)",
			err)

		errx := terror.NewTerror(ctx, err, constant.ErrorCodeUnknownServerAbnormal, errMsg)

		return nil, errx
	}

	if cores < 1 {
		cores = 1
	}

	avgStat, err := load.Avg()
	if err != nil {
		errMsg := tlog.E(ctx).Err(err).Msgf("CPU check err (load avg %v)",
			err)

		errx := terror.NewTerror(ctx, err, constant.ErrorCodeUnknownServerAbnormal, errMsg)

		return nil, errx
	}

	load1 := avgStat.Load1
	load5 := avgStat.Load5
	load15 := avgStat.Load15

	loadPerCore := load5 / float64(cores)

	statusCode := http.StatusOK
	status := "OK"

	if loadPerCore >= loadPerCoreCritical {
		statusCode = http.StatusInternalServerError
		status = "CRITICAL"
	} else if loadPerCore >= loadPerCoreWarning {
		statusCode = http.StatusTooManyRequests
		status = "WARNING"
	}

	cpuCheckRespData := &data.CpuCheckRespData{
		StatusCode: statusCode,
		Status:     status,
		Detail: fmt.Sprintf("%s - Load average: %.2f, %.2f, %.2f | Load/core: %.2f | Cores: %d",
			status, load1, load5, load15, loadPerCore, cores),
	}

	return cpuCheckRespData, nil
}

// RamCheck returns the current memory health status for the host.
func RamCheck(ctx context.Context) (*data.RamCheckRespData, *terror.Terror) {
	virtualMemoryStat, err := mem.VirtualMemory()
	if err != nil {
		errMsg := tlog.E(ctx).Err(err).Msgf("Ram check err (mem virtual memory %v)",
			err)

		errx := terror.NewTerror(ctx, err, constant.ErrorCodeUnknownServerAbnormal, errMsg)

		return nil, errx
	}

	usedMB := int(virtualMemoryStat.Used) / constant.MB
	usedGB := int(virtualMemoryStat.Used) / constant.GB

	totalMB := int(virtualMemoryStat.Total) / constant.MB
	totalGB := int(virtualMemoryStat.Total) / constant.GB

	usedPercent := int(virtualMemoryStat.UsedPercent)

	statusCode := http.StatusOK
	status := "OK"

	if usedPercent >= 95 {
		statusCode = http.StatusInternalServerError
		status = "CRITICAL"
	} else if usedPercent >= 90 {
		statusCode = http.StatusTooManyRequests
		status = "WARNING"
	}

	ramCheckRespData := &data.RamCheckRespData{
		StatusCode: statusCode,
		Status:     status,
		Detail: fmt.Sprintf("%s - Used: %dMB (%dGB) / Total: %dMB (%dGB) | Used: %d%%",
			status, usedMB, usedGB, totalMB, totalGB, usedPercent),
	}

	return ramCheckRespData, nil
}
