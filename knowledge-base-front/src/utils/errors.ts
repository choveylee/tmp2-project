import { ApiError } from '@/api/client'

export function formatErrorMessage(error: unknown, fallback = 'Request failed') {
  if (error instanceof ApiError) {
    return error.detail || error.message || fallback
  }

  if (error instanceof Error) {
    return error.message || fallback
  }

  if (typeof error === 'string' && error.trim()) {
    return error.trim()
  }

  return fallback
}
