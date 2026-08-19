package data

type UploadFileRespData struct {
	FileObjectId string `json:"file_object_id"`

	BucketName string `json:"bucket_name"`
	ObjectKey  string `json:"object_key"`

	FileName string `json:"file_name"`
	MimeType string `json:"mime_type"`
	FileExt  string `json:"file_ext"`

	SizeBytes uint64 `json:"size_bytes"`
	Sha256    string `json:"sha256"`

	StorageProvider int `json:"storage_provider"`

	CreatedAt string `json:"created_at"`
}

type DocumentData struct {
	DocumentId string `json:"document_id"`

	KnowledgeBaseId string `json:"knowledge_base_id,omitempty"`
	ChatSessionId   string `json:"chat_session_id,omitempty"`

	ScopeType  int `json:"scope_type"`
	SourceType int `json:"source_type"`

	Title   string   `json:"title"`
	Summary string   `json:"summary"`
	Tags    []string `json:"tags"`

	OwnerId string `json:"owner_id"`

	LangCode string `json:"lang_code"`

	CurVersionNo uint   `json:"cur_version_no"`
	CurVersionId string `json:"cur_version_id,omitempty"`

	ProcessStatus int `json:"process_status"`
	Status        int `json:"status"`

	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`

	CurrentVersion *DocumentVersionData `json:"current_version,omitempty"`
	FileObject     *FileObjectData      `json:"file_object,omitempty"`
	LatestJob      *IngestJobData       `json:"latest_job,omitempty"`
}

type ListDocumentsRespData struct {
	List []*DocumentData `json:"list"`

	Total int64 `json:"total"`
}

type GetDocumentRespData struct {
	DocumentId string `json:"document_id"`

	KnowledgeBaseId string `json:"knowledge_base_id,omitempty"`
	ChatSessionId   string `json:"chat_session_id,omitempty"`

	ScopeType  int `json:"scope_type"`
	SourceType int `json:"source_type"`

	Title   string   `json:"title"`
	Summary string   `json:"summary"`
	Tags    []string `json:"tags"`

	OwnerId string `json:"owner_id"`

	LangCode string `json:"lang_code"`

	CurVersionNo uint   `json:"cur_version_no"`
	CurVersionId string `json:"cur_version_id,omitempty"`

	ProcessStatus int `json:"process_status"`
	Status        int `json:"status"`

	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`

	CurrentVersion *DocumentVersionData `json:"current_version,omitempty"`
	FileObject     *FileObjectData      `json:"file_object,omitempty"`
	LatestJob      *IngestJobData       `json:"latest_job,omitempty"`
}

type CreateDocumentRequest struct {
	KnowledgeBaseId string `json:"knowledge_base_id"`
	ChatSessionId   string `json:"chat_session_id"`

	ScopeType  int `json:"scope_type"`
	SourceType int `json:"source_type"`

	Title   string   `json:"title"`
	Summary string   `json:"summary"`
	Tags    []string `json:"tags"`

	OwnerId string `json:"owner_id"`

	LangCode string `json:"lang_code"`

	ParseStrategy int    `json:"parse_strategy"`
	FileObjectId  string `json:"file_object_id"`
}

type CreateDocumentRespData struct {
	DocumentId string `json:"document_id"`

	VersionId    string `json:"version_id"`
	FileObjectId string `json:"file_object_id"`
	JobId        string `json:"job_id"`

	BucketName string `json:"bucket_name"`
	ObjectKey  string `json:"object_key"`

	ProcessStatus int `json:"process_status"`
	Status        int `json:"status"`
}

type UpdateDocumentRequest struct {
	Title string `json:"title"`

	Summary string   `json:"summary"`
	Tags    []string `json:"tags"`

	LangCode string `json:"lang_code"`

	Status int `json:"status"`
}

type FileObjectData struct {
	FileObjectId string `json:"file_object_id"`

	BucketName string `json:"bucket_name"`
	ObjectKey  string `json:"object_key"`

	FileName string `json:"file_name"`
	MimeType string `json:"mime_type"`
	FileExt  string `json:"file_ext"`

	SizeBytes uint64 `json:"size_bytes"`
	Sha256    string `json:"sha256"`

	StorageProvider int `json:"storage_provider"`

	CreatedAt string `json:"created_at"`
}

type DocumentVersionData struct {
	VersionId string `json:"version_id"`

	DocumentId string `json:"document_id"`
	VersionNo  uint   `json:"version_no"`

	FileObjectId string `json:"file_object_id"`

	ParseStrategy int `json:"parse_strategy"`
	ParserType    int `json:"parser_type"`

	ContentSha256 string `json:"content_sha256"`

	PageCount  uint `json:"page_count"`
	TokenCount uint `json:"token_count"`
	ChunkCount uint `json:"chunk_count"`

	ParseStatus int    `json:"parse_status"`
	ParseError  string `json:"parse_error,omitempty"`

	OcrStatus int    `json:"ocr_status"`
	OcrError  string `json:"ocr_error,omitempty"`

	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type IngestJobData struct {
	JobId string `json:"job_id"`

	DocumentId string `json:"document_id"`
	VersionId  string `json:"version_id"`

	JobType   int `json:"job_type"`
	JobStatus int `json:"job_status"`

	RetryCount uint `json:"retry_count"`

	WorkerName   string `json:"worker_name,omitempty"`
	ErrorMessage string `json:"error_message,omitempty"`
	Payload      string `json:"payload,omitempty"`

	StartedAt  string `json:"started_at,omitempty"`
	FinishedAt string `json:"finished_at,omitempty"`

	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}
