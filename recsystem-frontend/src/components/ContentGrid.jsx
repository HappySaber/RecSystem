import { ContentItem } from './ContentItem'

export function ContentGrid({ items, emptyText = 'Ничего не найдено' }) {
  if (!items?.length) {
    return <p className="empty-state">{emptyText}</p>
  }

  return (
    <div className="content-grid">
      {items.map((item) => (
        <ContentItem key={item.id} content={item} />
      ))}
    </div>
  )
}
