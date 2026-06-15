export function parseLoginError(err) {
  const raw = err?.response?.data ?? err?.message ?? ''
  const text = typeof raw === 'string' ? raw : String(raw)

  if (
    text.includes('Неверный логин или пароль') ||
    text.includes('invalid credentials') ||
    text.includes('wrong arguments') ||
    text.includes('rpc error')
  ) {
    return 'Неверный логин или пароль'
  }

  if (err?.response?.status === 401) {
    return 'Неверный логин или пароль'
  }

  return 'Не удалось войти. Попробуйте позже'
}

export function parseRegisterError(err) {
  const raw = err?.response?.data ?? err?.message ?? ''
  const text = typeof raw === 'string' ? raw : String(raw)

  if (text.includes('уже зарегистрирован') || text.includes('already exists')) {
    return 'Пользователь с таким email уже зарегистрирован'
  }

  if (text.includes('Проверьте правильность') || text.includes('wrong arguments')) {
    return 'Проверьте правильность заполнения полей (имя и фамилия — латиницей, пароль от 6 символов)'
  }

  return 'Не удалось зарегистрироваться. Попробуйте позже'
}
