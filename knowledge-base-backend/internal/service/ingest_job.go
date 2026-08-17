package service

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/choveylee/terror"
	"github.com/choveylee/tlog"
	"github.com/choveylee/tutil"
	"gorm.io/gorm"

	constant "dev.choveylee.top/knowledge-base-backend/internal/const"
	"dev.choveylee.top/knowledge-base-backend/internal/lib"
	dbmodel "dev.choveylee.top/knowledge-base-backend/internal/model/mysql"
	redmodel "dev.choveylee.top/knowledge-base-backend/internal/model/redis"
)

func StartIngestWorker(ctx context.Context) *terror.Terror {
	if !ingestWorkerEnabled {
		tlog.I(ctx).Msgf("Ingest worker disabled (worker name: %s, block timeout: %s, error sleep: %s)",
			ingestJobWorkerName, ingestWorkerBlockTimeout, ingestWorkerErrorSleep)

		return nil
	}

	go func() {
		tlog.I(ctx).Msgf("Ingest worker started (worker name: %s, block timeout: %s, error sleep: %s)",
			ingestJobWorkerName, ingestWorkerBlockTimeout, ingestWorkerErrorSleep)

		for {
			select {
			case <-ctx.Done():
				tlog.I(ctx).Msgf("Ingest worker stopped (worker name: %s)",
					ingestJobWorkerName)

				return
			default:
			}

			func() {
				defer func() {
					if panicValue := recover(); panicValue != nil {
						tlog.E(ctx).Msgf("Run ingest worker (worker name: %s, block timeout: %s, error sleep: %s) err (panic recovered %v)",
							ingestJobWorkerName, ingestWorkerBlockTimeout, ingestWorkerErrorSleep, panicValue)

						time.Sleep(ingestWorkerErrorSleep)
					}
				}()

				// Redis blocking pop is the main low-latency queue path.
				jobId, errx := redmodel.BlockPopParseIngestJobId(ctx, ingestWorkerBlockTimeout)
				if errx != nil {
					errMsg := tlog.E(ctx).Err(errx).Msgf("Run ingest worker (worker name: %s, block timeout: %s, error sleep: %s) err (redis block pop parse ingest job id %v)",
						ingestJobWorkerName, ingestWorkerBlockTimeout, ingestWorkerErrorSleep, errx)
					errx.AttachErrMsg(errMsg)

					time.Sleep(ingestWorkerErrorSleep)

					return
				}

				if jobId == "" {
					return
				}

				errx = processParseIngestJobById(ctx, ingestJobWorkerName, jobId)
				if errx != nil {
					errMsg := tlog.E(ctx).Err(errx).Msgf("Run ingest worker (worker name: %s, block timeout: %s, error sleep: %s, job id: %s) err (process parse ingest job by id %v)",
						ingestJobWorkerName, ingestWorkerBlockTimeout, ingestWorkerErrorSleep, jobId, errx)
					errx.AttachErrMsg(errMsg)

					time.Sleep(ingestWorkerErrorSleep)
				}
			}()
		}
	}()

	return nil
}

func RunPendingIngestJob(ctx context.Context, workerName string, limit int) *terror.Terror {
	if limit <= 0 {
		return nil
	}

	ingestJobsDB, errx := dbmodel.FindPendingIngestJobs(ctx, dbmodel.IngestJobTypeParse, limit)
	if errx != nil {
		errMsg := tlog.E(ctx).Err(errx).Msgf("Run pending ingest job (worker name: %s, limit: %d) err (db find pending ingest jobs %v)",
			workerName, limit, errx)
		errx.AttachErrMsg(errMsg)

		return errx
	}

	for _, ingestJobDB := range ingestJobsDB {
		jobId := ingestJobDB.Id

		errx = processParseIngestJobById(ctx, workerName, jobId)
		if errx != nil {
			errMsg := tlog.E(ctx).Err(errx).Msgf("Run pending ingest job (worker name: %s, limit: %d, job id: %s) err (process parse ingest job by id %v)",
				workerName, limit, jobId, errx)
			errx.AttachErrMsg(errMsg)

			return errx
		}
	}

	return nil
}

