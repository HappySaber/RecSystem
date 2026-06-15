import { useEffect, useState } from 'react'
import { ContentGrid } from '../components/ContentGrid'
import { GenreFilter } from '../components/GenreFilter'
import { PageHeader } from '../components/PageHeader'
import { useRecommendations } from '../hooks/useRecommendations'

export function PreferencesPage() {
  const [items, setItems] = useState([])
  const { getByGenres, loading, error } = useRecommendations()

  const handleGenresChange = async (genres) => {
    if (!genres.length) {
      setItems([])
      return
    }
    const result = await getByGenres(genres, 16)
    setItems(result)
  }

  useEffect(() => {
    handleGenresChange(['Action', 'Thriller'])
  }, [])

  return (
    <div className="content">
      <PageHeader
        title="Предпочтения"
        subtitle="Выберите жанры — получите подборку по вкусу"
      />
      <GenreFilter onChange={handleGenresChange} />
      {loading && <p className="status-text">Загрузка…</p>}
      {error && <p className="error-text">{error}</p>}
      <ContentGrid items={items} emptyText="Выберите один или несколько жанров" />
    </div>
  )
}
