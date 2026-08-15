package service

import (
	"context"
	"encoding/json"
	"mime/multipart"
	"path"
	"strings"
	"time"

	"github.com/choveylee/terror"
	"github.com/choveylee/tlog"
	"gorm.io/gorm"

	constant "dev.choveylee.top/knowledge-base-backend/internal/const"
	"dev.choveylee.top/knowledge-base-backend/internal/data"
	"dev.choveylee.top/knowledge-base-backend/internal/lib"
	dbmodel "dev.choveylee.top/knowledge-base-backend/internal/model/mysql"
)

func ListDocuments(ctx context.Context, userId, knowledgeBaseId, keyword string, scopeType, sourceType, processStatus, status int, pageNum, pageSize int) (*data.ListDocumentsRespData, *terror.Terror) {
	total, documentsDB, errx := dbmodel.FindDocuments(ctx, knowledgeBaseId, keyword, scopeType, sourceType, processStatus, status, pageNum, pageSize)
	if errx != nil {
		errMsg := tlog.E(ctx).Err(errx).Msgf("List documents (user id: %s, knowledge base id: %s, keyword: %s, scope type: %d, source type: %d, process status: %d, status: %d, page num: %d, page size: %d) err (db find documents %v)",
			userId, knowledgeBaseId, keyword, scopeType, sourceType, processStatus, status, pageNum, pageSize, errx)
		errx.AttachErrMsg(errMsg)

		return nil, errx
	}

	listDocumentsRespData := &data.ListDocumentsRespData{
		List: make([]*data.DocumentData, 0),

		Total: total,
	}

	// List responses include only document table data to avoid N+1 version/file queries.
	for _, documentDB := range documentsDB {
		tags := make([]string, 0)

		srcTags := strings.TrimSpace(documentDB.Tags)
		if srcTags != "" && srcTags != "null" {
			_ = json.Unmarshal([]byte(srcTags), &tags)
		}

		documentData := &data.DocumentData{
			DocumentId: documentDB.Id,

			KnowledgeBaseId: documentDB.KnowledgeBaseId,
			ChatSessionId:   documentDB.ChatSessionId,

			ScopeType:  documentDB.ScopeType,
			SourceType: documentDB.SourceType,

			Title:   documentDB.Title,
			Summary: documentDB.Summary,
			Tags:    tags,

			OwnerId: documentDB.OwnerId,

			LangCode: documentDB.LangCode,

			CurVersionNo: documentDB.CurVersionNo,
			CurVersionId: documentDB.CurVersionId,

			ProcessStatus: documentDB.ProcessStatus,
			Status:        documentDB.Status,

			CreatedAt: documentDB.CreatedAt.Format(time.RFC3339),
			UpdatedAt: documentDB.UpdatedAt.Format(time.RFC3339),
		}

		listDocumentsRespData.List = append(listDocumentsRespData.List, documentData)
	}

	return listDocumentsRespData, nil
}

