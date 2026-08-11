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

func HandleListKnowledgeBases(c *gin.Context) {
	ctx := c.Request.Context()

	userId := c.Request.Header.Get("user_id")

	keyword := strings.TrimSpace(c.Query("keyword"))

	visible := -1

	srcVisible := strings.TrimSpace(c.Query("visible"))
	if srcVisible != "" {
		desVisible, err := strconv.Atoi(srcVisible)
		if err != nil {
			errMsg := tlog.E(ctx).Err(err).Msgf("Handle list knowledge bases (visible: %s) err (strconv atoi %v)",
				srcVisible, err)

			SendFailResponse(c, constant.ErrorCodeRequestParamInvalid, errMsg)

			return
		}

		_, ok := dbmodel.KnowledgeBaseVisiblesMap[desVisible]
		if !ok {
			errMsg := tlog.E(ctx).Msgf("Handle list knowledge bases (visible: %d) err (visible invalid)",
				desVisible)

			SendFailResponse(c, constant.ErrorCodeRequestParamInvalid, errMsg)

			return
		}

		visible = desVisible
	}

	status := -1

	srcStatus := strings.TrimSpace(c.Query("status"))
	if srcStatus != "" {
		desStatus, err := strconv.Atoi(srcStatus)
		if err != nil {
			errMsg := tlog.E(ctx).Err(err).Msgf("Handle list knowledge bases (status: %s) err (strconv atoi %v)",
				srcStatus, err)

			SendFailResponse(c, constant.ErrorCodeRequestParamInvalid, errMsg)

			return
		}

		_, ok := dbmodel.KnowledgeBaseStatusesMap[desStatus]
		if !ok {
			errMsg := tlog.E(ctx).Msgf("Handle list knowledge bases (status: %d) err (status invalid)",
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
			errMsg := tlog.E(ctx).Err(err).Msgf("Handle list knowledge bases (page num: %s) err (strconv atoi %v)",
				srcPageNum, err)

			SendFailResponse(c, constant.ErrorCodeRequestParamInvalid, errMsg)

			return
		}

		if desPageNum <= 0 || desPageNum > constant.MaxPageNum {
			errMsg := tlog.E(ctx).Msgf("Handle list knowledge bases (page num: %d) err (page num invalid)",
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
			errMsg := tlog.E(ctx).Err(err).Msgf("Handle list knowledge bases (page size: %s) err (strconv atoi %v)",
				srcPageSize, err)

			SendFailResponse(c, constant.ErrorCodeRequestParamInvalid, errMsg)

			return
		}

		if desPageSize <= 0 || desPageSize > constant.MaxPageSize {
			errMsg := tlog.E(ctx).Msgf("Handle list knowledge bases (page size: %d) err (page size invalid)",
				desPageSize)

			SendFailResponse(c, constant.ErrorCodeRequestParamInvalid, errMsg)

			return
		}

		pageSize = desPageSize
	}

	listKnowledgeBasesRespData, errx := service.ListKnowledgeBases(ctx, userId, keyword, visible, status, pageNum, pageSize)
	if errx != nil {
		errMsg := tlog.E(ctx).Err(errx).Msgf("Handle list knowledge bases (user id: %s, keyword: %s, visible: %d, status: %d, page num: %d, page size: %d) err (list knowledge bases %v)",
			userId, keyword, visible, status, pageNum, pageSize, errx)

		SendFailResponse(c, errx.ErrCode(), errMsg)

		return
	}

	SendPassResponse(c, listKnowledgeBasesRespData)
}

func HandleGetKnowledgeBase(c *gin.Context) {
	ctx := c.Request.Context()

	userId := c.Request.Header.Get("user_id")

	knowledgeBaseId := strings.TrimSpace(c.Param("id"))
	if knowledgeBaseId == "" {
		errMsg := tlog.E(ctx).Msgf("Handle get knowledge base err (knowledge base id empty)")

		SendFailResponse(c, constant.ErrorCodeRequestParamInvalid, errMsg)

		return
	}

	getKnowledgeBaseRespData, errx := service.GetKnowledgeBase(ctx, userId, knowledgeBaseId)
	if errx != nil {
		errMsg := tlog.E(ctx).Err(errx).Msgf("Handle get knowledge base (user id: %s, knowledge base id: %s) err (get knowledge base %v)",
			userId, knowledgeBaseId, errx)

		SendFailResponse(c, errx.ErrCode(), errMsg)

		return
	}

	SendPassResponse(c, getKnowledgeBaseRespData)
}

func HandleCreateKnowledgeBase(c *gin.Context) {
	ctx := c.Request.Context()

	userId := c.Request.Header.Get("user_id")

	createKnowledgeBaseRequest := &data.CreateKnowledgeBaseRequest{}

	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, constant.RequestBodyMaxSize)

	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		errMsg := tlog.E(ctx).Err(err).Msgf("Handle create knowledge base (body: %s) err (request body read %v)",
			string(body), err)

		SendFailResponse(c, constant.ErrorCodeRequestBodyInvalid, errMsg)

		return
	}

	err = json.Unmarshal(body, createKnowledgeBaseRequest)
	if err != nil {
		errMsg := tlog.E(ctx).Err(err).Msgf("Handle create knowledge base (body: %s) err (request body unmarshal %v)",
			string(body), err)

		SendFailResponse(c, constant.ErrorCodeRequestBodyInvalid, errMsg)

		return
	}

	code := strings.TrimSpace(createKnowledgeBaseRequest.Code)
	if code == "" {
		errMsg := tlog.E(ctx).Msgf("Handle create knowledge base err (code empty)")

		SendFailResponse(c, constant.ErrorCodeRequestParamInvalid, errMsg)

		return
	}

	if len(code) > dbmodel.KnowledgeBaseCodeLen {
		errMsg := tlog.E(ctx).Msgf("Handle create knowledge base (code: %s) err (code len limit)",
			code)

		SendFailResponse(c, constant.ErrorCodeRequestParamInvalid, errMsg)

		return
	}

	name := strings.TrimSpace(createKnowledgeBaseRequest.Name)
	if name == "" {
		errMsg := tlog.E(ctx).Msgf("Handle create knowledge base err (name empty)")

		SendFailResponse(c, constant.ErrorCodeRequestParamInvalid, errMsg)

		return
	}

	if len(name) > dbmodel.KnowledgeBaseNameLen {
		errMsg := tlog.E(ctx).Msgf("Handle create knowledge base (name: %s) err (name len limit)",
			name)

		SendFailResponse(c, constant.ErrorCodeRequestParamInvalid, errMsg)

		return
	}

	description := strings.TrimSpace(createKnowledgeBaseRequest.Description)
	if len(description) > dbmodel.KnowledgeBaseDescriptionLen {
		errMsg := tlog.E(ctx).Msgf("Handle create knowledge base (description: %s) err (description len limit)",
			description)

		SendFailResponse(c, constant.ErrorCodeRequestParamInvalid, errMsg)

		return
	}

	visible := createKnowledgeBaseRequest.Visible

	_, ok := dbmodel.KnowledgeBaseVisiblesMap[visible]
	if !ok {
		errMsg := tlog.E(ctx).Msgf("Handle create knowledge base (visible: %d) err (visible invalid)",
			visible)

		SendFailResponse(c, constant.ErrorCodeRequestParamInvalid, errMsg)

		return
	}

	status := createKnowledgeBaseRequest.Status

	_, ok = dbmodel.KnowledgeBaseStatusesMap[status]
	if !ok {
		errMsg := tlog.E(ctx).Msgf("Handle create knowledge base (status: %d) err (status invalid)",
			status)

		SendFailResponse(c, constant.ErrorCodeRequestParamInvalid, errMsg)

		return
	}

	createKnowledgeBaseRespData, errx := service.CreateKnowledgeBase(ctx, userId, code, name, description, visible, status)
	if errx != nil {
		errMsg := tlog.E(ctx).Err(errx).Msgf("Handle create knowledge base (user id: %s, code: %s, name: %s, description: %s, visible: %d, status: %d) err (create knowledge base %v)",
			userId, code, name, description, visible, status, errx)

		SendFailResponse(c, errx.ErrCode(), errMsg)

		return
	}

	SendPassResponse(c, createKnowledgeBaseRespData)
}

