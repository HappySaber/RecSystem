import { memo, useState } from 'react'
import { Link } from 'react-router-dom'
import { ContentPoster } from './ContentPoster'
import { useUserActions } from '../hooks/useUserActions'
import { isFavorite } from '../utils/favorites'

export const ContentItem = memo(function ContentItem({ content, eagerPoster = false }) {
  const [liked, setLiked] = useState(false)
  const [bookmarked, setBookmarked] = useState(() => isFavorite(content.id))
  const { trackAction, toggleFavorite } = useUserActions()

  const year = content.release_date?.slice(0, 4) ?? '—'
  const isAi = content.aiOnly

  const handleLike = async (e) => {
    e.preventDefault()
    e.stopPropagation()
    if (isAi) return
    const action = liked ? 'DISLIKE' : 'LIKE'
    await trackAction(content.id, action)
    setLiked(!liked)
  }

  const handleBookmark = async (e) => {
    e.preventDefault()
    e.stopPropagation()
    if (isAi) return
    await toggleFavorite(content.id, bookmarked)
    setBookmarked(!bookmarked)
  }

  const body = (
    <>
      <ContentPoster content={content} size="card" eager={eagerPoster} />
      <div className="content-info">
        <h3 className="content-title-text">
          {content.title}
          {isAi && <span className="ai-badge">ИИ</span>}
        </h3>
        <p className="content-meta">{year}</p>
        {content.description && (
          <p className="content-desc">
            {content.description.length > 100
              ? `${content.description.slice(0, 100)}…`
              : content.description}
          </p>
        )}
      </div>
      {!isAi && (
        <div className="badge-action">
          <button
            type="button"
            className={`btn-icon ${liked ? 'liked' : ''}`}
            onClick={handleLike}
            aria-label="лайк"
          >
            <i className="ti ti-heart" aria-hidden="true" />
          </button>
          <button
            type="button"
            className={`btn-icon ${bookmarked ? 'bookmarked' : ''}`}
            onClick={handleBookmark}
            aria-label="в избранное"
          >
            <i className="ti ti-bookmark" aria-hidden="true" />
          </button>
        </div>
      )}
    </>
  )

  if (isAi) {
    return <article className="content-item content-item--card">{body}</article>
  }

  return (
    <Link to={`/content/${content.id}`} className="content-item content-item--card">
      {body}
    </Link>
  )
})
