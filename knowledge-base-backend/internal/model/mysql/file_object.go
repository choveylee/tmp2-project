package dbmodel

import (
	"context"
	"strings"
	"time"

	"github.com/choveylee/terror"
	"github.com/choveylee/tlog"
	"github.com/choveylee/tutil"
	"gorm.io/gorm"

	constant "dev.choveylee.top/knowledge-base-backend/internal/const"
)

const (
	FileObjectBucketNameLen = 128
	FileObjectKeyLen        = 512
	FileObjectFileNameLen   = 255
	FileObjectMimeTypeLen   = 128
	FileObjectFileExtLen    = 32
)

const (
	FileObjectStorageProviderSeaweedFS = 0
)

type FileObject struct {
	Id string `gorm:"column:id"`

	BucketName string `gorm:"column:bucket_name"`
	ObjectKey  string `gorm:"column:object_key"`

	FileName string `gorm:"column:file_name"`
	MimeType string `gorm:"column:mime_type"`
	FileExt  string `gorm:"column:file_ext"`

	SizeBytes uint64 `gorm:"column:size_bytes"`
	Sha256    string `gorm:"column:sha256"`

	StorageProvider int `gorm:"column:storage_provider"`

	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

func (*FileObject) TableName() string {
	return "file_objects"
}

func CreateFileObject(ctx context.Context, bucketName, objectKey string, fileName, mimeType, fileExt string, sizeBytes uint64, sha256Value string, storageProvider int) (*FileObject, *terror.Terror) {
	fileObjectId := tutil.NewOid().String()
	if strings.TrimSpace(objectKey) == "" {
		objectKey = "pending/" + fileObjectId
	}

	fileObjectDB := &FileObject{
		Id: fileObjectId,

		BucketName: bucketName,
		ObjectKey:  objectKey,

		FileName: fileName,
		MimeType: mimeType,
		FileExt:  fileExt,

		SizeBytes: sizeBytes,
		Sha256:    sha256Value,

		StorageProvider: storageProvider,
	}

	retGorm := DB(ctx).Create(fileObjectDB)
	if retGorm.Error != nil {
		errMsg := tlog.E(ctx).Err(retGorm.Error).Msgf("Create file object (id: %s, bucket: %s, object key: %s, file name: %s, mime type: %s, file ext: %s, size bytes: %d, sha256: %s, storage provider: %d) err (db create %v)",
			fileObjectDB.Id, bucketName, fileObjectDB.ObjectKey, fileName, mimeType, fileExt, sizeBytes, sha256Value, storageProvider, retGorm.Error)

		errx := terror.NewTerror(ctx, retGorm.Error, constant.ErrorCodeMysqlServerAbnormal, errMsg)

		return nil, errx
	}

	return fileObjectDB, nil
}

func CreateFileObjectTx(ctx context.Context, tx *gorm.DB, bucketName, objectKey string, fileName, mimeType, fileExt string, sizeBytes uint64, sha256Value string, storageProvider int) (*FileObject, *terror.Terror) {
	fileObjectId := tutil.NewOid().String()
	if strings.TrimSpace(objectKey) == "" {
		objectKey = "pending/" + fileObjectId
	}

	fileObjectDB := &FileObject{
		Id: fileObjectId,

		BucketName: bucketName,
		ObjectKey:  objectKey,

		FileName: fileName,
		MimeType: mimeType,
		FileExt:  fileExt,

		SizeBytes: sizeBytes,
		Sha256:    sha256Value,

		StorageProvider: storageProvider,
	}

	retGorm := tx.Create(fileObjectDB)
	if retGorm.Error != nil {
		errMsg := tlog.E(ctx).Err(retGorm.Error).Msgf("Create file object tx (id: %s, bucket: %s, object key: %s, file name: %s, mime type: %s, file ext: %s, size bytes: %d, sha256: %s, storage provider: %d) err (db create %v)",
			fileObjectDB.Id, bucketName, fileObjectDB.ObjectKey, fileName, mimeType, fileExt, sizeBytes, sha256Value, storageProvider, retGorm.Error)

		errx := terror.NewTerror(ctx, retGorm.Error, constant.ErrorCodeMysqlServerAbnormal, errMsg)

		return nil, errx
	}

	return fileObjectDB, nil
}

func FindFileObject(ctx context.Context, fileObjectId string) (*FileObject, *terror.Terror) {
	fileObjectsDB := make([]*FileObject, 0)

	retGorm := DB(ctx).Where("id = ?", fileObjectId).Limit(1).Find(&fileObjectsDB)
	if retGorm.Error != nil {
		errMsg := tlog.E(ctx).Err(retGorm.Error).Msgf("Find file object (id: %s) err (db find %v)",
			fileObjectId, retGorm.Error)

		errx := terror.NewTerror(ctx, retGorm.Error, constant.ErrorCodeMysqlServerAbnormal, errMsg)

		return nil, errx
	}

	if len(fileObjectsDB) == 0 {
		return nil, nil
	}

	return fileObjectsDB[0], nil
}

func UpdateFileObjectStorageInfoTx(ctx context.Context, tx *gorm.DB, fileObjectId, objectKey, sha256Value string) *terror.Terror {
	params := map[string]any{
		"object_key": objectKey,
		"sha256":     sha256Value,

		"updated_at": time.Now(),
	}

	retGorm := tx.Model(&FileObject{}).Where("id = ?", fileObjectId).Updates(params)
	if retGorm.Error != nil {
		errMsg := tlog.E(ctx).Err(retGorm.Error).Msgf("Update file object storage info tx (id: %s, object key: %s, sha256: %s) err (db updates %v)",
			fileObjectId, objectKey, sha256Value, retGorm.Error)

		errx := terror.NewTerror(ctx, retGorm.Error, constant.ErrorCodeMysqlServerAbnormal, errMsg)

		return errx
	}

	return nil
}
