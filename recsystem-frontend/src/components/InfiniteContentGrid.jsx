import { useEffect, useRef } from 'react'
import { ContentGrid } from './ContentGrid'

export function InfiniteContentGrid({
  items,
  loading,
  loadingMore,
  hasMore,
  onLoadMore,
  emptyText = 'Ничего не найдено',
}) {
  const sentinelRef = useRef(null)
  const onLoadMoreRef = useRef(onLoadMore)
  const loadingMoreRef = useRef(loadingMore)

  onLoadMoreRef.current = onLoadMore
  loadingMoreRef.current = loadingMore

  useEffect(() => {
    const node = sentinelRef.current
    if (!node || !hasMore) return undefined

    const observer = new IntersectionObserver(
      (entries) => {
        if (entries[0]?.isIntersecting && !loadingMoreRef.current) {
          onLoadMoreRef.current()
        }
      },
      { rootMargin: '300px' }
    )

    observer.observe(node)
    return () => observer.disconnect()
  }, [hasMore])

  return (
    <>
      <ContentGrid items={items} emptyText={loading ? 'Загрузка…' : emptyText} />
      {hasMore && (
        <div ref={sentinelRef} className="load-sentinel">
          {loadingMore && <p className="status-text">Загружаем ещё…</p>}
        </div>
      )}
    </>
  )
}
