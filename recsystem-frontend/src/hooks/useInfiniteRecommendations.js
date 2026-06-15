import { useCallback, useRef, useState } from 'react'
import { getPersonalRecommendations } from '../api/recommendationApi'
import { getContentsByIds } from '../api/catalogApi'

const PAGE_SIZE = 12

function formatError(e) {
  return typeof e.response?.data === 'string'
    ? e.response.data
    : e.response?.data?.message ?? e.message ?? 'Ошибка загрузки'
}

export function useInfinitePersonalRecommendations() {
  const [items, setItems] = useState([])
  const [loading, setLoading] = useState(false)
  const [loadingMore, setLoadingMore] = useState(false)
  const [error, setError] = useState(null)
  const [hasMore, setHasMore] = useState(true)
  const offsetRef = useRef(0)
  const loadingRef = useRef(false)

  const appendUnique = useCallback((prev, next) => {
    const existing = new Map(prev.map((item) => [item.id, item]))
    for (const item of next) {
      if (!existing.has(item.id)) {
        existing.set(item.id, item)
      }
    }
    return Array.from(existing.values())
  }, [])

  const loadPage = useCallback(async (reset = false) => {
    if (loadingRef.current) return
    loadingRef.current = true

    if (reset) {
      setLoading(true)
      offsetRef.current = 0
      setHasMore(true)
      setError(null)
    } else {
      setLoadingMore(true)
    }

    try {
      const currentOffset = reset ? 0 : offsetRef.current
      const ids = await getPersonalRecommendations(PAGE_SIZE, currentOffset)

      if (ids.length < PAGE_SIZE) {
        setHasMore(false)
      }

      if (ids.length === 0) {
        if (reset) setItems([])
        return
      }

      const contents = await getContentsByIds(ids)
      offsetRef.current = currentOffset + ids.length
      setItems((prev) => (reset ? contents : appendUnique(prev, contents)))
    } catch (e) {
      setError(formatError(e))
      if (reset) setItems([])
    } finally {
      loadingRef.current = false
      setLoading(false)
      setLoadingMore(false)
    }
  }, [appendUnique])

  const loadMore = useCallback(() => {
    if (!hasMore || loadingRef.current) return
    loadPage(false)
  }, [hasMore, loadPage])

  const reload = useCallback(() => loadPage(true), [loadPage])

  return { items, loading, loadingMore, error, hasMore, loadMore, reload }
}
