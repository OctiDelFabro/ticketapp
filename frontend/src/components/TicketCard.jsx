import { useState } from 'react'
import Badge from './Badge.jsx'
import Button from './Button.jsx'
import { fallbackEventImage, formatEventDate, formatEventTime } from '../utils/events.js'

export default function TicketCard({ ticket, onCancel, onTransfer }) {
  const [targetEmail, setTargetEmail] = useState('')
  const [showTransfer, setShowTransfer] = useState(false)
  const isActive = ticket.status === 'ACTIVE'

  const submitTransfer = () => {
    const email = targetEmail.trim()
    if (!email) return
    onTransfer(ticket.id, email)
    setTargetEmail('')
    setShowTransfer(false)
  }

  return (
    <article className="overflow-hidden rounded-3xl border border-ticket-border bg-ticket-card shadow-soft">
      <div className={`flex flex-col border-l-4 ${isActive ? 'border-ticket-purple' : 'border-gray-600'} sm:flex-row`}>
        <img src={fallbackEventImage} alt={ticket.event_title} className="h-48 w-full object-cover sm:h-auto sm:w-44" />
        <div className="flex-1 p-5">
          <div className="flex flex-wrap items-start justify-between gap-3"><div><h3 className="text-2xl font-black">{ticket.event_title}</h3><p className="mt-2 text-violet-200">{formatEventDate(ticket.event_start_date)} · {formatEventTime(ticket.event_start_date)}</p><p className="text-gray-400">{ticket.event_location}</p></div><Badge tone={isActive ? 'purple' : 'red'}>{ticket.status}</Badge></div>
          <div className="mt-5 grid gap-2 text-sm text-gray-400 sm:grid-cols-2">
            <p>ID ticket: <span className="font-bold text-white">{ticket.id}</span></p>
            <p>ID evento: <span className="font-bold text-white">{ticket.event_id}</span></p>
            <p>Compra: <span className="font-bold text-white">{formatEventDate(ticket.purchase_date)}</span></p>
            <p>Email: <span className="font-bold text-white">{ticket.user_email}</span></p>
            <p>Tipo: <span className="font-bold text-white">General</span></p>
          </div>
          {isActive && <div className="mt-5 flex flex-wrap gap-3"><Button onClick={() => onCancel(ticket.id)} variant="danger" type="button">× Cancelar</Button><Button onClick={() => setShowTransfer((value) => !value)} variant="secondary" type="button">Transferir</Button></div>}
          {isActive && showTransfer && (
            <div className="mt-4 flex flex-col gap-3 sm:flex-row">
              <input className="input-dark" onChange={(event) => setTargetEmail(event.target.value)} placeholder="Email destino" type="email" value={targetEmail} />
              <Button onClick={submitTransfer} type="button">Confirmar</Button>
            </div>
          )}
        </div>
      </div>
    </article>
  )
}