func HandleUpdateKnowledgeBase(c *gin.Context) {
	ctx := c.Request.Context()

	userId := c.Request.Header.Get("user_id")

	knowledgeBaseId := strings.TrimSpace(c.Param("id"))
	if knowledgeBaseId == "" {
		errMsg := tlog.E(ctx).Msgf("Handle update knowledge base err (knowledge base id empty)")

		SendFailResponse(c, constant.ErrorCodeRequestParamInvalid, errMsg)

		return
	}

	updateKnowledgeBaseRequest := &data.UpdateKnowledgeBaseRequest{}

	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, constant.RequestBodyMaxSize)

	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		errMsg := tlog.E(ctx).Err(err).Msgf("Handle update knowledge base (body: %s) err (request body read %v)",
			string(body), err)

		SendFailResponse(c, constant.ErrorCodeRequestBodyInvalid, errMsg)

		return
	}

	err = json.Unmarshal(body, updateKnowledgeBaseRequest)
	if err != nil {
		errMsg := tlog.E(ctx).Err(err).Msgf("Handle update knowledge base (body: %s) err (request body unmarshal %v)",
			string(body), err)

		SendFailResponse(c, constant.ErrorCodeRequestBodyInvalid, errMsg)

		return
	}

	name := strings.TrimSpace(updateKnowledgeBaseRequest.Name)
	if name == "" {
		errMsg := tlog.E(ctx).Msgf("Handle update knowledge base err (name empty)")

		SendFailResponse(c, constant.ErrorCodeRequestParamInvalid, errMsg)

		return
	}

	if len(name) > dbmodel.KnowledgeBaseNameLen {
		errMsg := tlog.E(ctx).Msgf("Handle update knowledge base (name: %s) err (name len limit)",
			name)

		SendFailResponse(c, constant.ErrorCodeRequestParamInvalid, errMsg)

		return
	}

	description := strings.TrimSpace(updateKnowledgeBaseRequest.Description)
	if len(description) > dbmodel.KnowledgeBaseDescriptionLen {
		errMsg := tlog.E(ctx).Msgf("Handle update knowledge base (description: %s) err (description len limit)",
			description)

		SendFailResponse(c, constant.ErrorCodeRequestParamInvalid, errMsg)

		return
	}

	visible := updateKnowledgeBaseRequest.Visible

	_, ok := dbmodel.KnowledgeBaseVisiblesMap[visible]
	if !ok {
		errMsg := tlog.E(ctx).Msgf("Handle update knowledge base (visible: %d) err (visible invalid)",
			visible)

		SendFailResponse(c, constant.ErrorCodeRequestParamInvalid, errMsg)

		return
	}

	status := updateKnowledgeBaseRequest.Status

	_, ok = dbmodel.KnowledgeBaseStatusesMap[status]
	if !ok {
		errMsg := tlog.E(ctx).Msgf("Handle update knowledge base (status: %d) err (status invalid)",
			status)

		SendFailResponse(c, constant.ErrorCodeRequestParamInvalid, errMsg)

		return
	}

	errx := service.UpdateKnowledgeBase(ctx, userId, knowledgeBaseId, name, description, visible, status)
	if errx != nil {
		errMsg := tlog.E(ctx).Err(errx).Msgf("Handle update knowledge base (user id: %s, knowledge base id: %s, name: %s, description: %s, visible: %d, status: %d) err (update knowledge base %v)",
			userId, knowledgeBaseId, name, description, visible, status, errx)

		SendFailResponse(c, errx.ErrCode(), errMsg)

		return
	}

	SendPassResponse(c, nil)
}

func HandleDeleteKnowledgeBase(c *gin.Context) {
	ctx := c.Request.Context()

	userId := c.Request.Header.Get("user_id")

	knowledgeBaseId := strings.TrimSpace(c.Param("id"))
	if knowledgeBaseId == "" {
		errMsg := tlog.E(ctx).Msgf("Handle delete knowledge base err (knowledge base id empty)")

		SendFailResponse(c, constant.ErrorCodeRequestParamInvalid, errMsg)

		return
	}

	errx := service.DeleteKnowledgeBase(ctx, userId, knowledgeBaseId)
	if errx != nil {
		errMsg := tlog.E(ctx).Err(errx).Msgf("Handle delete knowledge base (user id: %s, knowledge base id: %s) err (delete knowledge base %v)",
			userId, knowledgeBaseId, errx)

		SendFailResponse(c, errx.ErrCode(), errMsg)

		return
	}

	SendPassResponse(c, nil)
}
