import { api } from '../services/api'

export async function getPersonalRecommendations(limit = 12, offset = 0) {
  const { data } = await api.get('/api/recommendations', { params: { limit, offset } })
  return data.content_ids ?? []
}

export async function getExplicitRecommendations(query, limit = 12) {
  const { data } = await api.post('/api/recommendations/explicit', { query, limit })
  const titles = data.titles ?? data.content_ids ?? []
  return Array.isArray(titles) ? titles.filter(Boolean) : []
}

export async function getRecommendationsByGenres(genres, limit = 12) {
  const { data } = await api.post('/api/recommendations/genres', { genres, limit })
  return data.content_ids ?? []
}

export async function getSimilarContent(contentId, limit = 12) {
  const { data } = await api.post('/api/recommendations/similar', {
    content_id: contentId,
    limit,
  })
  return data.content_ids ?? []
}

export async function getTrendingContent(limit = 12) {
  const { data } = await api.get('/api/recommendations/trending', { params: { limit } })
  return data.content_ids ?? []
}

export async function getPopularContent(limit = 12) {
  const { data } = await api.get('/api/recommendations/popular', { params: { limit } })
  return data.content_ids ?? []
}
