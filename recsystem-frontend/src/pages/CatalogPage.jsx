import { useState } from 'react'
import { getContent } from '../api/catalogApi'
import { ContentGrid } from '../components/ContentGrid'
import { PageHeader } from '../components/PageHeader'
import { useRecommendations } from '../hooks/useRecommendations'

export function CatalogPage() {
  const [searchId, setSearchId] = useState('')
  const [items, setItems] = useState([])
  const [error, setError] = useState('')
  const [loading, setLoading] = useState(false)
  const { getPopular } = useRecommendations()

  const loadPopular = async () => {
    setLoading(true)
    setError('')
    const result = await getPopular(20)
    setItems(result)
    setLoading(false)
  }

  const searchById = async (e) => {
    e.preventDefault()
    if (!searchId.trim()) return
    setLoading(true)
    setError('')
    try {
      const content = await getContent(searchId.trim())
      setItems([content])
    } catch {
      setError('Контент с таким ID не найден')
      setItems([])
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="content">
      <PageHeader
        title="Каталог"
        subtitle="Поиск по ID или просмотр популярного контента"
      />

      <div className="catalog-toolbar">
        <form className="search-row" onSubmit={searchById}>
          <input
            className="search-input"
            placeholder="UUID контента…"
            value={searchId}
            onChange={(e) => setSearchId(e.target.value)}
          />
          <button type="submit" className="btn-secondary" disabled={loading}>
            Найти
          </button>
        </form>
        <button type="button" className="btn-secondary" onClick={loadPopular} disabled={loading}>
          Показать популярное
        </button>
      </div>

      {loading && <p className="status-text">Загрузка…</p>}
      {error && <p className="error-text">{error}</p>}
      <ContentGrid items={items} emptyText="Найдите контент по ID или загрузите популярное" />
    </div>
  )
}
