import { shallowRef } from 'vue'

const apiBaseUrl = shallowRef((import.meta.env.VITE_API_BASE_URL || 'http://127.0.0.1:8080').trim())

export function useAppConfig() {
  return {
    apiBaseUrl,
  }
}
