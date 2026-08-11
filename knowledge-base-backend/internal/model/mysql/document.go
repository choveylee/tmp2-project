package dbmodel

import (
	"context"
	"time"

	"github.com/choveylee/terror"
	"github.com/choveylee/tlog"
	"github.com/choveylee/tutil"
	"gorm.io/gorm"

	constant "dev.choveylee.top/knowledge-base-backend/internal/const"
	"dev.choveylee.top/knowledge-base-backend/internal/data"
)

const (
	DocumentTitleLenLimit   = 255
	DocumentSummaryLenLimit = 65535
	DocumentTagsLenLimit    = 4096
	DocumentLangCodeLen     = 16
)

const (
	DocumentScopeTypeKnowledge  = 0
	DocumentScopeTypeAttachment = 1
)

var (
	DocumentScopeTypesMap = map[int]bool{
		DocumentScopeTypeKnowledge:  true,
		DocumentScopeTypeAttachment: true,
	}
)

const (
	DocumentSourceTypeUser = 0
	DocumentSourceTypeApi  = 1
)

var (
	DocumentSourceTypesMap = map[int]bool{
		DocumentSourceTypeUser: true,
		DocumentSourceTypeApi:  true,
	}
)

const (
	DocumentProcessStatusWaiting    = 0
	DocumentProcessStatusUploaded   = 1
	DocumentProcessStatusProcessing = 2
	DocumentProcessStatusProcessed  = 3
	DocumentProcessStatusFailed     = 4
)

var (
	DocumentProcessStatusesMap = map[int]bool{
		DocumentProcessStatusWaiting:    true,
		DocumentProcessStatusUploaded:   true,
		DocumentProcessStatusProcessing: true,
		DocumentProcessStatusProcessed:  true,
		DocumentProcessStatusFailed:     true,
	}
)

const (
	DocumentStatusDisabled = 0
	DocumentStatusNormal   = 1
)

var (
	DocumentStatusesMap = map[int]bool{
		DocumentStatusDisabled: true,
		DocumentStatusNormal:   true,
	}
)

type Document struct {
	Id string `gorm:"column:id"`

	KnowledgeBaseId string `gorm:"column:knowledge_base_id"`
	ChatSessionId   string `gorm:"column:chat_session_id"`

	ScopeType  int `gorm:"column:scope_type"`
	SourceType int `gorm:"column:source_type"`

	Title   string `gorm:"column:title"`
	Summary string `gorm:"column:summary"`
	Tags    string `gorm:"column:tags"`

	OwnerId string `gorm:"column:owner_id"`

	LangCode string `gorm:"column:lang_code"`

	CurVersionNo uint   `gorm:"column:cur_version_no"`
	CurVersionId string `gorm:"column:cur_version_id"`

	ProcessStatus int `gorm:"column:process_status"`
	Status        int `gorm:"column:status"`

	CreatedAt time.Time      `gorm:"column:created_at"`
	UpdatedAt time.Time      `gorm:"column:updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at"`
}

func CreateDocument(ctx context.Context, tx *gorm.DB, knowledgeBaseId, chatSessionId string, scopeType, sourceType int,
	title, summary string, tags []string, ownerId, langCode string, curVersionNo uint, curVersionId string, processStatus, status int) (*Document, *terror.Terror) {
	documentDB := &Document{
		Id: tutil.NewOid().String(),

		KnowledgeBaseId: knowledgeBaseId,
		ChatSessionId:   chatSessionId,

		ScopeType:  scopeType,
		SourceType: sourceType,

		Title:   title,
		Summary: summary,
		Tags:    data.MarshalData(tags),

		OwnerId: ownerId,

		LangCode: langCode,

		CurVersionNo: curVersionNo,
		CurVersionId: curVersionId,

		ProcessStatus: processStatus,
		Status:        status,
	}

	retGorm := tx.Create(documentDB)
	if retGorm.Error != nil {
		errMsg := tlog.E(ctx).Err(retGorm.Error).Msgf("Create document (id: %s, knowledge base id: %s, chat session id: %s, scope type: %d, source type: %d, title: %s, summary: %s, tags: %v, owner id: %s, lang code: %s, cur version no: %d, cur version id: %s, process status: %d, status: %d) err (db create %v)",
			documentDB.Id, knowledgeBaseId, chatSessionId, scopeType, sourceType, title, summary, tags, ownerId, langCode, curVersionNo, curVersionId, processStatus, status, retGorm.Error)

		errx := terror.NewTerror(ctx, retGorm.Error, constant.ErrorCodeMysqlServerAbnormal, errMsg)

		return nil, errx
	}

	return documentDB, nil
}

func FindDocument(ctx context.Context, documentId string) (*Document, *terror.Terror) {
	documentsDB := make([]*Document, 0)

	retGorm := serverClient.DB(ctx, runMode).Where("id = ?", documentId).Limit(1).Find(&documentsDB)
	if retGorm.Error != nil {
		errMsg := tlog.E(ctx).Err(retGorm.Error).Msgf("Find document (id: %s) err (db find %v)",
			documentId, retGorm.Error)

		errx := terror.NewTerror(ctx, retGorm.Error, constant.ErrorCodeMysqlServerAbnormal, errMsg)

		return nil, errx
	}

	if len(documentsDB) == 0 {
		return nil, nil
	}

	return documentsDB[0], nil
}

