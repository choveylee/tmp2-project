<script setup lang="ts">
import type { DocumentData } from '@/types/api'
import { formatDateTime } from '@/utils/time'
import { labelFromOptions, documentProcessStatusOptions, documentScopeTypeOptions, documentSourceTypeOptions, documentStatusOptions } from '@/constants/lookups'

const props = defineProps<{
  items: DocumentData[]
  loading: boolean
  selectedId: string
  total: number
  pageNum: number
  pageCount: number
}>()

const emit = defineEmits<{
  (event: 'select', documentId: string): void
  (event: 'edit', item: DocumentData): void
  (event: 'remove', item: DocumentData): void
  (event: 'page', pageNum: number): void
}>()
</script>

<template>
  <section class="panel stack">
    <div class="section-head">
      <div>
        <h3 class="section-title">Documents</h3>
        <p class="muted">{{ props.total }} records</p>
      </div>
      <div class="toolbar">
        <button type="button" :disabled="props.pageNum <= 1" @click="emit('page', props.pageNum - 1)">Prev</button>
        <button type="button" :disabled="props.pageNum >= props.pageCount" @click="emit('page', props.pageNum + 1)">Next</button>
      </div>
    </div>

    <div v-if="props.loading" class="notice">Loading documents...</div>
    <div v-else-if="props.items.length === 0" class="empty-state">No documents found.</div>
    <div v-else class="table-wrap">
      <table>
        <thead>
          <tr>
            <th>Title</th>
            <th>Knowledge Base</th>
            <th>Scope</th>
            <th>Source</th>
            <th>Process</th>
            <th>Status</th>
            <th>Updated</th>
            <th></th>
          </tr>
        </thead>
        <tbody>
          <tr
            v-for="item in props.items"
            :key="item.document_id"
            :class="{ 'row-selected': item.document_id === props.selectedId }"
            @click="emit('select', item.document_id)"
          >
            <td>
              <div class="stack-tight">
                <strong>{{ item.title }}</strong>
                <span class="muted mono">{{ item.document_id }}</span>
              </div>
            </td>
            <td>{{ item.knowledge_base_id || '-' }}</td>
            <td>{{ labelFromOptions(documentScopeTypeOptions, item.scope_type) }}</td>
            <td>{{ labelFromOptions(documentSourceTypeOptions, item.source_type) }}</td>
            <td>{{ labelFromOptions(documentProcessStatusOptions, item.process_status) }}</td>
            <td>
              <span :class="['badge', item.status === 1 ? 'good' : 'bad']">
                {{ labelFromOptions(documentStatusOptions, item.status) }}
              </span>
            </td>
            <td>{{ formatDateTime(item.updated_at) }}</td>
            <td>
              <div class="toolbar toolbar-tight" @click.stop>
                <button type="button" @click="emit('edit', item)">Edit</button>
                <button type="button" class="danger" @click="emit('remove', item)">Delete</button>
              </div>
            </td>
          </tr>
        </tbody>
      </table>
    </div>
  </section>
</template>
