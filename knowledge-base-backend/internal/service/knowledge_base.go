package service

import (
	"context"
	"time"

	"github.com/choveylee/terror"
	"github.com/choveylee/tlog"

	constant "dev.choveylee.top/knowledge-base-backend/internal/const"
	"dev.choveylee.top/knowledge-base-backend/internal/data"
	dbmodel "dev.choveylee.top/knowledge-base-backend/internal/model/mysql"
)

func ListKnowledgeBases(ctx context.Context, userId string, keyword string, visible, status int, pageNum, pageSize int) (*data.ListKnowledgeBasesRespData, *terror.Terror) {
	total, knowledgeBasesDB, errx := dbmodel.FindKnowledgeBases(ctx, keyword, visible, status, pageNum, pageSize)
	if errx != nil {
		errMsg := tlog.E(ctx).Err(errx).Msgf("List knowledge bases (user id: %s, keyword: %s, visible: %d, status: %d, page num: %d, page size: %d) err (db find knowledge bases %v)",
			userId, keyword, visible, status, pageNum, pageSize, errx)
		errx.AttachErrMsg(errMsg)

		return nil, errx
	}

	listKnowledgeBasesRespData := &data.ListKnowledgeBasesRespData{
		List: make([]*data.KnowledgeBaseData, 0),

		Total: total,
	}

	for _, knowledgeBaseDB := range knowledgeBasesDB {
		knowledgeBaseData := &data.KnowledgeBaseData{
			KnowledgeBaseId: knowledgeBaseDB.Id,

			Code: knowledgeBaseDB.Code,
			Name: knowledgeBaseDB.Name,

			OwnerId: knowledgeBaseDB.OwnerId,

			Description: knowledgeBaseDB.Description,

			Visible: knowledgeBaseDB.Visible,

			Status: knowledgeBaseDB.Status,

			CreatedAt: knowledgeBaseDB.CreatedAt.Format(time.RFC3339),
			UpdatedAt: knowledgeBaseDB.UpdatedAt.Format(time.RFC3339),
		}

		listKnowledgeBasesRespData.List = append(listKnowledgeBasesRespData.List, knowledgeBaseData)
	}

	return listKnowledgeBasesRespData, nil
}

func GetKnowledgeBase(ctx context.Context, userId string, knowledgeBaseId string) (*data.GetKnowledgeBaseRespData, *terror.Terror) {
	knowledgeBaseDB, errx := dbmodel.FindKnowledgeBase(ctx, knowledgeBaseId)
	if errx != nil {
		errMsg := tlog.E(ctx).Err(errx).Msgf("Get knowledge base (user id: %s, knowledge base id: %s) err (db find knowledge base %v)",
			userId, knowledgeBaseId, errx)
		errx.AttachErrMsg(errMsg)

		return nil, errx
	}

	if knowledgeBaseDB == nil {
		errMsg := tlog.E(ctx).Msgf("Get knowledge base (user id: %s, knowledge base id: %s) err (knowledge base not found)",
			userId, knowledgeBaseId)

		errx := terror.NewTerror(ctx, terror.ErrParamInvalid("knowledge base id"), constant.ErrorCodeKnowledgeBaseNotFound, errMsg)

		return nil, errx
	}

	getKnowledgeBaseRespData := &data.GetKnowledgeBaseRespData{
		KnowledgeBaseId: knowledgeBaseDB.Id,

		Code: knowledgeBaseDB.Code,
		Name: knowledgeBaseDB.Name,

		OwnerId: knowledgeBaseDB.OwnerId,

		Description: knowledgeBaseDB.Description,

		Visible: knowledgeBaseDB.Visible,

		Status: knowledgeBaseDB.Status,

		CreatedAt: knowledgeBaseDB.CreatedAt.Format(time.RFC3339),
		UpdatedAt: knowledgeBaseDB.UpdatedAt.Format(time.RFC3339),
	}

	return getKnowledgeBaseRespData, nil
}

