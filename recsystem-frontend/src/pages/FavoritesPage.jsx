import { useEffect, useState } from 'react'
import { getContentsByIds } from '../api/catalogApi'
import { ContentGrid } from '../components/ContentGrid'
import { PageHeader } from '../components/PageHeader'
import { getFavoriteIds } from '../utils/favorites'

export function FavoritesPage() {
  const [items, setItems] = useState([])
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    const ids = getFavoriteIds()
    if (!ids.length) {
      setLoading(false)
      return
    }
    getContentsByIds(ids).then((data) => {
      setItems(data)
      setLoading(false)
    })
  }, [])

  return (
    <div className="content">
      <PageHeader title="Избранное" subtitle="Контент, который вы сохранили" />
      {loading && <p className="status-text">Загрузка…</p>}
      <ContentGrid items={items} emptyText="Добавьте контент в избранное с карточки" />
    </div>
  )
}
