package handler

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/choveylee/tlog"
	"github.com/gin-gonic/gin"

	constant "dev.choveylee.top/knowledge-base-backend/internal/const"
	"dev.choveylee.top/knowledge-base-backend/internal/data"
	dbmodel "dev.choveylee.top/knowledge-base-backend/internal/model/mysql"
	"dev.choveylee.top/knowledge-base-backend/internal/service"
)

func HandleListDocuments(c *gin.Context) {
	ctx := c.Request.Context()

	userId := strings.TrimSpace(c.Request.Header.Get("user_id"))

	knowledgeBaseId := strings.TrimSpace(c.Query("knowledge_base_id"))

	keyword := strings.TrimSpace(c.Query("keyword"))

	scopeType := -1

	srcScopeType := strings.TrimSpace(c.Query("scope_type"))
	if srcScopeType != "" {
		desScopeType, err := strconv.Atoi(srcScopeType)
		if err != nil {
			errMsg := tlog.E(ctx).Err(err).Msgf("Handle list documents (scope type: %s) err (strconv atoi %v)",
				srcScopeType, err)

			SendFailResponse(c, constant.ErrorCodeRequestParamInvalid, errMsg)

			return
		}

		_, ok := dbmodel.DocumentScopeTypesMap[desScopeType]
		if !ok {
			errMsg := tlog.E(ctx).Msgf("Handle list documents (scope type: %d) err (scope type invalid)",
				desScopeType)

			SendFailResponse(c, constant.ErrorCodeRequestParamInvalid, errMsg)

			return
		}

		scopeType = desScopeType
	}

	sourceType := -1

	srcSourceType := strings.TrimSpace(c.Query("source_type"))
	if srcSourceType != "" {
		desSourceType, err := strconv.Atoi(srcSourceType)
		if err != nil {
			errMsg := tlog.E(ctx).Err(err).Msgf("Handle list documents (source type: %s) err (strconv atoi %v)",
				srcSourceType, err)

			SendFailResponse(c, constant.ErrorCodeRequestParamInvalid, errMsg)

			return
		}

		_, ok := dbmodel.DocumentSourceTypesMap[desSourceType]
		if !ok {
			errMsg := tlog.E(ctx).Msgf("Handle list documents (source type: %d) err (source type invalid)",
				desSourceType)

			SendFailResponse(c, constant.ErrorCodeRequestParamInvalid, errMsg)

			return
		}

		sourceType = desSourceType
	}

	processStatus := -1

	srcProcessStatus := strings.TrimSpace(c.Query("process_status"))
	if srcProcessStatus != "" {
		desProcessStatus, err := strconv.Atoi(srcProcessStatus)
		if err != nil {
			errMsg := tlog.E(ctx).Err(err).Msgf("Handle list documents (process status: %s) err (strconv atoi %v)",
				srcProcessStatus, err)

			SendFailResponse(c, constant.ErrorCodeRequestParamInvalid, errMsg)

			return
		}

		_, ok := dbmodel.DocumentProcessStatusesMap[desProcessStatus]
		if !ok {
			errMsg := tlog.E(ctx).Msgf("Handle list documents (process status: %d) err (process status invalid)",
				desProcessStatus)

			SendFailResponse(c, constant.ErrorCodeRequestParamInvalid, errMsg)

			return
		}

		processStatus = desProcessStatus
	}

	status := -1

	srcStatus := strings.TrimSpace(c.Query("status"))
	if srcStatus != "" {
		desStatus, err := strconv.Atoi(srcStatus)
		if err != nil {
			errMsg := tlog.E(ctx).Err(err).Msgf("Handle list documents (status: %s) err (strconv atoi %v)",
				srcStatus, err)

			SendFailResponse(c, constant.ErrorCodeRequestParamInvalid, errMsg)

			return
		}

		_, ok := dbmodel.DocumentStatusesMap[desStatus]
		if !ok {
			errMsg := tlog.E(ctx).Msgf("Handle list documents (status: %d) err (status invalid)",
				desStatus)

			SendFailResponse(c, constant.ErrorCodeRequestParamInvalid, errMsg)

			return
		}

		status = desStatus
	}

	pageNum := 1

	srcPageNum := strings.TrimSpace(c.Query("page_num"))
	if srcPageNum != "" {
		desPageNum, err := strconv.Atoi(srcPageNum)
		if err != nil {
			errMsg := tlog.E(ctx).Err(err).Msgf("Handle list documents (page num: %s) err (strconv atoi %v)",
				srcPageNum, err)

			SendFailResponse(c, constant.ErrorCodeRequestParamInvalid, errMsg)

			return
		}

		if desPageNum <= 0 || desPageNum > constant.MaxPageNum {
			errMsg := tlog.E(ctx).Msgf("Handle list documents (page num: %d) err (page num invalid)",
				desPageNum)

			SendFailResponse(c, constant.ErrorCodeRequestParamInvalid, errMsg)

			return
		}

		pageNum = desPageNum
	}

	pageSize := 10

	srcPageSize := strings.TrimSpace(c.Query("page_size"))
	if srcPageSize != "" {
		desPageSize, err := strconv.Atoi(srcPageSize)
		if err != nil {
			errMsg := tlog.E(ctx).Err(err).Msgf("Handle list documents (page size: %s) err (strconv atoi %v)",
				srcPageSize, err)

			SendFailResponse(c, constant.ErrorCodeRequestParamInvalid, errMsg)

			return
		}

		if desPageSize <= 0 || desPageSize > constant.MaxPageSize {
			errMsg := tlog.E(ctx).Msgf("Handle list documents (page size: %d) err (page size invalid)",
				desPageSize)

			SendFailResponse(c, constant.ErrorCodeRequestParamInvalid, errMsg)

			return
		}

		pageSize = desPageSize
	}

	listDocumentsRespData, errx := service.ListDocuments(ctx, userId, knowledgeBaseId, keyword, scopeType, sourceType, processStatus, status, pageNum, pageSize)
	if errx != nil {
		errMsg := tlog.E(ctx).Err(errx).Msgf("Handle list documents (user id: %s, knowledge base id: %s, keyword: %s, scope type: %d, source type: %d, process status: %d, status: %d, page num: %d, page size: %d) err (list documents %v)",
			userId, knowledgeBaseId, keyword, scopeType, sourceType, processStatus, status, pageNum, pageSize, errx)

		SendFailResponse(c, errx.ErrCode(), errMsg)

		return
	}

	SendPassResponse(c, listDocumentsRespData)
}

