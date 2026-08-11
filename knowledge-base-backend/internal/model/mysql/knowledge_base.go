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
	KnowledgeBaseCodeLen = 64
	KnowledgeBaseNameLen = 128

	KnowledgeBaseDescriptionLen = 65535
)

const (
	KnowledgeBaseVisiblePrivate  = 0
	KnowledgeBaseVisibleInternal = 1
	KnowledgeBaseVisiblePublic   = 2
)

var (
	KnowledgeBaseVisiblesMap = map[int]bool{
		KnowledgeBaseVisiblePrivate:  true,
		KnowledgeBaseVisibleInternal: true,
		KnowledgeBaseVisiblePublic:   true,
	}
)

const (
	KnowledgeBaseStatusDisabled = 0
	KnowledgeBaseStatusNormal   = 1
)

var (
	KnowledgeBaseStatusesMap = map[int]bool{
		KnowledgeBaseStatusDisabled: true,
		KnowledgeBaseStatusNormal:   true,
	}
)

type KnowledgeBase struct {
	Id string

	Code string
	Name string

	OwnerId string

	Description string

	Visible int

	Status int

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt
}

func (*KnowledgeBase) TableName() string {
	return "knowledge_bases"
}

func CreateKnowledgeBase(ctx context.Context, code, name string, ownerId string, description string, visible int, status int) (*KnowledgeBase, *terror.Terror) {
	knowledgeBaseDB := &KnowledgeBase{
		Id: tutil.NewOid().String(),

		Code: code,
		Name: name,

		OwnerId: ownerId,

		Description: description,

		Visible: visible,

		Status: status,
	}

	retGorm := serverClient.DB(ctx, runMode).Create(knowledgeBaseDB)
	if retGorm.Error != nil {
		errMsg := tlog.E(ctx).Err(retGorm.Error).Msgf("Create knowledge base (code: %s, name: %s, owner id: %s, description: %s, visible: %d, status: %d) err (db create %v)",
			code, name, ownerId, description, visible, status, retGorm.Error)

		errx := terror.NewTerror(ctx, retGorm.Error, constant.ErrorCodeMysqlServerAbnormal, errMsg)

		return nil, errx
	}

	return knowledgeBaseDB, nil
}

func FindKnowledgeBase(ctx context.Context, knowledgeBaseId string) (*KnowledgeBase, *terror.Terror) {
	knowledgeBasesDB := make([]*KnowledgeBase, 0)

	retGorm := serverClient.DB(ctx, runMode).Where("id = ?", knowledgeBaseId).Limit(1).Find(&knowledgeBasesDB)
	if retGorm.Error != nil {
		errMsg := tlog.E(ctx).Err(retGorm.Error).Msgf("Find knowledge base (id: %s) err (db find %v)",
			knowledgeBaseId, retGorm.Error)

		errx := terror.NewTerror(ctx, retGorm.Error, constant.ErrorCodeMysqlServerAbnormal, errMsg)

		return nil, errx
	}

	if len(knowledgeBasesDB) == 0 {
		return nil, nil
	}

	return knowledgeBasesDB[0], nil
}

func FindKnowledgeBaseByCode(ctx context.Context, code string) (*KnowledgeBase, *terror.Terror) {
	knowledgeBasesDB := make([]*KnowledgeBase, 0)

	retGorm := serverClient.DB(ctx, runMode).Where("code = ?", code).Limit(1).Find(&knowledgeBasesDB)
	if retGorm.Error != nil {
		errMsg := tlog.E(ctx).Err(retGorm.Error).Msgf("Find knowledge base (code: %s) err (db find %v)",
			code, retGorm.Error)

		errx := terror.NewTerror(ctx, retGorm.Error, constant.ErrorCodeMysqlServerAbnormal, errMsg)

		return nil, errx
	}

	if len(knowledgeBasesDB) == 0 {
		return nil, nil
	}

	return knowledgeBasesDB[0], nil
}

func FindKnowledgeBases(ctx context.Context, keyword string, visible, status int, pageNum, pageSize int) (int64, []*KnowledgeBase, *terror.Terror) {
	query := serverClient.DB(ctx, runMode)

	if keyword != "" {
		query = query.Where("code like ? OR name like ?", "%"+keyword+"%", "%"+keyword+"%")
	}

	if visible != -1 {
		query = query.Where("visible = ?", visible)
	}

	if status != -1 {
		query = query.Where("status = ?", status)
	}

	total := int64(0)

	retGorm := query.Model(&KnowledgeBase{}).Count(&total)
	if retGorm.Error != nil {
		errMsg := tlog.E(ctx).Err(retGorm.Error).Msgf("Find knowledge bases (keyword: %s, visible: %d, status: %d, page num: %d, page size: %d) err (db count %v)",
			keyword, visible, status, pageNum, pageSize, retGorm.Error)

		errx := terror.NewTerror(ctx, retGorm.Error, constant.ErrorCodeMysqlServerAbnormal, errMsg)

		return -1, nil, errx
	}

	if total == 0 {
		return 0, nil, nil
	}

	knowledgeBasesDB := make([]*KnowledgeBase, 0)

	retGorm = query.Offset((pageNum - 1) * pageSize).Limit(pageSize).Order("created_at DESC").Find(&knowledgeBasesDB)
	if retGorm.Error != nil {
		errMsg := tlog.E(ctx).Err(retGorm.Error).Msgf("Find knowledge bases (keyword: %s, visible: %d, status: %d, page num: %d, page size: %d) err (db find %v)",
			keyword, visible, status, pageNum, pageSize, retGorm.Error)

		errx := terror.NewTerror(ctx, retGorm.Error, constant.ErrorCodeMysqlServerAbnormal, errMsg)

		return -1, nil, errx
	}

	return total, knowledgeBasesDB, nil
}

func UpdateKnowledgeBase(ctx context.Context, knowledgeBaseId string, name string, description string, visible int, status int) *terror.Terror {
	params := map[string]any{
		"name": name,

		"description": description,

		"visible": visible,

		"status": status,

		"updated_at": time.Now(),
	}

	retGorm := serverClient.DB(ctx, runMode).Model(&KnowledgeBase{}).Where("id = ?", knowledgeBaseId).Updates(params)
	if retGorm.Error != nil {
		errMsg := tlog.E(ctx).Err(retGorm.Error).Msgf("Update knowledge base (id: %s, name: %s, description: %s, visible: %d, status: %d) err (db updates %v)",
			knowledgeBaseId, name, description, visible, status, retGorm.Error)

		errx := terror.NewTerror(ctx, retGorm.Error, constant.ErrorCodeMysqlServerAbnormal, errMsg)

		return errx
	}

	return nil
}

func DeleteKnowledgeBase(ctx context.Context, knowledgeBaseId string) *terror.Terror {
	retGorm := serverClient.DB(ctx, runMode).Where("id = ?", knowledgeBaseId).Delete(&KnowledgeBase{})
	if retGorm.Error != nil {
		errMsg := tlog.E(ctx).Err(retGorm.Error).Msgf("Delete knowledge base (id: %s) err (db delete %v)",
			knowledgeBaseId, retGorm.Error)

		errx := terror.NewTerror(ctx, retGorm.Error, constant.ErrorCodeMysqlServerAbnormal, errMsg)

		return errx
	}

	return nil
}
