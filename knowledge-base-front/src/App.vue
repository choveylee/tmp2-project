<script setup lang="ts">
import { RouterLink, RouterView } from 'vue-router'
</script>

<template>
  <div class="app-shell">
    <header class="app-header">
      <div class="app-brand">
        <div class="app-title">Knowledge Base Front</div>
        <div class="app-subtitle">Knowledge bases, documents, and health checks</div>
      </div>

    <nav class="app-nav">
      <RouterLink to="/knowledge-bases">Knowledge Bases</RouterLink>
      <RouterLink to="/documents">Documents</RouterLink>
      <RouterLink to="/monitor">Monitor</RouterLink>
    </nav>
  </header>

  <main class="app-main">
    <RouterView />
  </main>
  </div>
</template>
