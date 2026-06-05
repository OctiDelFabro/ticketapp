import { useState } from 'react'
import { useNavigate } from 'react-router-dom'
import Button from '../components/Button.jsx'
import TicketCard from '../components/TicketCard.jsx'
import { events } from '../data/mockData.js'

const tickets = {
  upcoming: [
    { event: 'Arctic Monkeys', date: '15 Jul 2026', time: '21:00 hs', venue: 'Estadio Unico La Plata', code: 'TICKETAPP-2026-AM-0047819', quantity: 2, type: 'General', image: events[0].image },
    { event: 'Bad Bunny', date: '12 Ago 2026', time: '21:30 hs', venue: 'Estadio Monumental', code: 'TICKETAPP-2026-BB-1029481', quantity: 1, type: 'General', image: events[4].image },
  ],
  past: [
    { event: 'Bizarrap Music Sessions', date: '20 Mar 2026', time: '20:30 hs', venue: 'Movistar Arena', code: 'TICKETAPP-2026-BZ-3819201', quantity: 2, type: 'General', image: events[1].image },
    { event: 'Hernan Cattaneo', date: '2 Feb 2026', time: '23:00 hs', venue: 'Mandarine Park', code: 'TICKETAPP-2026-HC-1928301', quantity: 4, type: 'General', image: events[3].image },
  ],
}

export default function MyTickets({ isLoggedIn }) {
  const [tab, setTab] = useState('upcoming')
  const navigate = useNavigate()

  if (!isLoggedIn) {
    return (
      <div className="bg-ticket-alt">
        <main className="container-page grid min-h-[70vh] place-items-center py-10">
          <section className="glass-card max-w-xl rounded-3xl p-8 text-center sm:p-10">
            <div className="mx-auto grid h-16 w-16 place-items-center rounded-2xl bg-ticket-purple text-2xl shadow-glow">🎟️</div>
            <h1 className="mt-6 text-3xl font-black">Iniciá sesión para ver tus entradas</h1>
            <p className="mt-3 leading-7 text-gray-400">
              Accedé a tu cuenta para consultar tus compras, entradas próximas e historial.
            </p>
            <Button onClick={() => navigate('/login')} className="mt-6">Iniciar sesión</Button>
          </section>
        </main>
      </div>
    )
  }

  return (
    <div className="bg-ticket-alt">
      <main className="container-page py-8">
        <button onClick={() => navigate('/')} className="mb-6 font-bold text-gray-300 hover:text-white">← Inicio</button>
        <section className="glass-card rounded-3xl p-6 sm:p-8">
          <div className="flex flex-col gap-5 sm:flex-row sm:items-center">
            <div className="grid h-20 w-20 place-items-center rounded-3xl bg-ticket-purple text-2xl font-black shadow-glow">JP</div>
            <div>
              <p className="text-gray-400">Bienvenido de vuelta</p>
              <h1 className="text-3xl font-black">Juan Perez</h1>
              <p className="text-violet-200">juan.perez@email.com</p>
            </div>
          </div>
        </section>
        <div className="mt-8 flex gap-3 rounded-2xl border border-ticket-border bg-ticket-card p-2 sm:w-fit">
          <button onClick={() => setTab('upcoming')} className={`flex-1 rounded-xl px-5 py-3 font-black transition sm:flex-none ${tab === 'upcoming' ? 'bg-ticket-purple text-white shadow-glow' : 'text-gray-400 hover:text-white'}`}>Próximas 2</button>
          <button onClick={() => setTab('past')} className={`flex-1 rounded-xl px-5 py-3 font-black transition sm:flex-none ${tab === 'past' ? 'bg-ticket-purple text-white shadow-glow' : 'text-gray-400 hover:text-white'}`}>Pasadas 2</button>
        </div>
        <section className="mt-6 space-y-5">
          {tickets[tab].map((ticket) => <TicketCard key={ticket.code} ticket={ticket} />)}
        </section>
      </main>
    </div>
  )
}