func GetDocument(ctx context.Context, userId string, documentId string) (*data.GetDocumentRespData, *terror.Terror) {
	documentDB, errx := dbmodel.FindDocument(ctx, documentId)
	if errx != nil {
		errMsg := tlog.E(ctx).Err(errx).Msgf("Get document (user id: %s, document id: %s) err (db find document %v)",
			userId, documentId, errx)
		errx.AttachErrMsg(errMsg)

		return nil, errx
	}

	if documentDB == nil {
		errMsg := tlog.E(ctx).Msgf("Get document (user id: %s, document id: %s) err (document not found)",
			userId, documentId)

		errx := terror.NewTerror(ctx, terror.ErrParamInvalid("document id"), constant.ErrorCodeDocumentNotFound, errMsg)

		return nil, errx
	}

	tags := make([]string, 0)

	srcTags := strings.TrimSpace(documentDB.Tags)
	if srcTags != "" && srcTags != "null" {
		_ = json.Unmarshal([]byte(srcTags), &tags)
	}

	curVersionId := documentDB.CurVersionId

	// Build the main document payload first, then enrich it with current version details.
	documentData := &data.DocumentData{
		DocumentId: documentDB.Id,

		KnowledgeBaseId: documentDB.KnowledgeBaseId,
		ChatSessionId:   documentDB.ChatSessionId,

		ScopeType:  documentDB.ScopeType,
		SourceType: documentDB.SourceType,

		Title:   documentDB.Title,
		Summary: documentDB.Summary,
		Tags:    tags,

		OwnerId: documentDB.OwnerId,

		LangCode: documentDB.LangCode,

		CurVersionNo: documentDB.CurVersionNo,
		CurVersionId: curVersionId,

		ProcessStatus: documentDB.ProcessStatus,
		Status:        documentDB.Status,

		CreatedAt: documentDB.CreatedAt.Format(time.RFC3339),
		UpdatedAt: documentDB.UpdatedAt.Format(time.RFC3339),
	}

	if curVersionId != "" {
		// Detail responses include the current version and its original file object.
		documentVersionDB, errx := dbmodel.FindDocumentVersion(ctx, curVersionId)
		if errx != nil {
			errMsg := tlog.E(ctx).Err(errx).Msgf("Get document (user id: %s, document id: %s, version id: %s) err (db find document version %v)",
				userId, documentId, curVersionId, errx)
			errx.AttachErrMsg(errMsg)

			return nil, errx
		}

		if documentVersionDB != nil {
			fileObjectId := documentVersionDB.FileObjectId

			documentVersionData := &data.DocumentVersionData{
				VersionId: documentVersionDB.Id,

				DocumentId: documentVersionDB.DocumentId,
				VersionNo:  documentVersionDB.VersionNo,

				FileObjectId: fileObjectId,

				ParseStrategy: documentVersionDB.ParseStrategy,
				ParserType:    documentVersionDB.ParserType,

				ContentSha256: documentVersionDB.ContentSha256,

				PageCount:  documentVersionDB.PageCount,
				TokenCount: documentVersionDB.TokenCount,
				ChunkCount: documentVersionDB.ChunkCount,

				ParseStatus: documentVersionDB.ParseStatus,
				ParseError:  documentVersionDB.ParseError,

				OcrStatus: documentVersionDB.OcrStatus,
				OcrError:  documentVersionDB.OcrError,

				CreatedAt: documentVersionDB.CreatedAt.Format(time.RFC3339),
				UpdatedAt: documentVersionDB.UpdatedAt.Format(time.RFC3339),
			}

			documentData.CurrentVersion = documentVersionData

			fileObjectDB, errx := dbmodel.FindFileObject(ctx, fileObjectId)
			if errx != nil {
				errMsg := tlog.E(ctx).Err(errx).Msgf("Get document (user id: %s, document id: %s, version id: %s, file object id: %s) err (db find file object %v)",
					userId, documentId, curVersionId, fileObjectId, errx)
				errx.AttachErrMsg(errMsg)

				return nil, errx
			}

			if fileObjectDB != nil {
				fileObjectData := &data.FileObjectData{
					FileObjectId: fileObjectDB.Id,

					BucketName: fileObjectDB.BucketName,
					ObjectKey:  fileObjectDB.ObjectKey,

					FileName: fileObjectDB.FileName,
					MimeType: fileObjectDB.MimeType,
					FileExt:  fileObjectDB.FileExt,

					SizeBytes: fileObjectDB.SizeBytes,
					Sha256:    fileObjectDB.Sha256,

					StorageProvider: fileObjectDB.StorageProvider,

					CreatedAt: fileObjectDB.CreatedAt.Format(time.RFC3339),
				}

				documentData.FileObject = fileObjectData
			}
		}
	}

	// Ingest jobs run asynchronously, so expose only the latest job for the current version.
	ingestJobDB, errx := dbmodel.FindLatestIngestJob(ctx, documentId, curVersionId)
	if errx != nil {
		errMsg := tlog.E(ctx).Err(errx).Msgf("Get document (user id: %s, document id: %s, version id: %s) err (db find latest ingest job %v)",
			userId, documentId, curVersionId, errx)
		errx.AttachErrMsg(errMsg)

		return nil, errx
	}

	if ingestJobDB != nil {
		ingestJobData := &data.IngestJobData{
			JobId: ingestJobDB.Id,

			DocumentId: ingestJobDB.DocumentId,
			VersionId:  ingestJobDB.VersionId,

			JobType:   ingestJobDB.JobType,
			JobStatus: ingestJobDB.JobStatus,

			RetryCount: ingestJobDB.RetryCount,

			WorkerName:   ingestJobDB.WorkerName,
			ErrorMessage: ingestJobDB.ErrorMessage,
			Payload:      ingestJobDB.Payload,

			CreatedAt: ingestJobDB.CreatedAt.Format(time.RFC3339),
			UpdatedAt: ingestJobDB.UpdatedAt.Format(time.RFC3339),
		}

		if ingestJobDB.StartedAt != nil {
			ingestJobData.StartedAt = ingestJobDB.StartedAt.Format(time.RFC3339)
		}

		if ingestJobDB.FinishedAt != nil {
			ingestJobData.FinishedAt = ingestJobDB.FinishedAt.Format(time.RFC3339)
		}

		documentData.LatestJob = ingestJobData
	}

	getDocumentRespData := &data.GetDocumentRespData{
		DocumentData: documentData,
	}

	return getDocumentRespData, nil
}

