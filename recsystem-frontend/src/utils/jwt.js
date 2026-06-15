export function parseJwt(token) {
  try {
    const payload = token.split('.')[1]
    const json = atob(payload.replace(/-/g, '+').replace(/_/g, '/'))
    return JSON.parse(json)
  } catch {
    return null
  }
}

export function getUserFromToken(token) {
  const claims = parseJwt(token)
  if (!claims) return null
  return {
    id: claims.uid,
    email: claims.email ?? '',
    name: claims.email?.split('@')[0] ?? 'Пользователь',
  }
}
