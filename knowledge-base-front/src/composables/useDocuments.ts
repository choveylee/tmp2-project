import { computed, onMounted, reactive, shallowRef, watch } from 'vue'

import {
  createDocument,
  deleteDocument,
  getDocument,
  listDocuments,
  listKnowledgeBases,
  updateDocument,
  uploadFile,
} from '@/api'
import type {
  CreateDocumentRequest,
  DocumentData,
  GetDocumentRespData,
  KnowledgeBaseData,
} from '@/types/api'
import type { DocumentFormDraft } from '@/types/forms'
import { formatErrorMessage } from '@/utils/errors'
import { formatFileSize } from '@/utils/number'
import { splitTagsText } from '@/utils/tags'

const knowledgeBaseOptionsPageSize = 100
const documentScopeTypeKnowledge = 0
const documentScopeTypeAttachment = 1
const maxDocumentUploadSize = 500 * 1024 * 1024

function createDraft(defaultKnowledgeBaseId = ''): DocumentFormDraft {
  return {
    knowledgeBaseId: defaultKnowledgeBaseId,
    chatSessionId: '',
    scopeType: 0,
    sourceType: 0,
    title: '',
    summary: '',
    tagsText: '',
    langCode: 'zh-CN',
    parseStrategy: 0,
    status: 1,
  }
}

