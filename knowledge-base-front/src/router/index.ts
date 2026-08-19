import { createRouter, createWebHistory } from 'vue-router'

import DocumentsPage from '@/views/DocumentsPage.vue'
import KnowledgeBasesPage from '@/views/KnowledgeBasesPage.vue'
import MonitorPage from '@/views/MonitorPage.vue'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    { path: '/', redirect: '/knowledge-bases' },
    { path: '/knowledge-bases', component: KnowledgeBasesPage },
    { path: '/documents', component: DocumentsPage },
    { path: '/monitor', component: MonitorPage },
  ],
})

export default router
