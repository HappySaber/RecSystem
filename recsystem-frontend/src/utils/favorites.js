const KEY = 'recsystem_favorites'

export function getFavoriteIds() {
  try {
    return JSON.parse(localStorage.getItem(KEY) ?? '[]')
  } catch {
    return []
  }
}

export function addFavorite(id) {
  const ids = getFavoriteIds()
  if (!ids.includes(id)) {
    localStorage.setItem(KEY, JSON.stringify([...ids, id]))
  }
}

export function removeFavorite(id) {
  localStorage.setItem(
    KEY,
    JSON.stringify(getFavoriteIds().filter((x) => x !== id))
  )
}

export function isFavorite(id) {
  return getFavoriteIds().includes(id)
}