export function useDocuments() {
  const filters = reactive({
    knowledgeBaseId: '',
    keyword: '',
    scopeType: '',
    sourceType: '',
    processStatus: '',
    status: '',
    pageNum: 1,
    pageSize: 10,
  })

  const items = shallowRef<DocumentData[]>([])
  const total = shallowRef(0)
  const loading = shallowRef(false)
  const error = shallowRef('')

  const selectedId = shallowRef('')
  const detail = shallowRef<GetDocumentRespData | null>(null)
  const detailLoading = shallowRef(false)
  const detailError = shallowRef('')

  const knowledgeBases = shallowRef<KnowledgeBaseData[]>([])
  const knowledgeBaseLoading = shallowRef(false)
  const knowledgeBaseError = shallowRef('')

  const isFormOpen = shallowRef(false)
  const formMode = shallowRef<'create' | 'edit'>('create')
  const saving = shallowRef(false)
  const uploading = shallowRef(false)
  const savingError = shallowRef('')
  const draft = reactive<DocumentFormDraft>(createDraft())
  let draftBaseline = createDraft()
  const selectedFile = shallowRef<File | null>(null)
  let detailRequestSeq = 0

  const selectedItem = computed(
    () => detail.value ?? items.value.find((item) => item.document_id === selectedId.value) ?? null,
  )

  const selectedFileLabel = computed(() => {
    if (!selectedFile.value) {
      return ''
    }

    return `${selectedFile.value.name} (${formatFileSize(selectedFile.value.size)})`
  })

  const pageCount = computed(() => {
    const count = Math.ceil(total.value / filters.pageSize)
    return Math.max(1, count || 1)
  })

  const knowledgeBaseOptions = computed(() => {
    return knowledgeBases.value.map((item) => ({
      value: item.knowledge_base_id,
      label: `${item.code} · ${item.name}`,
    }))
  })

  function applyDraft(next: DocumentFormDraft) {
    draft.knowledgeBaseId = next.knowledgeBaseId
    draft.chatSessionId = next.chatSessionId
    draft.scopeType = next.scopeType
    draft.sourceType = next.sourceType
    draft.title = next.title
    draft.summary = next.summary
    draft.tagsText = next.tagsText
    draft.langCode = next.langCode
    draft.parseStrategy = next.parseStrategy
    draft.status = next.status
  }

  function resetDraft(defaultKnowledgeBaseId = '') {
    applyDraft(createDraft(defaultKnowledgeBaseId))
    selectedFile.value = null
  }

  async function loadKnowledgeBaseOptions() {
    knowledgeBaseLoading.value = true
    knowledgeBaseError.value = ''

    try {
      const loadedKnowledgeBases: KnowledgeBaseData[] = []
      let pageNum = 1
      let total = 0

      while (true) {
        const response = await listKnowledgeBases({
          keyword: '',
          visible: '',
          status: '',
          pageNum,
          pageSize: knowledgeBaseOptionsPageSize,
        })

        loadedKnowledgeBases.push(...(response.list || []))
        total = response.total || 0

        if (loadedKnowledgeBases.length >= total || (response.list || []).length === 0) {
          break
        }

        pageNum += 1
      }

      knowledgeBases.value = loadedKnowledgeBases

      if (!draft.knowledgeBaseId && knowledgeBases.value.length > 0) {
        draft.knowledgeBaseId = knowledgeBases.value[0].knowledge_base_id
      }
    } catch (err) {
      knowledgeBaseError.value = formatErrorMessage(err)
      knowledgeBases.value = []
    } finally {
      knowledgeBaseLoading.value = false
    }
  }

  async function loadList() {
    loading.value = true
    error.value = ''

    try {
      const response = await listDocuments({
        knowledgeBaseId: filters.knowledgeBaseId,
        keyword: filters.keyword,
        scopeType: filters.scopeType,
        sourceType: filters.sourceType,
        processStatus: filters.processStatus,
        status: filters.status,
        pageNum: filters.pageNum,
        pageSize: filters.pageSize,
      })

      items.value = response.list || []
      total.value = response.total || 0

      if (items.value.length === 0) {
        selectedId.value = ''
        detail.value = null
        return
      }

      const selectedExists = items.value.some((item) => item.document_id === selectedId.value)
      if (!selectedId.value || !selectedExists) {
        selectedId.value = items.value[0].document_id
      }
    } catch (err) {
      error.value = formatErrorMessage(err)
      items.value = []
      total.value = 0
      selectedId.value = ''
      detail.value = null
    } finally {
      loading.value = false
    }
  }

  async function loadDetail(documentId: string) {
    const requestSeq = ++detailRequestSeq

    if (!documentId) {
      detail.value = null
      detailError.value = ''
      detailLoading.value = false
      return
    }

    detailLoading.value = true
    detailError.value = ''

    try {
      const nextDetail = await getDocument(documentId)
      if (requestSeq !== detailRequestSeq) {
        return
      }

      detail.value = nextDetail
    } catch (err) {
      if (requestSeq !== detailRequestSeq) {
        return
      }

      detail.value = null
      detailError.value = formatErrorMessage(err)
    } finally {
      if (requestSeq !== detailRequestSeq) {
        return
      }

      detailLoading.value = false
    }
  }

  watch(
    selectedId,
    (documentId) => {
      void loadDetail(documentId)
    },
    { immediate: true },
  )

  function selectItem(documentId: string) {
    selectedId.value = documentId
  }

  function startCreate() {
    formMode.value = 'create'
    savingError.value = ''
    uploading.value = false
    draftBaseline = createDraft(filters.knowledgeBaseId || knowledgeBaseOptions.value[0]?.value || '')
    resetDraft(draftBaseline.knowledgeBaseId)
    isFormOpen.value = true
  }

  function startEdit(item: DocumentData) {
    formMode.value = 'edit'
    savingError.value = ''
    uploading.value = false
    selectedId.value = item.document_id
    detail.value = item as GetDocumentRespData
    detailLoading.value = false
    detailError.value = ''
    draftBaseline = {
      knowledgeBaseId: item.knowledge_base_id || '',
      chatSessionId: item.chat_session_id || '',
      scopeType: item.scope_type,
      sourceType: item.source_type,
      title: item.title,
      summary: item.summary,
      tagsText: (item.tags || []).join(', '),
      langCode: item.lang_code,
      parseStrategy: item.current_version?.parse_strategy ?? 0,
      status: item.status,
    }
    applyDraft(draftBaseline)
    selectedFile.value = null
    isFormOpen.value = true
  }

  function closeForm() {
    isFormOpen.value = false
    savingError.value = ''
    uploading.value = false
    selectedFile.value = null
  }

  async function submitForm() {
    saving.value = true
    uploading.value = false
    savingError.value = ''

    try {
      if (formMode.value === 'create') {
        if (draft.scopeType === documentScopeTypeKnowledge) {
          if (!draft.knowledgeBaseId.trim()) {
            throw new Error('Knowledge base is required')
          }

          const knowledgeBaseExists = knowledgeBaseOptions.value.some(
            (item) => item.value === draft.knowledgeBaseId.trim(),
          )
          if (!knowledgeBaseExists) {
            throw new Error('Selected knowledge base is invalid')
          }
        }

        if (draft.scopeType === documentScopeTypeAttachment && !draft.chatSessionId.trim()) {
          throw new Error('Chat session id is required')
        }

        const file = selectedFile.value
        if (!file) {
          throw new Error('Document file is required')
        }
        if (file.size <= 0) {
          throw new Error('Document file is empty')
        }
        if (file.size > maxDocumentUploadSize) {
          throw new Error('Document file must be 500 MB or smaller')
        }

        uploading.value = true
        const uploadResp = await uploadFile(file)

        const payload = {
          knowledge_base_id: draft.knowledgeBaseId.trim(),
          chat_session_id: draft.chatSessionId.trim(),
          scope_type: draft.scopeType,
          source_type: draft.sourceType,
          title: draft.title.trim(),
          summary: draft.summary.trim(),
          tags: splitTagsText(draft.tagsText),
          owner_id: '',
          lang_code: draft.langCode.trim(),
          parse_strategy: draft.parseStrategy,
          file_object_id: uploadResp.file_object_id,
        } satisfies CreateDocumentRequest

        const response = await createDocument(payload)
        selectedId.value = response.document_id
        draftBaseline = createDraft(draft.knowledgeBaseId)
        detail.value = null
      } else {
        const documentId = selectedId.value
        if (!documentId) {
          throw new Error('Document not selected')
        }

        await updateDocument(documentId, {
          title: draft.title.trim(),
          summary: draft.summary.trim(),
          tags: splitTagsText(draft.tagsText),
          lang_code: draft.langCode.trim(),
          status: draft.status,
        })

        await loadDetail(documentId)
      }

      await loadList()
      closeForm()
    } catch (err) {
      savingError.value = formatErrorMessage(err)
    } finally {
      saving.value = false
      uploading.value = false
    }
  }

  async function removeItem(document: DocumentData) {
    if (!window.confirm(`Delete document ${document.title}?`)) {
      return
    }

    try {
      await deleteDocument(document.document_id)
      if (selectedId.value === document.document_id) {
        selectedId.value = ''
      }
      await loadList()
    } catch (err) {
      error.value = formatErrorMessage(err)
    }
  }

  function clearFilters() {
    filters.knowledgeBaseId = ''
    filters.keyword = ''
    filters.scopeType = ''
    filters.sourceType = ''
    filters.processStatus = ''
    filters.status = ''
    filters.pageNum = 1
  }

  function search() {
    filters.pageNum = 1
    void loadList()
  }

  function refresh() {
    void loadList()
  }

  onMounted(() => {
    void loadKnowledgeBaseOptions()
    void loadList()
  })

  return {
    filters,
    items,
    total,
    loading,
    error,
    selectedId,
    selectedItem,
    detailLoading,
    detailError,
    knowledgeBases,
    knowledgeBaseOptions,
    knowledgeBaseLoading,
    knowledgeBaseError,
    isFormOpen,
    formMode,
    saving,
    savingError,
    uploading,
    draft,
    selectedFile,
    selectedFileLabel,
    pageCount,
    loadList,
    selectItem,
    startCreate,
    startEdit,
    closeForm,
    submitForm,
    removeItem,
    clearFilters,
    search,
    refresh,
    resetDraft,
  }
}
