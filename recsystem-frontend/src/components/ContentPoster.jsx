import { memo, useMemo, useState } from 'react'
import { resolvePosterUrl } from '../utils/poster'

const loadedPosters = new Set()

export const ContentPoster = memo(function ContentPoster({
  content,
  size = 'card',
  className = '',
  eager = false,
}) {
  const src = useMemo(() => {
    const resolved = resolvePosterUrl(content?.poster_url, content?.external_source)
    if (resolved) return resolved
    if (content?.poster_url?.startsWith('http')) return content.poster_url
    return null
  }, [content?.id, content?.poster_url, content?.external_source])

  const [failed, setFailed] = useState(false)
  const [loaded, setLoaded] = useState(() => (src ? loadedPosters.has(src) : false))

  const cls = `content-poster content-poster--${size} ${className}`.trim()

  if (!src || failed) {
    return (
      <div className={`${cls} content-poster--placeholder`}>
        {content?.title?.[0]?.toUpperCase() ?? '?'}
      </div>
    )
  }

  return (
    <img
      className={`${cls} ${loaded ? 'content-poster--loaded' : 'content-poster--loading'}`}
      src={src}
      alt={content?.title ?? 'Постер'}
      loading={eager ? 'eager' : 'lazy'}
      decoding="async"
      referrerPolicy="no-referrer"
      onLoad={() => {
        loadedPosters.add(src)
        setLoaded(true)
      }}
      onError={() => setFailed(true)}
    />
  )
})