func FindDocuments(ctx context.Context, knowledgeBaseId, keyword string, scopeType, sourceType, processStatus, status int, pageNum, pageSize int) (int64, []*Document, *terror.Terror) {
	query := serverClient.DB(ctx, runMode).Model(&Document{})

	if knowledgeBaseId != "" {
		query = query.Where("knowledge_base_id = ?", knowledgeBaseId)
	}

	if keyword != "" {
		query = query.Where("title LIKE ? OR summary LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}

	if scopeType != -1 {
		query = query.Where("scope_type = ?", scopeType)
	}

	if sourceType != -1 {
		query = query.Where("source_type = ?", sourceType)
	}

	if processStatus != -1 {
		query = query.Where("process_status = ?", processStatus)
	}

	if status != -1 {
		query = query.Where("status = ?", status)
	}

	total := int64(0)

	retGorm := query.Count(&total)
	if retGorm.Error != nil {
		errMsg := tlog.E(ctx).Err(retGorm.Error).Msgf("Find documents (knowledge base id: %s, keyword: %s, scope type: %d, source type: %d, process status: %d, status: %d, page num: %d, page size: %d) err (db count %v)",
			knowledgeBaseId, keyword, scopeType, sourceType, processStatus, status, pageNum, pageSize, retGorm.Error)

		errx := terror.NewTerror(ctx, retGorm.Error, constant.ErrorCodeMysqlServerAbnormal, errMsg)

		return -1, nil, errx
	}

	if total == 0 {
		return 0, make([]*Document, 0), nil
	}

	documentsDB := make([]*Document, 0)

	retGorm = query.Offset((pageNum - 1) * pageSize).Limit(pageSize).Order("created_at DESC").Find(&documentsDB)
	if retGorm.Error != nil {
		errMsg := tlog.E(ctx).Err(retGorm.Error).Msgf("Find documents (knowledge base id: %s, keyword: %s, scope type: %d, source type: %d, process status: %d, status: %d, page num: %d, page size: %d) err (db find %v)",
			knowledgeBaseId, keyword, scopeType, sourceType, processStatus, status, pageNum, pageSize, retGorm.Error)

		errx := terror.NewTerror(ctx, retGorm.Error, constant.ErrorCodeMysqlServerAbnormal, errMsg)

		return -1, nil, errx
	}

	return total, documentsDB, nil
}

func UpdateDocument(ctx context.Context, documentId string, title, summary string, tags []string, langCode string, status int) *terror.Terror {
	params := map[string]any{
		"title":   title,
		"summary": summary,

		"tags": data.MarshalData(tags),

		"lang_code": langCode,

		"status": status,

		"updated_at": time.Now(),
	}

	retGorm := serverClient.DB(ctx, runMode).Model(&Document{}).Where("id = ?", documentId).Updates(params)
	if retGorm.Error != nil {
		errMsg := tlog.E(ctx).Err(retGorm.Error).Msgf("Update document (id: %s, title: %s, summary: %s, tags: %v, lang code: %s, status: %d) err (db updates %v)",
			documentId, title, summary, tags, langCode, status, retGorm.Error)

		errx := terror.NewTerror(ctx, retGorm.Error, constant.ErrorCodeMysqlServerAbnormal, errMsg)

		return errx
	}

	return nil
}

func DisableDocument(ctx context.Context, tx *gorm.DB, documentId string) *terror.Terror {
	params := map[string]any{
		"status": DocumentStatusDisabled,

		"updated_at": time.Now(),
	}

	retGorm := tx.Model(&Document{}).Where("id = ?", documentId).Updates(params)
	if retGorm.Error != nil {
		errMsg := tlog.E(ctx).Err(retGorm.Error).Msgf("Disable document (id: %s) err (db updates %v)",
			documentId, retGorm.Error)

		errx := terror.NewTerror(ctx, retGorm.Error, constant.ErrorCodeMysqlServerAbnormal, errMsg)

		return errx
	}

	return nil
}

func UpdateDocumentCurrentVersion(ctx context.Context, tx *gorm.DB, documentId string, curVersionNo uint, curVersionId string) *terror.Terror {
	params := map[string]any{
		"cur_version_no": curVersionNo,
		"cur_version_id": curVersionId,

		"updated_at": time.Now(),
	}

	retGorm := tx.Model(&Document{}).Where("id = ?", documentId).Updates(params)
	if retGorm.Error != nil {
		errMsg := tlog.E(ctx).Err(retGorm.Error).Msgf("Update document current version (id: %s, cur version no: %d, cur version id: %s) err (db updates %v)",
			documentId, curVersionNo, curVersionId, retGorm.Error)

		errx := terror.NewTerror(ctx, retGorm.Error, constant.ErrorCodeMysqlServerAbnormal, errMsg)

		return errx
	}

	return nil
}

func UpdateDocumentProcessStatus(ctx context.Context, tx *gorm.DB, documentId string, processStatus int) *terror.Terror {
	params := map[string]any{
		"process_status": processStatus,

		"updated_at": time.Now(),
	}

	retGorm := tx.Model(&Document{}).Where("id = ?", documentId).Updates(params)
	if retGorm.Error != nil {
		errMsg := tlog.E(ctx).Err(retGorm.Error).Msgf("Update document process status (id: %s, process status: %d) err (db updates %v)",
			documentId, processStatus, retGorm.Error)

		errx := terror.NewTerror(ctx, retGorm.Error, constant.ErrorCodeMysqlServerAbnormal, errMsg)

		return errx
	}

	return nil
}

func DeleteDocument(ctx context.Context, tx *gorm.DB, documentId string) *terror.Terror {
	retGorm := tx.Where("id = ?", documentId).Delete(&Document{})
	if retGorm.Error != nil {
		errMsg := tlog.E(ctx).Err(retGorm.Error).Msgf("Delete document (id: %s) err (db delete %v)",
			documentId, retGorm.Error)

		errx := terror.NewTerror(ctx, retGorm.Error, constant.ErrorCodeMysqlServerAbnormal, errMsg)

		return errx
	}

	return nil
}
