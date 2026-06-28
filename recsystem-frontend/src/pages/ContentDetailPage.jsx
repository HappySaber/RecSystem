import { useEffect, useState } from 'react'
import { Link, useParams } from 'react-router-dom'
import { getContentDetails } from '../api/catalogApi'
import { ContentPoster } from '../components/ContentPoster'
import { useUserActions } from '../hooks/useUserActions'
import { isFavorite } from '../utils/favorites'

export function ContentDetailPage() {
  const { id } = useParams()
  const [data, setData] = useState(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState(null)
  const [liked, setLiked] = useState(false)
  const [bookmarked, setBookmarked] = useState(false)
  const { trackAction, toggleFavorite } = useUserActions()

  useEffect(() => {
    setLoading(true)
    setError(null)
    getContentDetails(id)
      .then((res) => {
        setData(res)
        setBookmarked(isFavorite(id))
      })
      .catch((e) => setError(e.response?.data ?? e.message ?? 'Не удалось загрузить'))
      .finally(() => setLoading(false))
  }, [id])

  if (loading) {
    return (
      <div className="content">
        <p className="status-text">Загрузка…</p>
      </div>
    )
  }

  if (error || !data?.content) {
    return (
      <div className="content">
        <Link to="/" className="back-link">← Назад</Link>
        <p className="error-text">{error ?? 'Контент не найден'}</p>
      </div>
    )
  }

  const { content, movie } = data
  const year = content.release_date?.slice(0, 4)

  const handleLike = async () => {
    const action = liked ? 'DISLIKE' : 'LIKE'
    await trackAction(content.id, action)
    setLiked(!liked)
  }

  const handleBookmark = async () => {
    await toggleFavorite(content.id, bookmarked)
    setBookmarked(!bookmarked)
  }

  return (
    <div className="content">
      <Link to="/" className="back-link">← Назад к рекомендациям</Link>

      <div className="detail-layout">
        <div className="detail-poster">
          <ContentPoster content={content} size="detail" />
        </div>

        <div className="detail-main">
          <p className="detail-type">{content.type || content.external_source}</p>
          <h1 className="detail-title">{content.title}</h1>
          {movie?.tagline && <p className="detail-tagline">{movie.tagline}</p>}

          <div className="detail-meta-row">
            {year && <span>{year}</span>}
            {movie?.runtime ? <span>{movie.runtime} мин</span> : null}
            {movie?.language && <span>{movie.language.toUpperCase()}</span>}
          </div>

          {movie?.genres?.length > 0 && (
            <div className="detail-genres">
              {movie.genres.map((genre) => (
                <span key={genre} className="genre-chip active">{genre}</span>
              ))}
            </div>
          )}

          <div className="detail-actions">
            <button type="button" className={`btn-secondary ${liked ? 'liked' : ''}`} onClick={handleLike}>
              <i className="ti ti-heart" aria-hidden="true" /> {liked ? 'Нравится' : 'Лайк'}
            </button>
            <button
              type="button"
              className={`btn-secondary ${bookmarked ? 'bookmarked' : ''}`}
              onClick={handleBookmark}
            >
              <i className="ti ti-bookmark" aria-hidden="true" /> {bookmarked ? 'В избранном' : 'В избранное'}
            </button>
          </div>

          {content.description && (
            <section className="detail-section">
              <h2>Описание</h2>
              <p>{content.description}</p>
            </section>
          )}

          {movie?.cast?.length > 0 && (
            <section className="detail-section">
              <h2>В ролях</h2>
              <p className="detail-cast">{movie.cast.slice(0, 12).join(', ')}</p>
            </section>
          )}

          {movie && (
            <section className="detail-section detail-stats">
              {movie.budget > 0 && <div><span>Бюджет</span><strong>${movie.budget.toLocaleString()}</strong></div>}
              {movie.revenue > 0 && <div><span>Сборы</span><strong>${movie.revenue.toLocaleString()}</strong></div>}
              {movie.status && <div><span>Статус</span><strong>{movie.status}</strong></div>}
            </section>
          )}
        </div>
      </div>
    </div>
  )
}
