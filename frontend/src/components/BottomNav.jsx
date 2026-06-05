import { NavLink } from 'react-router-dom'

const tabs = [
  { to: '/', label: 'Home', icon: '🏠' },
  { to: '/evento/arctic-monkeys', label: 'Evento', icon: '🎸' },
  { to: '/login', label: 'Login', icon: '👤' },
  { to: '/checkout', label: 'Checkout', icon: '🛒' },
  { to: '/mis-entradas', label: 'Tickets', icon: '🎟️' },
]

export default function BottomNav() {
  return (
    <nav className="fixed bottom-4 left-1/2 z-50 flex -translate-x-1/2 gap-1 rounded-2xl border border-ticket-border bg-ticket-card/80 p-2 shadow-glow backdrop-blur-xl">
      {tabs.map((tab) => (
        <NavLink key={tab.to} to={tab.to} className={({ isActive }) => `flex min-w-14 flex-col items-center gap-1 rounded-xl px-2 py-2 text-[11px] font-bold transition sm:min-w-20 sm:px-3 ${isActive ? 'bg-ticket-purple text-white' : 'text-gray-400 hover:bg-ticket-card2 hover:text-white'}`}>
          <span className="text-base">{tab.icon}</span><span>{tab.label}</span>
        </NavLink>
      ))}
    </nav>
  )
}
