export function splitTagsText(value: string) {
  const seen = new Set<string>()
  const tags: string[] = []

  for (const part of value.split(',')) {
    const tag = part.trim()
    if (!tag || seen.has(tag)) {
      continue
    }

    seen.add(tag)
    tags.push(tag)
  }

  return tags
}

export function joinTagsText(tags: string[] | undefined | null) {
  return (tags || [])
    .map((tag) => tag.trim())
    .filter(Boolean)
    .join(', ')
}
