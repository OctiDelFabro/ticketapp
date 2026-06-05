import Badge from './Badge.jsx'
import Button from './Button.jsx'

export default function TicketCard({ ticket }) {
  return (
    <article className="overflow-hidden rounded-3xl border border-ticket-border bg-ticket-card shadow-soft">
      <div className="flex flex-col border-l-4 border-ticket-purple sm:flex-row">
        <img src={ticket.image} alt={ticket.event} className="h-48 w-full object-cover sm:h-auto sm:w-44" />
        <div className="flex-1 p-5">
          <div className="flex flex-wrap items-start justify-between gap-3"><div><h3 className="text-2xl font-black">{ticket.event}</h3><p className="mt-2 text-violet-200">{ticket.date} · {ticket.time}</p><p className="text-gray-400">{ticket.venue}</p></div><Badge tone="purple">{ticket.type}</Badge></div>
          <div className="mt-5 grid gap-2 text-sm text-gray-400 sm:grid-cols-2"><p>Código: <span className="font-bold text-white">{ticket.code}</span></p><p>Entradas: <span className="font-bold text-white">{ticket.quantity}</span></p></div>
          <div className="mt-5 flex flex-wrap gap-3"><Button variant="danger" type="button">× Cancelar</Button><Button variant="secondary" type="button">Transferir</Button></div>
        </div>
      </div>
    </article>
  )
}
