import { useEffect } from 'react'
import { InfiniteContentGrid } from '../components/InfiniteContentGrid'
import { PageHeader } from '../components/PageHeader'
import { useInfinitePersonalRecommendations } from '../hooks/useInfiniteRecommendations'

export function HomePage() {
  const { items, loading, loadingMore, error, hasMore, loadMore, reload } =
    useInfinitePersonalRecommendations()

  useEffect(() => {
    reload()
  }, [reload])

  return (
    <div className="content">
      <PageHeader
        title="Рекомендации для вас"
        subtitle="Прокрутите вниз — подгрузим ещё. Для новых пользователей показываем популярное"
      />
      {error && <p className="error-text">{error}</p>}
      <InfiniteContentGrid
        items={items}
        loading={loading}
        loadingMore={loadingMore}
        hasMore={hasMore}
        onLoadMore={loadMore}
        emptyText="Пока нет рекомендаций — оцените несколько фильмов"
      />
    </div>
  )
}
