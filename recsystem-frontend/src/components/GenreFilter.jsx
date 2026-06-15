import { useState } from 'react'

const GENRES = [
  { label: 'Боевик', value: 'Action' },
  { label: 'Триллер', value: 'Thriller' },
  { label: 'Драма', value: 'Drama' },
  { label: 'Фантастика', value: 'Science Fiction' },
  { label: 'Комедия', value: 'Comedy' },
  { label: 'Ужасы', value: 'Horror' },
  { label: 'Романтика', value: 'Romance' },
  { label: 'Приключения', value: 'Adventure' },
]

export function GenreFilter({ onChange }) {
  const [selected, setSelected] = useState(['Action', 'Thriller'])

  const toggle = (value) => {
    const next = selected.includes(value)
      ? selected.filter((g) => g !== value)
      : [...selected, value]
    setSelected(next)
    onChange(next)
  }

  return (
    <div className="genres-row">
      {GENRES.map(({ label, value }) => (
        <button
          key={value}
          type="button"
          className={`genre-chip ${selected.includes(value) ? 'active' : ''}`}
          onClick={() => toggle(value)}
        >
          {label}
        </button>
      ))}
    </div>
  )
}
