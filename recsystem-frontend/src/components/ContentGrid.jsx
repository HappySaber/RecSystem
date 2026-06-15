import { ContentItem } from './ContentItem'

const EAGER_POSTER_COUNT = 24

export function ContentGrid({ items, emptyText = 'Ничего не найдено' }) {
  if (!items?.length) {
    return <p className="empty-state">{emptyText}</p>
  }

  return (
    <div className="content-grid">
      {items.map((item, index) => (
        <ContentItem
          key={item.id}
          content={item}
          eagerPoster={index < EAGER_POSTER_COUNT}
        />
      ))}
    </div>
  )
}
