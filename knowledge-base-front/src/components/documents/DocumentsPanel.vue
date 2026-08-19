<script setup lang="ts">
import { computed } from 'vue'

import DocumentDetail from './DocumentDetail.vue'
import DocumentModal from './DocumentModal.vue'
import DocumentTable from './DocumentTable.vue'
import { useDocuments } from '@/composables/useDocuments'
import { documentProcessStatusOptions, documentScopeTypeOptions, documentSourceTypeOptions, documentStatusOptions } from '@/constants/lookups'

const {
  filters,
  items,
  total,
  loading,
  error,
  knowledgeBaseError,
  selectedId,
  selectedItem,
  detailLoading,
  detailError,
  knowledgeBaseOptions,
  isFormOpen,
  formMode,
  saving,
  savingError,
  draft,
  selectedFile,
  pageCount,
  selectItem,
  startCreate,
  startEdit,
  closeForm,
  submitForm,
  removeItem,
  clearFilters,
  search,
  refresh,
} = useDocuments()

const selectedSummary = computed(() => {
  if (!selectedItem.value) {
    return 'No row selected'
  }

  return `Selected: ${selectedItem.value.title}`
})
</script>

<template>
  <section class="page-grid kb-page">
    <section class="page-section stack kb-hero">
      <div class="section-head">
        <div class="stack-tight">
          <p class="eyebrow">Document library</p>
          <h2 class="section-title">Documents</h2>
          <p class="muted">Upload files, tune metadata, and inspect parse and ingest state.</p>
        </div>
        <div class="toolbar">
          <button type="button" class="primary" @click="startCreate">New Document</button>
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
          <span>Knowledge Base</span>
          <select v-model="filters.knowledgeBaseId">
            <option value="">All</option>
            <option v-for="item in knowledgeBaseOptions" :key="item.value" :value="item.value">
              {{ item.label }}
            </option>
          </select>
        </label>

        <label class="field">
          <span>Keyword</span>
          <input v-model="filters.keyword" type="text" placeholder="Title or summary" />
        </label>

        <label class="field">
          <span>Scope</span>
          <select v-model="filters.scopeType">
            <option value="">All</option>
            <option v-for="item in documentScopeTypeOptions" :key="item.value" :value="String(item.value)">
              {{ item.label }}
            </option>
          </select>
        </label>

        <label class="field">
          <span>Source</span>
          <select v-model="filters.sourceType">
            <option value="">All</option>
            <option v-for="item in documentSourceTypeOptions" :key="item.value" :value="String(item.value)">
              {{ item.label }}
            </option>
          </select>
        </label>

        <label class="field">
          <span>Process</span>
          <select v-model="filters.processStatus">
            <option value="">All</option>
            <option v-for="item in documentProcessStatusOptions" :key="item.value" :value="String(item.value)">
              {{ item.label }}
            </option>
          </select>
        </label>

        <label class="field">
          <span>Status</span>
          <select v-model="filters.status">
            <option value="">All</option>
            <option v-for="item in documentStatusOptions" :key="item.value" :value="String(item.value)">
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
    <div v-if="knowledgeBaseError" class="notice error">{{ knowledgeBaseError }}</div>

    <div class="page-columns">
      <DocumentTable
        :items="items"
        :loading="loading"
        :selected-id="selectedId"
        :total="total"
        :page-num="filters.pageNum"
        :page-count="pageCount"
        @select="selectItem"
        @edit="startEdit"
        @remove="removeItem"
        @page="filters.pageNum = $event"
      />

      <div class="stack">
        <DocumentDetail :item="selectedItem" :loading="detailLoading" :error="detailError" />
      </div>
    </div>

    <DocumentModal
      :open="isFormOpen"
      :mode="formMode"
      :saving="saving"
      :error="savingError"
      :knowledge-base-id="draft.knowledgeBaseId"
      :chat-session-id="draft.chatSessionId"
      :scope-type="draft.scopeType"
      :source-type="draft.sourceType"
      :title="draft.title"
      :summary="draft.summary"
      :tags-text="draft.tagsText"
      :lang-code="draft.langCode"
      :parse-strategy="draft.parseStrategy"
      :status="draft.status"
      :file-name="selectedFile?.name || ''"
      :file-required="formMode === 'create'"
      :knowledge-base-options="knowledgeBaseOptions"
      @close="closeForm"
      @submit="submitForm"
      @reset="startCreate"
      @update:knowledgeBaseId="draft.knowledgeBaseId = $event"
      @update:chatSessionId="draft.chatSessionId = $event"
      @update:scopeType="draft.scopeType = $event"
      @update:sourceType="draft.sourceType = $event"
      @update:title="draft.title = $event"
      @update:summary="draft.summary = $event"
      @update:tagsText="draft.tagsText = $event"
      @update:langCode="draft.langCode = $event"
      @update:parseStrategy="draft.parseStrategy = $event"
      @update:status="draft.status = $event"
      @update:file="selectedFile = $event"
    />
  </section>
</template>
