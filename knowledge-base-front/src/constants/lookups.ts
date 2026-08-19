export interface OptionItem {
  value: number
  label: string
}

export const knowledgeBaseVisibleOptions: OptionItem[] = [
  { value: 0, label: 'Private' },
  { value: 1, label: 'Internal' },
  { value: 2, label: 'Public' },
]

export const knowledgeBaseStatusOptions: OptionItem[] = [
  { value: 0, label: 'Disabled' },
  { value: 1, label: 'Normal' },
]

export const documentScopeTypeOptions: OptionItem[] = [
  { value: 0, label: 'Knowledge' },
  { value: 1, label: 'Attachment' },
]

export const documentSourceTypeOptions: OptionItem[] = [
  { value: 0, label: 'User' },
  { value: 1, label: 'API' },
]

export const documentParseStrategyOptions: OptionItem[] = [
  { value: 0, label: 'Auto' },
  { value: 1, label: 'Tika' },
  { value: 2, label: 'OCR' },
  { value: 3, label: 'Tika OCR' },
]

export const documentProcessStatusOptions: OptionItem[] = [
  { value: 0, label: 'Waiting' },
  { value: 1, label: 'Uploaded' },
  { value: 2, label: 'Processing' },
  { value: 3, label: 'Processed' },
  { value: 4, label: 'Failed' },
]

export const documentStatusOptions: OptionItem[] = [
  { value: 0, label: 'Disabled' },
  { value: 1, label: 'Normal' },
]

export const ingestJobStatusOptions: OptionItem[] = [
  { value: 0, label: 'Pending' },
  { value: 1, label: 'Processing' },
  { value: 2, label: 'Finished' },
  { value: 3, label: 'Failed' },
  { value: 4, label: 'Canceled' },
]

export const ingestJobTypeOptions: OptionItem[] = [
  { value: 0, label: 'Parse' },
  { value: 1, label: 'OCR' },
  { value: 2, label: 'Split' },
  { value: 3, label: 'Embedding' },
  { value: 4, label: 'Index' },
  { value: 5, label: 'Rebuild index' },
]

export function labelFromOptions(options: OptionItem[], value: number | string | undefined) {
  const normalized = typeof value === 'string' ? Number(value) : value
  const found = options.find((item) => item.value === normalized)
  return found ? found.label : normalized === undefined ? '-' : String(normalized)
}