func processParseIngestJobById(ctx context.Context, workerName string, jobId string) *terror.Terror {
	jobId = strings.TrimSpace(jobId)
	if jobId == "" {
		errMsg := tlog.E(ctx).Msgf("Process parse ingest job by id (worker name: %s) err (job id empty)",
			workerName)

		errx := terror.NewTerror(ctx, terror.ErrParamInvalid("job_id"), constant.ErrorCodeRequestParamInvalid, errMsg)

		return errx
	}

	ingestJobDB, errx := dbmodel.FindIngestJob(ctx, jobId)
	if errx != nil {
		errMsg := tlog.E(ctx).Err(errx).Msgf("Process parse ingest job by id (worker name: %s, job id: %s) err (db find ingest job %v)",
			workerName, jobId, errx)
		errx.AttachErrMsg(errMsg)

		return errx
	}

	if ingestJobDB == nil || ingestJobDB.JobStatus != dbmodel.IngestJobStatusPending {
		return nil
	}

	// Claim the job before parsing so multiple workers do not process it twice.
	started, errx := dbmodel.TryStartIngestJob(ctx, jobId, workerName)
	if errx != nil {
		errMsg := tlog.E(ctx).Err(errx).Msgf("Process parse ingest job by id (worker name: %s, job id: %s) err (db try start ingest job %v)",
			workerName, jobId, errx)
		errx.AttachErrMsg(errMsg)

		return errx
	}

	if !started {
		return nil
	}

	ingestJobDB.JobStatus = dbmodel.IngestJobStatusProcessing
	ingestJobDB.WorkerName = workerName

	errx = processParseIngestJob(ctx, workerName, ingestJobDB)
	if errx != nil {
		errMsg := tlog.E(ctx).Err(errx).Msgf("Process parse ingest job by id (worker name: %s, job id: %s) err (process parse ingest job %v)",
			workerName, jobId, errx)
		errx.AttachErrMsg(errMsg)

		return errx
	}

	return nil
}

