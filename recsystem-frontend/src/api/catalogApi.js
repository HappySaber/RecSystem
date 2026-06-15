import { api } from '../services/api'
import { normalizeContent, normalizeContentList } from '../utils/content'

export async function getContent(contentId) {
  const { data } = await api.get('/api/catalog/get_content', {
    params: { content_id: contentId },
  })
  return normalizeContent(data)
}

export async function getContentDetails(contentId) {
  const { data } = await api.get(`/api/catalog/content/${contentId}`)
  return {
    ...data,
    content: normalizeContent(data.content),
  }
}

export async function getContentsByIds(ids) {
  if (!ids?.length) return []

  const results = await Promise.allSettled(ids.map((id) => getContent(id)))

  const ordered = ids
    .map((id, index) => {
      const result = results[index]
      return result.status === 'fulfilled' ? result.value : null
    })
    .filter(Boolean)

  return normalizeContentList(ordered)
}
