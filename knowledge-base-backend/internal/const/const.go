// Package constant defines shared constants used across the service.
package constant

const (
	RunModeDebug   = "debug"
	RunModeRelease = "release"
)

const (
	B  = 1
	KB = 1024 * B
	MB = 1024 * KB
	GB = 1024 * MB
)

const (
	RequestBodyMaxSize int64 = 20 * MB
)

const (
	MaxPageNum  = 100
	MaxPageSize = 100
)

const (
	SortSeqAsc  = "asc"
	SortSeqDesc = "desc"
)
