<script setup lang="ts">
import type { GetDocumentRespData } from '@/types/api'
import { formatDateTime } from '@/utils/time'
import {
  documentParseStrategyOptions,
  documentProcessStatusOptions,
  documentScopeTypeOptions,
  documentSourceTypeOptions,
  documentStatusOptions,
  ingestJobStatusOptions,
  ingestJobTypeOptions,
  labelFromOptions,
} from '@/constants/lookups'

defineProps<{
  item: GetDocumentRespData | null
  loading: boolean
  error: string
}>()
</script>

<template>
  <section class="panel stack">
    <div class="section-head">
      <div>
        <h3 class="section-title">Document Detail</h3>
        <p class="muted">Selected document snapshot and latest processing state.</p>
      </div>
    </div>

    <div v-if="loading" class="notice">Loading detail...</div>
    <div v-else-if="error" class="notice error">{{ error }}</div>
    <div v-else-if="!item" class="empty-state">No document selected.</div>

    <div v-else class="detail-grid">
      <div class="summary-row">
        <span class="badge soft">{{ item.document_id }}</span>
        <span class="badge soft">{{ item.title }}</span>
        <span :class="['badge', item.status === 1 ? 'good' : 'bad']">
          {{ labelFromOptions(documentStatusOptions, item.status) }}
        </span>
      </div>

      <div class="mini-grid">
        <div class="field">
          <span>Knowledge Base</span>
          <strong>{{ item.knowledge_base_id || '-' }}</strong>
        </div>
        <div class="field">
          <span>Scope</span>
          <strong>{{ labelFromOptions(documentScopeTypeOptions, item.scope_type) }}</strong>
        </div>
        <div class="field">
          <span>Source</span>
          <strong>{{ labelFromOptions(documentSourceTypeOptions, item.source_type) }}</strong>
        </div>
        <div class="field">
          <span>Process</span>
          <strong>{{ labelFromOptions(documentProcessStatusOptions, item.process_status) }}</strong>
        </div>
      </div>

      <div class="mini-grid">
        <div class="field">
          <span>Lang</span>
          <strong>{{ item.lang_code || '-' }}</strong>
        </div>
        <div class="field">
          <span>Version</span>
          <strong>{{ item.cur_version_no || '-' }}</strong>
        </div>
        <div class="field">
          <span>Tags</span>
          <strong>{{ (item.tags || []).join(', ') || '-' }}</strong>
        </div>
        <div class="field">
          <span>Updated</span>
          <strong>{{ formatDateTime(item.updated_at) }}</strong>
        </div>
      </div>

      <div class="field">
        <span>Summary</span>
        <div class="detail-text">{{ item.summary || '-' }}</div>
      </div>

      <div v-if="item.file_object" class="panel-soft detail-grid">
        <strong>File Object</strong>
        <div class="mini-grid">
          <div class="field"><span>Name</span><strong>{{ item.file_object.file_name }}</strong></div>
          <div class="field"><span>Size</span><strong>{{ item.file_object.size_bytes }}</strong></div>
          <div class="field"><span>Bucket</span><strong>{{ item.file_object.bucket_name }}</strong></div>
          <div class="field"><span>Key</span><strong class="mono">{{ item.file_object.object_key }}</strong></div>
        </div>
      </div>

      <div v-if="item.current_version" class="panel-soft detail-grid">
        <strong>Current Version</strong>
        <div class="mini-grid">
          <div class="field"><span>Parse Strategy</span><strong>{{ labelFromOptions(documentParseStrategyOptions, item.current_version.parse_strategy) }}</strong></div>
          <div class="field"><span>Parser Type</span><strong>{{ item.current_version.parser_type }}</strong></div>
          <div class="field"><span>Page Count</span><strong>{{ item.current_version.page_count }}</strong></div>
          <div class="field"><span>Chunk Count</span><strong>{{ item.current_version.chunk_count }}</strong></div>
        </div>
      </div>

      <div v-if="item.latest_job" class="panel-soft detail-grid">
        <strong>Latest Job</strong>
        <div class="mini-grid">
          <div class="field"><span>Job ID</span><strong class="mono">{{ item.latest_job.job_id }}</strong></div>
          <div class="field">
            <span>Type</span>
            <strong>{{ labelFromOptions(ingestJobTypeOptions, item.latest_job.job_type) }}</strong>
          </div>
          <div class="field">
            <span>Status</span>
            <strong>{{ labelFromOptions(ingestJobStatusOptions, item.latest_job.job_status) }}</strong>
          </div>
          <div class="field"><span>Retry</span><strong>{{ item.latest_job.retry_count }}</strong></div>
        </div>
      </div>
    </div>
  </section>
</template>
