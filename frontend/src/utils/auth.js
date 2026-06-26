const TOKEN_KEY = 'ticketapp_token'
const USER_KEY = 'ticketapp_user'
const LEGACY_LOGIN_KEY = 'ticketapp_isLoggedIn'

const removeLegacyLocalAuth = () => {
  localStorage.removeItem(TOKEN_KEY)
  localStorage.removeItem(USER_KEY)
  localStorage.removeItem(LEGACY_LOGIN_KEY)
}

export const saveAuthSession = (response) => {
  if (!response?.token) return

  removeLegacyLocalAuth()
  sessionStorage.setItem(TOKEN_KEY, response.token)
  sessionStorage.setItem(USER_KEY, JSON.stringify(response.user ?? null))
  sessionStorage.setItem(LEGACY_LOGIN_KEY, 'true')
}

export const getStoredToken = () => {
  removeLegacyLocalAuth()
  return sessionStorage.getItem(TOKEN_KEY)
}

export const getStoredUser = () => {
  removeLegacyLocalAuth()
  try {
    const value = sessionStorage.getItem(USER_KEY)
    return value ? JSON.parse(value) : null
  } catch {
    return null
  }
}

export const clearAuthSession = () => {
  removeLegacyLocalAuth()
  sessionStorage.removeItem(TOKEN_KEY)
  sessionStorage.removeItem(USER_KEY)
  sessionStorage.removeItem(LEGACY_LOGIN_KEY)
}

export const isAuthenticated = () => Boolean(getStoredToken())
