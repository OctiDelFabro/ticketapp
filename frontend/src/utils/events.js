import heroImage from '../assets/hero.png'

export const fallbackEventImage = heroImage

const dateFormatter = new Intl.DateTimeFormat('es-AR', {
  day: 'numeric',
  month: 'short',
  year: 'numeric',
})

const timeFormatter = new Intl.DateTimeFormat('es-AR', {
  hour: '2-digit',
  minute: '2-digit',
})

export const parseEventDate = (value) => {
  const date = value ? new Date(value) : null
  return date && !Number.isNaN(date.getTime()) ? date : null
}

export const formatEventDate = (value) => {
  const date = parseEventDate(value)
  return date ? dateFormatter.format(date) : 'Fecha a confirmar'
}

export const formatEventTime = (value) => {
  const date = parseEventDate(value)
  return date ? `${timeFormatter.format(date)} hs` : 'Horario a confirmar'
}

export const formatDuration = (minutes) => {
  if (!minutes) return 'Duración a confirmar'
  const hours = Math.floor(minutes / 60)
  const remainingMinutes = minutes % 60
  if (!hours) return `${remainingMinutes} min`
  return remainingMinutes ? `${hours} h ${remainingMinutes} min` : `${hours} h`
}

export const getEventImage = (event) => event?.image_url || event?.image || fallbackEventImage

export const normalizeCartEvent = (event) => ({
  ...event,
  image: getEventImage(event),
  date: formatEventDate(event?.start_date),
  time: formatEventTime(event?.start_date),
  venue: event?.location || 'Lugar a confirmar',
  city: '',
  price: event?.price ?? 0,
})
