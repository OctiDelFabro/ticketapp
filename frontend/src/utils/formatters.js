export const normalizePrice = (price) => {
  if (price === null || price === undefined || price === '') return 0
  const numericPrice = Number(price)
  return Number.isFinite(numericPrice) && numericPrice >= 0 ? numericPrice : 0
}

export const formatPrice = (price) => {
  const numericPrice = normalizePrice(price)
  if (numericPrice === 0) return 'Gratis'

  return `$${new Intl.NumberFormat('es-AR', {
    maximumFractionDigits: 0,
  }).format(numericPrice)}`
}
