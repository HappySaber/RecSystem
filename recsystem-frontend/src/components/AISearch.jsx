import { useState } from 'react'
import { useRecommendations } from '../hooks/useRecommendations'

export function AISearch({ onResults }) {
  const [query, setQuery] = useState('')
  const { getExplicit, loading, error } = useRecommendations()

  const handleSearch = async () => {
    if (!query.trim()) return
    const results = await getExplicit(query.trim(), 12)
    onResults(results)
  }

  const handleKeyDown = (e) => {
    if (e.key === 'Enter') handleSearch()
  }

  return (
    <div className="full-card">
      <div className="ai-box-title">
        <i className="ti ti-sparkles" aria-hidden="true" />
        ИИ-рекомендации по запросу
      </div>
      <div className="ai-input-row">
        <input
          className="ai-input"
          type="text"
          value={query}
          onChange={(e) => setQuery(e.target.value)}
          onKeyDown={handleKeyDown}
          placeholder="Например: атмосферный фильм про космос…"
        />
        <button type="button" className="ai-btn" onClick={handleSearch} disabled={loading}>
          {loading ? 'Ищу…' : 'Найти'}
        </button>
      </div>
      {error && (
        <p className="error-text" role="alert">
          Ошибка: {error}
        </p>
      )}
    </div>
  )
}
