import { useNavigate } from 'react-router-dom'
import { useAuth } from '../context/AuthContext'

export function Topbar() {
  const { user, logout } = useAuth()
  const navigate = useNavigate()

  const handleLogout = () => {
    logout()
    navigate('/login')
  }

  return (
    <div className="topbar">
      {/* <div className="topbar-title">RecSystem</div>
      <div className="topbar-actions">
        {user?.email && <span className="topbar-user">{user.email}</span>}
        <button type="button" className="btn-logout-top" onClick={handleLogout}>
          <i className="ti ti-logout" aria-hidden="true" />
          Выйти
        </button>
      </div> */}
    </div>
  )
}