func CreateKnowledgeBase(ctx context.Context, userId string, code, name string, description string, visible, status int) (*data.CreateKnowledgeBaseRespData, *terror.Terror) {
	knowledgeBaseDB, errx := dbmodel.FindKnowledgeBaseByCode(ctx, code)
	if errx != nil {
		errMsg := tlog.E(ctx).Err(errx).Msgf("Create knowledge base (user id: %s, code: %s, name: %s, description: %s, visible: %d, status: %d) err (db find knowledge base by code %v)",
			userId, code, name, description, visible, status, errx)
		errx.AttachErrMsg(errMsg)

		return nil, errx
	}

	if knowledgeBaseDB != nil {
		errMsg := tlog.E(ctx).Msgf("Create knowledge base (user id: %s, code: %s, name: %s, description: %s, visible: %d, status: %d) err (knowledge base code exist)",
			userId, code, name, description, visible, status)

		errx := terror.NewTerror(ctx, terror.ErrParamInvalid("code"), constant.ErrorCodeKnowledgeBaseCodeExist, errMsg)

		return nil, errx
	}

	knowledgeBaseDB, errx = dbmodel.CreateKnowledgeBase(ctx, code, name, userId, description, visible, status)
	if errx != nil {
		errMsg := tlog.E(ctx).Err(errx).Msgf("Create knowledge base (user id: %s, code: %s, name: %s, description: %s, visible: %d, status: %d) err (db create knowledge base %v)",
			userId, code, name, description, visible, status, errx)
		errx.AttachErrMsg(errMsg)

		return nil, errx
	}

	createKnowledgeBaseRespData := &data.CreateKnowledgeBaseRespData{
		KnowledgeBaseId: knowledgeBaseDB.Id,
	}

	return createKnowledgeBaseRespData, nil
}

func UpdateKnowledgeBase(ctx context.Context, userId string, knowledgeBaseId string, name string, description string, visible, status int) *terror.Terror {
	knowledgeBaseDB, errx := dbmodel.FindKnowledgeBase(ctx, knowledgeBaseId)
	if errx != nil {
		errMsg := tlog.E(ctx).Err(errx).Msgf("Update knowledge base (user id: %s, knowledge base id: %s, name: %s, description: %s, visible: %d, status: %d) err (db find knowledge base %v)",
			userId, knowledgeBaseId, name, description, visible, status, errx)
		errx.AttachErrMsg(errMsg)

		return errx
	}

	if knowledgeBaseDB == nil {
		errMsg := tlog.E(ctx).Msgf("Update knowledge base (user id: %s, knowledge base id: %s, name: %s, description: %s, visible: %d, status: %d) err (knowledge base not found)",
			userId, knowledgeBaseId, name, description, visible, status)

		errx := terror.NewTerror(ctx, terror.ErrParamInvalid("knowledge base id"), constant.ErrorCodeKnowledgeBaseNotFound, errMsg)

		return errx
	}

	errx = dbmodel.UpdateKnowledgeBase(ctx, knowledgeBaseId, name, description, visible, status)
	if errx != nil {
		errMsg := tlog.E(ctx).Err(errx).Msgf("Update knowledge base (user id: %s, knowledge base id: %s, name: %s, description: %s, visible: %d, status: %d) err (db update knowledge base %v)",
			userId, knowledgeBaseId, name, description, visible, status, errx)
		errx.AttachErrMsg(errMsg)

		return errx
	}

	return nil
}

func DeleteKnowledgeBase(ctx context.Context, userId string, knowledgeBaseId string) *terror.Terror {
	knowledgeBaseDB, errx := dbmodel.FindKnowledgeBase(ctx, knowledgeBaseId)
	if errx != nil {
		errMsg := tlog.E(ctx).Err(errx).Msgf("Delete knowledge base (user id: %s, knowledge base id: %s) err (db find knowledge base %v)",
			userId, knowledgeBaseId, errx)
		errx.AttachErrMsg(errMsg)

		return errx
	}

	if knowledgeBaseDB == nil {
		errMsg := tlog.E(ctx).Msgf("Delete knowledge base (user id: %s, knowledge base id: %s) err (knowledge base not found)",
			userId, knowledgeBaseId)

		errx := terror.NewTerror(ctx, terror.ErrParamInvalid("knowledge base id"), constant.ErrorCodeKnowledgeBaseNotFound, errMsg)

		return errx
	}

	errx = dbmodel.DeleteKnowledgeBase(ctx, knowledgeBaseId)
	if errx != nil {
		errMsg := tlog.E(ctx).Err(errx).Msgf("Delete knowledge base (user id: %s, knowledge base id: %s) err (db delete knowledge base %v)",
			userId, knowledgeBaseId, errx)
		errx.AttachErrMsg(errMsg)

		return errx
	}

	return nil
}
