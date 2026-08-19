export interface KnowledgeBaseFormDraft {
  code: string
  name: string
  description: string
  visible: number
  status: number
}

export interface DocumentFormDraft {
  knowledgeBaseId: string
  chatSessionId: string
  scopeType: number
  sourceType: number
  title: string
  summary: string
  tagsText: string
  langCode: string
  parseStrategy: number
  status: number
}
