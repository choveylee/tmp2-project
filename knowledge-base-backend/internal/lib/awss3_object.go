package lib

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"mime/multipart"
	"path"
	"regexp"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/choveylee/terror"
	"github.com/choveylee/tlog"

	constant "dev.choveylee.top/knowledge-base-backend/internal/const"
)

var awsS3InvalidFileNameRegexp = regexp.MustCompile(`[^a-zA-Z0-9._-]+`)

func SanitizeAwsS3FileName(fileName string) string {
	fileName = strings.TrimSpace(fileName)
	if fileName == "" {
		return "unnamed"
	}

	fileName = strings.ReplaceAll(fileName, " ", "-")
	fileName = awsS3InvalidFileNameRegexp.ReplaceAllString(fileName, "-")
	fileName = strings.Trim(fileName, "-")
	if fileName == "" {
		return "unnamed"
	}

	return fileName
}

func BuildAwsS3RawObjectKey(knowledgeBaseId string, chatSessionId string, documentId string, versionId string, fileObjectName string, now time.Time) string {
	ownerPath := knowledgeBaseId
	if strings.TrimSpace(ownerPath) == "" {
		ownerPath = path.Join("chat", chatSessionId)
	}

	rawDir := path.Join("raw", ownerPath, documentId, versionId, now.Format("2006"), now.Format("01"), now.Format("02"))

	return path.Join(rawDir, fileObjectName)
}

func UploadAwsS3File(ctx context.Context, bucketName string, objectKey string, objectFileName string, fileHeader *multipart.FileHeader) (string, *terror.Terror) {
	fileName := ""
	fileSize := int64(0)
	if fileHeader != nil {
		fileName = fileHeader.Filename
		fileSize = fileHeader.Size
	}

	if fileHeader == nil {
		errMsg := tlog.E(ctx).Msgf("Upload aws s3 file (bucket: %s, object key: %s, object file name: %s, file name: %s, file size: %d) err (file header nil)",
			bucketName, objectKey, objectFileName, fileName, fileSize)

		errx := terror.NewTerror(ctx, terror.ErrParamInvalid("file"), constant.ErrorCodeObjectStorageUploadFailed, errMsg)

		return "", errx
	}

	if awsS3Client == nil {
		errMsg := tlog.E(ctx).Msgf("Upload aws s3 file (bucket: %s, object key: %s, object file name: %s, file name: %s, file size: %d) err (aws s3 client nil)",
			bucketName, objectKey, objectFileName, fileName, fileSize)

		errx := terror.NewTerror(ctx, terror.ErrConfInvalid("AWSS3_CLIENT"), constant.ErrorCodeObjectStorageUploadFailed, errMsg)

		return "", errx
	}

	bucketName = strings.Trim(bucketName, "/")
	if bucketName == "" {
		errMsg := tlog.E(ctx).Msgf("Upload aws s3 file (bucket: %s, object key: %s, object file name: %s, file name: %s, file size: %d) err (bucket empty)",
			bucketName, objectKey, objectFileName, fileName, fileSize)

		errx := terror.NewTerror(ctx, terror.ErrConfInvalid("AWSS3_BUCKET"), constant.ErrorCodeObjectStorageUploadFailed, errMsg)

		return "", errx
	}

	objectKey = strings.Trim(objectKey, "/")
	if objectKey == "" {
		errMsg := tlog.E(ctx).Msgf("Upload aws s3 file (bucket: %s, object key: %s, object file name: %s, file name: %s, file size: %d) err (object key empty)",
			bucketName, objectKey, objectFileName, fileName, fileSize)

		errx := terror.NewTerror(ctx, terror.ErrParamInvalid("object key"), constant.ErrorCodeObjectStorageUploadFailed, errMsg)

		return "", errx
	}

	objectFileName = strings.TrimSpace(objectFileName)

	file, err := fileHeader.Open()
	if err != nil {
		errMsg := tlog.E(ctx).Err(err).Msgf("Upload aws s3 file (bucket: %s, object key: %s, object file name: %s, file name: %s, file size: %d) err (open file %v)",
			bucketName, objectKey, objectFileName, fileName, fileSize, err)

		errx := terror.NewTerror(ctx, err, constant.ErrorCodeObjectStorageUploadFailed, errMsg)

		return "", errx
	}
	defer file.Close()

	hash := sha256.New()

	_, err = io.Copy(hash, file)
	if err != nil {
		errMsg := tlog.E(ctx).Err(err).Msgf("Upload aws s3 file (bucket: %s, object key: %s, object file name: %s, file name: %s, file size: %d) err (copy file %v)",
			bucketName, objectKey, objectFileName, fileName, fileSize, err)

		errx := terror.NewTerror(ctx, err, constant.ErrorCodeObjectStorageUploadFailed, errMsg)

		return "", errx
	}

	_, err = file.Seek(0, io.SeekStart)
	if err != nil {
		errMsg := tlog.E(ctx).Err(err).Msgf("Upload aws s3 file (bucket: %s, object key: %s, object file name: %s, file name: %s, file size: %d) err (seek file %v)",
			bucketName, objectKey, objectFileName, fileName, fileSize, err)

		errx := terror.NewTerror(ctx, err, constant.ErrorCodeObjectStorageUploadFailed, errMsg)

		return "", errx
	}

	putObjectInput := &s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectKey),

		Body:          file,
		ContentLength: aws.Int64(fileSize),
	}

	contentType := strings.TrimSpace(fileHeader.Header.Get("Content-Type"))
	if contentType != "" {
		putObjectInput.ContentType = aws.String(contentType)
	}

	_, err = awsS3Client.PutObject(ctx, putObjectInput)
	if err != nil {
		errMsg := tlog.E(ctx).Err(err).Msgf("Upload aws s3 file (bucket: %s, object key: %s, object file name: %s, file name: %s, file size: %d, content type: %s) err (put object %v)",
			bucketName, objectKey, objectFileName, fileName, fileSize, contentType, err)

		errx := terror.NewTerror(ctx, err, constant.ErrorCodeObjectStorageUploadFailed, errMsg)

		return "", errx
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}
