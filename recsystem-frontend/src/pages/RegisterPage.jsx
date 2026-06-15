import { useState } from 'react'
import { Link, useNavigate } from 'react-router-dom'
import { register as registerApi, login as loginApi } from '../api/authApi'
import { useAuth } from '../context/AuthContext'
import { parseLoginError, parseRegisterError } from '../utils/authErrors'

export function RegisterPage() {
  const [form, setForm] = useState({
    name: '',
    surname: '',
    email: '',
    password: '',
  })
  const [error, setError] = useState('')
  const [loading, setLoading] = useState(false)
  const { login } = useAuth()
  const navigate = useNavigate()

  const update = (field) => (e) => setForm({ ...form, [field]: e.target.value })

  const handleSubmit = async (e) => {
    e.preventDefault()
    setError('')
    setLoading(true)
    try {
      await registerApi(form)
      const data = await loginApi({ email: form.email, password: form.password })
      login(data.access_token, {
        email: form.email,
        name: `${form.name} ${form.surname}`.trim(),
      })
      navigate('/')
    } catch (err) {
      setError(parseRegisterError(err) || parseLoginError(err))
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="auth-page">
      <form className="auth-card" onSubmit={handleSubmit}>
        <h1>Регистрация</h1>
        <p className="auth-subtitle">Создайте аккаунт RecSystem</p>

        {error && <div className="auth-error">{String(error)}</div>}

        <div className="field-row">
          <label className="field">
            <span>Имя</span>
            <input value={form.name} onChange={update('name')} required />
          </label>
          <label className="field">
            <span>Фамилия</span>
            <input value={form.surname} onChange={update('surname')} required />
          </label>
        </div>

        <label className="field">
          <span>Email</span>
          <input type="email" value={form.email} onChange={update('email')} required />
        </label>

        <label className="field">
          <span>Пароль</span>
          <input
            type="password"
            value={form.password}
            onChange={update('password')}
            required
            minLength={6}
          />
        </label>

        <button type="submit" className="btn-primary" disabled={loading}>
          {loading ? 'Создание…' : 'Зарегистрироваться'}
        </button>

        <p className="auth-footer">
          Уже есть аккаунт? <Link to="/login">Войти</Link>
        </p>
      </form>
    </div>
  )
}
