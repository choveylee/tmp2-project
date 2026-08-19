<script setup lang="ts">
import {
  documentParseStrategyOptions,
  documentScopeTypeOptions,
  documentSourceTypeOptions,
  documentStatusOptions,
} from '@/constants/lookups'

const props = defineProps<{
  mode: 'create' | 'edit'
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
  fileName: string
  fileRequired: boolean
  knowledgeBaseOptions: { value: string; label: string }[]
}>()

const emit = defineEmits<{
  (event: 'update:knowledgeBaseId', value: string): void
  (event: 'update:chatSessionId', value: string): void
  (event: 'update:scopeType', value: number): void
  (event: 'update:sourceType', value: number): void
  (event: 'update:title', value: string): void
  (event: 'update:summary', value: string): void
  (event: 'update:tagsText', value: string): void
  (event: 'update:langCode', value: string): void
  (event: 'update:parseStrategy', value: number): void
  (event: 'update:status', value: number): void
  (event: 'update:file', value: File | null): void
}>()
</script>

<template>
  <div class="detail-grid">
    <div class="mini-grid">
      <label class="field">
        <span>Knowledge Base</span>
        <select
          :value="props.knowledgeBaseId"
          @change="emit('update:knowledgeBaseId', ($event.target as HTMLSelectElement).value)"
        >
          <option value="">Select one</option>
          <option v-for="item in props.knowledgeBaseOptions" :key="item.value" :value="item.value">
            {{ item.label }}
          </option>
        </select>
      </label>

      <label class="field">
        <span>Chat Session ID</span>
        <input
          :value="props.chatSessionId"
          type="text"
          placeholder="Optional"
          @input="emit('update:chatSessionId', ($event.target as HTMLInputElement).value)"
        />
      </label>
    </div>

    <div class="mini-grid">
      <label class="field">
        <span>Scope Type</span>
        <select
          :value="props.scopeType"
          @change="emit('update:scopeType', Number(($event.target as HTMLSelectElement).value))"
        >
          <option v-for="item in documentScopeTypeOptions" :key="item.value" :value="item.value">
            {{ item.label }}
          </option>
        </select>
      </label>

      <label class="field">
        <span>Source Type</span>
        <select
          :value="props.sourceType"
          @change="emit('update:sourceType', Number(($event.target as HTMLSelectElement).value))"
        >
          <option v-for="item in documentSourceTypeOptions" :key="item.value" :value="item.value">
            {{ item.label }}
          </option>
        </select>
      </label>
    </div>

    <label class="field">
      <span>Title</span>
      <input
        :value="props.title"
        type="text"
        placeholder="Document title"
        @input="emit('update:title', ($event.target as HTMLInputElement).value)"
      />
    </label>

    <label class="field">
      <span>Summary</span>
      <textarea
        :value="props.summary"
        placeholder="Document summary"
        @input="emit('update:summary', ($event.target as HTMLTextAreaElement).value)"
      />
    </label>

    <div class="mini-grid">
      <label class="field">
        <span>Tags</span>
        <input
          :value="props.tagsText"
          type="text"
          placeholder="tag1, tag2"
          @input="emit('update:tagsText', ($event.target as HTMLInputElement).value)"
        />
      </label>

      <label class="field">
        <span>Lang Code</span>
        <input
          :value="props.langCode"
          type="text"
          placeholder="zh-CN"
          @input="emit('update:langCode', ($event.target as HTMLInputElement).value)"
        />
      </label>
    </div>

    <div class="mini-grid">
      <label class="field">
        <span>Parse Strategy</span>
        <select
          :value="props.parseStrategy"
          @change="emit('update:parseStrategy', Number(($event.target as HTMLSelectElement).value))"
        >
          <option v-for="item in documentParseStrategyOptions" :key="item.value" :value="item.value">
            {{ item.label }}
          </option>
        </select>
      </label>

      <label class="field">
        <span>Status</span>
        <select
          :value="props.status"
          @change="emit('update:status', Number(($event.target as HTMLSelectElement).value))"
        >
          <option v-for="item in documentStatusOptions" :key="item.value" :value="item.value">
            {{ item.label }}
          </option>
        </select>
      </label>
    </div>

    <label class="field" v-if="props.mode === 'create'">
      <span>File</span>
      <input type="file" @change="emit('update:file', ($event.target as HTMLInputElement).files?.[0] || null)" />
    </label>

    <div class="notice soft-row">
      <span>File</span>
      <strong>{{ props.fileName || (props.fileRequired ? 'Required' : '-') }}</strong>
    </div>
  </div>
</template>
