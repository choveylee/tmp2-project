<script setup lang="ts">
import { computed } from 'vue'

import KnowledgeBaseDetail from './KnowledgeBaseDetail.vue'
import KnowledgeBaseModal from './KnowledgeBaseModal.vue'
import KnowledgeBaseTable from './KnowledgeBaseTable.vue'
import { useKnowledgeBases } from '@/composables/useKnowledgeBases'
import {
  knowledgeBaseStatusOptions,
  knowledgeBaseVisibleOptions,
  listPageSizeOptions,
} from '@/constants/lookups'

const {
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
  selectItem,
  startCreate,
  startEdit,
  closeForm,
  resetDraft,
  submitForm,
  removeItem,
  clearFilters,
  search,
  refresh,
} = useKnowledgeBases()

function handlePage(pageNum: number) {
  filters.pageNum = pageNum
  void refresh()
}

const selectedSummary = computed(() => {
  if (!selectedItem.value) {
    return 'No row selected'
  }

  return `Selected: ${selectedItem.value.code}`
})
</script>

<template>
  <section class="page-grid kb-page">
    <section class="page-section stack kb-hero">
      <div class="section-head">
        <div class="stack-tight">
          <p class="eyebrow">Knowledge catalog</p>
          <h2 class="section-title">Knowledge Bases</h2>
          <p class="muted">Manage the libraries that power documents, parsing, and ingest jobs.</p>
        </div>
        <div class="toolbar">
          <button type="button" class="primary" @click="startCreate">New Knowledge Base</button>
          <button type="button" @click="refresh">Refresh</button>
        </div>
      </div>

      <div class="summary-row">
        <span class="badge soft">{{ total }} records</span>
        <span class="badge soft">Page {{ filters.pageNum }} / {{ pageCount }}</span>
        <span class="badge soft">{{ selectedSummary }}</span>
      </div>

      <div class="toolbar filter-grid">
        <label class="field">
          <span>Keyword</span>
          <input v-model="filters.keyword" type="text" placeholder="Search code or name" />
        </label>

        <label class="field">
          <span>Visible</span>
          <select v-model="filters.visible">
            <option value="">All</option>
            <option v-for="item in knowledgeBaseVisibleOptions" :key="item.value" :value="String(item.value)">
              {{ item.label }}
            </option>
          </select>
        </label>

        <label class="field">
          <span>Status</span>
          <select v-model="filters.status">
            <option value="">All</option>
            <option v-for="item in knowledgeBaseStatusOptions" :key="item.value" :value="String(item.value)">
              {{ item.label }}
            </option>
          </select>
        </label>

        <label class="field">
          <span>Page Size</span>
          <select v-model.number="filters.pageSize" @change="search">
            <option v-for="item in listPageSizeOptions" :key="item.value" :value="item.value">
              {{ item.label }}
            </option>
          </select>
        </label>

        <div class="toolbar filter-actions">
          <button type="button" class="primary" @click="search">Search</button>
          <button type="button" @click="clearFilters">Clear</button>
        </div>
      </div>
    </section>

    <div v-if="error" class="notice error">{{ error }}</div>

    <div class="page-columns">
      <KnowledgeBaseTable
        :items="items"
        :loading="loading"
        :selected-id="selectedId"
        :total="total"
        :page-num="filters.pageNum"
        :page-count="pageCount"
        @select="selectItem"
        @edit="startEdit"
        @remove="removeItem"
        @page="handlePage"
      />

      <div class="stack">
        <KnowledgeBaseDetail :item="selectedItem" :loading="detailLoading" :error="detailError" />
      </div>
    </div>

    <KnowledgeBaseModal
      :open="isFormOpen"
      :mode="formMode"
      :saving="saving"
      :error="savingError"
      :code="draft.code"
      :name="draft.name"
      :description="draft.description"
      :visible="draft.visible"
      :status="draft.status"
      @close="closeForm"
      @submit="submitForm"
      @reset="resetDraft"
      @update:code="draft.code = $event"
      @update:name="draft.name = $event"
      @update:description="draft.description = $event"
      @update:visible="draft.visible = $event"
      @update:status="draft.status = $event"
    />
  </section>
</template>
