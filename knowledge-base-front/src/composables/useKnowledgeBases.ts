import { computed, onMounted, reactive, shallowRef, watch } from 'vue'

import {
  createKnowledgeBase,
  deleteKnowledgeBase,
  getKnowledgeBase,
  listKnowledgeBases,
  updateKnowledgeBase,
} from '@/api'
import type { KnowledgeBaseData } from '@/types/api'
import type { KnowledgeBaseFormDraft } from '@/types/forms'
import { formatErrorMessage } from '@/utils/errors'

function createDraft(): KnowledgeBaseFormDraft {
  return {
    code: '',
    name: '',
    description: '',
    visible: 1,
    status: 1,
  }
}

export function useKnowledgeBases() {
  const filters = reactive({
    keyword: '',
    visible: '',
    status: '',
    pageNum: 1,
    pageSize: 10,
  })

  const items = shallowRef<KnowledgeBaseData[]>([])
  const total = shallowRef(0)
  const loading = shallowRef(false)
  const error = shallowRef('')

  const selectedId = shallowRef('')
  const detail = shallowRef<KnowledgeBaseData | null>(null)
  const detailLoading = shallowRef(false)
  const detailError = shallowRef('')

  const isFormOpen = shallowRef(false)
  const formMode = shallowRef<'create' | 'edit'>('create')
  const saving = shallowRef(false)
  const savingError = shallowRef('')
  const draft = reactive<KnowledgeBaseFormDraft>(createDraft())
  let draftBaseline = createDraft()
  let detailRequestSeq = 0

  const selectedItem = computed(() => {
    return detail.value ?? items.value.find((item) => item.knowledge_base_id === selectedId.value) ?? null
  })

  const pageCount = computed(() => {
    const count = Math.ceil(total.value / filters.pageSize)
    return Math.max(1, count || 1)
  })

  function applyDraft(next: KnowledgeBaseFormDraft) {
    draft.code = next.code
    draft.name = next.name
    draft.description = next.description
    draft.visible = next.visible
    draft.status = next.status
  }

  function resetDraft() {
    applyDraft(draftBaseline)
  }

  async function loadList() {
    loading.value = true
    error.value = ''

    try {
      const response = await listKnowledgeBases({
        keyword: filters.keyword,
        visible: filters.visible,
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

      const selectedExists = items.value.some((item) => item.knowledge_base_id === selectedId.value)
      if (!selectedId.value || !selectedExists) {
        selectedId.value = items.value[0].knowledge_base_id
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

  async function loadDetail(knowledgeBaseId: string) {
    const requestSeq = ++detailRequestSeq

    if (!knowledgeBaseId) {
      detail.value = null
      detailError.value = ''
      detailLoading.value = false
      return
    }

    detailLoading.value = true
    detailError.value = ''

    try {
      const nextDetail = await getKnowledgeBase(knowledgeBaseId)
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
    (knowledgeBaseId) => {
      void loadDetail(knowledgeBaseId)
    },
    { immediate: true },
  )

  function selectItem(knowledgeBaseId: string) {
    selectedId.value = knowledgeBaseId
  }

  function startCreate() {
    formMode.value = 'create'
    savingError.value = ''
    draftBaseline = createDraft()
    resetDraft()
    isFormOpen.value = true
  }

  function startEdit(item: KnowledgeBaseData) {
    formMode.value = 'edit'
    savingError.value = ''
    selectedId.value = item.knowledge_base_id
    detail.value = item
    detailLoading.value = false
    detailError.value = ''
    draftBaseline = {
      code: item.code,
      name: item.name,
      description: item.description,
      visible: item.visible,
      status: item.status,
    }
    resetDraft()
    isFormOpen.value = true
  }

  function closeForm() {
    isFormOpen.value = false
    savingError.value = ''
  }

  async function submitForm() {
    saving.value = true
    savingError.value = ''

    try {
      let targetId = selectedId.value

      if (formMode.value === 'create') {
        const response = await createKnowledgeBase({
          code: draft.code.trim(),
          name: draft.name.trim(),
          description: draft.description.trim(),
          visible: draft.visible,
          status: draft.status,
        })

        selectedId.value = response.knowledge_base_id
        targetId = response.knowledge_base_id
      } else {
        const knowledgeBaseId = selectedId.value
        if (!knowledgeBaseId) {
          throw new Error('Knowledge base not selected')
        }

        targetId = knowledgeBaseId
        await updateKnowledgeBase(knowledgeBaseId, {
          name: draft.name.trim(),
          description: draft.description.trim(),
          visible: draft.visible,
          status: draft.status,
        })
      }

      await loadList()

      if (formMode.value === 'edit') {
        await loadDetail(targetId)
      }

      closeForm()
    } catch (err) {
      savingError.value = formatErrorMessage(err)
    } finally {
      saving.value = false
    }
  }

  async function removeItem(knowledgeBase: KnowledgeBaseData) {
    if (!window.confirm(`Delete knowledge base ${knowledgeBase.name}?`)) {
      return
    }

    try {
      await deleteKnowledgeBase(knowledgeBase.knowledge_base_id)
      if (selectedId.value === knowledgeBase.knowledge_base_id) {
        selectedId.value = ''
      }
      await loadList()
    } catch (err) {
      error.value = formatErrorMessage(err)
    }
  }

  function clearFilters() {
    filters.keyword = ''
    filters.visible = ''
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
    isFormOpen,
    formMode,
    saving,
    savingError,
    draft,
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
