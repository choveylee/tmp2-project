import type {
  CreateDocumentRequest,
  CreateKnowledgeBaseRequest,
  CreateKnowledgeBaseRespData,
  CreateDocumentRespData,
  DocumentData,
  GetDocumentRespData,
  GetKnowledgeBaseRespData,
  HealthCheckRespData,
  KnowledgeBaseData,
  ListDocumentsRespData,
  ListKnowledgeBasesRespData,
  UpdateDocumentRequest,
  UpdateKnowledgeBaseRequest,
  UploadFileRespData,
} from '@/types/api'

import { requestFormData, requestJson } from './client'

export async function listKnowledgeBases(params: {
  keyword: string
  visible: string
  status: string
  pageNum: number
  pageSize: number
}) {
  const query = new URLSearchParams()
  if (params.keyword.trim()) query.set('keyword', params.keyword.trim())
  if (params.visible.trim()) query.set('visible', params.visible.trim())
  if (params.status.trim()) query.set('status', params.status.trim())
  query.set('page_num', String(params.pageNum))
  query.set('page_size', String(params.pageSize))

  return requestJson<ListKnowledgeBasesRespData>(`/api/v1/knowledge-bases?${query.toString()}`)
}

export async function getKnowledgeBase(knowledgeBaseId: string) {
  return requestJson<GetKnowledgeBaseRespData>(`/api/v1/knowledge-bases/${encodeURIComponent(knowledgeBaseId)}`)
}

export async function createKnowledgeBase(payload: CreateKnowledgeBaseRequest) {
  return requestJson<CreateKnowledgeBaseRespData>('/api/v1/knowledge-bases', {
    method: 'POST',
    body: payload,
  })
}

export async function updateKnowledgeBase(
  knowledgeBaseId: string,
  payload: UpdateKnowledgeBaseRequest,
) {
  return requestJson<void>(`/api/v1/knowledge-bases/${encodeURIComponent(knowledgeBaseId)}`, {
    method: 'PUT',
    body: payload,
  })
}

export async function deleteKnowledgeBase(knowledgeBaseId: string) {
  return requestJson<void>(`/api/v1/knowledge-bases/${encodeURIComponent(knowledgeBaseId)}`, {
    method: 'DELETE',
  })
}

export async function listDocuments(params: {
  knowledgeBaseId: string
  keyword: string
  scopeType: string
  sourceType: string
  processStatus: string
  status: string
  pageNum: number
  pageSize: number
}) {
  const query = new URLSearchParams()
  if (params.knowledgeBaseId.trim()) query.set('knowledge_base_id', params.knowledgeBaseId.trim())
  if (params.keyword.trim()) query.set('keyword', params.keyword.trim())
  if (params.scopeType.trim()) query.set('scope_type', params.scopeType.trim())
  if (params.sourceType.trim()) query.set('source_type', params.sourceType.trim())
  if (params.processStatus.trim()) query.set('process_status', params.processStatus.trim())
  if (params.status.trim()) query.set('status', params.status.trim())
  query.set('page_num', String(params.pageNum))
  query.set('page_size', String(params.pageSize))

  return requestJson<ListDocumentsRespData>(`/api/v1/documents?${query.toString()}`)
}

export async function getDocument(documentId: string) {
  return requestJson<GetDocumentRespData>(`/api/v1/documents/${encodeURIComponent(documentId)}`)
}

export async function uploadFile(file: File) {
  const formData = new FormData()
  formData.set('file', file)

  return requestFormData<UploadFileRespData>('/api/v1/files', formData)
}

export async function createDocument(payload: CreateDocumentRequest) {
  return requestJson<CreateDocumentRespData>('/api/v1/documents', {
    method: 'POST',
    body: payload,
  })
}

export async function updateDocument(documentId: string, payload: UpdateDocumentRequest) {
  return requestJson<void>(`/api/v1/documents/${encodeURIComponent(documentId)}`, {
    method: 'PUT',
    body: payload,
  })
}

export async function deleteDocument(documentId: string) {
  return requestJson<void>(`/api/v1/documents/${encodeURIComponent(documentId)}`, {
    method: 'DELETE',
  })
}

export async function getCpuCheck() {
  return requestJson<HealthCheckRespData>('/cpu-check')
}

export async function getRamCheck() {
  return requestJson<HealthCheckRespData>('/ram-check')
}
