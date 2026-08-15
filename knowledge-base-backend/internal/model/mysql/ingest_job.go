package dbmodel

import (
	"context"
	"time"

	"github.com/choveylee/terror"
	"github.com/choveylee/tlog"
	"github.com/choveylee/tutil"
	"gorm.io/gorm"

	constant "dev.choveylee.top/knowledge-base-backend/internal/const"
)

const (
	IngestJobTypeParse        = 0
	IngestJobTypeOCR          = 1
	IngestJobTypeSplit        = 2
	IngestJobTypeEmbedding    = 3
	IngestJobTypeIndex        = 4
	IngestJobTypeRebuildIndex = 5
)

const (
	IngestJobStatusPending    = 0
	IngestJobStatusProcessing = 1
	IngestJobStatusFinished   = 2
	IngestJobStatusFailed     = 3
	IngestJobStatusCanceled   = 4
)

type IngestJob struct {
	Id string `gorm:"column:id"`

	DocumentId string `gorm:"column:document_id"`
	VersionId  string `gorm:"column:version_id"`

	JobType   int `gorm:"column:job_type"`
	JobStatus int `gorm:"column:job_status"`

	RetryCount uint `gorm:"column:retry_count"`

	WorkerName   string `gorm:"column:worker_name"`
	ErrorMessage string `gorm:"column:error_message"`
	Payload      string `gorm:"column:payload"`

	StartedAt  *time.Time `gorm:"column:started_at"`
	FinishedAt *time.Time `gorm:"column:finished_at"`

	CreatedAt time.Time      `gorm:"column:created_at"`
	UpdatedAt time.Time      `gorm:"column:updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at"`
}

func (*IngestJob) TableName() string {
	return "ingest_jobs"
}

func CreateIngestJobTx(ctx context.Context, tx *gorm.DB, documentId, versionId string, jobType, jobStatus int, retryCount uint, workerName, errorMessage, payload string) (*IngestJob, *terror.Terror) {
	ingestJobDB := &IngestJob{
		Id: tutil.NewOid().String(),

		DocumentId: documentId,
		VersionId:  versionId,

		JobType:   jobType,
		JobStatus: jobStatus,

		RetryCount: retryCount,

		WorkerName:   workerName,
		ErrorMessage: errorMessage,
		Payload:      payload,
	}

	retGorm := tx.Create(ingestJobDB)
	if retGorm.Error != nil {
		errMsg := tlog.E(ctx).Err(retGorm.Error).Msgf("Create ingest job tx (id: %s, document id: %s, version id: %s, job type: %d, job status: %d, retry count: %d, worker name: %s, error message: %s, payload: %s) err (db create %v)",
			ingestJobDB.Id, documentId, versionId, jobType, jobStatus, retryCount, workerName, errorMessage, payload, retGorm.Error)

		errx := terror.NewTerror(ctx, retGorm.Error, constant.ErrorCodeMysqlServerAbnormal, errMsg)

		return nil, errx
	}

	return ingestJobDB, nil
}

func FindLatestIngestJob(ctx context.Context, documentId string, versionId string) (*IngestJob, *terror.Terror) {
	ingestJobsDB := make([]*IngestJob, 0)

	query := DB(ctx).Where("document_id = ?", documentId)

	if versionId != "" {
		query = query.Where("version_id = ?", versionId)
	}

	retGorm := query.Order("created_at DESC").Limit(1).Find(&ingestJobsDB)
	if retGorm.Error != nil {
		errMsg := tlog.E(ctx).Err(retGorm.Error).Msgf("Find latest ingest job (document id: %s, version id: %s) err (db find %v)",
			documentId, versionId, retGorm.Error)

		errx := terror.NewTerror(ctx, retGorm.Error, constant.ErrorCodeMysqlServerAbnormal, errMsg)

		return nil, errx
	}

	if len(ingestJobsDB) == 0 {
		return nil, nil
	}

	return ingestJobsDB[0], nil
}
