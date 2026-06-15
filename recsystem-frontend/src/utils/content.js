import { resolvePosterUrl } from './poster'

export function normalizeContent(content) {
  if (!content) return content
  return {
    ...content,
    poster_url: resolvePosterUrl(content.poster_url, content.external_source),
  }
}

export function normalizeContentList(items) {
  return (items ?? []).map(normalizeContent)
}

export function titlesToItems(titles) {
  return (titles ?? []).map((title, index) => ({
    id: `ai-${index}-${title}`,
    title,
    aiOnly: true,
  }))
}
