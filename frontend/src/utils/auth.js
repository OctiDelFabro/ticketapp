const TOKEN_KEY = 'ticketapp_token'
const USER_KEY = 'ticketapp_user'
const LEGACY_LOGIN_KEY = 'ticketapp_isLoggedIn'

export const saveAuthSession = (response) => {
  if (!response?.token) return

  localStorage.setItem(TOKEN_KEY, response.token)
  localStorage.setItem(USER_KEY, JSON.stringify(response.user ?? null))
  localStorage.setItem(LEGACY_LOGIN_KEY, 'true')
}

export const getStoredToken = () => localStorage.getItem(TOKEN_KEY)

export const getStoredUser = () => {
  try {
    const value = localStorage.getItem(USER_KEY)
    return value ? JSON.parse(value) : null
  } catch {
    return null
  }
}

export const clearAuthSession = () => {
  localStorage.removeItem(TOKEN_KEY)
  localStorage.removeItem(USER_KEY)
  localStorage.setItem(LEGACY_LOGIN_KEY, 'false')
  sessionStorage.removeItem(TOKEN_KEY)
  sessionStorage.removeItem(USER_KEY)
  sessionStorage.removeItem(LEGACY_LOGIN_KEY)
}

export const isAuthenticated = () => Boolean(getStoredToken())
