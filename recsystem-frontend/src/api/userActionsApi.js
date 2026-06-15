import { api } from '../services/api'

const ACTION_PATHS = {
  LIKE: 'like',
  DISLIKE: 'dislike',
  ADD_TO_FAVORITES: 'favorite',
  VIEW: 'view',
}

export async function trackUserAction(contentId, action) {
  const path = ACTION_PATHS[action] ?? action.toLowerCase()
  await api.post(`/api/content/${contentId}/${path}`)
}
