import { memo } from 'react'

export const ContentPoster = memo(function ContentPoster({
  content,
  size = 'card',
  className = '',
}) {
  const src = content?.poster_url || null
  const letter = content?.title?.[0]?.toUpperCase() ?? '?'
  const wrapCls = `content-poster-wrap content-poster-wrap--${size} ${className}`.trim()

  return (
    <div className={wrapCls}>
      <span className="content-poster-fallback" aria-hidden="true">
        {letter}
      </span>
      {src && (
        <img
          className="content-poster"
          src={src}
          alt={content?.title ?? 'Постер'}
          referrerPolicy="no-referrer"
          loading="lazy"
          decoding="async"
          draggable={false}
        />
      )}
    </div>
  )
})
