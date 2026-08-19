<script setup lang="ts">
import { computed, onBeforeUnmount, watch } from 'vue'

import DocumentForm from './DocumentForm.vue'

const props = defineProps<{
  open: boolean
  mode: 'create' | 'edit'
  saving: boolean
  error: string
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
  fileLabel: string
  fileRequired: boolean
  uploading: boolean
  knowledgeBaseOptions: { value: string; label: string }[]
}>()

const emit = defineEmits<{
  (event: 'close'): void
  (event: 'submit'): void
  (event: 'reset'): void
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

const title = computed(() => (props.mode === 'create' ? 'Create document' : 'Edit document'))
const subtitle = computed(() =>
  props.mode === 'create'
    ? 'Upload the source file first, then create the document record with its file object.'
    : 'Adjust metadata only. The uploaded file is preserved.',
)
const submitLabel = computed(() => (props.mode === 'create' ? 'Create document' : 'Save changes'))
const resetLabel = computed(() => (props.mode === 'create' ? 'Reset' : 'Reset changes'))

function onKeydown(event: KeyboardEvent) {
  if (event.key === 'Escape' && props.open) {
    emit('close')
  }
}

watch(
  () => props.open,
  (open) => {
    if (open) {
      window.addEventListener('keydown', onKeydown)
      return
    }

    window.removeEventListener('keydown', onKeydown)
  },
  { immediate: true },
)

onBeforeUnmount(() => {
  window.removeEventListener('keydown', onKeydown)
})
</script>

<template>
  <Teleport to="body">
    <Transition name="modal-shell">
      <div v-if="props.open" class="modal-overlay" @click.self="emit('close')">
        <section
          class="modal-panel"
          role="dialog"
          aria-modal="true"
          aria-labelledby="document-modal-title"
          aria-describedby="document-modal-desc"
        >
          <header class="modal-head">
            <div class="stack-tight">
              <p class="eyebrow">Document</p>
              <h3 id="document-modal-title" class="modal-title">{{ title }}</h3>
              <p id="document-modal-desc" class="muted">{{ subtitle }}</p>
            </div>

            <button type="button" class="modal-close" aria-label="Close dialog" @click="emit('close')">×</button>
          </header>

          <div class="modal-body stack">
            <div v-if="props.error" class="notice error">{{ props.error }}</div>
            <DocumentForm
              :mode="props.mode"
              :knowledge-base-id="props.knowledgeBaseId"
              :chat-session-id="props.chatSessionId"
              :scope-type="props.scopeType"
              :source-type="props.sourceType"
              :title="props.title"
              :summary="props.summary"
              :tags-text="props.tagsText"
              :lang-code="props.langCode"
              :parse-strategy="props.parseStrategy"
              :status="props.status"
              :file-name="props.fileName"
              :file-label="props.fileLabel"
              :file-required="props.fileRequired"
              :knowledge-base-options="props.knowledgeBaseOptions"
              @update:knowledgeBaseId="emit('update:knowledgeBaseId', $event)"
              @update:chatSessionId="emit('update:chatSessionId', $event)"
              @update:scopeType="emit('update:scopeType', $event)"
              @update:sourceType="emit('update:sourceType', $event)"
              @update:title="emit('update:title', $event)"
              @update:summary="emit('update:summary', $event)"
              @update:tagsText="emit('update:tagsText', $event)"
              @update:langCode="emit('update:langCode', $event)"
              @update:parseStrategy="emit('update:parseStrategy', $event)"
              @update:status="emit('update:status', $event)"
              @update:file="emit('update:file', $event)"
            />
          </div>

          <footer class="modal-footer">
            <button type="button" :disabled="props.saving" @click="emit('reset')">{{ resetLabel }}</button>
            <div class="toolbar">
              <button type="button" :disabled="props.saving" @click="emit('close')">Cancel</button>
              <button type="button" class="primary" :disabled="props.saving" @click="emit('submit')">
                {{ props.saving ? (props.uploading ? 'Uploading...' : 'Saving...') : submitLabel }}
              </button>
            </div>
          </footer>
        </section>
      </div>
    </Transition>
  </Teleport>
</template>
