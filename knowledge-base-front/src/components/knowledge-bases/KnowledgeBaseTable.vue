<script setup lang="ts">
import type { KnowledgeBaseData } from '@/types/api'
import { formatDateTime } from '@/utils/time'
import { labelFromOptions, knowledgeBaseStatusOptions, knowledgeBaseVisibleOptions } from '@/constants/lookups'

const props = defineProps<{
  items: KnowledgeBaseData[]
  loading: boolean
  selectedId: string
  total: number
  pageNum: number
  pageCount: number
}>()

const emit = defineEmits<{
  (event: 'select', knowledgeBaseId: string): void
  (event: 'edit', item: KnowledgeBaseData): void
  (event: 'remove', item: KnowledgeBaseData): void
  (event: 'page', pageNum: number): void
}>()
</script>

<template>
  <section class="panel stack">
    <div class="section-head">
      <div>
        <h3 class="section-title">Knowledge Bases</h3>
        <p class="muted">{{ props.total }} records</p>
      </div>
      <div class="toolbar">
        <button type="button" :disabled="props.pageNum <= 1" @click="emit('page', props.pageNum - 1)">
          Prev
        </button>
        <button type="button" :disabled="props.pageNum >= props.pageCount" @click="emit('page', props.pageNum + 1)">
          Next
        </button>
      </div>
    </div>

    <div v-if="props.loading" class="notice">Loading knowledge bases...</div>
    <div v-else-if="props.items.length === 0" class="empty-state">No knowledge bases found.</div>
    <div v-else class="table-wrap">
      <table>
        <thead>
          <tr>
            <th>Code</th>
            <th>Name</th>
            <th>Visible</th>
            <th>Status</th>
            <th>Owner</th>
            <th>Updated</th>
            <th></th>
          </tr>
        </thead>
        <tbody>
          <tr
            v-for="item in props.items"
            :key="item.knowledge_base_id"
            :class="{ 'row-selected': item.knowledge_base_id === props.selectedId }"
            @click="emit('select', item.knowledge_base_id)"
          >
            <td class="mono">{{ item.code }}</td>
            <td>
              <div class="stack-tight">
                <strong>{{ item.name }}</strong>
                <span class="muted">{{ item.knowledge_base_id }}</span>
              </div>
            </td>
            <td>{{ labelFromOptions(knowledgeBaseVisibleOptions, item.visible) }}</td>
            <td>
              <span :class="['badge', item.status === 1 ? 'good' : 'bad']">
                {{ labelFromOptions(knowledgeBaseStatusOptions, item.status) }}
              </span>
            </td>
            <td>{{ item.owner_id || '-' }}</td>
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
