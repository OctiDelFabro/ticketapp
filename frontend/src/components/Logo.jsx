import { Link } from 'react-router-dom'

export default function Logo({ centered = false }) {
  return (
    <Link to="/" className={`flex items-center gap-3 ${centered ? 'justify-center' : ''}`}>
      <span className="grid h-11 w-11 place-items-center rounded-2xl bg-gradient-to-br from-ticket-purple to-ticket-purple2 text-xl shadow-glow">⚡</span>
      <span className="text-2xl font-black tracking-tight">Ticket<span className="text-ticket-purple2">ar</span></span>
    </Link>
  )
}
