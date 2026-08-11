// Package service implements service-layer logic exposed through HTTP handlers and scheduled jobs.
package service

import (
	"context"

	"github.com/choveylee/terror"

	"dev.choveylee.top/knowledge-base-backend/internal/lib"
)

var (
	awsS3Endpoint string
	awsS3Bucket   string
)

// InitService initializes service-layer dependencies.
func InitService(ctx context.Context) *terror.Terror {
	awsS3Endpoint = lib.AwsS3Endpoint2()
	awsS3Bucket = lib.AwsS3Bucket()

	return nil
}
