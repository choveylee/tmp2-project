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
