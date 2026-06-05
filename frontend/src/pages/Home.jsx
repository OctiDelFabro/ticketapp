import { useState } from 'react'
import { useNavigate } from 'react-router-dom'
import Badge from '../components/Badge.jsx'
import Button from '../components/Button.jsx'
import EventCard from '../components/EventCard.jsx'
import { events } from '../data/mockData.js'

const categories = ['Todos', 'Rock', 'Pop', 'Electrónica', 'Teatro', 'Deportes']

export default function Home({ searchQuery = '' }) {
  const [category, setCategory] = useState('Todos')
  const navigate = useNavigate()
  const featured = events[0]
  const normalizedSearch = searchQuery.trim().toLowerCase()
  const filteredByCategory = category === 'Todos' ? events : events.filter((event) => event.category === category)
  const filtered = normalizedSearch
    ? filteredByCategory.filter((event) =>
        [event.title, event.venue, event.city, event.category].some((value) =>
          value.toLowerCase().includes(normalizedSearch),
        ),
      )
    : filteredByCategory

  return (
    <div>
      <main className="container-page py-8">
        <section className="relative overflow-hidden rounded-[2rem] border border-ticket-purple2/30 bg-ticket-card shadow-glow">
          <img src={featured.image} alt={featured.title} className="absolute inset-0 h-full w-full object-cover opacity-45" />
          <div className="absolute inset-0 bg-gradient-to-r from-ticket-bg via-ticket-bg/80 to-transparent" />
          <div className="relative max-w-3xl p-8 sm:p-12 lg:p-16">
            <div className="flex flex-wrap gap-3"><Badge>Destacado</Badge></div>
            <p className="mt-8 text-sm font-black uppercase tracking-[0.35em] text-violet-200">Evento del año</p>
            <h1 className="mt-3 text-5xl font-black tracking-tight sm:text-7xl">{featured.title}</h1>
            <div className="mt-6 grid gap-3 text-gray-200 sm:grid-cols-3"><span>📅 {featured.date}</span><span>⏰ {featured.time}</span><span>📍 {featured.venue}</span></div>
            <div className="mt-8 flex flex-wrap gap-4"><Button to="/evento/arctic-monkeys">Comprar entradas</Button><Button to="/evento/arctic-monkeys" variant="secondary">Más info &gt;</Button></div>
          </div>
        </section>

        <section className="mt-10">
          <div className="flex flex-col justify-between gap-5 sm:flex-row sm:items-end"><div><p className="text-sm font-black uppercase tracking-[0.25em] text-ticket-purple2">Agenda 2026</p><h2 className="mt-2 text-3xl font-black">Próximos eventos</h2></div></div>
          <div className="mt-6 flex gap-3 overflow-x-auto pb-2">
            {categories.map((item) => <button key={item} onClick={() => setCategory(item)} className={`whitespace-nowrap rounded-full border px-5 py-2 text-sm font-black transition ${category === item ? 'border-ticket-purple bg-ticket-purple text-white shadow-glow' : 'border-ticket-border bg-ticket-card text-gray-300 hover:border-ticket-purple2/60'}`}>{item}</button>)}
          </div>
          <div className="mt-8 grid gap-6 md:grid-cols-2 lg:grid-cols-3">
            {filtered.length ? filtered.map((event) => <EventCard key={event.id} event={event} />) : <div className="col-span-full rounded-3xl border border-ticket-border bg-ticket-card p-10 text-center text-gray-400">No encontramos eventos para tu búsqueda.</div>}
          </div>
        </section>
      </main>
    </div>
  )
}