func CreateDocument(ctx context.Context, userId, knowledgeBaseId, chatSessionId string, scopeType, sourceType int,
	title, summary string, tags []string, ownerId, langCode string, parseStrategy int, fileHeader *multipart.FileHeader) (*data.CreateDocumentRespData, *terror.Terror) {
	originFileName := strings.TrimSpace(fileHeader.Filename)
	fileSize := fileHeader.Size

	if scopeType == dbmodel.DocumentScopeTypeKnowledge {
		// Knowledge documents must reference an existing knowledge base; attachments do not.
		knowledgeBaseDB, errx := dbmodel.FindKnowledgeBase(ctx, knowledgeBaseId)
		if errx != nil {
			errMsg := tlog.E(ctx).Err(errx).Msgf("Create document (user id: %s, knowledge base id: %s, chat session id: %s, scope type: %d, source type: %d, title: %s, summary: %s, tags: %v, owner id: %s, lang code: %s, parse strategy: %d, file name: %s, file size: %d) err (db find knowledge base %v)",
				userId, knowledgeBaseId, chatSessionId, scopeType, sourceType, title, summary, tags, ownerId, langCode, parseStrategy, originFileName, fileSize, errx)
			errx.AttachErrMsg(errMsg)

			return nil, errx
		}

		if knowledgeBaseDB == nil {
			errMsg := tlog.E(ctx).Msgf("Create document (user id: %s, knowledge base id: %s, chat session id: %s, scope type: %d, source type: %d, title: %s, summary: %s, tags: %v, owner id: %s, lang code: %s, parse strategy: %d, file name: %s, file size: %d) err (knowledge base not found)",
				userId, knowledgeBaseId, chatSessionId, scopeType, sourceType, title, summary, tags, ownerId, langCode, parseStrategy, originFileName, fileSize)

			errx := terror.NewTerror(ctx, terror.ErrParamInvalid("knowledge_base_id"), constant.ErrorCodeKnowledgeBaseNotFound, errMsg)

			return nil, errx
		}
	}

	if title == "" {
		title = strings.TrimSpace(strings.TrimSuffix(originFileName, path.Ext(originFileName)))
	}
	if title == "" {
		title = "unnamed"
	}

	bucketName := lib.AwsS3Bucket()

	mimeType := strings.TrimSpace(fileHeader.Header.Get("Content-Type"))
	if mimeType == "" {
		mimeType = "application/octet-stream"
	}
	if len(mimeType) > dbmodel.FileObjectMimeTypeLen {
		mimeType = mimeType[:dbmodel.FileObjectMimeTypeLen]
	}

	fileExt := strings.TrimPrefix(strings.ToLower(path.Ext(originFileName)), ".")
	if len(fileExt) > dbmodel.FileObjectFileExtLen {
		fileExt = fileExt[:dbmodel.FileObjectFileExtLen]
	}
	if fileExt == "" {
		fileExt = "bin"
	}

	ocrStatus := dbmodel.DocumentVersionOcrStatusNotRequired
	if parseStrategy == dbmodel.DocumentParseStrategyOCR || parseStrategy == dbmodel.DocumentParseStrategyTikaOCR {
		ocrStatus = dbmodel.DocumentVersionOcrStatusPending
	}

	// The first transaction creates placeholder records so document/version/file IDs are available.
	var errx *terror.Terror

	var documentDB *dbmodel.Document
	var documentVersionDB *dbmodel.DocumentVersion
	var fileObjectDB *dbmodel.FileObject

	documentId := ""
	versionId := ""
	fileObjectId := ""

	err := dbmodel.DB(ctx).Transaction(func(tx *gorm.DB) error {
		documentDB, errx = dbmodel.CreateDocumentTx(ctx, tx, knowledgeBaseId, chatSessionId, scopeType, sourceType,
			title, summary, tags, ownerId, langCode, 0, "", dbmodel.DocumentProcessStatusWaiting, dbmodel.DocumentStatusNormal)
		if errx != nil {
			errMsg := tlog.E(ctx).Err(errx).Msgf("Create document tx (user id: %s, knowledge base id: %s, chat session id: %s, scope type: %d, source type: %d, title: %s, summary: %s, tags: %v, owner id: %s, lang code: %s, parse strategy: %d, file name: %s, file size: %d) err (db prepare document records transaction %v)",
				userId, knowledgeBaseId, chatSessionId, scopeType, sourceType, title, summary, tags, ownerId, langCode, parseStrategy, originFileName, fileSize, errx)
			errx.AttachErrMsg(errMsg)

			return errx
		}
		documentId = documentDB.Id

		fileObjectDB, errx = dbmodel.CreateFileObjectTx(ctx, tx, bucketName, "", originFileName, mimeType, fileExt,
			uint64(fileSize), "", dbmodel.FileObjectStorageProviderSeaweedFS)
		if errx != nil {
			errMsg := tlog.E(ctx).Err(errx).Msgf("Create document tx (user id: %s, knowledge base id: %s, chat session id: %s, scope type: %d, source type: %d, title: %s, summary: %s, tags: %v, owner id: %s, lang code: %s, parse strategy: %d, document id: %s, file name: %s, file size: %d) err (db prepare document records transaction %v)",
				userId, knowledgeBaseId, chatSessionId, scopeType, sourceType, title, summary, tags, ownerId, langCode, parseStrategy, documentId, originFileName, fileSize, errx)
			errx.AttachErrMsg(errMsg)

			return errx
		}
		fileObjectId = fileObjectDB.Id

		documentVersionDB, errx = dbmodel.CreateDocumentVersionTx(ctx, tx, documentId, 1, fileObjectId, parseStrategy,
			dbmodel.DocumentParserTypeUnknown, "", 0, 0, 0, dbmodel.DocumentVersionParseStatusPending, "", ocrStatus, "")
		if errx != nil {
			errMsg := tlog.E(ctx).Err(errx).Msgf("Create document tx (user id: %s, knowledge base id: %s, chat session id: %s, scope type: %d, source type: %d, title: %s, summary: %s, tags: %v, owner id: %s, lang code: %s, parse strategy: %d, document id: %s, file object id: %s, file name: %s, file size: %d) err (db prepare document records transaction %v)",
				userId, knowledgeBaseId, chatSessionId, scopeType, sourceType, title, summary, tags, ownerId, langCode, parseStrategy, documentId, fileObjectId, originFileName, fileSize, errx)
			errx.AttachErrMsg(errMsg)

			return errx
		}
		versionId = documentVersionDB.Id

		errx = dbmodel.UpdateDocumentCurrentVersionTx(ctx, tx, documentId, documentVersionDB.VersionNo, versionId)
		if errx != nil {
			errMsg := tlog.E(ctx).Err(errx).Msgf("Create document tx (user id: %s, knowledge base id: %s, chat session id: %s, scope type: %d, source type: %d, title: %s, summary: %s, tags: %v, owner id: %s, lang code: %s, parse strategy: %d, document id: %s, version id: %s, file object id: %s, file name: %s, file size: %d) err (db prepare document records transaction %v)",
				userId, knowledgeBaseId, chatSessionId, scopeType, sourceType, title, summary, tags, ownerId, langCode, parseStrategy, documentId, versionId, fileObjectId, originFileName, fileSize, errx)
			errx.AttachErrMsg(errMsg)

			return errx
		}

		return nil
	})
	if err != nil {
		errMsg := tlog.E(ctx).Err(err).Msgf("Create document (user id: %s, knowledge base id: %s, chat session id: %s, scope type: %d, source type: %d, title: %s, summary: %s, tags: %v, owner id: %s, lang code: %s, parse strategy: %d, document id: %s, version id: %s, file object id: %s, bucket: %s, file name: %s, file size: %d) err (db prepare document records transaction %v)",
			userId, knowledgeBaseId, chatSessionId, scopeType, sourceType, title, summary, tags, ownerId, langCode, parseStrategy, documentId, versionId, fileObjectId, bucketName, originFileName, fileSize, err)

		if errx != nil {
			errx.AttachErrMsg(errMsg)
		} else {
			errx = terror.NewTerror(ctx, err, constant.ErrorCodeMysqlServerAbnormal, errMsg)
		}

		return nil, errx
	}

	objectFileName := fileObjectId + "." + fileExt

	// Object storage names use generated IDs; the original filename stays in DB for display.
	objectKey := lib.BuildAwsS3RawObjectKey(knowledgeBaseId, chatSessionId, documentId, versionId, objectFileName, time.Now())

	sha256Value, errx := lib.UploadAwsS3File(ctx, bucketName, objectKey, objectFileName, fileHeader)
	if errx != nil {
		errMsg := tlog.E(ctx).Err(errx).Msgf("Create document (user id: %s, knowledge base id: %s, chat session id: %s, scope type: %d, source type: %d, title: %s, summary: %s, tags: %v, owner id: %s, lang code: %s, parse strategy: %d, document id: %s, version id: %s, file object id: %s, bucket: %s, object key: %s, object file name: %s, file name: %s, file size: %d) err (lib upload aws s3 file %v)",
			userId, knowledgeBaseId, chatSessionId, scopeType, sourceType, title, summary, tags, ownerId, langCode, parseStrategy, documentId, versionId, fileObjectId, bucketName, objectKey, objectFileName, originFileName, fileSize, errx)
		errx.AttachErrMsg(errMsg)

		// Upload failure only needs one status update, so avoid wrapping it in a transaction.
		errx2 := dbmodel.UpdateDocumentProcessStatus(ctx, documentId, dbmodel.DocumentProcessStatusFailed)
		if errx2 != nil {
			errMsg := tlog.E(ctx).Err(errx2).Msgf("Create document (user id: %s, knowledge base id: %s, chat session id: %s, scope type: %d, source type: %d, title: %s, summary: %s, tags: %v, owner id: %s, lang code: %s, parse strategy: %d, document id: %s, version id: %s, file object id: %s, bucket: %s, object key: %s, file name: %s, file size: %d) err (db mark document upload failed %v)",
				userId, knowledgeBaseId, chatSessionId, scopeType, sourceType, title, summary, tags, ownerId, langCode, parseStrategy, documentId, versionId, fileObjectId, bucketName, objectKey, originFileName, fileSize, errx2)
			errx2.AttachErrMsg(errMsg)
		}

		return nil, errx
	}

	jobPayloadBytes, _ := json.Marshal(map[string]string{
		"bucket_name": bucketName,
		"object_key":  objectKey,
	})

	jobPayload := string(jobPayloadBytes)

	var ingestJobDB *dbmodel.IngestJob

	// The second transaction completes storage metadata and creates the async parse job.
	err = dbmodel.DB(ctx).Transaction(func(tx *gorm.DB) error {
		errx = dbmodel.UpdateFileObjectStorageInfoTx(ctx, tx, fileObjectId, objectKey, sha256Value)
		if errx != nil {
			errMsg := tlog.E(ctx).Err(errx).Msgf("Create document tx (user id: %s, knowledge base id: %s, chat session id: %s, scope type: %d, source type: %d, title: %s, summary: %s, tags: %v, owner id: %s, lang code: %s, parse strategy: %d, document id: %s, version id: %s, file object id: %s, bucket: %s, object key: %s, sha256: %s, file name: %s, file size: %d) err (db update file object storage info %v)",
				userId, knowledgeBaseId, chatSessionId, scopeType, sourceType, title, summary, tags, ownerId, langCode, parseStrategy, documentId, versionId, fileObjectId, bucketName, objectKey, sha256Value, originFileName, fileSize, errx)
			errx.AttachErrMsg(errMsg)

			return errx
		}

		errx = dbmodel.UpdateDocumentVersionContentSha256Tx(ctx, tx, versionId, sha256Value)
		if errx != nil {
			errMsg := tlog.E(ctx).Err(errx).Msgf("Create document tx (user id: %s, knowledge base id: %s, chat session id: %s, scope type: %d, source type: %d, title: %s, summary: %s, tags: %v, owner id: %s, lang code: %s, parse strategy: %d, document id: %s, version id: %s, file object id: %s, bucket: %s, object key: %s, sha256: %s, file name: %s, file size: %d) err (db update document version content sha256 %v)",
				userId, knowledgeBaseId, chatSessionId, scopeType, sourceType, title, summary, tags, ownerId, langCode, parseStrategy, documentId, versionId, fileObjectId, bucketName, objectKey, sha256Value, originFileName, fileSize, errx)
			errx.AttachErrMsg(errMsg)

			return errx
		}

		errx = dbmodel.UpdateDocumentProcessStatusTx(ctx, tx, documentId, dbmodel.DocumentProcessStatusUploaded)
		if errx != nil {
			errMsg := tlog.E(ctx).Err(errx).Msgf("Create document tx (user id: %s, knowledge base id: %s, chat session id: %s, scope type: %d, source type: %d, title: %s, summary: %s, tags: %v, owner id: %s, lang code: %s, parse strategy: %d, document id: %s, version id: %s, file object id: %s, bucket: %s, object key: %s, sha256: %s, file name: %s, file size: %d) err (db update document process status %v)",
				userId, knowledgeBaseId, chatSessionId, scopeType, sourceType, title, summary, tags, ownerId, langCode, parseStrategy, documentId, versionId, fileObjectId, bucketName, objectKey, sha256Value, originFileName, fileSize, errx)
			errx.AttachErrMsg(errMsg)

			return errx
		}

		ingestJobDB, errx = dbmodel.CreateIngestJobTx(ctx, tx, documentId, versionId, dbmodel.IngestJobTypeParse,
			dbmodel.IngestJobStatusPending, 0, "", "", jobPayload)
		if errx != nil {
			errMsg := tlog.E(ctx).Err(errx).Msgf("Create document tx (user id: %s, knowledge base id: %s, chat session id: %s, scope type: %d, source type: %d, title: %s, summary: %s, tags: %v, owner id: %s, lang code: %s, parse strategy: %d, document id: %s, version id: %s, file object id: %s, bucket: %s, object key: %s, sha256: %s, file name: %s, file size: %d) err (db create ingest job %v)",
				userId, knowledgeBaseId, chatSessionId, scopeType, sourceType, title, summary, tags, ownerId, langCode, parseStrategy, documentId, versionId, fileObjectId, bucketName, objectKey, sha256Value, originFileName, fileSize, errx)
			errx.AttachErrMsg(errMsg)

			return errx
		}

		return nil
	})
	if err != nil {
		jobId := ""
		if ingestJobDB != nil {
			jobId = ingestJobDB.Id
		}

		errMsg := tlog.E(ctx).Err(err).Msgf("Create document (user id: %s, knowledge base id: %s, chat session id: %s, scope type: %d, source type: %d, title: %s, summary: %s, tags: %v, owner id: %s, lang code: %s, parse strategy: %d, document id: %s, version id: %s, file object id: %s, job id: %s, bucket: %s, object key: %s, sha256: %s, file name: %s, file size: %d) err (db complete document upload transaction %v)",
			userId, knowledgeBaseId, chatSessionId, scopeType, sourceType, title, summary, tags, ownerId, langCode, parseStrategy, documentId, versionId, fileObjectId, jobId, bucketName, objectKey, sha256Value, originFileName, fileSize, err)

		if errx != nil {
			errx.AttachErrMsg(errMsg)
		} else {
			errx = terror.NewTerror(ctx, err, constant.ErrorCodeMysqlServerAbnormal, errMsg)
		}

		return nil, errx
	}

	createDocumentRespData := &data.CreateDocumentRespData{
		DocumentId: documentId,

		VersionId:    versionId,
		FileObjectId: fileObjectId,
		JobId:        ingestJobDB.Id,

		BucketName: bucketName,
		ObjectKey:  objectKey,

		ProcessStatus: dbmodel.DocumentProcessStatusUploaded,
		Status:        dbmodel.DocumentStatusNormal,
	}

	return createDocumentRespData, nil
}

