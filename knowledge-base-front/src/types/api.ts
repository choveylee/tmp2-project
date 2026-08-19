export interface ResponseEnvelope<T = unknown> {
  code: number
  message?: string
  detail?: string
  data?: T
  ex_data?: unknown
}

export interface KnowledgeBaseData {
  knowledge_base_id: string
  code: string
  name: string
  owner_id: string
  description: string
  visible: number
  status: number
  created_at: string
  updated_at: string
}

export interface ListKnowledgeBasesRespData {
  list: KnowledgeBaseData[]
  total: number
}

export interface GetKnowledgeBaseRespData {
  knowledge_base_id: string
  code: string
  name: string
  owner_id: string
  description: string
  visible: number
  status: number
  created_at: string
  updated_at: string
}

export interface CreateKnowledgeBaseRequest {
  code: string
  name: string
  description: string
  visible: number
  status: number
}

export interface CreateKnowledgeBaseRespData {
  knowledge_base_id: string
}

export interface UpdateKnowledgeBaseRequest {
  name: string
  description: string
  visible: number
  status: number
}

export interface DocumentData {
  document_id: string
  knowledge_base_id?: string
  chat_session_id?: string
  scope_type: number
  source_type: number
  title: string
  summary: string
  tags: string[]
  owner_id: string
  lang_code: string
  cur_version_no: number
  cur_version_id?: string
  process_status: number
  status: number
  created_at: string
  updated_at: string
  current_version?: DocumentVersionData
  file_object?: FileObjectData
  latest_job?: IngestJobData
}

export interface CreateDocumentRespData {
  document_id: string
  version_id: string
  file_object_id: string
  job_id: string
  bucket_name: string
  object_key: string
  process_status: number
  status: number
}

export interface UploadFileRespData {
  file_object_id: string
  bucket_name: string
  object_key: string
  file_name: string
  mime_type: string
  file_ext: string
  size_bytes: number
  sha256: string
  storage_provider: number
  created_at: string
}

export interface FileObjectData {
  file_object_id: string
  bucket_name: string
  object_key: string
  file_name: string
  mime_type: string
  file_ext: string
  size_bytes: number
  sha256: string
  storage_provider: number
  created_at: string
}

export interface DocumentVersionData {
  version_id: string
  document_id: string
  version_no: number
  file_object_id: string
  parse_strategy: number
  parser_type: number
  content_sha256: string
  page_count: number
  token_count: number
  chunk_count: number
  parse_status: number
  parse_error?: string
  ocr_status: number
  ocr_error?: string
  created_at: string
  updated_at: string
}

export interface IngestJobData {
  job_id: string
  document_id: string
  version_id: string
  job_type: number
  job_status: number
  retry_count: number
  worker_name?: string
  error_message?: string
  payload?: string
  started_at?: string
  finished_at?: string
  created_at: string
  updated_at: string
}

export interface ListDocumentsRespData {
  list: DocumentData[]
  total: number
}

export interface GetDocumentRespData {
  document_id: string
  knowledge_base_id?: string
  chat_session_id?: string
  scope_type: number
  source_type: number
  title: string
  summary: string
  tags: string[]
  owner_id: string
  lang_code: string
  cur_version_no: number
  cur_version_id?: string
  process_status: number
  status: number
  created_at: string
  updated_at: string
  current_version?: DocumentVersionData
  file_object?: FileObjectData
  latest_job?: IngestJobData
}

export interface UpdateDocumentRequest {
  title: string
  summary: string
  tags: string[]
  lang_code: string
  status: number
}

export interface CreateDocumentRequest {
  knowledge_base_id: string
  chat_session_id: string
  scope_type: number
  source_type: number
  title: string
  summary: string
  tags: string[]
  owner_id: string
  lang_code: string
  parse_strategy: number
  file_object_id: string
}

export interface UpdateDocumentRespData {
  ok: true
}

export interface HealthCheckRespData {
  status: string
  detail: string
}
