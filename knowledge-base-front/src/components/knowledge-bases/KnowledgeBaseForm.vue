<script setup lang="ts">
import { knowledgeBaseStatusOptions, knowledgeBaseVisibleOptions } from '@/constants/lookups'

const props = defineProps<{
  mode: 'create' | 'edit'
  code: string
  name: string
  description: string
  visible: number
  status: number
}>()

const emit = defineEmits<{
  (event: 'update:code', value: string): void
  (event: 'update:name', value: string): void
  (event: 'update:description', value: string): void
  (event: 'update:visible', value: number): void
  (event: 'update:status', value: number): void
}>()
</script>

<template>
  <div class="detail-grid">
    <label class="field">
      <span>Code</span>
      <input
        :value="props.code"
        type="text"
        :disabled="props.mode === 'edit'"
        placeholder="kb-code"
        @input="emit('update:code', ($event.target as HTMLInputElement).value)"
      />
    </label>

    <label class="field">
      <span>Name</span>
      <input
        :value="props.name"
        type="text"
        placeholder="Knowledge base name"
        @input="emit('update:name', ($event.target as HTMLInputElement).value)"
      />
    </label>

    <label class="field">
      <span>Description</span>
      <textarea
        :value="props.description"
        placeholder="Short description"
        @input="emit('update:description', ($event.target as HTMLTextAreaElement).value)"
      />
    </label>

    <div class="mini-grid">
      <label class="field">
        <span>Visible</span>
        <select
          :value="props.visible"
          @change="emit('update:visible', Number(($event.target as HTMLSelectElement).value))"
        >
          <option v-for="item in knowledgeBaseVisibleOptions" :key="item.value" :value="item.value">
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
          <option v-for="item in knowledgeBaseStatusOptions" :key="item.value" :value="item.value">
            {{ item.label }}
          </option>
        </select>
      </label>
    </div>
  </div>
</template>
