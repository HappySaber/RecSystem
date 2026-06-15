import { useState, useCallback } from 'react'
import {
  getPersonalRecommendations,
  getExplicitRecommendations,
  getRecommendationsByGenres,
  getTrendingContent,
  getPopularContent,
} from '../api/recommendationApi'
import { getContentsByIds } from '../api/catalogApi'
import { titlesToItems } from '../utils/content'

async function resolveContentIds(ids) {
  if (!ids?.length) return []
  return getContentsByIds(ids)
}

function formatError(e) {
  return typeof e.response?.data === 'string'
    ? e.response.data
    : e.response?.data?.message ?? e.message ?? 'Ошибка загрузки'
}

export function useRecommendations() {
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState(null)

  const wrapWithCatalog = useCallback(async (fn) => {
    setLoading(true)
    setError(null)
    try {
      const ids = await fn()
      return resolveContentIds(ids)
    } catch (e) {
      setError(formatError(e))
      return []
    } finally {
      setLoading(false)
    }
  }, [])

  const getPersonal = useCallback(
    (limit = 12) => wrapWithCatalog(() => getPersonalRecommendations(limit)),
    [wrapWithCatalog]
  )

  const getExplicit = useCallback(async (query, limit = 12) => {
    setLoading(true)
    setError(null)
    try {
      const titles = await getExplicitRecommendations(query, limit)
      return titlesToItems(titles)
    } catch (e) {
      setError(formatError(e))
      return []
    } finally {
      setLoading(false)
    }
  }, [])

  const getByGenres = useCallback(
    (genres, limit = 12) => wrapWithCatalog(() => getRecommendationsByGenres(genres, limit)),
    [wrapWithCatalog]
  )

  const getTrending = useCallback(
    (limit = 12) => wrapWithCatalog(() => getTrendingContent(limit)),
    [wrapWithCatalog]
  )

  const getPopular = useCallback(
    (limit = 12) => wrapWithCatalog(() => getPopularContent(limit)),
    [wrapWithCatalog]
  )

  return { getPersonal, getExplicit, getByGenres, getTrending, getPopular, loading, error }
}
