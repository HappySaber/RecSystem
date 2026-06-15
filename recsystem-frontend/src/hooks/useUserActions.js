import { useCallback } from 'react'
import { trackUserAction } from '../api/userActionsApi'
import { addFavorite, removeFavorite } from '../utils/favorites'

export function useUserActions() {
  const trackAction = useCallback(async (contentId, action) => {
    try {
      await trackUserAction(contentId, action)
      if (action === 'ADD_TO_FAVORITES') {
        addFavorite(contentId)
      }
    } catch (e) {
      console.error('Failed to track action:', e)
    }
  }, [])

  const toggleFavorite = useCallback(async (contentId, isFav) => {
    if (isFav) {
      removeFavorite(contentId)
    } else {
      await trackUserAction(contentId, 'ADD_TO_FAVORITES')
      addFavorite(contentId)
    }
  }, [])

  return { trackAction, toggleFavorite }
}