func HandleGetDocument(c *gin.Context) {
	ctx := c.Request.Context()

	userId := strings.TrimSpace(c.Request.Header.Get("user_id"))

	documentId := strings.TrimSpace(c.Param("id"))
	if documentId == "" {
		errMsg := tlog.E(ctx).Msgf("Handle get document err (document id empty)")

		SendFailResponse(c, constant.ErrorCodeRequestParamInvalid, errMsg)

		return
	}

	getDocumentRespData, errx := service.GetDocument(ctx, userId, documentId)
	if errx != nil {
		errMsg := tlog.E(ctx).Err(errx).Msgf("Handle get document (user id: %s, document id: %s) err (get document %v)",
			userId, documentId, errx)

		SendFailResponse(c, errx.ErrCode(), errMsg)

		return
	}

	SendPassResponse(c, getDocumentRespData)
}

func HandleCreateDocument(c *gin.Context) {
	ctx := c.Request.Context()

	userId := strings.TrimSpace(c.Request.Header.Get("user_id"))

	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, constant.MaxDocumentUploadSize+constant.MB)

	scopeType := dbmodel.DocumentScopeTypeKnowledge

	srcScopeType := strings.TrimSpace(c.PostForm("scope_type"))
	if srcScopeType != "" {
		desScopeType, err := strconv.Atoi(srcScopeType)
		if err != nil {
			errMsg := tlog.E(ctx).Err(err).Msgf("Handle create document (scope type: %s) err (strconv atoi %v)",
				srcScopeType, err)

			SendFailResponse(c, constant.ErrorCodeRequestParamInvalid, errMsg)

			return
		}

		_, ok := dbmodel.DocumentScopeTypesMap[desScopeType]
		if !ok {
			errMsg := tlog.E(ctx).Msgf("Handle create document (scope type: %d) err (scope type invalid)",
				desScopeType)

			SendFailResponse(c, constant.ErrorCodeRequestParamInvalid, errMsg)

			return
		}

		scopeType = desScopeType
	}

	sourceType := dbmodel.DocumentSourceTypeUser

	srcSourceType := strings.TrimSpace(c.PostForm("source_type"))
	if srcSourceType != "" {
		desSourceType, err := strconv.Atoi(srcSourceType)
		if err != nil {
			errMsg := tlog.E(ctx).Err(err).Msgf("Handle create document (source type: %s) err (strconv atoi %v)",
				srcSourceType, err)

			SendFailResponse(c, constant.ErrorCodeRequestParamInvalid, errMsg)

			return
		}

		_, ok := dbmodel.DocumentSourceTypesMap[desSourceType]
		if !ok {
			errMsg := tlog.E(ctx).Msgf("Handle create document (source type: %d) err (source type invalid)",
				desSourceType)

			SendFailResponse(c, constant.ErrorCodeRequestParamInvalid, errMsg)

			return
		}

		sourceType = desSourceType
	}

	parseStrategy := dbmodel.DocumentParseStrategyAuto

	srcParseStrategy := strings.TrimSpace(c.PostForm("parse_strategy"))
	if srcParseStrategy != "" {
		desParseStrategy, err := strconv.Atoi(srcParseStrategy)
		if err != nil {
			errMsg := tlog.E(ctx).Err(err).Msgf("Handle create document (parse strategy: %s) err (strconv atoi %v)",
				srcParseStrategy, err)

			SendFailResponse(c, constant.ErrorCodeRequestParamInvalid, errMsg)

			return
		}

		_, ok := dbmodel.DocumentParseStrategiesMap[desParseStrategy]
		if !ok {
			errMsg := tlog.E(ctx).Msgf("Handle create document (parse strategy: %d) err (parse strategy invalid)",
				desParseStrategy)

			SendFailResponse(c, constant.ErrorCodeRequestParamInvalid, errMsg)

			return
		}

		parseStrategy = desParseStrategy
	}

	knowledgeBaseId := strings.TrimSpace(c.PostForm("knowledge_base_id"))
	if scopeType == dbmodel.DocumentScopeTypeKnowledge && knowledgeBaseId == "" {
		errMsg := tlog.E(ctx).Msgf("Handle create document err (knowledge base id empty)")

		SendFailResponse(c, constant.ErrorCodeRequestParamInvalid, errMsg)

		return
	}

	chatSessionId := strings.TrimSpace(c.PostForm("chat_session_id"))
	if scopeType == dbmodel.DocumentScopeTypeAttachment && chatSessionId == "" {
		errMsg := tlog.E(ctx).Msgf("Handle create document err (chat session id empty)")

		SendFailResponse(c, constant.ErrorCodeRequestParamInvalid, errMsg)

		return
	}

	title := strings.TrimSpace(c.PostForm("title"))
	if len(title) > dbmodel.DocumentTitleLenLimit {
		errMsg := tlog.E(ctx).Msgf("Handle create document (title: %s) err (title len limit)",
			title)

		SendFailResponse(c, constant.ErrorCodeRequestParamInvalid, errMsg)

		return
	}

	summary := strings.TrimSpace(c.PostForm("summary"))
	if len(summary) > dbmodel.DocumentSummaryLenLimit {
		errMsg := tlog.E(ctx).Msgf("Handle create document (summary: %s) err (summary len limit)",
			summary)

		SendFailResponse(c, constant.ErrorCodeRequestParamInvalid, errMsg)

		return
	}

	tags := make([]string, 0)

	srcTags := strings.TrimSpace(c.PostForm("tags"))
	if srcTags != "" {
		if strings.HasPrefix(srcTags, "[") {
			err := json.Unmarshal([]byte(srcTags), &tags)
			if err != nil {
				errMsg := tlog.E(ctx).Err(err).Msgf("Handle create document (tags: %s) err (tags unmarshal %v)",
					srcTags, err)

				SendFailResponse(c, constant.ErrorCodeRequestParamInvalid, errMsg)

				return
			}
		} else {
			tags = strings.Split(srcTags, ",")
		}
	}

	tagsMap := make(map[string]bool)

	for index, tag := range tags {
		tag = strings.TrimSpace(tag)
		tags[index] = tag

		if tag == "" {
			errMsg := tlog.E(ctx).Msgf("Handle create document err (tag empty)")

			SendFailResponse(c, constant.ErrorCodeRequestParamInvalid, errMsg)

			return
		}

		_, ok := tagsMap[tag]
		if ok {
			errMsg := tlog.E(ctx).Msgf("Handle create document (tag: %s) err (tag duplicate)",
				tag)

			SendFailResponse(c, constant.ErrorCodeRequestParamInvalid, errMsg)

			return
		}

		tagsMap[tag] = true
	}

	tagsData, _ := json.Marshal(tags)
	if len(tagsData) > dbmodel.DocumentTagsLenLimit {
		errMsg := tlog.E(ctx).Msgf("Handle create document (tags: %v) err (tags len limit)",
			tags)

		SendFailResponse(c, constant.ErrorCodeRequestParamInvalid, errMsg)

		return
	}

	ownerId := strings.TrimSpace(c.PostForm("owner_id"))
	if userId != "" {
		ownerId = userId
	}
	if ownerId == "" {
		ownerId = constant.DefaultAnonymousOwnerId
	}

	langCode := strings.TrimSpace(c.PostForm("lang_code"))
	if len(langCode) > dbmodel.DocumentLangCodeLen {
		errMsg := tlog.E(ctx).Msgf("Handle create document (lang code: %s) err (lang code len limit)",
			langCode)

		SendFailResponse(c, constant.ErrorCodeRequestParamInvalid, errMsg)

		return
	}
	if langCode == "" {
		langCode = constant.DefaultDocumentLangCode
	}

	fileHeader, err := c.FormFile("file")
	if err != nil {
		errMsg := tlog.E(ctx).Err(err).Msgf("Handle create document err (form file %v)",
			err)

		SendFailResponse(c, constant.ErrorCodeDocumentFileInvalid, errMsg)

		return
	}

	if fileHeader.Size <= 0 {
		errMsg := tlog.E(ctx).Msgf("Handle create document err (file empty)")

		SendFailResponse(c, constant.ErrorCodeDocumentFileInvalid, errMsg)

		return
	}

	if fileHeader.Size > constant.MaxDocumentUploadSize {
		errMsg := tlog.E(ctx).Msgf("Handle create document (file size: %d) err (file size limit)",
			fileHeader.Size)

		SendFailResponse(c, constant.ErrorCodeDocumentFileInvalid, errMsg)

		return
	}

	fileHeader.Filename = strings.TrimSpace(fileHeader.Filename)
	if len(fileHeader.Filename) > dbmodel.FileObjectFileNameLen {
		errMsg := tlog.E(ctx).Msgf("Handle create document (file name: %s) err (file name len limit)",
			fileHeader.Filename)

		SendFailResponse(c, constant.ErrorCodeRequestParamInvalid, errMsg)

		return
	}

	createDocumentRespData, errx := service.CreateDocument(ctx, userId, knowledgeBaseId, chatSessionId, scopeType, sourceType,
		title, summary, tags, ownerId, langCode, parseStrategy, fileHeader)
	if errx != nil {
		errMsg := tlog.E(ctx).Err(errx).Msgf("Handle create document (user id: %s, knowledge base id: %s, chat session id: %s, scope type: %d, source type: %d, title: %s, summary: %s, tags: %v, owner id: %s, lang code: %s, parse strategy: %d, file name: %s, file size: %d) err (create document %v)",
			userId, knowledgeBaseId, chatSessionId, scopeType, sourceType, title, summary, tags, ownerId, langCode, parseStrategy, fileHeader.Filename, fileHeader.Size, errx)

		SendFailResponse(c, errx.ErrCode(), errMsg)

		return
	}

	SendPassResponse(c, createDocumentRespData)
}

