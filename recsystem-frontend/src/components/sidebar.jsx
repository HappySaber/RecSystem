import { NavLink, useNavigate } from 'react-router-dom'
import { useAuth } from '../context/AuthContext'

const navItems = [
  {
    section: 'Главное',
    items: [
      { to: '/', icon: 'ti-layout-dashboard', label: 'Главная' },
      { to: '/catalog', icon: 'ti-stack-2', label: 'Каталог' },
      { to: '/favorites', icon: 'ti-star', label: 'Избранное' },
    ],
  },
  {
    section: 'ИИ-Рекомендации',
    items: [
      { to: '/ai', icon: 'ti-sparkles', label: 'ИИ-подборка' },
      // { to: '/trending', icon: 'ti-flame', label: 'Трендовое' },
      // { to: '/popular', icon: 'ti-chart-bar', label: 'Популярное' },
    ],
  },
  {
    section: 'Аккаунт',
    items: [
      { to: '/preferences', icon: 'ti-adjustments-horizontal', label: 'Предпочтения' },
    ],
  },
]

export function Sidebar() {
  const { user, logout } = useAuth()
  const navigate = useNavigate()

  const handleLogout = () => {
    logout()
    navigate('/login')
  }

  const displayName = user?.name ?? 'Гость'
  const displayEmail = user?.email ?? ''
  const initials = displayName.slice(0, 2).toUpperCase()

  return (
    <aside className="sidebar">
      <div className="sidebar-logo">
        RecSystem
        <span>рекомендации контента</span>
      </div>

      {navItems.map(({ section, items }) => (
        <div key={section}>
          <p className="nav-section">{section}</p>
          {items.map(({ to, icon, label }) => (
            <NavLink
              key={to}
              to={to}
              end={to === '/'}
              className={({ isActive }) => `nav-item ${isActive ? 'active' : ''}`}
            >
              <i className={`ti ${icon}`} aria-hidden="true" />
              {label}
            </NavLink>
          ))}
        </div>
      ))}

      <div className="sidebar-bottom">
        <div className="user-row">
          <div className="avatar">{initials}</div>
          <div>
            <div className="user-name">{displayName}</div>
            <div className="user-email">{displayEmail}</div>
          </div>
        </div>
        <button type="button" className="btn-logout" onClick={handleLogout}>
          Выйти
        </button>
      </div>
    </aside>
  )
}
