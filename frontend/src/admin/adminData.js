export const baseCategories = ['Todas', 'Rock', 'Pop', 'Electrónica', 'Teatro', 'Deportes']
export const money = new Intl.NumberFormat('es-AR', { style: 'currency', currency: 'ARS', maximumFractionDigits: 0 })
export const normalizeEvent = (event = {}) => {
  const sold = Number(event.sold_count ?? event.tickets_sold ?? event.sold ?? 0)
  const capacity = Number(event.capacity ?? 0)
  const price = Number(event.price ?? event.ticket_price ?? 0)
  return { ...event, sold_count: sold, capacity, price, occupation: capacity > 0 ? Math.round((sold / capacity) * 100) : 0 }
}
export const normalizeEvents = (events = []) => events.map(normalizeEvent)
export const getId = (item) => item?.id ?? item?.ID
export const toInputDate = (value) => {
  const date = value ? new Date(value) : null
  return date && !Number.isNaN(date.getTime()) ? date.toISOString().slice(0, 16) : ''
}
