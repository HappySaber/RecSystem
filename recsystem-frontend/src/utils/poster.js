const TMDB_IMAGE_BASE = 'https://image.tmdb.org/t/p/w500'

const INVALID_VALUES = new Set(['', '<nil>', 'null', 'undefined', 'none'])

export function resolvePosterUrl(posterUrl, externalSource) {
  const raw = String(posterUrl ?? '').trim()

  if (INVALID_VALUES.has(raw.toLowerCase())) {
    return null
  }

  if (raw.startsWith('http://') || raw.startsWith('https://')) {
    return raw
  }

  if (raw.startsWith('//')) {
    return `https:${raw}`
  }

  const isTmdb =
    externalSource === 'tmdb' ||
    externalSource === '' ||
    !externalSource

  if (isTmdb) {
    const path = raw.startsWith('/') ? raw : `/${raw}`
    return `${TMDB_IMAGE_BASE}${path}`
  }

  return null
}
