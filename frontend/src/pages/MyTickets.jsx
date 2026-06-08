import { useEffect, useMemo, useState } from 'react'
import { useNavigate } from 'react-router-dom'
import Button from '../components/Button.jsx'
import TicketCard from '../components/TicketCard.jsx'
import { cancelTicket, getMyTickets, transferTicket } from '../services/api.js'
import { getStoredUser } from '../utils/auth.js'

export default function MyTickets({ isLoggedIn }) {
  const [tickets, setTickets] = useState([])
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState('')
  const [actionError, setActionError] = useState('')
  const navigate = useNavigate()
  const user = getStoredUser()

  const loadTickets = async () => {
    if (!isLoggedIn) return
    setLoading(true)
    setError('')
    try {
      const data = await getMyTickets()
      setTickets(data)
    } catch (err) {
      setError(err.message)
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    loadTickets()
  }, [isLoggedIn])

  const initials = useMemo(() => {
    const name = user?.name || user?.email || 'TU'
    return name.split(' ').map((part) => part[0]).join('').slice(0, 2).toUpperCase()
  }, [user])

  const handleCancel = async (ticketId) => {
    setActionError('')
    try {
      await cancelTicket(ticketId)
      await loadTickets()
    } catch (err) {
      setActionError(err.message)
    }
  }

  const handleTransfer = async (ticketId, targetEmail) => {
    setActionError('')
    try {
      await transferTicket(ticketId, targetEmail)
      await loadTickets()
    } catch (err) {
      setActionError(err.message)
    }
  }

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
            <div className="grid h-20 w-20 place-items-center rounded-3xl bg-ticket-purple text-2xl font-black shadow-glow">{initials}</div>
            <div>
              <p className="text-gray-400">Bienvenido de vuelta</p>
              <h1 className="text-3xl font-black">{user?.name ?? 'Usuario TicketApp'}</h1>
              <p className="text-violet-200">{user?.email}</p>
            </div>
          </div>
        </section>
        <div className="mt-8 flex gap-3 rounded-2xl border border-ticket-border bg-ticket-card p-2 sm:w-fit">
          <div className="rounded-xl bg-ticket-purple px-5 py-3 font-black text-white shadow-glow">Mis entradas {tickets.length}</div>
        </div>
        {error && <div className="mt-6 rounded-2xl border border-red-500/40 bg-red-500/10 p-4 text-red-200">{error}</div>}
        {actionError && <div className="mt-6 rounded-2xl border border-red-500/40 bg-red-500/10 p-4 text-red-200">{actionError}</div>}
        <section className="mt-6 space-y-5">
          {loading ? (
            <div className="rounded-3xl border border-ticket-border bg-ticket-card p-10 text-center text-gray-400">Cargando entradas...</div>
          ) : tickets.length ? (
            tickets.map((ticket) => <TicketCard key={ticket.id} ticket={ticket} onCancel={handleCancel} onTransfer={handleTransfer} />)
          ) : (
            <div className="rounded-3xl border border-ticket-border bg-ticket-card p-10 text-center text-gray-400">Todavía no tenés entradas.</div>
          )}
        </section>
      </main>
    </div>
  )
}
