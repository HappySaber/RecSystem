import { useEffect, useState } from 'react'
import { ContentGrid } from '../components/ContentGrid'
import { PageHeader } from '../components/PageHeader'
import { useRecommendations } from '../hooks/useRecommendations'

export function PopularPage() {
  const [items, setItems] = useState([])
  const { getPopular, loading, error } = useRecommendations()

  useEffect(() => {
    getPopular(16).then(setItems)
  }, [getPopular])

  return (
    <div className="content">
      <PageHeader title="Популярное" subtitle="Контент с наибольшим числом просмотров" />
      {loading && <p className="status-text">Загрузка…</p>}
      {error && <p className="error-text">{error}</p>}
      <ContentGrid items={items} />
    </div>
  )
}
