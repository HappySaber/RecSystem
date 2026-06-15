import { useState } from 'react'
import { AISearch } from '../components/AISearch'
import { ContentGrid } from '../components/ContentGrid'
import { PageHeader } from '../components/PageHeader'

export function AIPage() {
  const [items, setItems] = useState([])
  const [searched, setSearched] = useState(false)

  const handleResults = (results) => {
    setSearched(true)
    setItems(results)
  }

  return (
    <div className="content">
      <PageHeader
        title="ИИ-подборка"
        subtitle="Опишите настроение или тему — система подберёт контент"
      />
      <AISearch onResults={handleResults} />
      <ContentGrid
        items={items}
        emptyText={
          searched
            ? 'ИИ не вернул названий — проверьте логи recommendation-сервиса'
            : 'Введите запрос и нажмите «Найти»'
        }
      />
    </div>
  )
}