func HandleUpdateDocument(c *gin.Context) {
	ctx := c.Request.Context()

	userId := strings.TrimSpace(c.Request.Header.Get("user_id"))

	documentId := strings.TrimSpace(c.Param("id"))
	if documentId == "" {
		errMsg := tlog.E(ctx).Msgf("Handle update document err (document id empty)")

		SendFailResponse(c, constant.ErrorCodeRequestParamInvalid, errMsg)

		return
	}

	updateDocumentRequest := &data.UpdateDocumentRequest{}

	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, constant.RequestBodyMaxSize)

	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		errMsg := tlog.E(ctx).Err(err).Msgf("Handle update document (body: %s) err (request body read %v)",
			string(body), err)

		SendFailResponse(c, constant.ErrorCodeRequestBodyInvalid, errMsg)

		return
	}

	err = json.Unmarshal(body, updateDocumentRequest)
	if err != nil {
		if strings.TrimSpace(string(body)) == "" {
			errMsg := tlog.E(ctx).Err(err).Msgf("Handle update document (body: %s) err (request body empty or unmarshal %v)",
				string(body), err)

			SendFailResponse(c, constant.ErrorCodeRequestBodyInvalid, errMsg)

			return
		}

		errMsg := tlog.E(ctx).Err(err).Msgf("Handle update document (body: %s) err (request body unmarshal %v)",
			string(body), err)

		SendFailResponse(c, constant.ErrorCodeRequestBodyInvalid, errMsg)

		return
	}

	title := strings.TrimSpace(updateDocumentRequest.Title)
	if title == "" {
		errMsg := tlog.E(ctx).Msgf("Handle update document err (title empty)")

		SendFailResponse(c, constant.ErrorCodeRequestParamInvalid, errMsg)

		return
	}

	if len(title) > dbmodel.DocumentTitleLenLimit {
		errMsg := tlog.E(ctx).Msgf("Handle update document (title: %s) err (title len limit)",
			title)

		SendFailResponse(c, constant.ErrorCodeRequestParamInvalid, errMsg)

		return
	}

	summary := strings.TrimSpace(updateDocumentRequest.Summary)
	if len(summary) > dbmodel.DocumentSummaryLenLimit {
		errMsg := tlog.E(ctx).Msgf("Handle update document (summary: %s) err (summary len limit)",
			summary)

		SendFailResponse(c, constant.ErrorCodeRequestParamInvalid, errMsg)

		return
	}

	tags := updateDocumentRequest.Tags
	if tags == nil {
		tags = make([]string, 0)
	}

	tagsMap := make(map[string]bool)

	for index, tag := range tags {
		tag = strings.TrimSpace(tag)
		tags[index] = tag

		if tag == "" {
			errMsg := tlog.E(ctx).Msgf("Handle update document err (tag empty)")

			SendFailResponse(c, constant.ErrorCodeRequestParamInvalid, errMsg)

			return
		}

		_, ok := tagsMap[tag]
		if ok {
			errMsg := tlog.E(ctx).Msgf("Handle update document (tag: %s) err (tag duplicate)",
				tag)

			SendFailResponse(c, constant.ErrorCodeRequestParamInvalid, errMsg)

			return
		}

		tagsMap[tag] = true
	}

	tagsData, _ := json.Marshal(tags)
	if len(tagsData) > dbmodel.DocumentTagsLenLimit {
		errMsg := tlog.E(ctx).Msgf("Handle update document (tags: %v) err (tags len limit)",
			tags)

		SendFailResponse(c, constant.ErrorCodeRequestParamInvalid, errMsg)

		return
	}

	langCode := strings.TrimSpace(updateDocumentRequest.LangCode)
	if len(langCode) > dbmodel.DocumentLangCodeLen {
		errMsg := tlog.E(ctx).Msgf("Handle update document (lang code: %s) err (lang code len limit)",
			langCode)

		SendFailResponse(c, constant.ErrorCodeRequestParamInvalid, errMsg)

		return
	}
	if langCode == "" {
		langCode = constant.DefaultDocumentLangCode
	}

	status := updateDocumentRequest.Status

	_, ok := dbmodel.DocumentStatusesMap[status]
	if !ok {
		errMsg := tlog.E(ctx).Msgf("Handle update document (status: %d) err (status invalid)",
			status)

		SendFailResponse(c, constant.ErrorCodeRequestParamInvalid, errMsg)

		return
	}

	errx := service.UpdateDocument(ctx, userId, documentId, title, summary, tags, langCode, status)
	if errx != nil {
		errMsg := tlog.E(ctx).Err(errx).Msgf("Handle update document (user id: %s, document id: %s, title: %s, summary: %s, tags: %v, lang code: %s, status: %d) err (update document %v)",
			userId, documentId, title, summary, tags, langCode, status, errx)

		SendFailResponse(c, errx.ErrCode(), errMsg)

		return
	}

	SendPassResponse(c, nil)
}

func HandleDeleteDocument(c *gin.Context) {
	ctx := c.Request.Context()

	userId := strings.TrimSpace(c.Request.Header.Get("user_id"))

	documentId := strings.TrimSpace(c.Param("id"))
	if documentId == "" {
		errMsg := tlog.E(ctx).Msgf("Handle delete document err (document id empty)")

		SendFailResponse(c, constant.ErrorCodeRequestParamInvalid, errMsg)

		return
	}

	errx := service.DeleteDocument(ctx, userId, documentId)
	if errx != nil {
		errMsg := tlog.E(ctx).Err(errx).Msgf("Handle delete document (user id: %s, document id: %s) err (delete document %v)",
			userId, documentId, errx)

		SendFailResponse(c, errx.ErrCode(), errMsg)

		return
	}

	SendPassResponse(c, nil)
}
