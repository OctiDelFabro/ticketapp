import { useState } from 'react'
import AlertMessage from './AlertMessage.jsx'
import Badge from './Badge.jsx'
import Button from './Button.jsx'
import { formatPrice } from '../utils/formatters.js'
import { formatEventDate, formatEventTime, getEventImage } from '../utils/events.js'

const emailPattern = /^[^\s@]+@[^\s@]+\.[^\s@]+$/

export default function TicketCard({ ticket, onCancel, onTransfer, cancelling = false, disabled = false, transferring = false }) {
  const [targetEmail, setTargetEmail] = useState('')
  const [showTransfer, setShowTransfer] = useState(false)
  const [emailError, setEmailError] = useState('')
  const isActive = ticket.status === 'ACTIVE'
  const imageSrc = getEventImage(ticket)
  const ticketPrice = ticket.event_price ?? ticket.event?.price

  const submitTransfer = () => {
    if (disabled) return
    const email = targetEmail.trim()
    if (!emailPattern.test(email)) {
      setEmailError('Ingresá un email destino válido.')
      return
    }
    onTransfer(ticket.id, email)
    setTargetEmail('')
    setEmailError('')
    setShowTransfer(false)
  }

  const toggleTransfer = () => {
    if (disabled) return
    setShowTransfer((value) => !value)
    setEmailError('')
  }

  return (
    <article className="overflow-hidden rounded-3xl border border-ticket-border bg-ticket-card shadow-soft">
      <div className={`flex flex-col border-l-4 ${isActive ? 'border-ticket-purple' : 'border-gray-600'} sm:flex-row`}>
        <img src={imageSrc} alt={ticket.event_title || 'Imagen del evento'} className="h-48 w-full object-cover sm:h-auto sm:w-44" />
        <div className="flex-1 p-5">
          <div className="flex flex-wrap items-start justify-between gap-3"><div><h3 className="text-2xl font-black">{ticket.event_title}</h3><p className="mt-2 text-violet-200">{formatEventDate(ticket.event_start_date)} · {formatEventTime(ticket.event_start_date)}</p><p className="text-gray-400">{ticket.event_location}</p></div><Badge tone={isActive ? 'purple' : 'red'}>{ticket.status}</Badge></div>
          <div className="mt-5 grid gap-2 text-sm text-gray-400 sm:grid-cols-2">
            <p>ID ticket: <span className="font-bold text-white">{ticket.id}</span></p>
            <p>ID evento: <span className="font-bold text-white">{ticket.event_id}</span></p>
            <p>Compra: <span className="font-bold text-white">{formatEventDate(ticket.purchase_date)}</span></p>
            <p>Email: <span className="font-bold text-white">{ticket.user_email}</span></p>
            <p>Tipo: <span className="font-bold text-white">General</span></p>
            <p>Precio: <span className="font-bold text-white">{ticketPrice === undefined || ticketPrice === null ? 'Precio no disponible' : formatPrice(ticketPrice)}</span></p>
          </div>
          {isActive && <div className="mt-5 flex flex-wrap gap-3"><Button onClick={() => onCancel(ticket.id)} variant="danger" type="button" disabled={disabled}>{cancelling ? 'Cancelando entrada...' : '× Cancelar'}</Button><Button onClick={toggleTransfer} variant="secondary" type="button" disabled={disabled}>{transferring ? 'Transfiriendo entrada...' : 'Transferir'}</Button></div>}
          {isActive && showTransfer && (
            <div className="mt-4 space-y-3">
              <label className="block text-sm font-bold text-gray-300" htmlFor={`transfer-${ticket.id}`}>Email destino para transferir la entrada</label>
              <div className="flex flex-col gap-3 sm:flex-row">
                <input id={`transfer-${ticket.id}`} className={`input-dark ${emailError ? 'border-red-500' : ''}`} disabled={disabled} onChange={(event) => { setTargetEmail(event.target.value); setEmailError('') }} placeholder="email@destino.com" type="email" value={targetEmail} />
                <Button onClick={submitTransfer} type="button" disabled={disabled}>Confirmar transferencia</Button>
              </div>
              <AlertMessage type="error" message={emailError} />
            </div>
          )}
        </div>
      </div>
    </article>
  )
}
