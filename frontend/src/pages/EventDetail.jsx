import { useMemo, useState } from 'react'
import { Link, useNavigate, useParams } from 'react-router-dom'
import Badge from '../components/Badge.jsx'
import Button from '../components/Button.jsx'
import EventCard from '../components/EventCard.jsx'
import { events, formatPrice, getEventById, serviceFee } from '../data/mockData.js'

export default function EventDetail() {
  const { id } = useParams()
  const navigate = useNavigate()
  const [quantity, setQuantity] = useState(1)
  const event = getEventById(id)
  const related = useMemo(() => event ? event.relatedEvents.map(getEventById).filter(Boolean) : [], [event])

  if (!event) {
    return (
      <div className="app-shell grid place-items-center px-4">
        <div className="glass-card max-w-lg rounded-3xl p-10 text-center"><h1 className="text-4xl font-black">Evento no encontrado</h1><p className="mt-4 text-gray-400">El evento que buscás no existe o ya no está disponible.</p><Button to="/" className="mt-8">Volver al Home</Button></div>
      </div>
    )
  }

  const subtotal = event.price * quantity
  const fee = serviceFee(event.price, quantity)

  return (
    <div className="app-shell bg-ticket-alt">
      <div className="container-page py-6"><button onClick={() => navigate(-1)} className="rounded-2xl border border-ticket-border bg-ticket-card px-4 py-2 font-bold text-gray-200 transition hover:border-ticket-purple2">← Volver</button></div>
      <section className="relative -mt-20 h-[560px] overflow-hidden">
        <img src={event.image} alt={event.title} className="h-full w-full object-cover opacity-55" />
        <div className="absolute inset-0 bg-gradient-to-t from-ticket-alt via-ticket-alt/55 to-ticket-bg/40" />
      </section>
      <main className="container-page relative -mt-72">
        <div className="max-w-5xl">
          <div className="flex flex-wrap gap-3">{event.tags.map((tag) => <Badge key={tag}>{tag}</Badge>)}</div>
          <h1 className="mt-5 text-5xl font-black tracking-tight sm:text-7xl">{event.title}</h1>
          <div className="mt-7 grid gap-4 text-gray-200 sm:grid-cols-2 lg:grid-cols-4"><span>📅 {event.date}</span><span>⏰ {event.time} · puertas {event.doorsTime}</span><span>📍 {event.venue}</span><span>👥 {event.attendees} asistentes</span></div>
        </div>

        <div className="mt-12 grid gap-8 lg:grid-cols-[1fr_380px]">
          <div className="space-y-8">
            <section className="glass-card rounded-3xl p-6 sm:p-8"><h2 className="text-2xl font-black">Sobre el evento</h2>{event.description.map((paragraph) => <p key={paragraph} className="mt-4 leading-8 text-gray-300">{paragraph}</p>)}</section>
            <section className="glass-card rounded-3xl p-6 sm:p-8"><h2 className="text-2xl font-black">Lineup</h2><div className="mt-5 space-y-3">{event.lineup.map((item) => <div key={`${item.artist}-${item.time}`} className="flex flex-wrap items-center justify-between gap-3 rounded-2xl border border-ticket-border bg-ticket-card2 p-4"><div><p className="font-black">{item.artist}</p><p className="text-sm text-gray-400">{item.role}</p></div><Badge>{item.time}</Badge></div>)}</div></section>
            <section className="glass-card overflow-hidden rounded-3xl"><img src={event.venueImage} alt={event.venue} className="h-64 w-full object-cover opacity-80" /><div className="p-6 sm:p-8"><h2 className="text-2xl font-black">Venue</h2><p className="mt-3 text-lg text-violet-200">{event.venue}</p><p className="mt-1 text-gray-400">{event.address}</p></div></section>
            <section><h2 className="text-2xl font-black">También te puede interesar</h2><div className="mt-5 grid gap-5 sm:grid-cols-2 xl:grid-cols-4">{related.map((item) => <EventCard key={item.id} event={item} compact />)}</div></section>
          </div>

          <aside className="h-fit rounded-3xl border border-ticket-purple2/35 bg-ticket-card p-6 shadow-glow lg:sticky lg:top-8">
            <h2 className="text-2xl font-black">Comprar entrada</h2>
            <div className="mt-5 rounded-2xl border border-ticket-border bg-ticket-card2 p-4"><div className="flex justify-between"><div><p className="font-black">General</p><p className="text-sm text-gray-400">Acceso campo general</p></div><p className="text-xl font-black text-violet-200">{formatPrice(event.price)}</p></div></div>
            <div className="mt-5 flex items-center justify-between"><span className="font-bold">Cantidad</span><div className="flex items-center gap-3"><button onClick={() => setQuantity(Math.max(1, quantity - 1))} className="h-10 w-10 rounded-full border border-ticket-border bg-ticket-card2 font-black">−</button><span className="w-8 text-center text-xl font-black">{quantity}</span><button onClick={() => setQuantity(Math.min(6, quantity + 1))} className="h-10 w-10 rounded-full border border-ticket-border bg-ticket-card2 font-black">+</button></div></div>
            <div className="my-6 space-y-3 border-y border-ticket-border py-5"><div className="flex justify-between text-gray-300"><span>General × {quantity}</span><span>{formatPrice(subtotal)}</span></div><div className="flex justify-between text-gray-300"><span>Cargo por servicio</span><span>{formatPrice(fee)}</span></div><div className="flex justify-between text-xl font-black"><span>Total</span><span>{formatPrice(subtotal + fee)}</span></div></div>
            <Button to="/checkout" className="w-full">Comprar ahora →</Button><p className="mt-4 text-center text-xs text-gray-500">Compra 100% segura · Reembolso garantizado</p>
          </aside>
        </div>
      </main>
    </div>
  )
}
