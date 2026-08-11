package data

type KnowledgeBaseData struct {
	KnowledgeBaseId string `json:"knowledge_base_id"`

	Code string `json:"code"`
	Name string `json:"name"`

	OwnerId string `json:"owner_id"`

	Description string `json:"description"`

	Visible int `json:"visible"`

	Status int `json:"status"`

	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type ListKnowledgeBasesRespData struct {
	List []*KnowledgeBaseData `json:"list"`

	Total int64 `json:"total"`
}

type GetKnowledgeBaseRespData struct {
	KnowledgeBaseId string `json:"knowledge_base_id"`

	Code string `json:"code"`
	Name string `json:"name"`

	OwnerId string `json:"owner_id"`

	Description string `json:"description"`

	Visible int `json:"visible"`

	Status int `json:"status"`

	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type CreateKnowledgeBaseRequest struct {
	Code string `json:"code"`

	Name string `json:"name"`

	Description string `json:"description"`

	Visible int `json:"visible"`

	Status int `json:"status"`
}

type CreateKnowledgeBaseRespData struct {
	KnowledgeBaseId string `json:"knowledge_base_id"`
}

type UpdateKnowledgeBaseRequest struct {
	Name string `json:"name"`

	Description string `json:"description"`

	Visible int `json:"visible"`

	Status int `json:"status"`
}