func processParseIngestJob(ctx context.Context, workerName string, ingestJobDB *dbmodel.IngestJob) *terror.Terror {
	if ingestJobDB == nil {
		errMsg := tlog.E(ctx).Msgf("Process parse ingest job (worker name: %s) err (ingest job nil)",
			workerName)

		errx := terror.NewTerror(ctx, terror.ErrParamInvalid("ingest_job"), constant.ErrorCodeRequestParamInvalid, errMsg)

		return errx
	}

	jobId := ingestJobDB.Id
	documentId := ingestJobDB.DocumentId
	versionId := ingestJobDB.VersionId
	retryCount := ingestJobDB.RetryCount
	payload := ingestJobDB.Payload

	var errx *terror.Terror

	err := dbmodel.DB(ctx).Transaction(func(tx *gorm.DB) error {
		errx = dbmodel.UpdateDocumentProcessStatusTx(ctx, tx, documentId, dbmodel.DocumentProcessStatusProcessing)
		if errx != nil {
			errMsg := tlog.E(ctx).Err(errx).Msgf("Process parse ingest job tx (worker name: %s, job id: %s, document id: %s, version id: %s, retry count: %d, payload: %s) err (db update document process status %v)",
				workerName, jobId, documentId, versionId, retryCount, payload, errx)
			errx.AttachErrMsg(errMsg)

			return errx
		}

		errx = dbmodel.UpdateDocumentVersionParseStatusTx(ctx, tx, versionId, dbmodel.DocumentVersionParseStatusProcessing, "")
		if errx != nil {
			errMsg := tlog.E(ctx).Err(errx).Msgf("Process parse ingest job tx (worker name: %s, job id: %s, document id: %s, version id: %s, retry count: %d, payload: %s) err (db update document version parse status %v)",
				workerName, jobId, documentId, versionId, retryCount, payload, errx)
			errx.AttachErrMsg(errMsg)

			return errx
		}

		return nil
	})
	if err != nil {
		errMsg := tlog.E(ctx).Err(err).Msgf("Process parse ingest job (worker name: %s, job id: %s, document id: %s, version id: %s, retry count: %d, payload: %s) err (db mark parse processing transaction %v)",
			workerName, jobId, documentId, versionId, retryCount, payload, err)

		if errx != nil {
			errx.AttachErrMsg(errMsg)
		} else {
			errx = terror.NewTerror(ctx, err, constant.ErrorCodeMysqlServerAbnormal, errMsg)
		}

		return errx
	}

	documentVersionDB, errx := dbmodel.FindDocumentVersion(ctx, versionId)
	if errx != nil {
		errMsg := tlog.E(ctx).Err(errx).Msgf("Process parse ingest job (worker name: %s, job id: %s, document id: %s, version id: %s, retry count: %d, payload: %s) err (db find document version %v)",
			workerName, jobId, documentId, versionId, retryCount, payload, errx)
		errx.AttachErrMsg(errMsg)

		_ = markParseIngestJobFailed(ctx, workerName, ingestJobDB, errx.Error(), dbmodel.DocumentVersionOcrStatusNotRequired, "")

		return errx
	}

	if documentVersionDB == nil {
		errMsg := tlog.E(ctx).Msgf("Process parse ingest job (worker name: %s, job id: %s, document id: %s, version id: %s, retry count: %d, payload: %s) err (document version not found)",
			workerName, jobId, documentId, versionId, retryCount, payload)

		errx := terror.NewTerror(ctx, terror.ErrParamInvalid("version_id"), constant.ErrorCodeDocumentNotFound, errMsg)

		_ = markParseIngestJobFailed(ctx, workerName, ingestJobDB, errMsg, dbmodel.DocumentVersionOcrStatusNotRequired, "")

		return errx
	}

	parseStrategy := documentVersionDB.ParseStrategy
	if parseStrategy == dbmodel.DocumentParseStrategyOCR || parseStrategy == dbmodel.DocumentParseStrategyTikaOCR {
		err := errors.New("ocr parse strategy is not implemented yet")
		errMsg := tlog.E(ctx).Err(err).Msgf("Process parse ingest job (worker name: %s, job id: %s, document id: %s, version id: %s, retry count: %d, payload: %s, parse strategy: %d) err (ocr parser unsupported %v)",
			workerName, jobId, documentId, versionId, retryCount, payload, parseStrategy, err)

		errx := terror.NewTerror(ctx, err, constant.ErrorCodeDocumentParseFailed, errMsg)

		_ = markParseIngestJobFailed(ctx, workerName, ingestJobDB, errMsg, dbmodel.DocumentVersionOcrStatusFailed, errMsg)

		return errx
	}

	fileObjectId := documentVersionDB.FileObjectId

	fileObjectDB, errx := dbmodel.FindFileObject(ctx, fileObjectId)
	if errx != nil {
		errMsg := tlog.E(ctx).Err(errx).Msgf("Process parse ingest job (worker name: %s, job id: %s, document id: %s, version id: %s, retry count: %d, payload: %s, file object id: %s) err (db find file object %v)",
			workerName, jobId, documentId, versionId, retryCount, payload, fileObjectId, errx)
		errx.AttachErrMsg(errMsg)

		_ = markParseIngestJobFailed(ctx, workerName, ingestJobDB, errx.Error(), dbmodel.DocumentVersionOcrStatusNotRequired, "")

		return errx
	}

	if fileObjectDB == nil {
		errMsg := tlog.E(ctx).Msgf("Process parse ingest job (worker name: %s, job id: %s, document id: %s, version id: %s, retry count: %d, payload: %s, file object id: %s) err (file object not found)",
			workerName, jobId, documentId, versionId, retryCount, payload, fileObjectId)

		errx := terror.NewTerror(ctx, terror.ErrParamInvalid("file_object_id"), constant.ErrorCodeDocumentFileInvalid, errMsg)

		_ = markParseIngestJobFailed(ctx, workerName, ingestJobDB, errMsg, dbmodel.DocumentVersionOcrStatusNotRequired, "")

		return errx
	}

	bucketName := fileObjectDB.BucketName
	objectKey := fileObjectDB.ObjectKey

	if strings.TrimSpace(payload) != "" {
		var payloadData struct {
			BucketName string `json:"bucket_name"`
			ObjectKey  string `json:"object_key"`
		}

		err := json.Unmarshal([]byte(payload), &payloadData)
		if err != nil {
			errMsg := tlog.E(ctx).Err(err).Msgf("Process parse ingest job (worker name: %s, job id: %s, document id: %s, version id: %s, retry count: %d, file object id: %s, payload: %s) err (payload unmarshal %v)",
				workerName, jobId, documentId, versionId, retryCount, fileObjectId, payload, err)

			errx := terror.NewTerror(ctx, err, constant.ErrorCodeDocumentParseFailed, errMsg)

			_ = markParseIngestJobFailed(ctx, workerName, ingestJobDB, errMsg, dbmodel.DocumentVersionOcrStatusNotRequired, "")

			return errx
		}

		payloadData.BucketName = strings.TrimSpace(payloadData.BucketName)
		if payloadData.BucketName != "" {
			bucketName = payloadData.BucketName
		}

		payloadData.ObjectKey = strings.TrimSpace(payloadData.ObjectKey)
		if payloadData.ObjectKey != "" {
			objectKey = payloadData.ObjectKey
		}
	}

	fileBytes, errx := lib.DownloadAwsS3File(ctx, bucketName, objectKey)
	if errx != nil {
		errMsg := tlog.E(ctx).Err(errx).Msgf("Process parse ingest job (worker name: %s, job id: %s, document id: %s, version id: %s, file object id: %s, bucket: %s, object key: %s) err (lib download aws s3 file %v)",
			workerName, jobId, documentId, versionId, fileObjectId, bucketName, objectKey, errx)
		errx.AttachErrMsg(errMsg)

		_ = markParseIngestJobFailed(ctx, workerName, ingestJobDB, errx.Error(), dbmodel.DocumentVersionOcrStatusNotRequired, "")

		return errx
	}

	plainText, errx := lib.ParseTikaFile(ctx, fileObjectDB.FileName, fileObjectDB.MimeType, fileBytes)
	if errx != nil {
		errMsg := tlog.E(ctx).Err(errx).Msgf("Process parse ingest job (worker name: %s, job id: %s, document id: %s, version id: %s, file object id: %s, bucket: %s, object key: %s, file name: %s, mime type: %s, file size: %d) err (lib parse tika file %v)",
			workerName, jobId, documentId, versionId, fileObjectId, bucketName, objectKey, fileObjectDB.FileName, fileObjectDB.MimeType, len(fileBytes), errx)
		errx.AttachErrMsg(errMsg)

		_ = markParseIngestJobFailed(ctx, workerName, ingestJobDB, errx.Error(), dbmodel.DocumentVersionOcrStatusNotRequired, "")

		return errx
	}

	plainText = lib.NormalizeDocumentText(plainText)
	if plainText == "" {
		err := errors.New("tika parsed text empty")
		errMsg := tlog.E(ctx).Err(err).Msgf("Process parse ingest job (worker name: %s, job id: %s, document id: %s, version id: %s, file object id: %s, bucket: %s, object key: %s, file name: %s, mime type: %s, file size: %d) err (parsed text empty %v)",
			workerName, jobId, documentId, versionId, fileObjectId, bucketName, objectKey, fileObjectDB.FileName, fileObjectDB.MimeType, len(fileBytes), err)

		errx := terror.NewTerror(ctx, err, constant.ErrorCodeDocumentParseFailed, errMsg)

		_ = markParseIngestJobFailed(ctx, workerName, ingestJobDB, errMsg, dbmodel.DocumentVersionOcrStatusNotRequired, "")

		return errx
	}

	chunks, errx := lib.SplitEinoText(ctx, documentId, plainText)
	if errx != nil {
		errMsg := tlog.E(ctx).Err(errx).Msgf("Process parse ingest job (worker name: %s, job id: %s, document id: %s, version id: %s, file object id: %s, bucket: %s, object key: %s, parsed text len: %d) err (lib split eino text %v)",
			workerName, jobId, documentId, versionId, fileObjectId, bucketName, objectKey, lib.EstimateDocumentTokenCount(plainText), errx)
		errx.AttachErrMsg(errMsg)

		_ = markParseIngestJobFailed(ctx, workerName, ingestJobDB, errx.Error(), dbmodel.DocumentVersionOcrStatusNotRequired, "")

		return errx
	}

	if len(chunks) == 0 {
		err := errors.New("document chunks empty")
		errMsg := tlog.E(ctx).Err(err).Msgf("Process parse ingest job (worker name: %s, job id: %s, document id: %s, version id: %s, file object id: %s, bucket: %s, object key: %s, parsed text len: %d) err (document chunks empty %v)",
			workerName, jobId, documentId, versionId, fileObjectId, bucketName, objectKey, lib.EstimateDocumentTokenCount(plainText), err)

		errx := terror.NewTerror(ctx, err, constant.ErrorCodeDocumentParseFailed, errMsg)

		_ = markParseIngestJobFailed(ctx, workerName, ingestJobDB, errMsg, dbmodel.DocumentVersionOcrStatusNotRequired, "")

		return errx
	}

	documentChunksDB := make([]*dbmodel.DocumentChunk, 0, len(chunks))
	totalTokenCount := uint(0)
	charStart := uint(0)

	for index, chunk := range chunks {
		chunkText := strings.TrimSpace(chunk.Content)
		tokenCount := lib.EstimateDocumentTokenCount(chunkText)
		charEnd := charStart + tokenCount

		headingItems := make([]string, 0)
		for _, key := range []string{"h1", "h2", "h3"} {
			value, ok := chunk.MetaData[key].(string)
			if ok && strings.TrimSpace(value) != "" {
				headingItems = append(headingItems, strings.TrimSpace(value))
			}
		}

		headingPath := strings.Join(headingItems, " / ")
		if len(headingPath) > dbmodel.DocumentChunkHeadingPathLen {
			headingPath = headingPath[:dbmodel.DocumentChunkHeadingPathLen]
		}

		documentChunkDB := &dbmodel.DocumentChunk{
			Id: tutil.NewOid().String(),

			DocumentId: documentId,
			VersionId:  versionId,
			ChunkNo:    uint(index + 1),

			ChunkType:   dbmodel.DocumentChunkTypeParagraph,
			HeadingPath: headingPath,

			ChunkText: chunkText,
			TextHash:  lib.Sha256Text(chunkText),

			CharStart: charStart,
			CharEnd:   charEnd,

			TokenCount: tokenCount,

			IndexStatus: dbmodel.DocumentChunkIndexStatusPending,
		}

		documentChunksDB = append(documentChunksDB, documentChunkDB)

		totalTokenCount += tokenCount
		charStart = charEnd
	}

	errx = nil
	err = dbmodel.DB(ctx).Transaction(func(tx *gorm.DB) error {
		errx = dbmodel.HardDeleteDocumentChunksByVersionTx(ctx, tx, versionId)
		if errx != nil {
			errMsg := tlog.E(ctx).Err(errx).Msgf("Process parse ingest job tx (worker name: %s, job id: %s, document id: %s, version id: %s, chunk count: %d) err (db hard delete document chunks %v)",
				workerName, jobId, documentId, versionId, len(documentChunksDB), errx)
			errx.AttachErrMsg(errMsg)

			return errx
		}

		errx = dbmodel.BatchCreateDocumentChunksTx(ctx, tx, documentChunksDB)
		if errx != nil {
			errMsg := tlog.E(ctx).Err(errx).Msgf("Process parse ingest job tx (worker name: %s, job id: %s, document id: %s, version id: %s, chunk count: %d) err (db batch create document chunks %v)",
				workerName, jobId, documentId, versionId, len(documentChunksDB), errx)
			errx.AttachErrMsg(errMsg)

			return errx
		}

		errx = dbmodel.UpdateDocumentVersionParseResultTx(ctx, tx, versionId, dbmodel.DocumentParserTypeTika, plainText, 0,
			totalTokenCount, uint(len(documentChunksDB)), dbmodel.DocumentVersionParseStatusSuccess, "", dbmodel.DocumentVersionOcrStatusNotRequired, "")
		if errx != nil {
			errMsg := tlog.E(ctx).Err(errx).Msgf("Process parse ingest job tx (worker name: %s, job id: %s, document id: %s, version id: %s, token count: %d, chunk count: %d) err (db update document version parse result %v)",
				workerName, jobId, documentId, versionId, totalTokenCount, len(documentChunksDB), errx)
			errx.AttachErrMsg(errMsg)

			return errx
		}

		errx = dbmodel.UpdateDocumentProcessStatusTx(ctx, tx, documentId, dbmodel.DocumentProcessStatusProcessed)
		if errx != nil {
			errMsg := tlog.E(ctx).Err(errx).Msgf("Process parse ingest job tx (worker name: %s, job id: %s, document id: %s, version id: %s) err (db update document process status %v)",
				workerName, jobId, documentId, versionId, errx)
			errx.AttachErrMsg(errMsg)

			return errx
		}

		errx = dbmodel.UpdateIngestJobFinishedTx(ctx, tx, jobId)
		if errx != nil {
			errMsg := tlog.E(ctx).Err(errx).Msgf("Process parse ingest job tx (worker name: %s, job id: %s, document id: %s, version id: %s) err (db update ingest job finished %v)",
				workerName, jobId, documentId, versionId, errx)
			errx.AttachErrMsg(errMsg)

			return errx
		}

		return nil
	})
	if err != nil {
		errMsg := tlog.E(ctx).Err(err).Msgf("Process parse ingest job (worker name: %s, job id: %s, document id: %s, version id: %s, retry count: %d, chunk count: %d) err (db complete parse job transaction %v)",
			workerName, jobId, documentId, versionId, retryCount, len(documentChunksDB), err)

		if errx != nil {
			errx.AttachErrMsg(errMsg)
		} else {
			errx = terror.NewTerror(ctx, err, constant.ErrorCodeMysqlServerAbnormal, errMsg)
		}

		return errx
	}

	return nil
}

