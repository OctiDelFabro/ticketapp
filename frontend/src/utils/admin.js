export const ADMIN_ROLE_VALUES = ['admin', 'administrador']
export const CLIENT_ROLE_VALUES = ['client', 'cliente']

export const getUserRole = (user) => String(user?.role ?? user?.Role ?? user?.rol ?? user?.type ?? '').trim()

export const isAdminUser = (user) => ADMIN_ROLE_VALUES.includes(getUserRole(user).toLowerCase())

export const isClientUser = (user) => CLIENT_ROLE_VALUES.includes(getUserRole(user).toLowerCase())

export const getUserInitials = (user) => {
  const source = user?.name || user?.email || 'Admin'
  return source.split(/\s|@/).filter(Boolean).slice(0, 2).map((part) => part[0]?.toUpperCase()).join('') || 'AD'
}
