import { useEffect, useState } from 'react'
import { ContentGrid } from '../components/ContentGrid'
import { PageHeader } from '../components/PageHeader'
import { useRecommendations } from '../hooks/useRecommendations'

export function TrendingPage() {
  const [items, setItems] = useState([])
  const { getTrending, loading, error } = useRecommendations()

  useEffect(() => {
    getTrending(16).then(setItems)
  }, [getTrending])

  return (
    <div className="content">
      <PageHeader title="Трендовое" subtitle="Самый обсуждаемый контент прямо сейчас" />
      {loading && <p className="status-text">Загрузка…</p>}
      {error && <p className="error-text">{error}</p>}
      <ContentGrid items={items} />
    </div>
  )
}