func UpdateDocument(ctx context.Context, userId, documentId string, title, summary string, tags []string, langCode string, status int) *terror.Terror {
	documentDB, errx := dbmodel.FindDocument(ctx, documentId)
	if errx != nil {
		errMsg := tlog.E(ctx).Err(errx).Msgf("Update document (user id: %s, document id: %s, title: %s, summary: %s, tags: %v, lang code: %s, status: %d) err (db find document %v)",
			userId, documentId, title, summary, tags, langCode, status, errx)
		errx.AttachErrMsg(errMsg)

		return errx
	}

	if documentDB == nil {
		errMsg := tlog.E(ctx).Msgf("Update document (user id: %s, document id: %s, title: %s, summary: %s, tags: %v, lang code: %s, status: %d) err (document not found)",
			userId, documentId, title, summary, tags, langCode, status)

		errx := terror.NewTerror(ctx, terror.ErrParamInvalid("document id"), constant.ErrorCodeDocumentNotFound, errMsg)

		return errx
	}

	errx = dbmodel.UpdateDocument(ctx, documentId, title, summary, tags, langCode, status)
	if errx != nil {
		errMsg := tlog.E(ctx).Err(errx).Msgf("Update document (user id: %s, document id: %s, title: %s, summary: %s, tags: %v, lang code: %s, status: %d) err (db update document %v)",
			userId, documentId, title, summary, tags, langCode, status, errx)
		errx.AttachErrMsg(errMsg)

		return errx
	}

	return nil
}

