import { useNavigate } from 'react-router-dom'
import Badge from './Badge.jsx'
import Button from './Button.jsx'
import { formatEventDate, formatEventTime, getEventImage } from '../utils/events.js'

export default function EventCard({ event, compact = false }) {
  const navigate = useNavigate()
  return (
    <article onClick={() => navigate(`/evento/${event.id}`)} className={`group cursor-pointer overflow-hidden rounded-3xl border border-ticket-border bg-ticket-card shadow-soft transition hover:-translate-y-1 hover:border-ticket-purple2/50 ${compact ? '' : 'hover:shadow-glow'}`}>
      <div className={`relative ${compact ? 'h-32' : 'h-56'} overflow-hidden bg-gradient-to-br from-ticket-purple/30 to-ticket-card2`}>
        <img src={getEventImage(event)} alt={event.title} className="h-full w-full object-cover opacity-80 transition duration-500 group-hover:scale-105" />
        <div className="absolute inset-0 bg-gradient-to-t from-ticket-card via-transparent to-transparent" />
        {event.category && <Badge tone="purple" className="absolute left-4 top-4">{event.category}</Badge>}
      </div>
      <div className={compact ? 'p-4' : 'p-5'}>
        <h3 className={`${compact ? 'text-base' : 'text-xl'} font-black`}>{event.title}</h3>
        <p className="mt-2 text-sm font-bold text-violet-200">{formatEventDate(event.start_date)} · {formatEventTime(event.start_date)}</p>
        <p className="mt-1 text-sm text-gray-400">{event.location}</p>
        <div className="mt-5 flex items-end justify-between gap-3">
          <div><p className="text-xs uppercase tracking-widest text-gray-500">Tipo</p><p className="text-2xl font-black text-white">General</p></div>
          {!compact && <Button onClick={(e) => { e.stopPropagation(); navigate(`/evento/${event.id}`) }} className="px-4 py-2">Comprar</Button>}
        </div>
      </div>
    </article>
  )
}