func markParseIngestJobFailed(ctx context.Context, workerName string, ingestJobDB *dbmodel.IngestJob, parseError string, ocrStatus int, ocrError string) *terror.Terror {
	if ingestJobDB == nil {
		errMsg := tlog.E(ctx).Msgf("Mark parse ingest job failed (worker name: %s, parse error: %s, ocr status: %d, ocr error: %s) err (ingest job nil)",
			workerName, parseError, ocrStatus, ocrError)

		errx := terror.NewTerror(ctx, terror.ErrParamInvalid("ingest_job"), constant.ErrorCodeRequestParamInvalid, errMsg)

		return errx
	}

	jobId := ingestJobDB.Id
	documentId := ingestJobDB.DocumentId
	versionId := ingestJobDB.VersionId
	jobType := ingestJobDB.JobType
	jobStatus := ingestJobDB.JobStatus
	retryCount := ingestJobDB.RetryCount + 1
	payload := ingestJobDB.Payload

	var errx *terror.Terror

	err := dbmodel.DB(ctx).Transaction(func(tx *gorm.DB) error {
		errx = dbmodel.UpdateDocumentProcessStatusTx(ctx, tx, documentId, dbmodel.DocumentProcessStatusFailed)
		if errx != nil {
			errMsg := tlog.E(ctx).Err(errx).Msgf("Mark parse ingest job failed tx (worker name: %s, job id: %s, document id: %s, version id: %s, job type: %d, job status: %d, retry count: %d, payload: %s, parse error: %s, ocr status: %d, ocr error: %s) err (db update document process status %v)",
				workerName, jobId, documentId, versionId, jobType, jobStatus, retryCount, payload, parseError, ocrStatus, ocrError, errx)
			errx.AttachErrMsg(errMsg)

			return errx
		}

		errx = dbmodel.UpdateDocumentVersionParseResultTx(ctx, tx, versionId, dbmodel.DocumentParserTypeUnknown, "", 0, 0, 0,
			dbmodel.DocumentVersionParseStatusFailed, parseError, ocrStatus, ocrError)
		if errx != nil {
			errMsg := tlog.E(ctx).Err(errx).Msgf("Mark parse ingest job failed tx (worker name: %s, job id: %s, document id: %s, version id: %s, job type: %d, job status: %d, retry count: %d, payload: %s, parse error: %s, ocr status: %d, ocr error: %s) err (db update document version parse result %v)",
				workerName, jobId, documentId, versionId, jobType, jobStatus, retryCount, payload, parseError, ocrStatus, ocrError, errx)
			errx.AttachErrMsg(errMsg)

			return errx
		}

		errx = dbmodel.UpdateIngestJobFailedTx(ctx, tx, jobId, retryCount, parseError)
		if errx != nil {
			errMsg := tlog.E(ctx).Err(errx).Msgf("Mark parse ingest job failed tx (worker name: %s, job id: %s, document id: %s, version id: %s, job type: %d, job status: %d, retry count: %d, payload: %s, parse error: %s, ocr status: %d, ocr error: %s) err (db update ingest job failed %v)",
				workerName, jobId, documentId, versionId, jobType, jobStatus, retryCount, payload, parseError, ocrStatus, ocrError, errx)
			errx.AttachErrMsg(errMsg)

			return errx
		}

		return nil
	})
	if err != nil {
		errMsg := tlog.E(ctx).Err(err).Msgf("Mark parse ingest job failed (worker name: %s, job id: %s, document id: %s, version id: %s, job type: %d, job status: %d, retry count: %d, payload: %s, parse error: %s, ocr status: %d, ocr error: %s) err (db mark parse failed transaction %v)",
			workerName, jobId, documentId, versionId, jobType, jobStatus, retryCount, payload, parseError, ocrStatus, ocrError, err)

		if errx != nil {
			errx.AttachErrMsg(errMsg)
		} else {
			errx = terror.NewTerror(ctx, err, constant.ErrorCodeMysqlServerAbnormal, errMsg)
		}

		return errx
	}

	return nil
}
