import { useState } from 'react'
import { Link } from 'react-router-dom'
import Logo from '../components/Logo.jsx'
import TicketCard from '../components/TicketCard.jsx'
import { events } from '../data/mockData.js'

const tickets = {
  upcoming: [
    { event: 'Arctic Monkeys', date: '15 Jul 2026', time: '21:00 hs', venue: 'Estadio Único La Plata', code: 'TICKETAR-2026-AM-0047819', quantity: 2, type: 'General', image: events[0].image },
    { event: 'Bad Bunny', date: '12 Ago 2026', time: '21:30 hs', venue: 'Estadio Monumental', code: 'TICKETAR-2026-BB-1029481', quantity: 1, type: 'VIP', image: events[4].image },
  ],
  past: [
    { event: 'Bizarrap Music Sessions', date: '20 Mar 2026', time: '20:30 hs', venue: 'Movistar Arena', code: 'TICKETAR-2026-BZ-3819201', quantity: 2, type: 'General', image: events[1].image },
    { event: 'Hernán Cattáneo', date: '2 Feb 2026', time: '23:00 hs', venue: 'Mandarine Park', code: 'TICKETAR-2026-HC-1928301', quantity: 4, type: 'VIP', image: events[3].image },
  ],
}

export default function MyTickets() {
  const [tab, setTab] = useState('upcoming')
  return (
    <div className="app-shell bg-ticket-alt">
      <header className="border-b border-ticket-border bg-ticket-bg/70 backdrop-blur-xl"><div className="container-page flex items-center justify-between py-5"><Link to="/" className="font-bold text-gray-300 hover:text-white">← Inicio</Link><Logo /></div></header>
      <main className="container-page py-8">
        <section className="glass-card rounded-3xl p-6 sm:p-8">
          <div className="flex flex-col gap-5 sm:flex-row sm:items-center"><div className="grid h-20 w-20 place-items-center rounded-3xl bg-ticket-purple text-2xl font-black shadow-glow">JP</div><div><p className="text-gray-400">Bienvenido de vuelta</p><h1 className="text-3xl font-black">Juan Pérez</h1><p className="text-violet-200">juan.perez@email.com</p></div></div>
        </section>
        <div className="mt-8 flex gap-3 rounded-2xl border border-ticket-border bg-ticket-card p-2 sm:w-fit">
          <button onClick={() => setTab('upcoming')} className={`flex-1 rounded-xl px-5 py-3 font-black transition sm:flex-none ${tab === 'upcoming' ? 'bg-ticket-purple text-white shadow-glow' : 'text-gray-400 hover:text-white'}`}>Próximas 2</button>
          <button onClick={() => setTab('past')} className={`flex-1 rounded-xl px-5 py-3 font-black transition sm:flex-none ${tab === 'past' ? 'bg-ticket-purple text-white shadow-glow' : 'text-gray-400 hover:text-white'}`}>Pasadas 2</button>
        </div>
        <section className="mt-6 space-y-5">{tickets[tab].map((ticket) => <TicketCard key={ticket.code} ticket={ticket} />)}</section>
      </main>
    </div>
  )
}
