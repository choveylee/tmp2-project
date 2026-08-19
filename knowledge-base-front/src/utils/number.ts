export function toOptionalNumber(value: string): number | undefined {
  const trimmed = value.trim()
  if (!trimmed) {
    return undefined
  }

  const parsed = Number(trimmed)
  return Number.isFinite(parsed) ? parsed : undefined
}

export function formatMaybeNumber(value: number | string | undefined, fallback = '-') {
  if (value === undefined || value === null || value === '') {
    return fallback
  }

  return String(value)
}

export function formatFileSize(bytes: number, fractionDigits = 1) {
  if (!Number.isFinite(bytes) || bytes < 0) {
    return '-'
  }

  if (bytes < 1024) {
    return `${bytes} B`
  }

  const units = ['KB', 'MB', 'GB', 'TB']
  let size = bytes / 1024
  let unitIndex = 0

  while (size >= 1024 && unitIndex < units.length - 1) {
    size /= 1024
    unitIndex += 1
  }

  return `${size.toFixed(fractionDigits)} ${units[unitIndex]}`
}
