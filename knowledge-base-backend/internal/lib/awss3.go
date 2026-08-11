package lib

import (
	"context"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/choveylee/tcfg"
	"github.com/choveylee/terror"
	"github.com/choveylee/tlog"
)

var (
	awsS3Client *s3.Client

	awsS3PresignClient *s3.PresignClient

	awsS3Endpoint2 string
	awsS3Bucket    string
)

func AwsS3Endpoint2() string {
	return awsS3Endpoint2
}

func AwsS3Bucket() string {
	return awsS3Bucket
}

func initAmazonS3(ctx context.Context) *terror.Terror {
	awsS3Endpoint := strings.TrimSpace(tcfg.DefaultString(tcfg.LocalKey("AWSS3_ENDPOINT"), ""))
	if awsS3Endpoint == "" {
		errMsg := tlog.E(ctx).Msg("Init amazon s3 err (asw s3 endpoint illegal)")

		errx := terror.NewRawTerror(ctx, terror.ErrConfInvalid("asw s3 endpoint"), errMsg)

		return errx
	}

	awsS3Endpoint2 = strings.TrimSpace(tcfg.DefaultString(tcfg.LocalKey("AWSS3_ENDPOINT2"), ""))
	if awsS3Endpoint2 == "" {
		errMsg := tlog.E(ctx).Msg("Init amazon s3 err (asw s3 endpoint2 illegal)")

		errx := terror.NewRawTerror(ctx, terror.ErrConfInvalid("asw s3 endpoint2"), errMsg)

		return errx
	}

	awsS3AccessKey := strings.TrimSpace(tcfg.DefaultString(tcfg.LocalKey("AWSS3_ACCESS_KEY"), ""))
	if awsS3AccessKey == "" {
		errMsg := tlog.E(ctx).Msg("Init amazon s3 err (asw s3 access key illegal)")

		errx := terror.NewRawTerror(ctx, terror.ErrConfInvalid("asw s3 access key"), errMsg)

		return errx
	}

	awsS3SecretKey := strings.TrimSpace(tcfg.DefaultString(tcfg.LocalKey("AWSS3_SECRET_KEY"), ""))
	if awsS3SecretKey == "" {
		errMsg := tlog.E(ctx).Msg("Init amazon s3 err (asw s3 secret key illegal)")

		errx := terror.NewRawTerror(ctx, terror.ErrConfInvalid("asw s3 secret key"), errMsg)

		return errx
	}

	awsS3Bucket = strings.Trim(strings.TrimSpace(tcfg.DefaultString(tcfg.LocalKey("AWSS3_BUCKET"), "")), "/")
	if awsS3Bucket == "" {
		errMsg := tlog.E(ctx).Msg("Init amazon s3 err (asw s3 bucket illegal)")

		errx := terror.NewRawTerror(ctx, terror.ErrConfInvalid("asw s3 bucket"), errMsg)

		return errx
	}

	awsS3Region := strings.TrimSpace(tcfg.DefaultString(tcfg.LocalKey("AWSS3_REGION"), "us-east-1"))
	if awsS3Region == "" {
		errMsg := tlog.E(ctx).Msg("Init amazon s3 err (asw s3 region illegal)")

		errx := terror.NewRawTerror(ctx, terror.ErrConfInvalid("asw s3 region"), errMsg)

		return errx
	}

	awsConfig, err := awsconfig.LoadDefaultConfig(ctx,
		awsconfig.WithRegion(awsS3Region),
		awsconfig.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(awsS3AccessKey, awsS3SecretKey, "")),
	)
	if err != nil {
		errMsg := tlog.E(ctx).Err(err).Msgf("Init S3 object storage (asw s3 endpoint: %s, asw s3 bucket: %s, asw s3 region: %s) err (load aws config %v)",
			awsS3Endpoint, awsS3Bucket, awsS3Region, err)

		errx := terror.NewRawTerror(ctx, err, errMsg)

		return errx
	}

	awsS3Client = s3.NewFromConfig(awsConfig, func(options *s3.Options) {
		options.BaseEndpoint = aws.String(awsS3Endpoint)
		options.UsePathStyle = true
		options.RequestChecksumCalculation = aws.RequestChecksumCalculationWhenRequired
		options.ResponseChecksumValidation = aws.ResponseChecksumValidationWhenRequired
	})

	awsS3PresignClient = s3.NewPresignClient(
		s3.NewFromConfig(awsConfig, func(options *s3.Options) {
			options.BaseEndpoint = aws.String(awsS3Endpoint2)
			options.UsePathStyle = true
			options.RequestChecksumCalculation = aws.RequestChecksumCalculationWhenRequired
			options.ResponseChecksumValidation = aws.ResponseChecksumValidationWhenRequired
		}))

	checkCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	_, err = awsS3Client.HeadBucket(checkCtx, &s3.HeadBucketInput{
		Bucket: aws.String(awsS3Bucket),
	})
	if err != nil {
		errMsg := tlog.E(ctx).Err(err).Msgf("Init amazon s3 (asw s3 endpoint: %s, asw s3 bucket: %s, asw s3 region: %s) err (head bucket %v)",
			awsS3Endpoint, awsS3Bucket, awsS3Region, err)

		errx := terror.NewRawTerror(ctx, err, errMsg)

		return errx
	}

	return nil
}
