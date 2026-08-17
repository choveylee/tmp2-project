package dbmodel

import (
	"context"
	"time"

	"github.com/choveylee/terror"
	"github.com/choveylee/tlog"
	"gorm.io/gorm"

	constant "dev.choveylee.top/knowledge-base-backend/internal/const"
)

const (
	DocumentChunkHeadingPathLen     = 512
	DocumentChunkOpenSearchDocIdLen = 128
)

const (
	DocumentChunkTypeParagraph = 0
	DocumentChunkTypeHeading   = 1
	DocumentChunkTypeTable     = 2
	DocumentChunkTypeList      = 3
	DocumentChunkTypeCode      = 4
	DocumentChunkTypeOCR       = 5
)

const (
	DocumentChunkIndexStatusPending    = 0
	DocumentChunkIndexStatusProcessing = 1
	DocumentChunkIndexStatusSuccess    = 2
	DocumentChunkIndexStatusFailed     = 3
)

type DocumentChunk struct {
	Id string `gorm:"column:id"`

	DocumentId string `gorm:"column:document_id"`
	VersionId  string `gorm:"column:version_id"`
	ChunkNo    uint   `gorm:"column:chunk_no"`

	ChunkType   int    `gorm:"column:chunk_type"`
	HeadingPath string `gorm:"column:heading_path"`

	ChunkText string `gorm:"column:chunk_text"`
	TextHash  string `gorm:"column:text_hash"`

	PageFrom  uint `gorm:"column:page_from"`
	PageTo    uint `gorm:"column:page_to"`
	CharStart uint `gorm:"column:char_start"`
	CharEnd   uint `gorm:"column:char_end"`

	TokenCount uint `gorm:"column:token_count"`

	IndexStatus     int    `gorm:"column:index_status"`
	OpenSearchDocId string `gorm:"column:opensearch_doc_id"`

	CreatedAt time.Time      `gorm:"column:created_at"`
	UpdatedAt time.Time      `gorm:"column:updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at"`
}

func (*DocumentChunk) TableName() string {
	return "document_chunks"
}

func BatchCreateDocumentChunksTx(ctx context.Context, tx *gorm.DB, documentChunksDB []*DocumentChunk) *terror.Terror {
	retGorm := tx.Create(documentChunksDB)
	if retGorm.Error != nil {
		errMsg := tlog.E(ctx).Err(retGorm.Error).Msgf("Batch create document chunks tx (chunk count: %d) err (db create %v)",
			len(documentChunksDB), retGorm.Error)

		errx := terror.NewTerror(ctx, retGorm.Error, constant.ErrorCodeMysqlServerAbnormal, errMsg)

		return errx
	}

	return nil
}

func HardDeleteDocumentChunksByVersionTx(ctx context.Context, tx *gorm.DB, versionId string) *terror.Terror {
	retGorm := tx.Unscoped().Where("version_id = ?", versionId).Delete(&DocumentChunk{})
	if retGorm.Error != nil {
		errMsg := tlog.E(ctx).Err(retGorm.Error).Msgf("Hard delete document chunks by version tx (version id: %s) err (db delete %v)",
			versionId, retGorm.Error)

		errx := terror.NewTerror(ctx, retGorm.Error, constant.ErrorCodeMysqlServerAbnormal, errMsg)

		return errx
	}

	return nil
}
