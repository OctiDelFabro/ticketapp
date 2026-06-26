import { adminEventCategoryFilters } from '../constants/eventCategories.js'

export const baseCategories = adminEventCategoryFilters
export const money = new Intl.NumberFormat('es-AR', { style: 'currency', currency: 'ARS', maximumFractionDigits: 0 })
export const normalizeEvent = (event = {}) => {
  const sold = Number(event.tickets_sold ?? event.sold_count ?? event.sold ?? event.active_tickets ?? 0)
  const capacity = Number(event.capacity ?? 0)
  const availableCapacity = Number(event.available_capacity ?? Math.max(capacity - sold, 0))
  const price = Number(event.price ?? event.ticket_price ?? 0)
  return { ...event, sold_count: sold, tickets_sold: sold, available_capacity: availableCapacity, capacity, price, occupation: capacity > 0 ? Math.round((sold / capacity) * 100) : 0 }
}
export const normalizeEvents = (events = []) => events.map(normalizeEvent)
export const normalizeStatsEvent = (event = {}) => {
  const capacity = Number(event.capacity ?? 0)
  const activeTickets = Number(event.active_tickets ?? 0)
  const availableCapacity = Number(event.available_capacity ?? Math.max(capacity - activeTickets, 0))
  const occupancy = Number(event.occupancy_rate_percent ?? 0)
  const price = Number(event.price ?? 0)
  const estimatedRevenue = Number(event.estimated_revenue ?? activeTickets * price)

  return {
    ...event,
    id: event.event_id ?? event.id ?? event.ID,
    event_id: event.event_id ?? event.id ?? event.ID,
    capacity,
    price,
    active_tickets: activeTickets,
    cancelled_tickets: Number(event.cancelled_tickets ?? 0),
    total_tickets: Number(event.total_tickets ?? activeTickets),
    available_capacity: availableCapacity,
    occupancy_rate_percent: occupancy,
    estimated_revenue: estimatedRevenue,
    sold_count: activeTickets,
    occupation: occupancy,
  }
}
export const normalizeStatsEvents = (events = []) => events.map(normalizeStatsEvent)
export const getId = (item) => item?.id ?? item?.event_id ?? item?.ID
export const toInputDate = (value) => {
  const date = value ? new Date(value) : null
  return date && !Number.isNaN(date.getTime()) ? date.toISOString().slice(0, 16) : ''
}
