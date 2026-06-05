import { Link } from 'react-router-dom'

export default function Logo({ centered = false }) {
  return (
    <Link to="/" className={`flex items-center gap-3 ${centered ? 'justify-center' : ''}`}>
      <span className="grid h-11 w-11 place-items-center rounded-2xl bg-gradient-to-br from-ticket-purple to-ticket-purple2 shadow-glow">
        <svg viewBox="0 0 24 24" aria-hidden="true" className="h-6 w-6 text-white">
          <path
            d="M5 7.5A2.5 2.5 0 0 1 7.5 5h9A2.5 2.5 0 0 1 19 7.5v1.25a2.25 2.25 0 0 0 0 4.5v1.25A2.5 2.5 0 0 1 16.5 17h-9A2.5 2.5 0 0 1 5 14.5v-1.25a2.25 2.25 0 0 0 0-4.5V7.5Z"
            fill="none"
            stroke="currentColor"
            strokeLinejoin="round"
            strokeWidth="1.8"
          />
          <path d="M12.8 7.4 9.9 12h2.4l-1.1 4.5 3.2-5h-2.5l.9-4.1Z" fill="currentColor" />
        </svg>
      </span>
      <span className="text-2xl font-black tracking-tight">
        Ticket<span className="text-ticket-purple2">App</span>
      </span>
    </Link>
  )
}
