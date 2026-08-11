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

	if documentDB.CurVersionId != "" {
		documentVersionDB, errx := dbmodel.FindDocumentVersion(ctx, documentDB.CurVersionId)
		if errx != nil {
			errMsg := tlog.E(ctx).Err(errx).Msgf("Get document (user id: %s, document id: %s, version id: %s) err (db find document version %v)",
				userId, documentId, documentDB.CurVersionId, errx)
			errx.AttachErrMsg(errMsg)

			return nil, errx
		}

		if documentVersionDB != nil {
			documentData.CurrentVersion = &data.DocumentVersionData{
				VersionId: documentVersionDB.Id,

				DocumentId: documentVersionDB.DocumentId,
				VersionNo:  documentVersionDB.VersionNo,

				FileObjectId: documentVersionDB.FileObjectId,

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

			fileObjectDB, errx := dbmodel.FindFileObject(ctx, documentVersionDB.FileObjectId)
			if errx != nil {
				errMsg := tlog.E(ctx).Err(errx).Msgf("Get document (user id: %s, document id: %s, version id: %s, file object id: %s) err (db find file object %v)",
					userId, documentId, documentDB.CurVersionId, documentVersionDB.FileObjectId, errx)
				errx.AttachErrMsg(errMsg)

				return nil, errx
			}

			if fileObjectDB != nil {
				documentData.FileObject = &data.FileObjectData{
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
			}
		}
	}

	ingestJobDB, errx := dbmodel.FindLatestIngestJob(ctx, documentId, documentDB.CurVersionId)
	if errx != nil {
		errMsg := tlog.E(ctx).Err(errx).Msgf("Get document (user id: %s, document id: %s, version id: %s) err (db find latest ingest job %v)",
			userId, documentId, documentDB.CurVersionId, errx)
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

	return &data.GetDocumentRespData{
		DocumentData: documentData,
	}, nil
}

func CreateDocument(ctx context.Context, userId, knowledgeBaseId, chatSessionId string, scopeType, sourceType int,
	title, summary string, tags []string, ownerId, langCode string, parseStrategy int, fileHeader *multipart.FileHeader) (*data.CreateDocumentRespData, *terror.Terror) {
	if scopeType == dbmodel.DocumentScopeTypeKnowledge {
		knowledgeBaseDB, errx := dbmodel.FindKnowledgeBase(ctx, knowledgeBaseId)
		if errx != nil {
			errMsg := tlog.E(ctx).Err(errx).Msgf("Create document (user id: %s, knowledge base id: %s, chat session id: %s, scope type: %d, source type: %d, title: %s, summary: %s, tags: %v, owner id: %s, lang code: %s, parse strategy: %d, file name: %s, file size: %d) err (db find knowledge base %v)",
				userId, knowledgeBaseId, chatSessionId, scopeType, sourceType, title, summary, tags, ownerId, langCode, parseStrategy, fileHeader.Filename, fileHeader.Size, errx)
			errx.AttachErrMsg(errMsg)

			return nil, errx
		}

		if knowledgeBaseDB == nil {
			errMsg := tlog.E(ctx).Msgf("Create document (user id: %s, knowledge base id: %s, chat session id: %s, scope type: %d, source type: %d, title: %s, summary: %s, tags: %v, owner id: %s, lang code: %s, parse strategy: %d, file name: %s, file size: %d) err (knowledge base not found)",
				userId, knowledgeBaseId, chatSessionId, scopeType, sourceType, title, summary, tags, ownerId, langCode, parseStrategy, fileHeader.Filename, fileHeader.Size)

			errx := terror.NewTerror(ctx, terror.ErrParamInvalid("knowledge_base_id"), constant.ErrorCodeKnowledgeBaseNotFound, errMsg)

			return nil, errx
		}
	}

	originFileName := strings.TrimSpace(fileHeader.Filename)
	safeFileName := lib.SanitizeAwsS3FileName(originFileName)

	if title == "" {
		title = strings.TrimSpace(strings.TrimSuffix(originFileName, path.Ext(originFileName)))
	}
	if title == "" {
		title = safeFileName
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

	ownerId = strings.TrimSpace(ownerId)
	if ownerId == "" {
		ownerId = constant.DefaultAnonymousOwnerId
	}

	langCode = strings.TrimSpace(langCode)
	if langCode == "" {
		langCode = constant.DefaultDocumentLangCode
	}

	ocrStatus := dbmodel.DocumentVersionOcrStatusNotRequired
	if parseStrategy == dbmodel.DocumentParseStrategyOCR || parseStrategy == dbmodel.DocumentParseStrategyTikaOCR {
		ocrStatus = dbmodel.DocumentVersionOcrStatusPending
	}

	var dbErrx *terror.Terror
	var documentDB *dbmodel.Document
	var documentVersionDB *dbmodel.DocumentVersion
	var fileObjectDB *dbmodel.FileObject

	documentId := ""
	versionId := ""
	fileObjectId := ""

	err := dbmodel.DB(ctx).Transaction(func(tx *gorm.DB) error {
		documentDB, dbErrx = dbmodel.CreateDocument(ctx, tx, knowledgeBaseId, chatSessionId, scopeType, sourceType,
			title, summary, tags, ownerId, langCode, 0, "", dbmodel.DocumentProcessStatusWaiting, dbmodel.DocumentStatusNormal)
		if dbErrx != nil {
			return dbErrx
		}
		documentId = documentDB.Id

		fileObjectDB, dbErrx = dbmodel.CreateFileObject(ctx, tx, bucketName, "", originFileName, mimeType, fileExt,
			uint64(fileHeader.Size), "", dbmodel.FileObjectStorageProviderSeaweedFS)
		if dbErrx != nil {
			return dbErrx
		}
		fileObjectId = fileObjectDB.Id

		documentVersionDB, dbErrx = dbmodel.CreateDocumentVersion(ctx, tx, documentId, 1, fileObjectId, parseStrategy,
			dbmodel.DocumentParserTypeUnknown, "", 0, 0, 0, dbmodel.DocumentVersionParseStatusPending, "", ocrStatus, "")
		if dbErrx != nil {
			return dbErrx
		}
		versionId = documentVersionDB.Id

		dbErrx = dbmodel.UpdateDocumentCurrentVersion(ctx, tx, documentId, documentVersionDB.VersionNo, versionId)
		if dbErrx != nil {
			return dbErrx
		}

		return nil
	})
	if err != nil {
		if dbErrx != nil {
			errMsg := tlog.E(ctx).Err(dbErrx).Msgf("Create document (user id: %s, knowledge base id: %s, chat session id: %s, scope type: %d, source type: %d, title: %s, summary: %s, tags: %v, owner id: %s, lang code: %s, parse strategy: %d, document id: %s, version id: %s, file object id: %s, bucket: %s, file name: %s, file size: %d) err (db prepare document records transaction %v)",
				userId, knowledgeBaseId, chatSessionId, scopeType, sourceType, title, summary, tags, ownerId, langCode, parseStrategy, documentId, versionId, fileObjectId, bucketName, fileHeader.Filename, fileHeader.Size, dbErrx)
			dbErrx.AttachErrMsg(errMsg)

			return nil, dbErrx
		}

		errMsg := tlog.E(ctx).Err(err).Msgf("Create document (user id: %s, knowledge base id: %s, chat session id: %s, scope type: %d, source type: %d, title: %s, summary: %s, tags: %v, owner id: %s, lang code: %s, parse strategy: %d, document id: %s, version id: %s, file object id: %s, bucket: %s, file name: %s, file size: %d) err (db prepare document records transaction %v)",
			userId, knowledgeBaseId, chatSessionId, scopeType, sourceType, title, summary, tags, ownerId, langCode, parseStrategy, documentId, versionId, fileObjectId, bucketName, fileHeader.Filename, fileHeader.Size, err)

		errx := terror.NewTerror(ctx, err, constant.ErrorCodeMysqlServerAbnormal, errMsg)

		return nil, errx
	}

	objectKey := lib.BuildAwsS3RawObjectKey(knowledgeBaseId, chatSessionId, documentId, versionId, fileObjectId, safeFileName, time.Now())

	sha256Value, errx := lib.UploadAwsS3File(ctx, bucketName, objectKey, safeFileName, fileHeader)
	if errx != nil {
		errMsg := tlog.E(ctx).Err(errx).Msgf("Create document (user id: %s, knowledge base id: %s, chat session id: %s, scope type: %d, source type: %d, title: %s, summary: %s, tags: %v, owner id: %s, lang code: %s, parse strategy: %d, document id: %s, version id: %s, file object id: %s, bucket: %s, object key: %s, file name: %s, file size: %d) err (lib upload aws s3 file %v)",
			userId, knowledgeBaseId, chatSessionId, scopeType, sourceType, title, summary, tags, ownerId, langCode, parseStrategy, documentId, versionId, fileObjectId, bucketName, objectKey, fileHeader.Filename, fileHeader.Size, errx)
		errx.AttachErrMsg(errMsg)

		dbErrx = nil
		err = dbmodel.DB(ctx).Transaction(func(tx *gorm.DB) error {
			dbErrx = dbmodel.UpdateDocumentProcessStatus(ctx, tx, documentId, dbmodel.DocumentProcessStatusFailed)
			if dbErrx != nil {
				return dbErrx
			}

			return nil
		})
		if err != nil {
			if dbErrx != nil {
				errMsg = tlog.E(ctx).Err(dbErrx).Msgf("Create document (user id: %s, knowledge base id: %s, chat session id: %s, scope type: %d, source type: %d, title: %s, summary: %s, tags: %v, owner id: %s, lang code: %s, parse strategy: %d, document id: %s, version id: %s, file object id: %s, bucket: %s, object key: %s, file name: %s, file size: %d) err (db mark document upload failed %v)",
					userId, knowledgeBaseId, chatSessionId, scopeType, sourceType, title, summary, tags, ownerId, langCode, parseStrategy, documentId, versionId, fileObjectId, bucketName, objectKey, fileHeader.Filename, fileHeader.Size, dbErrx)
				dbErrx.AttachErrMsg(errMsg)
			} else {
				_ = tlog.E(ctx).Err(err).Msgf("Create document (user id: %s, knowledge base id: %s, chat session id: %s, scope type: %d, source type: %d, title: %s, summary: %s, tags: %v, owner id: %s, lang code: %s, parse strategy: %d, document id: %s, version id: %s, file object id: %s, bucket: %s, object key: %s, file name: %s, file size: %d) err (db mark document upload failed %v)",
					userId, knowledgeBaseId, chatSessionId, scopeType, sourceType, title, summary, tags, ownerId, langCode, parseStrategy, documentId, versionId, fileObjectId, bucketName, objectKey, fileHeader.Filename, fileHeader.Size, err)
			}
		}

		return nil, errx
	}

	jobPayloadBytes, _ := json.Marshal(map[string]string{
		"bucket_name": bucketName,
		"object_key":  objectKey,
	})

	jobPayload := string(jobPayloadBytes)

	var ingestJobDB *dbmodel.IngestJob

	err = dbmodel.DB(ctx).Transaction(func(tx *gorm.DB) error {
		dbErrx = dbmodel.UpdateFileObjectStorageInfo(ctx, tx, fileObjectId, objectKey, sha256Value)
		if dbErrx != nil {
			return dbErrx
		}

		dbErrx = dbmodel.UpdateDocumentVersionContentSha256(ctx, tx, versionId, sha256Value)
		if dbErrx != nil {
			return dbErrx
		}

		dbErrx = dbmodel.UpdateDocumentProcessStatus(ctx, tx, documentId, dbmodel.DocumentProcessStatusUploaded)
		if dbErrx != nil {
			return dbErrx
		}

		ingestJobDB, dbErrx = dbmodel.CreateIngestJob(ctx, tx, documentId, versionId, dbmodel.IngestJobTypeParse,
			dbmodel.IngestJobStatusPending, 0, "", "", jobPayload)
		if dbErrx != nil {
			return dbErrx
		}

		return nil
	})
	if err != nil {
		jobId := ""
		if ingestJobDB != nil {
			jobId = ingestJobDB.Id
		}

		if dbErrx != nil {
			errMsg := tlog.E(ctx).Err(dbErrx).Msgf("Create document (user id: %s, knowledge base id: %s, chat session id: %s, scope type: %d, source type: %d, title: %s, summary: %s, tags: %v, owner id: %s, lang code: %s, parse strategy: %d, document id: %s, version id: %s, file object id: %s, job id: %s, bucket: %s, object key: %s, sha256: %s, file name: %s, file size: %d) err (db complete document upload transaction %v)",
				userId, knowledgeBaseId, chatSessionId, scopeType, sourceType, title, summary, tags, ownerId, langCode, parseStrategy, documentId, versionId, fileObjectId, jobId, bucketName, objectKey, sha256Value, fileHeader.Filename, fileHeader.Size, dbErrx)
			dbErrx.AttachErrMsg(errMsg)

			return nil, dbErrx
		}

		errMsg := tlog.E(ctx).Err(err).Msgf("Create document (user id: %s, knowledge base id: %s, chat session id: %s, scope type: %d, source type: %d, title: %s, summary: %s, tags: %v, owner id: %s, lang code: %s, parse strategy: %d, document id: %s, version id: %s, file object id: %s, job id: %s, bucket: %s, object key: %s, sha256: %s, file name: %s, file size: %d) err (db complete document upload transaction %v)",
			userId, knowledgeBaseId, chatSessionId, scopeType, sourceType, title, summary, tags, ownerId, langCode, parseStrategy, documentId, versionId, fileObjectId, jobId, bucketName, objectKey, sha256Value, fileHeader.Filename, fileHeader.Size, err)

		errx := terror.NewTerror(ctx, err, constant.ErrorCodeMysqlServerAbnormal, errMsg)

		return nil, errx
	}

	return &data.CreateDocumentRespData{
		DocumentId: documentId,

		VersionId:    versionId,
		FileObjectId: fileObjectId,
		JobId:        ingestJobDB.Id,

		BucketName: bucketName,
		ObjectKey:  objectKey,

		ProcessStatus: dbmodel.DocumentProcessStatusUploaded,
		Status:        dbmodel.DocumentStatusNormal,
	}, nil
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

	var dbErrx *terror.Terror

	err := dbmodel.DB(ctx).Transaction(func(tx *gorm.DB) error {
		dbErrx = dbmodel.DisableDocument(ctx, tx, documentId)
		if dbErrx != nil {
			return dbErrx
		}

		dbErrx = dbmodel.DeleteDocument(ctx, tx, documentId)
		if dbErrx != nil {
			return dbErrx
		}

		return nil
	})
	if err != nil {
		if dbErrx != nil {
			errMsg := tlog.E(ctx).Err(dbErrx).Msgf("Delete document (user id: %s, document id: %s) err (db transaction %v)",
				userId, documentId, dbErrx)
			dbErrx.AttachErrMsg(errMsg)

			return dbErrx
		}

		errMsg := tlog.E(ctx).Err(err).Msgf("Delete document (user id: %s, document id: %s) err (db transaction %v)",
			userId, documentId, err)

		errx := terror.NewTerror(ctx, err, constant.ErrorCodeMysqlServerAbnormal, errMsg)

		return errx
	}

	return nil
}
