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
	DocumentParseStrategyAuto    = 0
	DocumentParseStrategyTika    = 1
	DocumentParseStrategyOCR     = 2
	DocumentParseStrategyTikaOCR = 3
)

var (
	DocumentParseStrategiesMap = map[int]bool{
		DocumentParseStrategyAuto:    true,
		DocumentParseStrategyTika:    true,
		DocumentParseStrategyOCR:     true,
		DocumentParseStrategyTikaOCR: true,
	}
)

const (
	DocumentParserTypeUnknown = 0
	DocumentParserTypeTika    = 1
	DocumentParserTypeOCR     = 2
	DocumentParserTypeMixed   = 3
	DocumentParserTypeManual  = 4
)

const (
	DocumentVersionParseStatusPending    = 0
	DocumentVersionParseStatusProcessing = 1
	DocumentVersionParseStatusSuccess    = 2
	DocumentVersionParseStatusFailed     = 3
)

const (
	DocumentVersionOcrStatusNotRequired = 0
	DocumentVersionOcrStatusPending     = 1
	DocumentVersionOcrStatusProcessing  = 2
	DocumentVersionOcrStatusSuccess     = 3
	DocumentVersionOcrStatusFailed      = 4
)

type DocumentVersion struct {
	Id string `gorm:"column:id"`

	DocumentId string `gorm:"column:document_id"`
	VersionNo  uint   `gorm:"column:version_no"`

	FileObjectId string `gorm:"column:file_object_id"`

	ParseStrategy int `gorm:"column:parse_strategy"`
	ParserType    int `gorm:"column:parser_type"`

	PlainText     string `gorm:"column:plain_text"`
	ContentSha256 string `gorm:"column:content_sha256"`

	PageCount  uint `gorm:"column:page_count"`
	TokenCount uint `gorm:"column:token_count"`
	ChunkCount uint `gorm:"column:chunk_count"`

	ParseStatus int    `gorm:"column:parse_status"`
	ParseError  string `gorm:"column:parse_error"`

	OcrStatus int    `gorm:"column:ocr_status"`
	OcrError  string `gorm:"column:ocr_error"`

	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

func (*DocumentVersion) TableName() string {
	return "document_versions"
}

func CreateDocumentVersion(ctx context.Context, tx *gorm.DB, documentId string, versionNo uint, fileObjectId string,
	parseStrategy, parserType int, contentSha256 string, pageCount, tokenCount, chunkCount uint, parseStatus int, parseError string, ocrStatus int, ocrError string) (*DocumentVersion, *terror.Terror) {
	documentVersionDB := &DocumentVersion{
		Id: tutil.NewOid().String(),

		DocumentId: documentId,
		VersionNo:  versionNo,

		FileObjectId: fileObjectId,

		ParseStrategy: parseStrategy,
		ParserType:    parserType,

		ContentSha256: contentSha256,

		PageCount:  pageCount,
		TokenCount: tokenCount,
		ChunkCount: chunkCount,

		ParseStatus: parseStatus,
		ParseError:  parseError,

		OcrStatus: ocrStatus,
		OcrError:  ocrError,
	}

	retGorm := tx.Create(documentVersionDB)
	if retGorm.Error != nil {
		errMsg := tlog.E(ctx).Err(retGorm.Error).Msgf("Create document version (id: %s, document id: %s, version no: %d, file object id: %s, parse strategy: %d, parser type: %d, content sha256: %s, page count: %d, token count: %d, chunk count: %d, parse status: %d, parse error: %s, ocr status: %d, ocr error: %s) err (db create %v)",
			documentVersionDB.Id, documentId, versionNo, fileObjectId, parseStrategy, parserType, contentSha256, pageCount, tokenCount, chunkCount, parseStatus, parseError, ocrStatus, ocrError, retGorm.Error)

		errx := terror.NewTerror(ctx, retGorm.Error, constant.ErrorCodeMysqlServerAbnormal, errMsg)

		return nil, errx
	}

	return documentVersionDB, nil
}

func FindDocumentVersion(ctx context.Context, versionId string) (*DocumentVersion, *terror.Terror) {
	documentVersionsDB := make([]*DocumentVersion, 0)

	retGorm := serverClient.DB(ctx, runMode).Where("id = ?", versionId).Limit(1).Find(&documentVersionsDB)
	if retGorm.Error != nil {
		errMsg := tlog.E(ctx).Err(retGorm.Error).Msgf("Find document version (id: %s) err (db find %v)",
			versionId, retGorm.Error)

		errx := terror.NewTerror(ctx, retGorm.Error, constant.ErrorCodeMysqlServerAbnormal, errMsg)

		return nil, errx
	}

	if len(documentVersionsDB) == 0 {
		return nil, nil
	}

	return documentVersionsDB[0], nil
}

func UpdateDocumentVersionContentSha256(ctx context.Context, tx *gorm.DB, versionId, contentSha256 string) *terror.Terror {
	params := map[string]any{
		"content_sha256": contentSha256,

		"updated_at": time.Now(),
	}

	retGorm := tx.Model(&DocumentVersion{}).Where("id = ?", versionId).Updates(params)
	if retGorm.Error != nil {
		errMsg := tlog.E(ctx).Err(retGorm.Error).Msgf("Update document version content sha256 (id: %s, content sha256: %s) err (db updates %v)",
			versionId, contentSha256, retGorm.Error)

		errx := terror.NewTerror(ctx, retGorm.Error, constant.ErrorCodeMysqlServerAbnormal, errMsg)

		return errx
	}

	return nil
}
