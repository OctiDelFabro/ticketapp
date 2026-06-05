const CART_KEY = 'ticketapp_cart'

export const getStoredCart = () => {
  try {
    const value = localStorage.getItem(CART_KEY)
    return value ? JSON.parse(value) : null
  } catch {
    return null
  }
}

export const saveStoredCart = (cartItem) => {
  if (!cartItem) {
    localStorage.removeItem(CART_KEY)
    return
  }

  localStorage.setItem(CART_KEY, JSON.stringify(cartItem))
}