func DeleteDocument(ctx context.Context, userId string, documentId string) *terror.Terror {
	documentDB, errx := dbmodel.FindDocument(ctx, documentId)
	if errx != nil {
		errMsg := tlog.E(ctx).Err(errx).Msgf("Delete document (user id: %s, document id: %s) err (db find document %v)",
			userId, documentId, errx)
		errx.AttachErrMsg(errMsg)

		return errx
	}

	if documentDB == nil {
		errMsg := tlog.E(ctx).Msgf("Delete document (user id: %s, document id: %s) err (document not found)",
			userId, documentId)

		errx := terror.NewTerror(ctx, terror.ErrParamInvalid("document id"), constant.ErrorCodeDocumentNotFound, errMsg)

		return errx
	}

	// Delete disables business status and soft-deletes the row atomically.
	err := dbmodel.DB(ctx).Transaction(func(tx *gorm.DB) error {
		errx = dbmodel.DisableDocumentTx(ctx, tx, documentId)
		if errx != nil {
			errMsg := tlog.E(ctx).Err(errx).Msgf("Delete document tx (user id: %s, document id: %s) err (db transaction %v)",
				userId, documentId, errx)
			errx.AttachErrMsg(errMsg)

			return errx
		}

		errx = dbmodel.DeleteDocumentTx(ctx, tx, documentId)
		if errx != nil {
			errMsg := tlog.E(ctx).Err(errx).Msgf("Delete document tx (user id: %s, document id: %s) err (db transaction %v)",
				userId, documentId, errx)
			errx.AttachErrMsg(errMsg)

			return errx
		}

		return nil
	})
	if err != nil {
		errMsg := tlog.E(ctx).Err(err).Msgf("Delete document (user id: %s, document id: %s) err (db transaction %v)",
			userId, documentId, err)

		if errx != nil {
			errx.AttachErrMsg(errMsg)
		} else {
			errx = terror.NewTerror(ctx, err, constant.ErrorCodeMysqlServerAbnormal, errMsg)
		}

		return errx
	}

	return nil
}
