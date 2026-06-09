import { useEffect, useState } from 'react'
import { useNavigate, useParams } from 'react-router-dom'
import AlertMessage from '../components/AlertMessage.jsx'
import Badge from '../components/Badge.jsx'
import Button from '../components/Button.jsx'
import { getEventById } from '../services/api.js'
import { isAuthenticated } from '../utils/auth.js'
import { formatDuration, formatEventDate, formatEventTime, getEventImage, normalizeCartEvent } from '../utils/events.js'

export default function EventDetail({ onAddToCart }) {
  const { id } = useParams()
  const navigate = useNavigate()
  const [event, setEvent] = useState(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')

  useEffect(() => {
    let ignore = false

    const loadEvent = async () => {
      setLoading(true)
      setError('')
      try {
        const data = await getEventById(id)
        if (!ignore) setEvent(data)
      } catch (err) {
        if (!ignore) setError(err.message)
      } finally {
        if (!ignore) setLoading(false)
      }
    }

    loadEvent()
    return () => {
      ignore = true
    }
  }, [id])

  const checkout = () => {
    if (!isAuthenticated()) {
      navigate('/login', { state: { message: 'Necesitás iniciar sesión para comprar una entrada.' } })
      return
    }

    onAddToCart({ eventId: event.id, event: normalizeCartEvent(event), quantity: 1 })
    navigate('/checkout')
  }

  if (loading) {
    return (
      <div className="app-shell grid place-items-center px-4">
        <div className="glass-card max-w-lg rounded-3xl p-10 text-center"><h1 className="text-4xl font-black">Cargando evento...</h1></div>
      </div>
    )
  }

  if (error || !event) {
    return (
      <div className="app-shell grid place-items-center px-4">
        <div className="glass-card max-w-lg rounded-3xl p-10 text-center"><h1 className="text-4xl font-black">Evento no encontrado</h1><div className="mt-4"><AlertMessage type="error" message={error || 'El evento que buscás no existe o ya no está disponible.'} /></div><Button to="/" className="mt-8">Volver al Home</Button></div>
      </div>
    )
  }

  return (
    <div className="bg-ticket-alt">
      <div className="container-page py-6"><button onClick={() => navigate('/')} className="rounded-2xl border border-ticket-border bg-ticket-card px-4 py-2 font-bold text-gray-200 transition hover:border-ticket-purple2">← Volver</button></div>
      <section className="relative -mt-20 h-[560px] overflow-hidden">
        <img src={getEventImage(event)} alt={event.title} className="h-full w-full object-cover opacity-55" />
        <div className="absolute inset-0 bg-gradient-to-t from-ticket-alt via-ticket-alt/55 to-ticket-bg/40" />
      </section>
      <main className="container-page relative -mt-72">
        <div className="max-w-5xl">
          <div className="flex flex-wrap gap-3"><Badge>{event.category}</Badge><Badge>{event.active ? 'Activo' : 'Inactivo'}</Badge></div>
          <h1 className="mt-5 text-5xl font-black tracking-tight sm:text-7xl">{event.title}</h1>
          <div className="mt-7 grid gap-4 text-gray-200 sm:grid-cols-2 lg:grid-cols-4"><span>📅 {formatEventDate(event.start_date)}</span><span>⏰ {formatEventTime(event.start_date)}</span><span>📍 {event.location}</span><span>👥 {event.available_capacity} disponibles</span></div>
        </div>

        <div className="mt-12 grid gap-8 lg:grid-cols-[1fr_380px]">
          <div className="space-y-8">
            <section className="glass-card rounded-3xl p-6 sm:p-8"><h2 className="text-2xl font-black">Sobre el evento</h2><p className="mt-4 leading-8 text-gray-300">{event.description}</p></section>
            <section className="glass-card rounded-3xl p-6 sm:p-8">
              <h2 className="text-2xl font-black">Información</h2>
              <div className="mt-5 grid gap-4 text-gray-300 sm:grid-cols-2">
                <p>Categoría: <b className="text-white">{event.category}</b></p>
                <p>Duración: <b className="text-white">{formatDuration(event.duration_minutes)}</b></p>
                <p>Capacidad total: <b className="text-white">{event.capacity}</b></p>
                <p>Disponibles: <b className="text-white">{event.available_capacity}</b></p>
              </div>
            </section>
            <section className="glass-card overflow-hidden rounded-3xl"><img src={getEventImage(event)} alt={event.location} className="h-64 w-full object-cover opacity-80" /><div className="p-6 sm:p-8"><h2 className="text-2xl font-black">Venue</h2><p className="mt-3 text-lg text-violet-200">{event.location}</p><p className="mt-1 text-gray-400">{formatEventDate(event.start_date)} · {formatEventTime(event.start_date)}</p></div></section>
          </div>

          <aside className="h-fit rounded-3xl border border-ticket-purple2/35 bg-ticket-card p-6 shadow-glow lg:sticky lg:top-8">
            <h2 className="text-2xl font-black">Comprar entrada</h2>
            <div className="mt-5 rounded-2xl border border-ticket-border bg-ticket-card2 p-4"><div className="flex justify-between"><div><p className="font-black">General</p><p className="text-sm text-gray-400">Acceso general</p></div><p className="text-xl font-black text-violet-200">Disponible</p></div></div>
            <div className="my-6 space-y-3 border-y border-ticket-border py-5"><div className="flex justify-between text-gray-300"><span>General × 1</span><span>Entrada</span></div><div className="flex justify-between text-xl font-black"><span>Disponibles</span><span>{event.available_capacity}</span></div></div>
            <Button onClick={checkout} className="w-full" disabled={event.available_capacity <= 0}>Comprar ahora →</Button><p className="mt-4 text-center text-xs text-gray-500">Compra 100% segura · Tipo General</p>
          </aside>
        </div>
      </main>
    </div>
  )
}
