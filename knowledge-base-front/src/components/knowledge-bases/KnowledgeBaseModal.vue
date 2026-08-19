<script setup lang="ts">
import { computed, onBeforeUnmount, watch } from 'vue'

import KnowledgeBaseForm from './KnowledgeBaseForm.vue'

const props = defineProps<{
  open: boolean
  mode: 'create' | 'edit'
  saving: boolean
  error: string
  code: string
  name: string
  description: string
  visible: number
  status: number
}>()

const emit = defineEmits<{
  (event: 'close'): void
  (event: 'submit'): void
  (event: 'reset'): void
  (event: 'update:code', value: string): void
  (event: 'update:name', value: string): void
  (event: 'update:description', value: string): void
  (event: 'update:visible', value: number): void
  (event: 'update:status', value: number): void
}>()

const title = computed(() => (props.mode === 'create' ? 'Create knowledge base' : 'Edit knowledge base'))
const subtitle = computed(() =>
  props.mode === 'create'
    ? 'Set the code first, then refine visibility and lifecycle settings.'
    : 'Adjust the mutable fields and save the updated record.',
)
const submitLabel = computed(() => (props.mode === 'create' ? 'Create' : 'Save changes'))
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
        <section class="modal-panel" role="dialog" aria-modal="true" aria-labelledby="kb-modal-title" aria-describedby="kb-modal-desc">
          <header class="modal-head">
            <div class="stack-tight">
              <p class="eyebrow">Knowledge base</p>
              <h3 id="kb-modal-title" class="modal-title">{{ title }}</h3>
              <p id="kb-modal-desc" class="muted">{{ subtitle }}</p>
            </div>

            <button type="button" class="modal-close" aria-label="Close dialog" @click="emit('close')">×</button>
          </header>

          <div class="modal-body stack">
            <div v-if="props.error" class="notice error">{{ props.error }}</div>
            <KnowledgeBaseForm
              :mode="props.mode"
              :code="props.code"
              :name="props.name"
              :description="props.description"
              :visible="props.visible"
              :status="props.status"
              @update:code="emit('update:code', $event)"
              @update:name="emit('update:name', $event)"
              @update:description="emit('update:description', $event)"
              @update:visible="emit('update:visible', $event)"
              @update:status="emit('update:status', $event)"
            />
          </div>

          <footer class="modal-footer">
            <button type="button" :disabled="props.saving" @click="emit('reset')">{{ resetLabel }}</button>
            <div class="toolbar">
              <button type="button" :disabled="props.saving" @click="emit('close')">Cancel</button>
              <button type="button" class="primary" :disabled="props.saving" @click="emit('submit')">
                {{ props.saving ? 'Saving...' : submitLabel }}
              </button>
            </div>
          </footer>
        </section>
      </div>
    </Transition>
  </Teleport>
</template>
