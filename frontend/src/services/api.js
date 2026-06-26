import { getStoredToken } from '../utils/auth.js'

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080/api'

const jsonHeaders = { 'Content-Type': 'application/json' }

const getErrorMessage = (payload, fallback) => {
  if (!payload) return fallback
  if (typeof payload === 'string') return payload
  return payload.error || payload.message || fallback
}

const normalizeListResponse = (payload, key) => {
  if (Array.isArray(payload)) return payload
  if (Array.isArray(payload?.[key])) return payload[key]
  if (Array.isArray(payload?.data)) return payload.data
  return []
}

const normalizeObjectResponse = (payload, key) => payload?.[key] ?? payload?.data ?? payload ?? {}

export async function request(path, { method = 'GET', body, auth = false, params } = {}) {
  const url = new URL(`${API_BASE_URL}${path}`)

  if (params) {
    Object.entries(params).forEach(([key, value]) => {
      if (value !== undefined && value !== null && String(value).trim() !== '') {
        url.searchParams.set(key, value)
      }
    })
  }

  const headers = { ...jsonHeaders }
  if (auth) {
    const token = getStoredToken()
    if (token) headers.Authorization = `Bearer ${token}`
  }

  const response = await fetch(url, {
    method,
    headers,
    body: body ? JSON.stringify(body) : undefined,
  })

  const contentType = response.headers.get('content-type') || ''
  const payload = contentType.includes('application/json') ? await response.json() : null

  if (!response.ok) {
    const error = new Error(getErrorMessage(payload, 'No pudimos completar la operación.'))
    error.status = response.status
    throw error
  }

  return payload
}

export const register = (payload) => request('/auth/register', { method: 'POST', body: payload })
export const login = (payload) => request('/auth/login', { method: 'POST', body: payload })
export const getEvents = (params) => request('/events', { params })
export const getEventById = (id) => request(`/events/${id}`)
export const purchaseTicket = (eventId, quantity = 1) => request('/tickets/purchase', { method: 'POST', body: { event_id: eventId, quantity }, auth: true })
export const giftTicket = (eventId, targetEmail, message = '') => request('/tickets/gift', { method: 'POST', body: { event_id: eventId, target_email: targetEmail, message }, auth: true })
export const getMyTickets = () => request('/tickets/me', { auth: true })
export const cancelTicket = (ticketId) => request(`/tickets/${ticketId}/cancel`, { method: 'PATCH', auth: true })
export const transferTicket = (ticketId, targetEmail) => request(`/tickets/${ticketId}/transfer`, { method: 'PATCH', body: { target_email: targetEmail }, auth: true })


export const getAdminEvents = async (params) => {
  try {
    return normalizeListResponse(await request('/admin/events', { params, auth: true }), 'events')
  } catch (error) {
    if (error.status !== 404 && !/not found|404/i.test(error.message)) throw error
    return normalizeListResponse(await request('/events', { params }), 'events')
  }
}
export const createAdminEvent = (payload) => request('/admin/events', { method: 'POST', body: payload, auth: true })
export const updateAdminEvent = (id, payload) => request(`/admin/events/${id}`, { method: 'PATCH', body: payload, auth: true })
export const deleteAdminEvent = (id) => request(`/admin/events/${id}`, { method: 'DELETE', auth: true })
export const getAdminStatsSummary = async () => normalizeObjectResponse(await request('/admin/stats/summary', { auth: true }), 'summary')
export const getAdminEventStats = async () => normalizeListResponse(await request('/admin/stats/events', { auth: true }), 'events')
export const getAdminEventReport = async (id) => normalizeObjectResponse(await request(`/admin/events/${id}/report`, { auth: true }), 'report')
