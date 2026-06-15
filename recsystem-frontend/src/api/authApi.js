import { api } from '../services/api'

export async function login({ email, password }) {
  const { data } = await api.post('/auth/login', { email, password })
  return data
}

export async function register({ email, password, name, surname, role = 'user' }) {
  const { data } = await api.post('/auth/register', {
    email,
    password,
    name,
    surname,
    role,
  })
  return data
}
