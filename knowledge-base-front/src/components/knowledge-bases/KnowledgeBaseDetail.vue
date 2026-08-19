<script setup lang="ts">
import type { KnowledgeBaseData } from '@/types/api'
import { formatDateTime } from '@/utils/time'
import { labelFromOptions, knowledgeBaseStatusOptions, knowledgeBaseVisibleOptions } from '@/constants/lookups'

defineProps<{
  item: KnowledgeBaseData | null
  loading: boolean
  error: string
}>()
</script>

<template>
  <section class="panel stack">
    <div class="section-head">
      <div>
        <h3 class="section-title">Details</h3>
        <p class="muted">Selected knowledge base snapshot.</p>
      </div>
    </div>

    <div v-if="loading" class="notice">Loading detail...</div>
    <div v-else-if="error" class="notice error">{{ error }}</div>
    <div v-else-if="!item" class="empty-state">No knowledge base selected.</div>

    <div v-else class="detail-grid">
      <div class="summary-row">
        <span class="badge soft">{{ item.knowledge_base_id }}</span>
        <span class="badge soft">{{ item.code }}</span>
        <span :class="['badge', item.status === 1 ? 'good' : 'bad']">
          {{ labelFromOptions(knowledgeBaseStatusOptions, item.status) }}
        </span>
      </div>

      <div class="mini-grid">
        <div class="field">
          <span>Name</span>
          <strong>{{ item.name }}</strong>
        </div>
        <div class="field">
          <span>Owner</span>
          <strong>{{ item.owner_id || '-' }}</strong>
        </div>
        <div class="field">
          <span>Visible</span>
          <strong>{{ labelFromOptions(knowledgeBaseVisibleOptions, item.visible) }}</strong>
        </div>
        <div class="field">
          <span>Updated</span>
          <strong>{{ formatDateTime(item.updated_at) }}</strong>
        </div>
      </div>

      <div class="field">
        <span>Description</span>
        <div class="detail-text">{{ item.description || '-' }}</div>
      </div>
    </div>
  </section>
</template>
