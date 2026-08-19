import { useAppConfig } from '@/composables/useAppConfig'
import type { ResponseEnvelope } from '@/types/api'

export class ApiError extends Error {
  status: number
  code: number
  detail: string
  raw: unknown

  constructor(status: number, code: number, detail: string, raw: unknown) {
    super(detail)
    this.name = 'ApiError'
    this.status = status
    this.code = code
    this.detail = detail
    this.raw = raw
  }
}

function joinUrl(baseUrl: string, path: string) {
  const normalizedBase = baseUrl.trim().replace(/\/+$/, '') || 'http://127.0.0.1:8080'
  const normalizedPath = path.startsWith('/') ? path : `/${path}`

  return `${normalizedBase}${normalizedPath}`
}

async function parseResponse<T>(response: Response): Promise<T> {
  const text = await response.text()

  let payload: ResponseEnvelope<T> | null = null
  if (text) {
    try {
      payload = JSON.parse(text) as ResponseEnvelope<T>
    } catch {
      throw new ApiError(response.status, response.status, text || response.statusText, text)
    }
  }

  if (!payload) {
    throw new ApiError(response.status, response.status, response.statusText, text)
  }

  if (!response.ok || payload.code !== 0) {
    throw new ApiError(
      response.status,
      payload.code ?? response.status,
      payload.detail || payload.message || response.statusText,
      payload,
    )
  }

  return payload.data as T
}

export async function requestJson<T>(
  path: string,
  options: {
    method?: string
    body?: unknown
    headers?: Record<string, string>
  } = {},
): Promise<T> {
  const { apiBaseUrl } = useAppConfig()
  const headers: Record<string, string> = { ...(options.headers || {}) }

  let body: BodyInit | undefined
  if (options.body !== undefined) {
    headers['Content-Type'] = 'application/json'
    body = JSON.stringify(options.body)
  }

  const response = await fetch(joinUrl(apiBaseUrl.value, path), {
    method: options.method || 'GET',
    headers,
    body,
  })

  return parseResponse<T>(response)
}

export async function requestFormData<T>(
  path: string,
  formData: FormData,
  method = 'POST',
): Promise<T> {
  const { apiBaseUrl } = useAppConfig()
  const headers: Record<string, string> = {}

  const response = await fetch(joinUrl(apiBaseUrl.value, path), {
    method,
    headers,
    body: formData,
  })

  return parseResponse<T>(response)
}
