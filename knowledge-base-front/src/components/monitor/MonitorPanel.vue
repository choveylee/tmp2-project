<script setup lang="ts">
import { computed } from 'vue'

import { useMonitor } from '@/composables/useMonitor'
import { formatDateTime } from '@/utils/time'

const { cpu, ram, loading, error, lastCheckedAt, refresh } = useMonitor()

const cpuLabel = computed(() => cpu.value?.status || 'idle')
const ramLabel = computed(() => ram.value?.status || 'idle')
</script>

<template>
  <section class="page-grid kb-page">
    <section class="page-section stack kb-hero">
      <div class="section-head">
        <div class="stack-tight">
          <p class="eyebrow">System watch</p>
          <h2 class="section-title">Monitor</h2>
          <p class="muted">Quick CPU and RAM checks from the backend.</p>
        </div>
        <div class="toolbar">
          <button type="button" class="primary" :disabled="loading" @click="refresh">
            {{ loading ? 'Checking...' : 'Refresh' }}
          </button>
        </div>
      </div>

      <div v-if="error" class="notice error">{{ error }}</div>

      <div class="summary-row">
        <span class="badge soft">CPU: {{ cpuLabel }}</span>
        <span class="badge soft">RAM: {{ ramLabel }}</span>
        <span class="badge soft">Checked: {{ formatDateTime(lastCheckedAt) }}</span>
      </div>

      <div class="metric-grid">
        <section class="panel metric-card stack">
          <div class="section-head">
            <div class="stack-tight">
              <p class="eyebrow">CPU</p>
              <h3 class="section-title">CPU Check</h3>
              <p class="muted">`/cpu-check`</p>
            </div>
          </div>
          <div v-if="cpu" class="metric-value">
            <span :class="['badge', cpu.status === 'ok' ? 'good' : 'bad']">{{ cpu.status }}</span>
            <div class="detail-text">{{ cpu.detail }}</div>
          </div>
          <div v-else class="empty-state">No CPU response yet.</div>
        </section>

        <section class="panel metric-card stack">
          <div class="section-head">
            <div class="stack-tight">
              <p class="eyebrow">RAM</p>
              <h3 class="section-title">RAM Check</h3>
              <p class="muted">`/ram-check`</p>
            </div>
          </div>
          <div v-if="ram" class="metric-value">
            <span :class="['badge', ram.status === 'ok' ? 'good' : 'bad']">{{ ram.status }}</span>
            <div class="detail-text">{{ ram.detail }}</div>
          </div>
          <div v-else class="empty-state">No RAM response yet.</div>
        </section>
      </div>
    </section>
  </section>
</template>
